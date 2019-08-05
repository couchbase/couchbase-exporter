package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/couchbase/couchbase_exporter/util"
	"sync"
)

type MetaCollector struct {
	mutex  sync.Mutex
	client util.Client

	up             *prometheus.Desc
	scrapeDuration *prometheus.Desc
}
