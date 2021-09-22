package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/config"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/test/mocks"
	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestBucketStatsDescribeReturnsAppropriateValuesBasedOnDefaultConfig(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketStats.Namespace,
			defaultConfig.Collectors.BucketStats.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketStats.Namespace,
			defaultConfig.Collectors.BucketStats.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.BucketStats.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.BucketStats.Namespace,
				defaultConfig.Collectors.BucketStats.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewBucketStatsCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan *prometheus.Desc, 9)

	defer close(c)

	go testCollector.Describe(c)

	count := 0

	for {
		select {
		case x := <-c:
			found := false

			for _, check := range possibleValues {
				if x.String() == check {
					found = true
					break
				}
			}

			assert.True(t, found, x.String())
			count++
		case <-time.After(1 * time.Second):
			if count >= len(possibleValues) {
				return
			}
		}
	}
}

func TestBucketStatsDescribeReturnsAppropriateValuesWithNameOverride(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	// set up the name overrides
	metrics := make(map[string]objects.MetricInfo)
	now := time.Now()

	for key, val := range defaultConfig.Collectors.BucketStats.Metrics {
		val.NameOverride = fmt.Sprintf("%s-%s", val.Name, now.Format("YYYY-MM-DD"))
		metrics[key] = val
	}

	defaultConfig.Collectors.BucketStats.Metrics = metrics
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketStats.Namespace,
			defaultConfig.Collectors.BucketStats.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketStats.Namespace,
			defaultConfig.Collectors.BucketStats.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.BucketStats.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.BucketStats.Namespace,
				defaultConfig.Collectors.BucketStats.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan *prometheus.Desc, 9)

	defer close(c)

	go testCollector.Describe(c)

	count := 0

	for {
		select {
		case x := <-c:
			found := false

			for _, check := range possibleValues {
				if x.String() == check {
					found = true
					break
				}
			}

			assert.True(t, found, x.String())
			count++
		case <-time.After(1 * time.Second):
			if count >= len(possibleValues) {
				return
			}
		}
	}
}

func TestBucketStatsCollectReturnsDownIfClientReturnsError(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrDummy)
	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestBucketStatsCollectReturnsDownIfClientReturnsErrorOnBuckets(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, ErrDummy)

	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestBucketStatsCollectReturnsDownIfClientReturnsErrorOnBucketStats(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	singleBucket := test.GenerateBucket("wawa-bucket")
	buckets = append(buckets, singleBucket)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	stats := test.GenerateBucketStats()

	mockClient.EXPECT().BucketStats(singleBucket.Name).Times(1).Return(stats, ErrDummy)

	testCollector := collectors.NewBucketStatsCollector(mockClient, defaultConfig.Collectors.BucketStats)

	c := make(chan prometheus.Metric, 10)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestBucketStatsCollectReturnsUpWithNoErrors(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	testCollector := collectors.NewBucketStatsCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan prometheus.Metric, 2)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		if strings.Contains(m.Desc().String(), "Couchbase cluster API is responding") {
			gauge, err := test.GetGaugeValue(m)
			assert.Nil(t, err)
			assert.Equal(t, 1.0, gauge)
		}
	}
}

func TestBucketStatsCollectReturnsOneOfEachMetricForASingleBucketWithCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	singleBucket := test.GenerateBucket("wawa-bucket")
	buckets = append(buckets, singleBucket)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	stats := test.GenerateBucketStats()

	mockClient.EXPECT().BucketStats(singleBucket.Name).Times(1).Return(stats, nil)

	testCollector := collectors.NewBucketStatsCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan prometheus.Metric, 9)
	count := 0

	defer close(c)

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbbucketstat_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbbucketstat_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.BucketStats, fqName)
				name := defaultConfig.Collectors.BucketStats.Metrics[key].Name
				bucketName, err := test.GetBucketIfPresent(m)

				assert.Nil(t, err)
				assert.Equal(t, singleBucket.Name, bucketName)

				gauge, err := test.GetGaugeValue(m)

				testValue := test.Last(stats.Op.Samples[name])

				if key == "AvgBgWaitTime" {
					testValue /= 1000000
				}

				if key == "EpCacheMissRate" {
					testValue = test.Min(testValue, 100)
				}

				assert.Nil(t, err)
				assert.Equal(t, testValue, gauge, fqName, name, key, stats.Op.Samples[name])
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= len(defaultConfig.Collectors.BucketStats.Metrics)+2 {
				return
			}
		}
	}
}

func TestBucketStatsCollectReturnsCorrectValuesWithMultipleBuckets(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)

	firstBucket := test.GenerateBucket("wawa-bucket")
	buckets = append(buckets, firstBucket)

	secondBucket := test.GenerateBucket("super-america")
	buckets = append(buckets, secondBucket)

	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	firstBucketStats := test.GenerateBucketStats()
	secondBucketStats := test.GenerateBucketStats()

	mockClient.EXPECT().BucketStats(firstBucket.Name).Times(1).Return(firstBucketStats, nil)
	mockClient.EXPECT().BucketStats(secondBucket.Name).Times(1).Return(secondBucketStats, nil)

	testCollector := collectors.NewBucketStatsCollector(mockClient, defaultConfig.Collectors.BucketStats)
	c := make(chan prometheus.Metric)
	count := 0

	defer close(c)

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbbucketstat_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbbucketstat_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.BucketStats, fqName)
				name := defaultConfig.Collectors.BucketStats.Metrics[key].Name
				bucketName, err := test.GetBucketIfPresent(m)

				assert.Nil(t, err)

				gauge, err := test.GetGaugeValue(m)
				testValue := test.GetBucketStatsTestValue(name, bucketName, map[string]objects.BucketStats{
					firstBucket.Name:  firstBucketStats,
					secondBucket.Name: secondBucketStats,
				})

				assert.Nil(t, err)
				assert.Equal(t, testValue, gauge, fqName, name, key, firstBucketStats.Op.Samples[name])
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= 2*len(defaultConfig.Collectors.BucketStats.Metrics)+2 {
				return
			}
		}
	}
}
