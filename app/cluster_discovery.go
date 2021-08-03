package app

import (
	"time"

	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/modules/slog"
)

const (
	DiscoveryServiceWritePing = 60 * time.Second
)

type ClusterDiscoveryService struct {
	cluster.ClusterDiscovery
	srv  *Server
	stop chan bool
}

func (s *Server) NewClusterDiscoveryService() *ClusterDiscoveryService {
	ds := &ClusterDiscoveryService{
		ClusterDiscovery: cluster.ClusterDiscovery{},
		srv:              s,
		stop:             make(chan bool),
	}

	return ds
}

func (a *App) NewClusterDiscoveryService() *ClusterDiscoveryService {
	return a.Srv().NewClusterDiscoveryService()
}

func (cds *ClusterDiscoveryService) Start() {
	err := cds.srv.Store.ClusterDiscovery().Cleanup()
	if err != nil {
		slog.Warn("ClusterDiscoveryService failed to cleanup the outdated cluster discovery information", slog.Err(err))
	}

	exists, err := cds.srv.Store.ClusterDiscovery().Exists(&cds.ClusterDiscovery)
	if err != nil {
		slog.Warn("ClusterDiscoveryService failed to check if row exists", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id), slog.Err(err))
	} else if exists {
		if _, err := cds.srv.Store.ClusterDiscovery().Delete(&cds.ClusterDiscovery); err != nil {
			slog.Warn("ClusterDiscoveryService failed to start clean", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id), slog.Err(err))
		}
	}

	if err := cds.srv.Store.ClusterDiscovery().Save(&cds.ClusterDiscovery); err != nil {
		slog.Error("ClusterDiscoveryService failed to save", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id), slog.Err(err))
		return
	}

	go func() {
		slog.Debug("ClusterDiscoveryService ping writer started", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id))
		ticker := time.NewTicker(DiscoveryServiceWritePing)
		defer func() {
			ticker.Stop()
			if _, err := cds.srv.Store.ClusterDiscovery().Delete(&cds.ClusterDiscovery); err != nil {
				slog.Warn("ClusterDiscoveryService failed to cleanup", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id), slog.Err(err))
			}
			slog.Debug("ClusterDiscoveryService ping writer stopped", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id))
		}()

		for {
			select {
			case <-ticker.C:
				if err := cds.srv.Store.ClusterDiscovery().SetLastPingAt(&cds.ClusterDiscovery); err != nil {
					slog.Error("ClusterDiscoveryService failed to write ping", slog.String("ClusterDiscoveryID", cds.ClusterDiscovery.Id), slog.Err(err))
				}
			case <-cds.stop:
				return
			}
		}
	}()
}

func (cds *ClusterDiscoveryService) Stop() {
	cds.stop <- true
}

func (s *Server) IsLeader() bool {
	if *s.Config().ClusterSettings.Enable && s.Cluster != nil {
		return s.Cluster.IsLeader()
	}
	return true
}

func (a *App) IsLeader() bool {
	return a.Srv().IsLeader()
}

func (a *App) GetClusterId() string {
	if a.Cluster() == nil {
		return ""
	}

	return a.Cluster().GetClusterId()
}
