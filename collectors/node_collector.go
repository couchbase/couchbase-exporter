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
	"github.com/couchbase/couchbase_exporter/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type nodesCollector struct {
	m       MetaCollector
	healthy *prometheus.Desc

	systemStatsCPUUtilizationRate *prometheus.Desc
	systemStatsSwapTotal          *prometheus.Desc
	systemStatsSwapUsed           *prometheus.Desc
	systemStatsMemTotal           *prometheus.Desc
	systemStatsMemFree            *prometheus.Desc

	interestingStatsCmdGet                   *prometheus.Desc
	interestingStatsCouchDocsActualDiskSize  *prometheus.Desc
	interestingStatsCouchDocsDataSize        *prometheus.Desc
	interestingStatsCouchSpatialDataSize     *prometheus.Desc
	interestingStatsCouchSpatialDiskSize     *prometheus.Desc
	interestingStatsCouchViewsActualDiskSize *prometheus.Desc
	interestingStatsCouchViewsDataSize       *prometheus.Desc
	interestingStatsCurrItems                *prometheus.Desc
	interestingStatsCurrItemsTot             *prometheus.Desc
	interestingStatsEpBgFetched              *prometheus.Desc
	interestingStatsGetHits                  *prometheus.Desc
	interestingStatsMemUsed                  *prometheus.Desc
	interestingStatsOps                      *prometheus.Desc
	interestingStatsVbActiveNumNonResident   *prometheus.Desc // ??
	interestingStatsVbReplicaCurrItems       *prometheus.Desc

	uptime             *prometheus.Desc
	memoryTotal        *prometheus.Desc
	memoryFree         *prometheus.Desc
	mcdMemoryReserved  *prometheus.Desc
	mcdMemoryAllocated *prometheus.Desc
	//thisNode
	clusterCompatibility *prometheus.Desc
}

func NewNodesCollector(client util.Client) prometheus.Collector {
	const subsystem = "node"
	// nolint: lll
	return &nodesCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "up"),
				"Couchbase cluster API is responding",
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
		healthy: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "healthy"),
			"Is this node healthy",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchDocsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_docs_actual_disk_size"),
			"interestingstats_couch_docs_actual_disk_size",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchDocsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_docs_data_size"),
			"interestingstats_couch_docs_data_size",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchViewsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_views_actual_disk_size"),
			"interestingstats_couch_views_actual_disk_size",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchViewsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_views_data_size"),
			"interestingstats_couch_views_data_size",
			[]string{"node"},
			nil,
		),
		interestingStatsMemUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_mem_used"),
			"interestingstats_mem_used",
			[]string{"node"},
			nil,
		),
		interestingStatsOps: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_ops"),
			"interestingstats_ops",
			[]string{"node"},
			nil,
		),
		interestingStatsCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_curr_items"),
			"Current number of unique items in Couchbase",
			[]string{"node"},
			nil,
		),
		interestingStatsCurrItemsTot: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_curr_items_tot"),
			"Current number of items in Couchbase including replicas",
			[]string{"node"},
			nil,
		),
		interestingStatsVbReplicaCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_vb_replica_curr_items"),
			"interestingstats_vb_replica_curr_items",
			[]string{"node"},
			nil,
		),
		interestingStatsVbActiveNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_vb_active_number_non_resident"),
			"interestingstats_vb_active_number_non_resident",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchSpatialDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_spatial_disk_size"),
			"interestingstats_couch_spatial_disk_size",
			[]string{"node"},
			nil,
		),
		interestingStatsCouchSpatialDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_couch_spatial_data_size"),
			"interestingstats_couch_spatial_data_size",
			[]string{"node"},
			nil,
		),
		interestingStatsCmdGet: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_cmd_get"),
			"interestingstats_cmd_get",
			[]string{"node"},
			nil,
		),
		interestingStatsGetHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_get_hits"),
			"interestingstats_get_hits",
			[]string{"node"},
			nil,
		),
		interestingStatsEpBgFetched: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "interestingstats_ep_bg_fetched"),
			"interestingstats_ep_bg_fetched",
			[]string{"node"},
			nil,
		),
	}
}

// Describe all metrics
func (c *nodesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.healthy
	ch <- c.interestingStatsCouchDocsActualDiskSize
	ch <- c.interestingStatsCouchDocsDataSize
	ch <- c.interestingStatsCouchViewsActualDiskSize
	ch <- c.interestingStatsCouchViewsDataSize
	ch <- c.interestingStatsMemUsed
	ch <- c.interestingStatsOps
	ch <- c.interestingStatsCurrItems
	ch <- c.interestingStatsCurrItemsTot
	ch <- c.interestingStatsVbReplicaCurrItems
	ch <- c.interestingStatsCouchSpatialDiskSize
	ch <- c.interestingStatsCouchSpatialDataSize
	ch <- c.interestingStatsCmdGet
	ch <- c.interestingStatsGetHits
	ch <- c.interestingStatsEpBgFetched
}

// Collect all metrics
func (c *nodesCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting nodes metrics...")

	nodes, err := c.m.client.Nodes()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape nodes")
		return
	}

	// nolint: lll
	for _, node := range nodes.Nodes {
		log.Debugf("Collecting %s node metrics...", node.Hostname)
		//ch <- prometheus.MustNewConstMetric(c.healthy, prometheus.GaugeValue, fromBool(node.Status == "healthy"), node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchDocsActualDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchDocsActualDiskSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchDocsDataSize, prometheus.GaugeValue, node.InterestingStats.CouchDocsDataSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchViewsActualDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchViewsActualDiskSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchViewsDataSize, prometheus.GaugeValue, node.InterestingStats.CouchViewsDataSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsMemUsed, prometheus.GaugeValue, node.InterestingStats.MemUsed, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsOps, prometheus.GaugeValue, node.InterestingStats.Ops, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCurrItems, prometheus.GaugeValue, node.InterestingStats.CurrItems, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCurrItemsTot, prometheus.GaugeValue, node.InterestingStats.CurrItemsTot, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsVbReplicaCurrItems, prometheus.GaugeValue, node.InterestingStats.VbReplicaCurrItems, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchSpatialDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchSpatialDiskSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchSpatialDataSize, prometheus.GaugeValue, node.InterestingStats.CouchSpatialDataSize, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCmdGet, prometheus.GaugeValue, node.InterestingStats.CmdGet, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsGetHits, prometheus.GaugeValue, node.InterestingStats.GetHits, node.Hostname)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsEpBgFetched, prometheus.GaugeValue, node.InterestingStats.EpBgFetched, node.Hostname)
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	// nolint: lll
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
