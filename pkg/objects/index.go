//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

const (
	// Samples Keys.
	IndexMemoryQuota  = "index_memory_quota"
	IndexMemoryUsed   = "index_memory_used"
	IndexRAMPercent   = "index_ram_percent"
	IndexRemainingRAM = "index_remaining_ram"

	// these are const keys for the Indexer Stats.
	IndexDocsIndexed          = "num_docs_indexed"
	IndexItemsCount           = "items_count"
	IndexFragPercent          = "frag_percent"
	IndexNumDocsPendingQueued = "num_docs_pending_queued"
	IndexNumRequests          = "num_requests"
	IndexCacheMisses          = "cache_misses"
	IndexCacheHits            = "cache_hits"
	IndexCacheHitPercent      = "cache_hit_percent"
	IndexNumRowsReturned      = "num_rows_returned"
	IndexResidentPercent      = "resident_percent"
	IndexAvgScanLatency       = "avg_scan_latency"
)

type Index struct {
	Op struct {
		Samples map[string][]float64 `json:"samples"`
	} `json:"op"`
}
