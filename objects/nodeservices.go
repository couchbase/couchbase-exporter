package objects

type NodeService struct {
	Services           []string            `json:"services"`
	ThisNode           bool                `json:"thisNode,omitempty"`
	Hostname           string              `json:"hostname"`
	AlternateAddresses *AlternateAddresses `json:"alternateAddresses,omitempty"`
}
