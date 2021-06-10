package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ChannelCreate(ctx context.Context, input ChannelCreateInput) (*ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelUpdate(ctx context.Context, id string, input ChannelUpdateInput) (*ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDelete(ctx context.Context, id string, input *ChannelDeleteInput) (*ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelActivate(ctx context.Context, id string) (*ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDeactivate(ctx context.Context, id string) (*ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channel(ctx context.Context, id *string) (*Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channels(ctx context.Context) ([]Channel, error) {
	panic(fmt.Errorf("not implemented"))
}
