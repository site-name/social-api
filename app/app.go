package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/timezones"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/modules/util/api"
	"github.com/sitename/sitename/services/httpservice"
	"github.com/sitename/sitename/services/imageproxy"
	"github.com/sitename/sitename/services/searchengine"
)

type App struct {
	srv *Server
}

// New creates new app and returns it
func New(options ...AppOption) *App {
	app := &App{}

	for _, option := range options {
		option(app)
	}

	return app
}

func (a *App) Handle404(w http.ResponseWriter, r *http.Request) {
	ipAddress := util.GetIPAddress(r, a.Config().ServiceSettings.TrustedProxyIPHeader)
	slog.Debug("not found handler triggered", slog.String("path", r.URL.Path), slog.Int("code", 404), slog.String("ip", ipAddress))

	if *a.Config().ServiceSettings.WebserverMode == "disabled" {
		http.NotFound(w, r)
		return
	}

	api.RenderWebAppError(a.Config(), w, r, model_helper.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound), a.AsymmetricSigningKey())
}

func (a *App) GetSystemInstallDate() (int64, *model_helper.AppError) {
	return a.Srv().getSystemInstallDate()
}

func (s *Server) getSystemInstallDate() (int64, *model_helper.AppError) {
	systemData, err := s.Store.System().GetByName(model_helper.SystemInstallationDateKey)
	if err != nil {
		return 0, model_helper.NewAppError("getSystemInstallDate", "app.system.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	value, err := strconv.ParseInt(systemData.Value, 10, 64)
	if err != nil {
		return 0, model_helper.NewAppError("getSystemInstallDate", "app.system_install_date.parse_int.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return value, nil
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

func (a *App) AccountMigration() einterfaces.AccountMigrationInterface {
	return a.srv.AccountMigration
}

func (a *App) Cluster() einterfaces.ClusterInterface {
	return a.srv.Cluster
}

func (a *App) Compliance() einterfaces.ComplianceInterface {
	return a.srv.Compliance
}

func (a *App) DataRetention() einterfaces.DataRetentionInterface {
	return a.srv.DataRetention
}

func (a *App) SearchEngine() *searchengine.Broker {
	return a.srv.SearchEngine
}

func (a *App) Ldap() einterfaces.LdapInterface {
	return a.srv.Ldap
}

func (a *App) Saml() einterfaces.SamlInterface {
	return a.srv.Saml
}

func (a *App) Metrics() einterfaces.MetricsInterface {
	return a.srv.Metrics
}

func (a *App) HTTPService() httpservice.HTTPService {
	return a.srv.HTTPService
}

func (a *App) ImageProxy() *imageproxy.ImageProxy {
	return a.srv.ImageProxy
}

func (a *App) Timezones() *timezones.Timezones {
	return a.srv.timezones
}

func (a *App) DBHealthCheckWrite() error {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)

	return a.Srv().Store.System().SaveOrUpdate(model.System{
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

func (a *App) SetServer(srv *Server) {
	a.srv = srv
}

// func (a *App) Notification() einterfaces.NotificationInterface {
// 	return a.srv.Notification
// }

//	func (a *App) Cloud() einterfaces.CloudInterface {
//		return a.srv.Cloud
//	}

// func (a *App) TelemetryId() string {
// 	return a.Srv().TelemetryId()
// }
