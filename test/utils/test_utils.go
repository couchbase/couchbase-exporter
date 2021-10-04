package test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

func GetDescString(namespace string, subsystem string, name string, help string, labels []string) string {
	return fmt.Sprintf("Desc{fqName: \"%s%s_%s\", help: \"%s\", constLabels: {}, variableLabels: %+v}", namespace, subsystem, name, help, labels)
}

func GetMetricValue(metric io_prometheus_client.Metric) float64 {
	var gauge float64

	if metric.Gauge != nil {
		gauge = metric.Gauge.GetValue()
	}

	if metric.Counter != nil {
		gauge = metric.Counter.GetValue()
	}

	return gauge
}

func GetGaugeValue(m prometheus.Metric) (float64, error) {
	obj := new(io_prometheus_client.Metric)

	err := m.Write(obj)

	return GetMetricValue(*obj), err
}

func GetKeyspaceLabelIfPresent(m prometheus.Metric) (string, error) {
	obj := new(io_prometheus_client.Metric)

	if err := m.Write(obj); err != nil {
		return "", err
	}

	for _, label := range obj.Label {
		if *label.Name == "keyspace" {
			return label.GetValue(), nil
		}
	}

	return "", nil
}

func GetBucketIfPresent(m prometheus.Metric) (string, error) {
	obj := new(io_prometheus_client.Metric)

	if err := m.Write(obj); err != nil {
		return "", err
	}

	for _, label := range obj.Label {
		if *label.Name == "bucket" {
			return label.GetValue(), nil
		}
	}

	return "", nil
}

func GetFQNameFromDesc(desc *prometheus.Desc) string {
	i := reflect.ValueOf(*desc)
	t := i.FieldByName("fqName")

	return t.String()
}

func GetKeyFromFQName(config *objects.CollectorConfig, fqName string) string {
	for key, val := range config.Metrics {
		valFqName := GetFQNameFromDesc(val.GetPrometheusDescription(config.Namespace, config.Subsystem))
		if valFqName == fqName {
			return key
		}
	}

	return ""
}

func GetBucketStatsTestValue(key string, bucket string, stats map[string]objects.BucketStats) float64 {
	for bucketName, stat := range stats {
		if bucket != bucketName {
			continue
		}

		testValue := Last(stat.Op.Samples[key])

		if key == objects.AvgBgWaitTime {
			testValue /= 1000000
		}

		if key == objects.EpCacheMissRate {
			testValue = Min(testValue, 100)
		}

		return testValue
	}

	return 0
}

func GetName(name string, nameOverride string) string {
	n := name
	if nameOverride != "" {
		n = nameOverride
	}

	return n
}

func GetRandomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetRandomIntAsFloat64(min int, max int) float64 {
	return float64(GetRandomInt(min, max))
}

func GetRandomInt64(min int, max int) int64 {
	return int64(GetRandomInt(min, max))
}

func GetRandomInt(min int, max int) int {
	val := rand.Int()

	if val < min {
		return val + min
	}

	return val % max
}

func GetRandomFloatSlice(min float64, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = GetRandomFloat64(min, max)
	}

	return res
}

func Last(stats []float64) float64 {
	if len(stats) == 0 {
		panic(stats)
	}

	return stats[len(stats)-1]
}

func contains(s []string, t string) bool {
	for _, x := range s {
		if x == t {
			return true
		}
	}

	return false
}

func Min(x, y float64) float64 {
	if x > y {
		return y
	}

	return x
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}

	return 0.0
}

func getUptimeValue(uptime string, bitSize int) float64 {
	up, err := strconv.ParseFloat(uptime, bitSize)

	if err != nil {
		return 0
	}

	return up
}

func handleSpecialCases(key string, node objects.Node) float64 {
	switch key {
	case "healthy":
		return boolToFloat64(node.Status == "healthy")
	case "uptime":
		return getUptimeValue(node.Uptime, 64)
	case "clusterMembership":
		return boolToFloat64(node.ClusterMembership == "active")
	case "memoryTotal":
		return node.MemoryTotal
	case "memoryFree":
		return node.MemoryFree
	case "mcdMemoryAllocated":
		return node.McdMemoryAllocated
	case "mcdMemoryReserved":
		return node.McdMemoryReserved
	default:
		return 0
	}
}

