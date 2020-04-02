//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

type Query struct {
	Op struct {
		Samples struct {
			QueryAvgReqTime      []float64 `json:"query_avg_req_time"`
			QueryAvgSvcTime      []float64 `json:"query_avg_svc_time"`
			QueryAvgResponseSize []float64 `json:"query_avg_response_size"`
			QueryAvgResultCount  []float64 `json:"query_avg_result_count"`
			QueryActiveRequests  []float64 `json:"query_active_requests"`
			QueryErrors          []float64 `json:"query_errors"`
			QueryInvalidRequests []float64 `json:"query_invalid_requests"`
			QueryQueuedRequests  []float64 `json:"query_queued_requests"`
			QueryRequestTime     []float64 `json:"query_request_time"`
			QueryRequests        []float64 `json:"query_requests"`
			QueryRequests1000Ms  []float64 `json:"query_requests_1000ms"`
			QueryRequests250Ms   []float64 `json:"query_requests_250ms"`
			QueryRequests5000Ms  []float64 `json:"query_requests_5000ms"`
			QueryRequests500Ms   []float64 `json:"query_requests_500ms"`
			QueryResultCount     []float64 `json:"query_result_count"`
			QueryResultSize      []float64 `json:"query_result_size"`
			QuerySelects         []float64 `json:"query_selects"`
			QueryServiceTime     []float64 `json:"query_service_time"`
			QueryWarnings        []float64 `json:"query_warnings"`
		} `json:"samples"`
	} `json:"op"`
}
