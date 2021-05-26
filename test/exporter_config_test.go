package test

import (
	"os"
	"strings"
	"testing"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
)

const fileNotExists = "file not found"

func TestReturnsErrorOnNoFileFound(t *testing.T) {
	var config objects.ExporterConfig
	err := config.ParseConfigFile("../example/config_tests/not_found.json")

	if !(strings.Contains(err.Error(), fileNotExists)) {
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
