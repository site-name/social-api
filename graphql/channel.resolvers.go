package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) ChannelCreate(ctx context.Context, input gqlmodel.ChannelCreateInput) (*gqlmodel.ChannelCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelUpdate(ctx context.Context, id string, input gqlmodel.ChannelUpdateInput) (*gqlmodel.ChannelUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDelete(ctx context.Context, id string, input *gqlmodel.ChannelDeleteInput) (*gqlmodel.ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelActivate(ctx context.Context, id string) (*gqlmodel.ChannelActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChannelDeactivate(ctx context.Context, id string) (*gqlmodel.ChannelDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channel(ctx context.Context, id *string) (*gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Channels(ctx context.Context) ([]*gqlmodel.Channel, error) {
	panic(fmt.Errorf("not implemented"))
}
