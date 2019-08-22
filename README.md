# couchbase_exporter
Couchbase Prometheus Exporter (WIP)

## Usage

### Docker
Run `docker build --tag couchbase_exporter` if running Couchbase locally.

Edit the Dockerfile CMD if needed to pass arguments for your configuration.

If running Couchbase as a seperate Docker container, use `docker network inspect <network>` to get the address of that Couchbase Server.  

### Standalone
Run `couchbase_exporter` with the any of the following optional arguments.

| Arg | Default | Description |
| ------- | ------- | ------------|
| `couchbase_address` | localhost | The address where Couchbase Server is running |
| `couchPort` | 8091 | The port where Couchbase Server is running
| `couchbase_username` | Administrator | Couchbase Server Username |
| `couchbase_password` | password | Couchbase Server Password |
| `server_address` | 127.0.0.1 | The address to host the server on |
| `server_port` | 9091 | The port to host the server on |
| `per_node_refresh` | 5 | How frequently to collect `per_node_bucket_stats` collector in seconds |
