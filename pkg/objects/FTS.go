//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

type FTS struct {
	Op struct {
		Samples struct {
			FtsCurrBatchesBlockedByHerder   []float64 `json:"fts_curr_batches_blocked_by_herder"`
			FtsNumBytesUsedRAM              []float64 `json:"fts_num_bytes_used_ram"`
			FtsTotalQueriesRejectedByHerder []float64 `json:"fts_total_queries_rejected_by_herder"`
		} `json:"samples"`
	} `json:"op"`
}
