package config

import (
	"fmt"
	"os"

	"github.com/couchbase/couchbase-exporter/pkg/collectors"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	envConfigFile = "COUCHBASE_CONFIG_FILE"
)

func New(configFile string) (*objects.ExporterConfig, error) {
	exporterConfig := new(objects.ExporterConfig)
	config := ""

	if configFile != "" {
		config = configFile
	}

	if os.Getenv(envConfigFile) != "" {
		config = os.Getenv(envConfigFile)
	}

	if config != "" {
		err := exporterConfig.ParseConfigFile(config)
		if err != nil {
			return exporterConfig, fmt.Errorf("error parsing config file: %w", err)
		}
	} else {
		exporterConfig.SetDefaults()
	}

	return exporterConfig, nil
}

func GetDefaultConfig() *objects.ExporterConfig {
	exporterConfig := new(objects.ExporterConfig)
	exporterConfig.SetDefaults()

	return exporterConfig
}

func RegisterCollectors(config *objects.ExporterConfig, client util.CbClient) {
	prometheus.MustRegister(collectors.NewNodesCollector(client, config.Collectors.Node))
	prometheus.MustRegister(collectors.NewBucketInfoCollector(client, config.Collectors.BucketInfo))
	prometheus.MustRegister(collectors.NewBucketStatsCollector(client, config.Collectors.BucketStats))
	prometheus.MustRegister(collectors.NewTaskCollector(client, config.Collectors.Task))

	prometheus.MustRegister(collectors.NewQueryCollector(client, config.Collectors.Query))
	prometheus.MustRegister(collectors.NewIndexCollector(client, config.Collectors.Index))
	prometheus.MustRegister(collectors.NewFTSCollector(client, config.Collectors.Search))
	prometheus.MustRegister(collectors.NewCbasCollector(client, config.Collectors.Analytics))
	prometheus.MustRegister(collectors.NewEventingCollector(client, config.Collectors.Eventing))
}
