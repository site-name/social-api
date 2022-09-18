package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) WebhookCreate(ctx context.Context, args struct{ Input WebhookCreateInput }) (*WebhookCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookDelete(ctx context.Context, args struct{ Id string }) (*WebhookDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookUpdate(ctx context.Context, args struct {
	Id    string
	Input WebhookUpdateInput
}) (*WebhookUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Webhook(ctx context.Context, args struct{ Id string }) (*Webhook, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookEvents(ctx context.Context) ([]*WebhookEvent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) WebhookSamplePayload(ctx context.Context, args struct {
	EventType WebhookSampleEventTypeEnum
}) (JSONString, error) {
	panic(fmt.Errorf("not implemented"))
}
