package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func NewHandler(interface{}) http.Handler {
	c := Config{
		Resolvers: &Resolver{
			// app: app,
		},
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(c))

	return srv
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
