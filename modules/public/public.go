package public

import (
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/modules/httpcache"
	"github.com/sitename/sitename/modules/setting"
)

// Options represents the available options to configure the handler.
type Options struct {
	Directory   string
	IndexFile   string
	SkipLogging bool
	FileSystem  http.FileSystem
	Prefix      string
}

// KnownPublicEntries list all direct children in the `public` directory
var KnownPublicEntries = []string{
	"css",
	"fonts",
	"img",
	"js",
	"serviceworker.js",
	"vendor",
	"favicon.ico",
}

// Custom implements the static handler for serving custom assets.
func Custom(opts *Options) func(next http.Handler) http.Handler {
	return opts.staticHandler(path.Join(setting.CustomPath, "public"))
}

// staticFileSystem implements http.FileSystem interface.
type staticFileSystem struct {
	dir *http.Dir
}

func newStaticFileSystem(dir string) staticFileSystem {
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(setting.AppWorkPath, dir)
	}
	direc := http.Dir(dir)
	return staticFileSystem{&direc}
}

func (fs staticFileSystem) Open(name string) (http.File, error) {
	return fs.dir.Open(name)
}

func (opts *Options) staticHandler(dir string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// defaults
		if len(opts.IndexFile) == 0 {
			opts.IndexFile = "index.html"
		}
		// Normalize the prefix if provided
		if opts.Prefix != "" {
			// Ensure we have a leading '/'
			if opts.Prefix[0] != '/' {
				opts.Prefix = "/" + opts.Prefix
			}

			// Remove any trailing '/'
			opts.Prefix = strings.TrimRight(opts.Prefix, "/")
		}
		if opts.FileSystem == nil {
			opts.FileSystem = newStaticFileSystem(dir)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !opts.handle(w, r, opts) {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// parseAcceptEncoding parse Accept-Encoding: deflate, gzip;q=1.0, *;q=0.5 as compress methods
func parseAcceptEncoding(val string) map[string]bool {
	parts := strings.Split(val, ";")
	var types = make(map[string]bool)
	for _, v := range strings.Split(parts[0], ",") {
		types[strings.TrimSpace(v)] = true
	}
	return types
}

func (opts *Options) handle(w http.ResponseWriter, r *http.Request, opt *Options) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}

	file := r.URL.Path
	// i we have a prefix, filter requests by stripping the prefix
	if opt.Prefix != "" {
		if !strings.HasPrefix(file, opts.Prefix) {
			return false
		}
		file = file[len(opt.Prefix):]
		if file != "" && file[0] != '/' {
			return false
		}
	}

	f, err := opt.FileSystem.Open(file)
	if err != nil {
		// 404 request to any known entries in 'public'
		if path.Base(opts.Directory) == "public" {
			parts := strings.Split(file, "/")
			if len(parts) < 2 {
				return false
			}
			for _, entry := range KnownPublicEntries {
				if entry == parts[1] {
					w.WriteHeader(http.StatusNotFound)
					return true
				}
			}
		}
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("[Static] %q exists, but failes to open: %v", file, err)
		return true
	}

	// Try to serve index file
	if fi.IsDir() {
		// Redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, path.Clean(r.URL.Path+"/"), http.StatusFound)
			return true
		}

		f, err = opt.FileSystem.Open(file)
		if err != nil {
			return false // Discard error
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return false
		}
	}

	if !opt.SkipLogging {
		log.Println("[Static] Serving " + file)
	}

	if httpcache.HandleEtagCache(r, w, fi) {
		return true
	}

	ServeContent(w, r, fi, fi.ModTime(), f)
	return true
}
