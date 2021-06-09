package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) WebhookCreate(ctx context.Context, input WebhookCreateInput) (*WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookDelete(ctx context.Context, id string) (*WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) WebhookUpdate(ctx context.Context, id string, input WebhookUpdateInput) (*WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Webhook(ctx context.Context, id string) (*Webhook, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookEvents(ctx context.Context) ([]*WebhookEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) WebhookSamplePayload(ctx context.Context, eventType WebhookSampleEventTypeEnum) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
