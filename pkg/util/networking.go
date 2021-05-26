//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

//  Portions Copyright (c) 2018 TOTVS Labs

package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/pkg/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// CaError Error appending CA cert.
	CaError         = "failed to append CA certificate"
	certAndKeyError = "please specify both cert and key arguments"
	caAppendError   = "failed to append CA"
	x509Error       = "failed to create X509 KeyPair"
)

var (
	errCertAndKey = fmt.Errorf(certAndKeyError)
	errCaAppend   = fmt.Errorf(caAppendError)
	errX509       = fmt.Errorf(x509Error)
	log           = logf.Log.WithName("util")
)

type CbClient interface {
	URL(string) string
	Get(string, interface{}) error
	Buckets() ([]objects.BucketInfo, error)
	BucketStats(string) (objects.BucketStats, error)
	BucketPerNodeStats(string, string) (objects.BucketStats, error)
	Nodes() (objects.Nodes, error)
	ClusterName() (string, error)
	NodesNodes() (objects.Nodes, error)
	BucketNodes(string) ([]interface{}, error)
	Tasks() ([]objects.Task, error)
	Servers(string) (objects.Servers, error)
	Query() (objects.Query, error)
	Index() (objects.Index, error)
	Fts() (objects.FTS, error)
	Cbas() (objects.Analytics, error)
	Eventing() (objects.Eventing, error)
	QueryNode(string) (objects.Query, error)
	IndexNode(string) (objects.Index, error)
	GetCurrentNode() (string, error)
}

// Client is the couchbase client.
type Client struct {
	domain string
	Client http.Client
}

// NewClient creates a new couchbase client.
func NewClient(exporterConfig *objects.ExporterConfig) (Client, error) {
	var client = Client{
		domain: "",
		Client: http.Client{
			Transport: &AuthTransport{
				Username:  exporterConfig.CouchbaseUser,
				Password:  exporterConfig.CouchbasePassword,
				config:    nil,
				Transport: nil,
			},
		},
	}

	// Default to nil.
	var tlsClientConfig = tls.Config{ //nolint:gosec
		RootCAs: x509.NewCertPool(),
	}

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
		log.Error(errCertAndKey, "please specify both clientCert and clientKey")
		var certError = errCertAndKey

		return client, certError
	}

	couchFullAddress := fmt.Sprintf("%v://%v:%v", scheme, exporterConfig.CouchbaseAddress, exporterConfig.CouchbasePort)
	log.Info("dial CB Server", "couchAddress", couchFullAddress)

	client.domain = couchFullAddress
	client.Client = http.Client{
		Transport: &AuthTransport{
			Username:  exporterConfig.CouchbaseUser,
			Password:  exporterConfig.CouchbasePassword,
			config:    &tlsClientConfig,
			Transport: nil,
		},
	}

	return client, nil
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
		log.Error(errX509, "could not read client key")
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

// configTLS examines the configuration and creates a TLS configuration.
func ConfigClientTLS(cacert, chain, key string) *tls.Config {
	tlsClientConfig := &tls.Config{ //nolint:gosec
		RootCAs: x509.NewCertPool(),
	}

	caCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		log.Error(err, "error reading cert file")
	}

	if ok := tlsClientConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Error(errors.New(CaError), "error appending cert")
	}

	cert, err := tls.LoadX509KeyPair(chain, key)
	if err != nil {
		log.Error(err, "error loading x509 keypair")
	}

	tlsClientConfig.Certificates = []tls.Certificate{cert}

	return tlsClientConfig
}

func (c Client) URL(path string) string {
	return c.domain + "/" + path
}

const statusOk = 200

