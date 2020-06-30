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
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	username = "Administrator"
	passwd   = "password"
	logger   = logf.Log.WithName("metrics")
	client   = http.Client{}
)

const (
	subsystem = "pernodebucket"
)

var (
	AvgDiskUpdateTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "avg_disk_update_time",
		"Average disk update time in microseconds as from disk_update histogram of timings",
		nil,
	},
		[]string{"bucket", "node"},
	)
	AvgDiskCommitTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "avg_disk_commit_time",
		"Average disk commit time in seconds as from disk_update histogram of timings",
		nil,
	},
		[]string{"bucket", "node"},
	)
	AvgBgWaitTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "avg_bg_wait_seconds",
		" ",
		nil,
	},
		[]string{"bucket", "node"},
	)
	AvgActiveTimestampDrift = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "avg_active_timestamp_drift",
		"  ",
		nil,
	},
		[]string{"bucket", "node"},
	)
	AvgReplicaTimestampDrift = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "avg_replica_timestamp_drift",
		"  ",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CouchTotalDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_total_disk_size",
		"The total size on disk of all data and view files for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchDocsFragmentation = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_docs_fragmentation",
		"How much fragmented data there is to be compacted compared to real data for the data files in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchViewsFragmentation = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_views_fragmentation",
		"How much fragmented data there is to be compacted compared to real data for the view index files in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchDocsActualDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_docs_actual_disk_size",
		"The size of all data files for this bucket, including the data itself, meta data and temporary files",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchDocsDataSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_docs_data_size",
		"The size of active data in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchDocsDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_docs_disk_size",
		"The size of all data files for this bucket, including the data itself, meta data and temporary files",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchSpatialDataSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_spatial_data_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchSpatialDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_spatial_disk_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchSpatialOps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_spatial_ops",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchViewsActualDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_views_actual_disk_size",
		"The size of all active items in all the indexes for this bucket on disk",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchViewsDataSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_views_data_size",
		"The size of active data on for all the indexes in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchViewsDiskSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_views_disk_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	CouchViewsOps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "couch_views_ops",
		"All the view reads for all design documents including scatter gather",
		nil,
	},
		[]string{"bucket", "node"},
	)

	HitRatio = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "hit_ratio",
		"Hit ratio",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpCacheMissRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_cache_miss_rate",
		"Percentage of reads per second to this bucket from disk as opposed to RAM",
		nil,
	},
		[]string{"bucket", "node"},
	)
	EpResidentItemsRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_resident_items_rate",
		"Percentage of all items cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsIndexesCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	EpDcpViewsIndexesItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_items_remaining",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	EpDcpViewsIndexesProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_producer_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)
	EpDcpViewsIndexesTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsIndexesItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_items_sent",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsIndexesTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_total_bytes",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsIndexesBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_indexes_backoff",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	BgWaitCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "bg_wait_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	BgWaitTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "bg_wait_total",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	BytesRead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "bytes_read",
		"Bytes Read",
		nil,
	},
		[]string{"bucket", "node"},
	)

	BytesWritten = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "bytes_written",
		"Bytes written",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CasBadVal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cas_bad_val",
		"Compare and Swap bad values",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CasHits = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cas_hits",
		"Number of operations with a CAS id per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CasMisses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cas_misses",
		"Compare and Swap misses",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CmdGet = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cmd_get",
		"Number of reads (get operations) per second from this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CmdSet = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cmd_set",
		"Number of writes (set operations) per second to this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CurrConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "curr_connections",
		"Number of connections to this server including connections from external client SDKs, proxies, DCP requests and internal statistic gathering",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CurrItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "curr_items",
		"Number of items in active vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CurrItemsTot = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "curr_items_tot",
		"Total number of items in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DecrHits = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "decr_hits",
		"Decrement hits",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DecrMisses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "decr_misses",
		"Decrement misses",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DeleteHits = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "delete_hits",
		"Number of delete operations per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DeleteMisses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "delete_misses",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DiskCommitCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "disk_commit_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DiskCommitTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "disk_commit_total",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DiskUpdateCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "disk_update_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DiskUpdateTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "disk_update_total",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	DiskWriteQueue = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "disk_write_queue",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpActiveAheadExceptions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_active_ahead_exceptions",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpActiveHlcDrift = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_active_hlc_drift",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpActiveHlcDriftCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_active_hlc_drift_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpBgFetched = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_bg_fetched",
		"Number of reads per second from disk for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpClockCasDriftTheresholExceeded = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_clock_cas_drift_threshold_exceeded",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDataReadFailed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_data_read_failed",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDataWriteFailed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_data_write_failed",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_backoff",
		"Number of backoffs for indexes DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_count",
		"Number of indexes DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_items_remaining",
		"Number of indexes items remaining to be sent",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_items_sent",
		"Number of indexes items sent",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_producers",
		"Number of indexes producers",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcp2iTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_2i_total_bytes",
		"Number of bytes per second being sent for indexes DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_backoff",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_items_remaining",
		"Number of items remaining to be sent to consumer in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_items_sent",
		"Number of items per second being sent for a producer for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_producer_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_cbas_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpCbasTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_total_bytes",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_backoff",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_items_remaining",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_items_sent",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_producer_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpFtsTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_fts_total_bytes",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_backoff",
		"Number of backoffs for other DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_count",
		"Number of other DCP connections in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_items_remaining",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_items_sent",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_producer_count",
		"Number of other senders for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpOtherTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_other_total_bytes",
		"Number of bytes per second being sent for other DCP connections for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_backoff",
		"Number of backoffs for replication DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_count",
		"Number of internal replication DCP connections in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_items_remaining",
		"Number of items remaining to be sent to consumer in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_items_sent",
		"Number of items per second being sent for a producer for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_producer_count",
		"Number of replication senders for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpReplicaTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_replica_total_bytes",
		"Number of bytes per second being sent for replication DCP connections for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_backoff",
		"Number of backoffs for views DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_count",
		"Number of views DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_items_remaining",
		"Number of views items remaining to be sent",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_items_sent",
		"Number of views items sent",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_producer_count",
		"Number of views producers",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpViewsTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_views_total_bytes",
		"Number bytes per second being sent for views DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrBackoff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_backoff",
		"Number of backoffs for XDCR DCP connections",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_count",
		"Number of internal XDCR DCP connections in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrItemsRemaining = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_items_remaining",
		"Number of items remaining to be sent to consumer in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrItemsSent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_items_sent",
		"Number of items per second being sent for a producer for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrProducerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_producer_count",
		"Number of XDCR senders for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrTotalBacklogSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_total_backlog_size",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDcpXdcrTotalBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_dcp_xdcr_total_bytes",
		"Number of bytes per second being sent for XDCR DCP connections for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDiskqueueDrain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_diskqueue_drain",
		"Total number of items per second being written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDiskqueueFill = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_diskqueue_fill",
		"Total number of items per second being put on the disk queue in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpDiskqueueItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_diskqueue_items",
		"Total number of items waiting to be written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpFlusherTodo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_flusher_todo",
		"Number of items currently being written",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpItemCommitFailed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_item_commit_failed",
		"Number of times a transaction failed to commit due to storage errors",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpKvSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_kv_size",
		"Total amount of user data cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpMaxSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_max_size",
		"The maximum amount of memory this bucket can use",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpMemHighWat = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_mem_high_wat",
		"High water mark for auto-evictions",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpMemLowWat = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_mem_low_wat",
		"Low water mark for auto-evictions",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpMetaDataMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_meta_data_memory",
		"Total amount of item metadata consuming RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumNonResident = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_non_resident",
		"Number of non-resident items",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumOpsDelMeta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_ops_del_meta",
		"Number of delete operations per second for this bucket as the target for XDCR",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumOpsDelRetMeta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_ops_del_ret_meta",
		"Number of delRetMeta operations per second for this bucket as the target for XDCR",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumOpsGetMeta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_ops_get_meta",
		"Number of metadata read operations per second for this bucket as the target for XDCR",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumOpsSetMeta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_ops_set_meta",
		"Number of set operations per second for this bucket as the target for XDCR",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumOpsSetRetMeta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_ops_set_ret_meta",
		"Number of setRetMeta operations per second for this bucket as the target for XDCR",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpNumValueEjects = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_num_value_ejects",
		"Total number of items per second being ejected to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpOomErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_oom_errors",
		"Number of times unrecoverable OOMs happened while processing operations",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpOpsCreate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_ops_create",
		"Total number of new items being inserted into this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpOpsUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_ops_update",
		"Number of items updated on disk per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpOverhead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_overhead",
		"Extra memory used by transient data like persistence queues or checkpoints",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_queue_size",
		"Number of items queued for storage",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpReplicaAheadExceptions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_replica_ahead_exceptions",
		"Percentage of all items cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpReplicaHlcDrift = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_replica_hlc_drift",
		"The sum of the total Absolute Drift, which is the accumulated drift observed by the vBucket. Drift is always accumulated as an absolute value.",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpReplicaHlcDriftCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_replica_hlc_drift_count",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpTmpOomErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_tmp_oom_errors",
		"Number of back-offs sent per second to client SDKs due to OOM situations from this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	EpVbTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ep_vb_total",
		"Total number of vBuckets for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	Evictions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "evictions",
		"Number of evictions",
		nil,
	},
		[]string{"bucket", "node"},
	)

	GetHits = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "get_hits",
		"Number of get hits",
		nil,
	},
		[]string{"bucket", "node"},
	)

	GetMisses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "get_misses",
		"Number of get misses",
		nil,
	},
		[]string{"bucket", "node"},
	)

	IncrHits = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "incr_hits",
		"Number of increment hits",
		nil,
	},
		[]string{"bucket", "node"},
	)

	IncrMisses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "incr_misses",
		"Number of increment misses",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_used",
		"Amount of memory used",
		nil,
	},
		[]string{"bucket", "node"},
	)

	Misses = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "misses",
		"Number of misses",
		nil,
	},
		[]string{"bucket", "node"},
	)

	Ops = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "ops",
		"Total amount of operations per second to this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	// lol Timestamp

	VbActiveEject = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_eject",
		"Number of items per second being ejected to disk from active vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveItmMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_itm_memory",
		"Amount of active user data cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveMetaDataMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_meta_data_memory",
		"Amount of active item metadata consuming RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveNum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_num",
		"Number of vBuckets in the active state for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveNumNonresident = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_num_non_resident",
		"Number of non resident vBuckets in the active state for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveOpsCreate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_ops_create",
		"New items per second being inserted into active vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveOpsUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_ops_update",
		"Number of items updated on active vBucket per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_queue_age",
		"Sum of disk queue item age in milliseconds",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveQueueDrain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_queue_drain",
		"Number of active items per second being written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveQueueFill = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_queue_fill",
		"Number of active items per second being put on the active item disk queue in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_queue_size",
		"Number of active items waiting to be written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveQueueItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_queue_items",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingCurrItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_curr_items",
		"Number of items in pending vBuckets in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingEject = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_eject",
		"Number of items per second being ejected to disk from pending vBuckets in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingItmMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_itm_memory",
		"Amount of pending user data cached in RAM in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingMetaDataMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_meta_data_memory",
		"Amount of pending item metadata consuming RAM in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingNum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_num",
		"Number of vBuckets in the pending state for this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingNumNonResident = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_num_non_resident",
		"Number of non resident vBuckets in the pending state for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingOpsCreate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_ops_create",
		"New items per second being instead into pending vBuckets in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingOpsUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_ops_update",
		"Number of items updated on pending vBucket per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_queue_age",
		"Sum of disk pending queue item age in milliseconds",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingQueueDrain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_queue_drain",
		"Number of pending items per second being written to disk in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingQueueFill = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_queue_fill",
		"Number of pending items per second being put on the pending item disk queue in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_queue_size",
		"Number of pending items waiting to be written to disk in this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaCurrItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_curr_items",
		"Number of items in replica vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaEject = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_eject",
		"Number of items per second being ejected to disk from replica vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaItmMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_itm_memory",
		"Amount of replica user data cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaMetaDataMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_meta_data_memory",
		"Amount of replica item metadata consuming in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaNum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_num",
		"Number of vBuckets in the replica state for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaNumNonResident = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_num_non_resident",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaOpsCreate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_ops_create",
		"New items per second being inserted into replica vBuckets in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaOpsUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_ops_update",
		"Number of items updated on replica vBucket per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_queue_age",
		"Sum of disk replica queue item age in milliseconds",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaQueueDrain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_queue_drain",
		"Number of replica items per second being written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaQueueFill = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_queue_fill",
		"Number of replica items per second being put on the replica item disk queue in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_queue_size",
		"Number of replica items waiting to be written to disk in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbTotalQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_total_queue_age",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbAvgActiveQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_avg_active_queue_age",
		"Sum of disk queue item age in milliseconds",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbAvgReplicaQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_avg_replica_queue_age",
		"Average age in seconds of replica items in the replica item queue for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbAvgPendingQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_avg_pending_queue_age",
		"Average age in seconds of pending items in the pending item queue for this bucket and should be transient during rebalancing",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbAvgTotalQueueAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_avg_total_queue_age",
		"Average age in seconds of all items in the disk write queue for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbActiveResidentItemsRatio = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_active_resident_items_ratio",
		"Percentage of active items cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbReplicaResidentItemsRatio = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_replica_resident_items_ratio",
		"Percentage of active items cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	VbPendingResidentItemsRatio = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "vb_pending_resident_items_ratio",
		"Percentage of items in pending state vbuckets cached in RAM in this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	XdcOps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "xdc_ops",
		"Total XDCR operations per second for this bucket",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CpuIdleMs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cpu_idle_ms",
		"CPU idle milliseconds",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CpuLocalMs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cpu_local_ms",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	CpuUtilizationRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "cpu_utilization_rate",
		"Percentage of CPU in use across all available cores on this server",
		nil,
	},
		[]string{"bucket", "node"},
	)

	HibernatedRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "hibernated_requests",
		"Number of streaming requests on port 8091 now idle",
		nil,
	},
		[]string{"bucket", "node"},
	)

	HibernatedWaked = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "hibernated_waked",
		"Rate of streaming request wakeups on port 8091",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemActualFree = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_actual_free",
		"Amount of RAM available on this server",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemActualUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_actual_used",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemFree = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_free",
		"Amount of Memory free",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_total",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	MemUsedSys = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "mem_used_sys",
		"",
		nil,
	},
		[]string{"bucket", "node"},
	)

	RestRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "rest_requests",
		"Rate of http requests on port 8091",
		nil,
	},
		[]string{"bucket", "node"},
	)

	SwapTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "swap_total",
		"Total amount of swap available",
		nil,
	},
		[]string{"bucket", "node"},
	)

	SwapUsed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		FQ_NAMESPACE + subsystem, "", "swap_used",
		"Amount of swap space in use on this server",
		nil,
	},
		[]string{"bucket", "node"},
	)
)