func GetNodeTestValue(key string, name string, config objects.CollectorConfig, nodes objects.Nodes) float64 {
	metric := config.Metrics[key]
	if contains(metric.Labels, "node") {
		switch {
		case strings.HasPrefix(name, "systemstats_"):
			newKey := strings.Replace(name, "systemstats_", "", 1)
			return nodes.Nodes[0].SystemStats[newKey]
		case strings.HasPrefix(name, "interestingstats_"):
			newKey := strings.Replace(name, "interestingstats_", "", 1)
			return nodes.Nodes[0].InterestingStats[newKey]
		default:
			return handleSpecialCases(key, nodes.Nodes[0])
		}
	} else {
		return nodes.Counters[name]
	}
}

func GenerateServers() objects.Servers {
	server := objects.Server{
		Hostname: "localhost",
		URI:      "/pools/default/buckets/wawa-bucket/nodes/%5B%3A%3A1%5D%3A8091",
		Stats: map[string]string{
			"uri": "/pools/default/buckets/wawa-bucket/nodes/%5B%3A%3A1%5D%3A8091/stats",
		},
	}

	return objects.Servers{
		Servers: []objects.Server{server},
	}
}

func GenerateNodes(name string, nodes []objects.Node) objects.Nodes {
	cluster := objects.Nodes{
		Name:  name,
		Nodes: nodes,
		Counters: map[string]float64{
			objects.RebalanceSuccess:        GetRandomIntAsFloat64(0, 99999),
			objects.RebalanceStart:          GetRandomIntAsFloat64(0, 99999),
			objects.RebalanceFail:           GetRandomIntAsFloat64(0, 99999),
			objects.RebalanceStop:           GetRandomIntAsFloat64(0, 99999),
			objects.FailoverNode:            GetRandomIntAsFloat64(0, 99999),
			objects.Failover:                GetRandomIntAsFloat64(0, 99999),
			objects.FailoverComplete:        GetRandomIntAsFloat64(0, 99999),
			objects.FailoverIncomplete:      GetRandomIntAsFloat64(0, 99999),
			objects.GracefulFailoverStart:   GetRandomIntAsFloat64(0, 99999),
			objects.GracefulFailoverSuccess: GetRandomIntAsFloat64(0, 99999),
			objects.GracefulFailoverFail:    GetRandomIntAsFloat64(0, 99999),
		},
	}

	return cluster
}

func GenerateNode() objects.Node {
	node := objects.Node{
		SystemStats: map[string]float64{
			objects.CPUUtilizationRate: GetRandomFloat64(0, 99999),
			objects.SwapTotal:          GetRandomFloat64(0, 99999),
			objects.SwapUsed:           GetRandomFloat64(0, 99999),
			objects.MemTotal:           GetRandomFloat64(0, 99999),
			objects.MemFree:            GetRandomFloat64(0, 99999),
		},
		InterestingStats: map[string]float64{
			objects.CmdGet:                                 GetRandomFloat64(0, 99999),
			objects.CouchDocsActualDiskSize:                GetRandomFloat64(0, 99999),
			objects.CouchDocsDataSize:                      GetRandomFloat64(0, 99999),
			objects.CouchSpatialDataSize:                   GetRandomFloat64(0, 99999),
			objects.CouchSpatialDiskSize:                   GetRandomFloat64(0, 99999),
			objects.CouchViewsActualDiskSize:               GetRandomFloat64(0, 99999),
			objects.CouchViewsDataSize:                     GetRandomFloat64(0, 99999),
			objects.CurrItems:                              GetRandomFloat64(0, 99999),
			objects.CurrItemsTot:                           GetRandomFloat64(0, 99999),
			objects.EpBgFetched:                            GetRandomFloat64(0, 99999),
			objects.GetHits:                                GetRandomFloat64(0, 99999),
			objects.InterestingStatsMemUsed:                GetRandomFloat64(0, 99999),
			objects.Ops:                                    GetRandomFloat64(0, 99999),
			objects.InterestingStatsVbActiveNumNonResident: GetRandomFloat64(0, 99999),
			objects.VbReplicaCurrItems:                     GetRandomFloat64(0, 99999),
		},
		Uptime:               "34234323",
		MemoryTotal:          GetRandomFloat64(0, 999999),
		MemoryFree:           GetRandomFloat64(0, 999999),
		CouchAPIBaseHTTPS:    "",
		CouchAPIBase:         "",
		McdMemoryReserved:    GetRandomFloat64(0, 999999),
		McdMemoryAllocated:   GetRandomFloat64(0, 99999),
		Replication:          GetRandomIntAsFloat64(0, 9999),
		ClusterMembership:    "active",
		RecoveryType:         "",
		Status:               "healthy",
		OtpNode:              "",
		ThisNode:             true,
		OtpCookie:            "",
		Hostname:             "localhost",
		ClusterCompatibility: 5,
		Version:              "",
		Os:                   "",
		CPUCount:             5,
		Ports:                &objects.Ports{},
		Services:             []string{""},
		AlternateAddresses:   &objects.AlternateAddressesExternal{},
	}

	return node
}

