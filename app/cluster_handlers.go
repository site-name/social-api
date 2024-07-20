package app

import (
	"bytes"
	"encoding/json"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) clusterInstallPluginHandler(msg *model_helper.ClusterMessage) {
	var data model_helper.PluginEventData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		slog.Warn("Failed to decode from JSON", slog.Err(err))
	}
	s.Plugin.InstallPluginFromData(data)
}

func (s *Server) clusterRemovePluginHandler(msg *model_helper.ClusterMessage) {
	var data model_helper.PluginEventData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		slog.Warn("Failed to decode from JSON", slog.Err(err))
	}
	s.Plugin.RemovePluginFromData(data)
}

func (s *Server) clusterPluginEventHandler(msg *model_helper.ClusterMessage) {
	env, appErr := s.Plugin.GetPluginsEnvironment()
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

	hooks.OnPluginClusterEvent(&plugin.Context{}, model_helper.PluginClusterEvent{
		Id:   eventID,
		Data: msg.Data,
	})
}

// registerClusterHandlers registers the cluster message handlers that are handled by the server.
//
// The cluster event handlers are spread across this function and NewLocalCacheLayer.
// Be careful to not have duplicated handlers here and there.
func (s *Server) registerClusterHandlers() {
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventPublish, s.clusterPublishHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventUpdateStatus, s.clusterUpdateStatusHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateAllCaches, s.clusterInvalidateAllCachesHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForUser, s.clusterInvalidateCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventBusyStateChanged, s.clusterBusyStateChgHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventClearSessionCacheForUser, s.clusterClearSessionCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventClearSessionCacheForAllUsers, s.clusterClearSessionCacheForAllUsersHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInstallPlugin, s.clusterInstallPluginHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventRemovePlugin, s.clusterRemovePluginHandler)
	s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventPluginEvent, s.clusterPluginEventHandler)
	// s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForChannelMembersNotifyProps, s.clusterInvalidateCacheForChannelMembersNotifyPropHandler)
	// s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForChannelByName, s.clusterInvalidateCacheForChannelByNameHandler)
	// s.Cluster.RegisterClusterMessageHandler(model_helper.ClusterEventInvalidateCacheForUserTeams, s.clusterInvalidateCacheForUserTeamsHandler)
}

func (s *Server) clusterBusyStateChgHandler(msg *model_helper.ClusterMessage) {
	var sbs model_helper.ServerBusyState
	if jsonErr := json.Unmarshal(msg.Data, &sbs); jsonErr != nil {
		slog.Warn("Failed to decode server busy state from JSON", slog.Err(jsonErr))
	}
	s.serverBusyStateChanged(&sbs)
}

func (s *Server) clusterPublishHandler(msg *model_helper.ClusterMessage) {
	event := model_helper.WebSocketEventFromJson(bytes.NewReader(msg.Data))
	if event == nil {
		return
	}
	s.PublishSkipClusterSend(event)
}

func (s *Server) clusterInvalidateCacheForUserHandler(msg *model_helper.ClusterMessage) {
	s.invalidateCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterUpdateStatusHandler(msg *model_helper.ClusterMessage) {
	status := model_helper.StatusFromJson(msg.Data)
	s.Account.AddStatusCache(status)
}

func (s *Server) clearSessionCacheForAllUsersSkipClusterSend() {
	slog.Info("Purging sessions cache")
	s.Account.ClearAllUsersSessionCacheLocal()
}

func (s *Server) clusterInvalidateAllCachesHandler(msg *model_helper.ClusterMessage) {
	s.InvalidateAllCachesSkipSend()
}

// invalidateCacheForUserSkipClusterSend
func (s *Server) invalidateCacheForUserSkipClusterSend(userID string) {

	s.invalidateWebConnSessionCacheForUser(userID)
}

func (s *Server) clusterClearSessionCacheForUserHandler(msg *model_helper.ClusterMessage) {
	s.clearSessionCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterClearSessionCacheForAllUsersHandler(msg *model_helper.ClusterMessage) {
	s.clearSessionCacheForAllUsersSkipClusterSend()
}

// func (s *Server) clusterBusyStateChgHandler(msg *model_helper.ClusterMessage) {
// 	s.serverBusyStateChanged(model.ServerBusyStateFromJson(bytes.NewReader(msg.Data)))
// }

// invalidateWebConnSessionCacheForUser
func (s *Server) invalidateWebConnSessionCacheForUser(userID string) {
	panic("not implt")
}

func (s *Server) clearSessionCacheForUserSkipClusterSend(userID string) {
	s.Account.ClearUserSessionCacheLocal(userID)
	s.invalidateWebConnSessionCacheForUser(userID)
}