func (c Client) Get(path string, v interface{}) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, c.URL(path), nil)
	if err != nil {
		newErr := errors.Wrapf(err, "failed to create request for %s", path)

		return newErr
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		newErr := errors.Wrapf(err, "failed to Get %s", path)

		return newErr
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		newErr := errors.Wrapf(err, "failed to read response body from %s", path)

		return newErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != statusOk {
		newErr := errors.Wrapf(err, "%s failed to Get metrics: %s %d", path, string(bts), resp.StatusCode)

		return newErr
	}

	if err := json.Unmarshal(bts, v); err != nil {
		newErr := errors.Wrapf(err, "failed to unmarshall %s output: %s", path, string(bts))

		return newErr
	}

	return nil
}

// AuthTransport is a http.RoundTripper that does the authentication.
type AuthTransport struct {
	Username string
	Password string
	config   *tls.Config

	Transport http.RoundTripper
}

const tlsHandshakeTimeout = 10
const maxIdleConnections = 100
const idleConnectionTimeout = 90

func (t *AuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSClientConfig:       t.config,
		TLSHandshakeTimeout:   tlsHandshakeTimeout * time.Second,
		MaxIdleConns:          maxIdleConnections,
		IdleConnTimeout:       idleConnectionTimeout * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}

// RoundTrip implements the RoundTripper interface.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header, len(req.Header))

	for k, s := range req.Header {
		req2.Header[k] = append([]string(nil), s...)
	}

	req2.SetBasicAuth(t.Username, t.Password)

	tripper, err := t.transport().RoundTrip(req2)
	if err != nil {
		return tripper, fmt.Errorf("error setting round trip %w", err)
	}

	return tripper, nil
}

// Buckets returns the results of /pools/default/buckets.
func (c Client) Buckets() ([]objects.BucketInfo, error) {
	var buckets []objects.BucketInfo
	err := c.Get("pools/default/buckets", &buckets)

	return buckets, errors.Wrap(err, "failed to Get buckets")
}

// BucketStats returns the results of /pools/default/buckets/<bucket_name>/stats.
func (c Client) BucketStats(name string) (objects.BucketStats, error) {
	var stats objects.BucketStats
	err := c.Get(fmt.Sprintf("pools/default/buckets/%s/stats", name), &stats)

	return stats, errors.Wrap(err, "failed to Get bucket stats")
}

func (c Client) BucketPerNodeStats(bucket, node string) (objects.BucketStats, error) {
	var stats objects.BucketStats
	err := c.Get(fmt.Sprintf("pools/default/buckets/%s/nodes/%s/stats", bucket, node), &stats)

	return stats, errors.Wrap(err, "failed to Get bucket stats")
}

// Nodes returns the results of /pools/default/.
func (c Client) Nodes() (objects.Nodes, error) {
	var nodes objects.Nodes
	err := c.Get("pools/default", &nodes)

	return nodes, errors.Wrap(err, "failed to Get nodes")
}

// ClusterName returns the name of the Cluster.
func (c Client) ClusterName() (string, error) {
	var nodes objects.Nodes
	err := c.Get("pools/default", &nodes)

	return nodes.ClusterName, errors.Wrap(err, "failed to retrieve ClusterName")
}

// NodesNodes returns the results of /pools/nodes/.
func (c Client) NodesNodes() (objects.Nodes, error) {
	var nodes objects.Nodes
	err := c.Get("pools/nodes", &nodes)

	return nodes, errors.Wrap(err, "failed to Get nodes")
}

// BucketNodes returns the nodes that this bucket spans.
func (c Client) BucketNodes(bucket string) ([]interface{}, error) {
	var nodes []interface{}
	err := c.Get(fmt.Sprintf("pools/default/buckets/%s/nodes", bucket), nodes)

	return nodes, errors.Wrap(err, "failed to Get nodes")
}

// Tasks returns the results of /pools/default/tasks.
func (c Client) Tasks() ([]objects.Task, error) {
	var tasks []objects.Task
	err := c.Get("pools/default/tasks", &tasks)

	return tasks, errors.Wrap(err, "failed to Get tasks")
}

