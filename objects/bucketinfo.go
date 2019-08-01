package objects

// /pools/default/buckets/  to list all buckets
// /pools/default/buckets/<bucket-name>

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
	ReplicaIndex           bool   `json:"replicaIndex"`
	AutoCompactionSettings bool   `json:"autoCompactionSettings"`
	UUID                   string `json:"uuid"`
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
	BucketBasicStats       BucketBasicStats `json:"basicStats"`
	EvictionPolicy         string           `json:"evictionPolicy"`
	ConflictResolutionType string           `json:"conflictResolutionType"`
	BucketCapabilitiesVer  string           `json:"bucketCapabilitiesVer"`
	BucketCapabilities     []string         `json:"bucketCapabilities"`
}

type BucketBasicStats struct {
	QuotaPercentUsed       float64 `json:"quotaPercentUsed"`
	OpsPerSec              float64 `json:"opsPerSec"`
	DiskFetches            float64 `json:"diskFetches"`
	ItemCount              float64 `json:"itemCount"`
	DiskUsed               float64 `json:"diskUsed"`
	DataUsed               float64 `json:"dataUsed"`
	MemUsed                float64 `json:"memUsed"`
	VbActiveNumNonResident float64 `json:"vbActiveNumNonResident"`
}
