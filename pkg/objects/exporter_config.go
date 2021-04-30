//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	operatorUser = "COUCHBASE_OPERATOR_USER"
	operatorPass = "COUCHBASE_OPERATOR_PASS"

	envUser = "COUCHBASE_USER"
	envPass = "COUCHBASE_PASS"

	bearerToken = "AUTH_BEARER_TOKEN"
)

type ExporterConfig struct {
	CouchbaseAddress  string             `json:"couchbaseAddress"`
	CouchbasePort     int                `json:"couchbasePort"`
	CouchbaseUser     string             `json:"couchbaseUser"`
	CouchbasePassword string             `json:"couchbasePassword"`
	ServerAddress     string             `json:"serverAddress"`
	ServerPort        int                `json:"serverPort"`
	RefreshRate       int                `json:"refreshRate"`
	BackoffLimit      int                `json:"backoffLimit"`
	LogLevel          string             `json:"logLevel"`
	LogJSON           bool               `json:"logJson"`
	Token             string             `json:"token"`
	Certificate       string             `json:"certificate"`
	Key               string             `json:"key"`
	Ca                string             `json:"ca"`
	ClientCertificate string             `json:"clientCertificate"`
	ClientKey         string             `json:"clientKey"`
	Collectors        ExporterCollectors `json:"collectors"`
}

type ExporterCollectors struct {
	BucketInfo         *CollectorConfig `json:"bucketInfo"`
	BucketStats        *CollectorConfig `json:"bucketStats"`
	Analytics          *CollectorConfig `json:"analytics"`
	Eventing           *CollectorConfig `json:"eventing"`
	Index              *CollectorConfig `json:"index"`
	Node               *CollectorConfig `json:"node"`
	Query              *CollectorConfig `json:"query"`
	Search             *CollectorConfig `json:"search"`
	Task               *CollectorConfig `json:"task"`
	PerNodeBucketStats *CollectorConfig `json:"perNodeBucketStats"`
}

func (e *ExporterConfig) ParseConfigFile(configFilePath string) error {
	if _, err := os.Stat(configFilePath); err != nil {
		return err
	}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		defer jsonFile.Close()
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		defer jsonFile.Close()
		return err
	}

	err = json.Unmarshal(byteValue, &e)
	if err != nil {
		defer jsonFile.Close()
		return err
	}

	defer jsonFile.Close()

	return nil
}

func (e *ExporterConfig) SetDefaults() {
	e.BackoffLimit = 5
	e.Ca = ""
	e.Certificate = ""
	e.ClientCertificate = ""
	e.ClientKey = ""
	e.Collectors = ExporterCollectors{
		BucketInfo:         GetBucketInfoCollectorDefaultConfig(),
		BucketStats:        GetBucketStatsCollectorDefaultConfig(),
		Analytics:          GetAnalyticsCollectorDefaultConfig(),
		Eventing:           GetEventingCollectorDefaultConfig(),
		Index:              GetIndexCollectorDefaultConfig(),
		Node:               GetNodeCollectorDefaultConfig(),
		Query:              GetQueryCollectorDefaultConfig(),
		Search:             GetSearchCollectorDefaultConfig(),
		Task:               GetTaskCollectorDefaultConfig(),
		PerNodeBucketStats: GetPerNodeBucketStatsCollectorDefaultConfig(),
	}
	e.CouchbaseAddress = "localhost"
	e.CouchbasePort = 8091
	e.CouchbaseUser = "Administrator"
	e.CouchbasePassword = "password"
	e.Key = ""
	e.LogJSON = true
	e.LogLevel = "info"
	e.RefreshRate = 5
	e.ServerAddress = "localhost"
	e.ServerPort = 9091
	e.Token = ""
}

func (e *ExporterConfig) SetOrDefaultLogJSON(logJSON bool) {
	if !logJSON {
		e.LogJSON = logJSON
	}
}

func (e *ExporterConfig) SetOrDefaultLogLevel(logLevel string) {
	if logLevel != "" {
		e.LogLevel = logLevel
	}
}

func (e *ExporterConfig) SetOrDefaultCouchAddress(couchAddr string) {
	if couchAddr != "" {
		e.CouchbaseAddress = couchAddr
	}
}

func (e *ExporterConfig) SetOrDefaultCouchPort(couchPort string) {
	if couchPort != "" && isInt(couchPort) {
		e.CouchbasePort, _ = strconv.Atoi(couchPort)
	}
}

func (e *ExporterConfig) SetOrDefaultCouchUser(couchUser string) {
	if couchUser != "" {
		e.CouchbaseUser = couchUser
	}

	// override defaults and CLI parameters with ENV Vars.
	// get couchbase server credentials.
	if os.Getenv(envUser) != "" {
		e.CouchbaseUser = os.Getenv(envUser)
	}

	// for operator only, override both plain-text CLI flags and other env-vars.
	// get couchbase server credentials
	if os.Getenv(operatorUser) != "" {
		e.CouchbaseUser = os.Getenv(operatorUser)
	}
}

func (e *ExporterConfig) SetOrDefaultCouchPassword(couchPass string) {
	if couchPass != "" {
		e.CouchbasePassword = couchPass
	}

	// override defaults and CLI parameters with ENV Vars.
	// get couchbase server credentials.
	if os.Getenv(envPass) != "" {
		e.CouchbasePassword = os.Getenv(envPass)
	}

	// for operator only, override both plain-text CLI flags and other env-vars.
	// get couchbase server credentials
	if os.Getenv(operatorPass) != "" {
		e.CouchbasePassword = os.Getenv(operatorPass)
	}
}

func (e *ExporterConfig) SetOrDefaultServerAddress(svrAddr string) {
	if svrAddr != "" {
		e.ServerAddress = svrAddr
	}
}

func (e *ExporterConfig) SetOrDefaultServerPort(svrPort string) {
	if svrPort != "" && isInt(svrPort) {
		e.ServerPort, _ = strconv.Atoi(svrPort)
	}
}

func (e *ExporterConfig) SetOrDefaultRefreshRate(refreshRate string) {
	if refreshRate != "" && isInt(refreshRate) {
		e.RefreshRate, _ = strconv.Atoi(refreshRate)
	}
}

func (e *ExporterConfig) SetOrDefaultBackoffLimit(backoffLimit string) {
	if backoffLimit != "" && isInt(backoffLimit) {
		e.BackoffLimit, _ = strconv.Atoi(backoffLimit)
	}
}

func (e *ExporterConfig) SetOrDefaultToken(token string) {
	if token != "" {
		e.Token = token
	}

	// override passed value with ENV var value if it has one.
	if os.Getenv(bearerToken) != "" {
		e.Token = os.Getenv(bearerToken)
	}
}

func (e *ExporterConfig) SetOrDefaultCa(ca string) {
	if ca != "" {
		e.Ca = ca
	}
}

func (e *ExporterConfig) SetOrDefaultCertificate(certificate string) {
	if certificate != "" {
		e.Certificate = certificate
	}
}

func (e *ExporterConfig) SetOrDefaultKey(key string) {
	if key != "" {
		e.Key = key
	}
}

func (e *ExporterConfig) SetOrDefaultClientCertificate(clientCertificate string) {
	if clientCertificate != "" {
		e.ClientCertificate = clientCertificate
	}
}

func (e *ExporterConfig) SetOrDefaultClientKey(clientKey string) {
	if clientKey != "" {
		e.ClientKey = clientKey
	}
}

func (e *ExporterConfig) ValidateConfig() {

}

func isInt(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}

	return false
}
