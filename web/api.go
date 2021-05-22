package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graph"
)

const (
	rootApiPath string = "/api"
	graphqlPath string = "/graphql"
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

	playgroundHandler := graph.NewPlaygroundHandler(graphqlPath)
	logicHandler := graph.NewHandler(w.app)

	// register routes to graphql api
	api.BaseRoutes.GraphqlAPI = root.PathPrefix(rootApiPath).Subrouter()
	// playground handler
	api.
		BaseRoutes.
		GraphqlAPI.
		Handle("", w.NewHandler(graphqlHanlerWrapper(playgroundHandler))).
		Methods(http.MethodGet, http.MethodOptions)

	// api request handler
	api.
		BaseRoutes.
		GraphqlAPI.
		Handle(graphqlPath, w.NewHandler(graphqlHanlerWrapper(logicHandler))).
		Methods(http.MethodPost, http.MethodOptions)

	return api
}

func graphqlHanlerWrapper(playgroundHandler http.Handler) func(c *Context, w http.ResponseWriter, r *http.Request) {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {
		playgroundHandler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), graph.ApiContextKey, c)))
	}
}
