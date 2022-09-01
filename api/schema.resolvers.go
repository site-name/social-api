package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*gqlmodel.VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*gqlmodel.VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AssignNavigation(ctx context.Context, menu *string, navigationType gqlmodel.NavigationType) (*gqlmodel.AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) FileUpload(ctx context.Context, file graphql.Upload) (*gqlmodel.FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExternalNotificationTrigger(ctx context.Context, channel string, input gqlmodel.ExternalNotificationTriggerInput, pluginID *string) (*gqlmodel.ExternalNotificationTrigger, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ReportProductSales(ctx context.Context, period gqlmodel.ReportingPeriod, channel string, before *string, after *string, first *int, last *int) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) HomepageEvents(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) TaxTypes(ctx context.Context) ([]*gqlmodel.TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}
