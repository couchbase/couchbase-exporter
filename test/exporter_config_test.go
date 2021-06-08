package test

import (
	"os"
	"testing"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
)

func TestReturnsErrorOnNoFileFound(t *testing.T) {
	var config objects.ExporterConfig
	err := config.ParseConfigFile("../example/config_tests/not_found.json")

	if !os.IsNotExist(err) {
		t.Error("Somehow file that should not exist was found.")
	}
}

func TestReturnsErrorOnBadFile(t *testing.T) {
	var config objects.ExporterConfig
	err := config.ParseConfigFile("../example/config_tests/improperly_formatted.json")

	if err == nil || os.IsNotExist(err) {
		t.Error("An impropertly formatted json file should have errored.")
	}
}

func TestParsesGoodConfigFile(t *testing.T) {
	var config objects.ExporterConfig
	err := config.ParseConfigFile("../example/config.json")

	if err != nil {
		t.Error("Error during parsing of config file.")
	}

	if config.CouchbaseAddress != "localhost" {
		t.Error("Parsing of config file did not work properly")
	}
}
