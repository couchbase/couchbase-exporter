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

type eventingCollector struct {
	m MetaCollector

	EventingBucketOpExceptionCount     *prometheus.Desc
	EventingCheckpointFailureCount     *prometheus.Desc
	EventingDcpBacklog                 *prometheus.Desc
	EventingFailedCount                *prometheus.Desc
	EventingN1QlOpExceptionCount       *prometheus.Desc
	EventingOnDeleteFailure            *prometheus.Desc
	EventingOnDeleteSuccess            *prometheus.Desc
	EventingOnUpdateFailure            *prometheus.Desc
	EventingOnUpdateSuccess            *prometheus.Desc
	EventingProcessedCount             *prometheus.Desc
	EventingTimeoutCount               *prometheus.Desc
	EventingTestBucketOpExceptionCount *prometheus.Desc
	EventingTestCheckpointFailureCount *prometheus.Desc
	EventingTestDcpBacklog             *prometheus.Desc
	EventingTestFailedCount            *prometheus.Desc
	EventingTestN1QlOpExceptionCount   *prometheus.Desc
	EventingTestOnDeleteFailure        *prometheus.Desc
	EventingTestOnDeleteSuccess        *prometheus.Desc
	EventingTestOnUpdateFailure        *prometheus.Desc
	EventingTestOnUpdateSuccess        *prometheus.Desc
	EventingTestProcessedCount         *prometheus.Desc
	EventingTestTimeoutCount           *prometheus.Desc
}

func NewEventingCollector(client util.Client) prometheus.Collector {
	const subsystem = "eventing"
	return &eventingCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "up"),
				"Couchbase cluster API is responding",
				nil,
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "scrape_duration_seconds"),
				"Scrape duration in seconds",
				nil,
				nil,
			),
		},
		EventingBucketOpExceptionCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "bucket_op_exception_count"),
			"",
			nil,
			nil,
		),
		EventingCheckpointFailureCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "checkpoint_failure_count"),
			"",
			nil,
			nil,
		),
		EventingDcpBacklog: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "dcp_backlog"),
			"Mutations yet to be processed by the function",
			nil,
			nil,
		),
		EventingFailedCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "failed_count"),
			"Mutations for which the function execution failed",
			nil,
			nil,
		),
		EventingN1QlOpExceptionCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "n1ql_op_exception_count"),
			"Number of disk bytes read on Analytics node per second",
			nil,
			nil,
		),
		EventingOnDeleteFailure: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "on_delete_failure"),
			"Number of disk bytes written on Analytics node per second",
			nil,
			nil,
		),
		EventingOnDeleteSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "on_delete_success"),
			"System load for Analytics node",
			nil,
			nil,
		),
		EventingOnUpdateFailure: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "on_update_failure"),
			"",
			nil,
			nil,
		),
		EventingOnUpdateSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "on_update_success"),
			"",
			nil,
			nil,
		),
		EventingProcessedCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "processed_count"),
			"Mutations for which the function has finished processing",
			nil,
			nil,
		),
		EventingTimeoutCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "timeout_count"),
			"Function execution timed-out while processing",
			nil,
			nil,
		),
		EventingTestBucketOpExceptionCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_bucket_op_exception_count"),
			"The total disk size used by Analytics",
			nil,
			nil,
		),
		EventingTestCheckpointFailureCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_checkpoint_failure_count"),
			"Number of JVM garbage collections for Analytics node",
			nil,
			nil,
		),
		EventingTestDcpBacklog: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_dcp_backlog"),
			"",
			nil,
			nil,
		),
		EventingTestFailedCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_failed_count"),
			"Amount of JVM heap used by Analytics on this server",
			nil,
			nil,
		),
		EventingTestN1QlOpExceptionCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_n1ql_op_exception_count"),
			"Number of disk bytes read on Analytics node per second",
			nil,
			nil,
		),
		EventingTestOnDeleteFailure: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_on_delete_failure"),
			"Number of disk bytes written on Analytics node per second",
			nil,
			nil,
		),
		EventingTestOnDeleteSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_on_delete_success"),
			"",
			nil,
			nil,
		),
		EventingTestOnUpdateFailure: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_on_update_failure"),
			"",
			nil,
			nil,
		),
		EventingTestOnUpdateSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_on_update_success"),
			"",
			nil,
			nil,
		),
		EventingTestProcessedCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_processed_count"),
			"",
			nil,
			nil,
		),
		EventingTestTimeoutCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "test_timeout_count"),
			"",
			nil,
			nil,
		),
	}
}

