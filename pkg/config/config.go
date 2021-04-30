package config

import (
	"os"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
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
			return exporterConfig, err
		}
	} else {
		exporterConfig.SetDefaults()
	}

	return exporterConfig, nil
}
