package web

import (
	"github.com/gorilla/mux"
	"github.com/sitename/sitename/app"
)

const (
	rootApiPath   string     = "/api"             // api endpoint
	graphqlPath   string     = "/graphql"         // graphql endpoint
	ApiContextKey ContextKey = "thisIsContextKey" // use this key to get context value
)

type ContextKey string

type Routes struct {
	GraphqlAPI *mux.Router
}

type API struct {
	app        app.AppIface
	BaseRoutes *Routes
}

// InitAPI setups graphql and other api routes
// func (w *Web) InitAPI(root *mux.Router) *API {
// 	api := &API{
// 		app:        w.app,
// 		BaseRoutes: &Routes{},
// 	}

// 	playgroundHandler := playground.Handler("GraphQL Playground", graphqlPath)
// 	logicHandler := handler.NewDefaultServer(NewExecutableSchema(Config{
// 		Resolvers: &Resolver{
// 			app: w.app,
// 		},
// 	}))

// 	// register routes to graphql api
// 	api.BaseRoutes.GraphqlAPI = root.PathPrefix(rootApiPath).Subrouter()
// 	// playground handler
// 	api.
// 		BaseRoutes.
// 		GraphqlAPI.
// 		Handle("", w.NewHandler(graphqlHanlerWrapper(playgroundHandler))).
// 		Methods(http.MethodGet, http.MethodOptions)

// 	// api request handler
// 	api.
// 		BaseRoutes.
// 		GraphqlAPI.
// 		Handle(graphqlPath, w.NewHandler(graphqlHanlerWrapper(logicHandler))).
// 		Methods(http.MethodPost, http.MethodOptions)

// 	return api
// }

// func graphqlHanlerWrapper(handler http.Handler) func(c *Context, w http.ResponseWriter, r *http.Request) {
// 	return func(c *Context, w http.ResponseWriter, r *http.Request) {
// 		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ApiContextKey, c)))
// 	}
// }
