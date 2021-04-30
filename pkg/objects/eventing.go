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
	EventingCheckpointFailureCount     = "eventing/checkpoint_failure_count"
	EventingBucketOpExceptionCount     = "eventing/bucket_op_exception_count"
	EventingDcpBacklog                 = "eventing/dcp_backlog"
	EventingFailedCount                = "eventing/failed_count"
	EventingN1QlOpExceptionCount       = "eventing/n1ql_op_exception_count"
	EventingOnDeleteFailure            = "eventing/on_delete_failure"
	EventingOnDeleteSuccess            = "eventing/on_delete_success"
	EventingOnUpdateFailure            = "eventing/on_update_failure"
	EventingOnUpdateSuccess            = "eventing/on_update_success"
	EventingProcessedCount             = "eventing/processed_count"
	EventingTestBucketOpExceptionCount = "eventing/test/bucket_op_exception_count"
	EventingTestCheckpointFailureCount = "eventing/test/checkpoint_failure_count"
	EventingTestDcpBacklog             = "eventing/test/dcp_backlog"
	EventingTestFailedCount            = "eventing/test/failed_count"
	EventingTestN1QlOpExceptionCount   = "eventing/test/n1ql_op_exception_count"
	EventingTestOnDeleteFailure        = "eventing/test/on_delete_failure"
	EventingTestOnDeleteSuccess        = "eventing/test/on_delete_success"
	EventingTestOnUpdateFailure        = "eventing/test/on_update_failure"
	EventingTestOnUpdateSuccess        = "eventing/test/on_update_success"
	EventingTestProcessedCount         = "eventing/test/processed_count"
	EventingTestTimeoutCount           = "eventing/test/timeout_count"
	EventingTimeoutCount               = "eventing/timeout_count"
)

type Eventing struct {
	Op struct {
		Samples map[string][]float64 `json:"samples"`
	} `json:"op"`
}
