package routes

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"gitea.com/go-chi/captcha"
	"github.com/NYTimes/gziphandler"
	"github.com/chi-middleware/proxy"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sitename/sitename/modules/context"
	"github.com/sitename/sitename/modules/httpcache"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/metrics"
	"github.com/sitename/sitename/modules/public"
	"github.com/sitename/sitename/modules/setting"
	"github.com/sitename/sitename/routers"
	"github.com/tstranex/u2f"
)

const (
	// GzipMinSize represents min size to compress for the body size of response
	GzipMinSize = 1400
)

func commonMiddlewares() []func(http.Handler) http.Handler {
	var handlers = []func(http.Handler) http.Handler{}

	if setting.ReverseProxyLimit > 0 {
		opt := proxy.NewForwardedHeadersOptions().
			WithForwardLimit(setting.ReverseProxyLimit).
			ClearTrustedProxies()
		for _, n := range setting.ReverseProxyTrustedProxies {
			if !strings.Contains(n, "/") {
				opt.AddTrustedProxy(n)
			} else {
				opt.AddTrustedNetwork(n)
			}
		}
		handlers = append(handlers, proxy.ForwardedHeaders(opt))
	}
	handlers = append(handlers, middleware.StripSlashes)

	// if !setting.DisableRouterLog && setting.RouterLogLevel != log.Level(log.NONE) {
	// 	if log.GetLogger("router").GetLevel() <= log.Level(setting.RouterLogLevel) {
	// 		handlers = append(handlers, )
	// 	}
	// }
	handlers = append(handlers, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// Why we need this? The Recovery() will try to render a beautiful
			// error page for user, but the process can still panic again, and other
			// middleware like session also may panic then we have to recover twice
			// and send a simple error page that should not panic any more.
			defer func() {
				if err := recover(); err != nil {
					combinedErr := fmt.Sprintf("PANIC: %v\n%s", err, string(log.Stack(2)))
					log.Error("%v", combinedErr)
					if setting.IsProd() {
						http.Error(resp, http.StatusText(500), 500)
					} else {
						http.Error(resp, combinedErr, 500)
					}
				}
			}()
			next.ServeHTTP(resp, req)
		})
	})

	return handlers
}

// NormalRoutes represents non install routes
func NormalRoutes() *chi.Mux {
	r := chi.NewRouter()

	for _, middleware := range commonMiddlewares() {
		r.Use(middleware)
	}
	r.Mount("/", WebRoutes())
	return r
}

// WebRoutes returns all web routes
func WebRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	// GetHead allows a HEAD request redirect to GET if HEAD method is not defined for that route
	r.Use(middleware.GetHead)

	r.Use(public.Custom(
		&public.Options{
			SkipLogging: setting.DisableRouterLog,
		},
	))
	r.Use(public.Static(
		&public.Options{
			Directory:   path.Join(setting.StaticRootPath, "public"),
			SkipLogging: setting.DisableRouterLog,
		},
	))
	r.Use(securityHeaders())
	if setting.CORSConfig.Enabled {
		r.Use(cors.Handler(cors.Options{
			//Scheme:           setting.CORSConfig.Scheme, // FIXME: the cors middleware needs scheme option
			AllowedOrigins: setting.CORSConfig.AllowDomain,
			//setting.CORSConfig.AllowSubdomain // FIXME: the cors middleware needs allowSubdomain option
			AllowedMethods:   setting.CORSConfig.Methods,
			AllowCredentials: setting.CORSConfig.AllowCredentials,
			MaxAge:           int(setting.CORSConfig.MaxAge.Seconds()),
		}))
	}

	gob.Register(&u2f.Challenge{})

	if setting.EnableGzip {
		h, err := gziphandler.GzipHandlerWithOpts(gziphandler.MinSize(GzipMinSize))
		if err != nil {
			log.Fatal("GzipHandlerWithOpts failed: %v", err)
		}
		r.Use(h)
	}

	if setting.Service.EnableCaptcha {
		r.Use(captcha.Captchaer(context.GetImageCaptcha()))
	}

	if setting.HasRobotsTxt {
		r.Get("/robots.txt", func(w http.ResponseWriter, req *http.Request) {
			filePath := path.Join(setting.CustomPath, "robots.txt")
			fi, err := os.Stat(filePath)
			if err == nil && httpcache.HandleTimeCache(req, w, fi) {
				return
			}
			http.ServeFile(w, req, filePath)
		})
	}

	r.Get("/apple-touch-icon.png", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, path.Join(setting.StaticURLPrefix, "img/apple-touch-icon.png"), 301)
	})

	// metrics endpoint
	if setting.Metrics.Enabled {
		c := metrics.NewCollector()
		prometheus.MustRegister(c)

		r.Get("/metrics", routers.Metrics)
	}

	// for health check
	r.Head("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func securityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// CORB: https://www.chromium.org/Home/chromium-security/corb-for-developers
			// http://stackoverflow.com/a/3146618/244009
			resp.Header().Set("x-content-type-options", "nosniff")
			next.ServeHTTP(resp, req)
		})
	}
}
