package cluster

import (
	"io"
	"os"

	"github.com/sitename/sitename/model"
)

const (
	CDS_OFFLINE_AFTER_MILLIS = 1000 * 60 * 30 // 30 minutes
	CDS_TYPE_APP             = "mattermost_app"
)

type ClusterDiscovery struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	ClusterName string `json:"cluster_name"`
	Hostname    string `json:"hostname"`
	GossipPort  int32  `json:"gossip_port"`
	Port        int32  `json:"port"`
	CreateAt    int64  `json:"create_at"`
	LastPingAt  int64  `json:"last_ping_at"`
}

func (o *ClusterDiscovery) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}

	if o.CreateAt == 0 {
		o.CreateAt = model.GetMillis()
		o.LastPingAt = o.CreateAt
	}
}

func (o *ClusterDiscovery) AutoFillHostname() {
	// attempt to set the hostname from the OS
	if o.Hostname == "" {
		if hn, err := os.Hostname(); err == nil {
			o.Hostname = hn
		}
	}
}

func (o *ClusterDiscovery) AutoFillIpAddress(iface string, ipAddress string) {
	// attempt to set the hostname to the first non-local IP address
	if o.Hostname == "" {
		if ipAddress != "" {
			o.Hostname = ipAddress
		} else {
			o.Hostname = model.GetServerIpAddress(iface)
		}
	}
}

func (o *ClusterDiscovery) IsEqual(in *ClusterDiscovery) bool {
	if in == nil {
		return false
	}

	if o.Type != in.Type {
		return false
	}

	if o.ClusterName != in.ClusterName {
		return false
	}

	if o.Hostname != in.Hostname {
		return false
	}

	return true
}

func FilterClusterDiscovery(vs []*ClusterDiscovery, f func(*ClusterDiscovery) bool) []*ClusterDiscovery {
	copy := make([]*ClusterDiscovery, 0)
	for _, v := range vs {
		if f(v) {
			copy = append(copy, v)
		}
	}

	return copy
}

func (o *ClusterDiscovery) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.cluster.is_valid.%s.app_error",
		"cluster_discovery_id=",
		"ClusterDiscovery.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.ClusterName == "" {
		return outer("cluster_name", &o.Id)
	}
	if o.Type == "" {
		return outer("type", &o.Id)
	}
	if o.Hostname == "" {
		return outer("host_name", &o.Id)
	}
	if o.CreateAt == 0 {
		return outer("create_at", &o.Id)
	}
	if o.LastPingAt == 0 {
		return outer("last_ping_at", &o.Id)
	}

	return nil
}

func (o *ClusterDiscovery) ToJson() string {
	return model.ModelToJson(o)
}

func ClusterDiscoveryFromJson(data io.Reader) *ClusterDiscovery {
	var me *ClusterDiscovery
	model.ModelFromJson(&me, data)
	return me
}
