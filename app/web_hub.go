package app

import "github.com/sitename/sitename/model"

func (a *App) InvalidateCacheForUser(userID string) {
	// a.Srv().
	panic("not implemented")
}

func (s *Server) Publish(message *model.WebSocketEvent) {
	if s.Metrics != nil {
		s.Metrics.IncrementWebsocketEvent(message.EventType())
	}

	s.PublishSkipClusterSend(message)

	if s.Cluster != nil {
		cm := &model.ClusterMessage{
			Event:    model.CLUSTER_EVENT_PUBLISH,
			SendType: model.CLUSTER_SEND_BEST_EFFORT,
			Data:     message.ToJson(),
		}

		switch message.EventType() {
		case model.WEBSOCKET_EVENT_POSTED, model.WEBSOCKET_EVENT_POST_EDITED, model.WEBSOCKET_EVENT_DIRECT_ADDED, model.WEBSOCKET_EVENT_GROUP_ADDED, model.WEBSOCKET_EVENT_ADDED_TO_TEAM:
			cm.SendType = model.CLUSTER_SEND_RELIABLE
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
