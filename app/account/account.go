package account

import (
	"errors"
	"runtime"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/services/cache"
)

type ServiceAccount struct {
	srv          *app.Server
	sessionPool  sync.Pool
	sessionCache cache.Cache
	statusCache  cache.Cache

	// optional fields
	metrics einterfaces.MetricsInterface
	cluster einterfaces.ClusterInterface
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		if s.CacheProvider == nil {
			return errors.New("s.CacheProvider must not be nil")
		}

		sessionCache, err := s.CacheProvider.NewCache(&cache.CacheOptions{
			Size:           model_helper.SESSION_CACHE_SIZE,
			Striped:        true,
			StripedBuckets: max(runtime.NumCPU()-1, 1),
		})
		if err != nil {
			return errors.New("could not create session cache")
		}

		statusCache, err := s.CacheProvider.NewCache(&cache.CacheOptions{
			Size:           model_helper.STATUS_CACHE_SIZE,
			Striped:        true,
			StripedBuckets: max(runtime.NumCPU()-1, 1),
		})
		if err != nil {
			return errors.New("could not create status cache")
		}

		s.Account = &ServiceAccount{
			srv:          s,
			sessionCache: sessionCache,
			statusCache:  statusCache,
			metrics:      s.Metrics,
			cluster:      s.Cluster,
			sessionPool: sync.Pool{
				New: func() interface{} {
					return &model.Session{}
				},
			},
		}

		return nil
	})
}
