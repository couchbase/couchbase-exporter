package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
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
		[]string{objects.ClusterLabel})
	scrapeVec = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        objects.DefaultScrapeDurationMetric,
			Help:        objects.DefaultScrapeDurationMetricHelp,
			ConstLabels: nil,
		},
		[]string{objects.ClusterLabel})
)

type PrometheusVecSetter interface {
	SetGaugeVec(prometheus.GaugeVec, float64, ...string)
}

type PerNodeBucketStatsCollector struct {
	config         *objects.CollectorConfig
	metrics        map[string]*prometheus.GaugeVec
	registry       *prometheus.Registry
	client         util.CbClient
	up             *prometheus.GaugeVec
	scrapeDuration *prometheus.GaugeVec
	labelManger    util.CbLabelManager
	// This is for TESTING purposes only.
	// By default PerNodeBucketStatsCollector implements and uses itself to
	// fulfill this functionality.
	Setter PrometheusVecSetter
}

func NewPerNodeBucketStatsCollector(client util.CbClient, config *objects.CollectorConfig, labelManager util.CbLabelManager) PerNodeBucketStatsCollector {
	collector := &PerNodeBucketStatsCollector{
		client:         client,
		metrics:        map[string]*prometheus.GaugeVec{},
		registry:       prometheus.NewRegistry(),
		config:         config,
		up:             upVec,
		scrapeDuration: scrapeVec,
		labelManger:    labelManager,
	}
	collector.Setter = collector

	return *collector
}

func (c *PerNodeBucketStatsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.metrics {
		metric.Collect(ch)
	}
}

func (c *PerNodeBucketStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(prometheus.BuildFQName(c.config.Namespace, c.config.Subsystem, objects.DefaultScrapeDurationMetric),
		objects.DefaultScrapeDurationMetricHelp, []string{objects.ClusterLabel}, nil)
	ch <- prometheus.NewDesc(prometheus.BuildFQName(c.config.Namespace, c.config.Subsystem, objects.DefaultUptimeMetric),
		objects.DefaultUptimeMetricHelp, []string{objects.ClusterLabel}, nil)

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}
		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

// Implements Worker interface for CycleController.
func (c *PerNodeBucketStatsCollector) DoWork() {
	c.CollectMetrics()
}

func (c *PerNodeBucketStatsCollector) CollectMetrics() {
	start := time.Now()

	log.Info("Begin collection of Node Stats")
	// get current node hostname and cache it as we'll need it later when we re-execute
	ctx, err := c.labelManger.GetBasicMetricContext()
	if err != nil {
		c.Setter.SetGaugeVec(*c.up, 0, objects.ClusterLabel)
		log.Error("%s", err)

		return
	}

	log.Info("Cluster name is: %s", ctx.ClusterName)

	rebalanced, err := getClusterBalancedStatus(c.client)
	if err != nil {
		c.Setter.SetGaugeVec(*c.up, 0, ctx.ClusterName)
		log.Error("Unable to get rebalance status %s", err)

		return
	}

	if !rebalanced {
		log.Info("Waiting for Rebalance... retrying...")
		return
	}

	buckets, err := c.client.Buckets()
	if err != nil {
		c.Setter.SetGaugeVec(*c.up, 0, ctx.ClusterName)
		log.Error("Unable to get buckets %s", err)

		return
	}

	for _, bucket := range buckets {
		log.Debug("Collecting per-node bucket stats, node=%s, bucket=%s", ctx.NodeHostname, bucket.Name)

		ctx, _ := c.labelManger.GetMetricContext(bucket.Name, "")
		samples, err := getPerNodeBucketStats(c.client, ctx)

		if err != nil {
			c.Setter.SetGaugeVec(*c.up, 0, ctx.ClusterName)

			return
		}

		for _, value := range c.config.Metrics {
			c.setMetric(value, samples, ctx)
		}
	}

	c.Setter.SetGaugeVec(*c.up, 1, ctx.ClusterName)
	c.Setter.SetGaugeVec(*c.scrapeDuration, time.Since(start).Seconds(), ctx.ClusterName)
	log.Info("Per node bucket stats is complete Duration: %v", time.Since(start))
}

func (c *PerNodeBucketStatsCollector) setMetric(metric objects.MetricInfo, samples map[string]interface{}, ctx util.MetricContext) {
	if !metric.Enabled {
		return
	}

	if mt, ok := c.metrics[metric.Name]; ok {
		c.Setter.SetGaugeVec(*mt, last(strToFloatArr(fmt.Sprint(samples[metric.Name]))), c.labelManger.GetLabelValues(metric.Labels, ctx)...)
	} else {
		mt := metric.GetPrometheusGaugeVec(c.registry, c.config.Namespace, c.config.Subsystem)
		c.metrics[metric.Name] = mt
		stats := strToFloatArr(fmt.Sprint(samples[metric.Name]))
		if len(stats) > 0 {
			c.Setter.SetGaugeVec(*mt, last(stats), c.labelManger.GetLabelValues(metric.Labels, ctx)...)
		}
	}
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
			log.Error("Error parsing %v", f)
		}
	}

	return floatsArr
}

func getPerNodeBucketStats(client util.CbClient, ctx util.MetricContext) (map[string]interface{}, error) {
	url, err := getSpecificNodeBucketStatsURL(client, ctx.BucketName, ctx.NodeHostname)

	if err != nil {
		log.Error("unable to GET PerNodeBucketStats %s", err)
		return nil, err
	}

	var bucketStats objects.PerNodeBucketStats
	err = client.Get(url, &bucketStats)

	if err != nil {
		log.Error("unable to GET PerNodeBucketStats %s", err)
		return nil, err
	}

	return bucketStats.Op.Samples, nil
}

func getSpecificNodeBucketStatsURL(client util.CbClient, bucket, node string) (string, error) {
	servers, err := client.Servers(bucket)
	if err != nil {
		log.Error("unable to retrieve Servers %s", err)
		return "", err
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
