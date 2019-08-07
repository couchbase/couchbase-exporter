package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/couchbase/couchbase_exporter/util"
	"time"
)

const FQ_NAMESPACE = "couchbase"

type bucketStatsCollector struct {
	m MetaCollector

	AvgBgWaitTime                    *prometheus.Desc
	AvgActiveTimestampDrift          *prometheus.Desc
	AvgReplicaTimestampDrift         *prometheus.Desc
	AvgDiskCommitTime                *prometheus.Desc
	AvgDiskUpdateTime                *prometheus.Desc
	BytesRead                        *prometheus.Desc
	BytesWritten                     *prometheus.Desc
	CasBadval                        *prometheus.Desc
	CasHits                          *prometheus.Desc
	CasMisses                        *prometheus.Desc
	CmdGet                           *prometheus.Desc
	CmdSet                           *prometheus.Desc
	CouchTotalDiskSize               *prometheus.Desc
	CouchViewsDataSize               *prometheus.Desc
	CouchViewsActualDiskSize         *prometheus.Desc
	CouchViewsFragmentation          *prometheus.Desc
	CouchViewsOps                    *prometheus.Desc
	CouchDocsDataSize                *prometheus.Desc
	CouchDocsDiskSize                *prometheus.Desc
	CouchDocsActualDiskSize          *prometheus.Desc
	CouchDocsFragmentation           *prometheus.Desc
	CPUIdleMs                        *prometheus.Desc
	CPULocalMs                       *prometheus.Desc
	CPUUtilizationRate               *prometheus.Desc
	CurrConnections                  *prometheus.Desc
	CurrItems                        *prometheus.Desc
	CurrItemsTot                     *prometheus.Desc
	DecrHits                         *prometheus.Desc
	DecrMisses                       *prometheus.Desc
	DeleteHits                       *prometheus.Desc
	DeleteMisses                     *prometheus.Desc
	DiskCommitCount                  *prometheus.Desc
	DiskUpdateCount                  *prometheus.Desc
	DiskWriteQueue                   *prometheus.Desc
	EpActiveAheadExceptions          *prometheus.Desc
	EpActiveHlcDrift                 *prometheus.Desc
	EpClockCasDriftThresholdExceeded *prometheus.Desc
	EpBgFetched                      *prometheus.Desc
	EpCacheMissRate                  *prometheus.Desc
	EpDcp2iBackoff                   *prometheus.Desc
	EpDcp2iCount                     *prometheus.Desc
	EpDcp2iItemsRemaining            *prometheus.Desc
	EpDcp2iItemsSent                 *prometheus.Desc
	EpDcp2iProducerCount             *prometheus.Desc
	EpDcp2iTotalBacklogSize          *prometheus.Desc
	EpDcp2iTotalBytes                *prometheus.Desc
	EpDcpOtherBackoff                *prometheus.Desc
	EpDcpOtherCount                  *prometheus.Desc
	EpDcpOtherItemsRemaining         *prometheus.Desc
	EpDcpOtherItemsSent              *prometheus.Desc
	EpDcpOtherProducerCount          *prometheus.Desc
	EpDcpOtherTotalBacklogSize       *prometheus.Desc
	EpDcpOtherTotalBytes             *prometheus.Desc
	EpDcpReplicaBackoff              *prometheus.Desc
	EpDcpReplicaCount                *prometheus.Desc
	EpDcpReplicaItemsRemaining       *prometheus.Desc
	EpDcpReplicaItemsSent            *prometheus.Desc
	EpDcpReplicaProducerCount        *prometheus.Desc
	EpDcpReplicaTotalBacklogSize     *prometheus.Desc
	EpDcpReplicaTotalBytes           *prometheus.Desc
	EpDcpViewsBackoff                *prometheus.Desc
	EpDcpViewsCount                  *prometheus.Desc
	EpDcpViewsItemsRemaining         *prometheus.Desc
	EpDcpViewsItemsSent              *prometheus.Desc
	EpDcpViewsProducerCount          *prometheus.Desc
	EpDcpViewsTotalBacklogSize       *prometheus.Desc
	EpDcpViewsTotalBytes             *prometheus.Desc
	EpDcpXdcrBackoff                 *prometheus.Desc
	EpDcpXdcrCount                   *prometheus.Desc
	EpDcpXdcrItemsRemaining          *prometheus.Desc
	EpDcpXdcrItemsSent               *prometheus.Desc
	EpDcpXdcrProducerCount           *prometheus.Desc
	EpDcpXdcrTotalBacklogSize        *prometheus.Desc
	EpDcpXdcrTotalBytes              *prometheus.Desc
	EpDiskqueueDrain                 *prometheus.Desc
	EpDiskqueueFill                  *prometheus.Desc
	EpDiskqueueItems                 *prometheus.Desc
	EpFlusherTodo                    *prometheus.Desc
	EpItemCommitFailed               *prometheus.Desc
	EpKvSize                         *prometheus.Desc
	EpMaxSize                        *prometheus.Desc
	EpMemHighWat                     *prometheus.Desc
	EpMemLowWat                      *prometheus.Desc
	EpMetaDataMemory                 *prometheus.Desc
	EpNumNonResident                 *prometheus.Desc
	EpNumOpsDelMeta                  *prometheus.Desc
	EpNumOpsDelRetMeta               *prometheus.Desc
	EpNumOpsGetMeta                  *prometheus.Desc
	EpNumOpsSetMeta                  *prometheus.Desc
	EpNumOpsSetRetMeta               *prometheus.Desc
	EpNumValueEjects                 *prometheus.Desc
	EpOomErrors                      *prometheus.Desc
	EpOpsCreate                      *prometheus.Desc
	EpOpsUpdate                      *prometheus.Desc
	EpOverhead                       *prometheus.Desc
	EpQueueSize                      *prometheus.Desc
	EpResidentItemsRate              *prometheus.Desc
	EpReplicaAheadExceptions         *prometheus.Desc
	EpReplicaHlcDrift                *prometheus.Desc
	EpTmpOomErrors                   *prometheus.Desc
	EpVbTotal                        *prometheus.Desc
	Evictions                        *prometheus.Desc
	GetHits                          *prometheus.Desc
	GetMisses                        *prometheus.Desc
	HibernatedRequests               *prometheus.Desc
	HibernatedWaked                  *prometheus.Desc
	HitRatio                         *prometheus.Desc
	IncrHits                         *prometheus.Desc
	IncrMisses                       *prometheus.Desc
	MemActualFree                    *prometheus.Desc
	MemActualUsed                    *prometheus.Desc
	MemFree                          *prometheus.Desc
	MemTotal                         *prometheus.Desc
	MemUsed                          *prometheus.Desc
	MemUsedSys                       *prometheus.Desc
	Misses                           *prometheus.Desc
	Ops                              *prometheus.Desc
	RestRequests                     *prometheus.Desc
	SwapTotal                        *prometheus.Desc
	SwapUsed                         *prometheus.Desc
	Timestamp                        *prometheus.Desc
	VbActiveEject                    *prometheus.Desc
	VbActiveItmMemory                *prometheus.Desc
	VbActiveMetaDataMemory           *prometheus.Desc
	VbActiveNum                      *prometheus.Desc
	VbActiveNumNonResident           *prometheus.Desc
	VbActiveOpsCreate                *prometheus.Desc
	VbActiveOpsUpdate                *prometheus.Desc
	VbActiveQueueAge                 *prometheus.Desc
	VbActiveQueueDrain               *prometheus.Desc
	VbActiveQueueFill                *prometheus.Desc
	VbActiveQueueSize                *prometheus.Desc
	VbActiveResidentItemsRatio       *prometheus.Desc
	VbAvgActiveQueueAge              *prometheus.Desc
	VbAvgPendingQueueAge             *prometheus.Desc
	VbAvgReplicaQueueAge             *prometheus.Desc
	VbAvgTotalQueueAge               *prometheus.Desc
	VbPendingCurrItems               *prometheus.Desc
	VbPendingEject                   *prometheus.Desc
	VbPendingItmMemory               *prometheus.Desc
	VbPendingMetaDataMemory          *prometheus.Desc
	VbPendingNum                     *prometheus.Desc
	VbPendingNumNonResident          *prometheus.Desc
	VbPendingOpsCreate               *prometheus.Desc
	VbPendingOpsUpdate               *prometheus.Desc
	VbPendingQueueAge                *prometheus.Desc
	VbPendingQueueDrain              *prometheus.Desc
	VbPendingQueueFill               *prometheus.Desc
	VbPendingQueueSize               *prometheus.Desc
	VbPendingResidentItemsRatio      *prometheus.Desc
	VbReplicaCurrItems               *prometheus.Desc
	VbReplicaEject                   *prometheus.Desc
	VbReplicaItmMemory               *prometheus.Desc
	VbReplicaMetaDataMemory          *prometheus.Desc
	VbReplicaNum                     *prometheus.Desc
	VbReplicaNumNonResident          *prometheus.Desc
	VbReplicaOpsCreate               *prometheus.Desc
	VbReplicaOpsUpdate               *prometheus.Desc
	VbReplicaQueueAge                *prometheus.Desc
	VbReplicaQueueDrain              *prometheus.Desc
	VbReplicaQueueFill               *prometheus.Desc
	VbReplicaQueueSize               *prometheus.Desc
	VbReplicaResidentItemsRatio      *prometheus.Desc
	VbTotalQueueAge                  *prometheus.Desc
	XdcOps                           *prometheus.Desc
}

