package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (r *Resolver) VariantMediaAssign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VariantMediaUnassign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignNavigation(ctx context.Context, args struct {
	Menu           *string
	NavigationType NavigationType
}) (*AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) FileUpload(ctx context.Context, args struct{ File graphql.Upload }) (*FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalNotificationTrigger(ctx context.Context, args struct {
	Channel  string
	Input    ExternalNotificationTriggerInput
	PluginID *string
}) (*ExternalNotificationTrigger, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ReportProductSales(ctx context.Context, args struct {
	Period  ReportingPeriod
	Channel string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) HomepageEvents(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TaxTypes(ctx context.Context) ([]*TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}
