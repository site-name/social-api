package model

import (
	"net/http"
	"os"

	"gorm.io/gorm"
)

const (
	CDS_OFFLINE_AFTER_MILLIS = 1000 * 60 * 30 // 30 minutes
	CDS_TYPE_APP             = "mattermost_app"
)

type ClusterDiscovery struct {
	Id          string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Type        string `json:"type" gorm:"type:varchar(64);column:Type"`
	ClusterName string `json:"cluster_name" gorm:"type:varchar(64);column:ClusterName"`
	Hostname    string `json:"hostname" gorm:"type:varchar(512);column:Hostname"`
	GossipPort  int32  `json:"gossip_port" gorm:"type:integer;column:GossipPort"`
	Port        int32  `json:"port" gorm:"column:Port"`
	CreateAt    int64  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	LastPingAt  int64  `json:"last_ping_at" gorm:"type:bigint;column:LastPingAt;autoCreateTime:milli"`
}

func (c *ClusterDiscovery) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *ClusterDiscovery) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *ClusterDiscovery) TableName() string             { return ClusterDiscoveryTableName }

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
			o.Hostname = GetServerIpAddress(iface)
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

func (o *ClusterDiscovery) IsValid() *AppError {
	if o.ClusterName == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.cluster_name.app_error", nil, "please provide cluster name", http.StatusBadRequest)
	}
	if o.Type == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.type.app_error", nil, "please provide cluster type", http.StatusBadRequest)
	}
	if o.Hostname == "" {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.host_name.app_error", nil, "please provide host name", http.StatusBadRequest)
	}

	return nil
}
