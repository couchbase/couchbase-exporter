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

const (
	taskRebalance                      = "rebalance"
	taskBucketCompaction               = "bucket_compaction"
	taskXdcr                           = "xdcr"
	taskClusterLogCollection           = "clusterLogsCollection"
	metricRebalancePerNode             = "rebalancePerNode"
	metricCompacting                   = "compacting"
	metricXdcrChangesLeft              = "xdcrChangesLeft"
	metricXdcrDocsChecked              = "xdcrDocsChecked"
	metricXdcrDocsWritten              = "xdcrDocsWritten"
	metricXdcrPaused                   = "xdcrPaused"
	metricXdcrErrors                   = "xdcrErrors"
	metricDocsTotal                    = "progressDocsTotal"
	metricDocsTransferred              = "progressDocsTransferred"
	metricDocsActiveVbucketsLeft       = "progressActiveVBucketsLeft"
	metricDocsTotalReplicaVBucketsLeft = "progressReplicaVBucketsLeft"
)

type taskCollector struct {
	m MetaCollector

	config *objects.CollectorConfig
}

func NewTaskCollector(client util.CbClient, config *objects.CollectorConfig) prometheus.Collector {
	if config == nil {
		config = objects.GetTaskCollectorDefaultConfig()
	}

	return &taskCollector{
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

func (c *taskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	for _, value := range c.config.Metrics {
		if !value.Enabled {
			continue
		}

		ch <- value.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem)
	}
}

func (c *taskCollector) collectTasks(ch chan<- prometheus.Metric, clusterName string, tasks []objects.Task) map[string]bool {
	var compactsReported = map[string]bool{}

	for _, task := range tasks {
		switch task.Type {
		case taskRebalance:
			c.addRebalance(ch, task, clusterName)
		case taskBucketCompaction:
			// XXX: there can be more than one compacting tasks for the same
			// bucket for now, let's report just the first.
			c.addBucketCompaction(ch, task, clusterName, compactsReported[task.Bucket])
			compactsReported[task.Bucket] = true
		case taskXdcr:
			log.Debug("found xdcr tasks from %s to %s", task.Source, task.Target)
			c.addXdcr(ch, task, clusterName)
		case taskClusterLogCollection:
			c.addClusterLogCollection(ch, task, clusterName)
		default:
			log.Warn("not implemented")
		}
	}

	return compactsReported
}
func (c *taskCollector) addBucketCompaction(ch chan<- prometheus.Metric, task objects.Task, clusterName string, compactsReported bool) {
	if cp, ok := c.config.Metrics[metricCompacting]; ok && cp.Enabled && !compactsReported {
		ch <- prometheus.MustNewConstMetric(
			cp.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			task.Progress,
			task.Bucket,
			clusterName)
	}
}

func (c *taskCollector) addClusterLogCollection(ch chan<- prometheus.Metric, task objects.Task, clusterName string) {
	if clc, ok := c.config.Metrics[taskClusterLogCollection]; ok && clc.Enabled {
		ch <- prometheus.MustNewConstMetric(
			clc.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			task.Progress,
			clusterName)
	}
}

func (c *taskCollector) addRebalance(ch chan<- prometheus.Metric, task objects.Task, clusterName string) {
	if rb, ok := c.config.Metrics[taskRebalance]; ok && rb.Enabled {
		ch <- prometheus.MustNewConstMetric(rb.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem), prometheus.GaugeValue, task.Progress, clusterName)
	}

	if rbPN, ok := c.config.Metrics[metricRebalancePerNode]; ok && rbPN.Enabled {
		for node, progress := range task.PerNode {
			ch <- prometheus.MustNewConstMetric(
				rbPN.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				progress.Progress,
				node,
				clusterName)
		}
	}
}

