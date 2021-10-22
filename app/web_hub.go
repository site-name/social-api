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

	a.srv.AccountService().InvalidateCacheForUser(userID)
}

// Publish push websocket event to all subscribers
func (s *Server) Publish(message *model.WebSocketEvent) {
	if s.Metrics != nil {
		s.Metrics.IncrementWebsocketEvent(message.EventType())
	}

	s.PublishSkipClusterSend(message)

	if s.Cluster != nil {
		cm := &cluster.ClusterMessage{
			Event:    cluster.ClusterEventPublish,
			SendType: cluster.ClusterSendBestEffort,
			Data:     message.ToJSON(),
		}

		switch message.EventType() {
		case model.WebsocketEventPosted, model.WebsocketEventPostEdited, model.WebsocketEventDirectAdded, model.WebsocketEventGroupAdded, model.WebsocketEventAddedToTeam:
			cm.SendType = cluster.ClusterSendReliable
		default:
		}

		s.Cluster.SendClusterMessage(cm)
	}
}

// Publish puplish websocket events
func (a *App) Publish(message *model.WebSocketEvent) {
	a.Srv().Publish(message)
}

func (s *Server) PublishSkipClusterSend(event *model.WebSocketEvent) {
	// if event.GetBroadcast().UserId != "" {

	// }
	panic("not implemented")
}
