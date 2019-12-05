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
	"github.com/couchbase/couchbase_exporter/pkg/collectors"
	"github.com/couchbase/couchbase_exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/prometheus/client_golang/prometheus"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	operatorUser = "COUCHBASE_OPERATOR_USER"
	operatorPass = "COUCHBASE_OPERATOR_PASS"
	bearerToken  = "AUTH_BEARER_TOKEN"
)

var log = logf.Log.WithName("metrics")

var (
	couchAddr   = flag.String("couchbase-address", "localhost", "The address where Couchbase Server is running")
	couchPort   = flag.String("couchbase-port", "", "The port where Couchbase Server is running.")
	userFlag    = flag.String("couchbase-username", "Administrator", "Couchbase Server Username")
	passFlag    = flag.String("couchbase-password", "password", "Couchbase Server Password")
	svrAddr     = flag.String("server-address", "127.0.0.1", "The address to host the server on")
	svrPort     = flag.String("server-port", "9091", "The port to host the server on")
	refreshTime = flag.String("per-node-refresh", "5", "How frequently to collect per_node_bucket_stats collector in seconds")

	tokenFlag  = flag.String("token", "", "bearer token that allows access to /metrics")
	cert       = flag.String("cert", "", "certificate file for exporter in order to serve metrics over TLS")
	key        = flag.String("key", "", "private key file for exporter in order to serve metrics over TLS")
	ca         = flag.String("ca", "", "PKI certificate authority file")
	clientCert = flag.String("client-cert", "", "client certificate file to authenticate this client with couchbase-server")
	clientKey  = flag.String("client-key", "", "client private key file to authenticate this client with couchbase-server")
)

func main() {
	logf.SetLogger(zap.Logger())
	log.Info("Starting metrics collection...")

	flag.Parse()

	validateInt(*svrPort, "server-port")
	validateInt(*refreshTime, "per-node-refresh")

	// by default take credentials from flags (standalone)
	username := *userFlag
	password := *passFlag

	// for operator only, override flags
	// get couchbase server credentials
	if os.Getenv(operatorUser) != "" {
		username = os.Getenv(operatorUser)
	}
	if os.Getenv(operatorPass) != "" {
		password = os.Getenv(operatorPass)
	}

	client := createClient(username, password)

	prometheus.MustRegister(collectors.NewBucketInfoCollector(client))
	prometheus.MustRegister(collectors.NewBucketStatsCollector(client))
	prometheus.MustRegister(collectors.NewNodesCollector(client))
	prometheus.MustRegister(collectors.NewTaskCollector(client))

	i, _ := strconv.Atoi(*refreshTime)
	collectors.RunPerNodeBucketStatsCollection(client, i)

	serveMetrics()
}

// serve the actual metrics
func serveMetrics() {
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

	if len(*cert) == 0 && len(*key) == 0 {
		if err := http.ListenAndServe(":"+*svrPort, &handler); err != nil {
			log.Error(err, "failed to start server:")
		}
	} else {
		if len(*cert) != 0 && len(*key) != 0 {
			metricsServer := *svrAddr + ":" + *svrPort
			log.Info("metrics server: " + metricsServer)

			if err := http.ListenAndServeTLS(":"+*svrPort, *cert, *key, &handler); err != nil {
				log.Error(err, "failed to start server:")
			}

			log.Info("server started listening on", "server", metricsServer)
		} else {
			log.Error(fmt.Errorf("please specify both cert and key arguments"), "")
			os.Exit(1)
		}
	}
}

// create and config client connection to Couchbase Server
func createClient(username, password string) util.Client {
	// Default to nil
	var tlsClientConfig tls.Config

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
			fmt.Printf("could not read CA")
			os.Exit(1)
		}
		if ok := tlsClientConfig.RootCAs.AppendCertsFromPEM(caContents); !ok {
			fmt.Printf("failed to append CA")
			os.Exit(1)
		}

		certContents, err := ioutil.ReadFile(*clientCert)
		if err != nil {
			fmt.Printf("could not read client cert")
			os.Exit(1)
		}
		key, err := ioutil.ReadFile(*clientKey)
		if err != nil {
			fmt.Printf("could not read client key")
			os.Exit(1)
		}
		cert, err := tls.X509KeyPair(certContents, key)
		if err != nil {
			fmt.Printf("failed to create X509 KeyPair")
			os.Exit(1)
		}
		tlsClientConfig.Certificates = append(tlsClientConfig.Certificates, cert)
	} else {
		if len(*clientCert) != 0 || len(*clientKey) != 0 {
			fmt.Printf("please specify both clientCert and clientKey")
			os.Exit(1)
		}
	}

	if len(*couchPort) != 0 {
		port = *couchPort
	}

	log.Info("dial CB Server at: " + scheme + "://" + *couchAddr + ":" + port)

	return util.NewClient(scheme+"://"+*couchAddr+":"+port, username, password, &tlsClientConfig)
}

func validateInt(str, param string) {
	if _, err := strconv.Atoi(str); err != nil {
		log.Error(fmt.Errorf("%q is not a valid integer, parameter: %s.\n", str, param), "")
		os.Exit(1)
	}
}
