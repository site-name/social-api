package app

import (
	"time"

	"github.com/sitename/sitename/model_helper"
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

	a.srv.Account.InvalidateCacheForUser(userID)
}

// Publish push websocket event to all subscribers
func (s *Server) Publish(message *model_helper.WebSocketEvent) {
	if s.Metrics != nil {
		s.Metrics.IncrementWebsocketEvent(message.EventType())
	}

	s.PublishSkipClusterSend(message)

	if s.Cluster != nil {
		cm := &model_helper.ClusterMessage{
			Event:    model_helper.ClusterEventPublish,
			SendType: model_helper.ClusterSendBestEffort,
			Data:     message.ToJSON(),
		}

		switch message.EventType() {
		case model_helper.WebsocketEventPosted, model_helper.WebsocketEventPostEdited, model_helper.WebsocketEventDirectAdded, model_helper.WebsocketEventGroupAdded, model_helper.WebsocketEventAddedToTeam:
			cm.SendType = model_helper.ClusterSendReliable
		default:
		}

		s.Cluster.SendClusterMessage(cm)
	}
}

// Publish puplish websocket events
func (a *App) Publish(message *model_helper.WebSocketEvent) {
	a.Srv().Publish(message)
}

func (s *Server) PublishSkipClusterSend(event *model_helper.WebSocketEvent) {
	// if event.GetBroadcast().UserId != "" {

	// }
	panic("not implemented")
}
