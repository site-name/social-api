package web

import (
	"context"

	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/sitename/sitename/app"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	app app.AppIface
}

type S struct{}

func (s *S) IsRepeatable(ctx context.Context, obj *introspection.Directive) (bool, error) {
	return false, nil
}

func (r *Resolver) __Directive() __DirectiveResolver {
	return new(S)
}
