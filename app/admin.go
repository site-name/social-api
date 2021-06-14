package app

import (
	"runtime/debug"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) InvalidateAllCaches() *model.AppError {
	debug.FreeOSMemory()
	s.InvalidateAllCachesSkipSend()

	if s.Cluster != nil {

		msg := &model.ClusterMessage{
			Event:            model.CLUSTER_EVENT_INVALIDATE_ALL_CACHES,
			SendType:         model.CLUSTER_SEND_RELIABLE,
			WaitForAllToSend: true,
		}

		s.Cluster.SendClusterMessage(msg)
	}

	return nil
}

func (s *Server) InvalidateAllCachesSkipSend() {
	slog.Info("Purging all caches")
	s.sessionCache.Purge()
	s.statusCache.Purge()
	// s.Store.Team().ClearCaches()
	// s.Store.Channel().ClearCaches()
	s.Store.User().ClearCaches()
	// s.Store.Post().ClearCaches()
	s.Store.FileInfo().ClearCaches()
	// s.Store.Webhook().ClearCaches()
	// s.LoadLicense()
}