func GenerateTasks() []objects.Task {
	tasks := make([]objects.Task, 0)
	rebalance := objects.Task{
		Type:     "rebalance",
		Progress: GetRandomIntAsFloat64(0, 100),
		PerNode: map[string]struct {
			Progress float64 "json:\"progress,omitempty\""
		}{
			"wawa-node": {
				Progress: GetRandomIntAsFloat64(0, 100),
			},
		},
	}
	xdcr := objects.Task{
		Type:           "xdcr",
		ChangesLeft:    GetRandomInt64(0, 99999),
		DocsChecked:    GetRandomInt64(0, 99999),
		DocsWritten:    GetRandomInt64(0, 99999),
		PauseRequested: false,
		Errors:         make([]interface{}, 5),
		Source:         "Foo",
		Target:         "Bar",
		DetailedProgress: struct {
			Bucket       string "json:\"bucket,omitempty\""
			BucketNumber int    "json:\"bucketNumber,omitempty\""
			BucketCount  int    "json:\"bucketCount,omitempty\""
			PerNode      map[string]struct {
				Ingoing  objects.NodeProgress "json:\"ingoing,omitempty\""
				Outgoing objects.NodeProgress "json:\"outgoing,omitempty\""
			} "json:\"perNode,omitempty\""
		}{
			Bucket:       "wawa-bucket",
			BucketNumber: 3,
			BucketCount:  3,
			PerNode: map[string]struct {
				Ingoing  objects.NodeProgress "json:\"ingoing,omitempty\""
				Outgoing objects.NodeProgress "json:\"outgoing,omitempty\""
			}{
				"wawa-node": {
					Ingoing: objects.NodeProgress{
						DocsTotal:           GetRandomInt64(0, 99999),
						DocsTransferred:     GetRandomInt64(0, 99999),
						ActiveVBucketsLeft:  GetRandomInt64(0, 3),
						ReplicaVBucketsLeft: GetRandomInt64(0, 5),
					},
				},
			},
		},
	}
	bucketCompaction := objects.Task{
		Type:     "bucket_compaction",
		Progress: GetRandomIntAsFloat64(0, 100),
		Bucket:   "wawa-bucket",
	}
	clusterLogCollection := objects.Task{
		Type:     "clusterLogsCollection",
		Progress: GetRandomIntAsFloat64(0, 100),
	}
	tasks = append(tasks, clusterLogCollection, rebalance, xdcr, bucketCompaction)

	return tasks
}

func getTask(tasks []objects.Task, taskType string) objects.Task {
	task := objects.Task{}

	for _, t := range tasks {
		if t.Type == taskType {
			task = t
		}
	}

	return task
}

