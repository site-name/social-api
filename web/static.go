package web

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/avct/uasurfer"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/templates"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/modules/util/fileutils"
)

var robotsTxt = []byte("User-agent: *\nDisallow: /\n")

func (w *Web) InitStatic() {
	if *w.app.Config().ServiceSettings.WebserverMode != "disabled" {
		if err := util.UpdateAssetsSubpathFromConfig(w.app.Config()); err != nil {
			slog.Error("Failed to update assets subpath from config", slog.Err(err))
		}

		staticDir, _ := fileutils.FindDir(model.CLIENT_DIR)
		slog.Debug("Using client directory", slog.String("clientDir", staticDir))

		subpath, _ := util.GetSubpathFromConfig(w.app.Config())

		staticHandler := staticFilesHandler(
			http.StripPrefix(
				path.Join(subpath, "static"),
				http.FileServer(http.Dir(staticDir)),
			),
		)
		// pluginHandler := staticFilesHandler(
		// 	http.StripPrefix(
		// 		path.Join(subpath, "static", "plugins"),
		// 		http.FileServer(
		// 			http.Dir(*w.app.Config().PluginSettings.ClientDirectory),
		// 		),
		// 	),
		// )

		if *w.app.Config().ServiceSettings.WebserverMode == "gzip" {
			staticHandler = gziphandler.GzipHandler(staticHandler)
			// pluginHandler = gziphandler.GzipHandler(pluginHandler)
		}

		// w.MainRouter.PathPrefix("/static/plugins/").Handler(pluginHandler)
		w.MainRouter.PathPrefix("/static/").Handler(staticHandler)
		w.MainRouter.Handle("/robots.txt", http.HandlerFunc(robotsHandler))
		w.MainRouter.Handle("/unsupported_browser.js", http.HandlerFunc(unsupportedBrowserScriptHandler))
		w.MainRouter.Handle("/{anything:.*}", w.NewStaticHandler(root)).Methods(http.MethodGet)

		// When a subpath is defined, it's necessary to handle redirects without a
		// trailing slash. We don't want to use StrictSlash on the w.MainRouter and affect
		// all routes, just /subpath -> /subpath/.
		w.MainRouter.HandleFunc("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path += "/"
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		}))
	}
}

func root(c *Context, w http.ResponseWriter, r *http.Request) {

	if !CheckClientCompatibility(r.UserAgent()) {
		renderUnsupportedBrowser(c.App, w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, max-age=31556926, public")

	staticDir, _ := fileutils.FindDir(model.CLIENT_DIR)
	http.ServeFile(w, r, filepath.Join(staticDir, "root.html"))
}

func staticFilesHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//wrap our ResponseWriter with our no-cache 404-handler
		w = &notFoundNoCacheResponseWriter{ResponseWriter: w}

		w.Header().Set("Cache-Control", "max-age=31556926, public")

		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

type notFoundNoCacheResponseWriter struct {
	http.ResponseWriter
}

func (w *notFoundNoCacheResponseWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusNotFound {
		w.Header().Set("Cache-Control", "no-cache, public")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	w.Write(robotsTxt)
}

func unsupportedBrowserScriptHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	templatesDir, _ := templates.GetTemplateDirectory()
	http.ServeFile(w, r, filepath.Join(templatesDir, "unsupported_browser.js"))
}

// Browser describes a browser with a download link
type Browser struct {
	LogoSrc                string
	Title                  string
	SupportedVersionString string
	Src                    string
	GetLatestString        string
}

// SystemBrowser describes a browser but includes 2 links: one to open the local browser, and one to make it default
type SystemBrowser struct {
	LogoSrc                string
	Title                  string
	SupportedVersionString string
	LabelOpen              string
	LinkOpen               string
	LinkMakeDefault        string
	OrString               string
	MakeDefaultString      string
}

func renderUnsupportedBrowser(app app.AppIface, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	data := templates.Data{
		Props: map[string]interface{}{
			"DownloadAppOrUpgradeBrowserString": app.T("web.error.unsupported_browser.download_app_or_upgrade_browser"),
			"LearnMoreString":                   app.T("web.error.unsupported_browser.learn_more"),
		},
	}

	// User Agent info
	ua := uasurfer.Parse(r.UserAgent())
	isWindows := ua.OS.Platform.String() == "PlatformWindows"
	isWindows10 := isWindows && ua.OS.Version.Major == 10
	// isMacOSX := ua.OS.Name.String() == "OSMacOSX" && ua.OS.Version.Major == 10
	isSafari := ua.Browser.Name.String() == "BrowserSafari"

	// Basic heading translations
	if isSafari {
		data.Props["NoLongerSupportString"] = app.T("web.error.unsupported_browser.no_longer_support_version")
	} else {
		data.Props["NoLongerSupportString"] = app.T("web.error.unsupported_browser.no_longer_support")
	}

	// Mattermost app version
	// if isWindows {
	// 	data.Props["App"] = renderMattermostAppWindows(app)
	// } else if isMacOSX {
	// 	data.Props["App"] = renderMattermostAppMac(app)
	// }

	// Browsers to download
	// Show a link to Safari if you're using safari and it's outdated
	// Can't show on Mac all the time because there's no way to open it via URI
	browsers := []Browser{renderBrowserChrome(app), renderBrowserFirefox(app)}
	if isSafari {
		browsers = append(browsers, renderBrowserSafari(app))
	}
	data.Props["Browsers"] = browsers

	// If on Windows 10, show link to Edge
	if isWindows10 {
		data.Props["SystemBrowser"] = renderSystemBrowserEdge(app, r)
	}

	app.Srv().TemplatesContainer().Render(w, "unsupported_browser", data)
}

func renderBrowserChrome(app app.AppIface) Browser {
	return Browser{
		"/static/images/browser-icons/chrome.svg",
		app.T("web.error.unsupported_browser.browser_title.chrome"),
		app.T("web.error.unsupported_browser.min_browser_version.chrome"),
		"http://www.google.com/chrome",
		app.T("web.error.unsupported_browser.browser_get_latest.chrome"),
	}
}

func renderBrowserFirefox(app app.AppIface) Browser {
	return Browser{
		"/static/images/browser-icons/firefox.svg",
		app.T("web.error.unsupported_browser.browser_title.firefox"),
		app.T("web.error.unsupported_browser.min_browser_version.firefox"),
		"https://www.mozilla.org/firefox/new/",
		app.T("web.error.unsupported_browser.browser_get_latest.firefox"),
	}
}

func renderBrowserSafari(app app.AppIface) Browser {
	return Browser{
		"/static/images/browser-icons/safari.svg",
		app.T("web.error.unsupported_browser.browser_title.safari"),
		app.T("web.error.unsupported_browser.min_browser_version.safari"),
		"macappstore://showUpdatesPage",
		app.T("web.error.unsupported_browser.browser_get_latest.safari"),
	}
}

func renderSystemBrowserEdge(app app.AppIface, r *http.Request) SystemBrowser {
	return SystemBrowser{
		"/static/images/browser-icons/edge.svg",
		app.T("web.error.unsupported_browser.browser_title.edge"),
		app.T("web.error.unsupported_browser.min_browser_version.edge"),
		app.T("web.error.unsupported_browser.open_system_browser.edge"),
		"microsoft-edge:http://" + r.Host + r.RequestURI, //TODO: Can we get HTTP or HTTPS? If someone's server doesn't have a redirect this won't work
		"ms-settings:defaultapps",
		app.T("web.error.unsupported_browser.system_browser_or"),
		app.T("web.error.unsupported_browser.system_browser_make_default"),
	}
}
