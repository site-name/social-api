package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) PluginUpdate(ctx context.Context, channelID *string, id string, input gqlmodel.PluginUpdateInput) (*gqlmodel.PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugin(ctx context.Context, id string) (*gqlmodel.Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugins(ctx context.Context, filter *gqlmodel.PluginFilterInput, sortBy *gqlmodel.PluginSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