func GetTaskTestValue(key string, name string, tasks []objects.Task) float64 {
	switch key {
	case "compacting":
		return getTask(tasks, "bucket_compaction").Progress
	case "clusterLogsCollection":
		return getTask(tasks, "clusterLogsCollection").Progress
	case "rebalance":
		return getTask(tasks, "rebalance").Progress
	case "rebalancePerNode":
		return getTask(tasks, "rebalance").PerNode["wawa-node"].Progress
	case "xdcrChangesLeft":
		return float64(getTask(tasks, "xdcr").ChangesLeft)
	case "xdcrDocsChecked":
		return float64(getTask(tasks, "xdcr").DocsChecked)
	case "xdcrDocsWritten":
		return float64(getTask(tasks, "xdcr").DocsWritten)
	case "xdcrPaused":
		return boolToFloat64(getTask(tasks, "xdcr").PauseRequested)
	case "xdcrErrors":
		return float64(len(getTask(tasks, "xdcr").Errors))
	case "progressDocsTotal":
		return float64(getTask(tasks, "xdcr").DetailedProgress.PerNode["wawa-node"].Ingoing.DocsTotal)
	case "progressDocsTransferred":
		return float64(getTask(tasks, "xdcr").DetailedProgress.PerNode["wawa-node"].Ingoing.DocsTransferred)
	case "progressActiveVBucketsLeft":
		return float64(getTask(tasks, "xdcr").DetailedProgress.PerNode["wawa-node"].Ingoing.ActiveVBucketsLeft)
	case "progressReplicaVBucketsLeft":
		return float64(getTask(tasks, "xdcr").DetailedProgress.PerNode["wawa-node"].Ingoing.ReplicaVBucketsLeft)
	default:
		return 0
	}
}

func GenerateBucketInfo(name string) objects.BucketInfo {
	singleBucket := objects.BucketInfo{
		Name:              name,
		BucketType:        "doesn't matter!",
		AuthType:          "saml",
		ProxyPort:         8091,
		URI:               "some uri",
		StreamingURI:      "some streaming uri",
		LocalRandomKeyURI: "some random key uri",
		Nodes:             []objects.Node{},
		BucketBasicStats: map[string]float64{
			objects.QuotaPercentUsed:       GetRandomFloat64(1, 99999),
			objects.OpsPerSec:              GetRandomFloat64(1, 99999),
			objects.DiskFetches:            GetRandomFloat64(1, 99999),
			objects.ItemCount:              GetRandomFloat64(1, 99999),
			objects.DiskUsed:               GetRandomFloat64(1, 99999),
			objects.DataUsed:               GetRandomFloat64(1, 99999),
			objects.MemUsed:                GetRandomFloat64(1, 99999),
			objects.VbActiveNumNonResident: GetRandomFloat64(1, 99999),
		},
	}

	return singleBucket
}

