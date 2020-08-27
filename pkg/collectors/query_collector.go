//  Copyright (c) 2020 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package collectors

import (
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

type queryCollector struct {
	m MetaCollector

	QueryAvgReqTime      *prometheus.Desc
	QueryAvgSvcTime      *prometheus.Desc
	QueryAvgResponseSize *prometheus.Desc
	QueryAvgResultCount  *prometheus.Desc
	QueryActiveRequests  *prometheus.Desc
	QueryErrors          *prometheus.Desc
	QueryInvalidRequests *prometheus.Desc
	QueryQueuedRequests  *prometheus.Desc
	QueryRequestTime     *prometheus.Desc
	QueryRequests        *prometheus.Desc
	QueryRequests1000Ms  *prometheus.Desc
	QueryRequests250Ms   *prometheus.Desc
	QueryRequests5000Ms  *prometheus.Desc
	QueryRequests500Ms   *prometheus.Desc
	QueryResultCount     *prometheus.Desc
	QueryResultSize      *prometheus.Desc
	QuerySelects         *prometheus.Desc
	QueryServiceTime     *prometheus.Desc
	QueryWarnings        *prometheus.Desc
}

func NewQueryCollector(client util.Client) prometheus.Collector {
	const subsystem = "query"
	return &queryCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "up"),
				"Couchbase cluster API is responding",
				[]string{"cluster"},
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "scrape_duration_seconds"),
				"Scrape duration in seconds",
				[]string{"cluster"},
				nil,
			),
		},
		QueryAvgReqTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "avg_req_time"),
			"average request time",
			[]string{"cluster"},
			nil,
		),
		QueryAvgSvcTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "avg_svc_time"),
			"average service time",
			[]string{"cluster"},
			nil,
		),
		QueryAvgResponseSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "avg_response_size"),
			"average response size",
			[]string{"cluster"},
			nil,
		),
		QueryAvgResultCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "avg_result_count"),
			"average request time",
			[]string{"cluster"},
			nil,
		),
		QueryActiveRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "avg_requests"),
			"average number of requests",
			[]string{"cluster"},
			nil,
		),
		QueryErrors: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "errors"),
			"number of query errors",
			[]string{"cluster"},
			nil,
		),
		QueryInvalidRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "invalid_requests"),
			"number of invalid requests",
			[]string{"cluster"},
			nil,
		),
		QueryQueuedRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "queued_requests"),
			"number of queued requests",
			[]string{"cluster"},
			nil,
		),
		QueryRequestTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "request_time"),
			"query request time",
			[]string{"cluster"},
			nil,
		),
		QueryRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "requests"),
			"number of query requests",
			[]string{"cluster"},
			nil,
		),
		QueryRequests1000Ms: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "requests_1000ms"),
			"number of requests that take longer than 1000 ms per second",
			[]string{"cluster"},
			nil,
		),
		QueryRequests250Ms: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "requests_250ms"),
			"number of requests that take longer than 250 ms per second",
			[]string{"cluster"},
			nil,
		),
		QueryRequests5000Ms: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "requests_5000ms"),
			"number of requests that take longer than 5000 ms per second",
			[]string{"cluster"},
			nil,
		),
		QueryRequests500Ms: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "requests_500ms"),
			"number of requests that take longer than 500 ms per second",
			[]string{"cluster"},
			nil,
		),
		QueryResultCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "result_count"),
			"query result count",
			[]string{"cluster"},
			nil,
		),
		QueryResultSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "result_size"),
			"query result size",
			[]string{"cluster"},
			nil,
		),
		QuerySelects: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "selects"),
			"number of queries involving SELECT",
			[]string{"cluster"},
			nil,
		),
		QueryServiceTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "service_time"),
			"query service time",
			[]string{"cluster"},
			nil,
		),
		QueryWarnings: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "warnings"),
			"number of query warnings",
			[]string{"cluster"},
			nil,
		),
	}
}

// Describe all metrics
func (c *queryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.QueryAvgReqTime
	ch <- c.QueryAvgSvcTime
	ch <- c.QueryAvgResponseSize
	ch <- c.QueryAvgResultCount
	ch <- c.QueryActiveRequests
	ch <- c.QueryErrors
	ch <- c.QueryInvalidRequests
	ch <- c.QueryQueuedRequests
	ch <- c.QueryRequestTime
	ch <- c.QueryRequests
	ch <- c.QueryRequests1000Ms
	ch <- c.QueryRequests250Ms
	ch <- c.QueryRequests5000Ms
	ch <- c.QueryRequests500Ms
	ch <- c.QueryResultCount
	ch <- c.QueryResultSize
	ch <- c.QuerySelects
	ch <- c.QueryServiceTime
	ch <- c.QueryWarnings
}

// Collect all metrics
func (c *queryCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting query metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("%s", err)
		return
	}

	queryStats, err := c.m.client.Query()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("failed to scrape query stats")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.QueryAvgReqTime, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryAvgReqTime), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryAvgSvcTime, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryAvgSvcTime), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryAvgResponseSize, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryAvgResponseSize), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryAvgResultCount, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryAvgResultCount), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryErrors, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryErrors), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryInvalidRequests, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryInvalidRequests), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryQueuedRequests, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryQueuedRequests), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequestTime, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequestTime), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequests, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequests), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequests1000Ms, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequests1000Ms), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequests250Ms, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequests250Ms), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequests5000Ms, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequests5000Ms), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryRequests500Ms, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryRequests500Ms), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryResultCount, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryResultCount), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryResultSize, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryResultSize), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QuerySelects, prometheus.GaugeValue, last(queryStats.Op.Samples.QuerySelects), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryServiceTime, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryServiceTime), clusterName)
	ch <- prometheus.MustNewConstMetric(c.QueryWarnings, prometheus.GaugeValue, last(queryStats.Op.Samples.QueryWarnings), clusterName)

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}
