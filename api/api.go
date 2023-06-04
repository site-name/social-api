package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/mattermost/gziphandler"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/web"
)

type API struct {
	srv    *app.Server
	schema *graphql.Schema
	Router *mux.Router
}

func Init(srv *app.Server) (*API, error) {
	api := &API{
		srv:    srv,
		Router: srv.Router,
	}

	if err := api.InitGraphql(); err != nil {
		return nil, err
	}

	return api, nil
}

type handlerFunc func(*web.Context, http.ResponseWriter, *http.Request)

// APIHandler provides a handler for API endpoints which do not require the user to be logged in order for access to be
// granted.
func (api *API) APIHandler(h handlerFunc) http.Handler {
	handler := &web.Handler{
		Srv:             api.srv,
		HandleFunc:      h,
		HandlerName:     web.GetHandlerName(h),
		RequireSession:  false,
		TrustRequester:  false,
		RequireMfa:      false,
		IsStatic:        false,
		IsLocal:         false,
		DisableWhenBusy: true,
	}
	if *api.srv.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler
}
