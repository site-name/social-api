package graphql

import (
	"github.com/sitename/sitename/app"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	app.AppIface
}

// __Directive
// func (r *Resolver) __Directive() __DirectiveResolver { return &A{} }

// type A struct{}

// func (r *A) IsRepeatable(ctx context.Context, obj *introspection.Directive) (bool, error) {
// 	return true, nil
// }
