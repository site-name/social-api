package app

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
)

const (
	broadcastQueueSize         = 4096
	inactiveConnReaperInterval = 5 * time.Minute
)

type webConnActivityMessage struct {
	userID       string
	sessionToken string
	activityAt   int64
}

// InvalidateCacheForUser
func (a *App) InvalidateCacheForUser(userID string) {
	a.Srv().invalidateCacheForUserSkipClusterSend(userID)

	a.Srv().Store.User().InvalidateProfileCacheForUser(userID)

	if a.Cluster() != nil {
		msg := &cluster.ClusterMessage{
			Event:    cluster.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_USER,
			SendType: cluster.CLUSTER_SEND_BEST_EFFORT,
			Data:     userID,
		}
		a.Cluster().SendClusterMessage(msg)
	}
}

// Publish push websocket event to all subscribers
func (s *Server) Publish(message *model.WebSocketEvent) {
	if s.Metrics != nil {
		s.Metrics.IncrementWebsocketEvent(message.EventType())
	}

	s.PublishSkipClusterSend(message)

	if s.Cluster != nil {
		cm := &cluster.ClusterMessage{
			Event:    cluster.CLUSTER_EVENT_PUBLISH,
			SendType: cluster.CLUSTER_SEND_BEST_EFFORT,
			Data:     message.ToJson(),
		}

		switch message.EventType() {
		case model.WEBSOCKET_EVENT_POSTED, model.WEBSOCKET_EVENT_POST_EDITED, model.WEBSOCKET_EVENT_DIRECT_ADDED, model.WEBSOCKET_EVENT_GROUP_ADDED, model.WEBSOCKET_EVENT_ADDED_TO_TEAM:
			cm.SendType = cluster.CLUSTER_SEND_RELIABLE
		default:
		}

		s.Cluster.SendClusterMessage(cm)
	}
}

func (a *App) Publish(message *model.WebSocketEvent) {
	a.Srv().Publish(message)
}

func (s *Server) PublishSkipClusterSend(event *model.WebSocketEvent) {
	// if event.GetBroadcast().UserId != "" {

	// }
	panic("not implemented")
}
