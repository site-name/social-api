package web

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/web/generated"
)

const (
	rootApiPath   string     = "/api"
	graphqlPath   string     = "/graphql"
	ApiContextKey ContextKey = "thisIsContextKey"
)

type ContextKey string

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

	playgroundHandler := playground.Handler("GraphQL Playground", graphqlPath)
	logicHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			app: w.app,
		},
	}))

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

func graphqlHanlerWrapper(handler http.Handler) func(c *Context, w http.ResponseWriter, r *http.Request) {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ApiContextKey, c)))
	}
}
