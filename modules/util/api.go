package util

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/sitename/sitename/model"
)

func RenderWebAppError(config *model.Config, w http.ResponseWriter, r *http.Request, err *model.AppError, s crypto.Signer) {
	RenderWebError(config, w, r, err.StatusCode, url.Values{
		"message": []string{err.Message},
	}, s)
}

func RenderWebError(config *model.Config, w http.ResponseWriter, r *http.Request, status int, params url.Values, s crypto.Signer) {
	queryString := params.Encode()

	subpath, _ := GetSubpathFromConfig(config)

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
