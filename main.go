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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/config"
	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	certAndKeyError = "please specify both cert and key arguments"
	caAppendError   = "failed to append CA"
	x509Error       = "failed to create X509 KeyPair"
)

var (
	exporterConfig = new(objects.ExporterConfig)
	couchAddr      *string
	couchPort      *string
	userFlag       *string
	passFlag       *string
	svrAddr        *string
	svrPort        *string
	refreshTime    *string
	tokenFlag      *string
	cert           *string
	key            *string
	ca             *string
	clientCert     *string
	clientKey      *string
	logLevel       *string
	logJSON        *bool
	backOffLimit   *string
	configFile     *string
	defaultConfig  *bool
	panics         = 0
	errCertAndKey  = fmt.Errorf(certAndKeyError)
	errCaAppend    = fmt.Errorf(caAppendError)
	errX509        = fmt.Errorf(x509Error)
)

func init() {
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
	logLevel = flag.String("log-level", "", "log level (debug/info/warn/error)")
	logJSON = flag.Bool("log-json", true, "if set to true, logs will be JSON formatted")

	backOffLimit = flag.String("backofflimit", "", "number of retries after panicking before exiting")
	configFile = flag.String("config", "", "The location of the PE configuration. Overridden by env-var COUCHBASE_CONFIG_FILE if set.")
	defaultConfig = flag.Bool("print-config", false, "Outputs the config file with CLI and ENV var override to stdout")
}

func main() {
	flag.Parse()

	// Load config from file, or load up defaults.
	exporterConfig, err := config.New(*configFile)
	if err != nil {
		log.Error("Error loading config file.")
		os.Exit(1)
	}

	// Get Logging settings and initialize log level.
	exporterConfig.SetOrDefaultLogJSON(*logJSON)
	exporterConfig.SetOrDefaultLogLevel(*logLevel)

	if exporterConfig.LogLevel != "" {
		log.SetLevel(exporterConfig.LogLevel)
	}

	if exporterConfig.LogJSON {
		log.SetFormat("json")
	}

	// Override defaults with values from CLI.
	exporterConfig.SetOrDefaultCouchAddress(*couchAddr)
	exporterConfig.SetOrDefaultCouchPort(*couchPort)
	exporterConfig.SetOrDefaultCouchUser(*userFlag)
	exporterConfig.SetOrDefaultCouchPassword(*passFlag)
	exporterConfig.SetOrDefaultServerAddress(*svrAddr)
	exporterConfig.SetOrDefaultServerPort(*svrPort)
	exporterConfig.SetOrDefaultRefreshRate(*refreshTime)
	exporterConfig.SetOrDefaultBackoffLimit(*backOffLimit)
	exporterConfig.SetOrDefaultToken(*tokenFlag)
	exporterConfig.SetOrDefaultCa(*ca)
	exporterConfig.SetOrDefaultCertificate(*cert)
	exporterConfig.SetOrDefaultKey(*key)
	exporterConfig.SetOrDefaultClientCertificate(*clientCert)
	exporterConfig.SetOrDefaultClientKey(*clientKey)

	// This is if we want to dump the config to stdout to generate a configuration file.
	if *defaultConfig {
		c, err := json.MarshalIndent(exporterConfig, "", "    ")
		if err != nil {
			log.Error("Error generating Json config file.  Exiting")
			writeToTerminationLog(err)
			os.Exit(1)
		}

		os.Stdout.WriteString(string(c))
		os.Exit(0)
	}

	log.Info("Couchbase Address:  %s:%v", exporterConfig.CouchbaseAddress, exporterConfig.CouchbasePort)

	log.Info("Starting metrics collection...")

	client, err := createClient(exporterConfig)
	if err != nil {
		log.Error("%s", err)
		writeToTerminationLog(err)
		os.Exit(1)
	}

	log.Info("Registering Collectors...")

	prometheus.MustRegister(collectors.NewNodesCollector(client, exporterConfig.Collectors.Node))
	prometheus.MustRegister(collectors.NewBucketInfoCollector(client, exporterConfig.Collectors.BucketInfo))
	prometheus.MustRegister(collectors.NewBucketStatsCollector(client, exporterConfig.Collectors.BucketStats))
	prometheus.MustRegister(collectors.NewTaskCollector(client, exporterConfig.Collectors.Task))

	prometheus.MustRegister(collectors.NewQueryCollector(client, exporterConfig.Collectors.Query))
	prometheus.MustRegister(collectors.NewIndexCollector(client, exporterConfig.Collectors.Index))
	prometheus.MustRegister(collectors.NewFTSCollector(client, exporterConfig.Collectors.Search))
	prometheus.MustRegister(collectors.NewCbasCollector(client, exporterConfig.Collectors.Analytics))
	prometheus.MustRegister(collectors.NewEventingCollector(client, exporterConfig.Collectors.Eventing))

	// moved to a goroutine to improve startup time.
	// Create my cycle controller with refreshrate (seconds) in milliseconds
	cycle := util.NewCycleController(exporterConfig.RefreshRate * 1000)
	perNodeBucketStatCollector := collectors.NewPerNodeBucketStatsCollector(client, exporterConfig.Collectors.PerNodeBucketStats)

	cycle.Subscribe(&perNodeBucketStatCollector)
	cycle.Start()

	log.Info("Serving Metrics")

	for {
		serveMetrics(exporterConfig)
	}
}

