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

func TestIndexDescribeReturnsAppropriateValuesBasedOnDefaultConfig(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Index.Namespace,
			defaultConfig.Collectors.Index.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Index.Namespace,
			defaultConfig.Collectors.Index.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Index.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Index.Namespace,
				defaultConfig.Collectors.Index.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
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

func TestIndexDescribeReturnsAppropriateValuesWithNameOverride(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	// set up the name overrides
	metrics := make(map[string]objects.MetricInfo)
	now := time.Now()

	for key, val := range defaultConfig.Collectors.Index.Metrics {
		val.NameOverride = fmt.Sprintf("%s-%s", val.Name, now.Format("YYYY-MM-DD"))
		metrics[key] = val
	}

	defaultConfig.Collectors.Index.Metrics = metrics
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	var possibleValues []string
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Index.Namespace,
			defaultConfig.Collectors.Index.Subsystem,
			"up",
			"Couchbase cluster API is responding",
			[]string{"cluster"}))
	possibleValues = append(possibleValues,
		test.GetDescString(
			defaultConfig.Collectors.Index.Namespace,
			defaultConfig.Collectors.Index.Subsystem,
			"scrape_duration_seconds",
			"Scrape duration in seconds",
			[]string{"cluster"}))

	for _, val := range defaultConfig.Collectors.Index.Metrics {
		possibleValues = append(possibleValues,
			test.GetDescString(
				defaultConfig.Collectors.Index.Namespace,
				defaultConfig.Collectors.Index.Subsystem,
				test.GetName(val.Name, val.NameOverride),
				val.HelpText,
				val.Labels))
	}

	mockClient := mocks.NewMockCbClient(mockCtrl)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
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

func TestIndexCollectReturnsDownIfClientReturnsError(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", ErrDummy)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestIndexCollectReturnsDownIfClientReturnsErrorOnIndex(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(node, nil)

	Index := objects.Index{}
	mockClient.EXPECT().Index().Times(1).Return(Index, ErrDummy)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
	c := make(chan prometheus.Metric, 1)
	testCollector.Collect(c)
	close(c)

	for m := range c {
		gauge, err := test.GetGaugeValue(m)
		assert.Nil(t, err)
		assert.Equal(t, 0.0, gauge)
	}
}

func TestIndexCollectReturnsUpWithNoErrors(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	metrics := make(map[string]objects.MetricInfo)

	for key, val := range defaultConfig.Collectors.Index.Metrics {
		val.Enabled = false
		metrics[key] = val
	}

	defaultConfig.Collectors.Index.Metrics = metrics

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(2).Return(Node, nil)

	Index := objects.Index{}
	mockClient.EXPECT().Index().Times(1).Return(Index, nil)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
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

func contains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

func TestIndexCollectReturnsOneOfEachMetricWithCorrectValues(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Index := test.GenerateIndex()
	mockClient.EXPECT().Index().Times(1).Return(Index, nil)

	Node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(2).Return(Node, nil)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
	c := make(chan prometheus.Metric, 9)
	count := 0

	defer close(c)

	metricCount := 0

	for _, val := range defaultConfig.Collectors.Index.Metrics {
		if !contains(val.Labels, objects.KeyspaceLabel) {
			metricCount++
		}
	}

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbindex_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbindex_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.Index, fqName)
				name := defaultConfig.Collectors.Index.Metrics[key].Name

				sampleName := "index_" + name

				gauge, err := test.GetGaugeValue(m)
				testValue := test.Last(Index.Op.Samples[sampleName])

				assert.Equal(t, testValue, gauge)
				assert.Nil(t, err)
				log.Debug("%s: %v", name, gauge)
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= metricCount {
				return
			}
		}
	}
}

func TestIndexCollectReturnsOneOfEachMetricWithCorrectValuesWithIndexer(t *testing.T) {
	defaultConfig := config.GetDefaultConfig()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	Index := test.GenerateIndex()
	mockClient.EXPECT().Index().Times(1).Return(Index, nil)

	Node := objects.Node{
		Services: []string{"index"},
	}
	mockClient.EXPECT().GetCurrentNode().Times(2).Return(Node, nil)

	Stats := test.GenerateIndexerStats()
	mockClient.EXPECT().IndexStats().Times(1).Return(Stats, nil)
	labelManager := util.NewLabelManager(mockClient, 600*time.Second)

	testCollector := collectors.NewIndexCollector(mockClient, defaultConfig.Collectors.Index, labelManager)
	c := make(chan prometheus.Metric, 9)
	count := 0

	defer close(c)

	metricCount := len(defaultConfig.Collectors.Index.Metrics) + 2

	go testCollector.Collect(c)

	for {
		select {
		case m := <-c:
			fqName := test.GetFQNameFromDesc(m.Desc())
			switch fqName {
			case "cbindex_up":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.Equal(t, 1.0, gauge, fqName)
			case "cbindex_scrape_duration_seconds":
				gauge, err := test.GetGaugeValue(m)
				assert.Nil(t, err)
				assert.True(t, gauge > 0, fqName)
			default:
				key := test.GetKeyFromFQName(defaultConfig.Collectors.Index, fqName)
				name := defaultConfig.Collectors.Index.Metrics[key].Name

				if !contains(defaultConfig.Collectors.Index.Metrics[key].Labels, "keyspace") {
					sampleName := "index_" + name

					gauge, err := test.GetGaugeValue(m)
					testValue := test.Last(Index.Op.Samples[sampleName])

					assert.Equal(t, testValue, gauge)
					assert.Nil(t, err)
					log.Debug("%s: %v", name, gauge)
				} else {
					gauge, err := test.GetGaugeValue(m)
					assert.Nil(t, err)

					keyspace, err := test.GetKeyspaceLabelIfPresent(m)
					assert.Nil(t, err)

					testValue, ok := Stats[keyspace][name].(float64)
					assert.True(t, ok)

					assert.Equal(t, testValue, gauge)
					log.Debug("%s: %v", name, gauge)
				}
			}
			count++
		case <-time.After(1 * time.Second):
			if count >= metricCount {
				return
			}
		}
	}
}
