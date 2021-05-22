package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/graph/generated"
)

func NewHandler(app app.AppIface) http.Handler {
	c := &generated.Config{
		Resolvers: &Resolver{
			app: app,
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(*c))

	return srv
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