// Describe all metrics
func (c *eventingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.EventingBucketOpExceptionCount
	ch <- c.EventingCheckpointFailureCount
	ch <- c.EventingDcpBacklog
	ch <- c.EventingFailedCount
	ch <- c.EventingN1QlOpExceptionCount
	ch <- c.EventingOnDeleteFailure
	ch <- c.EventingOnDeleteSuccess
	ch <- c.EventingOnUpdateFailure
	ch <- c.EventingOnUpdateSuccess
	ch <- c.EventingProcessedCount
	ch <- c.EventingTimeoutCount

	ch <- c.EventingTestBucketOpExceptionCount
	ch <- c.EventingTestCheckpointFailureCount
	ch <- c.EventingTestDcpBacklog
	ch <- c.EventingTestFailedCount
	ch <- c.EventingTestN1QlOpExceptionCount
	ch <- c.EventingTestOnDeleteFailure
	ch <- c.EventingTestOnDeleteSuccess
	ch <- c.EventingTestOnUpdateFailure
	ch <- c.EventingTestOnUpdateSuccess
	ch <- c.EventingTestProcessedCount
	ch <- c.EventingTestTimeoutCount
}

// Collect all metrics
func (c *eventingCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting query metrics...")

	ev, err := c.m.client.Eventing()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.Error("failed to scrape eventing stats")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.EventingBucketOpExceptionCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingBucketOpExceptionCount))
	ch <- prometheus.MustNewConstMetric(c.EventingCheckpointFailureCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingCheckpointFailureCount))
	ch <- prometheus.MustNewConstMetric(c.EventingDcpBacklog, prometheus.GaugeValue, last(ev.Op.Samples.EventingDcpBacklog))
	ch <- prometheus.MustNewConstMetric(c.EventingFailedCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingFailedCount))
	ch <- prometheus.MustNewConstMetric(c.EventingN1QlOpExceptionCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingN1QlOpExceptionCount))
	ch <- prometheus.MustNewConstMetric(c.EventingOnDeleteFailure, prometheus.GaugeValue, last(ev.Op.Samples.EventingOnDeleteFailure))
	ch <- prometheus.MustNewConstMetric(c.EventingOnDeleteSuccess, prometheus.GaugeValue, last(ev.Op.Samples.EventingOnDeleteSuccess))
	ch <- prometheus.MustNewConstMetric(c.EventingOnUpdateFailure, prometheus.GaugeValue, last(ev.Op.Samples.EventingOnUpdateFailure))
	ch <- prometheus.MustNewConstMetric(c.EventingOnUpdateSuccess, prometheus.GaugeValue, last(ev.Op.Samples.EventingOnUpdateSuccess))
	ch <- prometheus.MustNewConstMetric(c.EventingProcessedCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingProcessedCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTimeoutCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTimeoutCount))

	ch <- prometheus.MustNewConstMetric(c.EventingTestBucketOpExceptionCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestBucketOpExceptionCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTestCheckpointFailureCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestCheckpointFailureCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTestDcpBacklog, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestDcpBacklog))
	ch <- prometheus.MustNewConstMetric(c.EventingTestFailedCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestFailedCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTestN1QlOpExceptionCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestN1QlOpExceptionCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTestOnDeleteFailure, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestOnDeleteFailure))
	ch <- prometheus.MustNewConstMetric(c.EventingTestOnDeleteSuccess, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestOnDeleteSuccess))
	ch <- prometheus.MustNewConstMetric(c.EventingTestOnUpdateFailure, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestOnUpdateFailure))
	ch <- prometheus.MustNewConstMetric(c.EventingTestOnUpdateSuccess, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestOnUpdateSuccess))
	ch <- prometheus.MustNewConstMetric(c.EventingTestProcessedCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestProcessedCount))
	ch <- prometheus.MustNewConstMetric(c.EventingTestTimeoutCount, prometheus.GaugeValue, last(ev.Op.Samples.EventingTestTimeoutCount))

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
