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

type indexCollector struct {
	m MetaCollector

	IndexMemoryQuota  *prometheus.Desc
	IndexMemoryUsed   *prometheus.Desc
	IndexRAMPercent   *prometheus.Desc
	IndexRemainingRAM *prometheus.Desc
}

func NewIndexCollector(client util.Client) prometheus.Collector {
	const subsystem = "index"
	return &indexCollector{
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
		IndexMemoryQuota: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memory_quota"),
			"Index Service memory quota",
			nil,
			nil,
		),
		IndexMemoryUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "memory_used"),
			"Index Service memory used ",
			nil,
			nil,
		),
		IndexRAMPercent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "ram_percent"),
			"Percentage of Index RAM quota in use across all indexes on this server.",
			nil,
			nil,
		),
		IndexRemainingRAM: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "remaining_ram"),
			"Bytes of Index RAM quota still available on this server.",
			nil,
			nil,
		),
	}
}

// Describe all metrics
func (c *indexCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.IndexMemoryUsed
	ch <- c.IndexMemoryQuota
	ch <- c.IndexRAMPercent
	ch <- c.IndexRemainingRAM
}

// Collect all metrics
func (c *indexCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting query metrics...")

	indexStats, err := c.m.client.Index()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.Error("failed to scrape index stats")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.IndexMemoryQuota, prometheus.GaugeValue, last(indexStats.Op.Samples.IndexMemoryQuota))
	ch <- prometheus.MustNewConstMetric(c.IndexMemoryUsed, prometheus.GaugeValue, last(indexStats.Op.Samples.IndexMemoryUsed))
	ch <- prometheus.MustNewConstMetric(c.IndexRAMPercent, prometheus.GaugeValue, last(indexStats.Op.Samples.IndexRAMPercent))
	ch <- prometheus.MustNewConstMetric(c.IndexRemainingRAM, prometheus.GaugeValue, last(indexStats.Op.Samples.IndexRemainingRAM))

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}
