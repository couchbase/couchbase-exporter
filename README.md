# couchbase-exporter

## Summary

The Couchbase Prometheus Exporter is an official [Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) which supplies Couchbase Server metrics to Prometheus.
These include, but are not limited to, CPU and Memory usage, Ops per second, read and write rates, and queue sizes.
Metrics can be refined to a per node and also a per bucket basis.

### Technical Overview

We run four main collectors (as is the standard for prometheus exporters. see [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)) for bucketInfo, bucketStats, nodes, and tasks. These are items that collect cluster-wide statistics on demand when requested by something like say for example, Grafana.

What we do not run a collector for is for the specific “per-server” statistics that exist for each Couchbase node. For this we run a pernode_bucketstats program which, by default, every 5 seconds, reports to Prometheus statistics for the current node it is running on only, and does not concern itself with other nodes’ statistics.

Most of the useful statistics will be found in bucketStats, nodes and perNodeBucketStats.

## Limitations

This exporter is supported to Couchbase Enterprise subscribers only in conjunction with the Couchbase Autonomous (Kubernetes) Operator.
It is available as an Open Source project and documented to be useful as either a container or a built-from-source binary.
At the moment, it is not tested or documented for support in other kinds of environments.
Contact Couchbase if you would like this to be made available in the future.

## Usage

Firstly make sure to build the exporter binary using `make`.

### Couchbase Exporter Arguments
| Arg | Description | Default |
| ------- | ------- | ------------|
| `-couchbase-address` | The address where Couchbase Server is running | localhost  |
| `-couchbase-port` | The port where Couchbase Server is running | 8091  |
| `-couchbase-username` | Couchbase Server Username | Administrator |
| `-couchbase-password` | Couchbase Server Password | password |
| `-server-address` | The address to host the server on | 127.0.0.1 |
| `-server-port` | The port to host the server on | 9091 |
| `-per-node-refresh` | How frequently to collect `per_node_bucket_stats` in seconds | 5 |
| `-token` | bearer token that allows access to `/metrics` |
| `-cert`  | certificate file for exporter in order to serve metrics over TLS |
| `-key` | private key file for exporter in order to serve metrics over TLS |
| `-ca`  | PKI certificate authority file |
| `-clientCert` | client certificate file to authenticate this client with couchbase-server |
| `-clientKey`  | client private key file to authenticate this client with couchbase-server |
| `-logLevel` | log level (debug/info/warn/error) | info
| `-logJson` | if set to true, logs will be JSON formatted | true

### Docker

#### Local Setup

If you have Couchbase server running locally, run the command `docker build --tag couchbase-exporter:v1 .` to build the image.

Then run `docker run --name couchbase-exporter -d -p 9091:9091 couchbase-exporter:v1`.
This will run the exporter with the default arguments above.

Arguments can be specified by appending them to the docker run command. Note that specifying arguments is completely optional.
For example, if you are running Couchbase at the address 10.1.2.3 with a port 9000, and a new password,
and your Prometheus server is running at 10.0.0.1 on its default port, run:

`docker run --name couchbase-exporter -d -p 9091:9091 couchbase-exporter:v1 --couchbase-address 10.1.2.3 --couchbase-port 9000 --couchbase-password mrkrabs101 --server-address 10.0.0.1`

To see a complete list of available arguments, run `docker run --rm couchbase-exporter --help`

#### Docker Setup

If running Couchbase as a separate Docker container, use `docker network list` to
get the ID of the `bridge` network where Couchbase Server should be running.

Then use `docker network inspect <network-ID>` to get the address of Couchbase Server.

```
"Containers": {
            "d2e4d563998893ce457bd672f40a60b849533541ccc883052943e44469331f71": {
                "Name": "db",
                "EndpointID": "<endpointID>",
                "MacAddress": "<mac>",
                "IPv4Address": "172.17.0.2/16",
                "IPv6Address": ""
            }
        },
```

In this example, we would include `"--couchbase-address", "172.17.0.2"` in the CMD array when building.

In both cases, you should see the exporter displaying metrics at `localhost:9091/metrics`.

### Standalone
If running on Linux, simply run `couchbase-exporter` with the any of the optional arguments in the main directory.

Or navigate to `bin/darwin` to run on Mac.

## Customizing Metrics

The Couchbase Prometheus Exporter allows for customizations via a config file to change the namespace, subsystem, name, and help text for each metric.
It is also possible to enable or disable whether each metric is exported to Prometheus or not.

### Generating a Config File

To generate a config file, simply run the application with the `--print-config` flag to print the default config to standard console output:

```
`docker run --name couchbase-exporter couchbase/couchbase-exporter:v1 --print-config`
```

An example default config file is provided: [/example/config.json](./example/config.json).

### Using a config file

To use a config file, use the `--config-file` flag when you run the executable:

