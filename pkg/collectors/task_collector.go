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
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

type taskCollector struct {
	m MetaCollector

	rebalance             *prometheus.Desc
	rebalancePerNode      *prometheus.Desc
	compacting            *prometheus.Desc
	clusterLogsCollection *prometheus.Desc
	xdcrChangesLeft       *prometheus.Desc
	xdcrDocsChecked       *prometheus.Desc
	xdcrDocsWritten       *prometheus.Desc
	xdcrPaused            *prometheus.Desc
	xdcrErrors            *prometheus.Desc

	ProgressDocsTotal           *prometheus.Desc
	ProgressDocsTransferred     *prometheus.Desc
	ProgressActiveVBucketsLeft  *prometheus.Desc
	ProgressReplicaVBucketsLeft *prometheus.Desc
}

func NewTaskCollector(client util.Client) prometheus.Collector {
	const subsystem = "task"

	return &taskCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "up"),
				"Couchbase task API is responding",
				[]string{"cluster"},
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "scrape_duration_seconds"),
				"Scrape duration in seconds",
				[]string{"cluster"},
				nil,
			),
		},
		rebalance: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "rebalance_progress"),
			"Progress of a rebalance task",
			[]string{"cluster"},
			nil,
		),
		rebalancePerNode: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "node_rebalance_progress"),
			"Progress of a rebalance task per node",
			[]string{"node", "cluster"},
			nil,
		),
		compacting: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "compacting_progress"),
			"Progress of a bucket compaction task",
			[]string{"bucket", "cluster"},
			nil,
		),
		clusterLogsCollection: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "cluster_logs_collection_progress"),
			"Progress of a cluster logs collection task",
			[]string{"cluster"},
			nil,
		),
		xdcrChangesLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "xdcr_changes_left"),
			"Number of updates still pending replication",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		xdcrDocsChecked: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "xdcr_docs_checked"),
			"Number of documents checked for changes",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		xdcrDocsWritten: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "xdcr_docs_written"),
			"Number of documents written to the destination cluster",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		xdcrPaused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "xdcr_paused"),
			"Is this replication paused",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		xdcrErrors: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "xdcr_errors"),
			"Number of errors",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		ProgressDocsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "docs_total"),
			"docs_total",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		ProgressDocsTransferred: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "docs_transferred"),
			"docs_transferred",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		ProgressActiveVBucketsLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "active_vbuckets_left"),
			"Number of Active VBuckets remaining",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
		ProgressReplicaVBucketsLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "replica_vbuckets_left"),
			"Number of Replica VBuckets remaining",
			[]string{"bucket", "target", "cluster"},
			nil,
		),
	}
}

func (c *taskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.rebalance
	ch <- c.rebalancePerNode
	ch <- c.compacting
	ch <- c.clusterLogsCollection
	ch <- c.xdcrChangesLeft
	ch <- c.xdcrDocsChecked
	ch <- c.xdcrDocsWritten
	ch <- c.xdcrPaused
	ch <- c.xdcrErrors
}

func (c *taskCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting tasks metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
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

	var compactsReported = map[string]bool{}
	// nolint: lll
	for _, task := range tasks {
		switch task.Type {
		case "rebalance":
			ch <- prometheus.MustNewConstMetric(c.rebalance, prometheus.GaugeValue, task.Progress, clusterName)
			for node, progress := range task.PerNode {
				ch <- prometheus.MustNewConstMetric(c.rebalancePerNode, prometheus.GaugeValue, progress.Progress, node, clusterName)
			}
		case "bucket_compaction":
			// XXX: there can be more than one compacting tasks for the same
			// bucket for now, let's report just the first.
			if ok := compactsReported[task.Bucket]; !ok {
				ch <- prometheus.MustNewConstMetric(c.compacting, prometheus.GaugeValue, task.Progress, task.Bucket, clusterName)
			}
			compactsReported[task.Bucket] = true
		case "xdcr":
			log.Debug("found xdcr tasks from %s to %s", task.Source, task.Target)
			ch <- prometheus.MustNewConstMetric(c.xdcrChangesLeft, prometheus.GaugeValue, float64(task.ChangesLeft), task.Source, task.Target, clusterName)
			ch <- prometheus.MustNewConstMetric(c.xdcrDocsChecked, prometheus.GaugeValue, float64(task.DocsChecked), task.Source, task.Target, clusterName)
			ch <- prometheus.MustNewConstMetric(c.xdcrDocsWritten, prometheus.GaugeValue, float64(task.DocsWritten), task.Source, task.Target, clusterName)
			ch <- prometheus.MustNewConstMetric(c.xdcrPaused, prometheus.GaugeValue, boolToFloat64(task.PauseRequested), task.Source, task.Target, clusterName)
			ch <- prometheus.MustNewConstMetric(c.xdcrErrors, prometheus.GaugeValue, float64(len(task.Errors)), task.Source, task.Target, clusterName)
		case "clusterLogsCollection":
			ch <- prometheus.MustNewConstMetric(c.clusterLogsCollection, prometheus.GaugeValue, task.Progress, clusterName)
		default:
			log.Warn("not implemented")
		}
	}
	// always report the compacting task, even if it is not happening
	// this is to not break dashboards and make it easier to test alert rule
	// and etc.
	for _, bucket := range buckets {
		if ok := compactsReported[bucket.Name]; !ok {
			// nolint: lll
			ch <- prometheus.MustNewConstMetric(c.compacting, prometheus.GaugeValue, 0, bucket.Name, clusterName)
		}
		compactsReported[bucket.Name] = true
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	// nolint: lll
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}
