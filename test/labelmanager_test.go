package test

import (
	"testing"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/couchbase/couchbase-exporter/test/mocks"
	test "github.com/couchbase/couchbase-exporter/test/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLabelManagerCallsClusterNameOnce(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(node, nil)

	manager := util.NewLabelManager(mockClient)

	ctx, err := manager.GetMetricContext("a", "b")

	assert.Nil(t, err)
	assert.Equal(t, "dummy-cluster", ctx.ClusterName)
}

func TestLabelManagerCallsClusterNameOnceEvenOnSubsequentRequests(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(node, nil)

	manager := util.NewLabelManager(mockClient)

	ctx, err := manager.GetMetricContext("a", "b")
	assert.Nil(t, err)

	ctx2, err := manager.GetMetricContext("x", "d")

	assert.Nil(t, err)
	assert.Equal(t, "dummy-cluster", ctx.ClusterName)
	assert.Equal(t, "a", ctx.BucketName)
	assert.Equal(t, "x", ctx2.BucketName)
}

func TestLabelManagerReturnsErrorIfClientErrors(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("", ErrDummy)

	manager := util.NewLabelManager(mockClient)

	_, err := manager.GetMetricContext("a", "b")
	assert.NotNil(t, err)
}

func TestLabelManagerGetsAppropriateValuesFromCTX(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	manager := util.NewLabelManager(mockClient)

	ctx := util.MetricContext{
		ClusterName:  "dummy-cluster",
		BucketName:   "travel-sample",
		NodeHostname: "localhost",
		Keyspace:     "travel-sample:testIndex",
	}

	labelValues := manager.GetLabelValues([]string{objects.ClusterLabel, objects.BucketLabel, objects.KeyspaceLabel, objects.NodeLabel}, ctx)

	assert.Contains(t, labelValues, "dummy-cluster")
	assert.Contains(t, labelValues, "travel-sample")
	assert.Contains(t, labelValues, "localhost")
	assert.Contains(t, labelValues, "travel-sample:testIndex")
}

func TestLabelManagerSplitsValuesWithColonAndUsesSecondForValue(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	manager := util.NewLabelManager(mockClient)

	ctx := util.MetricContext{
		ClusterName:  "dummy-cluster",
		BucketName:   "travel-sample",
		NodeHostname: "localhost",
		Keyspace:     "travel-sample:testIndex",
	}

	labelValues := manager.GetLabelValues([]string{objects.ClusterLabel, objects.BucketLabel, objects.KeyspaceLabel, objects.NodeLabel, "new:value"}, ctx)

	assert.Contains(t, labelValues, "dummy-cluster")
	assert.Contains(t, labelValues, "travel-sample")
	assert.Contains(t, labelValues, "localhost")
	assert.Contains(t, labelValues, "travel-sample:testIndex")
	assert.Contains(t, labelValues, "value")
}

func TestLabelManagerTakesLabelForValueIfUnrecognized(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	manager := util.NewLabelManager(mockClient)

	ctx := util.MetricContext{
		ClusterName:  "dummy-cluster",
		BucketName:   "travel-sample",
		NodeHostname: "localhost",
		Keyspace:     "travel-sample:testIndex",
	}

	labelValues := manager.GetLabelValues([]string{objects.ClusterLabel, objects.BucketLabel, objects.KeyspaceLabel, objects.NodeLabel, "new:value", "foobarbaz"}, ctx)

	assert.Contains(t, labelValues, "dummy-cluster")
	assert.Contains(t, labelValues, "travel-sample")
	assert.Contains(t, labelValues, "localhost")
	assert.Contains(t, labelValues, "travel-sample:testIndex")
	assert.Contains(t, labelValues, "value")
	assert.Contains(t, labelValues, "foobarbaz")
}

func TestLabelManagerSplitsValueOnColonAndRetursFirstValueAsLabel(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	manager := util.NewLabelManager(mockClient)

	labelValues := manager.GetLabelKeys([]string{objects.ClusterLabel, objects.BucketLabel, objects.KeyspaceLabel, objects.NodeLabel, "new:value", "foobarbaz"})

	assert.Contains(t, labelValues, "cluster")
	assert.Contains(t, labelValues, "bucket")
	assert.Contains(t, labelValues, "node")
	assert.Contains(t, labelValues, "keyspace")
	assert.Contains(t, labelValues, "new")
	assert.Contains(t, labelValues, "foobarbaz")
}
func TestLabelManagerGetMetricContextCacheRace(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockClient := mocks.NewMockCbClient(mockCtrl)
	mockClient.EXPECT().ClusterName().Times(1).Return("dummy-cluster", nil)

	node := test.GenerateNode()
	mockClient.EXPECT().GetCurrentNode().Times(1).Return(node, nil)

	manager := util.NewLabelManager(mockClient)

	for i := 0; i < 1000; i++ {
		go func() {
			ctx, err := manager.GetMetricContext("a", "b")

			assert.Nil(t, err)
			assert.Equal(t, "dummy-cluster", ctx.ClusterName)
		}()
	}
}
