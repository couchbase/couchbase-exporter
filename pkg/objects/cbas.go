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
	// Keys for samples structures.
	CbasDiskUsed          = "cbas_disk_used"
	CbasGcCount           = "cbas_gc_count"
	CbasGcTime            = "cbas_gc_time"
	CbasHeapUsed          = "cbas_heap_used"
	CbasIoReads           = "cbas_io_reads"
	CbasIoWrites          = "cbas_io_writes"
	CbasSystemLoadAverage = "cbas_system_load_average"
	CbasThreadCount       = "cbas_thread_count"
)

type Analytics struct {
	Op struct {
		Samples      map[string][]float64 `json:"samples"`
		SamplesCount int                  `json:"samplesCount"`
		IsPersistent bool                 `json:"isPersistent"`
		LastTStamp   int64                `json:"lastTStamp"`
		Interval     int                  `json:"interval"`
	} `json:"op"`
}
