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
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"time"
)

type cbasCollector struct {
	m MetaCollector

	CbasDiskUsed          *prometheus.Desc
	CbasGcCount           *prometheus.Desc
	CbasGcTime            *prometheus.Desc
	CbasHeapUsed          *prometheus.Desc
	CbasIoReads           *prometheus.Desc
	CbasIoWrites          *prometheus.Desc
	CbasSystemLoadAverage *prometheus.Desc
	CbasThreadCount       *prometheus.Desc
}

func NewCbasCollector(client util.Client) prometheus.Collector {
	const subsystem = "cbas"
	return &cbasCollector{
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
		CbasDiskUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "disk_used"),
			"The total disk size used by Analytics",
			nil,
			nil,
		),
		CbasGcCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "gc_count"),
			"Number of JVM garbage collections for Analytics node",
			nil,
			nil,
		),
		CbasGcTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "gc_time"),
			"The amount of time in milliseconds spent performing JVM garbage collections for Analytics node",
			nil,
			nil,
		),
		CbasHeapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "heap_used"),
			"Amount of JVM heap used by Analytics on this server",
			nil,
			nil,
		),
		CbasIoReads: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "io_reads"),
			"Number of disk bytes read on Analytics node per second",
			nil,
			nil,
		),
		CbasIoWrites: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "io_writes"),
			"Number of disk bytes written on Analytics node per second",
			nil,
			nil,
		),
		CbasSystemLoadAverage: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "system_load_avg"),
			"System load for Analytics node",
			nil,
			nil,
		),
		CbasThreadCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "thread_count"),
			"Number of threads for Analytics node",
			nil,
			nil,
		),
	}
}

// Describe all metrics
func (c *cbasCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.CbasDiskUsed
	ch <- c.CbasGcCount
	ch <- c.CbasGcTime
	ch <- c.CbasHeapUsed
	ch <- c.CbasIoReads
	ch <- c.CbasIoWrites
	ch <- c.CbasSystemLoadAverage
	ch <- c.CbasThreadCount
}

// Collect all metrics
func (c *cbasCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting query metrics...")

	cbas, err := c.m.client.Cbas()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape Analytics stats")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.CbasDiskUsed, prometheus.GaugeValue, last(cbas.Op.Samples.CbasDiskUsed))
	ch <- prometheus.MustNewConstMetric(c.CbasGcCount, prometheus.GaugeValue, last(cbas.Op.Samples.CbasGcCount))
	ch <- prometheus.MustNewConstMetric(c.CbasGcTime, prometheus.GaugeValue, last(cbas.Op.Samples.CbasGcTime))
	ch <- prometheus.MustNewConstMetric(c.CbasHeapUsed, prometheus.GaugeValue, last(cbas.Op.Samples.CbasHeapUsed))
	ch <- prometheus.MustNewConstMetric(c.CbasIoReads, prometheus.GaugeValue, last(cbas.Op.Samples.CbasIoReads))
	ch <- prometheus.MustNewConstMetric(c.CbasIoWrites, prometheus.GaugeValue, last(cbas.Op.Samples.CbasIoWrites))
	ch <- prometheus.MustNewConstMetric(c.CbasSystemLoadAverage, prometheus.GaugeValue, last(cbas.Op.Samples.CbasSystemLoadAverage))
	ch <- prometheus.MustNewConstMetric(c.CbasThreadCount, prometheus.GaugeValue, last(cbas.Op.Samples.CbasThreadCount))

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
