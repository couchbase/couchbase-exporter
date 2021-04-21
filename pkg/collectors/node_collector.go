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
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
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
	clusterMembership    *prometheus.Desc

	ctrRebalanceSuccess        *prometheus.Desc
	ctrRebalanceStart          *prometheus.Desc
	ctrRebalanceFail           *prometheus.Desc
	ctrRebalanceStop           *prometheus.Desc
	ctrFailoverNode            *prometheus.Desc
	ctrFailover                *prometheus.Desc
	ctrFailoverComplete        *prometheus.Desc
	ctrFailoverIncomplete      *prometheus.Desc
	ctrGracefulFailoverStart   *prometheus.Desc
	ctrGracefulFailoverSuccess *prometheus.Desc
	ctrGracefulFailoverFail    *prometheus.Desc
}

func NewNodesCollector(client util.Client) prometheus.Collector {
	const subsystem = "node"
	// nolint: lll
	return &nodesCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "up"),
				"Couchbase cluster API is responding",
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
		healthy: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "healthy"),
			"Is this node healthy",
			[]string{"node", "cluster"},
			nil,
		),
		systemStatsCPUUtilizationRate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "systemstats_cpu_utilization_rate"),
			"Percentage of CPU in use across all available cores on this server.",
			[]string{"node", "cluster"},
			nil,
		),
		systemStatsSwapTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "systemstats_swap_total"),
			"Bytes of total swap space available on this server.",
			[]string{"node", "cluster"},
			nil,
		),
		systemStatsSwapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "systemstats_swap_used"),
			"Bytes of swap space in use on this server.",
			[]string{"node", "cluster"},
			nil,
		),
		systemStatsMemTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "systemstats_mem_total"),
			"Bytes of total memory available on this server.",
			[]string{"node", "cluster"},
			nil,
		),
		systemStatsMemFree: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "systemstats_mem_free"),
			"Bytes of memory not in use on this server.",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchDocsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_docs_actual_disk_size"),
			"The size of all data service files on disk for this bucket, including the data itself, metadata, and temporary files. (measured from couch_docs_actual_disk_size)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchDocsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_docs_data_size"),
			"Bytes of active data in this bucket. (measured from couch_docs_data_size)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchViewsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_views_actual_disk_size"),
			"Bytes of active items in all the views for this bucket on disk (measured from couch_views_actual_disk_size)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchViewsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_views_data_size"),
			"Bytes of active data for all the views in this bucket. (measured from couch_views_data_size)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsMemUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_mem_used"),
			"Total memory used in bytes. (as measured from mem_used)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsOps: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_ops"),
			"Total operations per second (including XDCR) to this bucket. (measured from cmd_get + cmd_set + incr_misses + incr_hits + decr_misses + decr_hits + delete_misses + delete_hits + ep_num_ops_del_meta + ep_num_ops_get_meta + ep_num_ops_set_meta)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_curr_items"),
			"Current number of unique items in Couchbase",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCurrItemsTot: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_curr_items_tot"),
			"Current number of items in Couchbase including replicas",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsVbReplicaCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_vb_replica_curr_items"),
			"Number of items in replica vBuckets in this bucket. (measured from vb_replica_curr_items)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsVbActiveNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_vb_active_number_non_resident"),
			"interestingstats_vb_active_number_non_resident",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchSpatialDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_spatial_disk_size"),
			"interestingstats_couch_spatial_disk_size",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCouchSpatialDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_couch_spatial_data_size"),
			"interestingstats_couch_spatial_data_size",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsCmdGet: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_cmd_get"),
			"Number of reads (get operations) per second from this bucket. (measured from cmd_get)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsGetHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_get_hits"),
			"Number of get operations per second for data that this bucket contains. (measured from get_hits)",
			[]string{"node", "cluster"},
			nil,
		),
		interestingStatsEpBgFetched: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "interestingstats_ep_bg_fetched"),
			"Number of reads per second from disk for this bucket. (measured from ep_bg_fetched)",
			[]string{"node", "cluster"},
			nil,
		),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "uptime"),
			"uptime",
			[]string{"node", "cluster"},
			nil,
		),
		memoryTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memory_total"),
			"memory_total",
			[]string{"node", "cluster"},
			nil,
		),
		memoryFree: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memory_free"),
			"memory_free",
			[]string{"node", "cluster"},
			nil,
		),
		mcdMemoryAllocated: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memcached_memory_allocated"),
			"memcached_memory_allocated",
			[]string{"node", "cluster"},
			nil,
		),
		mcdMemoryReserved: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memcached_memory_reserved"),
			"memcached_memory_reserved",
			[]string{"node", "cluster"},
			nil,
		),
		clusterMembership: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "cluster_membership"),
			"whether or not node is part of the CB cluster",
			[]string{"node", "cluster"},
			nil,
		),
		ctrFailover: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "failover"),
			"failover",
			[]string{"cluster"},
			nil,
		),
		ctrFailoverNode: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "failover_node"),
			"failover_node",
			[]string{"cluster"},
			nil,
		),
		ctrFailoverComplete: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "failover_complete"),
			"failover_complete",
			[]string{"cluster"},
			nil,
		),
		ctrFailoverIncomplete: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "failover_incomplete"),
			"failover_incomplete",
			[]string{"cluster"},
			nil,
		),
		ctrRebalanceStart: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "rebalance_start"),
			"rebalance_start",
			[]string{"cluster"},
			nil,
		),
		ctrRebalanceStop: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "rebalance_stop"),
			"rebalance_stop",
			[]string{"cluster"},
			nil,
		),
		ctrRebalanceSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "rebalance_success"),
			"rebalance_success",
			[]string{"cluster"},
			nil,
		),
		ctrRebalanceFail: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "rebalance_failure"),
			"rebalance_failure",
			[]string{"cluster"},
			nil,
		),
		ctrGracefulFailoverStart: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "graceful_failover_start"),
			"graceful_failover_start",
			[]string{"cluster"},
			nil,
		),
		ctrGracefulFailoverSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "graceful_failover_success"),
			"graceful_failover_success",
			[]string{"cluster"},
			nil,
		),
		ctrGracefulFailoverFail: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "graceful_failover_fail"),
			"graceful_failover_fail",
			[]string{"cluster"},
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
	ch <- c.systemStatsCPUUtilizationRate
	ch <- c.systemStatsMemFree
	ch <- c.systemStatsMemTotal
	ch <- c.systemStatsSwapTotal
	ch <- c.systemStatsSwapUsed
	ch <- c.uptime
	ch <- c.memoryTotal
	ch <- c.memoryFree
	ch <- c.mcdMemoryAllocated
	ch <- c.mcdMemoryReserved
	ch <- c.clusterMembership
	ch <- c.ctrFailover
	ch <- c.ctrFailoverNode
	ch <- c.ctrFailoverComplete
	ch <- c.ctrFailoverIncomplete
	ch <- c.ctrRebalanceStart
	ch <- c.ctrRebalanceStop
	ch <- c.ctrRebalanceSuccess
	ch <- c.ctrRebalanceFail
	ch <- c.ctrGracefulFailoverSuccess
	ch <- c.ctrGracefulFailoverFail
	ch <- c.ctrGracefulFailoverStart
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

