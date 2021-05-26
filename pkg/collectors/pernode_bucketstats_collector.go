package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	rebalanceSuccess = "rebalance_success"
	notFound         = "node not found"
	namespace        = "cbpernode_bucketstats"
	subsystem        = ""
)

var (
	ErrNotFound = fmt.Errorf(notFound)
	upVec       = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        objects.DefaultUptimeMetric,
			Help:        objects.DefaultUptimeMetricHelp,
			ConstLabels: nil,
		},
		[]string{objects.NodeLabel})
	scrapeVec = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        objects.DefaultScrapeDurationMetric,
			Help:        objects.DefaultScrapeDurationMetricHelp,
			ConstLabels: nil,
		},
		[]string{objects.NodeLabel})
)

type PrometheusVecSetter interface {
	SetGaugeVec(prometheus.GaugeVec, float64, ...string)
}

type PerNodeBucketStatsCollector struct {
	config         *objects.CollectorConfig
	metrics        map[string]*prometheus.GaugeVec
	client         util.CbClient
	nodeCache      string
	up             *prometheus.GaugeVec
	scrapeDuration *prometheus.GaugeVec
	// This is for TESTING purposes only.
	// By default PerNodeBucketStatsCollector implements and uses itself to
	// fulfill this functionality.
	Setter PrometheusVecSetter
}

func NewPerNodeBucketStatsCollector(client util.CbClient, config *objects.CollectorConfig) PerNodeBucketStatsCollector {
	collector := &PerNodeBucketStatsCollector{
		config:         config,
		metrics:        map[string]*prometheus.GaugeVec{},
		client:         client,
		nodeCache:      "",
		up:             upVec,
		scrapeDuration: scrapeVec,
		Setter:         nil,
	}
	collector.Setter = collector

	return *collector
}

// DoWork Implements Worker interface for CycleController.
func (c *PerNodeBucketStatsCollector) DoWork() {
	c.CollectMetrics()
}

func (c *PerNodeBucketStatsCollector) recordError(currNode string, message string, err error) {
	c.Setter.SetGaugeVec(*c.up, 0, currNode)
	log.Error(err, message)
}

func (c *PerNodeBucketStatsCollector) getNodeCache() error {
	if c.nodeCache == "" {
		currNode, err := getCurrentNode(c.client)
		log.Info("Current node", "current_node", currNode)
		c.nodeCache = currNode

		return err
	}

	return nil
}

func (c *PerNodeBucketStatsCollector) CollectMetrics() {
	start := time.Now()

	log.Info("Begin collection of Node Stats")
	// get current node hostname and cache it as we'll need it later when we re-execute
	err := c.getNodeCache()
	if err != nil {
		c.recordError(c.nodeCache, "Could not retrieve current node information.", err)

		return
	}

	clusterName, err := c.client.ClusterName()
	if err != nil {
		c.recordError(c.nodeCache, "error retrieving clustername", err)

		return
	}

	log.Info("Cluster name", "clusterName", clusterName)

	rebalanced, err := getClusterBalancedStatus(c.client)
	if err != nil {
		c.recordError(c.nodeCache, "Unable to get rebalance status", err)

		return
	}

	if !rebalanced {
		log.Info("Waiting for Rebalance... retrying...")

		return
	}

	buckets, err := c.client.Buckets()
	if err != nil {
		c.recordError(c.nodeCache, "Unable to get buckets", err)

		return
	}

	for _, bucket := range buckets {
		log.Info("Collecting per-node bucket stats", "node", c.nodeCache, "bucket", bucket.Name)

		samples, err := getPerNodeBucketStats(c.client, bucket.Name, c.nodeCache)
		if err != nil {
			c.recordError(c.nodeCache, "Unable to get samples for "+bucket.Name, err)

			return
		}

		for _, value := range c.config.Metrics {
			c.setMetric(value, samples, bucket.Name, c.nodeCache)
		}
	}

	c.Setter.SetGaugeVec(*c.up, 1, c.nodeCache)
	c.Setter.SetGaugeVec(*c.scrapeDuration, time.Since(start).Seconds(), c.nodeCache)
	log.Info("Per node bucket stats is complete.", "duration", time.Since(start).String())
}

