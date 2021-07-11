package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
)

type ClusterMessageHandler func(msg *cluster.ClusterMessage)

type ClusterInterface interface {
	StartInterNodeCommunication()
	StopInterNodeCommunication()
	RegisterClusterMessageHandler(event string, crm ClusterMessageHandler)
	GetClusterId() string
	IsLeader() bool
	// HealthScore returns a number which is indicative of how well an instance is meeting
	// the soft real-time requirements of the protocol. Lower numbers are better,
	// and zero means "totally healthy".
	HealthScore() int
	GetMyClusterInfo() *cluster.ClusterInfo
	GetClusterInfos() []*cluster.ClusterInfo
	SendClusterMessage(cluster *cluster.ClusterMessage)
	NotifyMsg(buf []byte)
	GetClusterStats() ([]*cluster.ClusterStats, *model.AppError)
	GetLogs(page, perPage int) ([]string, *model.AppError)
	GetPluginStatuses() (plugins.PluginStatuses, *model.AppError)
	ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError
}