// nolint: cyclop
func (c *taskCollector) addXdcr(ch chan<- prometheus.Metric, task objects.Task, clusterName string) {
	if xcl, ok := c.config.Metrics[metricXdcrChangesLeft]; ok && xcl.Enabled {
		ch <- prometheus.MustNewConstMetric(
			xcl.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			float64(task.ChangesLeft),
			task.Source,
			task.Target,
			clusterName)
	}

	if xdc, ok := c.config.Metrics[metricXdcrDocsChecked]; ok && xdc.Enabled {
		ch <- prometheus.MustNewConstMetric(
			xdc.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			float64(task.DocsChecked),
			task.Source,
			task.Target,
			clusterName)
	}

	if xdw, ok := c.config.Metrics[metricXdcrDocsWritten]; ok && xdw.Enabled {
		ch <- prometheus.MustNewConstMetric(
			xdw.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			float64(task.DocsWritten),
			task.Source,
			task.Target,
			clusterName)
	}

	if xp, ok := c.config.Metrics[metricXdcrPaused]; ok && xp.Enabled {
		ch <- prometheus.MustNewConstMetric(
			xp.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			boolToFloat64(task.PauseRequested),
			task.Source,
			task.Target,
			clusterName)
	}

	if xe, ok := c.config.Metrics[metricXdcrErrors]; ok && xe.Enabled {
		ch <- prometheus.MustNewConstMetric(
			xe.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
			prometheus.GaugeValue,
			float64(len(task.Errors)),
			task.Source,
			task.Target,
			clusterName)
	}

	for _, data := range task.DetailedProgress.PerNode {
		// for each node grab these specific metrics from the config (if they exist)
		// then grab their data from the request and dump it into prometheus.
		if dt, ok := c.config.Metrics[metricDocsTotal]; ok && dt.Enabled {
			ch <- prometheus.MustNewConstMetric(
				dt.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				float64(data.Ingoing.DocsTotal),
				task.DetailedProgress.Bucket,
				task.Source,
				task.Target,
				clusterName)
		}

		if dtrans, ok := c.config.Metrics[metricDocsTransferred]; ok && dtrans.Enabled {
			ch <- prometheus.MustNewConstMetric(
				dtrans.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				float64(data.Ingoing.DocsTransferred),
				task.DetailedProgress.Bucket,
				task.Source,
				task.Target,
				clusterName)
		}

		if avbl, ok := c.config.Metrics[metricDocsActiveVbucketsLeft]; ok && avbl.Enabled {
			ch <- prometheus.MustNewConstMetric(
				avbl.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				float64(data.Ingoing.ActiveVBucketsLeft),
				task.DetailedProgress.Bucket,
				task.Source,
				task.Target,
				clusterName)
		}

		if rvbl, ok := c.config.Metrics[metricDocsTotalReplicaVBucketsLeft]; ok && rvbl.Enabled {
			ch <- prometheus.MustNewConstMetric(
				rvbl.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem),
				prometheus.GaugeValue,
				float64(data.Ingoing.ReplicaVBucketsLeft),
				task.DetailedProgress.Bucket,
				task.Source,
				task.Target,
				clusterName)
		}
	}
}

func (c *taskCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()

	log.Info("Collecting tasks metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error("%s", err)

		return
	}

	tasks, err := c.m.client.Tasks()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error("failed to scrape tasks")

		return
	}

	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)

		log.Error("failed to scrape tasks")

		return
	}

	// nolint: lll
	compactsReported := c.collectTasks(ch, clusterName, tasks)
	// always report the compacting task, even if it is not happening
	// this is to not break dashboards and make it easier to test alert rule
	// and etc.
	compact := c.config.Metrics[metricCompacting]

	for _, bucket := range buckets {
		if _, ok := compactsReported[bucket.Name]; !ok {
			// nolint: lll
			ch <- prometheus.MustNewConstMetric(compact.GetPrometheusDescription(c.config.Namespace, c.config.Subsystem), prometheus.GaugeValue, 0, bucket.Name, clusterName)
		}

		compactsReported[bucket.Name] = true
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	// nolint: lll
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}
