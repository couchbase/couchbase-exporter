//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package collectors

import (
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

type cbasCollector struct {
	m      MetaCollector
	config *objects.CollectorConfig
}

func NewCbasCollector(client util.CbClient, config *objects.CollectorConfig, labelManager util.CbLabelManager) prometheus.Collector {
	if config == nil {
		config = objects.GetAnalyticsCollectorDefaultConfig()
	}

	return &cbasCollector{
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
			labelManger: labelManager,
		},
		config: config,
	}
}

// Describe all metrics.
func (c *cbasCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}
		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

// Collect all metrics.
func (c *cbasCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()

	log.Info("Collecting cbas metrics...")

	ctx, err := c.m.labelManger.GetBasicMetricContext()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, objects.ClusterLabel)

		log.Error("%s", err)

		return
	}

	cbas, err := c.m.client.Cbas()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, ctx.ClusterName)

		log.Error("failed to scrape cbas stats")

		return
	}

	for _, value := range c.config.Metrics {
		if value.Enabled {
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				last(cbas.Op.Samples[objects.AnalyticsMetricPrefix+value.Name]),
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...,
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, ctx.ClusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), ctx.ClusterName)
}
