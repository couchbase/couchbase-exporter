//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

type Eventing struct {
	Op struct {
		Samples struct {
			EventingBucketOpExceptionCount     []float64 `json:"eventing/bucket_op_exception_count"`
			EventingCheckpointFailureCount     []float64 `json:"eventing/checkpoint_failure_count"`
			EventingDcpBacklog                 []float64 `json:"eventing/dcp_backlog"`
			EventingFailedCount                []float64 `json:"eventing/failed_count"`
			EventingN1QlOpExceptionCount       []float64 `json:"eventing/n1ql_op_exception_count"`
			EventingOnDeleteFailure            []float64 `json:"eventing/on_delete_failure"`
			EventingOnDeleteSuccess            []float64 `json:"eventing/on_delete_success"`
			EventingOnUpdateFailure            []float64 `json:"eventing/on_update_failure"`
			EventingOnUpdateSuccess            []float64 `json:"eventing/on_update_success"`
			EventingProcessedCount             []float64 `json:"eventing/processed_count"`
			EventingTestBucketOpExceptionCount []float64 `json:"eventing/test/bucket_op_exception_count"`
			EventingTestCheckpointFailureCount []float64 `json:"eventing/test/checkpoint_failure_count"`
			EventingTestDcpBacklog             []float64 `json:"eventing/test/dcp_backlog"`
			EventingTestFailedCount            []float64 `json:"eventing/test/failed_count"`
			EventingTestN1QlOpExceptionCount   []float64 `json:"eventing/test/n1ql_op_exception_count"`
			EventingTestOnDeleteFailure        []float64 `json:"eventing/test/on_delete_failure"`
			EventingTestOnDeleteSuccess        []float64 `json:"eventing/test/on_delete_success"`
			EventingTestOnUpdateFailure        []float64 `json:"eventing/test/on_update_failure"`
			EventingTestOnUpdateSuccess        []float64 `json:"eventing/test/on_update_success"`
			EventingTestProcessedCount         []float64 `json:"eventing/test/processed_count"`
			EventingTestTimeoutCount           []float64 `json:"eventing/test/timeout_count"`
			EventingTimeoutCount               []float64 `json:"eventing/timeout_count"`
		} `json:"samples"`
	} `json:"op"`
}
