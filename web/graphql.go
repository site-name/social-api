package web

import (
	"net/http"

	"github.com/sitename/sitename/graph"
)

func (w *Web) InitGraphqlAPI() {
	w.MainRouter.Handle("/graphql", graph.NewHandler(nil)).Methods(http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut)

	w.MainRouter.Handle("/", graph.NewPlaygroundHandler("/")).Methods(http.MethodGet)
}