func last(ss []float64) float64 {
	if len(ss) == 0 {
		return 0.0
	}
	return ss[len(ss)-1]
}

// Describe all metrics
func (c *bucketStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.m.up
	ch <- c.m.scrapeDuration

	ch <- c.AvgBgWaitTime
	ch <- c.AvgActiveTimestampDrift
	ch <- c.AvgReplicaTimestampDrift
	ch <- c.AvgDiskCommitTime
	ch <- c.AvgDiskUpdateTime
	ch <- c.BytesRead
	ch <- c.BytesWritten
	ch <- c.CasBadval
	ch <- c.CasHits
	ch <- c.CasMisses
	ch <- c.CmdGet
	ch <- c.CmdSet
	ch <- c.CouchDocsDataSize
	ch <- c.CouchDocsActualDiskSize
	ch <- c.CouchTotalDiskSize
	ch <- c.CouchViewsDataSize
	ch <- c.CouchViewsActualDiskSize
	ch <- c.CouchViewsFragmentation
	ch <- c.CouchViewsOps
	ch <- c.CouchDocsFragmentation
	ch <- c.CouchDocsDiskSize
	ch <- c.CPUIdleMs
	ch <- c.CPULocalMs
	ch <- c.CPUUtilizationRate
	ch <- c.CurrConnections
	ch <- c.CurrItems
	ch <- c.CurrItemsTot
	ch <- c.DecrHits
	ch <- c.DecrMisses
	ch <- c.DeleteHits
	ch <- c.DeleteMisses
	ch <- c.DiskCommitCount
	ch <- c.DiskUpdateCount
	ch <- c.DiskWriteQueue
	ch <- c.EpActiveAheadExceptions
	ch <- c.EpActiveHlcDrift
	ch <- c.EpClockCasDriftThresholdExceeded
	ch <- c.EpBgFetched
	ch <- c.EpCacheMissRate
	ch <- c.EpDcp2iBackoff
	ch <- c.EpDcp2iCount
	ch <- c.EpDcp2iItemsRemaining
	ch <- c.EpDcp2iItemsSent
	ch <- c.EpDcp2iProducerCount
	ch <- c.EpDcp2iTotalBacklogSize
	ch <- c.EpDcp2iTotalBytes
	ch <- c.EpDcpOtherBackoff
	ch <- c.EpDcpOtherCount
	ch <- c.EpDcpOtherItemsRemaining
	ch <- c.EpDcpOtherItemsSent
	ch <- c.EpDcpOtherProducerCount
	ch <- c.EpDcpOtherTotalBacklogSize
	ch <- c.EpDcpOtherTotalBytes
	ch <- c.EpDcpReplicaBackoff
	ch <- c.EpDcpReplicaCount
	ch <- c.EpDcpReplicaItemsRemaining
	ch <- c.EpDcpReplicaItemsSent
	ch <- c.EpDcpReplicaProducerCount
	ch <- c.EpDcpReplicaTotalBacklogSize
	ch <- c.EpDcpReplicaTotalBytes
	ch <- c.EpDcpViewsBackoff
	ch <- c.EpDcpViewsCount
	ch <- c.EpDcpViewsItemsRemaining
	ch <- c.EpDcpViewsItemsSent
	ch <- c.EpDcpViewsProducerCount
	ch <- c.EpDcpViewsTotalBacklogSize
	ch <- c.EpDcpViewsTotalBytes
	ch <- c.EpDcpXdcrBackoff
	ch <- c.EpDcpXdcrCount
	ch <- c.EpDcpXdcrItemsRemaining
	ch <- c.EpDcpXdcrItemsSent
	ch <- c.EpDcpXdcrProducerCount
	ch <- c.EpDcpXdcrTotalBacklogSize
	ch <- c.EpDcpXdcrTotalBytes
	ch <- c.EpDiskqueueDrain
	ch <- c.EpDiskqueueFill
	ch <- c.EpDiskqueueItems
	ch <- c.EpFlusherTodo
	ch <- c.EpItemCommitFailed
	ch <- c.EpKvSize
	ch <- c.EpMaxSize
	ch <- c.EpMemHighWat
	ch <- c.EpMemLowWat
	ch <- c.EpMetaDataMemory
	ch <- c.EpNumNonResident
	ch <- c.EpNumOpsDelMeta
	ch <- c.EpNumOpsDelRetMeta
	ch <- c.EpNumOpsGetMeta
	ch <- c.EpNumOpsSetMeta
	ch <- c.EpNumOpsSetRetMeta
	ch <- c.EpNumValueEjects
	ch <- c.EpOomErrors
	ch <- c.EpOpsCreate
	ch <- c.EpOpsUpdate
	ch <- c.EpOverhead
	ch <- c.EpQueueSize
	ch <- c.EpResidentItemsRate
	ch <- c.EpReplicaAheadExceptions
	ch <- c.EpReplicaHlcDrift
	ch <- c.EpTmpOomErrors
	ch <- c.EpVbTotal
	ch <- c.Evictions
	ch <- c.GetHits
	ch <- c.GetMisses
	ch <- c.HibernatedRequests
	ch <- c.HibernatedWaked
	ch <- c.HitRatio
	ch <- c.IncrHits
	ch <- c.IncrMisses
	ch <- c.MemActualFree
	ch <- c.MemActualUsed
	ch <- c.MemFree
	ch <- c.MemTotal
	ch <- c.MemUsed
	ch <- c.MemUsedSys
	ch <- c.Misses
	ch <- c.Ops
	ch <- c.RestRequests
	ch <- c.SwapTotal
	ch <- c.SwapUsed
	ch <- c.Timestamp
	ch <- c.VbActiveEject
	ch <- c.VbActiveItmMemory
	ch <- c.VbActiveMetaDataMemory
	ch <- c.VbActiveNum
	ch <- c.VbActiveNumNonResident
	ch <- c.VbActiveOpsCreate
	ch <- c.VbActiveOpsUpdate
	ch <- c.VbActiveQueueAge
	ch <- c.VbActiveQueueDrain
	ch <- c.VbActiveQueueFill
	ch <- c.VbActiveQueueSize
	ch <- c.VbActiveResidentItemsRatio
	ch <- c.VbAvgActiveQueueAge
	ch <- c.VbAvgPendingQueueAge
	ch <- c.VbAvgReplicaQueueAge
	ch <- c.VbAvgTotalQueueAge
	ch <- c.VbPendingCurrItems
	ch <- c.VbPendingEject
	ch <- c.VbPendingItmMemory
	ch <- c.VbPendingMetaDataMemory
	ch <- c.VbPendingNum
	ch <- c.VbPendingNumNonResident
	ch <- c.VbPendingOpsCreate
	ch <- c.VbPendingOpsUpdate
	ch <- c.VbPendingQueueAge
	ch <- c.VbPendingQueueDrain
	ch <- c.VbPendingQueueFill
	ch <- c.VbPendingQueueSize
	ch <- c.VbPendingResidentItemsRatio
	ch <- c.VbReplicaCurrItems
	ch <- c.VbReplicaEject
	ch <- c.VbReplicaItmMemory
	ch <- c.VbReplicaMetaDataMemory
	ch <- c.VbReplicaNum
	ch <- c.VbReplicaNumNonResident
	ch <- c.VbReplicaOpsCreate
	ch <- c.VbReplicaOpsUpdate
	ch <- c.VbReplicaQueueAge
	ch <- c.VbReplicaQueueDrain
	ch <- c.VbReplicaQueueFill
	ch <- c.VbReplicaQueueSize
	ch <- c.VbReplicaResidentItemsRatio
	ch <- c.VbTotalQueueAge
	ch <- c.XdcOps
}

