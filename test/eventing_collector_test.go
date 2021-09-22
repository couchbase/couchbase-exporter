package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/config"
	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/test/mocks"
	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestEventingDescribeReturnsAppropriateValuesBasedOnDefaultConfig(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Eventing.Namespace,
			defaultConfig.Collectors.Eventing.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Eventing.Namespace,
			defaultConfig.Collectors.Eventing.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Eventing.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Eventing.Namespace,
				defaultConfig.Collectors.Eventing.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
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

func TestEventingDescribeReturnsAppropriateValuesWithNameOverride(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	// set up the name overrides
	metrics := make(map[string]objects.MetricInfo)
	now := time.Now()

	for key, val := range defaultConfig.Collectors.Eventing.Metrics {
		val.NameOverride = fmt.Sprintf("%s-%s", val.Name, now.Format("YYYY-MM-DD"))
		metrics[key] = val
	}

	defaultConfig.Collectors.Eventing.Metrics = metrics
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Eventing.Namespace,
			defaultConfig.Collectors.Eventing.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Eventing.Namespace,
			defaultConfig.Collectors.Eventing.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Eventing.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Eventing.Namespace,
				defaultConfig.Collectors.Eventing.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
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

func TestEventingCollectReturnsDownIfClientReturnsError(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrDummy)
	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestEventingCollectReturnsDownIfClientReturnsErrorOnEventing(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	eventing := objects.Eventing{}
	mockClient.EXPECT().Eventing().Times(1).Return(eventing, ErrDummy)

	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestEventingCollectReturnsUpWithNoErrors(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	metrics := make(map[string]objects.MetricInfo)

	for key, val := range defaultConfig.Collectors.Eventing.Metrics {
		val.Enabled = false
		metrics[key] = val
	}

	defaultConfig.Collectors.Eventing.Metrics = metrics

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	eventing := objects.Eventing{}
	mockClient.EXPECT().Eventing().Times(1).Return(eventing, nil)

	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
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

func TestEventingCollectReturnsOneOfEachMetricWithCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	eventing := test.GenerateEventing()
	mockClient.EXPECT().Eventing().Times(1).Return(eventing, nil)

	testCollector := collectors.NewEventingCollector(mockClient, defaultConfig.Collectors.Eventing)
	c := make(chan prometheus.Metric, 9)
	count := 0

	defer close(c)

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbeventing_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbeventing_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.Eventing, fqName)
				name := defaultConfig.Collectors.Eventing.Metrics[key].Name

				sampleName := "eventing/"
				if strings.HasPrefix(name, "test_") {
					sampleName += strings.ReplaceAll(name, "test_", "test/")
				} else {
					sampleName += name
				}

				gauge, err := test.GetGaugeValue(m)
				testValue := test.Last(eventing.Op.Samples[sampleName])

				assert.Equal(t, testValue, gauge)
				assert.Nil(t, err)
				log.Debug("%s: %v", name, gauge)
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= len(defaultConfig.Collectors.Eventing.Metrics)+2 {
				return
			}
		}
	}
}
