package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graph"
)

const (
	rootApiPath = "/api"
	graphqlPath = "/graphql"
)

type Routes struct {
	GraphqlAPI *mux.Router
}

type API struct {
	app        app.AppIface
	BaseRoutes *Routes
}

// init graphql and other api routes
func (w *Web) InitAPI(root *mux.Router) *API {
	api := &API{
		app:        w.app,
		BaseRoutes: &Routes{},
	}

	// register routes to graphql api
	api.BaseRoutes.GraphqlAPI = root.PathPrefix(rootApiPath).Subrouter()
	api.
		BaseRoutes.
		GraphqlAPI.
		Handle("", graph.NewPlaygroundHandler(graphqlPath)).
		Methods(http.MethodGet, http.MethodOptions)
	api.
		BaseRoutes.
		GraphqlAPI.
		Handle(graphqlPath, graph.NewHandler(w.app)).
		Methods(http.MethodPost, http.MethodOptions)

	return api
}
