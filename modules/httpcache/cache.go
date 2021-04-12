package httpcache

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sitename/sitename/modules/setting"
)

// GetCacheControl returns a suitable "Cache-Control" header value
func GetCacheControl() string {
	if !setting.IsProd() {
		return "no-store"
	}
	return "private, max-age=" + strconv.FormatInt(int64(setting.StaticCacheTime.Seconds()), 10)
}

// generateETag generates an ETag based on size, filename and file modification time
func generateETag(fi os.FileInfo) string {
	etag := fmt.Sprint(fi.Size()) + fi.Name() + fi.ModTime().UTC().Format(http.TimeFormat)
	return base64.StdEncoding.EncodeToString([]byte(etag))
}

// HandleTimeCache handles time-based caching for a HTTP request
func HandleTimeCache(r *http.Request, w http.ResponseWriter, fi os.FileInfo) bool {
	ifModifiedSince := r.Header.Get("If-Modified_Since")
	if ifModifiedSince != "" {
		t, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err != nil && fi.ModTime().Unix() <= t.Unix() {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}

	w.Header().Set("Cache-Control", GetCacheControl())
	w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))
	return false
}

// HandleEtagCache handles ETag-based caching for a HTTP request
func HandleEtagCache(req *http.Request, w http.ResponseWriter, fi os.FileInfo) (handled bool) {
	etag := generateETag(fi)
	if req.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	w.Header().Set("Cache-Control", GetCacheControl())
	w.Header().Set("ETag", etag)
	return false
}
