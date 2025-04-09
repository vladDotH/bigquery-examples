package main

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func sqlExample() {
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

	log.Printf("Sample data created successfully, selecting...\n")

	query := client.Query(`select * from testdataset.testtable where num > 1;`)

	it, err := query.Read(ctx)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// Read sample data as map
	log.Printf("Total rows: %v, data:\n", it.TotalRows)
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating through results: %v", err)
		}
		log.Printf("\t%v\n", row)
	}

	// Cleanup test table
	clearSampleData(ctx, client)
}

func createSampleData(ctx context.Context, client *bigquery.Client) {
	query := client.Query(`
		create schema if not exists testdataset;

		create table if not exists testdataset.testtable(
			num int64,
			str string,
			arr Array<int64>,
			obj Struct<x int64, y string>
		);

		insert into testdataset.testtable(num, str, arr, obj)
		values 
			(1, '1 str', [1,2,3], Struct(100, '100str')), 
			(2, '2nd str', [], Struct<int64, string>(200, '200str')), 
			(33, '3rd string', [0,0,0], Struct(300 as x, '300str' as y));
	`)

	job, err := query.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	status, err := job.Wait(ctx)

	if err != nil {
		log.Fatalf("Failed to wait query: %v", err)
	}
	if err = status.Err(); err != nil {
		log.Fatalf("Errors during query execution: %v", err)
	}
}

func clearSampleData(ctx context.Context, client *bigquery.Client) {
	query := client.Query(`drop table testdataset.testtable; drop schema testdataset;`)

	job, err := query.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		log.Fatalf("Failed to wait query: %v", err)
	}
	if err = status.Err(); err != nil {
		log.Fatalf("Errors during query execution: %v", err)
	}
}
