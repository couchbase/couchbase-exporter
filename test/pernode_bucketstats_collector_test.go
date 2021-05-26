package test

import (
	"testing"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/config"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/test/mocks"
	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	metricPrefix = "cbpernode_bucketstats_"
)

func TestPerNodeBucketStatsReturnsDownIfCantGetCurrentNode(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})

	mockClient.EXPECT().Nodes().Times(1).Return(Nodes, ErrorDummy)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 0, ""))
}

func TestPerNodeBucketStatsReturnsDownIfCantGetClusterName(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})

	mockClient.EXPECT().Nodes().Times(1).Return(Nodes, nil)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrorDummy)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 0, Node.Hostname))
}

func TestPerNodeBucketStatsReturnsDownIfCantGetClusterBalanceStatus(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})

	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)
	first := mockClient.EXPECT().Nodes().Times(1).Return(Nodes, nil)
	mockClient.EXPECT().Nodes().Times(1).Return(Nodes, ErrorDummy).After(first)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 0, Node.Hostname))
}

func TestPerNodeBucketStatsReturnsDownIfCantGetBuckets(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})

	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)
	mockClient.EXPECT().Nodes().Times(2).Return(Nodes, nil)

	buckets := make([]objects.BucketInfo, 0)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, ErrorDummy)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 0, Node.Hostname))
}

func TestPerNodeBucketStatsReturnsDownIfCantGetBucketStats(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})

	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)
	mockClient.EXPECT().Nodes().Times(2).Return(Nodes, nil)

	buckets := []objects.BucketInfo{test.GenerateBucket("wawa-bucket")}
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	mockClient.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(ErrorDummy)

	servers := test.GenerateServers()
	mockClient.EXPECT().Servers(gomock.Any()).Times(1).Return(servers, nil)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 0, Node.Hostname))
}

func TestPerNodeBucketStatsReturnsUp(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	metrics := make(map[string]objects.MetricInfo)

	for key, val := range defaultConfig.Collectors.PerNodeBucketStats.Metrics {
		val.Enabled = false
		metrics[key] = val
	}

	defaultConfig.Collectors.PerNodeBucketStats.Metrics = metrics

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})
	stats := objects.PerNodeBucketStats{
		Op: struct {
			Samples      map[string]interface{} "json:\"samples\""
			SamplesCount int                    "json:\"samplesCount\""
			IsPersistent bool                   "json:\"isPersistent\""
			LastTStamp   int64                  "json:\"lastTStamp\""
			Interval     int                    "json:\"interval\""
		}{
			Samples: test.GenerateBucketStatSamples(),
		},
	}

	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)
	mockClient.EXPECT().Nodes().Times(2).Return(Nodes, nil)

	buckets := []objects.BucketInfo{test.GenerateBucket("wawa-bucket")}
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	mockClient.EXPECT().Get(gomock.Any(), gomock.Any()).SetArg(1, stats).Return(nil).Times(1)

	servers := test.GenerateServers()
	mockClient.EXPECT().Servers(gomock.Any()).Times(1).Return(servers, nil)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 1, Node.Hostname))
	assert.True(t, mockSetter.TestMetricGreaterThanOrEqual(metricPrefix+objects.DefaultScrapeDurationMetric, 0, Node.Hostname))
}

func TestPerNodeBucketStatsReturnsCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockSetter := mocks.NewMockSetter()
	Node := test.GenerateNode()
	Nodes := test.GenerateNodes("dummy-cluster", []objects.Node{Node})
	stats := objects.PerNodeBucketStats{
		Op: struct {
			Samples      map[string]interface{} "json:\"samples\""
			SamplesCount int                    "json:\"samplesCount\""
			IsPersistent bool                   "json:\"isPersistent\""
			LastTStamp   int64                  "json:\"lastTStamp\""
			Interval     int                    "json:\"interval\""
		}{
			Samples: test.GenerateBucketStatSamples(),
		},
	}

	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)
	mockClient.EXPECT().Nodes().Times(2).Return(Nodes, nil)

	buckets := []objects.BucketInfo{test.GenerateBucket("wawa-bucket")}
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	mockClient.EXPECT().Get(gomock.Any(), gomock.Any()).SetArg(1, stats).Return(nil).Times(1)

	servers := test.GenerateServers()
	mockClient.EXPECT().Servers(gomock.Any()).Times(1).Return(servers, nil)

	testCollector := collectors.NewPerNodeBucketStatsCollector(mockClient, defaultConfig.Collectors.PerNodeBucketStats)
	testCollector.Setter = &mockSetter

	testCollector.CollectMetrics()

	assert.True(t, mockSetter.TestMetric(metricPrefix+objects.DefaultUptimeMetric, 1, Node.Hostname))
	assert.True(t, mockSetter.TestMetricGreaterThanOrEqual(metricPrefix+objects.DefaultScrapeDurationMetric, 0, Node.Hostname))

	for _, value := range defaultConfig.Collectors.PerNodeBucketStats.Metrics {
		sample, ok := stats.Op.Samples[value.Name].([]float64)
		if sample == nil || !ok {
			continue
		}

		name := value.Name
		if value.NameOverride != "" {
			name = value.NameOverride
		}

		assert.True(t,
			mockSetter.TestMetric(
				defaultConfig.Collectors.PerNodeBucketStats.Namespace+defaultConfig.Collectors.PerNodeBucketStats.Subsystem+"_"+name,
				test.Last(sample),
				"wawa-bucket", Node.Hostname, "dummy-cluster",
			),
			value.Name, test.Last(sample),
		)
	}
}
