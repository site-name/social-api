package app

import (
	"bytes"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) clusterInstallPluginHandler(msg *cluster.ClusterMessage) {
	s.PluginService().InstallPluginFromData(plugins.PluginEventDataFromJson(bytes.NewReader(msg.Data)))
}

func (s *Server) clusterRemovePluginHandler(msg *cluster.ClusterMessage) {
	s.PluginService().RemovePluginFromData(plugins.PluginEventDataFromJson(bytes.NewReader(msg.Data)))
}

func (s *Server) clusterPluginEventHandler(msg *cluster.ClusterMessage) {
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

	hooks.OnPluginClusterEvent(&plugin.Context{}, plugins.PluginClusterEvent{
		Id:   eventID,
		Data: msg.Data,
	})
}

// registerClusterHandlers registers the cluster message handlers that are handled by the server.
//
// The cluster event handlers are spread across this function and NewLocalCacheLayer.
// Be careful to not have duplicated handlers here and there.
func (s *Server) registerClusterHandlers() {
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventPublish, s.clusterPublishHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventUpdateStatus, s.clusterUpdateStatusHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInvalidateAllCaches, s.clusterInvalidateAllCachesHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInvalidateCacheForUser, s.clusterInvalidateCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventBusyStateChanged, s.clusterBusyStateChgHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventClearSessionCacheForUser, s.clusterClearSessionCacheForUserHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventClearSessionCacheForAllUsers, s.clusterClearSessionCacheForAllUsersHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInstallPlugin, s.clusterInstallPluginHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventRemovePlugin, s.clusterRemovePluginHandler)
	s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventPluginEvent, s.clusterPluginEventHandler)
	// s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInvalidateCacheForChannelMembersNotifyProps, s.clusterInvalidateCacheForChannelMembersNotifyPropHandler)
	// s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInvalidateCacheForChannelByName, s.clusterInvalidateCacheForChannelByNameHandler)
	// s.Cluster.RegisterClusterMessageHandler(cluster.ClusterEventInvalidateCacheForUserTeams, s.clusterInvalidateCacheForUserTeamsHandler)
}

func (s *Server) clusterBusyStateChgHandler(msg *cluster.ClusterMessage) {
	var sbs model.ServerBusyState
	if jsonErr := json.JSON.Unmarshal(msg.Data, &sbs); jsonErr != nil {
		slog.Warn("Failed to decode server busy state from JSON", slog.Err(jsonErr))
	}
	s.serverBusyStateChanged(&sbs)
}

func (s *Server) clusterPublishHandler(msg *cluster.ClusterMessage) {
	event := model.WebSocketEventFromJson(bytes.NewReader(msg.Data))
	if event == nil {
		return
	}
	s.PublishSkipClusterSend(event)
}

func (s *Server) clusterInvalidateCacheForUserHandler(msg *cluster.ClusterMessage) {
	s.invalidateCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterUpdateStatusHandler(msg *cluster.ClusterMessage) {
	status := account.StatusFromJson(bytes.NewReader(msg.Data))
	s.StatusCache.Set(status.UserId, status)
}

func (s *Server) clearSessionCacheForAllUsersSkipClusterSend() {
	slog.Info("Purging sessions cache")
	s.AccountService().ClearAllUsersSessionCacheLocal()
}

func (s *Server) clusterInvalidateAllCachesHandler(msg *cluster.ClusterMessage) {
	s.InvalidateAllCachesSkipSend()
}

// invalidateCacheForUserSkipClusterSend
func (s *Server) invalidateCacheForUserSkipClusterSend(userID string) {

	s.invalidateWebConnSessionCacheForUser(userID)
}

func (s *Server) clusterClearSessionCacheForUserHandler(msg *cluster.ClusterMessage) {
	s.clearSessionCacheForUserSkipClusterSend(string(msg.Data))
}

func (s *Server) clusterClearSessionCacheForAllUsersHandler(msg *cluster.ClusterMessage) {
	s.clearSessionCacheForAllUsersSkipClusterSend()
}

// func (s *Server) clusterBusyStateChgHandler(msg *cluster.ClusterMessage) {
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