func (c Client) Servers(bucket string) (objects.Servers, error) {
	var servers objects.Servers
	err := c.Get(fmt.Sprintf("pools/default/buckets/%s/nodes", bucket), &servers)

	return servers, errors.Wrap(err, "failed to Get servers")
}

func (c Client) Query() (objects.Query, error) {
	var query objects.Query
	err := c.Get("pools/default/buckets/@query/stats", &query)

	return query, errors.Wrap(err, "failed to Get query stats")
}

func (c Client) Index() (objects.Index, error) {
	var index objects.Index
	err := c.Get("pools/default/buckets/@index/stats", &index)

	return index, errors.Wrap(err, "failed to Get index stats")
}

func (c Client) Fts() (objects.FTS, error) {
	var fts objects.FTS
	err := c.Get("pools/default/buckets/@fts/stats", &fts)

	return fts, errors.Wrap(err, "failed to Get FTS stats")
}

func (c Client) Cbas() (objects.Analytics, error) {
	var cbas objects.Analytics
	err := c.Get("pools/default/buckets/@cbas/stats", &cbas)

	return cbas, errors.Wrap(err, "failed to Get Analytics stats")
}

func (c Client) Eventing() (objects.Eventing, error) {
	var eventing objects.Eventing
	err := c.Get("pools/default/buckets/@eventing/stats", &eventing)

	return eventing, errors.Wrap(err, "failed to Get eventing stats")
}

func (c Client) QueryNode(node string) (objects.Query, error) {
	var query objects.Query
	err := c.Get(fmt.Sprintf("pools/default/buckets/@query/nodes/%s/stats", node), &query)

	return query, errors.Wrap(err, "failed to Get query stats")
}

func (c Client) IndexNode(node string) (objects.Index, error) {
	var index objects.Index
	err := c.Get("pools/default/buckets/@index/stats", &index)

	return index, errors.Wrap(err, "failed to Get index stats")
}

// potentially deprecated.
func (c Client) GetCurrentNode() (string, error) {
	nodes, err := c.Nodes()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve nodes, %w", err)
	}

	for _, node := range nodes.Nodes {
		if node.ThisNode {
			return node.Hostname, nil // hostname seems to work? just don't use for single node setups
		}
	}

	return "", errors.New("sidecar container cannot find Couchbase Hostname")
}

type AuthHandler struct {
	ServeMux      *http.ServeMux
	TokenLocation string
}

var (
	errorHeader        = fmt.Errorf("error reading headers")
	errorNoRead        = fmt.Errorf("500 Internal Server Error, unable to read bearer token")
	errorTokenInvalid  = fmt.Errorf("401 Unauthorized, bearer token found but incorrect")
	errorFailWrite     = fmt.Errorf("failed to write response body")
	errorTokenNotFound = fmt.Errorf("401 Unauthorized please supply a bearer token")
)

func (authHandler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(authHandler.TokenLocation) == 0 {
		authHandler.ServeMux.ServeHTTP(w, r)

		return
	}

	tokenAccepted := false

	for name := range r.Header {
		if strings.EqualFold(name, "Authorization") {
			if len(r.Header[name]) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				log.Error(errorHeader, "400 bad request")

				return
			}

			token, err := ioutil.ReadFile(authHandler.TokenLocation)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(errorNoRead, "error reading token")

				return
			}

			tokenString := string(token)
			tokenString = strings.TrimSpace(tokenString)

			if !strings.EqualFold(r.Header[name][0], "Bearer "+tokenString) {
				w.WriteHeader(http.StatusUnauthorized)
				log.Error(errorTokenInvalid, "Invalid bearer token")

				return
			}

			tokenAccepted = true
		}
	}

	if !tokenAccepted {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("401 Unauthorized please supply a bearer token"))

		if err != nil {
			log.Error(errorFailWrite, "Failed to write response")

			return
		}

		log.Error(errorTokenNotFound, "Unauthorized")

		return
	}
}
