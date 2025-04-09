package main

import (
	"flag"
	"log"
)

var kind = flag.String("kind", "sql", "kind of example (sql or grpc)")
var projectID = flag.String("project", "testproject", "id of project in google cloud")
var apiUrl = flag.String("apiUrl", "http://0.0.0.0:9050", "URL of BQ emulator API")
var grpcUrl = flag.String("grpcUrl", "0.0.0.0:9060", "URL of BQ emulator GRPC service")

func main() {
	flag.Parse()

	switch *kind {
	case "grpc":
		grpcExample()
	case "sql":
		sqlExample()
	default:
		log.Fatal("Unknown kind")
	}
}
