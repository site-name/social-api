package web

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/web/graphql"
	"github.com/sitename/sitename/web/shared"
)

const (
	graphqlApi        string = "/api/graphql"
	graphqlPlayground string = "/playground"
)

// InitGraphql registers graphql playground and graphql api endpoint routes
func (web *Web) InitGraphql() {
	graphqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{
			AppIface: web.app,
		},
	}))
	playgroundHandler := playground.Handler("Sitename", graphqlApi)

	web.MainRouter.Handle(graphqlPlayground, web.NewHandler(commonGraphHandler(playgroundHandler))).Methods(http.MethodGet)
	web.MainRouter.Handle(graphqlApi, web.NewHandler(commonGraphHandler(graphqlServer))).Methods(http.MethodPost)
}

// commonGraphHandler is used for both graphql playground/api
func commonGraphHandler(handler http.Handler) func(c *shared.Context, w http.ResponseWriter, r *http.Request) {
	return func(c *shared.Context, w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), shared.APIContextKey, c)))
	}
}
