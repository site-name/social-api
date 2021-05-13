package web

import (
	"net/http"
	"path"
	"strings"

	"github.com/avct/uasurfer"
	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/services/configservice"
)

type Web struct {
	GetGlobalAppOptions app.AppOptionCreator
	ConfigService       configservice.ConfigService
	MainRouter          *mux.Router
}

// New initializes web routes and returns web instance
func New(config configservice.ConfigService, globalOptions app.AppOptionCreator, root *mux.Router) *Web {
	slog.Debug("Initializing web routes")

	web := &Web{
		GetGlobalAppOptions: globalOptions,
		ConfigService:       config,
		MainRouter:          root,
	}

	// web.InitOAuth()
	// web.InitWebhooks()
	// web.InitSaml()
	web.InitStatic()

	return web
}

// Due to the complexities of UA detection and the ramifications of a misdetection
// only older Safari and IE browsers throw incompatibility errors.
// Map should be of minimum required browser version.
// -1 means that the browser is not supported in any version.
var browserMinimumSupported = map[string]int{
	"BrowserIE":     11,
	"BrowserSafari": 12,
}

func CheckClientCompatibility(agentString string) bool {
	ua := uasurfer.Parse(agentString)

	if version, exist := browserMinimumSupported[ua.Browser.Name.String()]; exist && (ua.Browser.Version.Major < version || version < 0) {
		return false
	}

	return true
}

func Handle404(config configservice.ConfigService, w http.ResponseWriter, r *http.Request) {
	err := model.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound)
	ipAddress := util.GetIPAddress(r, config.Config().ServiceSettings.TrustedProxyIPHeader)
	slog.Debug("not found handler triggered", slog.String("path", r.URL.Path), slog.Int("code", 404), slog.String("ip", ipAddress))

	if IsApiCall(config, r) {
		w.WriteHeader(err.StatusCode)
		err.DetailedError = "There doesn't appear to be an api call for the url='" + r.URL.Path + "'.  Typo? are you missing a team_id or user_id as part of the url?"
		w.Write([]byte(err.ToJson()))
	} else if *config.Config().ServiceSettings.WebserverMode == "disabled" {
		http.NotFound(w, r)
	} else {
		util.RenderWebAppError(config.Config(), w, r, err, config.AsymmetricSigningKey())
	}
}

func IsApiCall(config configservice.ConfigService, r *http.Request) bool {
	subpath, _ := util.GetSubpathFromConfig(config.Config())

	return strings.HasPrefix(r.URL.Path, path.Join(subpath, "api")+"/")
}

// func IsWebhookCall(a app.AppIface, r *http.Request) bool {
// 	subpath, _ := util.GetSubpathFromConfig(a.Config())

// 	return strings.HasPrefix(r.URL.Path, path.Join(subpath, "hooks")+"/")
// }

// func IsOAuthApiCall(config configservice.ConfigService, r *http.Request) bool {
// 	subpath, _ := util.GetSubpathFromConfig(config.Config())

// 	if r.Method == "POST" && r.URL.Path == path.Join(subpath, "oauth", "authorize") {
// 		return true
// 	}

// 	if r.URL.Path == path.Join(subpath, "oauth", "apps", "authorized") ||
// 		r.URL.Path == path.Join(subpath, "oauth", "deauthorize") ||
// 		r.URL.Path == path.Join(subpath, "oauth", "access_token") {
// 		return true
// 	}
// 	return false
// }

func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	w.Write([]byte(model.MapToJson(m)))
}