func (c *PerNodeBucketStatsCollector) setMetric(metric objects.MetricInfo, samples map[string]interface{}, bucketName string, clusterName string) {
	if !metric.Enabled {
		return
	}

	if mt, ok := c.metrics[metric.Name]; ok {
		c.Setter.SetGaugeVec(*mt, last(strToFloatArr(fmt.Sprint(samples[metric.Name]))), bucketName, c.nodeCache, clusterName)
	} else {
		mt := metric.GetPrometheusGaugeVec(c.config.Namespace, c.config.Subsystem)
		c.metrics[metric.Name] = mt
		stats := strToFloatArr(fmt.Sprint(samples[metric.Name]))
		if len(stats) > 0 {
			c.Setter.SetGaugeVec(*mt, last(stats), bucketName, c.nodeCache, clusterName)
		}
	}
}

func getCurrentNode(c util.CbClient) (string, error) {
	nodes, err := c.Nodes()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve nodes: %w", err)
	}

	for _, node := range nodes.Nodes {
		if node.ThisNode { // "ThisNode" is a boolean value indicating that it is the current node
			return node.Hostname, nil // hostname seems to work? just don't use for single node setups
		}
	}

	return "", ErrNotFound
}

func getClusterBalancedStatus(c util.CbClient) (bool, error) {
	node, err := c.Nodes()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve nodes, %w", err)
	}

	return node.Counters[rebalanceSuccess] > 0 || (node.Balanced && node.RebalanceStatus == "none"), nil
}

func (c *PerNodeBucketStatsCollector) SetGaugeVec(vec prometheus.GaugeVec, stat float64, labelValues ...string) {
	vec.WithLabelValues(labelValues...).Set(stat)
}

func strToFloatArr(floatsStr string) []float64 {
	floatsStrArr := strings.Split(floatsStr, " ")

	var floatsArr []float64

	for _, f := range floatsStrArr {
		parse := f

		if strings.HasPrefix(parse, "[") {
			parse = strings.Replace(parse, "[", "", 1)
		}

		if strings.HasSuffix(parse, "]") {
			parse = strings.Replace(parse, "]", "", 1)
		}
		// if the key is omitted from the results (Which we know happens depending on version of CBS), this could be <nil>.
		if strings.Contains(parse, "<nil>") {
			parse = "0.0"
		}

		i, err := strconv.ParseFloat(parse, 64)
		if err == nil {
			floatsArr = append(floatsArr, i)
		} else {
			log.Error(err, "Error parsing %v", f)
		}
	}

	return floatsArr
}

func getPerNodeBucketStats(client util.CbClient, bucketName, nodeName string) (map[string]interface{}, error) {
	url, err := getSpecificNodeBucketStatsURL(client, bucketName, nodeName)

	if err != nil {
		log.Error(err, "unable to GET PerNodeBucketStats")

		return nil, err
	}

	var bucketStats objects.PerNodeBucketStats
	err = client.Get(url, &bucketStats)

	if err != nil {
		log.Error(err, "unable to GET PerNodeBucketStats")

		return nil, fmt.Errorf("unable to GET PerNodeBucketStats: %w", err)
	}

	return bucketStats.Op.Samples, nil
}

func getSpecificNodeBucketStatsURL(client util.CbClient, bucket, node string) (string, error) {
	servers, err := client.Servers(bucket)
	if err != nil {
		log.Error(err, "unable to retrieve Servers %s")

		return "", fmt.Errorf("uanble to retrieve Servers: %w", err)
	}

	correctURI := ""

	for _, server := range servers.Servers {
		if server.Hostname == node {
			correctURI = server.Stats["uri"]

			break
		}
	}

	return correctURI, nil
}
