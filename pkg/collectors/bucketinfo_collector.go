//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package collectors

import (
	"sync"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("collectors")

type bucketInfoCollector struct {
	m      MetaCollector
	config *objects.CollectorConfig
}

func (c *bucketInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}
		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

func (c *bucketInfoCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()

	log.Info("Collecting bucketinfo metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error(err, "Error getting Clustername")

		return
	}

	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error(err, "failed to scrape buckets")

		return
	}

	for _, bucket := range buckets {
		for key, value := range c.config.Metrics {
			log.Info("Collecting for metric", "bucket", bucket.Name, "metric", value.Name, "value", bucket.BucketBasicStats[key])

			if value.Enabled {
				ch <- prometheus.MustNewConstMetric(
					value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
					prometheus.GaugeValue,
					bucket.BucketBasicStats[key],
					bucket.Name,
					clusterName,
				)
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}

func NewBucketInfoCollector(client util.CbClient, config *objects.CollectorConfig) prometheus.Collector {
	if config == nil {
		config = objects.GetBucketInfoCollectorDefaultConfig()
	}

	return &bucketInfoCollector{
		m: MetaCollector{
			mutex:  sync.Mutex{},
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