```
docker run --name couchbase-exporter -d -p 9091:9091 couchbase/couchbase-exporter:v1 --config-file ./sample/config.json` 
```

Alternatively, you can set the `COUCHBASE_CONFIG_FILE` environment variable.
**NOTE:** CLI Parameters take precedent over Env Vars and Configuration file values.

#### Using a Config File with Docker and the Couchbase Autonomous Operator

To use a config file with the Couchbase Autonomous Operator, you must build a custom `couchbase-exporter` Docker image that sets the `COUCHBASE_CONFIG_FILE` environment variable and copies the config file.

The following command generates a custom Docker image that incorporates the example config file [/example/Dockerfile](./example/Dockerfile):

```
make build container container-config
```

Once you've created the image and [made it available](https://kubernetes.io/docs/concepts/containers/images/) in your Kubernetes environment, you'll need to specify the image in [`couchbaseclusters.spec.monitoring.prometheus.image`](https://docs.couchbase.com/operator/current/resource/couchbasecluster.html#couchbaseclusters-spec-monitoring-prometheus-image).
Once you create/modify the Couchbase cluster deployment with the new custom image, the Autonomous Operator will deploy and manage the `couchbase-exporter` sidecar container with the new custom configuration.

## Setting up Monitoring tools

To consume the metrics, we need to set up Prometheus and optionally Grafana to create custom dashboards.

There are a couple of ways to install and configure both.
If you are familiar with Kubernetes, the easiest way by far is to setup and run the [kube-prometheus](https://github.com/coreos/kube-prometheus) package, which bundles
the Prometheus Operator, AlertManager, Grafana, and many others. This is great if you are running the Couchbase Automonous Operator.
Alternatively, you can check out the official [Prometheus Helm Chart](https://github.com/helm/charts/tree/master/stable/prometheus).

Otherwise see below.

### Prometheus

#### Local Setup

If you wish to download and run Prometheus locally, follow this [link](https://prometheus.io/download/) and
use our provided `prometheus.yml` which contains the scrape_config that includes a look-up for `localhost:9091` -
which is where the exporter metrics should be served.

Run prometheus and `http://localhost:9091/metrics` should appear on the list of Prometheus targets at `localhost:9090/targets`
and report as ACTIVE.

#### Docker Setup

If you have Couchbase Exporter running as a Docker container, the provided `prometheus.yml` configuration will need editing
as the `localhost:9091` will be reported as DOWN. As with the [previous Docker Setup section](#docker-setup), use the commands
`docker network list` and `docker network inspect <network-ID>` to find the IPv4 Address of the Couchbase Exporter container and
substitute this value for `localhost` in the list of scrape targets. The port will still, obviously, be `9091`.

For example:
```
 static_configs:
      - targets: ['localhost:9090', '172.17.0.3:9091']
```
Now run

`docker build -f promDocker -t my-prometheus .`

to build the image and

`docker run --name my-prometheus -d -p 9090:9090 my-prometheus`

to run the Prometheus Server with our custom config.
The Targets page should now look something like the following.

![targets](img/targets.png)

#### AlertManager

Alerts can be set up in the provided `alerts_rules.yml` or specify a different YAML file name under the
`rule_files` heading in `prometheus.yml`. The template for an alerting rule is as follows:

```
- alert: alertname
  expr: prometheus-expression
  annotations:
    summary: of the alert
    description: of the alert
  for: how long expr must be true for (optional)
  labels:
    optional labels
```

More information can be found over at the [official documentation](https://prometheus.io/docs/alerting/overview/),
which also provides good examples.

#### Queries

Again, the [official documentation](https://prometheus.io/docs/prometheus/latest/querying/basics/) provides more info and
examples.

The most common query you will use through Prometheus' PromQL involves the Instant vector selector, which is where
we shall specify any label values. For example:

```
couchbase_bucketstat_couch_docs_disk_size{bucket="beer-sample"}
couchbase_bucketstat_cpu_utilization_rate{pod="cb-example-0000", bucket="default"}
couchbase_per_node_bucket_curr_connections{node="172.17.0.2:8091", bucket="gamesim-sample", namespace="monitoring"}
```

Prometheus includes a simple graphing tool to display multiple queries at once, but for full
customisability and control, it is recommended to export Prometheus metrics into Grafana and use that.

### Grafana

The easiest way to setup [Grafana](https://github.com/grafana/grafana) is to run it as a docker container:
`docker run -d -p 3000:3000 grafana/grafana`
This hosts the dashboard on `localhost:3000`. To login the username and password are both "admin".

![grafana dashboard](img/grafana.png)

Example Grafana Dashboards have been included within the repository. To import the dashboard of your choice,
hover over to the plus icon on the left sidebar and select "Import". Then copy and paste the JSON from
one of the two files provided.


## Reporting Bugs and Issues
Please use our official [JIRA board](https://issues.couchbase.com/projects/PE/issues/?filter=allopenissues) to report any bugs and issues.

## License

Copyright 2019 Couchbase Inc.

Licensed under the Apache License, Version 2.0

See [LICENSE](https://github.com/couchbase/couchbase-exporter/blob/master/LICENSE) for further details.
