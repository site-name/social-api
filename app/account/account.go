package account

import (
	"runtime"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
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
func NewServiceAccount(config *ServiceAccountConfig) *ServiceAccount {
	config.validate()

	sessionCahce, err := config.CacheProvider.NewCache(&cache.CacheOptions{
		Size:           model.SESSION_CACHE_SIZE,
		Striped:        true,
		StripedBuckets: util.Max(runtime.NumCPU()-1, 1),
	})
	if err != nil {
		slog.Critical("could not create session cache", slog.Err(err))
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
	}
}

func (config *ServiceAccountConfig) validate() {
	if config.CacheProvider == nil {
		slog.Critical("config.CacheProvider must not be nil")
	}
	if config.Server == nil {
		slog.Critical("config.Server must not be nil")
	}
}

// func init() {
// 	app.RegisterAccountApp(func(a app.AppIface) sub_app_iface.AccountService {
// 		return &ServiceAccount{
// 			AppIface: a,
// 			sessionPool: sync.Pool{
// 				New: func() interface{} {
// 					return &model.Session{}
// 				},
// 			},
// 		}
// 	})
// }