func GenerateFTS() objects.FTS {
	fts := objects.FTS{
		Op: struct {
			Samples map[string][]float64 "json:\"samples\""
		}{
			Samples: map[string][]float64{
				objects.FtsCurrBatchesBlockedByHerder:   GetRandomFloatSlice(0, 1000, 10),
				objects.FtsNumBytesUsedRAM:              GetRandomFloatSlice(0, 1000, 10),
				objects.FtsTotalQueriesRejectedByHerder: GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return fts
}

func GenerateQuery() objects.Query {
	query := objects.Query{
		Op: struct {
			Samples map[string][]float64 "json:\"samples\""
		}{
			Samples: map[string][]float64{
				objects.QueryAvgReqTime:      GetRandomFloatSlice(0, 1000, 10),
				objects.QueryAvgSvcTime:      GetRandomFloatSlice(0, 1000, 10),
				objects.QueryAvgResponseSize: GetRandomFloatSlice(0, 1000, 10),
				objects.QueryAvgResultCount:  GetRandomFloatSlice(0, 1000, 10),
				objects.QueryActiveRequests:  GetRandomFloatSlice(0, 1000, 10),
				objects.QueryErrors:          GetRandomFloatSlice(0, 1000, 10),
				objects.QueryInvalidRequests: GetRandomFloatSlice(0, 1000, 10),
				objects.QueryQueuedRequests:  GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequestTime:     GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequests:        GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequests1000Ms:  GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequests250Ms:   GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequests5000Ms:  GetRandomFloatSlice(0, 1000, 10),
				objects.QueryRequests500Ms:   GetRandomFloatSlice(0, 1000, 10),
				objects.QueryResultCount:     GetRandomFloatSlice(0, 1000, 10),
				objects.QueryResultSize:      GetRandomFloatSlice(0, 1000, 10),
				objects.QuerySelects:         GetRandomFloatSlice(0, 1000, 10),
				objects.QueryServiceTime:     GetRandomFloatSlice(0, 1000, 10),
				objects.QueryWarnings:        GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return query
}

func GenerateIndex() objects.Index {
	index := objects.Index{
		Op: struct {
			Samples map[string][]float64 "json:\"samples\""
		}{
			Samples: map[string][]float64{
				objects.IndexMemoryQuota:  GetRandomFloatSlice(0, 1000, 10),
				objects.IndexMemoryUsed:   GetRandomFloatSlice(0, 1000, 10),
				objects.IndexRAMPercent:   GetRandomFloatSlice(0, 1000, 10),
				objects.IndexRemainingRAM: GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return index
}

func GenerateEventing() objects.Eventing {
	eventing := objects.Eventing{
		Op: struct {
			Samples map[string][]float64 "json:\"samples\""
		}{
			Samples: map[string][]float64{
				objects.EventingCheckpointFailureCount:     GetRandomFloatSlice(0, 1000, 10),
				objects.EventingBucketOpExceptionCount:     GetRandomFloatSlice(0, 1000, 10),
				objects.EventingDcpBacklog:                 GetRandomFloatSlice(0, 1000, 10),
				objects.EventingFailedCount:                GetRandomFloatSlice(0, 1000, 10),
				objects.EventingN1QlOpExceptionCount:       GetRandomFloatSlice(0, 1000, 10),
				objects.EventingOnDeleteFailure:            GetRandomFloatSlice(0, 1000, 10),
				objects.EventingOnDeleteSuccess:            GetRandomFloatSlice(0, 1000, 10),
				objects.EventingOnUpdateFailure:            GetRandomFloatSlice(0, 1000, 10),
				objects.EventingOnUpdateSuccess:            GetRandomFloatSlice(0, 1000, 10),
				objects.EventingProcessedCount:             GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestBucketOpExceptionCount: GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestCheckpointFailureCount: GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestDcpBacklog:             GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestFailedCount:            GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestN1QlOpExceptionCount:   GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestOnDeleteFailure:        GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestOnDeleteSuccess:        GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestOnUpdateFailure:        GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestOnUpdateSuccess:        GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestProcessedCount:         GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTestTimeoutCount:           GetRandomFloatSlice(0, 1000, 10),
				objects.EventingTimeoutCount:               GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return eventing
}

func GenerateAnalytics() objects.Analytics {
	anal := objects.Analytics{
		Op: struct {
			Samples      map[string][]float64 "json:\"samples\""
			SamplesCount int                  "json:\"samplesCount\""
			IsPersistent bool                 "json:\"isPersistent\""
			LastTStamp   int64                "json:\"lastTStamp\""
			Interval     int                  "json:\"interval\""
		}{
			SamplesCount: 10,
			IsPersistent: true,
			Interval:     5,
			LastTStamp:   500,
			Samples: map[string][]float64{
				objects.CbasDiskUsed:          GetRandomFloatSlice(0, 1000, 10),
				objects.CbasGcCount:           GetRandomFloatSlice(0, 1000, 10),
				objects.CbasGcTime:            GetRandomFloatSlice(0, 1000, 10),
				objects.CbasHeapUsed:          GetRandomFloatSlice(0, 1000, 10),
				objects.CbasIoReads:           GetRandomFloatSlice(0, 1000, 10),
				objects.CbasIoWrites:          GetRandomFloatSlice(0, 1000, 10),
				objects.CbasSystemLoadAverage: GetRandomFloatSlice(0, 1000, 10),
				objects.CbasThreadCount:       GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return anal
}

func GenerateBucketStatSamples() map[string]interface{} {
	samples := GenerateBucketStats().Op.Samples

	realSamples := make(map[string]interface{}, len(samples))

	for k, v := range samples {
		realSamples[k] = v
	}

	return realSamples
}

func GenerateBucketStats() objects.BucketStats {
	stats := objects.BucketStats{
		Op: struct {
			Samples      map[string][]float64 "json:\"samples\""
			SamplesCount float64              "json:\"samplesCount\""
			IsPersistent bool                 "json:\"isPersistent\""
			LastTStamp   float64              "json:\"lastTStamp\""
			Interval     float64              "json:\"interval\""
		}{
			Samples: map[string][]float64{
				objects.CouchTotalDiskSize:                  GetRandomFloatSlice(0, 1000, 10),
				objects.CouchDocsFragmentation:              GetRandomFloatSlice(0, 1000, 10),
				objects.CouchViewsFragmentation:             GetRandomFloatSlice(0, 1000, 10),
				objects.HitRatio:                            GetRandomFloatSlice(0, 1000, 10),
				objects.EpCacheMissRate:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpResidentItemsRate:                 GetRandomFloatSlice(0, 1000, 10),
				objects.VbAvgActiveQueueAge:                 GetRandomFloatSlice(0, 1000, 10),
				objects.VbAvgReplicaQueueAge:                GetRandomFloatSlice(0, 1000, 10),
				objects.VbAvgPendingQueueAge:                GetRandomFloatSlice(0, 1000, 10),
				objects.VbAvgTotalQueueAge:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveResidentItemsRatio:          GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaResidentItemsRatio:         GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingResidentItemsRatio:         GetRandomFloatSlice(0, 1000, 10),
				objects.AvgDiskUpdateTime:                   GetRandomFloatSlice(0, 1000, 10),
				objects.AvgDiskCommitTime:                   GetRandomFloatSlice(0, 1000, 10),
				objects.AvgBgWaitTime:                       GetRandomFloatSlice(0, 1000, 10),
				objects.AvgActiveTimestampDrift:             GetRandomFloatSlice(0, 1000, 10),
				objects.AvgReplicaTimestampDrift:            GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesCount:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesItemsRemaining:     GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesProducerCount:      GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesTotalBacklogSize:   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesItemsSent:          GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesTotalBytes:         GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsIndexesBackoff:            GetRandomFloatSlice(0, 1000, 10),
				objects.BgWaitCount:                         GetRandomFloatSlice(0, 1000, 10),
				objects.BgWaitTotal:                         GetRandomFloatSlice(0, 1000, 10),
				objects.BytesRead:                           GetRandomFloatSlice(0, 1000, 10),
				objects.BytesWritten:                        GetRandomFloatSlice(0, 1000, 10),
				objects.CasBadval:                           GetRandomFloatSlice(0, 1000, 10),
				objects.CasHits:                             GetRandomFloatSlice(0, 1000, 10),
				objects.CasMisses:                           GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCmdGet:                   GetRandomFloatSlice(0, 1000, 10),
				objects.CmdSet:                              GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchDocsActualDiskSize:  GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchDocsDataSize:        GetRandomFloatSlice(0, 1000, 10),
				objects.CouchDocsDiskSize:                   GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchSpatialDataSize:     GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchSpatialDiskSize:     GetRandomFloatSlice(0, 1000, 10),
				objects.CouchSpatialOps:                     GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchViewsActualDiskSize: GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCouchViewsDataSize:       GetRandomFloatSlice(0, 1000, 10),
				objects.CouchViewsDiskSize:                  GetRandomFloatSlice(0, 1000, 10),
				objects.CouchViewsOps:                       GetRandomFloatSlice(0, 1000, 10),
				objects.CurrConnections:                     GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCurrItems:                GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCurrItemsTot:             GetRandomFloatSlice(0, 1000, 10),
				objects.DecrHits:                            GetRandomFloatSlice(0, 1000, 10),
				objects.DecrMisses:                          GetRandomFloatSlice(0, 1000, 10),
				objects.DeleteHits:                          GetRandomFloatSlice(0, 1000, 10),
				objects.DeleteMisses:                        GetRandomFloatSlice(0, 1000, 10),
				objects.DiskCommitCount:                     GetRandomFloatSlice(0, 1000, 10),
				objects.DiskCommitTotal:                     GetRandomFloatSlice(0, 1000, 10),
				objects.DiskUpdateCount:                     GetRandomFloatSlice(0, 1000, 10),
				objects.DiskUpdateTotal:                     GetRandomFloatSlice(0, 1000, 10),
				objects.DiskWriteQueue:                      GetRandomFloatSlice(0, 1000, 10),
				objects.EpActiveAheadExceptions:             GetRandomFloatSlice(0, 1000, 10),
				objects.EpActiveHlcDrift:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpActiveHlcDriftCount:               GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsEpBgFetched:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpClockCasDriftThresholdExceeded:    GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2IBackoff:                      GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2ICount:                        GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2IItemsRemaining:               GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2IItemsSent:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2IProducerCount:                GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2ITotalBacklogSize:             GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcp2ITotalBytes:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsBackoff:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsCount:                       GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsItemsRemaining:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsItemsSent:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsProducerCount:               GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsTotalBacklogSize:            GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpFtsTotalBytes:                  GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherBackoff:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherCount:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherItemsRemaining:            GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherItemsSent:                 GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherProducerCount:             GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherTotalBacklogSize:          GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpOtherTotalBytes:                GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaBackoff:                 GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaCount:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaItemsRemaining:          GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaItemsSent:               GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaProducerCount:           GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaTotalBacklogSize:        GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpReplicaTotalBytes:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsBackoff:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsCount:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsItemsRemaining:            GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsItemsSent:                 GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsProducerCount:             GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsTotalBacklogSize:          GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpViewsTotalBytes:                GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrBackoff:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrCount:                      GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrItemsRemaining:             GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrItemsSent:                  GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrProducerCount:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrTotalBacklogSize:           GetRandomFloatSlice(0, 1000, 10),
				objects.EpDcpXdcrTotalBytes:                 GetRandomFloatSlice(0, 1000, 10),
				objects.EpDiskqueueDrain:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpDiskqueueFill:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpDiskqueueItems:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpFlusherTodo:                       GetRandomFloatSlice(0, 1000, 10),
				objects.EpItemCommitFailed:                  GetRandomFloatSlice(0, 1000, 10),
				objects.EpKvSize:                            GetRandomFloatSlice(0, 1000, 10),
				objects.EpMaxSize:                           GetRandomFloatSlice(0, 1000, 10),
				objects.EpMemHighWat:                        GetRandomFloatSlice(0, 1000, 10),
				objects.EpMemLowWat:                         GetRandomFloatSlice(0, 1000, 10),
				objects.EpMetaDataMemory:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumNonResident:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumOpsDelMeta:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumOpsDelRetMeta:                  GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumOpsGetMeta:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumOpsSetMeta:                     GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumOpsSetRetMeta:                  GetRandomFloatSlice(0, 1000, 10),
				objects.EpNumValueEjects:                    GetRandomFloatSlice(0, 1000, 10),
				objects.EpOomErrors:                         GetRandomFloatSlice(0, 1000, 10),
				objects.EpOpsCreate:                         GetRandomFloatSlice(0, 1000, 10),
				objects.EpOpsUpdate:                         GetRandomFloatSlice(0, 1000, 10),
				objects.EpOverhead:                          GetRandomFloatSlice(0, 1000, 10),
				objects.EpQueueSize:                         GetRandomFloatSlice(0, 1000, 10),
				objects.EpReplicaAheadExceptions:            GetRandomFloatSlice(0, 1000, 10),
				objects.EpReplicaHlcDrift:                   GetRandomFloatSlice(0, 1000, 10),
				objects.EpReplicaHlcDriftCount:              GetRandomFloatSlice(0, 1000, 10),
				objects.EpTmpOomErrors:                      GetRandomFloatSlice(0, 1000, 10),
				objects.EpVbTotal:                           GetRandomFloatSlice(0, 1000, 10),
				objects.Evictions:                           GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsGetHits:                  GetRandomFloatSlice(0, 1000, 10),
				objects.GetMisses:                           GetRandomFloatSlice(0, 1000, 10),
				objects.IncrHits:                            GetRandomFloatSlice(0, 1000, 10),
				objects.IncrMisses:                          GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsMemUsed:                  GetRandomFloatSlice(0, 1000, 10),
				objects.Misses:                              GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsOps:                      GetRandomFloatSlice(0, 1000, 10),
				objects.XdcOps:                              GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveEject:                       GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveItmMemory:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveMetaDataMemory:              GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveNum:                         GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsVbActiveNumNonResident:   GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveOpsCreate:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveOpsUpdate:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveQueueAge:                    GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveQueueDrain:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveQueueFill:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbActiveQueueSize:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingCurrItems:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingEject:                      GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingItmMemory:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingMetaDataMemory:             GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingNum:                        GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingNumNonResident:             GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingOpsCreate:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingOpsUpdate:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingQueueAge:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingQueueDrain:                 GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingQueueFill:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbPendingQueueSize:                  GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsVbReplicaCurrItems:       GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaEject:                      GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaItmMemory:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaMetaDataMemory:             GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaNum:                        GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaNumNonResident:             GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaOpsCreate:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaOpsUpdate:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaQueueAge:                   GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaQueueDrain:                 GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaQueueFill:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbReplicaQueueSize:                  GetRandomFloatSlice(0, 1000, 10),
				objects.VbTotalQueueAge:                     GetRandomFloatSlice(0, 1000, 10),
				objects.CPUIdleMs:                           GetRandomFloatSlice(0, 1000, 10),
				objects.CPULocalMs:                          GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsCPUUtilizationRate:       GetRandomFloatSlice(0, 1000, 10),
				objects.HibernatedRequests:                  GetRandomFloatSlice(0, 1000, 10),
				objects.HibernatedWaked:                     GetRandomFloatSlice(0, 1000, 10),
				objects.MemActualFree:                       GetRandomFloatSlice(0, 1000, 10),
				objects.MemActualUsed:                       GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsMemFree:                  GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsMemTotal:                 GetRandomFloatSlice(0, 1000, 10),
				objects.MemUsedSys:                          GetRandomFloatSlice(0, 1000, 10),
				objects.RestRequests:                        GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsSwapTotal:                GetRandomFloatSlice(0, 1000, 10),
				objects.BucketStatsSwapUsed:                 GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasBackoff:          GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasItemsRemaining:   GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpTotalBytes:           GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasTotalBacklogSize: GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDataWriteFailed:         GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDataReadFailed:          GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasProducerCount:    GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasCount:            GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDEpDcpCbasItemsSent:        GetRandomFloatSlice(0, 1000, 10),
				objects.DEPRECATEDVbActiveQuueItems:         GetRandomFloatSlice(0, 1000, 10),
			},
		},
	}

	return stats
}

func GenerateBucket(name string) objects.BucketInfo {
	singleBucket := objects.BucketInfo{
		Name:              name,
		BucketType:        "doesn't matter!",
		AuthType:          "saml",
		ProxyPort:         8091,
		URI:               "some uri",
		StreamingURI:      "some streaming uri",
		LocalRandomKeyURI: "some random key uri",
		Nodes:             []objects.Node{},
		BucketBasicStats: map[string]float64{
			objects.QuotaPercentUsed:       GetRandomFloat64(0, 1000),
			objects.OpsPerSec:              GetRandomFloat64(0, 1000),
			objects.DiskFetches:            GetRandomFloat64(0, 1000),
			objects.ItemCount:              GetRandomFloat64(0, 1000),
			objects.DiskUsed:               GetRandomFloat64(0, 1000),
			objects.DataUsed:               GetRandomFloat64(0, 1000),
			objects.MemUsed:                GetRandomFloat64(0, 1000),
			objects.VbActiveNumNonResident: GetRandomFloat64(0, 1000),
		},
	}

	return singleBucket
}

func GenerateIndexerStats() map[string]map[string]interface{} {
	stats := map[string]map[string]interface{}{
		"indexer": {},
		"mybucket:keyspace": {
			objects.IndexDocsIndexed:          GetRandomFloat64(0, 1000),
			objects.IndexItemsCount:           GetRandomFloat64(0, 1000),
			objects.IndexFragPercent:          GetRandomFloat64(0, 1000),
			objects.IndexNumDocsPendingQueued: GetRandomFloat64(0, 1000),
			objects.IndexNumRequests:          GetRandomFloat64(0, 1000),
			objects.IndexCacheMisses:          GetRandomFloat64(0, 1000),
			objects.IndexCacheHits:            GetRandomFloat64(0, 1000),
			objects.IndexCacheHitPercent:      GetRandomFloat64(0, 1000),
			objects.IndexNumRowsReturned:      GetRandomFloat64(0, 1000),
			objects.IndexResidentPercent:      GetRandomFloat64(0, 1000),
			objects.IndexAvgScanLatency:       GetRandomFloat64(0, 1000),
		},
	}

	return stats
}
