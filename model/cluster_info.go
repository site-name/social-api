package model

import (
	"io"
)

type ClusterInfo struct {
	Id         string `json:"id"`
	Version    string `json:"version"`
	ConfigHash string `json:"config_hash"`
	IpAddress  string `json:"ipaddress"`
	Hostname   string `json:"hostname"`
}

func (ci *ClusterInfo) ToJson() string {
	return ModelToJson(ci)
}

func ClusterInfoFromJson(data io.Reader) *ClusterInfo {
	var ci *ClusterInfo
	ModelFromJson(&ci, data)
	return ci
}

func ClusterInfosToJson(objmap []*ClusterInfo) string {
	return ModelToJson(&objmap)
}

func ClusterInfosFromJson(data io.Reader) []*ClusterInfo {
	var objmap []*ClusterInfo
	ModelFromJson(&objmap, data)
	return objmap
}
