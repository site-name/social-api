package plugin

import (
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
)

func (s *ServicePlugin) notifyClusterPluginEvent(event cluster.ClusterEvent, data plugins.PluginEventData) {
	if s.srv.Cluster != nil {
		s.srv.Cluster.SendClusterMessage(&cluster.ClusterMessage{
			Event:            event,
			SendType:         cluster.ClusterSendReliable,
			WaitForAllToSend: true,
			Data:             data.ToJson(),
		})
	}
}
