package einterfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type ClusterMessageHandler func(msg *model_helper.ClusterMessage)

type ClusterInterface interface {
	StartInterNodeCommunication()
	StopInterNodeCommunication()
	RegisterClusterMessageHandler(event model_helper.ClusterEvent, crm ClusterMessageHandler)
	GetClusterId() string
	IsLeader() bool
	// HealthScore returns a number which is indicative of how well an instance is meeting
	// the soft real-time requirements of the protocol. Lower numbers are better,
	// and zero means "totally healthy".
	HealthScore() int
	GetMyClusterInfo() *model_helper.ClusterInfo
	GetClusterInfos() []*model_helper.ClusterInfo
	SendClusterMessage(model_helper *model_helper.ClusterMessage)
	NotifyMsg(buf []byte)
	GetClusterStats() ([]*model_helper.ClusterStats, *model_helper.AppError)
	GetLogs(page, perPage int) ([]string, *model_helper.AppError)
	GetPluginStatuses() (model_helper.PluginStatuses, *model_helper.AppError)
	ConfigChanged(previousConfig *model_helper.Config, newConfig *model_helper.Config, sendToOtherServer bool) *model_helper.AppError
}
