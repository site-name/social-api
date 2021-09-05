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
			Event:            cluster.ClusterEventInvalidateAllCaches,
			SendType:         cluster.ClusterSendReliable,
			WaitForAllToSend: true,
		}

		s.Cluster.SendClusterMessage(msg)
	}

	return nil
}

func (s *Server) InvalidateAllCachesSkipSend() {
	slog.Info("Purging all caches")
	s.AccountService().ClearAllUsersSessionCacheLocal()
	s.StatusCache.Purge()
	s.Store.User().ClearCaches()
	s.Store.FileInfo().ClearCaches()
	// s.Store.Webhook().ClearCaches()
	// s.Store.Post().ClearCaches()
}

// serverBusyStateChanged is called when a CLUSTER_EVENT_BUSY_STATE_CHANGED is received.
func (s *Server) serverBusyStateChanged(sbs *model.ServerBusyState) {
	s.Busy.ClusterEventChanged(sbs)
	if sbs.Busy {
		slog.Warn("server busy state activitated via cluster event - non-critical services disabled", slog.Int64("expires_sec", sbs.Expires))
	} else {
		slog.Info("server busy state cleared via cluster event - non-critical services enabled")
	}
}
