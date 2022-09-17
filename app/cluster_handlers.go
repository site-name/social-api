package app

import (
	"bytes"
	"encoding/json"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) clusterInstallPluginHandler(msg *model.ClusterMessage) {
	var data model.PluginEventData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		slog.Warn("Failed to decode from JSON", slog.Err(err))
	}
	s.PluginService().InstallPluginFromData(data)
}

func (s *Server) clusterRemovePluginHandler(msg *model.ClusterMessage) {
	var data model.PluginEventData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		slog.Warn("Failed to decode from JSON", slog.Err(err))
	}
	s.PluginService().RemovePluginFromData(data)
}

func (s *Server) clusterPluginEventHandler(msg *model.ClusterMessage) {
	env, appErr := s.PluginService().GetPluginsEnvironment()
	if env == nil || appErr != nil {
		return
	}
	if msg.Props == nil {
		slog.Warn("ClusterMessage.Props for plugin event should not be nil")
		return
	}
	pluginID := msg.Props["PluginID"]
	eventID := msg.Props["EventID"]
	if pluginID == "" || eventID == "" {
		slog.Warn("Invalid ClusterMessage.Props values for plugin event",
			slog.String("plugin_id", pluginID), slog.String("event_id", eventID))
		return
	}

	hooks, err := env.HooksForPlugin(pluginID)
	if err != nil {
		slog.Warn("Getting hooks for plugin failed", slog.String("plugin_id", pluginID), slog.Err(err))
		return
	}

	hooks.OnPluginClusterEvent(&plugin.Context{}, model.PluginClusterEvent{
		Id:   eventID,
		Data: msg.Data,
	})
}

// registerClusterHandlers registers the cluster message handlers that are handled by the server.
//
// The cluster event handlers are spread across this function and NewLocalCacheLayer.
// Be careful to not have duplicated handlers here and there.
func (s *Server) registerClusterHandlers() {
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventPublish, s.clusterPublishHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventUpdateStatus, s.clusterUpdateStatusHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInvalidateAllCaches, s.clusterInvalidateAllCachesHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInvalidateCacheForUser, s.clusterInvalidateCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventBusyStateChanged, s.clusterBusyStateChgHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventClearSessionCacheForUser, s.clusterClearSessionCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventClearSessionCacheForAllUsers, s.clusterClearSessionCacheForAllUsersHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInstallPlugin, s.clusterInstallPluginHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventRemovePlugin, s.clusterRemovePluginHandler)
	s.Cluster.RegisterClusterMessageHandler(model.ClusterEventPluginEvent, s.clusterPluginEventHandler)
	// s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInvalidateCacheForChannelMembersNotifyProps, s.clusterInvalidateCacheForChannelMembersNotifyPropHandler)
	// s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInvalidateCacheForChannelByName, s.clusterInvalidateCacheForChannelByNameHandler)
	// s.Cluster.RegisterClusterMessageHandler(model.ClusterEventInvalidateCacheForUserTeams, s.clusterInvalidateCacheForUserTeamsHandler)
}

func (s *Server) clusterBusyStateChgHandler(msg *model.ClusterMessage) {
	var sbs model.ServerBusyState
	if jsonErr := json.Unmarshal(msg.Data, &sbs); jsonErr != nil {
		slog.Warn("Failed to decode server busy state from JSON", slog.Err(jsonErr))
	}
	s.serverBusyStateChanged(&sbs)
}

func (s *Server) clusterPublishHandler(msg *model.ClusterMessage) {
	event := model.WebSocketEventFromJson(bytes.NewReader(msg.Data))
	if event == nil {
		return
	}
	s.PublishSkipClusterSend(event)
}

func (s *Server) clusterInvalidateCacheForUserHandler(msg *model.ClusterMessage) {
	s.invalidateCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterUpdateStatusHandler(msg *model.ClusterMessage) {
	status := model.StatusFromJson(bytes.NewReader(msg.Data))
	s.StatusCache.Set(status.UserId, status)
}

func (s *Server) clearSessionCacheForAllUsersSkipClusterSend() {
	slog.Info("Purging sessions cache")
	s.AccountService().ClearAllUsersSessionCacheLocal()
}

func (s *Server) clusterInvalidateAllCachesHandler(msg *model.ClusterMessage) {
	s.InvalidateAllCachesSkipSend()
}

// invalidateCacheForUserSkipClusterSend
func (s *Server) invalidateCacheForUserSkipClusterSend(userID string) {

	s.invalidateWebConnSessionCacheForUser(userID)
}

func (s *Server) clusterClearSessionCacheForUserHandler(msg *model.ClusterMessage) {
	s.clearSessionCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterClearSessionCacheForAllUsersHandler(msg *model.ClusterMessage) {
	s.clearSessionCacheForAllUsersSkipClusterSend()
}

// func (s *Server) clusterBusyStateChgHandler(msg *model.ClusterMessage) {
// 	s.serverBusyStateChanged(model.ServerBusyStateFromJson(bytes.NewReader(msg.Data)))
// }

// invalidateWebConnSessionCacheForUser
func (s *Server) invalidateWebConnSessionCacheForUser(userID string) {
	panic("not implt")
}

func (s *Server) clearSessionCacheForUserSkipClusterSend(userID string) {
	s.AccountService().ClearUserSessionCacheLocal(userID)
	s.invalidateWebConnSessionCacheForUser(userID)
}
