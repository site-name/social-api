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

// clearSessionCacheForUserSkipClusterSend iterates through server's sessionCache, if it finds any session belong to given userID, removes that session.
func (s *Server) clearSessionCacheForUserSkipClusterSend(userID string) {
	if keys, err := s.sessionCache.Keys(); err == nil {
		var session *model.Session
		for _, key := range keys {
			if err := s.sessionCache.Get(key, &session); err == nil {
				if session.UserId == userID {
					s.sessionCache.Remove(key)
					if s.Metrics != nil {
						s.Metrics.IncrementMemCacheInvalidationCounterSession()
					}
				}
			}
		}
	}

	s.invalidateWebConnSessionCacheForUser(userID)
}
