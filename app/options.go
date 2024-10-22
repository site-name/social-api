package app

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

type Option func(s *Server) error

// By default, the app will use the store specified by the configuration. This allows you to
// construct an app with a different store.
//
// The override parameter must be either a store.Store or func(App) store.Store().
func StoreOverride(override any) Option {
	return func(s *Server) error {
		switch o := override.(type) {
		case store.Store:
			s.newStore = func() (store.Store, error) {
				return o, nil
			}
			return nil

		case func(*Server) store.Store:
			s.newStore = func() (store.Store, error) {
				return o(s), nil
			}
			return nil

		default:
			return errors.New("invalid StoreOverride")
		}
	}
}

// Config applies the given config dsn, whether a path to config.json
// or a database connection string. It receives as well a set of
// custom defaults that will be applied for any unset property of the
// config loaded from the dsn on top of the normal defaults
func Config(dsn string, readOnly bool, configDefaults *model_helper.Config) Option {
	return func(s *Server) error {
		configStore, err := config.NewStoreFromDSN(dsn, readOnly, configDefaults, true)
		if err != nil {
			return errors.Wrap(err, "failed to apply Config option")
		}

		s.ConfigStore = configStore
		return nil
	}
}

// ConfigStore applies the given config store, typically to replace the traditional sources with a memory store for testing.
func ConfigStore(configStore *config.Store) Option {
	return func(s *Server) error {
		s.ConfigStore = configStore

		return nil
	}
}

func RunEssentialJobs(s *Server) error {
	s.runEssentialJobs = true

	return nil
}

func JoinCluster(s *Server) error {
	s.joinCluster = true

	return nil
}

func StartMetrics(s *Server) error {
	s.startMetrics = true

	return nil
}

func StartSearchEngine(s *Server) error {
	s.startSearchEngine = true

	return nil
}

func SetLogger(logger *slog.Logger) Option {
	return func(s *Server) error {
		s.Log = logger
		return nil
	}
}

type AppOption func(a *App)
type AppOptionCreator func() []AppOption

// ServerConnector simply healps set Server and SearchEngine to App
func ServerConnector(s *Server) AppOption {
	return func(a *App) {
		a.srv = s
	}
}
