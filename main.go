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
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	operatorUser = "COUCHBASE_OPERATOR_USER"
	operatorPass = "COUCHBASE_OPERATOR_PASS"

	envUser = "COUCHBASE_USER"
	envPass = "COUCHBASE_PASS"

	bearerToken = "AUTH_BEARER_TOKEN"
)

var (
	couchAddr = flag.String("couchbase-address", "localhost", "The address where Couchbase Server is running")
	couchPort = flag.String("couchbase-port", "", "The port where Couchbase Server is running.")
	userFlag  = flag.String("couchbase-username", "Administrator", "Couchbase Server Username. Overridden by env-var COUCHBASE_USER if set.")
	passFlag  = flag.String("couchbase-password", "password", "Plaintext Couchbase Server Password. Recommended to pass value via env-ver COUCHBASE_PASS. Overridden by aforementioned env-var.")

	svrAddr     = flag.String("server-address", "", "The address to host the server on, default all interfaces")
	svrPort     = flag.String("server-port", "9091", "The port to host the server on")
	refreshTime = flag.String("per-node-refresh", "5", "How frequently to collect per_node_bucket_stats collector in seconds")

	tokenFlag  = flag.String("token", "", "bearer token that allows access to /metrics")
	cert       = flag.String("cert", "", "certificate file for exporter in order to serve metrics over TLS")
	key        = flag.String("key", "", "private key file for exporter in order to serve metrics over TLS")
	ca         = flag.String("ca", "", "PKI certificate authority file")
	clientCert = flag.String("client-cert", "", "client certificate file to authenticate this client with couchbase-server")
	clientKey  = flag.String("client-key", "", "client private key file to authenticate this client with couchbase-server")
	logLevel   = flag.String("log-level", "info", "log level (debug/info/warn/error)")
	logJson    = flag.Bool("log-json", true, "if set to true, logs will be JSON formatted")

	backOffLimit = flag.String("backofflimit", "5", "number of retries after panicking before exiting")
	panics = 0
	panicLimit = 0
)

func main() {
	flag.Parse()

	if *logLevel != "" {
		log.SetLevel(*logLevel)
	}
	if *logJson {
		log.SetFormat("json")
	}


	log.Info("Starting metrics collection...")

	err := validateInt(*svrPort, "server-port")
	if err != nil {
		log.Error("%s", err)
		writeToTerminationLog(err)
		os.Exit(1)
	}

	err = validateInt(*refreshTime, "per-node-refresh")
	if err != nil {
		log.Error("%s", err)
		writeToTerminationLog(err)
		os.Exit(1)
	}

	// by default take credentials from flags (standalone)
	username := *userFlag
	password := *passFlag

	// override flags if env vars exist.
	// get couchbase server credentials
	if os.Getenv(envUser) != "" {
		username = os.Getenv(envUser)
	}
	if os.Getenv(envPass) != "" {
		password = os.Getenv(envPass)
	}

	// for operator only, override both plain-text CLI flags and other env-vars.
	// get couchbase server credentials
	if os.Getenv(operatorUser) != "" {
		username = os.Getenv(operatorUser)
	}
	if os.Getenv(operatorPass) != "" {
		password = os.Getenv(operatorPass)
	}

	client, err := createClient(username, password)
	if err != nil {
		log.Error("%s", err)
		writeToTerminationLog(err)
		os.Exit(1)
	}

	prometheus.MustRegister(collectors.NewNodesCollector(client))
	prometheus.MustRegister(collectors.NewBucketInfoCollector(client))
	prometheus.MustRegister(collectors.NewBucketStatsCollector(client))
	prometheus.MustRegister(collectors.NewTaskCollector(client))

	prometheus.MustRegister(collectors.NewQueryCollector(client))
	prometheus.MustRegister(collectors.NewIndexCollector(client))
	prometheus.MustRegister(collectors.NewFTSCollector(client))
	prometheus.MustRegister(collectors.NewCbasCollector(client))
	prometheus.MustRegister(collectors.NewEventingCollector(client))

	i, err := strconv.Atoi(*refreshTime)
	if err != nil {
		log.Error("%s", err)
		writeToTerminationLog(err)
		os.Exit(1)
	}
	collectors.RunPerNodeBucketStatsCollection(client, i)

	for {
		serveMetrics()
	}
}

