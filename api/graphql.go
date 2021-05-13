package api

import (
	"net/http"

	"github.com/sitename/sitename/graph"
)

func (a *API) InitGraphql() {
	a.
		BaseRoutes.
		GraphqlAPI.
		Handle("", graph.NewPlaygroundHandler("/graphql")).
		Methods(http.MethodGet)
	a.
		BaseRoutes.
		GraphqlAPI.
		Handle("/graphql", graph.NewHandler(nil)).
		Methods(http.MethodPost, http.MethodPut, http.MethodOptions)
}
