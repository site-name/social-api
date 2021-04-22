package app

import (
	"context"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

type App struct {
	srv            *Server
	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string
	context        context.Context
}

func New(options ...AppOption) *App {
	app := new(App)
	for _, option := range options {
		option(app)
	}

	return app
}

func (a *App) InitServer() {
	a.srv.AppInitializedOnce.Do(func() {
		// a.initEnterprise()

		a.AddConfigListener(func(oldConfig, newConfig *model.Config) {
			if *oldConfig.GuestAccountsSettings.Enable && !*newConfig.GuestAccountsSettings.Enable {
				if appErr := a.DeactivateGuests(); appErr != nil {
					slog.Error("Unable to deactivate guest accounts", slog.Err(appErr))
				}
			}
		})

		// Disable active guest accounts on first run if guest accounts are disabled
		if !*a.Config().GuestAccountsSettings.Enable {
			if appErr := a.DeactivateGuests(); appErr != nil {
				slog.Error("Uable to deactivate guest accounts", slog.Err(appErr))
			}
		}

		// Scheduler must be started before cluster.
		a.initJobs()

		// if a.srv.joinCluster && a.srv.Cluster != nil {
		// 	a.registerAppClusterMessageHandlers()
		// }

		// a.DoAppMigrations()

		// a.InitPostMetadata()

		// a.InitPlugins(*a.Config().PluginSettings.Directory, *a.Config().PluginSettings.ClientDirectory)

		// a.AddConfigListener(func(prevCfg, cfg *model.Config) {
		// 	if *cfg.PluginSettings.Enable {
		// 		a.InitPlugins(*cfg.PluginSettings.Directory, *a.Config().PluginSettings.ClientDirectory)
		// 	} else {
		// 		a.srv.ShutDownPlugins()
		// 	}
		// })

		if a.Srv().runEssentialJobs {
			// a.Srv().Go(func() {

			// })
			a.Srv().runJobs()
		}
	})
}

func (a *App) initJobs() {
	if jobsLdapSyncInterface != nil {
		a.srv.Jobs.LdapSync = jobsLdapSyncInterface(a)
	}
}

func (a *App) Srv() *Server {
	return a.srv
}

func (a *App) Config() *model.Config {
	return a.Srv().Config()
}

func (a *App) Metrics() einterfaces.MetricsInterface {
	return a.srv.Metrics
}
