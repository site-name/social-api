package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ChannelCreate(ctx context.Context, args struct{ Input ChannelCreateInput }) (*ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelUpdate(ctx context.Context, args struct {
	Id    string
	Input ChannelUpdateInput
}) (*ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelDelete(ctx context.Context, args struct {
	Id    string
	Input *ChannelDeleteInput
}) (*ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelActivate(ctx context.Context, args struct{ Id string }) (*ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelDeactivate(ctx context.Context, args struct{ Id string }) (*ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Channel(ctx context.Context, args struct{ Id *string }) (*Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Channels(ctx context.Context) ([]Channel, error) {
	panic(fmt.Errorf("not implemented"))
}
