# couchbase_exporter
Couchbase Prometheus Exporter (WIP)

## Summary
We run four main collectors (as is the standard for prometheus exporters. see https://prometheus.io/docs/instrumenting/writing_exporters/) for bucketInfo, bucketStats, nodes, and tasks. These are items that collect cluster-wide statistics on demand when requested by something like say for example, Grafana. 

What we do not run a collector for is for the specific “per-server” statistics that exist for each Couchbase node. For this we run a pernode_bucketstats program which, by default, every 5 seconds, reports to Prometheus statistics for the current node it is running on only, and does not concern itself with other nodes’ statistics. 

Most of the useful statistics will be found in bucketStats, nodes and perNodeBucketStats.

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
