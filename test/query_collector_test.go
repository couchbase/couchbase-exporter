// nolint: dupl
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
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/couchbase/couchbase-exporter/test/mocks"
	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestQueryDescribeReturnsAppropriateValuesBasedOnDefaultConfig(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Query.Namespace,
			defaultConfig.Collectors.Query.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Query.Namespace,
			defaultConfig.Collectors.Query.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Query.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Query.Namespace,
				defaultConfig.Collectors.Query.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
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

func TestQueryDescribeReturnsAppropriateValuesWithNameOverride(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	// set up the name overrides
	metrics := make(map[string]objects.MetricInfo)
	now := time.Now()

	for key, val := range defaultConfig.Collectors.Query.Metrics {
		val.NameOverride = fmt.Sprintf("%s-%s", val.Name, now.Format("YYYY-MM-DD"))
		metrics[key] = val
	}

	defaultConfig.Collectors.Query.Metrics = metrics
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Query.Namespace,
			defaultConfig.Collectors.Query.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Query.Namespace,
			defaultConfig.Collectors.Query.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Query.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Query.Namespace,
				defaultConfig.Collectors.Query.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
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

func TestQueryCollectReturnsDownIfClientReturnsError(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrDummy)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestQueryCollectReturnsDownIfClientReturnsErrorOnQuery(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(Node, nil)

	Query := objects.Query{}
	mockClient.EXPECT().Query().Times(1).Return(Query, ErrDummy)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestQueryCollectReturnsUpWithNoErrors(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	metrics := make(map[string]objects.MetricInfo)

	for key, val := range defaultConfig.Collectors.Query.Metrics {
		val.Enabled = false
		metrics[key] = val
	}

	defaultConfig.Collectors.Query.Metrics = metrics

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(Node, nil)

	Query := objects.Query{}
	mockClient.EXPECT().Query().Times(1).Return(Query, nil)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
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

func TestQueryCollectReturnsOneOfEachMetricWithCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(Node, nil)

	Query := test.GenerateQuery()
	mockClient.EXPECT().Query().Times(1).Return(Query, nil)
	labelManager := util.NewLabelManager(mockClient)

	testCollector := collectors.NewQueryCollector(mockClient, defaultConfig.Collectors.Query, labelManager)
	c := make(chan prometheus.Metric, 9)
	count := 0

	defer close(c)

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbquery_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbquery_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.Query, fqName)
				name := defaultConfig.Collectors.Query.Metrics[key].Name

				sampleName := "query_" + name

				gauge, err := test.GetGaugeValue(m)
				testValue := test.Last(Query.Op.Samples[sampleName])

				assert.Equal(t, testValue, gauge)
				assert.Nil(t, err)
				log.Debug("%s: %v", name, gauge)
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= len(defaultConfig.Collectors.Query.Metrics)+2 {
				return
			}
		}
	}
}
