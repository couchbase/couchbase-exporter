//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

type Analytics struct {
	Op struct {
		Samples struct {
			CbasDiskUsed          []float64 `json:"cbas_disk_used"`
			CbasGcCount           []float64 `json:"cbas_gc_count"`
			CbasGcTime            []float64 `json:"cbas_gc_time"`
			CbasHeapUsed          []float64 `json:"cbas_heap_used"`
			CbasIoReads           []float64 `json:"cbas_io_reads"`
			CbasIoWrites          []float64 `json:"cbas_io_writes"`
			CbasSystemLoadAverage []float64 `json:"cbas_system_load_average"`
			CbasThreadCount       []float64 `json:"cbas_thread_count"`
		} `json:"samples"`
		SamplesCount int   `json:"samplesCount"`
		IsPersistent bool  `json:"isPersistent"`
		LastTStamp   int64 `json:"lastTStamp"`
		Interval     int   `json:"interval"`
	} `json:"op"`
}
