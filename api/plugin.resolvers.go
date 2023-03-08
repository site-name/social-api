package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) PluginUpdate(ctx context.Context, args struct {
	ChannelID *string
	Id        string
	Input     PluginUpdateInput
}) (*PluginUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Plugin(ctx context.Context, args struct{ Id string }) (*Plugin, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Plugins(ctx context.Context, args struct {
	Filter *PluginFilterInput
	SortBy *PluginSortingInput
	GraphqlParams
}) (*PluginCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
