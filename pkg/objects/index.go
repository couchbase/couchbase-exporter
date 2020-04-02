//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

type Index struct {
	Op struct {
		Samples struct {
			IndexMemoryQuota  []float64 `json:"index_memory_quota"`
			IndexMemoryUsed   []float64 `json:"index_memory_used"`
			IndexRAMPercent   []float64 `json:"index_ram_percent"`
			IndexRemainingRAM []float64 `json:"index_remaining_ram"`
		} `json:"samples"`
	} `json:"op"`
}
