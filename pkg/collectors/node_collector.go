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
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	healthyState         = "healthy"
	uptime               = "uptime"
	clusterMembership    = "clusterMembership"
	memoryTotal          = "memoryTotal"
	memoryFree           = "memoryFree"
	mcdMemoryAllocated   = "mcdMemoryAllocated"
	mcdMemoryReserved    = "mcdMemoryReserved"
	interestingStats     = "interestingStats"
	systemStats          = "systemStats"
	interestingStatsTrim = "interestingstats_"
	systemStatsTrim      = "systemstats_"
)

type nodesCollector struct {
	m MetaCollector

	config *objects.CollectorConfig
}

func NewNodesCollector(client util.CbClient, config *objects.CollectorConfig, labelManager util.CbLabelManager) prometheus.Collector {
	if config == nil {
		config = objects.GetNodeCollectorDefaultConfig()
	}

	// nolint: lll
	return &nodesCollector{
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
func (c *nodesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}

		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}

	return 0.0
}

func ifActive(s string) float64 {
	if s == "active" {
		return 1.0
	}

	return 0.0
}

func contains(haystack []string, needle string) bool {
	contained := false

	for _, i := range haystack {
		if i == needle {
			contained = true
			break
		}
	}

	return contained
}

// Collect all metrics.
func (c *nodesCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()

	log.Info("Collecting nodes metrics...")

	ctx, err := c.m.labelManger.GetBasicMetricContext()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, objects.ClusterLabel)

		log.Error("%s", err)

		return
	}

	nodes, err := c.m.client.Nodes()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, ctx.ClusterName)

		log.Error("failed to scrape nodes")

		return
	}

	for key, value := range c.config.Metrics {
		if contains(nodeSpecificStats, key) || strings.HasPrefix(key, interestingStats) || strings.HasPrefix(key, systemStats) {
			c.addNodeStats(ch, key, value, &nodes)
		} else {
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				nodes.Counters[value.Name],
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, ctx.ClusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), ctx.ClusterName)
}

func getUptimeValue(uptime string, bitSize int) float64 {
	up, err := strconv.ParseFloat(uptime, bitSize)

	if err != nil {
		return 0
	}

	return up
}

// These are the metrics that we collect per node.  Including metrics with "InterestingStats" and "SystemStats" prefixes.  This list allows us to check
// metrics to see if we should collect them per node, or not.
var nodeSpecificStats = []string{healthyState, uptime, clusterMembership, memoryTotal, memoryFree, mcdMemoryAllocated, mcdMemoryReserved}

func (c *nodesCollector) addNodeStats(ch chan<- prometheus.Metric, key string, value objects.MetricInfo, nodes *objects.Nodes) {
	for _, node := range nodes.Nodes {
		ctx, _ := c.m.labelManger.GetMetricContext("", "")
		ctx.NodeHostname = node.Hostname
		log.Debug("Collecting %s-%s node metrics for metric %s", ctx.ClusterName, ctx.NodeHostname, key)

		switch key {
		case healthyState:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				boolToFloat64(node.Status == healthyState),
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case uptime:
			up := getUptimeValue(node.Uptime, 64)
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				up,
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case clusterMembership:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				ifActive(node.ClusterMembership),
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case memoryTotal:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				node.MemoryTotal,
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case memoryFree:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				node.MemoryFree,
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case mcdMemoryAllocated:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				node.McdMemoryAllocated,
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		case mcdMemoryReserved:
			ch <- prometheus.MustNewConstMetric(
				value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.CounterValue,
				node.McdMemoryReserved,
				c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
		default:
			c.handleNonSpecificNodeMetrics(ch, key, value, node, ctx)
		}
	}
}

func (c *nodesCollector) handleNonSpecificNodeMetrics(ch chan<- prometheus.Metric, key string, value objects.MetricInfo, node objects.Node, ctx util.MetricContext) {
	if strings.HasPrefix(key, interestingStats) {
		ch <- prometheus.MustNewConstMetric(
			value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			node.InterestingStats[strings.TrimPrefix(value.Name, interestingStatsTrim)],
			c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
	} else if strings.HasPrefix(key, systemStats) {
		ch <- prometheus.MustNewConstMetric(
			value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			node.SystemStats[strings.TrimPrefix(value.Name, systemStatsTrim)],
			c.m.labelManger.GetLabelValues(value.Labels, ctx)...)
	}
}
