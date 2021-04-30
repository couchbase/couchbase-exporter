//  Copyright (c) 2021 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	DefaultNamespace                = "cb"
	DefaultUptimeMetric             = "up"
	DefaultScrapeDurationMetric     = "scrape_duration_seconds"
	DefaultUptimeMetricHelp         = "Couchbase cluster API is responding"
	DefaultScrapeDurationMetricHelp = "Scrape duration in seconds"
	ClusterLabel                    = "cluster"
	NodeLabel                       = "node"
	BucketLabel                     = "bucket"
)

type CollectorConfig struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	Subsystem string                `json:"subsystem"`
	Metrics   map[string]MetricInfo `json:"metrics"`
}

func (m *MetricInfo) GetPrometheusDescription(namespace string, subsystem string) *prometheus.Desc {
	name := m.Name
	if m.NameOverride != "" {
		name = m.NameOverride
	}

	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		m.HelpText,
		m.Labels,
		nil)
}

func (m *MetricInfo) GetPrometheusGaugeVec(namespace string, subsystem string) *prometheus.GaugeVec {
	name := m.Name
	if m.NameOverride != "" {
		name = m.NameOverride
	}

	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        name,
			Help:        m.HelpText,
			ConstLabels: nil,
		},
		m.Labels)
}

type MetricInfo struct {
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	NameOverride string   `json:"nameOverride"`
	HelpText     string   `json:"helpText"`
	Labels       []string `json:"labels"`
}

func GetQueryCollectorDefaultConfig() *CollectorConfig {
	return queryCollectorDefaultConfig()
}

func GetBucketInfoCollectorDefaultConfig() *CollectorConfig {
	return bucketInfoCollectorDefaultConfig()
}

func GetAnalyticsCollectorDefaultConfig() *CollectorConfig {
	return analyticsCollectorDefaultConfig()
}

func GetIndexCollectorDefaultConfig() *CollectorConfig {
	return indexCollectorDefaultConfig()
}

func GetSearchCollectorDefaultConfig() *CollectorConfig {
	return searchCollectorDefaultConfig()
}

func GetTaskCollectorDefaultConfig() *CollectorConfig {
	return taskCollectorDefaultConfig()
}

func GetEventingCollectorDefaultConfig() *CollectorConfig {
	return eventingCollectorDefaultConfig()
}

func GetNodeCollectorDefaultConfig() *CollectorConfig {
	return nodeCollectorDefaultConfig()
}

func GetBucketStatsCollectorDefaultConfig() *CollectorConfig {
	return bucketStatsCollectorDefaultConfig()
}

func GetPerNodeBucketStatsCollectorDefaultConfig() *CollectorConfig {
	return perNodeBucketStatsCollectorDefaultConfig()
}

func perNodeBucketStatsCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "PerNodeBucketStats",
		Namespace: DefaultNamespace,
		Subsystem: "pernodebucket",
		Metrics: map[string]MetricInfo{
			"AvgDiskUpdateTime": {
				NameOverride: "",
				Name:         "avg_disk_update_time",
				HelpText:     "Average disk update time in microseconds as from disk_update histogram of timings",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"AvgDiskCommitTime": {
				NameOverride: "",
				Name:         "avg_disk_commit_time",
				HelpText:     "Average disk commit time in seconds as from disk_update histogram of timings",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"AvgBgWaitTime": {
				NameOverride: "",
				Name:         "avg_bg_wait_seconds",
				HelpText:     " ",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"AvgActiveTimestampDrift": {
				NameOverride: "",
				Name:         "avg_active_timestamp_drift",
				HelpText:     "Average drift (in seconds) between mutation timestamps and the local time for active vBuckets. (measured from ep_active_hlc_drift and ep_active_hlc_drift_count)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"AvgReplicaTimestampDrift": {
				NameOverride: "",
				Name:         "avg_replica_timestamp_drift",
				HelpText:     "Average drift (in seconds) between mutation timestamps and the local time for replica vBuckets. (measured from ep_replica_hlc_drift and ep_replica_hlc_drift_count)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchTotalDiskSize": {
				NameOverride: "",
				Name:         "couch_total_disk_size",
				HelpText:     "The total size on disk of all data and view files for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchDocsFragmentation": {
				NameOverride: "",
				Name:         "couch_docs_fragmentation",
				HelpText:     "How much fragmented data there is to be compacted compared to real data for the data files in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchViewsFragmentation": {
				NameOverride: "",
				Name:         "couch_views_fragmentation",
				HelpText:     "How much fragmented data there is to be compacted compared to real data for the view index files in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchDocsActualDiskSize": {
				NameOverride: "",
				Name:         "couch_docs_actual_disk_size",
				HelpText:     "The size of all data files for this bucket, including the data itself, meta data and temporary files",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchDocsDataSize": {
				NameOverride: "",
				Name:         "couch_docs_data_size",
				HelpText:     "The size of active data in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchDocsDiskSize": {
				NameOverride: "",
				Name:         "couch_docs_disk_size",
				HelpText:     "The size of all data files for this bucket, including the data itself, meta data and temporary files",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchSpatialDataSize": {
				NameOverride: "",
				Name:         "couch_spatial_data_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchSpatialDiskSize": {
				NameOverride: "",
				Name:         "couch_spatial_disk_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchSpatialOps": {
				NameOverride: "",
				Name:         "couch_spatial_ops",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchViewsActualDiskSize": {
				NameOverride: "",
				Name:         "couch_views_actual_disk_size",
				HelpText:     "The size of all active items in all the indexes for this bucket on disk",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchViewsDataSize": {
				NameOverride: "",
				Name:         "couch_views_data_size",
				HelpText:     "The size of active data on for all the indexes in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchViewsDiskSize": {
				NameOverride: "",
				Name:         "couch_views_disk_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CouchViewsOps": {
				NameOverride: "",
				Name:         "couch_views_ops",
				HelpText:     "All the view reads for all design documents including scatter gather",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"HitRatio": {
				NameOverride: "",
				Name:         "hit_ratio",
				HelpText:     "Hit ratio",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpCacheMissRate": {
				NameOverride: "",
				Name:         "ep_cache_miss_rate",
				HelpText:     "Percentage of reads per second to this bucket from disk as opposed to RAM",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpResidentItemsRate": {
				NameOverride: "",
				Name:         "ep_resident_items_rate",
				HelpText:     "Percentage of all items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesCount": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_items_remaining",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_producer_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_items_sent",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_total_bytes",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsIndexesBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_views_indexes_backoff",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"BgWaitCount": {
				NameOverride: "",
				Name:         "bg_wait_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"BgWaitTotal": {
				NameOverride: "",
				Name:         "bg_wait_total",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"BytesRead": {
				NameOverride: "",
				Name:         "bytes_read",
				HelpText:     "Bytes Read",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"BytesWritten": {
				NameOverride: "",
				Name:         "bytes_written",
				HelpText:     "Bytes written",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CasBadVal": {
				NameOverride: "",
				Name:         "cas_bad_val",
				HelpText:     "Compare and Swap bad values",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CasHits": {
				NameOverride: "",
				Name:         "cas_hits",
				HelpText:     "Number of operations with a CAS id per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CasMisses": {
				NameOverride: "",
				Name:         "cas_misses",
				HelpText:     "Compare and Swap misses",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CmdGet": {
				NameOverride: "",
				Name:         "cmd_get",
				HelpText:     "Number of reads (get operations) per second from this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CmdSet": {
				NameOverride: "",
				Name:         "cmd_set",
				HelpText:     "Number of writes (set operations) per second to this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CurrConnections": {
				NameOverride: "",
				Name:         "curr_connections",
				HelpText:     "Number of connections to this server including connections from external client SDKs, proxies, DCP requests and internal statistic gathering",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CurrItems": {
				NameOverride: "",
				Name:         "curr_items",
				HelpText:     "Number of items in active vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CurrItemsTot": {
				NameOverride: "",
				Name:         "curr_items_tot",
				HelpText:     "Total number of items in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DecrHits": {
				NameOverride: "",
				Name:         "decr_hits",
				HelpText:     "Decrement hits",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DecrMisses": {
				NameOverride: "",
				Name:         "decr_misses",
				HelpText:     "Decrement misses",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DeleteHits": {
				NameOverride: "",
				Name:         "delete_hits",
				HelpText:     "Number of delete operations per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DeleteMisses": {
				NameOverride: "",
				Name:         "delete_misses",
				HelpText:     "Number of delete operations per second for data that this bucket does not contain. (measured from delete_misses)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DiskCommitCount": {
				NameOverride: "",
				Name:         "disk_commit_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DiskCommitTotal": {
				NameOverride: "",
				Name:         "disk_commit_total",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DiskUpdateCount": {
				NameOverride: "",
				Name:         "disk_update_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DiskUpdateTotal": {
				NameOverride: "",
				Name:         "disk_update_total",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"DiskWriteQueue": {
				NameOverride: "",
				Name:         "disk_write_queue",
				HelpText:     "Number of items waiting to be written to disk in this bucket. (measured from ep_queue_size+ep_flusher_todo)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpActiveAheadExceptions": {
				NameOverride: "",
				Name:         "ep_active_ahead_exceptions",
				HelpText:     "Total number of ahead exceptions (when timestamp drift between mutations and local time has exceeded 5000000 μs) per second for all active vBuckets.",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpActiveHlcDrift": {
				NameOverride: "",
				Name:         "ep_active_hlc_drift",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpActiveHlcDriftCount": {
				NameOverride: "",
				Name:         "ep_active_hlc_drift_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpBgFetched": {
				NameOverride: "",
				Name:         "ep_bg_fetched",
				HelpText:     "Number of reads per second from disk for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpClockCasDriftTheresholExceeded": {
				NameOverride: "",
				Name:         "ep_clock_cas_drift_threshold_exceeded",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDataReadFailed": {
				NameOverride: "",
				Name:         "ep_data_read_failed",
				HelpText:     "Number of disk read failures. (measured from ep_data_read_failed)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDataWriteFailed": {
				NameOverride: "",
				Name:         "ep_data_write_failed",
				HelpText:     "Number of disk write failures. (measured from ep_data_write_failed)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_2i_backoff",
				HelpText:     "Number of backoffs for indexes DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iCount": {
				NameOverride: "",
				Name:         "ep_dcp_2i_count",
				HelpText:     "Number of indexes DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_2i_items_remaining",
				HelpText:     "Number of indexes items remaining to be sent",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_2i_items_sent",
				HelpText:     "Number of indexes items sent",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_2i_producers",
				HelpText:     "Number of indexes producers",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_2i_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcp2iTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_2i_total_bytes",
				HelpText:     "Number of bytes per second being sent for indexes DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_backoff",
				HelpText:     "Number of backoffs per second for analytics DCP connections (measured from ep_dcp_cbas_backoff)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasCount": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_count",
				HelpText:     "Number of internal analytics DCP connections in this bucket (measured from ep_dcp_cbas_count)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_producer_count",
				HelpText:     "Number of analytics senders for this bucket (measured from ep_dcp_cbas_producer_count)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_cbas_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpCbasTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_total_bytes",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_fts_backoff",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsCount": {
				NameOverride: "",
				Name:         "ep_dcp_fts_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_fts_items_remaining",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_fts_items_sent",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_fts_producer_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_fts_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpFtsTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_fts_total_bytes",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_other_backoff",
				HelpText:     "Number of backoffs for other DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherCount": {
				NameOverride: "",
				Name:         "ep_dcp_other_count",
				HelpText:     "Number of other DCP connections in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_other_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket (measured from ep_dcp_other_items_remaining)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_other_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket (measured from ep_dcp_other_items_sent)",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_other_producer_count",
				HelpText:     "Number of other senders for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_other_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpOtherTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_other_total_bytes",
				HelpText:     "Number of bytes per second being sent for other DCP connections for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_replica_backoff",
				HelpText:     "Number of backoffs for replication DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaCount": {
				NameOverride: "",
				Name:         "ep_dcp_replica_count",
				HelpText:     "Number of internal replication DCP connections in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_replica_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_replica_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_replica_producer_count",
				HelpText:     "Number of replication senders for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_replica_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpReplicaTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_replica_total_bytes",
				HelpText:     "Number of bytes per second being sent for replication DCP connections for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_views_backoff",
				HelpText:     "Number of backoffs for views DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsCount": {
				NameOverride: "",
				Name:         "ep_dcp_views_count",
				HelpText:     "Number of views DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_views_items_remaining",
				HelpText:     "Number of views items remaining to be sent",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_views_items_sent",
				HelpText:     "Number of views items sent",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_views_producer_count",
				HelpText:     "Number of views producers",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_views_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpViewsTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_views_total_bytes",
				HelpText:     "Number bytes per second being sent for views DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrBackoff": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_backoff",
				HelpText:     "Number of backoffs for XDCR DCP connections",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrCount": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_count",
				HelpText:     "Number of internal XDCR DCP connections in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrItemsRemaining": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrItemsSent": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrProducerCount": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_producer_count",
				HelpText:     "Number of XDCR senders for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrTotalBacklogSize": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_total_backlog_size",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDcpXdcrTotalBytes": {
				NameOverride: "",
				Name:         "ep_dcp_xdcr_total_bytes",
				HelpText:     "Number of bytes per second being sent for XDCR DCP connections for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDiskqueueDrain": {
				NameOverride: "",
				Name:         "ep_diskqueue_drain",
				HelpText:     "Total number of items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDiskqueueFill": {
				NameOverride: "",
				Name:         "ep_diskqueue_fill",
				HelpText:     "Total number of items per second being put on the disk queue in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpDiskqueueItems": {
				NameOverride: "",
				Name:         "ep_diskqueue_items",
				HelpText:     "Total number of items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpFlusherTodo": {
				NameOverride: "",
				Name:         "ep_flusher_todo",
				HelpText:     "Number of items currently being written",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpItemCommitFailed": {
				NameOverride: "",
				Name:         "ep_item_commit_failed",
				HelpText:     "Number of times a transaction failed to commit due to storage errors",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpKvSize": {
				NameOverride: "",
				Name:         "ep_kv_size",
				HelpText:     "Total amount of user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpMaxSize": {
				NameOverride: "",
				Name:         "ep_max_size",
				HelpText:     "The maximum amount of memory this bucket can use",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpMemHighWat": {
				NameOverride: "",
				Name:         "ep_mem_high_wat",
				HelpText:     "High water mark for auto-evictions",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpMemLowWat": {
				NameOverride: "",
				Name:         "ep_mem_low_wat",
				HelpText:     "Low water mark for auto-evictions",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpMetaDataMemory": {
				NameOverride: "",
				Name:         "ep_meta_data_memory",
				HelpText:     "Total amount of item metadata consuming RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumNonResident": {
				NameOverride: "",
				Name:         "ep_num_non_resident",
				HelpText:     "Number of non-resident items",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumOpsDelMeta": {
				NameOverride: "",
				Name:         "ep_num_ops_del_meta",
				HelpText:     "Number of delete operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumOpsDelRetMeta": {
				NameOverride: "",
				Name:         "ep_num_ops_del_ret_meta",
				HelpText:     "Number of delRetMeta operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumOpsGetMeta": {
				NameOverride: "",
				Name:         "ep_num_ops_get_meta",
				HelpText:     "Number of metadata read operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumOpsSetMeta": {
				NameOverride: "",
				Name:         "ep_num_ops_set_meta",
				HelpText:     "Number of set operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumOpsSetRetMeta": {
				NameOverride: "",
				Name:         "ep_num_ops_set_ret_meta",
				HelpText:     "Number of setRetMeta operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpNumValueEjects": {
				NameOverride: "",
				Name:         "ep_num_value_ejects",
				HelpText:     "Total number of items per second being ejected to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpOomErrors": {
				NameOverride: "",
				Name:         "ep_oom_errors",
				HelpText:     "Number of times unrecoverable OOMs happened while processing operations",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpOpsCreate": {
				NameOverride: "",
				Name:         "ep_ops_create",
				HelpText:     "Total number of new items being inserted into this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpOpsUpdate": {
				NameOverride: "",
				Name:         "ep_ops_update",
				HelpText:     "Number of items updated on disk per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpOverhead": {
				NameOverride: "",
				Name:         "ep_overhead",
				HelpText:     "Extra memory used by transient data like persistence queues or checkpoints",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpQueueSize": {
				NameOverride: "",
				Name:         "ep_queue_size",
				HelpText:     "Number of items queued for storage",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpReplicaAheadExceptions": {
				NameOverride: "",
				Name:         "ep_replica_ahead_exceptions",
				HelpText:     "Percentage of all items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpReplicaHlcDrift": {
				NameOverride: "",
				Name:         "ep_replica_hlc_drift",
				HelpText:     "The sum of the total Absolute Drift, which is the accumulated drift observed by the vBucket. Drift is always accumulated as an absolute value.",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpReplicaHlcDriftCount": {
				NameOverride: "",
				Name:         "ep_replica_hlc_drift_count",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpTmpOomErrors": {
				NameOverride: "",
				Name:         "ep_tmp_oom_errors",
				HelpText:     "Number of back-offs sent per second to client SDKs due to OOM situations from this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"EpVbTotal": {
				NameOverride: "",
				Name:         "ep_vb_total",
				HelpText:     "Total number of vBuckets for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"Evictions": {
				NameOverride: "",
				Name:         "evictions",
				HelpText:     "Number of evictions",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},

			"GetHits": {
				NameOverride: "",
				Name:         "get_hits",
				HelpText:     "Number of get hits",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"GetMisses": {
				NameOverride: "",
				Name:         "get_misses",
				HelpText:     "Number of get misses",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"IncrHits": {
				NameOverride: "",
				Name:         "incr_hits",
				HelpText:     "Number of increment hits",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"IncrMisses": {
				NameOverride: "",
				Name:         "incr_misses",
				HelpText:     "Number of increment misses",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemUsed": {
				NameOverride: "",
				Name:         "mem_used",
				HelpText:     "Amount of memory used",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"Misses": {
				NameOverride: "",
				Name:         "misses",
				HelpText:     "Number of misses",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"Ops": {
				NameOverride: "",
				Name:         "ops",
				HelpText:     "Total amount of operations per second to this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			// lol Timestamp.
			"VbActiveEject": {
				NameOverride: "",
				Name:         "vb_active_eject",
				HelpText:     "Number of items per second being ejected to disk from active vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveItmMemory": {
				NameOverride: "",
				Name:         "vb_active_itm_memory",
				HelpText:     "Amount of active user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveMetaDataMemory": {
				NameOverride: "",
				Name:         "vb_active_meta_data_memory",
				HelpText:     "Amount of active item metadata consuming RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveNum": {
				NameOverride: "",
				Name:         "vb_active_num",
				HelpText:     "Number of vBuckets in the active state for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveNumNonresident": {
				NameOverride: "",
				Name:         "vb_active_num_non_resident",
				HelpText:     "Number of non resident vBuckets in the active state for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveOpsCreate": {
				NameOverride: "",
				Name:         "vb_active_ops_create",
				HelpText:     "New items per second being inserted into active vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveOpsUpdate": {
				NameOverride: "",
				Name:         "vb_active_ops_update",
				HelpText:     "Number of items updated on active vBucket per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveQueueAge": {
				NameOverride: "",
				Name:         "vb_active_queue_age",
				HelpText:     "Sum of disk queue item age in milliseconds",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveQueueDrain": {
				NameOverride: "",
				Name:         "vb_active_queue_drain",
				HelpText:     "Number of active items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveQueueFill": {
				NameOverride: "",
				Name:         "vb_active_queue_fill",
				HelpText:     "Number of active items per second being put on the active item disk queue in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveQueueSize": {
				NameOverride: "",
				Name:         "vb_active_queue_size",
				HelpText:     "Number of active items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveQueueItems": {
				NameOverride: "",
				Name:         "vb_active_queue_items",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingCurrItems": {
				NameOverride: "",
				Name:         "vb_pending_curr_items",
				HelpText:     "Number of items in pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingEject": {
				NameOverride: "",
				Name:         "vb_pending_eject",
				HelpText:     "Number of items per second being ejected to disk from pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingItmMemory": {
				NameOverride: "",
				Name:         "vb_pending_itm_memory",
				HelpText:     "Amount of pending user data cached in RAM in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingMetaDataMemory": {
				NameOverride: "",
				Name:         "vb_pending_meta_data_memory",
				HelpText:     "Amount of pending item metadata consuming RAM in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingNum": {
				NameOverride: "",
				Name:         "vb_pending_num",
				HelpText:     "Number of vBuckets in the pending state for this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingNumNonResident": {
				NameOverride: "",
				Name:         "vb_pending_num_non_resident",
				HelpText:     "Number of non resident vBuckets in the pending state for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingOpsCreate": {
				NameOverride: "",
				Name:         "vb_pending_ops_create",
				HelpText:     "New items per second being instead into pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingOpsUpdate": {
				NameOverride: "",
				Name:         "vb_pending_ops_update",
				HelpText:     "Number of items updated on pending vBucket per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingQueueAge": {
				NameOverride: "",
				Name:         "vb_pending_queue_age",
				HelpText:     "Sum of disk pending queue item age in milliseconds",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingQueueDrain": {
				NameOverride: "",
				Name:         "vb_pending_queue_drain",
				HelpText:     "Number of pending items per second being written to disk in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingQueueFill": {
				NameOverride: "",
				Name:         "vb_pending_queue_fill",
				HelpText:     "Number of pending items per second being put on the pending item disk queue in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingQueueSize": {
				NameOverride: "",
				Name:         "vb_pending_queue_size",
				HelpText:     "Number of pending items waiting to be written to disk in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaCurrItems": {
				NameOverride: "",
				Name:         "vb_replica_curr_items",
				HelpText:     "Number of items in replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaEject": {
				NameOverride: "",
				Name:         "vb_replica_eject",
				HelpText:     "Number of items per second being ejected to disk from replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaItmMemory": {
				NameOverride: "",
				Name:         "vb_replica_itm_memory",
				HelpText:     "Amount of replica user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaMetaDataMemory": {
				NameOverride: "",
				Name:         "vb_replica_meta_data_memory",
				HelpText:     "Amount of replica item metadata consuming in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaNum": {
				NameOverride: "",
				Name:         "vb_replica_num",
				HelpText:     "Number of vBuckets in the replica state for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaNumNonResident": {
				NameOverride: "",
				Name:         "vb_replica_num_non_resident",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaOpsCreate": {
				NameOverride: "",
				Name:         "vb_replica_ops_create",
				HelpText:     "New items per second being inserted into replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaOpsUpdate": {
				NameOverride: "",
				Name:         "vb_replica_ops_update",
				HelpText:     "Number of items updated on replica vBucket per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaQueueAge": {
				NameOverride: "",
				Name:         "vb_replica_queue_age",
				HelpText:     "Sum of disk replica queue item age in milliseconds",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaQueueDrain": {
				NameOverride: "",
				Name:         "vb_replica_queue_drain",
				HelpText:     "Number of replica items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaQueueFill": {
				NameOverride: "",
				Name:         "vb_replica_queue_fill",
				HelpText:     "Number of replica items per second being put on the replica item disk queue in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaQueueSize": {
				NameOverride: "",
				Name:         "vb_replica_queue_size",
				HelpText:     "Number of replica items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbTotalQueueAge": {
				NameOverride: "",
				Name:         "vb_total_queue_age",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbAvgActiveQueueAge": {
				NameOverride: "",
				Name:         "vb_avg_active_queue_age",
				HelpText:     "Sum of disk queue item age in milliseconds",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbAvgReplicaQueueAge": {
				NameOverride: "",
				Name:         "vb_avg_replica_queue_age",
				HelpText:     "Average age in seconds of replica items in the replica item queue for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbAvgPendingQueueAge": {
				NameOverride: "",
				Name:         "vb_avg_pending_queue_age",
				HelpText:     "Average age in seconds of pending items in the pending item queue for this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbAvgTotalQueueAge": {
				NameOverride: "",
				Name:         "vb_avg_total_queue_age",
				HelpText:     "Average age in seconds of all items in the disk write queue for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbActiveResidentItemsRatio": {
				NameOverride: "",
				Name:         "vb_active_resident_items_ratio",
				HelpText:     "Percentage of active items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbReplicaResidentItemsRatio": {
				NameOverride: "",
				Name:         "vb_replica_resident_items_ratio",
				HelpText:     "Percentage of active items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"VbPendingResidentItemsRatio": {
				NameOverride: "",
				Name:         "vb_pending_resident_items_ratio",
				HelpText:     "Percentage of items in pending state vbuckets cached in RAM in this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"XdcOps": {
				NameOverride: "",
				Name:         "xdc_ops",
				HelpText:     "Total XDCR operations per second for this bucket",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CPUIdleMs": {
				NameOverride: "",
				Name:         "cpu_idle_ms",
				HelpText:     "CPU idle milliseconds",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CPULocalMs": {
				NameOverride: "",
				Name:         "cpu_local_ms",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"CPUUtilizationRate": {
				NameOverride: "",
				Name:         "cpu_utilization_rate",
				HelpText:     "Percentage of CPU in use across all available cores on this server",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"HibernatedRequests": {
				NameOverride: "",
				Name:         "hibernated_requests",
				HelpText:     "Number of streaming requests on port 8091 now idle",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"HibernatedWaked": {
				NameOverride: "",
				Name:         "hibernated_waked",
				HelpText:     "Rate of streaming request wakeups on port 8091",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemActualFree": {
				NameOverride: "",
				Name:         "mem_actual_free",
				HelpText:     "Amount of RAM available on this server",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemActualUsed": {
				NameOverride: "",
				Name:         "mem_actual_used",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemFree": {
				NameOverride: "",
				Name:         "mem_free",
				HelpText:     "Amount of Memory free",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemTotal": {
				NameOverride: "",
				Name:         "mem_total",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"MemUsedSys": {
				NameOverride: "",
				Name:         "mem_used_sys",
				HelpText:     "",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"RestRequests": {
				NameOverride: "",
				Name:         "rest_requests",
				HelpText:     "Rate of http requests on port 8091",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"SwapTotal": {
				NameOverride: "",
				Name:         "swap_total",
				HelpText:     "Total amount of swap available",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
			"SwapUsed": {
				NameOverride: "",
				Name:         "swap_used",
				HelpText:     "Amount of swap space in use on this server",
				Labels:       []string{BucketLabel, NodeLabel, ClusterLabel},
			},
		},
	}

	return newConfig
}

func bucketStatsCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "BucketStats",
		Namespace: DefaultNamespace,
		Subsystem: "bucketstats",
		Metrics: map[string]MetricInfo{
			"AvgBgWaitTime": {
				Name:         "avg_bg_wait_seconds",
				HelpText:     "Average background fetch time in seconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"AvgActiveTimestampDrift": {
				Name:         "avg_active_timestamp_drift",
				HelpText:     "Average drift (in seconds) per mutation on active vBuckets",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"AvgReplicaTimestampDrift": {
				Name:         "avg_replica_timestamp_drift",
				HelpText:     "Average drift (in seconds) per mutation on replica vBuckets",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"AvgDiskCommitTime": {
				Name:         "avg_disk_commit_time",
				HelpText:     "Average disk commit time in seconds as from disk_update histogram of timings",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"AvgDiskUpdateTime": {
				Name:         "avg_disk_update_time",
				HelpText:     "Average disk update time in microseconds as from disk_update histogram of timings",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"BytesRead": {
				Name:         "read_bytes",
				HelpText:     "Bytes read",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"BytesWritten": {
				Name:         "written_bytes",
				HelpText:     "Bytes written",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CasBadval": {
				Name:         "cas_badval",
				HelpText:     "Compare and Swap bad values",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CasHits": {
				Name:         "cas_hits",
				HelpText:     "Number of operations with a CAS id per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CasMisses": {
				Name:         "cas_misses",
				HelpText:     "Compare and Swap misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CmdGet": {
				Name:         "cmd_get",
				HelpText:     "Number of reads (get operations) per second from this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CmdSet": {
				Name:         "cmd_set",
				HelpText:     "Number of writes (set operations) per second to this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchTotalDiskSize": {
				Name:         "couch_total_disk_size",
				HelpText:     "The total size on disk of all data and view files for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchViewsFragmentation": {
				Name:         "couch_views_fragmentation",
				HelpText:     "How much fragmented data there is to be compacted compared to real data for the view index files in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchViewsOps": {
				Name:         "couch_views_ops",
				HelpText:     "All the view reads for all design documents including scatter gather",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchViewsDataSize": {
				Name:         "couch_views_data_size",
				HelpText:     "The size of active data on for all the indexes in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchViewsActualDiskSize": {
				Name:         "couch_views_actual_disk_size",
				HelpText:     "The size of all active items in all the indexes for this bucket on disk",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchDocsFragmentation": {
				Name:         "couch_docs_fragmentation",
				HelpText:     "How much fragmented data there is to be compacted compared to real data for the data files in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchDocsActualDiskSize": {
				Name:         "couch_docs_actual_disk_size",
				HelpText:     "The size of all data files for this bucket, including the data itself, meta data and temporary files",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchDocsDataSize": {
				Name:         "couch_docs_data_size",
				HelpText:     "The size of active data in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CouchDocsDiskSize": {
				Name:         "couch_docs_disk_size",
				HelpText:     "The size of all data files for this bucket, including the data itself, meta data and temporary files",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CPUIdleMs": {
				Name:         "cpu_idle_ms",
				HelpText:     "CPU idle milliseconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CPULocalMs": {
				Name:         "cpu_local_ms",
				HelpText:     "_cpu_local_ms",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CPUUtilizationRate": {
				Name:         "cpu_utilization_rate",
				HelpText:     "Percentage of CPU in use across all available cores on this server",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CurrConnections": {
				Name:         "curr_connections",
				HelpText:     "Number of connections to this server including connections from external client SDKs, proxies, DCP requests and internal statistic gathering",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CurrItems": {
				Name:         "curr_items",
				HelpText:     "Number of items in active vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"CurrItemsTot": {
				Name:         "curr_items_tot",
				HelpText:     "Total number of items in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DecrHits": {
				Name:         "decr_hits",
				HelpText:     "Decrement hits",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DecrMisses": {
				Name:         "decr_misses",
				HelpText:     "Decrement misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DeleteHits": {
				Name:         "delete_hits",
				HelpText:     "Number of delete operations per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DeleteMisses": {
				Name:         "delete_misses",
				HelpText:     "Delete misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DiskCommitCount": {
				Name:         "disk_commits",
				HelpText:     "Disk commits",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DiskUpdateCount": {
				Name:         "disk_updates",
				HelpText:     "Disk updates",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"DiskWriteQueue": {
				Name:         "disk_write_queue",
				HelpText:     "Number of items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpActiveAheadExceptions": {
				Name:         "ep_active_ahead_exceptions",
				HelpText:     "Total number of ahead exceptions for  all active vBuckets",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpActiveHlcDrift": {
				Name:         "ep_active_hlc_drift",
				HelpText:     "_ep_active_hlc_drift",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpClockCasDriftThresholdExceeded": {
				Name:         "ep_clock_cas_drift_threshold_exceeded",
				HelpText:     "_ep_clock_cas_drift_threshold_exceeded",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpBgFetched": {
				Name:         "ep_bg_fetched",
				HelpText:     "Number of reads per second from disk for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpCacheMissRate": {
				Name:         "ep_cache_miss_rate",
				HelpText:     "Percentage of reads per second to this bucket from disk as opposed to RAM",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iBackoff": {
				Name:         "ep_dcp_2i_backoff",
				HelpText:     "Number of backoffs for indexes DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iCount": {
				Name:         "ep_dcp_2i_connections",
				HelpText:     "Number of indexes DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iItemsRemaining": {
				Name:         "ep_dcp_2i_items_remaining",
				HelpText:     "Number of indexes items remaining to be sent",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iItemsSent": {
				Name:         "ep_dcp_2i_items_sent",
				HelpText:     "Number of indexes items sent",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iProducerCount": {
				Name:         "ep_dcp_2i_producers",
				HelpText:     "Number of indexes producers",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iTotalBacklogSize": {
				Name:         "ep_dcp_2i_total_backlog_size",
				HelpText:     "ep_dcp_2i_total_backlog_size",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcp2iTotalBytes": {
				Name:         "ep_dcp_2i_total_bytes",
				HelpText:     "Number bytes per second being sent for indexes DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherBackoff": {
				Name:         "ep_dcp_other_backoff",
				HelpText:     "Number of backoffs for other DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherCount": {
				Name:         "ep_dcp_others",
				HelpText:     "Number of other DCP connections in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherItemsRemaining": {
				Name:         "ep_dcp_other_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherItemsSent": {
				Name:         "ep_dcp_other_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherProducerCount": {
				Name:         "ep_dcp_other_producers",
				HelpText:     "Number of other senders for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherTotalBacklogSize": {
				Name:         "ep_dcp_other_total_backlog_size",
				HelpText:     "ep_dcp_other_total_backlog_size",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpOtherTotalBytes": {
				Name:         "ep_dcp_other_total_bytes",
				HelpText:     "Number of bytes per second being sent for other DCP connections for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaBackoff": {
				Name:         "ep_dcp_replica_backoff",
				HelpText:     "Number of backoffs for replication DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaCount": {
				Name:         "ep_dcp_replicas",
				HelpText:     "Number of internal replication DCP connections in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaItemsRemaining": {
				Name:         "ep_dcp_replica_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaItemsSent": {
				Name:         "ep_dcp_replica_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaProducerCount": {
				Name:         "ep_dcp_replica_producers",
				HelpText:     "Number of replication senders for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaTotalBacklogSize": {
				Name:         "ep_dcp_replica_total_backlog_size",
				HelpText:     "ep_dcp_replica_total_backlog_size",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpReplicaTotalBytes": {
				Name:         "ep_dcp_replica_total_bytes",
				HelpText:     "Number of bytes per second being sent for replication DCP connections for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsBackoff": {
				Name:         "ep_dcp_views_backoffs",
				HelpText:     "Number of backoffs for views DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsCount": {
				Name:         "ep_dcp_view_connections",
				HelpText:     "Number of views DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsItemsRemaining": {
				Name:         "ep_dcp_views_items_remaining",
				HelpText:     "Number of views items remaining to be sent",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsItemsSent": {
				Name:         "ep_dcp_views_items_sent",
				HelpText:     "Number of views items sent",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsProducerCount": {
				Name:         "ep_dcp_views_producers",
				HelpText:     "Number of views producers",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsTotalBacklogSize": {
				Name:         "ep_dcp_views_total_backlog_size",
				HelpText:     "ep_dcp_views_total_backlog_size",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpViewsTotalBytes": {
				Name:         "ep_dcp_views_total_bytes",
				HelpText:     "Number bytes per second being sent for views DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrBackoff": {
				Name:         "ep_dcp_xdcr_backoff",
				HelpText:     "Number of backoffs for XDCR DCP connections",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrCount": {
				Name:         "ep_dcp_xdcr_connections",
				HelpText:     "Number of internal XDCR DCP connections in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrItemsRemaining": {
				Name:         "ep_dcp_xdcr_items_remaining",
				HelpText:     "Number of items remaining to be sent to consumer in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrItemsSent": {
				Name:         "ep_dcp_xdcr_items_sent",
				HelpText:     "Number of items per second being sent for a producer for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrProducerCount": {
				Name:         "ep_dcp_xdcr_producers",
				HelpText:     "Number of XDCR senders for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrTotalBacklogSize": {
				Name:         "ep_dcp_xdcr_total_backlog_size",
				HelpText:     "ep_dcp_xdcr_total_backlog_size",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDcpXdcrTotalBytes": {
				Name:         "ep_dcp_xdcr_total_bytes",
				HelpText:     "Number of bytes per second being sent for XDCR DCP connections for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDiskqueueDrain": {
				Name:         "ep_diskqueue_drain",
				HelpText:     "Total number of items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDiskqueueFill": {
				Name:         "ep_diskqueue_fill",
				HelpText:     "Total number of items per second being put on the disk queue in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpDiskqueueItems": {
				Name:         "ep_diskqueue_items",
				HelpText:     "Total number of items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpFlusherTodo": {
				Name:         "ep_flusher_todo",
				HelpText:     "Number of items currently being written",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpItemCommitFailed": {
				Name:         "ep_item_commit_failed",
				HelpText:     "Number of times a transaction failed to commit due to storage errors",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpKvSize": {
				Name:         "ep_kv_size",
				HelpText:     "Total amount of user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpMaxSize": {
				Name:         "ep_max_size_bytes",
				HelpText:     "The maximum amount of memory this bucket can use",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpMemHighWat": {
				Name:         "ep_mem_high_wat_bytes",
				HelpText:     "High water mark for auto-evictions",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpMemLowWat": {
				Name:         "ep_mem_low_wat_bytes",
				HelpText:     "Low water mark for auto-evictions",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpMetaDataMemory": {
				Name:         "ep_meta_data_memory",
				HelpText:     "Total amount of item metadata consuming RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumNonResident": {
				Name:         "ep_num_non_resident",
				HelpText:     "Number of non-resident items",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumOpsDelMeta": {
				Name:         "ep_num_ops_del_meta",
				HelpText:     "Number of delete operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumOpsDelRetMeta": {
				Name:         "ep_num_ops_del_ret_meta",
				HelpText:     "Number of delRetMeta operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumOpsGetMeta": {
				Name:         "ep_num_ops_get_meta",
				HelpText:     "Number of metadata read operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumOpsSetMeta": {
				Name:         "ep_num_ops_set_meta",
				HelpText:     "Number of set operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumOpsSetRetMeta": {
				Name:         "ep_num_ops_set_ret_meta",
				HelpText:     "Number of setRetMeta operations per second for this bucket as the target for XDCR",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpNumValueEjects": {
				Name:         "ep_num_value_ejects",
				HelpText:     "Total number of items per second being ejected to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpOomErrors": {
				Name:         "ep_oom_errors",
				HelpText:     "Number of times unrecoverable OOMs happened while processing operations",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpOpsCreate": {
				Name:         "ep_ops_create",
				HelpText:     "Total number of new items being inserted into this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpOpsUpdate": {
				Name:         "ep_ops_update",
				HelpText:     "Number of items updated on disk per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpOverhead": {
				Name:         "ep_overhead",
				HelpText:     "Extra memory used by transient data like persistence queues or checkpoints",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpQueueSize": {
				Name:         "ep_queue_size",
				HelpText:     "Number of items queued for storage",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpResidentItemsRate": {
				Name:         "ep_resident_items_rate",
				HelpText:     "Percentage of all items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpReplicaAheadExceptions": {
				Name:         "ep_replica_ahead_exceptions",
				HelpText:     "Total number of ahead exceptions (when timestamp drift between mutations and local time has exceeded 5000000 μs) per second for all replica vBuckets.",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpReplicaHlcDrift": {
				Name:         "ep_replica_hlc_drift",
				HelpText:     "The sum of the total Absolute Drift, which is the accumulated drift observed by the vBucket. Drift is always accumulated as an absolute value.",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpTmpOomErrors": {
				Name:         "ep_tmp_oom_errors",
				HelpText:     "Number of back-offs sent per second to client SDKs due to OOM situations from this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"EpVbTotal": {
				Name:         "ep_vbuckets",
				HelpText:     "Total number of vBuckets for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"Evictions": {
				Name:         "evictions",
				HelpText:     "Number of evictions",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"GetHits": {
				Name:         "get_hits",
				HelpText:     "Number of get hits",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"GetMisses": {
				Name:         "get_misses",
				HelpText:     "Number of get misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"HibernatedRequests": {
				Name:         "hibernated_requests",
				HelpText:     "Number of streaming requests on port 8091 now idle",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"HibernatedWaked": {
				Name:         "hibernated_waked",
				HelpText:     "Rate of streaming request wakeups on port 8091",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"HitRatio": {
				Name:         "hit_ratio",
				HelpText:     "Hit ratio",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"IncrHits": {
				Name:         "incr_hits",
				HelpText:     "Number of increment hits",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"IncrMisses": {
				Name:         "incr_misses",
				HelpText:     "Number of increment misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemActualFree": {
				Name:         "mem_actual_free",
				HelpText:     "Amount of RAM available on this server",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemActualUsed": {
				Name:         "mem_actual_used_bytes",
				HelpText:     "_mem_actual_used",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemFree": {
				Name:         "mem_free_bytes",
				HelpText:     "Amount of Memory free",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemTotal": {
				Name:         "mem_bytes",
				HelpText:     "Total amount of memory available",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemUsed": {
				Name:         "mem_used_bytes",
				HelpText:     "Amount of memory used",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"MemUsedSys": {
				Name:         "mem_used_sys_bytes",
				HelpText:     "_mem_used_sys",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"Misses": {
				Name:         "misses",
				HelpText:     "Number of misses",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"Ops": {
				Name:         "ops",
				HelpText:     "Total amount of operations per second to this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"RestRequests": {
				Name:         "rest_requests",
				HelpText:     "Rate of http requests on port 8091",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"SwapTotal": {
				Name:         "swap_bytes",
				HelpText:     "Total amount of swap available",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"SwapUsed": {
				Name:         "swap_used_bytes",
				HelpText:     "Amount of swap space in use on this server",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveEject": {
				Name:         "vbuckets_active_eject",
				HelpText:     "Number of items per second being ejected to disk from active vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveItmMemory": {
				Name:         "vbuckets_active_itm_memory",
				HelpText:     "Amount of active user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveMetaDataMemory": {
				Name:         "vbuckets_active_meta_data_memory",
				HelpText:     "Amount of active item metadata consuming RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveNum": {
				Name:         "vbuckets_active_num",
				HelpText:     "Number of vBuckets in the active state for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveNumNonResident": {
				Name:         "vbuckets_active_num_non_resident",
				HelpText:     "Number of non resident vBuckets in the active state for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveOpsCreate": {
				Name:         "vbuckets_active_ops_create",
				HelpText:     "New items per second being inserted into active vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveOpsUpdate": {
				Name:         "vbuckets_active_ops_update",
				HelpText:     "Number of items updated on active vBucket per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveQueueAge": {
				Name:         "vbuckets_active_queue_age",
				HelpText:     "Sum of disk queue item age in milliseconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveQueueDrain": {
				Name:         "vbuckets_active_queue_drain",
				HelpText:     "Number of active items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveQueueFill": {
				Name:         "vbuckets_active_queue_fill",
				HelpText:     "Number of active items per second being put on the active item disk queue in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveQueueSize": {
				Name:         "vbuckets_active_queue_size",
				HelpText:     "Number of active items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbActiveResidentItemsRatio": {
				Name:         "vbuckets_active_resident_items_ratio",
				HelpText:     "Percentage of active items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbAvgActiveQueueAge": {
				Name:         "vbuckets_avg_active_queue_age",
				HelpText:     "Average age in seconds of active items in the active item queue for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbAvgPendingQueueAge": {
				Name:         "vbuckets_avg_pending_queue_age",
				HelpText:     "Average age in seconds of pending items in the pending item queue for this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbAvgReplicaQueueAge": {
				Name:         "vbuckets_avg_replica_queue_age",
				HelpText:     "Average age in seconds of replica items in the replica item queue for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbAvgTotalQueueAge": {
				Name:         "vbuckets_avg_total_queue_age",
				HelpText:     "Average age in seconds of all items in the disk write queue for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingCurrItems": {
				Name:         "vbuckets_pending_curr_items",
				HelpText:     "Number of items in pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingEject": {
				Name:         "vbuckets_pending_eject",
				HelpText:     "Number of items per second being ejected to disk from pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingItmMemory": {
				Name:         "vbuckets_pending_itm_memory",
				HelpText:     "Amount of pending user data cached in RAM in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingMetaDataMemory": {
				Name:         "vbuckets_pending_meta_data_memory",
				HelpText:     "Amount of pending item metadata consuming RAM in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingNum": {
				Name:         "vbuckets_pending_num",
				HelpText:     "Number of vBuckets in the pending state for this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingNumNonResident": {
				Name:         "vbuckets_pending_num_non_resident",
				HelpText:     "Number of non resident vBuckets in the pending state for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingOpsCreate": {
				Name:         "vbuckets_pending_ops_create",
				HelpText:     "New items per second being instead into pending vBuckets in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingOpsUpdate": {
				Name:         "vbuckets_pending_ops_update",
				HelpText:     "Number of items updated on pending vBucket per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingQueueAge": {
				Name:         "vbuckets_pending_queue_age",
				HelpText:     "Sum of disk pending queue item age in milliseconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingQueueDrain": {
				Name:         "vbuckets_pending_queue_drain",
				HelpText:     "Number of pending items per second being written to disk in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingQueueFill": {
				Name:         "vbuckets_pending_queue_fill",
				HelpText:     "Number of pending items per second being put on the pending item disk queue in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingQueueSize": {
				Name:         "vbuckets_pending_queue_size",
				HelpText:     "Number of pending items waiting to be written to disk in this bucket and should be transient during rebalancing",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbPendingResidentItemsRatio": {
				Name:         "vbuckets_pending_resident_items_ratio",
				HelpText:     "Percentage of items in pending state vbuckets cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaCurrItems": {
				Name:         "vbuckets_replica_curr_items",
				HelpText:     "Number of items in replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaEject": {
				Name:         "vbuckets_replica_eject",
				HelpText:     "Number of items per second being ejected to disk from replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaItmMemory": {
				Name:         "vbuckets_replica_itm_memory",
				HelpText:     "Amount of replica user data cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaMetaDataMemory": {
				Name:         "vbuckets_replica_meta_data_memory",
				HelpText:     "Amount of replica item metadata consuming in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaNum": {
				Name:         "vbuckets_replica_num",
				HelpText:     "Number of vBuckets in the replica state for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaNumNonResident": {
				Name:         "vbuckets_replica_num_non_resident",
				HelpText:     "_vb_replica_num_non_resident",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaOpsCreate": {
				Name:         "vbuckets_replica_ops_create",
				HelpText:     "New items per second being inserted into replica vBuckets in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaOpsUpdate": {
				Name:         "vbuckets_replica_ops_update",
				HelpText:     "Number of items updated on replica vBucket per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaQueueAge": {
				Name:         "vbuckets_replica_queue_age",
				HelpText:     "Sum of disk replica queue item age in milliseconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaQueueDrain": {
				Name:         "vbuckets_replica_queue_drain",
				HelpText:     "Number of replica items per second being written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaQueueFill": {
				Name:         "vbuckets_replica_queue_fill",
				HelpText:     "Number of replica items per second being put on the replica item disk queue in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaQueueSize": {
				Name:         "vbuckets_replica_queue_size",
				HelpText:     "Number of replica items waiting to be written to disk in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbReplicaResidentItemsRatio": {
				Name:         "vbuckets_replica_resident_items_ratio",
				HelpText:     "Percentage of replica items cached in RAM in this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"VbTotalQueueAge": {
				Name:         "vbuckets_total_queue_age",
				HelpText:     "Sum of disk queue item age in milliseconds",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
			"XdcOps": {
				Name:         "xdc_ops",
				HelpText:     "Total XDCR operations per second for this bucket",
				Labels:       []string{BucketLabel, ClusterLabel},
				NameOverride: "",
			},
		},
	}

	return newConfig
}

func nodeCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      NodeLabel,
		Namespace: DefaultNamespace,
		Subsystem: NodeLabel,
		Metrics: map[string]MetricInfo{
			"healthy": {
				Name:         "healthy",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Is this node healthy",
				Labels:       []string{ClusterLabel, NodeLabel},
			},
			"systemStatsCPUUtilizationRate": {
				Name:         "systemstats_cpu_utilization_rate",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Percentage of CPU in use across all available cores on this server.",
				Labels:       []string{ClusterLabel, NodeLabel},
			},
			"systemStatsSwapTotal": {
				Name:         "systemstats_swap_total",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Bytes of total swap space available on this server.",
				Labels:       []string{ClusterLabel, NodeLabel},
			},
			"systemStatsSwapUsed": {
				Name:         "systemstats_swap_used",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Bytes of swap space in use on this server.",
				Labels:       []string{ClusterLabel, NodeLabel},
			},
			"systemStatsMemTotal": {
				Name:         "systemstats_mem_total",
				NameOverride: "",
				HelpText:     "Bytes of total memory available on this server.",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"systemStatsMemFree": {
				Name:         "systemstats_mem_free",
				NameOverride: "",
				HelpText:     "Bytes of memory not in use on this server.",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchDocsActualDiskSize": {
				Name:         "interestingstats_couch_docs_actual_disk_size",
				NameOverride: "",
				HelpText:     "The size of all data service files on disk for this bucket, including the data itself, metadata, and temporary files. (measured from couch_docs_actual_disk_size)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchDocsDataSize": {
				Name:         "interestingstats_couch_docs_data_size",
				NameOverride: "",
				HelpText:     "Bytes of active data in this bucket. (measured from couch_docs_data_size)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchViewsActualDiskSize": {
				Name:         "interestingstats_couch_views_actual_disk_size",
				NameOverride: "",
				HelpText:     "Bytes of active items in all the views for this bucket on disk (measured from couch_views_actual_disk_size)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchViewsDataSize": {
				Name:         "interestingstats_couch_views_data_size",
				NameOverride: "",
				HelpText:     "Bytes of active data for all the views in this bucket. (measured from couch_views_data_size)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsMemUsed": {
				Name:         "interestingstats_mem_used",
				NameOverride: "",
				HelpText:     "Total memory used in bytes. (as measured from mem_used)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsOps": {
				Name:         "interestingstats_ops",
				NameOverride: "",
				HelpText:     "Total operations per second (including XDCR) to this bucket. (measured from cmd_get + cmd_set + incr_misses + incr_hits + decr_misses + decr_hits + delete_misses + delete_hits + ep_num_ops_del_meta + ep_num_ops_get_meta + ep_num_ops_set_meta)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCurrItems": {
				Name:         "interestingstats_curr_items",
				NameOverride: "",
				HelpText:     "Current number of unique items in Couchbase",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCurrItemsTot": {
				Name:         "interestingstats_curr_items_tot",
				NameOverride: "",
				HelpText:     "Current number of items in Couchbase including replicas",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsVbReplicaCurrItems": {
				Name:         "interestingstats_vb_replica_curr_items",
				NameOverride: "",
				HelpText:     "Number of items in replica vBuckets in this bucket. (measured from vb_replica_curr_items)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsVbActiveNumNonResident": {
				Name:         "interestingstats_vb_active_number_non_resident",
				NameOverride: "",
				HelpText:     "interestingstats_vb_active_number_non_resident",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchSpatialDiskSize": {
				Name:         "interestingstats_couch_spatial_disk_size",
				NameOverride: "",
				HelpText:     "interestingstats_couch_spatial_disk_size",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCouchSpatialDataSize": {
				Name:         "interestingstats_couch_spatial_data_size",
				NameOverride: "",
				HelpText:     "interestingstats_couch_spatial_data_size",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsCmdGet": {
				Name:         "interestingstats_cmd_get",
				NameOverride: "",
				HelpText:     "Number of reads (get operations) per second from this bucket. (measured from cmd_get)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsGetHits": {
				Name:         "interestingstats_get_hits",
				NameOverride: "",
				HelpText:     "Number of get operations per second for data that this bucket contains. (measured from get_hits)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"interestingStatsEpBgFetched": {
				Name:         "interestingstats_ep_bg_fetched",
				NameOverride: "",
				HelpText:     "Number of reads per second from disk for this bucket. (measured from ep_bg_fetched)",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"uptime": {
				Name:         "uptime",
				HelpText:     "uptime",
				NameOverride: "",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"memoryTotal": {
				Name:         "memory_total",
				HelpText:     "memory_total",
				NameOverride: "",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"memoryFree": {
				Name:         "memory_free",
				HelpText:     "memory_free",
				NameOverride: "",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"mcdMemoryAllocated": {
				Name:         "memcached_memory_allocated",
				NameOverride: "",
				HelpText:     "memcached_memory_allocated",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"mcdMemoryReserved": {
				Name:         "memcached_memory_reserved",
				NameOverride: "",
				HelpText:     "memcached_memory_reserved",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"clusterMembership": {
				Name:         "cluster_membership",
				NameOverride: "",
				HelpText:     "whether or not node is part of the CB cluster",
				Labels:       []string{NodeLabel, ClusterLabel},
				Enabled:      true,
			},
			"ctrFailover": {
				Name:         "failover",
				NameOverride: "",
				HelpText:     "failover",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrFailoverNode": {
				Name:         "failover_node",
				NameOverride: "",
				HelpText:     "failover_node",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrFailoverComplete": {
				Name:         "failover_complete",
				NameOverride: "",
				HelpText:     "failover_complete",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrFailoverIncomplete": {
				Name:         "failover_incomplete",
				NameOverride: "",
				HelpText:     "failover_incomplete",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrRebalanceStart": {
				Name:         "rebalance_start",
				NameOverride: "",
				HelpText:     "rebalance_start",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrRebalanceStop": {
				Name:         "rebalance_stop",
				NameOverride: "",
				HelpText:     "rebalance_stop",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrRebalanceSuccess": {
				Name:         "rebalance_success",
				NameOverride: "",
				HelpText:     "rebalance_success",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrRebalanceFail": {
				Name:         "rebalance_failure",
				NameOverride: "",
				HelpText:     "rebalance_failure",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrGracefulFailoverStart": {
				Name:         "graceful_failover_start",
				NameOverride: "",
				HelpText:     "graceful_failover_start",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrGracefulFailoverSuccess": {
				Name:         "graceful_failover_success",
				NameOverride: "",
				HelpText:     "graceful_failover_success",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
			"ctrGracefulFailoverFail": {
				Name:         "graceful_failover_fail",
				NameOverride: "",
				HelpText:     "graceful_failover_fail",
				Labels:       []string{ClusterLabel},
				Enabled:      true,
			},
		},
	}

	return newConfig
}

func eventingCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "Eventing",
		Namespace: DefaultNamespace,
		Subsystem: "eventing",
		Metrics: map[string]MetricInfo{
			"eventingBucketOpExceptionCount": {
				Name:         "bucket_op_exception_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingCheckpointFailureCount": {
				Name:         "checkpoint_failure_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingDcpBacklog": {
				Name:         "dcp_backlog",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Mutations yet to be processed by the function",
				Labels:       []string{ClusterLabel},
			},
			"eventingFailedCount": {
				Name:         "failed_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Mutations for which the function execution failed",
				Labels:       []string{ClusterLabel},
			},
			"eventingN1QlOpExceptionCount": {
				Name:         "n1ql_op_exception_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of disk bytes read on Analytics node per second",
				Labels:       []string{ClusterLabel},
			},
			"eventingOnDeleteFailure": {
				Name:         "on_delete_failure",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of disk bytes written on Analytics node per second",
				Labels:       []string{ClusterLabel},
			},
			"eventingOnDeleteSuccess": {
				Name:         "on_delete_success",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "System load for Analytics node",
				Labels:       []string{ClusterLabel},
			},
			"eventingOnUpdateFailure": {
				Name:         "on_update_failure",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingOnUpdateSuccess": {
				Name:         "on_update_success",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingProcessedCount": {
				Name:         "processed_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Mutations for which the function has finished processing",
				Labels:       []string{ClusterLabel},
			},
			"eventingTimeoutCount": {
				Name:         "timeout_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Function execution timed-out while processing",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestBucketOpExceptionCount": {
				Name:         "test_bucket_op_exception_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestCheckpointFailureCount": {
				Name:         "test_checkpoint_failure_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestDcpBacklog": {
				Name:         "test_dcp_backlog",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestFailedCount": {
				Name:         "test_failed_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestN1QlOpExceptionCount": {
				Name:         "test_n1ql_op_exception_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestOnDeleteFailure": {
				Name:         "test_on_delete_failure",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestOnDeleteSuccess": {
				Name:         "test_on_delete_success",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestOnUpdateFailure": {
				Name:         "test_on_update_failure",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestOnUpdateSuccess": {
				Name:         "test_on_update_success",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestProcessedCount": {
				Name:         "test_processed_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
			"eventingTestTimeoutCount": {
				Name:         "test_timeout_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "",
				Labels:       []string{ClusterLabel},
			},
		},
	}

	return newConfig
}

func taskCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "Search",
		Namespace: DefaultNamespace,
		Subsystem: "task",
		Metrics: map[string]MetricInfo{
			"rebalance": {
				Name:         "rebalance_progress",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Progress of a rebalance task",
				Labels:       []string{ClusterLabel},
			},
			"rebalancePerNode": {
				Name:         "node_rebalance_progress",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Progress of a rebalance task per node",
				Labels:       []string{NodeLabel, ClusterLabel},
			},
			"compacting": {
				Name:         "compacting_progress",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Progress of a bucket compaction task",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"clusterLogsCollection": {
				Name:         "cluster_logs_collection_progress",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Progress of a cluster logs collection task",
				Labels:       []string{ClusterLabel},
			},
			"xdcrChangesLeft": {
				Name:         "xdcr_changes_left",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of updates still pending replication",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"xdcrDocsChecked": {
				Name:         "xdcr_docs_checked",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of documents checked for changes",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"xdcrDocsWritten": {
				Name:         "xdcr_docs_written",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of documents written to the destination cluster",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"xdcrPaused": {
				Name:         "xdcr_paused",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Is this replication paused",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"xdcrErrors": {
				Name:         "xdcr_errors",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of errors",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"progressDocsTotal": {
				Name:         "docs_total",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "docs_total",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"progressDocsTransferred": {
				Name:         "docs_transferred",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "docs_transferred",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"progressActiveVBucketsLeft": {
				Name:         "active_vbuckets_left",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of Active VBuckets remaining",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
			"progressReplicateVBucketsLeft": {
				Name:         "replica_vbuckets_left",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of Replica VBuckets remaining",
				Labels:       []string{BucketLabel, "target", ClusterLabel},
			},
		},
	}

	return newConfig
}

func searchCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "Search",
		Namespace: DefaultNamespace,
		Subsystem: "fts",
		Metrics: map[string]MetricInfo{
			"FtsCurrBatchesBlockedByHerder": {
				Name:         "curr_batches_blocked_by_herder",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of DCP batches blocked by the FTS throttler due to high memory consumption",
				Labels:       []string{ClusterLabel},
			},
			"FtsNumBytesUsedRAM": {
				Name:         "num_bytes_used_ram",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Amount of RAM used by FTS on this server",
				Labels:       []string{ClusterLabel},
			},
			"FtsTotalQueriesRejectedByHerder": {
				Name:         "total_queries_rejected_by_herder",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of fts queries rejected by the FTS throttler due to high memory consumption",
				Labels:       []string{ClusterLabel},
			},
		},
	}

	return newConfig
}

func indexCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "Index",
		Namespace: DefaultNamespace,
		Subsystem: "index",
		Metrics: map[string]MetricInfo{
			"IndexMemoryQuota": {
				Name:         "memory_quota",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Index Service memory quota",
				Labels:       []string{ClusterLabel},
			},
			"IndexMemoryUsed": {
				Name:         "memory_used",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Index Service memory used",
				Labels:       []string{ClusterLabel},
			},
			"IndexRAMPercent": {
				Name:         "ram_percent",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Percentage of Index RAM quota in use across all indexes on this server.",
				Labels:       []string{ClusterLabel},
			},
			"IndexRemainingRAM": {
				Name:         "remaining_ram",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Bytes of Index RAM quota still available on this server.",
				Labels:       []string{ClusterLabel},
			},
		},
	}

	return newConfig
}

func analyticsCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "Analytics",
		Namespace: DefaultNamespace,
		Subsystem: "cbas",
		Metrics: map[string]MetricInfo{
			"CbasDiskUsed": {
				Name:         "disk_used",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "The total disk size used by Analytics",
				Labels:       []string{ClusterLabel},
			},
			"CbasGcCount": {
				Name:         "gc_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of JVM garbage collections for Analytics node",
				Labels:       []string{ClusterLabel},
			},
			"CbasGcTime": {
				Name:         "gc_time",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "The amount of time in milliseconds spent performing JVM garbage collections for Analytics node",
				Labels:       []string{ClusterLabel},
			},
			"CbasHeapUsed": {
				Name:         "heap_used",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Amount of JVM heap used by Analytics on this server",
				Labels:       []string{ClusterLabel},
			},
			"CbasIoReads": {
				Name:         "io_reads",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of disk bytes read on Analytics node per second",
				Labels:       []string{ClusterLabel},
			},
			"CbasIoWrites": {
				Name:         "io_writes",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of disk bytes written on Analytics node per second",
				Labels:       []string{ClusterLabel},
			},
			"CbasSystemLoadAverage": {
				Name:         "system_load_avg",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "System load for Analytics node",
				Labels:       []string{ClusterLabel},
			},
			"CbasThreadCount": {
				Name:         "thread_count",
				NameOverride: "",
				Enabled:      true,
				HelpText:     "Number of threads for Analytics node",
				Labels:       []string{ClusterLabel},
			},
		},
	}

	return newConfig
}

func bucketInfoCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "BucketInfoCollector",
		Namespace: DefaultNamespace,
		Subsystem: "bucketinfo",
		Metrics: map[string]MetricInfo{
			"dataUsed": {
				Name:         "basic_dataused_bytes",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_dataused",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"diskFetches": {
				Name:         "basic_diskfetches",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_diskfetches",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"diskUsed": {
				Name:         "basic_diskused_bytes",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_diskused",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"itemCount": {
				Name:         "basic_itemcount",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_itemcount",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"memUsed": {
				Name:         "basic_memused_bytes",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_memused",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"opsPerSec": {
				Name:         "basic_opspersec",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_opspersec",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
			"quotaPercentUsed": {
				Name:         "basic_quota_user_percent",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "basic_quotapercentused",
				Labels:       []string{BucketLabel, ClusterLabel},
			},
		},
	}

	return newConfig
}

func queryCollectorDefaultConfig() *CollectorConfig {
	newConfig := &CollectorConfig{
		Name:      "QueryCollector",
		Namespace: DefaultNamespace,
		Subsystem: "query",
		Metrics: map[string]MetricInfo{
			"QueryAvgReqTime": {
				Name:         "avg_req_time",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "average request time",
				Labels:       []string{ClusterLabel},
			},
			"QueryAvgSvcTime": {
				Name:         "avg_svc_time",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "average service time",
				Labels:       []string{ClusterLabel},
			},
			"QueryAvgResponseSize": {
				Name:         "avg_response_size",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "average response size",
				Labels:       []string{ClusterLabel},
			},
			"QueryAvgResultCount": {
				Name:         "avg_result_count",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "average result count",
				Labels:       []string{ClusterLabel},
			},
			"QueryActiveRequests": {
				Name:         "active_requests",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "active number of requests",
				Labels:       []string{ClusterLabel},
			},
			"QueryErrors": {
				Name:         "errors",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of query errors",
				Labels:       []string{ClusterLabel},
			},
			"QueryInvalidRequests": {
				Name:         "invalid_requests",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of invalid requests",
				Labels:       []string{ClusterLabel},
			},
			"QueryQueuedRequests": {
				Name:         "queued_requests",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of queued requests",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequestTime": {
				Name:         "request_time",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "query request time",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequests": {
				Name:         "requests",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of query requests",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequests1000Ms": {
				Name:         "requests_1000ms",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of requests that take longer than 1000 ms per second",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequests250Ms": {
				Name:         "requests_250ms",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of requests that take longer than 250 ms per second",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequests5000Ms": {
				Name:         "requests_5000ms",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of requests that take longer than 5000 ms per second",
				Labels:       []string{ClusterLabel},
			},
			"QueryRequests500Ms": {
				Name:         "requests_500ms",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of requests that take longer than 500 ms per second",
				Labels:       []string{ClusterLabel},
			},
			"QueryResultCount": {
				Name:         "result_count",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "query result count",
				Labels:       []string{ClusterLabel},
			},
			"QueryResultSize": {
				Name:         "result_size",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "query result size",
				Labels:       []string{ClusterLabel},
			},
			"QuerySelects": {
				Name:         "selects",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of queries involving SELECT",
				Labels:       []string{ClusterLabel},
			},
			"QueryServiceTime": {
				Name:         "service_time",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "query service time",
				Labels:       []string{ClusterLabel},
			},
			"QueryWarnings": {
				Name:         "warnings",
				Enabled:      true,
				NameOverride: "",
				HelpText:     "number of query warnings",
				Labels:       []string{ClusterLabel},
			},
		},
	}

	return newConfig
}
