package api

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/sitename/sitename/model_helper"
)

func RenderWebAppError(config *model_helper.Config, w http.ResponseWriter, r *http.Request, err *model_helper.AppError, s crypto.Signer) {
	RenderWebError(config, w, r, err.StatusCode, url.Values{
		"message": []string{err.Message},
	}, s)
}

func RenderWebError(config *model_helper.Config, w http.ResponseWriter, r *http.Request, status int, params url.Values, s crypto.Signer) {
	queryString := params.Encode()

	subpath, _ := model_helper.GetSubpathFromConfig(config)

	h := crypto.SHA256
	sum := h.New()
	sum.Write([]byte(path.Join(subpath, "error") + "?" + queryString))
	signature, err := s.Sign(rand.Reader, sum.Sum(nil), h)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	destination := path.Join(subpath, "error") + "?" + queryString + "&s=" + base64.URLEncoding.EncodeToString(signature)

	if status >= 300 && status < 400 {
		http.Redirect(w, r, destination, status)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprintln(w, fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
			</head>
			<body onload="window.location = '%s'">
				<noscript>
					<meta http-equiv="refresh" content="0; url=%s">
				</noscript>
				<!-- web error message -->
				<a href="%s" style="color: #c0c0c0;">...</a>
			</body>
		</html>`,
		template.HTMLEscapeString(template.JSEscapeString(destination)),
		template.HTMLEscapeString(destination),
		template.HTMLEscapeString(destination),
	))
}

func OriginChecker(allowedOrigins string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		return CheckOrigin(r, allowedOrigins)
	}
}

// CheckOrigin check if the origin of r is in allowedOrigins or not
func CheckOrigin(r *http.Request, allowedOrigins string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	if allowedOrigins == "*" {
		return true
	}

	for _, allowed := range strings.Split(allowedOrigins, " ") {
		if allowed == origin {
			return true
		}
	}
	return false
}
