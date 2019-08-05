package collectors

import (
	"github.com/couchbase/couchbase_exporter/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
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
				prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "up"),
				"Couchbase task API is responding",
				nil,
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "scrape_duration_seconds"),
				"Scrape duration in seconds",
				nil,
				nil,
			),
		},
		rebalance: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "rebalance_progress"),
			"Progress of a rebalance task",
			nil,
			nil,
		),
		rebalancePerNode: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "node_rebalance_progress"),
			"Progress of a rebalance task per node",
			[]string{"node"},
			nil,
		),
		compacting: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "compacting_progress"),
			"Progress of a bucket compaction task",
			[]string{"bucket"},
			nil,
		),
		clusterLogsCollection: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cluster_logs_collection_progress"),
			"Progress of a cluster logs collection task",
			nil,
			nil,
		),
		xdcrChangesLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdcr_changes_left"),
			"Number of updates still pending replication",
			[]string{"bucket", "target"},
			nil,
		),
		xdcrDocsChecked: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdcr_docs_checked"),
			"Number of documents checked for changes",
			[]string{"bucket", "target"},
			nil,
		),
		xdcrDocsWritten: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdcr_docs_written"),
			"Number of documents written to the destination cluster",
			[]string{"bucket", "target"},
			nil,
		),
		xdcrPaused: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdcr_paused"),
			"Is this replication paused",
			[]string{"bucket", "target"},
			nil,
		),
		xdcrErrors: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdcr_errors"),
			"Number of errors",
			[]string{"bucket", "target"},
			nil,
		),
		ProgressDocsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "docs_total"),
			"docs_total",
			[]string{"bucket", "target"},
			nil,
		),
		ProgressDocsTransferred: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "docs_transferred"),
			"docs_transferred",
			[]string{"bucket", "target"},
			nil,
		),
		ProgressActiveVBucketsLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "active_vbuckets_left"),
			"Number of Active VBuckets remaining",
			[]string{"bucket", "target"},
			nil,
		),
		ProgressReplicaVBucketsLeft: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "replica_vbuckets_left"),
			"Number of Replica VBuckets remaining",
			[]string{"bucket", "target"},
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

	tasks, err := c.m.client.Tasks()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape tasks")
		return
	}
	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape tasks")
		return
	}

	var compactsReported = map[string]bool{}
	// nolint: lll
	for _, task := range tasks {
		switch task.Type {
		case "rebalance":
			ch <- prometheus.MustNewConstMetric(c.rebalance, prometheus.GaugeValue, task.Progress)
			for node, progress := range task.PerNode {
				ch <- prometheus.MustNewConstMetric(c.rebalancePerNode, prometheus.GaugeValue, progress.Progress, node)
			}
		case "bucket_compaction":
			// XXX: there can be more than one compacting tasks for the same
			// bucket for now, let's report just the first.
			if ok := compactsReported[task.Bucket]; !ok {
				ch <- prometheus.MustNewConstMetric(c.compacting, prometheus.GaugeValue, task.Progress, task.Bucket)
			}
			compactsReported[task.Bucket] = true
		case "xdcr":
			log.Debugf("found xdcr tasks from %s to %s", task.Source, task.Target)
			ch <- prometheus.MustNewConstMetric(c.xdcrChangesLeft, prometheus.GaugeValue, float64(task.ChangesLeft), task.Source, task.Target)
			ch <- prometheus.MustNewConstMetric(c.xdcrDocsChecked, prometheus.GaugeValue, float64(task.DocsChecked), task.Source, task.Target)
			ch <- prometheus.MustNewConstMetric(c.xdcrDocsWritten, prometheus.GaugeValue, float64(task.DocsWritten), task.Source, task.Target)
			//ch <- prometheus.MustNewConstMetric(c.xdcrPaused, prometheus.GaugeValue, fromBool(task.PauseRequested), task.Source, task.Target)
			ch <- prometheus.MustNewConstMetric(c.xdcrErrors, prometheus.GaugeValue, float64(len(task.Errors)), task.Source, task.Target)
		case "clusterLogsCollection":
			ch <- prometheus.MustNewConstMetric(c.clusterLogsCollection, prometheus.GaugeValue, task.Progress)
		default:
			log.With("type", task.Type).Warn("not implemented")
		}
	}
	// always report the compacting task, even if it is not happening
	// this is to not break dashboards and make it easier to test alert rule
	// and etc.
	for _, bucket := range buckets {
		if ok := compactsReported[bucket.Name]; !ok {
			// nolint: lll
			ch <- prometheus.MustNewConstMetric(c.compacting, prometheus.GaugeValue, 0, bucket.Name)
		}
		compactsReported[bucket.Name] = true
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	// nolint: lll
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
