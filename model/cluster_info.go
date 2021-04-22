package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

type ClusterInfo struct {
	Id         string `json:"id"`
	Version    string `json:"version"`
	ConfigHash string `json:"config_hash"`
	IpAddress  string `json:"ipaddress"`
	Hostname   string `json:"hostname"`
}

func (ci *ClusterInfo) ToJson() string {
	b, _ := json.JSON.Marshal(ci)
	return string(b)
}

func ClusterInfoFromJson(data io.Reader) *ClusterInfo {
	var ci *ClusterInfo
	json.JSON.NewDecoder(data).Decode(&ci)
	return ci
}

func ClusterInfosToJson(objmap []*ClusterInfo) string {
	b, _ := json.JSON.Marshal(objmap)
	return string(b)
}

func ClusterInfosFromJson(data io.Reader) []*ClusterInfo {
	decoder := json.JSON.NewDecoder(data)

	var objmap []*ClusterInfo
	if err := decoder.Decode(&objmap); err != nil {
		return make([]*ClusterInfo, 0)
	}
	return objmap
}
