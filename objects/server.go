package objects

type Server struct {
	Hostname string            `json:"hostname"`
	URI      string            `json:"uri"`
	Stats    map[string]string `json:"stats"`
}

type Servers struct {
	Servers []Server `json:"servers"`
}