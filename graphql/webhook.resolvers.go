package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) WebhookCreate(ctx context.Context, input gqlmodel.WebhookCreateInput) (*gqlmodel.WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookDelete(ctx context.Context, id string) (*gqlmodel.WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookUpdate(ctx context.Context, id string, input gqlmodel.WebhookUpdateInput) (*gqlmodel.WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Webhook(ctx context.Context, id string) (*gqlmodel.Webhook, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookEvents(ctx context.Context) ([]*gqlmodel.WebhookEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookSamplePayload(ctx context.Context, eventType gqlmodel.WebhookSampleEventTypeEnum) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
