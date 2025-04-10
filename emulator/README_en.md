# BigQuery Emulator

The BigQuery emulator allows you to run a local approximation of this storage system.
The main instructions from [cloud/README.md](../cloud/README.md) is appliable here as well, 
except for the authorization part.

Similar examples of connection and usage in Go can be found in this folder.

To run: `go run . [--kind=grpc]`

The `docker-compose.yml` file contains an example of running the emulator in a container.

Here are some discovered shortcomings to consider when working with the emulator::
1. When deserializing Arrow schema, you shouldn't write the schema to the buffer - each batch already contains it [./grpc.go#L246](./grpc.go#L246).
    Issue: https://github.com/goccy/bigquery-emulator/issues/398
2. You can't create a dataset in the emulator using `create schema`. 
    Instead, you need to specify the
    dataset when starting the emulator using the -dataset parameter.
    However, if you mount a volume, subsequent runs will fail.
    Workaround: Don't create a volume.
    Issue: https://github.com/goccy/bigquery-emulator/issues/397
3. If you create a field with an array type and try to read it in Arrow format via gRPC API, 
    the emulator will crash.
    Currently, arrays cannot be selected this way.
    Issue: https://github.com/goccy/bigquery-emulator/issues/399