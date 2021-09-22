// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/util/networking.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	objects "github.com/couchbase/couchbase-exporter/pkg/objects"
	gomock "github.com/golang/mock/gomock"
)

// MockCbClient is a mock of CbClient interface.
type MockCbClient struct {
	ctrl     *gomock.Controller
	recorder *MockCbClientMockRecorder
}

// MockCbClientMockRecorder is the mock recorder for MockCbClient.
type MockCbClientMockRecorder struct {
	mock *MockCbClient
}

// NewMockCbClient creates a new mock instance.
func NewMockCbClient(ctrl *gomock.Controller) *MockCbClient {
	mock := &MockCbClient{ctrl: ctrl}
	mock.recorder = &MockCbClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCbClient) EXPECT() *MockCbClientMockRecorder {
	return m.recorder
}

// BucketNodes mocks base method.
func (m *MockCbClient) BucketNodes(arg0 string) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BucketNodes", arg0)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BucketNodes indicates an expected call of BucketNodes.
func (mr *MockCbClientMockRecorder) BucketNodes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BucketNodes", reflect.TypeOf((*MockCbClient)(nil).BucketNodes), arg0)
}

// BucketPerNodeStats mocks base method.
func (m *MockCbClient) BucketPerNodeStats(arg0, arg1 string) (objects.BucketStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BucketPerNodeStats", arg0, arg1)
	ret0, _ := ret[0].(objects.BucketStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BucketPerNodeStats indicates an expected call of BucketPerNodeStats.
func (mr *MockCbClientMockRecorder) BucketPerNodeStats(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BucketPerNodeStats", reflect.TypeOf((*MockCbClient)(nil).BucketPerNodeStats), arg0, arg1)
}

// BucketStats mocks base method.
func (m *MockCbClient) BucketStats(arg0 string) (objects.BucketStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BucketStats", arg0)
	ret0, _ := ret[0].(objects.BucketStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BucketStats indicates an expected call of BucketStats.
func (mr *MockCbClientMockRecorder) BucketStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BucketStats", reflect.TypeOf((*MockCbClient)(nil).BucketStats), arg0)
}

// Buckets mocks base method.
func (m *MockCbClient) Buckets() ([]objects.BucketInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buckets")
	ret0, _ := ret[0].([]objects.BucketInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Buckets indicates an expected call of Buckets.
func (mr *MockCbClientMockRecorder) Buckets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buckets", reflect.TypeOf((*MockCbClient)(nil).Buckets))
}

// Cbas mocks base method.
func (m *MockCbClient) Cbas() (objects.Analytics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cbas")
	ret0, _ := ret[0].(objects.Analytics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cbas indicates an expected call of Cbas.
func (mr *MockCbClientMockRecorder) Cbas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cbas", reflect.TypeOf((*MockCbClient)(nil).Cbas))
}

// ClusterName mocks base method.
func (m *MockCbClient) ClusterName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClusterName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClusterName indicates an expected call of ClusterName.
func (mr *MockCbClientMockRecorder) ClusterName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClusterName", reflect.TypeOf((*MockCbClient)(nil).ClusterName))
}

// Eventing mocks base method.
func (m *MockCbClient) Eventing() (objects.Eventing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Eventing")
	ret0, _ := ret[0].(objects.Eventing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Eventing indicates an expected call of Eventing.
func (mr *MockCbClientMockRecorder) Eventing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eventing", reflect.TypeOf((*MockCbClient)(nil).Eventing))
}

// Fts mocks base method.
func (m *MockCbClient) Fts() (objects.FTS, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fts")
	ret0, _ := ret[0].(objects.FTS)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fts indicates an expected call of Fts.
func (mr *MockCbClientMockRecorder) Fts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fts", reflect.TypeOf((*MockCbClient)(nil).Fts))
}

// Get mocks base method.
func (m *MockCbClient) Get(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockCbClientMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCbClient)(nil).Get), arg0, arg1)
}

// GetCurrentNode mocks base method.
func (m *MockCbClient) GetCurrentNode() (objects.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentNode")
	ret0, _ := ret[0].(objects.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentNode indicates an expected call of GetCurrentNode.
func (mr *MockCbClientMockRecorder) GetCurrentNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentNode", reflect.TypeOf((*MockCbClient)(nil).GetCurrentNode))
}

// Index mocks base method.
func (m *MockCbClient) Index() (objects.Index, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index")
	ret0, _ := ret[0].(objects.Index)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Index indicates an expected call of Index.
func (mr *MockCbClientMockRecorder) Index() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockCbClient)(nil).Index))
}

// IndexNode mocks base method.
func (m *MockCbClient) IndexNode(arg0 string) (objects.Index, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexNode", arg0)
	ret0, _ := ret[0].(objects.Index)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexNode indicates an expected call of IndexNode.
func (mr *MockCbClientMockRecorder) IndexNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexNode", reflect.TypeOf((*MockCbClient)(nil).IndexNode), arg0)
}

// IndexStats mocks base method.
func (m *MockCbClient) IndexStats() (map[string]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexStats")
	ret0, _ := ret[0].(map[string]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexStats indicates an expected call of IndexStats.
func (mr *MockCbClientMockRecorder) IndexStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexStats", reflect.TypeOf((*MockCbClient)(nil).IndexStats))
}

// Nodes mocks base method.
func (m *MockCbClient) Nodes() (objects.Nodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nodes")
	ret0, _ := ret[0].(objects.Nodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Nodes indicates an expected call of Nodes.
func (mr *MockCbClientMockRecorder) Nodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nodes", reflect.TypeOf((*MockCbClient)(nil).Nodes))
}

// NodesNodes mocks base method.
func (m *MockCbClient) NodesNodes() (objects.Nodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodesNodes")
	ret0, _ := ret[0].(objects.Nodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NodesNodes indicates an expected call of NodesNodes.
func (mr *MockCbClientMockRecorder) NodesNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodesNodes", reflect.TypeOf((*MockCbClient)(nil).NodesNodes))
}

// Query mocks base method.
func (m *MockCbClient) Query() (objects.Query, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query")
	ret0, _ := ret[0].(objects.Query)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockCbClientMockRecorder) Query() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockCbClient)(nil).Query))
}

// QueryNode mocks base method.
func (m *MockCbClient) QueryNode(arg0 string) (objects.Query, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryNode", arg0)
	ret0, _ := ret[0].(objects.Query)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryNode indicates an expected call of QueryNode.
func (mr *MockCbClientMockRecorder) QueryNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryNode", reflect.TypeOf((*MockCbClient)(nil).QueryNode), arg0)
}

// Servers mocks base method.
func (m *MockCbClient) Servers(arg0 string) (objects.Servers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Servers", arg0)
	ret0, _ := ret[0].(objects.Servers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Servers indicates an expected call of Servers.
func (mr *MockCbClientMockRecorder) Servers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Servers", reflect.TypeOf((*MockCbClient)(nil).Servers), arg0)
}

// Tasks mocks base method.
func (m *MockCbClient) Tasks() ([]objects.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tasks")
	ret0, _ := ret[0].([]objects.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tasks indicates an expected call of Tasks.
func (mr *MockCbClientMockRecorder) Tasks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tasks", reflect.TypeOf((*MockCbClient)(nil).Tasks))
}

// URL mocks base method.
func (m *MockCbClient) URL(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockCbClientMockRecorder) URL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockCbClient)(nil).URL), arg0)
}
