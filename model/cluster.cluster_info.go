package model

type ClusterInfo struct {
	Id         string `json:"id"`
	Version    string `json:"version"`
	ConfigHash string `json:"config_hash"`
	IpAddress  string `json:"ipaddress"`
	Hostname   string `json:"hostname"`
}
