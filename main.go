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
	"os"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/config"
	logging "github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	klog "k8s.io/klog/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	millisecondsInASecond = 1000
)

var (
	couchAddr     *string
	couchPort     *string
	userFlag      *string
	passFlag      *string
	svrAddr       *string
	svrPort       *string
	refreshTime   *string
	tokenFlag     *string
	cert          *string
	key           *string
	ca            *string
	clientCert    *string
	clientKey     *string
	backOffLimit  *string
	configFile    *string
	defaultConfig *bool
	printFile     *string
	log                           = logf.Log.WithName("main")
	logOptions    logging.Options = logging.Options{Options: zap.Options{
		Development:          false,
		Encoder:              nil,
		EncoderConfigOptions: []zap.EncoderConfigOption{},
		NewEncoder:           nil,
		DestWriter:           nil,
		DestWritter:          nil,
		Level:                nil,
		StacktraceLevel:      nil,
		ZapOpts:              nil,
	}}
)

func setupFlags() {
	couchAddr = flag.String("couchbase-address", "", "The address where Couchbase Server is running")
	couchPort = flag.String("couchbase-port", "", "The port where Couchbase Server is running.")
	userFlag = flag.String("couchbase-username", "", "Couchbase Server Username. Overridden by env-var COUCHBASE_USER if set.")
	passFlag = flag.String("couchbase-password", "", "Plaintext Couchbase Server Password. Recommended to pass value via env-ver COUCHBASE_PASS. Overridden by aforementioned env-var.")

	svrAddr = flag.String("server-address", "", "The address to host the server on, default all interfaces")
	svrPort = flag.String("server-port", "", "The port to host the server on")
	refreshTime = flag.String("per-node-refresh", "", "How frequently to collect per_node_bucket_stats collector in seconds")

	tokenFlag = flag.String("token", "", "bearer token that allows access to /metrics")
	cert = flag.String("cert", "", "certificate file for exporter in order to serve metrics over TLS")
	key = flag.String("key", "", "private key file for exporter in order to serve metrics over TLS")
	ca = flag.String("ca", "", "PKI certificate authority file")
	clientCert = flag.String("client-cert", "", "client certificate file to authenticate this client with couchbase-server")
	clientKey = flag.String("client-key", "", "client private key file to authenticate this client with couchbase-server")

	backOffLimit = flag.String("backofflimit", "", "number of retries after panicking before exiting")
	configFile = flag.String("config", "", "The location of the PE configuration. Overridden by env-var COUCHBASE_CONFIG_FILE if set.")
	defaultConfig = flag.Bool("print-config", false, "Outputs the config file with CLI and ENV var override to specified output file")
	printFile = flag.String("output", "./config.json", "Where to output config file.  defaults to './config.json'")
}

func main() {
	logOptions.AddFlagSet(flag.CommandLine)
	logger := logging.New(&logOptions)

	setupFlags()
	flag.Parse()
	logf.SetLogger(logger)
	klog.SetLogger(logger)
	// Load config from file, or load up defaults.
	exporterConfig, err := config.New(*configFile)
	if err != nil {
		log.Error(err, "Error loading config file.")
		os.Exit(1)
	}

	// Override defaults with values from CLI.
	exporterConfig.SetValues(*couchAddr, *couchPort, *userFlag, *passFlag, *svrAddr,
		*svrPort, *refreshTime, *backOffLimit, *tokenFlag, *ca,
		*cert, *key, *clientCert, *clientKey)

	// This is if we want to dump the config to stdout to generate a configuration file.
	if *defaultConfig {
		exporterConfig.PrintJSON(*printFile)
		os.Exit(0)
	}

	log.Info("Couchbase Address", "address", exporterConfig.CouchbaseAddress, "port", exporterConfig.CouchbasePort)

	client, err := util.NewClient(exporterConfig)
	if err != nil {
		log.Error(err, "Error retrieving client.")
		os.Exit(1)
	}

	config.RegisterCollectors(exporterConfig, client)
	// Create my cycle controller with refreshrate (seconds) in milliseconds
	cycle := util.NewCycleController(exporterConfig.RefreshRate * millisecondsInASecond)
	perNodeBucketStatCollector := collectors.NewPerNodeBucketStatsCollector(client, exporterConfig.Collectors.PerNodeBucketStats)

	cycle.Subscribe(&perNodeBucketStatCollector)
	cycle.Start()

	log.Info("Serving Metrics")
	util.Serve(exporterConfig)
}