// Collect all metrics
func (c *nodesCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting nodes metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("%s", err)
		return
	}

	nodes, err := c.m.client.Nodes()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("failed to scrape nodes")
		return
	}

	for _, node := range nodes.Nodes {
		log.Debug("Collecting %s node metrics...", node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.healthy, prometheus.GaugeValue, boolToFloat64(node.Status == "healthy"), node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchDocsActualDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchDocsActualDiskSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchDocsDataSize, prometheus.GaugeValue, node.InterestingStats.CouchDocsDataSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchViewsActualDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchViewsActualDiskSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchViewsDataSize, prometheus.GaugeValue, node.InterestingStats.CouchViewsDataSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsMemUsed, prometheus.GaugeValue, node.InterestingStats.MemUsed, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsOps, prometheus.GaugeValue, node.InterestingStats.Ops, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCurrItems, prometheus.GaugeValue, node.InterestingStats.CurrItems, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCurrItemsTot, prometheus.GaugeValue, node.InterestingStats.CurrItemsTot, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsVbReplicaCurrItems, prometheus.GaugeValue, node.InterestingStats.VbReplicaCurrItems, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchSpatialDiskSize, prometheus.GaugeValue, node.InterestingStats.CouchSpatialDiskSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCouchSpatialDataSize, prometheus.GaugeValue, node.InterestingStats.CouchSpatialDataSize, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsCmdGet, prometheus.GaugeValue, node.InterestingStats.CmdGet, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsGetHits, prometheus.GaugeValue, node.InterestingStats.GetHits, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.interestingStatsEpBgFetched, prometheus.GaugeValue, node.InterestingStats.EpBgFetched, node.Hostname, clusterName)

		ch <- prometheus.MustNewConstMetric(c.systemStatsCPUUtilizationRate, prometheus.GaugeValue, node.SystemStats.CPUUtilizationRate, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.systemStatsSwapUsed, prometheus.GaugeValue, node.SystemStats.SwapUsed, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.systemStatsSwapTotal, prometheus.GaugeValue, node.SystemStats.SwapTotal, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.systemStatsMemTotal, prometheus.GaugeValue, node.SystemStats.MemTotal, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.systemStatsMemFree, prometheus.GaugeValue, node.SystemStats.MemFree, node.Hostname, clusterName)

		up, err := strconv.ParseFloat(node.Uptime, 64)
		if err != nil {
			return
		}

		ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.CounterValue, up, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.memoryTotal, prometheus.GaugeValue, node.MemoryTotal, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.memoryFree, prometheus.GaugeValue, node.MemoryFree, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.mcdMemoryAllocated, prometheus.GaugeValue, node.McdMemoryAllocated, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.mcdMemoryReserved, prometheus.GaugeValue, node.McdMemoryReserved, node.Hostname, clusterName)
		ch <- prometheus.MustNewConstMetric(c.clusterMembership, prometheus.CounterValue, ifActive(node.ClusterMembership), node.Hostname, clusterName)
	}

	ch <- prometheus.MustNewConstMetric(c.ctrFailover, prometheus.CounterValue, nodes.Counters.Failover, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrFailoverNode, prometheus.CounterValue, nodes.Counters.FailoverNode, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrFailoverComplete, prometheus.CounterValue, nodes.Counters.FailoverComplete, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrFailoverIncomplete, prometheus.CounterValue, nodes.Counters.FailoverIncomplete, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrRebalanceStart, prometheus.CounterValue, nodes.Counters.RebalanceStart, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrRebalanceStop, prometheus.CounterValue, nodes.Counters.RebalanceStop, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrRebalanceSuccess, prometheus.CounterValue, nodes.Counters.RebalanceSuccess, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrRebalanceFail, prometheus.CounterValue, nodes.Counters.RebalanceFail, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrGracefulFailoverSuccess, prometheus.CounterValue, nodes.Counters.GracefulFailoverSuccess, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrGracefulFailoverFail, prometheus.CounterValue, nodes.Counters.GracefulFailoverFail, clusterName)
	ch <- prometheus.MustNewConstMetric(c.ctrGracefulFailoverStart, prometheus.CounterValue, nodes.Counters.GracefulFailoverStart, clusterName)

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}