func strToFloatArr(floatsStr string) []float64 {
	floatsStrArr := strings.Split(floatsStr, " ")
	var floatsArr []float64

	for _, f := range floatsStrArr {
		i, err := strconv.ParseFloat(f, 64)
		if err == nil {
			floatsArr = append(floatsArr, i)
		}
	}

	return floatsArr
}

func setGaugeVec(vec prometheus.GaugeVec, stats []float64, bucketName, nodeName string) {
	if len(stats) > 0 {
		vec.WithLabelValues(bucketName, nodeName).Set(stats[len(stats)-1])
	}
}

func getClusterBalancedStatus(c util.Client) (bool, error) {
	node, err := c.Nodes()
	if err != nil {
		logger.Error(err, "bad")
		return false, err
	}

	return node.Counters.RebalanceSuccess > 0 || (node.Balanced && node.RebalanceStatus == "none"), nil
}

func getCurrentNode(c util.Client) (string, error) {
	nodes, err := c.Nodes()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve nodes: %s", err)
	}

	for _, node := range nodes.Nodes {
		if node.ThisNode { // "ThisNode" is a boolean value indicating that it is the current node
			return node.Hostname, nil // hostname seems to work? just don't use for single node setups
		}
	}

	return "", err
}

func getPerNodeBucketStats(client util.Client, bucketName, nodeName string) map[string]interface{} {
	url := getSpecificNodeBucketStatsURL(client, bucketName, nodeName)

	var bucketStats objects.PerNodeBucketStats
	err := client.Get(url, &bucketStats)
	if err != nil {
		logger.Error(err, "unable to GET PerNodeBucketStats")
	}

	return bucketStats.Op.Samples
}

