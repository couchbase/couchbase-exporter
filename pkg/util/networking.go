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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/pkg/errors"
)

const (
	CaError string = "failed to append CA certificate"
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
	GetCurrentNode() (objects.Node, error)
	IndexStats() (map[string]map[string]interface{}, error)
}

// Client is the couchbase client.
type Client struct {
	port   int
	domain string
	Client http.Client
}

// NewClient creates a new couchbase client.
func NewClient(domain string, port int, user, password string, config *tls.Config) Client {
	var client = Client{
		domain: domain,
		port:   port,
		Client: http.Client{
			Transport: &AuthTransport{
				Username: user,
				Password: password,
				config:   config,
			},
		},
	}

	return client
}

// configTLS examines the configuration and creates a TLS configuration.
func ConfigClientTLS(cacert, chain, key string) *tls.Config {
	tlsClientConfig := &tls.Config{
		RootCAs: x509.NewCertPool(),
	}

	caCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		log.Fatal(err)
	}

	if ok := tlsClientConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Fatal(errors.New(CaError))
	}

	cert, err := tls.LoadX509KeyPair(chain, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsClientConfig.Certificates = []tls.Certificate{cert}

	return tlsClientConfig
}

func (c Client) URL(path string) string {
	return fmt.Sprintf("%s:%d/%s", c.domain, c.port, path)
}

func (c Client) IndexerURL(path string) string {
	var url string

	switch c.port {
	case 18091:
		url = fmt.Sprintf("%s:%d/%s", c.domain, 19102, path)
	case 8091:
		url = fmt.Sprintf("%s:%d/%s", c.domain, 9102, path)
	default:
		url = fmt.Sprintf("%s:%d/%s", c.domain, 9102, path)
	}

	return url
}

func (c Client) IndexAPIGet(path string, v interface{}) error {
	resp, err := c.Client.Get(c.IndexerURL(path))
	if err != nil {
		return errors.Wrapf(err, "failed to Get %s", path)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response body from %s", path)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("failed to Get 200 response status: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(bts, v); err != nil {
		return errors.Wrapf(err, "failed to unmarshall %s output: %s", path, string(bts))
	}

	return nil
}

func (c Client) Get(path string, v interface{}) error {
	resp, err := c.Client.Get(c.URL(path))
	if err != nil {
		return errors.Wrapf(err, "failed to Get %s", path)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response body from %s", path)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("failed to Get 200 response status: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(bts, v); err != nil {
		return errors.Wrapf(err, "failed to unmarshall %s output: %s", path, string(bts))
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

func (t *AuthTransport) transport() http.RoundTripper {
	if nil != t.Transport {
		return t.Transport
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSClientConfig:       t.config,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
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

	return t.transport().RoundTrip(req2)
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

//
func (c Client) IndexNode(node string) (objects.Index, error) {
	var index objects.Index
	err := c.Get("pools/default/buckets/@index/stats", &index)

	return index, errors.Wrap(err, "failed to Get index stats")
}

func (c Client) IndexStats() (map[string]map[string]interface{}, error) {
	var data map[string]map[string]interface{}
	err := c.IndexAPIGet("/api/v1/stats", &data)

	if err != nil {
		return data, errors.Wrap(err, "Failed to get Indexer stats")
	}

	return data, nil
}

// potentially deprecated.
func (c Client) GetCurrentNode() (objects.Node, error) {
	nodes, err := c.Nodes()

	var retNode objects.Node

	if err != nil {
		return retNode, fmt.Errorf("unable to retrieve nodes, %w", err)
	}

	for _, node := range nodes.Nodes {
		if node.ThisNode {
			retNode = node
			return retNode, nil // hostname seems to work? just don't use for single node setups
		}
	}

	return retNode, errors.New("sidecar container cannot find Couchbase Hostname")
}

type AuthHandler struct {
	ServeMux      *http.ServeMux
	TokenLocation string
}

func (authHandler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(authHandler.TokenLocation) != 0 {
		tokenAccepted := false

		for name := range r.Header {
			if strings.EqualFold(name, "Authorization") {
				if len(r.Header[name]) != 1 {
					w.WriteHeader(http.StatusBadRequest)
					log.Println("400 bad request")

					return
				}

				token, err := ioutil.ReadFile(authHandler.TokenLocation)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("500 Internal Server Error, unable to read bearer token")

					return
				}

				tokenString := string(token)
				tokenString = strings.TrimSpace(tokenString)

				if !strings.EqualFold(r.Header[name][0], "Bearer "+tokenString) {
					w.WriteHeader(http.StatusUnauthorized)
					log.Println("401 Unauthorized, bearer token found but incorrect")

					return
				}

				tokenAccepted = true
			}
		}

		if !tokenAccepted {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("401 Unauthorized please supply a bearer token"))

			if err != nil {
				log.Println("failed to write response body")
				return
			}

			log.Println("401 Unauthorized please supply a bearer token")

			return
		}
	}

	authHandler.ServeMux.ServeHTTP(w, r)
}
