package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graph"
	"github.com/sitename/sitename/services/configservice"
)

const (
	rootApiPath = "/api"
	graphqlPath = "/graphql"
)

type Routes struct {
	GraphqlAPI *mux.Router
}

type API struct {
	ConfigService       configservice.ConfigService
	GetGlobalAppOptions app.AppOptionCreator
	BaseRoutes          *Routes
}

// init graphql and other api routes
func NewAPI(configservice configservice.ConfigService, globalOptionsFunc app.AppOptionCreator, root *mux.Router) *API {
	api := &API{
		ConfigService:       configservice,
		GetGlobalAppOptions: globalOptionsFunc,
		BaseRoutes:          &Routes{},
	}

	// register routes to graphql api
	api.BaseRoutes.GraphqlAPI = root.PathPrefix(rootApiPath).Subrouter()
	api.BaseRoutes.GraphqlAPI.Handle("", graph.NewPlaygroundHandler(graphqlPath)).Methods(http.MethodGet)
	api.BaseRoutes.GraphqlAPI.Handle(graphqlPath, graph.NewHandler(nil)).Methods(http.MethodPost, http.MethodPut, http.MethodOptions)

	return api
}
