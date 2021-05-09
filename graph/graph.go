package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graph/generated"
)

func NewHandler(app *app.App) http.Handler {
	c := generated.Config{
		Resolvers: &Resolver{
			app: app,
		},
	}

	srv := handler.New(generated.NewExecutableSchema(c))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(*app.Config().ServiceSettings.ApiCacheSize))

	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
