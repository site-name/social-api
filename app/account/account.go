/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package account

import (
	"errors"
	"runtime"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/cache"
)

type ServiceAccount struct {
	srv          *app.Server
	sessionPool  sync.Pool
	sessionCache cache.Cache

	// optional fields
	metrics einterfaces.MetricsInterface
	cluster einterfaces.ClusterInterface
}

type ServiceAccountConfig struct {
	CacheProvider cache.Provider
	Server        *app.Server

	// optional fields
	Metrics einterfaces.MetricsInterface
	Cluster einterfaces.ClusterInterface
}

// NewServiceAccount initializes account service
func NewServiceAccount(config *ServiceAccountConfig) (sub_app_iface.AccountService, error) {
	if config.CacheProvider == nil {
		return nil, errors.New("config.CacheProvider must not be nil")
	}
	if config.Server == nil {
		return nil, errors.New("config.Server must not be nil")
	}

	sessionCahce, err := config.CacheProvider.NewCache(&cache.CacheOptions{
		Size:           model.SESSION_CACHE_SIZE,
		Striped:        true,
		StripedBuckets: util.Max(runtime.NumCPU()-1, 1),
	})
	if err != nil {
		return nil, errors.New("could not create session cache")
	}

	return &ServiceAccount{
		srv:          config.Server,
		sessionCache: sessionCahce,
		metrics:      config.Metrics,
		cluster:      config.Cluster,
		sessionPool: sync.Pool{
			New: func() interface{} {
				return &model.Session{}
			},
		},
	}, nil
}

func init() {
	app.RegisterAccountApp(func(s *app.Server) (sub_app_iface.AccountService, error) {
		return NewServiceAccount(&ServiceAccountConfig{
			Server:        s,
			CacheProvider: s.CacheProvider,
			Metrics:       s.Metrics,
			Cluster:       s.Cluster,
		})
	})
}
