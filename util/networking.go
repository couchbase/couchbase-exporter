package util

import (
	"encoding/json"
	"fmt"
	"main.go/objects"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// Client is the couchbase client
type Client struct {
	domain string
	Client http.Client
}

// New creates a new couchbase client
func NewClient(domain, user, password string) Client {
	var client = Client{
		domain: domain,
		Client: http.Client{
			Transport: &AuthTransport{
				Username: user,
				Password: password,
			},
		},
	}
	return client
}

func (c Client) url(path string) string {
	return c.domain + "/" + path
}

func (c Client) get(path string, v interface{}) error {
	resp, err := c.Client.Get(c.url(path))
	if err != nil {
		return errors.Wrapf(err, "failed to get %s", path)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read response body from %s", path)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get %s metrics: %s %d", path, string(bts), resp.StatusCode)
	}

	if err := json.Unmarshal(bts, v); err != nil {
		return errors.Wrapf(err, "failed to unmarshall %s output: %s", path, string(bts))
	}
	return nil
}

// AuthTransport is a http.RoundTripper that does the authentication
type AuthTransport struct {
	Username string
	Password string

	Transport http.RoundTripper
}

func (t *AuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
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

// Buckets returns the results of /pools/default/buckets
func (c Client) Buckets() ([]objects.BucketInfo, error) {
	var buckets []objects.BucketInfo
	err := c.get("pools/default/buckets", &buckets)
	return buckets, errors.Wrap(err, "failed to get buckets")
}

// BucketStats returns the results of /pools/default/buckets/<bucket_name>/stats
func (c Client) BucketStats(name string) (objects.BucketStats, error) {
	var stats objects.BucketStats
	err := c.get(fmt.Sprintf("pools/default/buckets/%s/stats", name), &stats)
	return stats, errors.Wrap(err, "failed to get bucket stats")
}

func (c Client) BucketPerNodeStats(bucket, node string) (objects.BucketStats, error) {
	var stats objects.BucketStats
	err := c.get(fmt.Sprintf("pools/default/buckets/%s/nodes/%s/stats", bucket, node), &stats)
	return stats, errors.Wrap(err, "failed to get bucket stats")
}

// Nodes returns the results of /pools/default/
func (c Client) Nodes() (objects.Nodes, error) {
	var nodes objects.Nodes
	err := c.get("pools/default", &nodes)
	return nodes, errors.Wrap(err, "failed to get nodes")
}

func (c Client) NodesNodes() (objects.Nodes, error) {
	var nodes objects.Nodes
	err := c.get("pools/nodes", &nodes)
	return nodes, errors.Wrap(err, "failed to get nodes")
}

// BucketNodes returns the nodes that this bucket spans
func (c Client) BucketNodes(bucket string) ([]interface{}, error) {
	var nodes []interface{}
	err := c.get(fmt.Sprintf("pools/default/buckets/%s/nodes", bucket), nodes)
	return nodes, errors.Wrap(err, "failed to get nodes")
}

// Tasks returns the results of /pools/default/tasks
func (c Client) Tasks() ([]objects.Task, error) {
	var tasks []objects.Task
	err := c.get("pools/default/tasks", &tasks)
	return tasks, errors.Wrap(err, "failed to get tasks")
}
