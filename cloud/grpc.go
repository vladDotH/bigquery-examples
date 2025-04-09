package main

// Full read storage grpc API example: https://cloud.google.com/bigquery/docs/reference/storage/libraries#use

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	bqStorage "cloud.google.com/go/bigquery/storage/apiv1"
	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/ipc"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var rpcOpts = gax.WithGRPCOptions(
	grpc.MaxCallRecvMsgSize(1024 * 1024 * 129),
)

func grpcExample() {
	ctx := context.Background()

	// Establish BQ api connection
	client, err := bigquery.NewClient(ctx, *projectID, option.WithCredentialsFile(*credentialsFile))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connection to poject %v established!\n", client.Project())

	defer client.Close()

	// Create dataset and table, and fill it
	createSampleData(ctx, client)
	// Scheduling cleanup
	defer clearSampleData(ctx, client)

	log.Printf("Sample data created successfully, starting storage grpc read...\n")

	// Create GRPC BQ read-client
	bqReadClient, err := bqStorage.NewBigQueryReadClient(ctx, option.WithCredentialsFile(*credentialsFile))

	if err != nil {
		log.Fatalf("NewBigQueryStorageClient: %v", err)
	}
	defer bqReadClient.Close()

	log.Printf("Grpc connection established...\n")

	// Path to table (projects/<project-id>/datasets/<dataset-name>/tables/<table-name>)
	readTable := "projects/" + *projectID + "/datasets/testdataset/tables/testtable"

	// Read options (select specific field, and sql-like where clause)
	tableReadOptions := &storagepb.ReadSession_TableReadOptions{
		SelectedFields: []string{"num", "str", "arr", "obj"}, RowRestriction: "num > 1",
	}

	// Read session options (passing apache arrow fomat)
	createReadSessionRequest := &storagepb.CreateReadSessionRequest{
		Parent: "projects/" + *projectID,
		ReadSession: &storagepb.ReadSession{
			Table:       readTable,
			DataFormat:  storagepb.DataFormat_ARROW,
			ReadOptions: tableReadOptions,
		},
		MaxStreamCount: 1,
	}

	// Create the session from the request.
	session, err := bqReadClient.CreateReadSession(ctx, createReadSessionRequest, rpcOpts)

	if err != nil {
		log.Fatalf("CreateReadSession: %v", err)
	}

	if len(session.GetStreams()) == 0 {
		log.Fatalf("no streams in session.  if this was a small query result, consider writing to output to a named table.")
	}

	log.Printf("Read session created: %v\n", session.Name)

	readStream := session.GetStreams()[0].Name

	log.Printf("Session established, got read steam: %v", readStream)

	ch := make(chan *storagepb.ReadRowsResponse)

	var wg sync.WaitGroup

	// Start the reading in one goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := processStream(ctx, bqReadClient, readStream, ch); err != nil {
			log.Fatalf("processStream failure: %v", err)
		}
		close(ch)
	}()

	// Start Arrow processing and decoding in another goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := processArrow(ctx, session.GetArrowSchema().GetSerializedSchema(), ch)

		if err != nil {
			log.Fatalf("error processing: %v", err)
		}
	}()

	// Wait until both the reading and decoding goroutines complete.
	wg.Wait()
}

func processStream(ctx context.Context, client *bqStorage.BigQueryReadClient, st string, ch chan<- *storagepb.ReadRowsResponse) error {
	var offset int64

	// Streams may be long-running.  Rather than using a global retry for the
	// stream, implement a retry that resets once progress is made.
	retryLimit := 3
	retries := 0
	for {
		// Send the initiating request to start streaming row blocks.
		rowStream, err := client.ReadRows(ctx, &storagepb.ReadRowsRequest{
			ReadStream: st,
			Offset:     offset,
		}, rpcOpts)
		if err != nil {
			return fmt.Errorf("couldn't invoke ReadRows: %w", err)
		}

		// Process the streamed responses.
		for {
			r, err := rowStream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				// If there is an error, check whether it is a retryable
				// error with a retry delay and sleep instead of increasing
				// retries count.
				var retryDelayDuration time.Duration
				if errorStatus, ok := status.FromError(err); ok && errorStatus.Code() == codes.ResourceExhausted {
					for _, detail := range errorStatus.Details() {
						retryInfo, ok := detail.(*errdetails.RetryInfo)
						if !ok {
							continue
						}
						retryDelay := retryInfo.GetRetryDelay()
						retryDelayDuration = time.Duration(retryDelay.Seconds)*time.Second + time.Duration(retryDelay.Nanos)*time.Nanosecond
						break
					}
				}
				if retryDelayDuration != 0 {
					log.Printf("processStream failed with a retryable error, retrying in %v", retryDelayDuration)
					time.Sleep(retryDelayDuration)
				} else {
					retries++
					if retries >= retryLimit {
						return fmt.Errorf("processStream retries exhausted: %w", err)
					}
				}
				// break the inner loop, and try to recover by starting a new streaming
				// ReadRows call at the last known good offset.
				break
			} else {
				// Reset retries after a successful response.
				retries = 0
			}

			rc := r.GetRowCount()
			if rc > 0 {
				// Bookmark our progress in case of retries and send the rowblock on the channel.
				offset = offset + rc
				// We're making progress, reset retries.
				retries = 0
				ch <- r
			}
		}
	}
}

// processArrow receives row blocks from a channel, and uses the provided Arrow
// schema to decode the blocks into individual row messages for printing.  Will
// continue to run until the channel is closed or the provided context is
// cancelled.
func processArrow(ctx context.Context, schema []byte, ch <-chan *storagepb.ReadRowsResponse) error {
	mem := memory.NewGoAllocator()
	buf := bytes.NewBuffer(schema)
	r, err := ipc.NewReader(buf, ipc.WithAllocator(mem))
	if err != nil {
		return err
	}
	aschema := r.Schema()

	log.Print("Got arrow data schema: \n", aschema, "\n")

	for {
		select {
		case <-ctx.Done():
			// Context was cancelled.  Stop.
			return ctx.Err()
		case rows, ok := <-ch:
			if !ok {
				// Channel closed, no further arrow messages.  Stop.
				return nil
			}
			undecoded := rows.GetArrowRecordBatch().GetSerializedRecordBatch()

			if len(undecoded) > 0 {
				var buf bytes.Buffer
				buf.Write(schema)
				buf.Write(undecoded)

				r, err = ipc.NewReader(&buf, ipc.WithAllocator(mem), ipc.WithSchema(aschema))
				if err != nil {
					return err
				}
				for r.Next() {
					record := r.Record()

					err = printRecordBatch(record)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}

// printRecordBatch prints the arrow record batch
func printRecordBatch(record arrow.Record) error {
	out, err := record.MarshalJSON()
	if err != nil {
		return err
	}
	list := []map[string]interface{}{}
	err = json.Unmarshal(out, &list)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}

	log.Printf("Got %v records:", len(list))

	for _, rec := range list {
		log.Printf("\t%v", rec)
	}

	return nil
}
