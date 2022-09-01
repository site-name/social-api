package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) VariantMediaAssign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*gqlmodel.VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VariantMediaUnassign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*gqlmodel.VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignNavigation(ctx context.Context, args struct {
	Menu           *string
	NavigationType gqlmodel.NavigationType
}) (*gqlmodel.AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) FileUpload(ctx context.Context, args struct{ File graphql.Upload }) (*gqlmodel.FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalNotificationTrigger(ctx context.Context, args struct {
	Channel  string
	Input    gqlmodel.ExternalNotificationTriggerInput
	PluginID *string
}) (*gqlmodel.ExternalNotificationTrigger, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ReportProductSales(ctx context.Context, args struct {
	Period  gqlmodel.ReportingPeriod
	Channel string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) HomepageEvents(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TaxTypes(ctx context.Context) ([]*gqlmodel.TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}
