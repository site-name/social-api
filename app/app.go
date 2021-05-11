package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"
)

type App struct {
	srv *Server

	// XXX: This is required because removing this needs BleveEngine
	// to be registered in (h *MainHelper) setupStore, but that creates
	// a cyclic dependency as bleve tests themselves import testlib.
	searchEngine *searchengine.Broker

	t              i18n.TranslateFunc
	session        model.Session
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string
	context        context.Context
	requestId      string
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

		a.DoAppMigrations()

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

func (a *App) Handle404(w http.ResponseWriter, r *http.Request) {
	ipAddress := util.GetIPAddress(r, a.Config().ServiceSettings.TrustedProxyIPHeader)
	slog.Debug("not found handler triggered", slog.String("path", r.URL.Path), slog.Int("code", 404), slog.String("ip", ipAddress))

	if *a.Config().ServiceSettings.WebserverMode == "disabled" {
		http.NotFound(w, r)
	}

	util.RenderWebAppError(a.Config(), w, r, model.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound), a.AsymmetricSigningKey())
}

func (s *Server) getSystemInstallDate() (int64, *model.AppError) {
	systemData, err := s.Store.System().GetByName(model.SYSTEM_INSTALLATION_DATE_KEY)
	if err != nil {
		return 0, model.NewAppError("getSystemInstallDate", "app.system.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	value, err := strconv.ParseInt(systemData.Value, 10, 64)
	if err != nil {
		return 0, model.NewAppError("getSystemInstallDate", "app.system_install_date.parse_int.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return value, nil
}

func (a *App) Cluster() einterfaces.ClusterInterface {
	return a.srv.Cluster
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

func (a *App) T(translationID string, args ...interface{}) string {
	return a.t(translationID, args...)
}

func (a *App) Session() *model.Session {
	return &a.session
}

func (a *App) RequestId() string {
	return a.requestId
}

func (a *App) IpAddress() string {
	return a.ipAddress
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

func (a *App) Compliance() einterfaces.ComplianceInterface {
	return a.srv.Compliance
}

// func (a *App) DataRetention() einterfaces.DataRetentionInterface {
// 	return a.srv.DataRetention
// }
func (a *App) Metrics() einterfaces.MetricsInterface {
	return a.srv.Metrics
}

func (a *App) SearchEngine() *searchengine.Broker {
	return a.searchEngine
}

func (a *App) Ldap() einterfaces.LdapInterface {
	return a.srv.Ldap
}

// func (a *App) Notification() einterfaces.NotificationInterface {
// 	return a.srv.Notification
// }

func (a *App) HTTPService() httpservice.HTTPService {
	return a.srv.HTTPService
}

func (a *App) ImageProxy() *imageproxy.ImageProxy {
	return a.srv.ImageProxy
}

func (a *App) Timezones() *timezones.Timezones {
	return a.srv.timezones
}

func (a *App) Context() context.Context {
	return a.context
}

func (a *App) SetSession(s *model.Session) {
	a.session = *s
}

func (a *App) SetT(t i18n.TranslateFunc) {
	a.t = t
}

func (a *App) SetRequestId(s string) {
	a.requestId = s
}

func (a *App) SetIpAddress(s string) {
	a.ipAddress = s
}

func (a *App) SetUserAgent(s string) {
	a.userAgent = s
}

func (a *App) SetAcceptLanguage(s string) {
	a.acceptLanguage = s
}

func (a *App) SetPath(s string) {
	a.path = s
}

func (a *App) SetContext(c context.Context) {
	a.context = c
}

func (a *App) SetServer(srv *Server) {
	a.srv = srv
}

func (a *App) GetT() i18n.TranslateFunc {
	return a.t
}

func (a *App) DBHealthCheckWrite() error {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)

	return a.Srv().Store.System().SaveOrUpdate(&model.System{
		Name:  a.dbHealthCheckKey(),
		Value: currentTime,
	})
}

func (a *App) DBHealthCheckDelete() error {
	_, err := a.Srv().Store.System().PermanentDeleteByName(a.dbHealthCheckKey())
	return err
}

func (a *App) dbHealthCheckKey() string {
	return fmt.Sprintf("health_check_%s", a.GetClusterId())
}
