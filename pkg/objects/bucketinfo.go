//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package objects

// /pools/default/buckets/  to list all buckets
// /pools/default/buckets/<bucket-name>

const (
	// Bucket Basic Stats Keys.
	QuotaPercentUsed       = "quotaPercentUsed"
	OpsPerSec              = "opsPerSec"
	DiskFetches            = "diskFetches"
	ItemCount              = "itemCount"
	DiskUsed               = "diskUsed"
	DataUsed               = "dataUsed"
	MemUsed                = "memUsed"
	VbActiveNumNonResident = "vbActiveNumNonResident"
)

type BucketInfo struct {
	Name              string `json:"name"`
	BucketType        string `json:"bucketType"`
	AuthType          string `json:"authType"`
	ProxyPort         int    `json:"proxyPort"`
	URI               string `json:"uri"`
	StreamingURI      string `json:"streamingUri"`
	LocalRandomKeyURI string `json:"localRandomKeyUri"`
	Controllers       struct {
		Flush         string `json:"flush"`
		CompactAll    string `json:"compactAll"`
		CompactDB     string `json:"compactDB"`
		PurgeDeletes  string `json:"purgeDeletes"`
		StartRecovery string `json:"startRecovery"`
	} `json:"controllers"`
	Nodes []Node `json:"nodes"`
	Stats struct {
		URI              string `json:"uri"`
		DirectoryURI     string `json:"directoryURI"`
		NodeStatsListURI string `json:"nodeStatsListURI"`
	} `json:"stats"`
	NodeLocator  string `json:"nodeLocator"`
	SaslPassword string `json:"saslPassword"`
	Ddocs        struct {
		URI string `json:"uri"`
	} `json:"ddocs"`
	ReplicaIndex           bool        `json:"replicaIndex"`
	AutoCompactionSettings interface{} `json:"autoCompactionSettings"`
	UUID                   string      `json:"uuid"`
	VBucketServerMap       struct {
		HashAlgorithm string   `json:"hashAlgorithm"`
		NumReplicas   int      `json:"numReplicas"`
		ServerList    []string `json:"serverList"`
		VBucketMap    [][]int  `json:"vBucketMap"`
	} `json:"vBucketServerMap"`
	MaxTTL          int    `json:"maxTTL"`
	CompressionMode string `json:"compressionMode"`
	ReplicaNumber   int    `json:"replicaNumber"`
	ThreadsNumber   int    `json:"threadsNumber"`
	Quota           struct {
		RAM    int `json:"ram"`
		RawRAM int `json:"rawRAM"`
	} `json:"quota"`
	BucketBasicStats       map[string]float64 `json:"basicStats"`
	EvictionPolicy         string             `json:"evictionPolicy"`
	ConflictResolutionType string             `json:"conflictResolutionType"`
	BucketCapabilitiesVer  string             `json:"bucketCapabilitiesVer"`
	BucketCapabilities     []string           `json:"bucketCapabilities"`
}