// /pools/default/buckets/<bucket-name>/nodes/<node-name>/stats
func getSpecificNodeBucketStatsURL(client util.Client, bucket, node string) string {
	servers, err := client.Servers(bucket)
	if err != nil {
		logger.Error(err, "unable to retrieve Servers")
	}

	correctURI := ""
	for _, server := range servers.Servers {
		if server.Hostname == node {
			correctURI = server.Stats["uri"]
		}
	}

	return correctURI
}

func collectPerNodeBucketMetrics(client util.Client, node string, refreshTime int) {

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	outerErr := util.Retry(ctx, 20*time.Second, 10, func() (bool, error) {

		rebalanced, err := getClusterBalancedStatus(client)
		if err != nil {
			logger.Error(err, "Unable to get rebalance status")
		}

		if !rebalanced {
			logger.Info("Waiting for Rebalance... retrying...")
			return false, err
		} else {
			go func() {
				for {
					buckets, err := client.Buckets()
					if err != nil {
						logger.Error(err, "Unable to get buckets")
					}

					for _, bucket := range buckets {
						logger.Info("Collecting per-node bucket stats", "node", node, "bucket", bucket.Name)

						samples := getPerNodeBucketStats(client, bucket.Name, node)

						setGaugeVec(*AvgDiskUpdateTime, strToFloatArr(fmt.Sprint(samples["avg_disk_update_time"])), bucket.Name, node)
						setGaugeVec(*AvgDiskCommitTime, strToFloatArr(fmt.Sprint(samples["avg_disk_commit_time"])), bucket.Name, node)
						setGaugeVec(*AvgBgWaitTime, strToFloatArr(fmt.Sprint(samples["avg_bg_wait_seconds"])), bucket.Name, node)
						setGaugeVec(*AvgActiveTimestampDrift, strToFloatArr(fmt.Sprint(samples["avg_active_timestamp_drift"])), bucket.Name, node)
						setGaugeVec(*AvgReplicaTimestampDrift, strToFloatArr(fmt.Sprint(samples["avg_replica_timestamp_drift"])), bucket.Name, node)

						setGaugeVec(*CouchTotalDiskSize, strToFloatArr(fmt.Sprint(samples["couch_total_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchDocsFragmentation, strToFloatArr(fmt.Sprint(samples["couch_docs_fragmentation"])), bucket.Name, node)
						setGaugeVec(*CouchViewsFragmentation, strToFloatArr(fmt.Sprint(samples["couch_views_fragmentation"])), bucket.Name, node)
						setGaugeVec(*CouchDocsActualDiskSize, strToFloatArr(fmt.Sprint(samples["couch_docs_actual_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchDocsDataSize, strToFloatArr(fmt.Sprint(samples["couch_docs_data_size"])), bucket.Name, node)
						setGaugeVec(*CouchDocsDiskSize, strToFloatArr(fmt.Sprint(samples["couch_docs_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchSpatialDataSize, strToFloatArr(fmt.Sprint(samples["couch_docs_spatial_data_size"])), bucket.Name, node)
						setGaugeVec(*CouchSpatialDiskSize, strToFloatArr(fmt.Sprint(samples["couch_docs_spatial_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchSpatialOps, strToFloatArr(fmt.Sprint(samples["couch_spatial_ops"])), bucket.Name, node)
						setGaugeVec(*CouchViewsActualDiskSize, strToFloatArr(fmt.Sprint(samples["couch_views_actual_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchViewsDataSize, strToFloatArr(fmt.Sprint(samples["couch_views_data_size"])), bucket.Name, node)
						setGaugeVec(*CouchViewsDiskSize, strToFloatArr(fmt.Sprint(samples["couch_views_disk_size"])), bucket.Name, node)
						setGaugeVec(*CouchViewsOps, strToFloatArr(fmt.Sprint(samples["couch_views_ops"])), bucket.Name, node)

						setGaugeVec(*EpCacheMissRate, strToFloatArr(fmt.Sprint(samples["ep_cache_miss_rate"])), bucket.Name, node)
						setGaugeVec(*EpResidentItemsRate, strToFloatArr(fmt.Sprint(samples["ep_resident_items_rate"])), bucket.Name, node)

						setGaugeVec(*EpActiveAheadExceptions, strToFloatArr(fmt.Sprint(samples["ep_active_ahead_exceptions"])), bucket.Name, node)
						setGaugeVec(*EpActiveHlcDrift, strToFloatArr(fmt.Sprint(samples["ep_active_hlc_drift"])), bucket.Name, node)
						setGaugeVec(*EpActiveHlcDriftCount, strToFloatArr(fmt.Sprint(samples["ep_active_hlc_drift_count"])), bucket.Name, node)
						setGaugeVec(*EpBgFetched, strToFloatArr(fmt.Sprint(samples["ep_bg_fetched"])), bucket.Name, node)
						setGaugeVec(*EpClockCasDriftTheresholExceeded, strToFloatArr(fmt.Sprint(samples["ep_clock_cas_drift_threshold_exceeded"])), bucket.Name, node)
						setGaugeVec(*EpDataReadFailed, strToFloatArr(fmt.Sprint(samples["ep_data_read_failed"])), bucket.Name, node)
						setGaugeVec(*EpDataWriteFailed, strToFloatArr(fmt.Sprint(samples["ep_data_write_failed"])), bucket.Name, node)

						setGaugeVec(*EpDcp2iBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_count"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_producers"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcp2iTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_2i_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpCbasBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_items_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_items_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpCbasTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_cbas_items_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpFtsBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpFtsTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_fts_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpOtherBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpOtherTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_other_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpReplicaBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpReplicaTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_replica_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpViewsBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpViewsIndexesBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpViewsIndexesTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_views_indexes_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDcpXdcrBackoff, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_backoff"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrItemsRemaining, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_items_remaining"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrItemsSent, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_items_sent"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrProducerCount, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_producer_count"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrTotalBacklogSize, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_total_backlog_size"])), bucket.Name, node)
						setGaugeVec(*EpDcpXdcrTotalBytes, strToFloatArr(fmt.Sprint(samples["ep_dcp_xdcr_total_bytes"])), bucket.Name, node)

						setGaugeVec(*EpDiskqueueDrain, strToFloatArr(fmt.Sprint(samples["ep_diskqueue_drain"])), bucket.Name, node)
						setGaugeVec(*EpDiskqueueFill, strToFloatArr(fmt.Sprint(samples["ep_diskqueue_fill"])), bucket.Name, node)
						setGaugeVec(*EpDiskqueueItems, strToFloatArr(fmt.Sprint(samples["ep_diskqueue_items"])), bucket.Name, node)

						setGaugeVec(*EpFlusherTodo, strToFloatArr(fmt.Sprint(samples["ep_flusher_todo"])), bucket.Name, node)
						setGaugeVec(*EpItemCommitFailed, strToFloatArr(fmt.Sprint(samples["ep_item_commit_failed"])), bucket.Name, node)
						setGaugeVec(*EpKvSize, strToFloatArr(fmt.Sprint(samples["ep_kv_size"])), bucket.Name, node)
						setGaugeVec(*EpMaxSize, strToFloatArr(fmt.Sprint(samples["ep_max_size"])), bucket.Name, node)
						setGaugeVec(*EpMemHighWat, strToFloatArr(fmt.Sprint(samples["ep_mem_high_wat"])), bucket.Name, node)
						setGaugeVec(*EpMemLowWat, strToFloatArr(fmt.Sprint(samples["ep_mem_low_wat"])), bucket.Name, node)
						setGaugeVec(*EpMetaDataMemory, strToFloatArr(fmt.Sprint(samples["ep_meta_data_memory"])), bucket.Name, node)

						setGaugeVec(*EpNumNonResident, strToFloatArr(fmt.Sprint(samples["ep_num_non_resident"])), bucket.Name, node)
						setGaugeVec(*EpNumOpsDelMeta, strToFloatArr(fmt.Sprint(samples["ep_num_ops_del_meta"])), bucket.Name, node)
						setGaugeVec(*EpNumOpsDelRetMeta, strToFloatArr(fmt.Sprint(samples["ep_num_ops_del_ret_meta"])), bucket.Name, node)
						setGaugeVec(*EpNumOpsGetMeta, strToFloatArr(fmt.Sprint(samples["ep_num_ops_get_meta"])), bucket.Name, node)
						setGaugeVec(*EpNumOpsSetMeta, strToFloatArr(fmt.Sprint(samples["ep_num_ops_set_meta"])), bucket.Name, node)
						setGaugeVec(*EpNumOpsSetRetMeta, strToFloatArr(fmt.Sprint(samples["ep_num_ops_set_ret_meta"])), bucket.Name, node)
						setGaugeVec(*EpNumValueEjects, strToFloatArr(fmt.Sprint(samples["ep_num_value_ejects"])), bucket.Name, node)

						setGaugeVec(*EpOomErrors, strToFloatArr(fmt.Sprint(samples["ep_oom_errors"])), bucket.Name, node)
						setGaugeVec(*EpOpsCreate, strToFloatArr(fmt.Sprint(samples["ep_ops_create"])), bucket.Name, node)
						setGaugeVec(*EpOpsUpdate, strToFloatArr(fmt.Sprint(samples["ep_ops_update"])), bucket.Name, node)
						setGaugeVec(*EpOverhead, strToFloatArr(fmt.Sprint(samples["ep_overhead"])), bucket.Name, node)
						setGaugeVec(*EpQueueSize, strToFloatArr(fmt.Sprint(samples["ep_queue_size"])), bucket.Name, node)

						setGaugeVec(*EpReplicaAheadExceptions, strToFloatArr(fmt.Sprint(samples["ep_replica_ahead_exceptions"])), bucket.Name, node)
						setGaugeVec(*EpReplicaHlcDrift, strToFloatArr(fmt.Sprint(samples["ep_replica_hlc_drift"])), bucket.Name, node)
						setGaugeVec(*EpReplicaHlcDriftCount, strToFloatArr(fmt.Sprint(samples["ep_replica_hlc_drift_count"])), bucket.Name, node)
						setGaugeVec(*EpTmpOomErrors, strToFloatArr(fmt.Sprint(samples["ep_tmp_oom_errors"])), bucket.Name, node)
						setGaugeVec(*EpVbTotal, strToFloatArr(fmt.Sprint(samples["ep_vb_total"])), bucket.Name, node)

						setGaugeVec(*VbAvgActiveQueueAge, strToFloatArr(fmt.Sprint(samples["vb_avg_active_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbAvgReplicaQueueAge, strToFloatArr(fmt.Sprint(samples["vb_avg_replica_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbAvgPendingQueueAge, strToFloatArr(fmt.Sprint(samples["vb_avg_pending_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbAvgTotalQueueAge, strToFloatArr(fmt.Sprint(samples["vb_avg_total_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbActiveResidentItemsRatio, strToFloatArr(fmt.Sprint(samples["vb_active_resident_items_ratio"])), bucket.Name, node)
						setGaugeVec(*VbReplicaResidentItemsRatio, strToFloatArr(fmt.Sprint(samples["vb_replica_resident_items_ratio"])), bucket.Name, node)
						setGaugeVec(*VbPendingResidentItemsRatio, strToFloatArr(fmt.Sprint(samples["vb_pending_resident_items_ratio"])), bucket.Name, node)

						setGaugeVec(*VbActiveEject, strToFloatArr(fmt.Sprint(samples["vb_active_eject"])), bucket.Name, node)
						setGaugeVec(*VbActiveItmMemory, strToFloatArr(fmt.Sprint(samples["vb_active_itm_memory"])), bucket.Name, node)
						setGaugeVec(*VbActiveMetaDataMemory, strToFloatArr(fmt.Sprint(samples["vb_active_meta_data_memory"])), bucket.Name, node)
						setGaugeVec(*VbActiveNum, strToFloatArr(fmt.Sprint(samples["vb_active_num"])), bucket.Name, node)
						setGaugeVec(*VbActiveNumNonresident, strToFloatArr(fmt.Sprint(samples["vb_active_num_non_resident"])), bucket.Name, node)
						setGaugeVec(*VbActiveOpsCreate, strToFloatArr(fmt.Sprint(samples["vb_active_ops_create"])), bucket.Name, node)
						setGaugeVec(*VbActiveOpsUpdate, strToFloatArr(fmt.Sprint(samples["vb_active_ops_update"])), bucket.Name, node)
						setGaugeVec(*VbActiveQueueAge, strToFloatArr(fmt.Sprint(samples["vb_active_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbActiveQueueDrain, strToFloatArr(fmt.Sprint(samples["vb_active_queue_drain"])), bucket.Name, node)
						setGaugeVec(*VbActiveQueueFill, strToFloatArr(fmt.Sprint(samples["vb_active_queue_fill"])), bucket.Name, node)
						setGaugeVec(*VbActiveQueueSize, strToFloatArr(fmt.Sprint(samples["vb_active_queue_size"])), bucket.Name, node)
						setGaugeVec(*VbActiveQueueItems, strToFloatArr(fmt.Sprint(samples["vb_active_queue_items"])), bucket.Name, node)

						setGaugeVec(*VbPendingCurrItems, strToFloatArr(fmt.Sprint(samples["vb_pending_curr_items"])), bucket.Name, node)
						setGaugeVec(*VbPendingEject, strToFloatArr(fmt.Sprint(samples["vb_pending_eject"])), bucket.Name, node)
						setGaugeVec(*VbPendingItmMemory, strToFloatArr(fmt.Sprint(samples["vb_pending_itm_memory"])), bucket.Name, node)
						setGaugeVec(*VbPendingMetaDataMemory, strToFloatArr(fmt.Sprint(samples["vb_pending_meta_data_memory"])), bucket.Name, node)
						setGaugeVec(*VbPendingNum, strToFloatArr(fmt.Sprint(samples["vb_pending_num"])), bucket.Name, node)
						setGaugeVec(*VbPendingNumNonResident, strToFloatArr(fmt.Sprint(samples["vb_pending_num_non_resident"])), bucket.Name, node)
						setGaugeVec(*VbPendingOpsCreate, strToFloatArr(fmt.Sprint(samples["vb_pending_ops_create"])), bucket.Name, node)
						setGaugeVec(*VbPendingOpsUpdate, strToFloatArr(fmt.Sprint(samples["vb_pending_ops_update"])), bucket.Name, node)
						setGaugeVec(*VbPendingQueueAge, strToFloatArr(fmt.Sprint(samples["vb_pending_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbPendingQueueDrain, strToFloatArr(fmt.Sprint(samples["vb_pending_queue_drain"])), bucket.Name, node)
						setGaugeVec(*VbPendingQueueFill, strToFloatArr(fmt.Sprint(samples["vb_pending_queue_fill"])), bucket.Name, node)
						setGaugeVec(*VbPendingQueueSize, strToFloatArr(fmt.Sprint(samples["vb_pending_queue_size"])), bucket.Name, node)

						setGaugeVec(*VbReplicaCurrItems, strToFloatArr(fmt.Sprint(samples["vb_replica_curr_items"])), bucket.Name, node)
						setGaugeVec(*VbReplicaEject, strToFloatArr(fmt.Sprint(samples["vb_replica_eject"])), bucket.Name, node)
						setGaugeVec(*VbReplicaItmMemory, strToFloatArr(fmt.Sprint(samples["vb_replica_itm_memory"])), bucket.Name, node)
						setGaugeVec(*VbReplicaMetaDataMemory, strToFloatArr(fmt.Sprint(samples["vb_replica_meta_data_memory"])), bucket.Name, node)
						setGaugeVec(*VbReplicaNum, strToFloatArr(fmt.Sprint(samples["vb_replica_num"])), bucket.Name, node)
						setGaugeVec(*VbReplicaNumNonResident, strToFloatArr(fmt.Sprint(samples["vb_replica_num_non_resident"])), bucket.Name, node)
						setGaugeVec(*VbReplicaOpsCreate, strToFloatArr(fmt.Sprint(samples["vb_replica_ops_create"])), bucket.Name, node)
						setGaugeVec(*VbReplicaOpsUpdate, strToFloatArr(fmt.Sprint(samples["vb_replica_ops_update"])), bucket.Name, node)
						setGaugeVec(*VbReplicaQueueAge, strToFloatArr(fmt.Sprint(samples["vb_replica_queue_age"])), bucket.Name, node)
						setGaugeVec(*VbReplicaQueueDrain, strToFloatArr(fmt.Sprint(samples["vb_replica_queue_drain"])), bucket.Name, node)
						setGaugeVec(*VbReplicaQueueFill, strToFloatArr(fmt.Sprint(samples["vb_replica_queue_fill"])), bucket.Name, node)
						setGaugeVec(*VbReplicaQueueSize, strToFloatArr(fmt.Sprint(samples["vb_replica_queue_size"])), bucket.Name, node)

						setGaugeVec(*VbTotalQueueAge, strToFloatArr(fmt.Sprint(samples["vb_total_queue_age"])), bucket.Name, node)
						setGaugeVec(*HibernatedRequests, strToFloatArr(fmt.Sprint(samples["hibernated_requests"])), bucket.Name, node)
						setGaugeVec(*HibernatedRequests, strToFloatArr(fmt.Sprint(samples["hibernated_waked"])), bucket.Name, node)
						setGaugeVec(*XdcOps, strToFloatArr(fmt.Sprint(samples["xdc_ops"])), bucket.Name, node)
						setGaugeVec(*CpuIdleMs, strToFloatArr(fmt.Sprint(samples["cpu_idle_ms"])), bucket.Name, node)
						setGaugeVec(*CpuLocalMs, strToFloatArr(fmt.Sprint(samples["cpu_local_ms"])), bucket.Name, node)
						setGaugeVec(*CpuUtilizationRate, strToFloatArr(fmt.Sprint(samples["cpu_utilization_rate"])), bucket.Name, node)

						setGaugeVec(*BgWaitCount, strToFloatArr(fmt.Sprint(samples["bg_wait_count"])), bucket.Name, node)
						setGaugeVec(*BgWaitTotal, strToFloatArr(fmt.Sprint(samples["bg_wait_total"])), bucket.Name, node)
						setGaugeVec(*BytesRead, strToFloatArr(fmt.Sprint(samples["bytes_read"])), bucket.Name, node)
						setGaugeVec(*BytesWritten, strToFloatArr(fmt.Sprint(samples["bytes_written"])), bucket.Name, node)
						setGaugeVec(*CasBadVal, strToFloatArr(fmt.Sprint(samples["cas_bad_val"])), bucket.Name, node)
						setGaugeVec(*CasHits, strToFloatArr(fmt.Sprint(samples["cas_hits"])), bucket.Name, node)
						setGaugeVec(*CasMisses, strToFloatArr(fmt.Sprint(samples["cas_misses"])), bucket.Name, node)
						setGaugeVec(*CmdGet, strToFloatArr(fmt.Sprint(samples["cmd_get"])), bucket.Name, node)
						setGaugeVec(*CmdSet, strToFloatArr(fmt.Sprint(samples["cmd_set"])), bucket.Name, node)
						setGaugeVec(*HitRatio, strToFloatArr(fmt.Sprint(samples["hit_ratio"])), bucket.Name, node)

						setGaugeVec(*CurrConnections, strToFloatArr(fmt.Sprint(samples["curr_connections"])), bucket.Name, node)
						setGaugeVec(*CurrItems, strToFloatArr(fmt.Sprint(samples["curr_items"])), bucket.Name, node)
						setGaugeVec(*CurrItemsTot, strToFloatArr(fmt.Sprint(samples["curr_items_tot"])), bucket.Name, node)

						setGaugeVec(*DecrHits, strToFloatArr(fmt.Sprint(samples["decr_hits"])), bucket.Name, node)
						setGaugeVec(*DecrMisses, strToFloatArr(fmt.Sprint(samples["decr_misses"])), bucket.Name, node)
						setGaugeVec(*DeleteHits, strToFloatArr(fmt.Sprint(samples["delete_hits"])), bucket.Name, node)
						setGaugeVec(*DeleteMisses, strToFloatArr(fmt.Sprint(samples["delete_misses"])), bucket.Name, node)

						setGaugeVec(*DiskCommitCount, strToFloatArr(fmt.Sprint(samples["disk_commit_count"])), bucket.Name, node)
						setGaugeVec(*DiskCommitTotal, strToFloatArr(fmt.Sprint(samples["disk_commit_total"])), bucket.Name, node)
						setGaugeVec(*DiskUpdateCount, strToFloatArr(fmt.Sprint(samples["disk_update_count"])), bucket.Name, node)
						setGaugeVec(*DiskUpdateTotal, strToFloatArr(fmt.Sprint(samples["disk_update_total"])), bucket.Name, node)
						setGaugeVec(*DiskWriteQueue, strToFloatArr(fmt.Sprint(samples["disk_write_queue"])), bucket.Name, node)

						setGaugeVec(*Evictions, strToFloatArr(fmt.Sprint(samples["evictions"])), bucket.Name, node)
						setGaugeVec(*GetHits, strToFloatArr(fmt.Sprint(samples["get_hits"])), bucket.Name, node)
						setGaugeVec(*GetMisses, strToFloatArr(fmt.Sprint(samples["get_misses"])), bucket.Name, node)
						setGaugeVec(*IncrHits, strToFloatArr(fmt.Sprint(samples["incr_hits"])), bucket.Name, node)
						setGaugeVec(*IncrMisses, strToFloatArr(fmt.Sprint(samples["incr_misses"])), bucket.Name, node)
						setGaugeVec(*Misses, strToFloatArr(fmt.Sprint(samples["misses"])), bucket.Name, node)
						setGaugeVec(*Ops, strToFloatArr(fmt.Sprint(samples["ops"])), bucket.Name, node)

						setGaugeVec(*MemActualFree, strToFloatArr(fmt.Sprint(samples["mem_actual_free"])), bucket.Name, node)
						setGaugeVec(*MemActualUsed, strToFloatArr(fmt.Sprint(samples["mem_actual_used"])), bucket.Name, node)
						setGaugeVec(*MemFree, strToFloatArr(fmt.Sprint(samples["mem_free"])), bucket.Name, node)
						setGaugeVec(*MemUsed, strToFloatArr(fmt.Sprint(samples["mem_used"])), bucket.Name, node)
						setGaugeVec(*MemTotal, strToFloatArr(fmt.Sprint(samples["mem_total"])), bucket.Name, node)
						setGaugeVec(*MemUsedSys, strToFloatArr(fmt.Sprint(samples["mem_used_sys"])), bucket.Name, node)
						setGaugeVec(*RestRequests, strToFloatArr(fmt.Sprint(samples["rest_requests"])), bucket.Name, node)
						setGaugeVec(*SwapTotal, strToFloatArr(fmt.Sprint(samples["swap_total"])), bucket.Name, node)
						setGaugeVec(*SwapUsed, strToFloatArr(fmt.Sprint(samples["swap_used"])), bucket.Name, node)

					}
					time.Sleep(time.Second * time.Duration(refreshTime))
				}
			}()
			logger.Info("Per Node Bucket Stats Go Thread executed successfully")
			return true, nil
		}
	})
	if outerErr != nil {
		logger.Error(outerErr, "Getting Per Node Bucket Stats failed")
	}
}

func RunPerNodeBucketStatsCollection(client util.Client, refreshTime int) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	outerErr := util.Retry(ctx, 20*time.Second, 8, func() (bool, error) {
		if currNode, err := getCurrentNode(client); err != nil {
			log.Error("could not get current node, will retry. %s", err)
			return false, err
		} else {
			collectPerNodeBucketMetrics(client, currNode, refreshTime)
		}
		return true, nil
	})
	if outerErr != nil {
		fmt.Println("getting default stats failed")
		fmt.Println(outerErr)
	}
}
