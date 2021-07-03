package app

import (
	"runtime/debug"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) InvalidateAllCaches() *model.AppError {
	debug.FreeOSMemory()
	s.InvalidateAllCachesSkipSend()

	if s.Cluster != nil {

		msg := &cluster.ClusterMessage{
			Event:            cluster.CLUSTER_EVENT_INVALIDATE_ALL_CACHES,
			SendType:         cluster.CLUSTER_SEND_RELIABLE,
			WaitForAllToSend: true,
		}

		s.Cluster.SendClusterMessage(msg)
	}

	return nil
}

func (s *Server) InvalidateAllCachesSkipSend() {
	slog.Info("Purging all caches")
	s.SessionCache.Purge()
	s.StatusCache.Purge()
	s.Store.User().ClearCaches()
	// s.Store.Post().ClearCaches()
	s.Store.FileInfo().ClearCaches()
	// s.Store.Webhook().ClearCaches()
}
