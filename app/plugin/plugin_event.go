package plugin

import (
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
)

func (s *AppPlugin) notifyClusterPluginEvent(event string, data plugins.PluginEventData) {
	if s.Cluster() != nil {
		s.Cluster().SendClusterMessage(&cluster.ClusterMessage{
			Event:            event,
			SendType:         cluster.CLUSTER_SEND_RELIABLE,
			WaitForAllToSend: true,
			Data:             data.ToJson(),
		})
	}
}
