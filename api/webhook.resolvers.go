package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *Resolver) WebhookCreate(ctx context.Context, args struct{ input gqlmodel.WebhookCreateInput }) (*gqlmodel.WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.WebhookUpdateInput
}) (*gqlmodel.WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Webhook(ctx context.Context, args struct{ id string }) (*gqlmodel.Webhook, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookEvents(ctx context.Context) ([]*gqlmodel.WebhookEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookSamplePayload(ctx context.Context, args struct {
	eventType gqlmodel.WebhookSampleEventTypeEnum
}) (model.StringInterface, error) {
	panic(fmt.Errorf("not implemented"))
}
