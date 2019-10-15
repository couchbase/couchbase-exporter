//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/couchbase_exporter/collectors"
	"github.com/couchbase/couchbase_exporter/util"
	"os"
	"strconv"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("metrics")

var (
	couchAddr = flag.String("couchbase-address", "localhost", "The address where Couchbase Server is running")
	couchPort = flag.String("couchbase-port", "8091", "The port where Couchbase Server is running.")
	userFlag  = flag.String("couchbase-username", "Administrator", "Couchbase Server Username")
	passFlag  = flag.String("couchbase-password", "password", "Couchbase Server Password")
	svrAddr   = flag.String("server-address", "127.0.0.1", "The address to host the server on")
	svrPort   = flag.String("server-port", "9091", "The port to host the server on")
	refreshTime = flag.String("per-node-refresh", "5", "How frequently to collect per_node_bucket_stats collector in seconds")
)

func main() {
	logf.SetLogger(zap.Logger())
	log.Info("Starting metrics collection...")

	flag.Parse()

	validateInt(*couchPort, "couchbase-port")
	validateInt(*svrPort, "server-port")
	validateInt(*refreshTime, "per-node-refresh")

	couchbaseServer := "http://" + *couchAddr + ":" + *couchPort

	client := util.NewClient(couchbaseServer, *userFlag, *passFlag)

	prometheus.MustRegister(collectors.NewBucketInfoCollector(client))
	prometheus.MustRegister(collectors.NewBucketStatsCollector(client))
	prometheus.MustRegister(collectors.NewNodesCollector(client))
	prometheus.MustRegister(collectors.NewTaskCollector(client))

	i, _ := strconv.Atoi(*refreshTime);
	collectors.RunPerNodeBucketStatsCollection(client, i)

	metricsServer := *svrAddr + ":" + *svrPort

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Info("server started listening on", "server", metricsServer)
	if err := http.ListenAndServe(":" + *svrPort, nil); err != nil {
		log.Error(err, "failed to start server:")
	}

}

func validateInt(str, param string) {
	if _, err := strconv.Atoi(str); err != nil {
		fmt.Printf("%q is not a valid integer, parameter: %s.\n", str, param)
		os.Exit(1)
	}
}