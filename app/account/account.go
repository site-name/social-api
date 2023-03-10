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

func init() {
	app.RegisterAccountService(func(s *app.Server) (sub_app_iface.AccountService, error) {

		if s.CacheProvider == nil {
			return nil, errors.New("s.CacheProvider must not be nil")
		}

		sessionCache, err := s.CacheProvider.NewCache(&cache.CacheOptions{
			Size:           model.SESSION_CACHE_SIZE,
			Striped:        true,
			StripedBuckets: util.GetMinMax(runtime.NumCPU()-1, 1).Max,
		})
		if err != nil {
			return nil, errors.New("could not create session cache")
		}

		return &ServiceAccount{
			srv:          s,
			sessionCache: sessionCache,
			metrics:      s.Metrics,
			cluster:      s.Cluster,
			sessionPool: sync.Pool{
				New: func() interface{} {
					return &model.Session{}
				},
			},
		}, nil
	})
}
