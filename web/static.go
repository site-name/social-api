package web

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/avct/uasurfer"
	"github.com/mattermost/gziphandler"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/templates"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/modules/util/fileutils"
)

var robotsTxt = []byte("User-agent: *\nDisallow: /\n")

func (w *Web) InitStatic() {
	if *w.srv.Config().ServiceSettings.WebserverMode != "disabled" {
		if err := util.UpdateAssetsSubpathFromConfig(w.srv.Config()); err != nil {
			slog.Error("Failed to update assets subpath from config", slog.Err(err))
		}

		staticDir, _ := fileutils.FindDir(model.CLIENT_DIR)
		slog.Debug("Using client directory", slog.String("clientDir", staticDir))

		subpath, _ := util.GetSubpathFromConfig(w.srv.Config())

		staticHandler := staticFilesHandler(
			http.StripPrefix(
				path.Join(subpath, "static"),
				http.FileServer(http.Dir(staticDir)),
			),
		)

		if *w.srv.Config().ServiceSettings.WebserverMode == "gzip" {
			staticHandler = gziphandler.GzipHandler(staticHandler)
		}

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
		w.Header().Set("Cache-Control", "no-store")
		data := renderUnsupportedBrowser(c.AppContext, r)

		c.App.Srv().TemplatesContainer().Render(w, "unsupported_browser", data)
		return
	}

	if IsApiCall(c.App, r) {
		Handle404(c.App, w, r)
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

func renderUnsupportedBrowser(ctx *request.Context, r *http.Request) templates.Data {

	data := templates.Data{
		Props: map[string]interface{}{
			"DownloadAppOrUpgradeBrowserString": ctx.T("web.error.unsupported_browser.download_app_or_upgrade_browser"),
			"LearnMoreString":                   ctx.T("web.error.unsupported_browser.learn_more"),
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
		data.Props["NoLongerSupportString"] = ctx.T("web.error.unsupported_browser.no_longer_support_version")
	} else {
		data.Props["NoLongerSupportString"] = ctx.T("web.error.unsupported_browser.no_longer_support")
	}

	// // Mattermost app version
	// if isWindows {
	// 	data.Props["App"] = renderMattermostAppWindows(ctx)
	// } else if isMacOSX {
	// 	data.Props["App"] = renderMattermostAppMac(ctx)
	// }

	// Browsers to download
	// Show a link to Safari if you're using safari and it's outdated
	// Can't show on Mac all the time because there's no way to open it via URI
	browsers := []Browser{renderBrowserChrome(ctx), renderBrowserFirefox(ctx)}
	if isSafari {
		browsers = append(browsers, renderBrowserSafari(ctx))
	}
	data.Props["Browsers"] = browsers

	// If on Windows 10, show link to Edge
	if isWindows10 {
		data.Props["SystemBrowser"] = renderSystemBrowserEdge(ctx, r)
	}

	return data

}

func renderBrowserChrome(ctx *request.Context) Browser {
	return Browser{
		"/static/images/browser-icons/chrome.svg",
		ctx.T("web.error.unsupported_browser.browser_title.chrome"),
		ctx.T("web.error.unsupported_browser.min_browser_version.chrome"),
		"http://www.google.com/chrome",
		ctx.T("web.error.unsupported_browser.browser_get_latest.chrome"),
	}
}

func renderBrowserFirefox(ctx *request.Context) Browser {
	return Browser{
		"/static/images/browser-icons/firefox.svg",
		ctx.T("web.error.unsupported_browser.browser_title.firefox"),
		ctx.T("web.error.unsupported_browser.min_browser_version.firefox"),
		"https://www.mozilla.org/firefox/new/",
		ctx.T("web.error.unsupported_browser.browser_get_latest.firefox"),
	}
}

func renderBrowserSafari(ctx *request.Context) Browser {
	return Browser{
		"/static/images/browser-icons/safari.svg",
		ctx.T("web.error.unsupported_browser.browser_title.safari"),
		ctx.T("web.error.unsupported_browser.min_browser_version.safari"),
		"macappstore://showUpdatesPage",
		ctx.T("web.error.unsupported_browser.browser_get_latest.safari"),
	}
}

func renderSystemBrowserEdge(ctx *request.Context, r *http.Request) SystemBrowser {
	return SystemBrowser{
		"/static/images/browser-icons/edge.svg",
		ctx.T("web.error.unsupported_browser.browser_title.edge"),
		ctx.T("web.error.unsupported_browser.min_browser_version.edge"),
		ctx.T("web.error.unsupported_browser.open_system_browser.edge"),
		"microsoft-edge:http://" + r.Host + r.RequestURI, //TODO: Can we get HTTP or HTTPS? If someone's server doesn't have a redirect this won't work
		"ms-settings:defaultapps",
		ctx.T("web.error.unsupported_browser.system_browser_or"),
		ctx.T("web.error.unsupported_browser.system_browser_make_default"),
	}
}
