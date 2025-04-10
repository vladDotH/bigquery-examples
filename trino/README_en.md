# Trino

The basic Trino configuration follows standard setup: https://trino.io/docs/current/installation/deployment.html#configuring-trino

To connect BigQuery as a data source, place a `<filename>.properties` 
file in /etc/catalog with the following content:
```
connector.name=bigquery
bigquery.project-id=<project-id>
bigquery.credentials-file=<path-to-credentials.json>
```

Instead of a JSON file, you can use a base64-encoded string of the same file.
More details: https://trino.io/docs/current/connector/bigquery.html

This folder contains a minimal example of all required files to run Trino with BigQuery as a data source:

<pre>
├── docker-compose.yml          
├── etc                             // trino volume 
    ├── catalog                     // datasources directory
    │   ├── bigquery.properties     // configuration for bigquery datasource
    │   └── credentials.json        // google cloud service account token
    ├── config.properties           // <a href="https://trino.io/docs/current/installation/deployment.html#config-properties">server configuration</a>
    ├── jvm.config                  // <a href="https://trino.io/docs/current/installation/deployment.html#jvm-config">jvm config for trino</a>
    └── node.properties             // <a href="https://trino.io/docs/current/installation/deployment.html#node-properties">instance trino config</a>
</pre>

The BigQuery emulator is unfortunately 
[not yet supported](https://github.com/trinodb/trino/discussions/25541) 
by trino.

## Running and Usage

1. Run `docker compose up` or a single container with equivalent arguments
2. Connect to the container `docker exec -it trino trino` 
    (first "trino" is my container name, second is the CLI command)
3. List available sources: `show catalogs;`. If we named our source
    file `bigquery.properties` - the catalog will be named accordingly (`bigquery`).
4. View existing datasets: `show schemas from bigquery;`
5. View tables in specified dataset: `show tables from bigquery.<dataset>;`
6. Show colums in table: `show tables from bigquery.<dataset>.<table>;`
7. Other queries work similarly. For example, to select from a table:
    `select * from bigquery.<dataset>.<table>;`

You can also view query history, parameters, performance, execution plans etc. in the web interface: http://localhost:8080/