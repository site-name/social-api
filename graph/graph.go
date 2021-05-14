package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sitename/sitename/graph/generated"
)

func NewHandler(interface{}) http.Handler {
	c := &generated.Config{
		Resolvers: &Resolver{},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(*c))

	return srv
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
