//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

//  Portions Copyright (c) 2018 TOTVS Labs

package collectors

import (
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	bucketUpVec = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "bucketstats",
			Subsystem:   "",
			Name:        objects.DefaultUptimeMetric,
			Help:        objects.DefaultUptimeMetricHelp,
			ConstLabels: nil,
		},
		[]string{objects.ClusterLabel})
	bucketScrapeVec = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "bucketstats",
			Subsystem:   "",
			Name:        objects.DefaultScrapeDurationMetric,
			Help:        objects.DefaultScrapeDurationMetricHelp,
			ConstLabels: nil,
		},
		[]string{objects.ClusterLabel})
)

type BucketStatsCollector struct {
	config         *objects.CollectorConfig
	metrics        map[string]*prometheus.GaugeVec
	registry       *prometheus.Registry
	client         util.CbClient
	up             *prometheus.GaugeVec
	scrapeDuration *prometheus.GaugeVec
	labelManger    util.CbLabelManager
	// This is for TESTING purposes only.
	// By default bucketStatsCollector implements and uses itself to
	// fulfill this functionality.
	Setter PrometheusVecSetter
}

func (c *BucketStatsCollector) SetGaugeVec(vec prometheus.GaugeVec, stat float64, labelValues ...string) {
	vec.WithLabelValues(labelValues...).Set(stat)
}

// Implements Worker interface for CycleController.
func (c *BucketStatsCollector) DoWork() {
	c.CollectMetrics()
}

func last(stats []float64) float64 {
	if len(stats) == 0 {
		return 0
	}

	return stats[len(stats)-1]
}

func min(x, y float64) float64 {
	if x > y {
		return y
	}

	return x
}

func (c *BucketStatsCollector) setMetric(metric objects.MetricInfo, samples map[string][]float64, ctx util.MetricContext) {
	if !metric.Enabled {
		return
	}

	promMetric, ok := c.metrics[metric.Name]
	if !ok {
		promMetric = metric.GetPrometheusGaugeVec(c.registry, c.config.Namespace, c.config.Subsystem)
		c.metrics[metric.Name] = promMetric
	}

	switch metric.Name {
	case "avg_bg_wait_time":
		// comes across as microseconds.  Convert
		c.Setter.SetGaugeVec(*promMetric, last(samples[metric.Name])/1000000, c.labelManger.GetLabelValues(metric.Labels, ctx)...)
	case "ep_cache_miss_rate":
		c.Setter.SetGaugeVec(*promMetric, min(last(samples[metric.Name]), 100), c.labelManger.GetLabelValues(metric.Labels, ctx)...)
	default:
		c.Setter.SetGaugeVec(*promMetric, last(samples[metric.Name]), c.labelManger.GetLabelValues(metric.Labels, ctx)...)
	}
}

func (c *BucketStatsCollector) CollectMetrics() {
	start := time.Now()

	log.Info("Begin collecting bucketstats metrics...")

	ctx, err := c.labelManger.GetBasicMetricContext()
	if err != nil {
		c.Setter.SetGaugeVec(*c.up, 0, objects.ClusterLabel)
		log.Error("%s", err)

		return
	}

	buckets, err := c.client.Buckets()
	if err != nil {
		c.Setter.SetGaugeVec(*c.up, 0, objects.ClusterLabel)
		log.Error("failed to scrape buckets")

		return
	}

	for _, bucket := range buckets {
		log.Debug("Collecting %s bucket stats metrics...", bucket.Name)

		ctx, _ := c.labelManger.GetMetricContext(bucket.Name, "")

		stats, err := c.client.BucketStats(bucket.Name)

		if err != nil {
			c.Setter.SetGaugeVec(*c.up, 0, objects.ClusterLabel)
			log.Error("failed to scrape bucket stats")

			return
		}

		for _, value := range c.config.Metrics {
			log.Debug("Collecting bucket stats: %s", value.Name)

			if value.Enabled {
				c.setMetric(value, stats.Op.Samples, ctx)
			}
		}
	}

	c.Setter.SetGaugeVec(*c.up, 1, objects.ClusterLabel)
	c.Setter.SetGaugeVec(*c.scrapeDuration, time.Since(start).Seconds(), ctx.ClusterName)
	log.Info("Bucket stats is complete Duration: %v", time.Since(start))
}

func (c *BucketStatsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.metrics {
		metric.Collect(ch)
	}
}

func (c *BucketStatsCollector) Describe(ch chan<- *prometheus.Desc) {
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

func NewBucketStatsCollector(client util.CbClient, config *objects.CollectorConfig, labelManager util.CbLabelManager) BucketStatsCollector {
	if config == nil {
		config = objects.GetBucketStatsCollectorDefaultConfig()
	}
	// nolint: lll
	collector := &BucketStatsCollector{
		client:         client,
		labelManger:    labelManager,
		up:             bucketUpVec,
		scrapeDuration: bucketScrapeVec,
		registry:       prometheus.NewRegistry(),
		config:         config,
		metrics:        map[string]*prometheus.GaugeVec{},
	}

	collector.Setter = collector

	return *collector
}
