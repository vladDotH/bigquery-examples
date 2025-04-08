package main

import (
	"flag"
	"log"
)

var kind = flag.String("kind", "sql", "kind of example (sql or grpc)")
var projectID = flag.String("project", "", "id of project in google cloud")
var credentialsFile = flag.String("credentialsFile", "credentials.json", "file with gcloud token of service account")

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