func writeToTerminationLog(mainErr error) {
	if mainErr != nil {
		i, err := strconv.Atoi(*backOffLimit)
		if err != nil {
			log.Error("%s", err)
		}

		if panics <= i {
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
			os.Exit(1)
		}
	}
}

func check(mainErr error) {
	if mainErr != nil {
		writeToTerminationLog(mainErr)

		log.Error("panicking %s", mainErr)
		panic(fmt.Sprintf("%s", mainErr))
	}
}


// serve the actual metrics
func serveMetrics() {
	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered in serveMetrics(): %s", r)
		}
	}()

	token := *tokenFlag
	if os.Getenv(bearerToken) != "" {
		token = os.Getenv(bearerToken)
	}

	handler := util.AuthHandler{
		ServeMux: http.NewServeMux(),
	}
	if len(token) != 0 {
		handler.TokenLocation = token
	}
	handler.ServeMux.Handle("/metrics", promhttp.Handler())

	metricsServer := *svrAddr + ":" + *svrPort
	if len(*cert) == 0 && len(*key) == 0 {
		if err := http.ListenAndServe(metricsServer, &handler); err != nil {
			check(fmt.Errorf("failed to start server: %v", err))
		}
		log.Info("server started listening on", "server", metricsServer)
	} else {
		if len(*cert) != 0 && len(*key) != 0 {
			if err := http.ListenAndServeTLS(metricsServer, *cert, *key, &handler); err != nil {
				log.Error("failed to start server: %v", err)
				check(err)
			}
			log.Info("server started listening on", "server", metricsServer)
		} else {
			err := fmt.Errorf("please specify both cert and key arguments")
			log.Error("%s", err)
			writeToTerminationLog(err)
			os.Exit(1)
		}
	}
}

// create and config client connection to Couchbase Server
func createClient(username, password string) (util.Client, error) {
	// Default to nil
	var tlsClientConfig tls.Config
	var client util.Client

	// Default to insecure
	scheme := "http"
	port := "8091"

	// Update the TLS, scheme and port
	if len(*ca) != 0 && len(*clientCert) != 0 && len(*clientKey) != 0 {
		scheme = "https"
		port = "18091"

		tlsClientConfig = tls.Config{
			RootCAs: x509.NewCertPool(),
		}

		caContents, err := ioutil.ReadFile(*ca)
		if err != nil {
			return client, fmt.Errorf("could not read CA")
		}
		if ok := tlsClientConfig.RootCAs.AppendCertsFromPEM(caContents); !ok {
			return client, fmt.Errorf("failed to append CA")
		}

		certContents, err := ioutil.ReadFile(*clientCert)
		if err != nil {
			return client, fmt.Errorf("could not read client cert %s", err)
		}
		key, err := ioutil.ReadFile(*clientKey)
		if err != nil {
			fmt.Printf("could not read client key")
			os.Exit(1)
		}
		cert, err := tls.X509KeyPair(certContents, key)
		if err != nil {
			return client, fmt.Errorf("failed to create X509 KeyPair")
		}
		tlsClientConfig.Certificates = append(tlsClientConfig.Certificates, cert)
	} else {
		if len(*clientCert) != 0 || len(*clientKey) != 0 {
			fmt.Printf("please specify both clientCert and clientKey")
			return client, fmt.Errorf("please specify both clientCert and clientKey")
		}
	}

	if len(*couchPort) != 0 {
		port = *couchPort
	}

	log.Info("dial CB Server at: " + scheme + "://" + *couchAddr + ":" + port)

	client = util.NewClient(scheme+"://"+*couchAddr+":"+port, username, password, &tlsClientConfig)
	return client, nil
}

func validateInt(str, param string) error {
	if _, err := strconv.Atoi(str); err != nil {
		log.Error("%q is not a valid integer, parameter: %s.", str, param)
		return fmt.Errorf("%q is not a valid integer, parameter: %s.", str, param)
	}
	return nil
}
