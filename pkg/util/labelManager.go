package util

import (
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
)

type LabelKeyPair struct {
	Key   string
	Value string
}

type MetricContext struct {
	ClusterName  string
	NodeHostname string
	BucketName   string
	Keyspace     string
	Source       string
	Target       string
}

// CbLabelManager is an interface to be used to obtain a context for labels and to retrieve label values from a context and/or the metric configuration
// Attempts to centralize this code so it is testable and logic is segregated.
type CbLabelManager interface {
	GetBasicMetricContext() (MetricContext, error)
	GetMetricContext(bucket, keyspace string) (MetricContext, error)
	GetMetricContextWithSourceAndTarget(bucket, keyspace, source, target string) (MetricContext, error)
	GetLabelValues(labels []string, context MetricContext) []string
	GetLabelKeys(labels []string) []string
}

type labelManager struct {
	client     CbClient
	labelCache cache
}

func NewLabelManager(client CbClient) CbLabelManager {
	lbl := &labelManager{
		client:     client,
		labelCache: newCache(),
	}

	return lbl
}

func (l *labelManager) GetMetricContext(bucket, keyspace string) (MetricContext, error) {
	ctx := MetricContext{
		BucketName: bucket,
		Keyspace:   keyspace,
	}

	if !l.labelCache.isExpired(objects.ClusterLabel) {
		val, _ := l.labelCache.get(objects.ClusterLabel).(string)
		ctx.ClusterName = val
	} else {
		clusterName, err := l.client.ClusterName()
		if err != nil {
			return ctx, err
		}
		ctx.ClusterName = clusterName
		l.labelCache.set(objects.ClusterLabel, clusterName)
	}

	if !l.labelCache.isExpired(objects.NodeLabel) {
		val, _ := l.labelCache.get(objects.NodeLabel).(string)
		ctx.NodeHostname = val
	} else {
		node, err := l.client.GetCurrentNode()
		if err != nil {
			return ctx, err
		}
		ctx.NodeHostname = node.Hostname
		l.labelCache.set(objects.NodeLabel, node.Hostname)
	}

	return ctx, nil
}

func (l *labelManager) GetMetricContextWithSourceAndTarget(bucket, keyspace, source, target string) (MetricContext, error) {
	ctx, err := l.GetMetricContext(bucket, keyspace)
	ctx.Source = source
	ctx.Target = target

	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (l *labelManager) GetBasicMetricContext() (MetricContext, error) {
	return l.GetMetricContextWithSourceAndTarget("", "", "", "")
}

func (l *labelManager) GetLabelValues(labels []string, context MetricContext) []string {
	values := []string{}

	for _, label := range labels {
		switch label {
		case objects.ClusterLabel:
			values = append(values, context.ClusterName)
		case objects.BucketLabel:
			values = append(values, context.BucketName)
		case objects.KeyspaceLabel:
			values = append(values, context.Keyspace)
		case objects.NodeLabel:
			values = append(values, context.NodeHostname)
		case objects.TargetLabel:
			values = append(values, context.Target)
		case objects.SourceLabel:
			values = append(values, context.Source)
		default:
			if strings.Contains(label, ":") {
				splits := strings.Split(label, ":")
				values = append(values, splits[1])
			} else {
				values = append(values, label)
			}
		}
	}

	return values
}

func (l *labelManager) GetLabelKeys(labels []string) []string {
	return objects.GetLabelKeys(labels)
}

type cache struct {
	items map[string]entry
}

func newCache() cache {
	return cache{
		items: make(map[string]entry),
	}
}

const cacheDuration = 10

func (c *cache) set(key string, data interface{}) {
	c.items[key] = entry{
		data:       data,
		expiration: time.Now().Add(cacheDuration * time.Minute),
	}
}

func (c *cache) get(key string) interface{} {
	if e, ok := c.items[key]; ok {
		return e.data
	}

	return nil
}

func (c *cache) isExpired(key string) bool {
	now := time.Now()
	if e, ok := c.items[key]; ok {
		return e.expiration.Before(now)
	}

	return true
}

type entry struct {
	data       interface{}
	expiration time.Time
}
