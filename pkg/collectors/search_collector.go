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

type ftsCollector struct {
	m MetaCollector

	FtsCurrBatchesBlockedByHerder   *prometheus.Desc
	FtsNumBytesUsedRAM              *prometheus.Desc
	FtsTotalQueriesRejectedByHerder *prometheus.Desc
}

func NewFTSCollector(client util.Client) prometheus.Collector {
	const subsystem = "fts"
	return &ftsCollector{
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
		FtsCurrBatchesBlockedByHerder: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "curr_batches_blocked_by_herder"),
			"Number of DCP batches blocked by the FTS throttler due to high memory consumption",
			[]string{"cluster"},
			nil,
		),
		FtsNumBytesUsedRAM: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "num_bytes_used_ram"),
			"Amount of RAM used by FTS on this server",
			[]string{"cluster"},
			nil,
		),
		FtsTotalQueriesRejectedByHerder: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE+subsystem, "", "total_queries_rejected_by_herder"),
			"Number of fts queries rejected by the FTS throttler due to high memory consumption",
			[]string{"cluster"},
			nil,
		),
	}
}

// Describe all metrics
func (c *ftsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration
	ch <- c.FtsTotalQueriesRejectedByHerder
	ch <- c.FtsNumBytesUsedRAM
	ch <- c.FtsTotalQueriesRejectedByHerder
}

// Collect all metrics
func (c *ftsCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting fts metrics...")

	clusterName, err := c.m.client.ClusterName()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("%s", err)
		return
	}

	ftsStats, err := c.m.client.Fts()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0, clusterName)
		log.Error("failed to scrape FTS stats")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.FtsCurrBatchesBlockedByHerder, prometheus.GaugeValue, last(ftsStats.Op.Samples.FtsCurrBatchesBlockedByHerder), clusterName)
	ch <- prometheus.MustNewConstMetric(c.FtsNumBytesUsedRAM, prometheus.GaugeValue, last(ftsStats.Op.Samples.FtsNumBytesUsedRAM), clusterName)
	ch <- prometheus.MustNewConstMetric(c.FtsTotalQueriesRejectedByHerder, prometheus.GaugeValue, last(ftsStats.Op.Samples.FtsTotalQueriesRejectedByHerder), clusterName)

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1, clusterName)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds(), clusterName)
}
