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
)

type bucketStatsCollector struct {
	m MetaCollector

	config *objects.CollectorConfig
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

// Describe all metrics.
func (c *bucketStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}
		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

func (c *bucketStatsCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()

	log.Info("Collecting bucketstats metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error("%s", err)

		return
	}

	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error("failed to scrape buckets")

		return
	}

	for _, bucket := range buckets {
		log.Debug("Collecting %s bucket stats metrics...", bucket.Name)
		stats, err := c.m.client.BucketStats(bucket.Name)

		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

			log.Error("failed to scrape bucket stats")

			return
		}

		for key, value := range c.config.Metrics {
			log.Debug("Collecting bucket stats: %s", value.Name)

			if value.Enabled {
				switch key {
				case "AvgBgWaitTime":
					ch <- prometheus.MustNewConstMetric(
						value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
						prometheus.GaugeValue,
						last(stats.Op.Samples[objects.AvgBgWaitTime])/1000000, // this comes as microseconds from cb
						bucket.Name,
						clusterName,
					)
				case "EpCacheMissRate":
					ch <- prometheus.MustNewConstMetric(
						value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
						prometheus.GaugeValue,
						min(last(stats.Op.Samples[value.Name]), 100), // percentage can exceed 100 due to code within CB, so needs limiting to 100
						bucket.Name,
						clusterName,
					)
				default:
					ch <- prometheus.MustNewConstMetric(
						value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
						prometheus.GaugeValue,
						last(stats.Op.Samples[value.Name]),
						bucket.Name,
						clusterName,
					)
				}
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}

func NewBucketStatsCollector(client util.CbClient, config *objects.CollectorConfig) prometheus.Collector {
	if config == nil {
		config = objects.GetBucketStatsCollectorDefaultConfig()
	}
	// nolint: lll
	return &bucketStatsCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(config.Namespace, config.Subsystem, objects.DefaultUptimeMetric),
				objects.DefaultUptimeMetricHelp,
				[]string{objects.ClusterLabel},
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(config.Namespace, config.Subsystem, objects.DefaultScrapeDurationMetric),
				objects.DefaultScrapeDurationMetricHelp,
				[]string{objects.ClusterLabel},
				nil,
			),
		},
		config: config,
	}
}
