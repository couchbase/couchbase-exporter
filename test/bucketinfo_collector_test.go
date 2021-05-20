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

const (
	DummyError = "ooops, i did it again"
)

var (
	ErrorDummy = fmt.Errorf(DummyError)
)

func TestBucketInfoDescribeReturnsAppropriateValuesBasedOnDefaultConfig(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan *prometheus.Desc, 9)
	testCollector.Describe(c)
	close(c)

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketInfo.Namespace,
			defaultConfig.Collectors.BucketInfo.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketInfo.Namespace,
			defaultConfig.Collectors.BucketInfo.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.BucketInfo.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.BucketInfo.Namespace,
				defaultConfig.Collectors.BucketInfo.Subsystem,
				val.Name,
				val.HelpText,
				val.Labels))
	}

	for desc := range c {
		found := false

		for _, check := range possibleValues {
			if desc.String() == check {
				found = true
			}
		}

		assert.True(t, found, desc.String())
	}
}

func TestBucketInfoDescribeReturnsAppropriateValuesWithNameOverride(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	// set up the name overrides
	metrics := make(map[string]objects.MetricInfo)
	now := time.Now()

	for key, val := range defaultConfig.Collectors.BucketInfo.Metrics {
		val.NameOverride = fmt.Sprintf("%s-%s", val.Name, now.Format("YYYY-MM-DD"))
		metrics[key] = val
	}

	defaultConfig.Collectors.BucketInfo.Metrics = metrics
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan *prometheus.Desc, 9)
	testCollector.Describe(c)
	close(c)

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketInfo.Namespace,
			defaultConfig.Collectors.BucketInfo.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.BucketInfo.Namespace,
			defaultConfig.Collectors.BucketInfo.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.BucketInfo.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.BucketInfo.Namespace,
				defaultConfig.Collectors.BucketInfo.Subsystem,
				val.NameOverride,
				val.HelpText,
				val.Labels))
	}

	for desc := range c {
		found := false

		for _, check := range possibleValues {
			if desc.String() == check {
				found = true
				break
			}
		}

		assert.True(t, found, desc.String())
	}
}

func TestBucketInfoCollectReturnsDownIfClientReturnsError(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrorDummy)
	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestBucketInfoCollectReturnsDownIfClientReturnsErrorOnBuckets(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, ErrorDummy)

	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestBucketInfoCollectReturnsUpWithNoErrors(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
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

func TestBucketInfoCollectReturnsOneOfEachMetricForASingleBucketWithCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	singleBucket := test.GenerateBucketInfo("wawa-bucket")
	buckets = append(buckets, singleBucket)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan prometheus.Metric, 9)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		fqName := test.GetFQNameFromDesc(m.Desc())
		switch fqName {
		case "cbbucketinfo_up":
			gauge, err := test.GetGaugeValue(m)
			assert.Nil(t, err)
			assert.Equal(t, 1.0, gauge, fqName)
		case "cbbucketinfo_scrape_duration_seconds":
			gauge, err := test.GetGaugeValue(m)
			assert.Nil(t, err)
			assert.True(t, gauge > 0, fqName)
		default:
			key := test.GetKeyFromFQName(defaultConfig.Collectors.BucketInfo, fqName)
			gauge, err := test.GetGaugeValue(m)

			assert.Nil(t, err)
			assert.Equal(t, gauge, singleBucket.BucketBasicStats[key], fqName)
		}
	}
}

func TestBucketInfoCollectReturnsMultipleMetricsForMultipleBuckets(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	buckets := make([]objects.BucketInfo, 0)
	singleBucket := test.GenerateBucketInfo("wawa-bucket")
	buckets = append(buckets, singleBucket)

	secondBucket := test.GenerateBucketInfo("double-bucket")
	buckets = append(buckets, secondBucket)
	mockClient.EXPECT().Buckets().Times(1).Return(buckets, nil)

	testCollector := collectors.NewBucketInfoCollector(mockClient, defaultConfig.Collectors.BucketInfo)
	c := make(chan prometheus.Metric, 16)
	testCollector.Collect(c)
	close(c)

	count := 0

	for m := range c {
		count++

		fqName := test.GetFQNameFromDesc(m.Desc())
		bucket, err := test.GetBucketIfPresent(m)

		assert.Nil(t, err)

		switch fqName {
		case "cbbucketinfo_up":
			gauge, err := test.GetGaugeValue(m)
			assert.Nil(t, err)
			assert.Equal(t, 1.0, gauge, fqName)
		case "cbbucketinfo_scrape_duration_seconds":
			gauge, err := test.GetGaugeValue(m)
			assert.Nil(t, err)
			assert.True(t, gauge > 0, fqName)
		default:
			key := test.GetKeyFromFQName(defaultConfig.Collectors.BucketInfo, fqName)
			gauge, err := test.GetGaugeValue(m)

			assert.Nil(t, err)

			if bucket == "wawa-bucket" {
				assert.Equal(t, gauge, singleBucket.BucketBasicStats[key], fqName)
			} else {
				assert.Equal(t, gauge, secondBucket.BucketBasicStats[key], fqName)
			}
		}
	}

	assert.Equal(t, 16, count)
}
