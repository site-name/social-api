package app

import (
	"context"
	"crypto/ecdsa"
	"net/http"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

type App struct {
	srv            *Server
	t              i18n.TranslateFunc
	session        model.Session
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

	a.srv.Jobs.InitSchedulers()
	a.srv.Jobs.InitWorkers()
}

func (a *App) Srv() *Server {
	return a.srv
}

func (a *App) Log() *slog.Logger {
	return a.srv.Log
}

func (a *App) NotificationsLog() *slog.Logger {
	return a.srv.NotificationsLog
}

func (a *App) RequestId() string {
	return a.requestId
}

func (a *App) IpAddress() string {
	return a.ipAddress
}

func (a *App) Config() *model.Config {
	return a.Srv().Config()
}

func (a *App) Metrics() einterfaces.MetricsInterface {
	return a.srv.Metrics
}

func (a *App) Ldap() einterfaces.LdapInterface {
	return a.srv.Ldap
}

func (a *App) Path() string {
	return a.path
}

func (a *App) UserAgent() string {
	return a.userAgent
}

func (a *App) AcceptLanguage() string {
	return a.acceptLanguage
}

func (a *App) TelemetryId() string {
	return a.Srv().TelemetryId()
}

func (a *App) Handle404(w http.ResponseWriter, r *http.Request) {
	ipAddress := util.GetIPAddress(r, a.Config().ServiceSettings.TrustedProxyIPHeader)
	slog.Debug("not found handler triggered", slog.String("path", r.URL.Path), slog.Int("code", 404), slog.String("ip", ipAddress))

	if *a.Config().ServiceSettings.WebserverMode == "disabled" {
		http.NotFound(w, r)
	}

	util.RenderWebAppError(a.Config(), w, r, model.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound), a.AsymmetricSigningKey())
}

func (a *App) AsymmetricSigningKey() *ecdsa.PrivateKey {
	return a.Srv().AsymmetricSigningKey()
}
