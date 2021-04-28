package app

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func GetProtocol(r *http.Request) string {
	if r.Header.Get(model.HEADER_FORWARDED_PROTO) == "https" || r.TLS != nil {
		return "https"
	}
	return "http"
}
