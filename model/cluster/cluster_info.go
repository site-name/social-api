package cluster

import (
	"io"

	"github.com/sitename/sitename/model"
)

type ClusterInfo struct {
	Id         string `json:"id"`
	Version    string `json:"version"`
	ConfigHash string `json:"config_hash"`
	IpAddress  string `json:"ipaddress"`
	Hostname   string `json:"hostname"`
}

func (ci *ClusterInfo) ToJson() string {
	return model.ModelToJson(ci)
}

func ClusterInfoFromJson(data io.Reader) *ClusterInfo {
	var ci *ClusterInfo
	model.ModelFromJson(&ci, data)
	return ci
}

func ClusterInfosToJson(objmap []*ClusterInfo) string {
	return model.ModelToJson(&objmap)
}

func ClusterInfosFromJson(data io.Reader) []*ClusterInfo {
	var objmap []*ClusterInfo
	model.ModelFromJson(&objmap, data)
	return objmap
}
