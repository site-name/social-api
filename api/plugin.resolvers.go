package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PluginUpdate(ctx context.Context, args struct {
	channelID *string
	id        string
	input     gqlmodel.PluginUpdateInput
}) (*gqlmodel.PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Plugin(ctx context.Context, args struct{ id string }) (*gqlmodel.Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Plugins(ctx context.Context, args struct {
	filter *gqlmodel.PluginFilterInput
	sortBy *gqlmodel.PluginSortingInput
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
