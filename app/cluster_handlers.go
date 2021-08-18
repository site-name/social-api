package app

import (
	"github.com/sitename/sitename/model"
)

// func (s *Server) clusterInstallPluginHandler(msg *model.ClusterMessage) {
// 	s.installPluginFromData(plugins.PluginEventDataFromJson(bytes.NewReader(msg.Data)))
// }

// func (s *Server) clusterRemovePluginHandler(msg *model.ClusterMessage) {
// 	s.removePluginFromData(plugins.PluginEventDataFromJson(bytes.NewReader(msg.Data)))
// }

// registerClusterHandlers registers the cluster message handlers that are handled by the server.
//
// The cluster event handlers are spread across this function and NewLocalCacheLayer.
// Be careful to not have duplicated handlers here and there.
func (s *Server) registerClusterHandlers() {

}

// invalidateCacheForUserSkipClusterSend
func (s *Server) invalidateCacheForUserSkipClusterSend(userID string) {

	s.invalidateWebConnSessionCacheForUser(userID)
}

// invalidateWebConnSessionCacheForUser
func (s *Server) invalidateWebConnSessionCacheForUser(userID string) {
	// TODO: fixme

	// just return for now, need implementation
	return
}

// ClearSessionCacheForUserSkipClusterSend iterates through server's sessionCache, if it finds any session belong to given userID, removes that session.
func (s *Server) ClearSessionCacheForUserSkipClusterSend(userID string) {
	if keys, err := s.SessionCache.Keys(); err == nil {
		var session *model.Session
		for _, key := range keys {
			if err := s.SessionCache.Get(key, &session); err == nil {
				if session.UserId == userID {
					s.SessionCache.Remove(key)
					if s.Metrics != nil {
						s.Metrics.IncrementMemCacheInvalidationCounterSession()
					}
				}
			}
		}
	}

	s.invalidateWebConnSessionCacheForUser(userID)
}
