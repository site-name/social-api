package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ChannelCreate(ctx context.Context, args struct{ Input gqlmodel.ChannelCreateInput }) (*gqlmodel.ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ChannelUpdateInput
}) (*gqlmodel.ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelDelete(ctx context.Context, args struct {
	Id    string
	Input *gqlmodel.ChannelDeleteInput
}) (*gqlmodel.ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelActivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelDeactivate(ctx context.Context, args struct{ Id string }) (*gqlmodel.ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Channel(ctx context.Context, args struct{ Id *string }) (*gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Channels(ctx context.Context) ([]gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}