func writeToTerminationLog(mainErr error) {
	if mainErr != nil {
		if panics <= exporterConfig.BackoffLimit {
			f, err := os.OpenFile("/dev/termination-log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString(fmt.Sprintf("%s\n", mainErr)); err != nil {
				panic(err)
			}
			panics++
		} else {
			log.Error("backoffLimit reached")
			defer os.Exit(1)
		}
	}
}

// serve the actual metrics.
func serveMetrics(exporterConfig *objects.ExporterConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered in serveMetrics(): %s", r)
		}
	}()

	handler := util.AuthHandler{
		ServeMux: http.NewServeMux(),
	}

	if len(exporterConfig.Token) != 0 {
		handler.TokenLocation = exporterConfig.Token
	}

	handler.ServeMux.Handle("/metrics", promhttp.Handler())

	metricsServer := fmt.Sprintf("%v:%v", exporterConfig.ServerAddress, exporterConfig.ServerPort)
	log.Info("starting server on %s", metricsServer)

	util.Serve(metricsServer, handler, exporterConfig.Certificate, exporterConfig.Key)
}

func setTLSClientConfig(exporterConfig objects.ExporterConfig, tlsConfig *tls.Config) error {
	caContents, err := ioutil.ReadFile(exporterConfig.Ca)
	if err != nil {
		return fmt.Errorf("could not read CA: %w", err)
	}

	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(caContents); !ok {
		return errCaAppend
	}

	certContents, err := ioutil.ReadFile(exporterConfig.ClientCertificate)
	if err != nil {
		return fmt.Errorf("could not read client cert %w", err)
	}

	key, err := ioutil.ReadFile(exporterConfig.ClientKey)
	if err != nil {
		log.Error("could not read client key")
		os.Exit(1)
	}

	cert, err := tls.X509KeyPair(certContents, key)
	if err != nil {
		var keypairError = errX509
		return fmt.Errorf("%w", keypairError)
	}

	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)

	return nil
}

// create and config client connection to Couchbase Server.
func createClient(exporterConfig *objects.ExporterConfig) (util.Client, error) {
	// Default to nil.
	var tlsClientConfig = tls.Config{
		RootCAs: x509.NewCertPool(),
	}

	var client util.Client

	// Default to insecure.
	scheme := "http"

	// Update the TLS, scheme and port.
	if len(exporterConfig.Ca) != 0 && len(exporterConfig.ClientCertificate) != 0 && len(exporterConfig.ClientKey) != 0 {
		scheme = "https"
		exporterConfig.CouchbasePort = 18091

		err := setTLSClientConfig(*exporterConfig, &tlsClientConfig)
		if err != nil {
			return client, err
		}
	} else if len(exporterConfig.ClientCertificate) != 0 || len(exporterConfig.ClientKey) != 0 {
		log.Error("please specify both clientCert and clientKey")
		var certError = errCertAndKey
		return client, certError
	}

	couchFullAddress := fmt.Sprintf("%v://%v:%v", scheme, exporterConfig.CouchbaseAddress, exporterConfig.CouchbasePort)
	log.Info("dial CB Server at: %v", couchFullAddress)

	client = util.NewClient(couchFullAddress, exporterConfig.CouchbaseUser, exporterConfig.CouchbasePassword, &tlsClientConfig)

	return client, nil
}