func (c *bucketStatsCollector) Collect(ch chan<- prometheus.Metric) {
	c.m.mutex.Lock()
	defer c.m.mutex.Unlock()

	start := time.Now()
	log.Info("Collecting buckets metrics...")

	buckets, err := c.m.client.Buckets()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
		log.With("error", err).Error("failed to scrape buckets")
		return
	}

	// nolint: lll
	for _, bucket := range buckets {
		log.Debugf("Collecting %s bucket stats metrics...", bucket.Name)
		stats, err := c.m.client.BucketStats(bucket.Name)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 0)
			log.With("error", err).Error("failed to scrape bucket stats")
			return
		}

		ch <- prometheus.MustNewConstMetric(c.AvgBgWaitTime, prometheus.GaugeValue, last(stats.Op.Samples.AvgBgWaitTime)/1000000, bucket.Name) // this comes as microseconds from CB
		ch <- prometheus.MustNewConstMetric(c.AvgActiveTimestampDrift, prometheus.GaugeValue, last(stats.Op.Samples.AvgActiveTimestampDrift), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.AvgReplicaTimestampDrift, prometheus.GaugeValue, last(stats.Op.Samples.AvgReplicaTimestampDrift), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.AvgDiskCommitTime, prometheus.GaugeValue, last(stats.Op.Samples.AvgDiskCommitTime), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.AvgDiskUpdateTime, prometheus.GaugeValue, last(stats.Op.Samples.AvgDiskUpdateTime), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.BytesRead, prometheus.GaugeValue, last(stats.Op.Samples.BytesRead), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.BytesWritten, prometheus.GaugeValue, last(stats.Op.Samples.BytesWritten), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CasBadval, prometheus.GaugeValue, last(stats.Op.Samples.CasBadval), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CasHits, prometheus.GaugeValue, last(stats.Op.Samples.CasHits), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CasMisses, prometheus.GaugeValue, last(stats.Op.Samples.CasMisses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CmdGet, prometheus.GaugeValue, last(stats.Op.Samples.CmdGet), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CmdSet, prometheus.GaugeValue, last(stats.Op.Samples.CmdSet), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchTotalDiskSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchTotalDiskSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchViewsDataSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchViewsDataSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchViewsActualDiskSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchViewsActualDiskSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchViewsFragmentation, prometheus.GaugeValue, last(stats.Op.Samples.CouchViewsFragmentation), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchViewsOps, prometheus.GaugeValue, last(stats.Op.Samples.CouchViewsOps), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchDocsActualDiskSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchDocsActualDiskSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchDocsFragmentation, prometheus.GaugeValue, last(stats.Op.Samples.CouchDocsFragmentation), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchDocsDataSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchDocsDataSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CouchDocsDiskSize, prometheus.GaugeValue, last(stats.Op.Samples.CouchDocsDiskSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CPUIdleMs, prometheus.GaugeValue, last(stats.Op.Samples.CPUIdleMs), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CPULocalMs, prometheus.GaugeValue, last(stats.Op.Samples.CPULocalMs), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CPUUtilizationRate, prometheus.GaugeValue, last(stats.Op.Samples.CPUUtilizationRate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CurrConnections, prometheus.GaugeValue, last(stats.Op.Samples.CurrConnections), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CurrItems, prometheus.GaugeValue, last(stats.Op.Samples.CurrItems), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.CurrItemsTot, prometheus.GaugeValue, last(stats.Op.Samples.CurrItemsTot), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DecrHits, prometheus.GaugeValue, last(stats.Op.Samples.DecrHits), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DecrMisses, prometheus.GaugeValue, last(stats.Op.Samples.DecrMisses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DeleteHits, prometheus.GaugeValue, last(stats.Op.Samples.DeleteHits), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DeleteMisses, prometheus.GaugeValue, last(stats.Op.Samples.DeleteMisses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DiskCommitCount, prometheus.GaugeValue, last(stats.Op.Samples.DiskCommitCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DiskUpdateCount, prometheus.GaugeValue, last(stats.Op.Samples.DiskUpdateCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.DiskWriteQueue, prometheus.GaugeValue, last(stats.Op.Samples.DiskWriteQueue), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpActiveAheadExceptions, prometheus.GaugeValue, last(stats.Op.Samples.EpActiveAheadExceptions), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpActiveHlcDrift, prometheus.GaugeValue, last(stats.Op.Samples.EpActiveHlcDrift), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpClockCasDriftThresholdExceeded, prometheus.GaugeValue, last(stats.Op.Samples.EpClockCasDriftThresholdExceeded), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpBgFetched, prometheus.GaugeValue, last(stats.Op.Samples.EpBgFetched), bucket.Name)
		// for some reason, this ratio can be > 100, so we added a `min` function
		//ch <- prometheus.MustNewConstMetric(c.EpCacheMissRate, prometheus.GaugeValue, min(last(stats.Op.Samples.EpCacheMissRate), 100), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpCacheMissRate, prometheus.GaugeValue, last(stats.Op.Samples.EpCacheMissRate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iBackoff, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2IBackoff), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2ICount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iItemsRemaining, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2IItemsRemaining), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iItemsSent, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2IItemsSent), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iProducerCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2IProducerCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iTotalBacklogSize, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2ITotalBacklogSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcp2iTotalBytes, prometheus.GaugeValue, last(stats.Op.Samples.EpDcp2ITotalBytes), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherBackoff, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherBackoff), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherItemsRemaining, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherItemsRemaining), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherItemsSent, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherItemsSent), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherProducerCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherProducerCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherTotalBacklogSize, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherTotalBacklogSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpOtherTotalBytes, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpOtherTotalBytes), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaBackoff, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaBackoff), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaItemsRemaining, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaItemsRemaining), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaItemsSent, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaItemsSent), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaProducerCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaProducerCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaTotalBacklogSize, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaTotalBacklogSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpReplicaTotalBytes, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpReplicaTotalBytes), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsBackoff, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsBackoff), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsItemsRemaining, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsItemsRemaining), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsItemsSent, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsItemsSent), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsProducerCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsProducerCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsTotalBacklogSize, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsTotalBacklogSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpViewsTotalBytes, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpViewsTotalBytes), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrBackoff, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrBackoff), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrItemsRemaining, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrItemsRemaining), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrItemsSent, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrItemsSent), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrProducerCount, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrProducerCount), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrTotalBacklogSize, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrTotalBacklogSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDcpXdcrTotalBytes, prometheus.GaugeValue, last(stats.Op.Samples.EpDcpXdcrTotalBytes), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDiskqueueDrain, prometheus.GaugeValue, last(stats.Op.Samples.EpDiskqueueDrain), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDiskqueueFill, prometheus.GaugeValue, last(stats.Op.Samples.EpDiskqueueFill), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpDiskqueueItems, prometheus.GaugeValue, last(stats.Op.Samples.EpDiskqueueItems), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpFlusherTodo, prometheus.GaugeValue, last(stats.Op.Samples.EpFlusherTodo), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpItemCommitFailed, prometheus.GaugeValue, last(stats.Op.Samples.EpItemCommitFailed), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpKvSize, prometheus.GaugeValue, last(stats.Op.Samples.EpKvSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpMaxSize, prometheus.GaugeValue, last(stats.Op.Samples.EpMaxSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpMemHighWat, prometheus.GaugeValue, last(stats.Op.Samples.EpMemHighWat), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpMemLowWat, prometheus.GaugeValue, last(stats.Op.Samples.EpMemLowWat), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpMetaDataMemory, prometheus.GaugeValue, last(stats.Op.Samples.EpMetaDataMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumNonResident, prometheus.GaugeValue, last(stats.Op.Samples.EpNumNonResident), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumOpsDelMeta, prometheus.GaugeValue, last(stats.Op.Samples.EpNumOpsDelMeta), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumOpsDelRetMeta, prometheus.GaugeValue, last(stats.Op.Samples.EpNumOpsDelRetMeta), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumOpsGetMeta, prometheus.GaugeValue, last(stats.Op.Samples.EpNumOpsGetMeta), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumOpsSetMeta, prometheus.GaugeValue, last(stats.Op.Samples.EpNumOpsSetMeta), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumOpsSetRetMeta, prometheus.GaugeValue, last(stats.Op.Samples.EpNumOpsSetRetMeta), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpNumValueEjects, prometheus.GaugeValue, last(stats.Op.Samples.EpNumValueEjects), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpOomErrors, prometheus.GaugeValue, last(stats.Op.Samples.EpOomErrors), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpOpsCreate, prometheus.GaugeValue, last(stats.Op.Samples.EpOpsCreate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpOpsUpdate, prometheus.GaugeValue, last(stats.Op.Samples.EpOpsUpdate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpOverhead, prometheus.GaugeValue, last(stats.Op.Samples.EpOverhead), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpQueueSize, prometheus.GaugeValue, last(stats.Op.Samples.EpQueueSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpResidentItemsRate, prometheus.GaugeValue, last(stats.Op.Samples.EpResidentItemsRate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpReplicaAheadExceptions, prometheus.GaugeValue, last(stats.Op.Samples.EpReplicaAheadExceptions), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpReplicaHlcDrift, prometheus.GaugeValue, last(stats.Op.Samples.EpReplicaHlcDrift), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpTmpOomErrors, prometheus.GaugeValue, last(stats.Op.Samples.EpTmpOomErrors), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.EpVbTotal, prometheus.GaugeValue, last(stats.Op.Samples.EpVbTotal), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.Evictions, prometheus.GaugeValue, last(stats.Op.Samples.Evictions), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.GetHits, prometheus.GaugeValue, last(stats.Op.Samples.GetHits), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.GetMisses, prometheus.GaugeValue, last(stats.Op.Samples.GetMisses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.HibernatedRequests, prometheus.GaugeValue, last(stats.Op.Samples.HibernatedRequests), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.HibernatedWaked, prometheus.GaugeValue, last(stats.Op.Samples.HibernatedWaked), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.HitRatio, prometheus.GaugeValue, last(stats.Op.Samples.HitRatio), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.IncrHits, prometheus.GaugeValue, last(stats.Op.Samples.IncrHits), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.IncrMisses, prometheus.GaugeValue, last(stats.Op.Samples.IncrMisses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemActualFree, prometheus.GaugeValue, last(stats.Op.Samples.MemActualFree), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemActualUsed, prometheus.GaugeValue, last(stats.Op.Samples.MemActualUsed), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemFree, prometheus.GaugeValue, last(stats.Op.Samples.MemFree), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemTotal, prometheus.GaugeValue, last(stats.Op.Samples.MemTotal), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemUsed, prometheus.GaugeValue, last(stats.Op.Samples.MemUsed), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.MemUsedSys, prometheus.GaugeValue, last(stats.Op.Samples.MemUsedSys), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.Misses, prometheus.GaugeValue, last(stats.Op.Samples.Misses), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.Ops, prometheus.GaugeValue, last(stats.Op.Samples.Ops), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.RestRequests, prometheus.GaugeValue, last(stats.Op.Samples.RestRequests), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.SwapTotal, prometheus.GaugeValue, last(stats.Op.Samples.SwapTotal), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.SwapUsed, prometheus.GaugeValue, last(stats.Op.Samples.SwapUsed), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.Timestamp, prometheus.GaugeValue, last(stats.Op.Samples.Timestamp), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveEject, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveEject), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveItmMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveItmMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveMetaDataMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveMetaDataMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveNum, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveNum), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveNumNonResident, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveNumNonResident), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveOpsCreate, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveOpsCreate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveOpsUpdate, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveOpsUpdate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveQueueDrain, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveQueueDrain), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveQueueFill, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveQueueFill), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveQueueSize, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveQueueSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbActiveResidentItemsRatio, prometheus.GaugeValue, last(stats.Op.Samples.VbActiveResidentItemsRatio), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbAvgActiveQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbAvgActiveQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbAvgPendingQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbAvgPendingQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbAvgReplicaQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbAvgReplicaQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbAvgTotalQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbAvgTotalQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingCurrItems, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingCurrItems), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingEject, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingEject), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingItmMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingItmMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingMetaDataMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingMetaDataMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingNum, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingNum), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingNumNonResident, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingNumNonResident), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingOpsCreate, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingOpsCreate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingOpsUpdate, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingOpsUpdate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingQueueDrain, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingQueueDrain), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingQueueFill, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingQueueFill), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingQueueSize, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingQueueSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbPendingResidentItemsRatio, prometheus.GaugeValue, last(stats.Op.Samples.VbPendingResidentItemsRatio), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaCurrItems, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaCurrItems), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaEject, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaEject), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaItmMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaItmMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaMetaDataMemory, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaMetaDataMemory), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaNum, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaNum), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaNumNonResident, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaNumNonResident), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaOpsCreate, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaOpsCreate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaOpsUpdate, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaOpsUpdate), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaQueueDrain, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaQueueDrain), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaQueueFill, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaQueueFill), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaQueueSize, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaQueueSize), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbReplicaResidentItemsRatio, prometheus.GaugeValue, last(stats.Op.Samples.VbReplicaResidentItemsRatio), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.VbTotalQueueAge, prometheus.GaugeValue, last(stats.Op.Samples.VbTotalQueueAge), bucket.Name)
		ch <- prometheus.MustNewConstMetric(c.XdcOps, prometheus.GaugeValue, last(stats.Op.Samples.XdcOps), bucket.Name)
	}

	ch <- prometheus.MustNewConstMetric(c.m.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.m.scrapeDuration, prometheus.GaugeValue, time.Since(start).Seconds())
}

func NewBucketStatsCollector(client util.Client) prometheus.Collector {
	const subsystem = "bucketstat"

	// nolint: lll
	return &bucketStatsCollector{
		m: MetaCollector{
			client: client,
			up: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "up"),
				"Couchbase buckets API is responding",
				nil,
				nil,
			),
			scrapeDuration: prometheus.NewDesc(
				prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "scrape_duration_seconds"),
				"Scrape duration in seconds",
				nil,
				nil,
			),
		},
		AvgBgWaitTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "avg_bg_wait_seconds"),
			"Average background fetch time in seconds",
			[]string{"bucket"},
			nil,
		),
		AvgActiveTimestampDrift: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "avg_active_timestamp_drift"),
			"avg_active_timestamp_drift",
			[]string{"bucket"},
			nil,
		),
		AvgReplicaTimestampDrift: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "avg_replica_timestamp_drift"),
			"avg_replica_timestamp_drift",
			[]string{"bucket"},
			nil,
		),
		AvgDiskCommitTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "avg_disk_commit_time"),
			"Average disk commit time in seconds as from disk_update histogram of timings",
			[]string{"bucket"},
			nil,
		),
		AvgDiskUpdateTime: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "avg_disk_update_time"),
			"Average disk update time in microseconds as from disk_update histogram of timings",
			[]string{"bucket"},
			nil,
		),
		BytesRead: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "read_bytes"),
			"Bytes read",
			[]string{"bucket"},
			nil,
		),
		BytesWritten: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "written_bytes"),
			"Bytes written",
			[]string{"bucket"},
			nil,
		),
		CasBadval: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cas_badval"),
			"Compare and Swap bad values",
			[]string{"bucket"},
			nil,
		),
		CasHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cas_hits"),
			"Number of operations with a CAS id per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		CasMisses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cas_misses"),
			"Compare and Swap misses",
			[]string{"bucket"},
			nil,
		),
		CmdGet: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cmd_get"),
			"Number of reads (get operations) per second from this bucket",
			[]string{"bucket"},
			nil,
		),
		CmdSet: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cmd_set"),
			"Number of writes (set operations) per second to this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchTotalDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_total_disk_size"),
			"The total size on disk of all data and view files for this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchViewsFragmentation: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_views_fragmentation"),
			"How much fragmented data there is to be compacted compared to real data for the view index files in this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchViewsOps: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_views_ops"),
			"All the view reads for all design documents including scatter gather",
			[]string{"bucket"},
			nil,
		),
		CouchViewsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_views_data_size"),
			"The size of active data on for all the indexes in this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchViewsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_views_actual_disk_size"),
			"The size of all active items in all the indexes for this bucket on disk",
			[]string{"bucket"},
			nil,
		),
		CouchDocsFragmentation: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_docs_fragmentation"),
			"How much fragmented data there is to be compacted compared to real data for the data files in this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchDocsActualDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_docs_actual_disk_size"),
			"The size of all data files for this bucket, including the data itself, meta data and temporary files",
			[]string{"bucket"},
			nil,
		),
		CouchDocsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_docs_data_size"),
			"The size of active data in this bucket",
			[]string{"bucket"},
			nil,
		),
		CouchDocsDiskSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "couch_docs_disk_size"),
			"The size of all data files for this bucket, including the data itself, meta data and temporary files",
			[]string{"bucket"},
			nil,
		),
		CPUIdleMs: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cpu_idle_ms"),
			"CPU idle milliseconds",
			[]string{"bucket"},
			nil,
		),
		CPULocalMs: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cpu_local_ms"),
			"_cpu_local_ms",
			[]string{"bucket"},
			nil,
		),
		CPUUtilizationRate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "cpu_utilization_rate"),
			"Percentage of CPU in use across all available cores on this server",
			[]string{"bucket"},
			nil,
		),
		CurrConnections: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "urr_connections"),
			"Number of connections to this server including connections from external client SDKs, proxies, DCP requests and internal statistic gathering",
			[]string{"bucket"},
			nil,
		),
		CurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "curr_items"),
			"Number of items in active vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		CurrItemsTot: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "curr_items_tot"),
			"Total number of items in this bucket",
			[]string{"bucket"},
			nil,
		),
		DecrHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "decr_hits"),
			"Decrement hits",
			[]string{"bucket"},
			nil,
		),
		DecrMisses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "decr_misses"),
			"Decrement misses",
			[]string{"bucket"},
			nil,
		),
		DeleteHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "delete_hits"),
			"Number of delete operations per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		DeleteMisses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "delete_misses"),
			"Delete misses",
			[]string{"bucket"},
			nil,
		),
		DiskCommitCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "disk_commits"),
			"Disk commits",
			[]string{"bucket"},
			nil,
		),
		DiskUpdateCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "disk_updates"),
			"Disk updates",
			[]string{"bucket"},
			nil,
		),
		DiskWriteQueue: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "disk_write_queue"),
			"Number of items waiting to be written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpActiveAheadExceptions: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_active_ahead_exceptions"),
			"_ep_active_ahead_exceptions",
			[]string{"bucket"},
			nil,
		),
		EpActiveHlcDrift: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_active_hlc_drift"),
			"_ep_active_hlc_drift",
			[]string{"bucket"},
			nil,
		),
		EpClockCasDriftThresholdExceeded: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_clock_cas_drift_threshold_exceeded"),
			"_ep_clock_cas_drift_threshold_exceeded",
			[]string{"bucket"},
			nil,
		),
		EpBgFetched: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_bg_fetched"),
			"Number of reads per second from disk for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpCacheMissRate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_cache_miss_rate"),
			"Percentage of reads per second to this bucket from disk as opposed to RAM",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iBackoff: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_backoff"),
			"Number of backoffs for indexes DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_connections"),
			"Number of indexes DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iItemsRemaining: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_items_remaining"),
			"Number of indexes items remaining to be sent",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iItemsSent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_items_sent"),
			"Number of indexes items sent",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iProducerCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_producers"),
			"Number of indexes producers",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iTotalBacklogSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_total_backlog_size"),
			"ep_dcp_2i_total_backlog_size",
			[]string{"bucket"},
			nil,
		),
		EpDcp2iTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_2i_total_bytes"),
			"Number bytes per second being sent for indexes DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherBackoff: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_backoff"),
			"Number of backoffs for other DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_others"),
			"Number of other DCP connections in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherItemsRemaining: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_items_remaining"),
			"Number of items remaining to be sent to consumer in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherItemsSent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_items_sent"),
			"Number of items per second being sent for a producer for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherProducerCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_producers"),
			"Number of other senders for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherTotalBacklogSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_total_backlog_size"),
			"ep_dcp_other_total_backlog_size",
			[]string{"bucket"},
			nil,
		),
		EpDcpOtherTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_other_total_bytes"),
			"Number of bytes per second being sent for other DCP connections for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaBackoff: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_backoff"),
			"Number of backoffs for replication DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replicas"),
			"Number of internal replication DCP connections in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaItemsRemaining: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_items_remaining"),
			"Number of items remaining to be sent to consumer in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaItemsSent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_items_sent"),
			"Number of items per second being sent for a producer for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaProducerCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_producers"),
			"Number of replication senders for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaTotalBacklogSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_total_backlog_size"),
			"ep_dcp_replica_total_backlog_size",
			[]string{"bucket"},
			nil,
		),
		EpDcpReplicaTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_replica_total_bytes"),
			"Number of bytes per second being sent for replication DCP connections for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsBackoff: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_backoffs"),
			"Number of backoffs for views DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_view_connections"),
			"Number of views DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsItemsRemaining: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_items_remaining"),
			"Number of views items remaining to be sent",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsItemsSent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_items_sent"),
			"Number of views items sent",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsProducerCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_producers"),
			"Number of views producers",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsTotalBacklogSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_total_backlog_size"),
			"ep_dcp_views_total_backlog_size",
			[]string{"bucket"},
			nil,
		),
		EpDcpViewsTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_views_total_bytes"),
			"Number bytes per second being sent for views DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrBackoff: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_backoff"),
			"Number of backoffs for XDCR DCP connections",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_connections"),
			"Number of internal XDCR DCP connections in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrItemsRemaining: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_items_remaining"),
			"Number of items remaining to be sent to consumer in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrItemsSent: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_items_sent"),
			"Number of items per second being sent for a producer for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrProducerCount: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_producers"),
			"Number of XDCR senders for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrTotalBacklogSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_total_backlog_size"),
			"ep_dcp_xdcr_total_backlog_size",
			[]string{"bucket"},
			nil,
		),
		EpDcpXdcrTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_dcp_xdcr_total_bytes"),
			"Number of bytes per second being sent for XDCR DCP connections for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDiskqueueDrain: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_diskqueue_drain"),
			"Total number of items per second being written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDiskqueueFill: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_diskqueue_fill"),
			"Total number of items per second being put on the disk queue in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpDiskqueueItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_diskqueue_items"),
			"Total number of items waiting to be written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpFlusherTodo: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_flusher_todo"),
			"Number of items currently being written",
			[]string{"bucket"},
			nil,
		),
		EpItemCommitFailed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_item_commit_failed"),
			"Number of times a transaction failed to commit due to storage errors",
			[]string{"bucket"},
			nil,
		),
		EpKvSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_kv_size"),
			"Total amount of user data cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpMaxSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_max_size_bytes"),
			"The maximum amount of memory this bucket can use",
			[]string{"bucket"},
			nil,
		),
		EpMemHighWat: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_mem_high_wat_bytes"),
			"High water mark for auto-evictions",
			[]string{"bucket"},
			nil,
		),
		EpMemLowWat: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_mem_low_wat_bytes"),
			"Low water mark for auto-evictions",
			[]string{"bucket"},
			nil,
		),
		EpMetaDataMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_meta_data_memory"),
			"Total amount of item metadata consuming RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_non_resident"),
			"Number of non-resident items",
			[]string{"bucket"},
			nil,
		),
		EpNumOpsDelMeta: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_ops_del_meta"),
			"Number of delete operations per second for this bucket as the target for XDCR",
			[]string{"bucket"},
			nil,
		),
		EpNumOpsDelRetMeta: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_ops_del_ret_meta"),
			"Number of delRetMeta operations per second for this bucket as the target for XDCR",
			[]string{"bucket"},
			nil,
		),
		EpNumOpsGetMeta: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_ops_get_meta"),
			"Number of metadata read operations per second for this bucket as the target for XDCR",
			[]string{"bucket"},
			nil,
		),
		EpNumOpsSetMeta: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_ops_set_meta"),
			"Number of set operations per second for this bucket as the target for XDCR",
			[]string{"bucket"},
			nil,
		),
		EpNumOpsSetRetMeta: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_ops_set_ret_meta"),
			"Number of setRetMeta operations per second for this bucket as the target for XDCR",
			[]string{"bucket"},
			nil,
		),
		EpNumValueEjects: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_num_value_ejects"),
			"Total number of items per second being ejected to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpOomErrors: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_oom_errors"),
			"Number of times unrecoverable OOMs happened while processing operations",
			[]string{"bucket"},
			nil,
		),
		EpOpsCreate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_ops_create"),
			"Total number of new items being inserted into this bucket",
			[]string{"bucket"},
			nil,
		),
		EpOpsUpdate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_ops_update"),
			"Number of items updated on disk per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		EpOverhead: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_overhead"),
			"Extra memory used by transient data like persistence queues or checkpoints",
			[]string{"bucket"},
			nil,
		),
		EpQueueSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_queue_size"),
			"Number of items queued for storage",
			[]string{"bucket"},
			nil,
		),
		EpResidentItemsRate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_resident_items_rate"),
			"Percentage of all items cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		EpReplicaAheadExceptions: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_replica_ahead_exceptions"),
			"ep_replica_ahead_exceptions",
			[]string{"bucket"},
			nil,
		),
		EpReplicaHlcDrift: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_replica_hlc_drift"),
			"The sum of the total Absolute Drift, which is the accumulated drift observed by the vBucket. Drift is always accumulated as an absolute value.",
			[]string{"bucket"},
			nil,
		),
		EpTmpOomErrors: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_tmp_oom_errors"),
			"Number of back-offs sent per second to client SDKs due to OOM situations from this bucket",
			[]string{"bucket"},
			nil,
		),
		EpVbTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ep_vbuckets"),
			"Total number of vBuckets for this bucket",
			[]string{"bucket"},
			nil,
		),
		Evictions: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "evictions"),
			"Number of evictions",
			[]string{"bucket"},
			nil,
		),
		GetHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "get_hits"),
			"Number of get hits",
			[]string{"bucket"},
			nil,
		),
		GetMisses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "get_misses"),
			"Number of get misses",
			[]string{"bucket"},
			nil,
		),
		HibernatedRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "hibernated_requests"),
			"Number of streaming requests on port 8091 now idle",
			[]string{"bucket"},
			nil,
		),
		HibernatedWaked: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "hibernated_waked"),
			"Rate of streaming request wakeups on port 8091",
			[]string{"bucket"},
			nil,
		),
		HitRatio: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "hit_ratio"),
			"Hit ratio",
			[]string{"bucket"},
			nil,
		),
		IncrHits: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "incr_hits"),
			"Number of increment hits",
			[]string{"bucket"},
			nil,
		),
		IncrMisses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "incr_misses"),
			"Number of increment misses",
			[]string{"bucket"},
			nil,
		),
		MemActualFree: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_actual_free"),
			"Amount of RAM available on this server",
			[]string{"bucket"},
			nil,
		),
		MemActualUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_actual_used_bytes"),
			"_mem_actual_used",
			[]string{"bucket"},
			nil,
		),
		MemFree: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_free_bytes"),
			"Amount of Memory free",
			[]string{"bucket"},
			nil,
		),
		MemTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_bytes"),
			"Total amount of memory available",
			[]string{"bucket"},
			nil,
		),
		MemUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_used_bytes"),
			"Amount of memory used",
			[]string{"bucket"},
			nil,
		),
		MemUsedSys: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "mem_used_sys_bytes"),
			"_mem_used_sys",
			[]string{"bucket"},
			nil,
		),
		Misses: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "misses"),
			"Number of misses",
			[]string{"bucket"},
			nil,
		),
		Ops: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "ops"),
			"Total amount of operations per second to this bucket",
			[]string{"bucket"},
			nil,
		),
		RestRequests: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "rest_requests"),
			"Rate of http requests on port 8091",
			[]string{"bucket"},
			nil,
		),
		SwapTotal: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "swap_bytes"),
			"Total amount of swap available",
			[]string{"bucket"},
			nil,
		),
		SwapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "swap_used_bytes"),
			"Amount of swap space in use on this server",
			[]string{"bucket"},
			nil,
		),
		Timestamp: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "timestamp"),
			"_timestamp",
			[]string{"bucket"},
			nil,
		),
		VbActiveEject: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_eject"),
			"Number of items per second being ejected to disk from active vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveItmMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_itm_memory"),
			"Amount of active user data cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveMetaDataMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_meta_data_memory"),
			"Amount of active item metadata consuming RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveNum: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_num"),
			"Number of vBuckets in the active state for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_num_non_resident"),
			"Number of non resident vBuckets in the active state for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveOpsCreate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_ops_create"),
			"New items per second being inserted into active vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveOpsUpdate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_ops_update"),
			"Number of items updated on active vBucket per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_queue_age"),
			"Sum of disk queue item age in milliseconds",
			[]string{"bucket"},
			nil,
		),
		VbActiveQueueDrain: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_queue_drain"),
			"Number of active items per second being written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveQueueFill: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_queue_fill"),
			"Number of active items per second being put on the active item disk queue in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveQueueSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_queue_size"),
			"Number of active items waiting to be written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbActiveResidentItemsRatio: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_active_resident_items_ratio"),
			"Percentage of active items cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbAvgActiveQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_avg_active_queue_age"),
			"Average age in seconds of active items in the active item queue for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbAvgPendingQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_avg_pending_queue_age"),
			"Average age in seconds of pending items in the pending item queue for this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbAvgReplicaQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_avg_replica_queue_age"),
			"Average age in seconds of replica items in the replica item queue for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbAvgTotalQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_avg_total_queue_age"),
			"Average age in seconds of all items in the disk write queue for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbPendingCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_curr_items"),
			"Number of items in pending vBuckets in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingEject: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_eject"),
			"Number of items per second being ejected to disk from pending vBuckets in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingItmMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_itm_memory"),
			"Amount of pending user data cached in RAM in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingMetaDataMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_meta_data_memory"),
			"Amount of pending item metadata consuming RAM in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingNum: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_num"),
			"Number of vBuckets in the pending state for this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_num_non_resident"),
			"Number of non resident vBuckets in the pending state for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbPendingOpsCreate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_ops_create"),
			"New items per second being instead into pending vBuckets in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingOpsUpdate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_ops_update"),
			"Number of items updated on pending vBucket per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbPendingQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_queue_age"),
			"Sum of disk pending queue item age in milliseconds",
			[]string{"bucket"},
			nil,
		),
		VbPendingQueueDrain: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_queue_drain"),
			"Number of pending items per second being written to disk in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingQueueFill: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_queue_fill"),
			"Number of pending items per second being put on the pending item disk queue in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingQueueSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_queue_size"),
			"Number of pending items waiting to be written to disk in this bucket and should be transient during rebalancing",
			[]string{"bucket"},
			nil,
		),
		VbPendingResidentItemsRatio: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_pending_resident_items_ratio"),
			"Percentage of items in pending state vbuckets cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaCurrItems: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_curr_items"),
			"Number of items in replica vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaEject: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_eject"),
			"Number of items per second being ejected to disk from replica vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaItmMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_itm_memory"),
			"Amount of replica user data cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaMetaDataMemory: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_meta_data_memory"),
			"Amount of replica item metadata consuming in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaNum: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_num"),
			"Number of vBuckets in the replica state for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaNumNonResident: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_num_non_resident"),
			"_vb_replica_num_non_resident",
			[]string{"bucket"},
			nil,
		),
		VbReplicaOpsCreate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_ops_create"),
			"New items per second being inserted into replica vBuckets in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaOpsUpdate: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_ops_update"),
			"Number of items updated on replica vBucket per second for this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_queue_age"),
			"Sum of disk replica queue item age in milliseconds",
			[]string{"bucket"},
			nil,
		),
		VbReplicaQueueDrain: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_queue_drain"),
			"Number of replica items per second being written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaQueueFill: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_queue_fill"),
			"Number of replica items per second being put on the replica item disk queue in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaQueueSize: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_queue_size"),
			"Number of replica items waiting to be written to disk in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbReplicaResidentItemsRatio: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_replica_resident_items_ratio"),
			"Percentage of replica items cached in RAM in this bucket",
			[]string{"bucket"},
			nil,
		),
		VbTotalQueueAge: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "vbuckets_total_queue_age"),
			"Sum of disk queue item age in milliseconds",
			[]string{"bucket"},
			nil,
		),
		XdcOps: prometheus.NewDesc(
			prometheus.BuildFQName(FQ_NAMESPACE, subsystem, "xdc_ops"),
			"Total XDCR operations per second for this bucket",
			[]string{"bucket"},
			nil,
		),
	}
}
