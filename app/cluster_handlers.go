package app

import "github.com/sitename/sitename/model"

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
