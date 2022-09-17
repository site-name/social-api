package plugin

import (
	"github.com/sitename/sitename/model"
)

func (s *ServicePlugin) notifyClusterPluginEvent(event model.ClusterEvent, data model.PluginEventData) {
	if s.srv.Cluster != nil {
		s.srv.Cluster.SendClusterMessage(&model.ClusterMessage{
			Event:            event,
			SendType:         model.ClusterSendReliable,
			WaitForAllToSend: true,
			Data:             data.ToJSON(),
		})
	}
}
