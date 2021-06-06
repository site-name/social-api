package web

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/web/consts"
	"github.com/sitename/sitename/web/graphql"
)

const (
	graphqlApi        string = "/api/graphql"
	graphqlPlayground string = "/playground"
)

func (web *Web) InitGraphql() {
	graphqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{},
	}))
	playgroundHandler := playground.Handler("Sitename", graphqlApi)

	web.MainRouter.Handle(graphqlPlayground, web.NewHandler(commonGraphhanler(playgroundHandler))).Methods(http.MethodGet)
	web.MainRouter.Handle(graphqlApi, web.NewHandler(commonGraphhanler(graphqlServer))).Methods(http.MethodPost)
}

func commonGraphhanler(handler http.Handler) func(c *Context, w http.ResponseWriter, r *http.Request) {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), consts.APIContextKey, c)))
	}
}
