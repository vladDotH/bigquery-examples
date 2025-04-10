# Working with BigQuery in Cloud

## Начало работы
First, Google Cloud will prompt you to enable the service and provide the necessary data:
https://cloud.google.com/bigquery/docs/sandbox?hl=en


## (Optional) Authorization with gcloud CLI
1. Install gcloud CLI: https://cloud.google.com/sdk/docs/install-sdk
2. Initialize gcloud CLI on your local machine: gcloud init
3. (Optional) Set up default authorization on your local machine: `gcloud auth application-default login`. In this case, a token won’t be needed.

The gcloud CLI will also install the BigQuery CLI: 
https://cloud.google.com/bigquery/docs/bq-command-line-tool,
which can be used to interact with data in projects (it also works with the emulator).


## Project
Projects contain datasets (analogous to a schema in database terms), which in turn contain tables.
BigQuery can also work with other entities, which we won’t cover here.
Learn more: https://cloud.google.com/bigquery?#common-uses

You can create a project in the following ways:
1. Via the console: https://console.cloud.google.com/cloud-resource-manager
2. Using the CLI: `gcloud projects create <id>`, where id is a unique string of 6 to 30 characters.

## Authorization
To access Google Cloud APIs in your code, you’ll need a service account and a credentials file for authentication.

You can create a service account in two ways:
1. In the Google Cloud Console
2. Using the CLI

### Creating a Service Account and Obtaining a Key in the Google Cloud Console
1. Go to Service Accounts: https://console.cloud.google.com/iam-admin/serviceaccounts
2. Select your project.
3. Click "+ Create service account".
4. Choose a name and assign BigQuery roles (e.g., `BigQuery Admin` for full access).
5. Under the account list, select "Manage keys" for the created account.
6. Create new key → JSON.
7. A JSON file with the token will be downloaded.

### Creating a Service Account and Obtaining a Key via gcloud CLI
Create a service account:
```bash
gcloud iam service-accounts create <unique-name>
```

View the service account email in the console or via: 
`gcloud iam service-accounts list`

Grant BigQuery access:
```bash
gcloud projects add-iam-policy-binding <project-id> \
    --member=serviceAccount:<service-account-email> \
    --role=roles/bigquery.admin
```

Download the token file:
```bash
gcloud iam service-accounts keys create <output-file>.json --iam-account=<service-account-email>
```

## Querying

### SQL 
BigQuery supports Google SQL queries via the CLI and libraries.

Libraries: https://cloud.google.com/bigquery/docs/reference/libraries

Google SQL reference: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax

An example of SQL queries is in `sql.go`. To run the example:
```bash
go run . --project=<project-id> [--kind=sql] [--credentialsFile=<credentials-file>] 
```

By default `kind=sql`, `credentialsFile=credentials.json`.  The credentialsFile is token from previous step. 

### GRPC
BigQuery API supports a [GRPC service](https://cloud.google.com/bigquery/docs/reference/storage/rpc) 
for streaming reads in Apache Arrow and Avro formats.

API: https://cloud.google.com/bigquery/docs/reference/storage/

Libraries: https://cloud.google.com/bigquery/docs/reference/storage/libraries

An example of gRPC usage is in `grpc.go`. To run the example:
```bash
go run . --project=<project-id> --kind=grpc [--credentialsFile=<credentials-file>] 
```
