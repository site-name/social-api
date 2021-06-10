package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) PluginUpdate(ctx context.Context, channelID *string, id string, input PluginUpdateInput) (*PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugin(ctx context.Context, id string) (*Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Plugins(ctx context.Context, filter *PluginFilterInput, sortBy *PluginSortingInput, before *string, after *string, first *int, last *int) (*PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
