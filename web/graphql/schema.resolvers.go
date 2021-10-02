package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	graphql1 "github.com/sitename/sitename/web/graphql/generated"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*gqlmodel.VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*gqlmodel.VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignNavigation(ctx context.Context, menu *string, navigationType gqlmodel.NavigationType) (*gqlmodel.AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*gqlmodel.FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExternalNotificationTrigger(ctx context.Context, channel string, input gqlmodel.ExternalNotificationTriggerInput, pluginID *string) (*gqlmodel.ExternalNotificationTrigger, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ReportProductSales(ctx context.Context, period gqlmodel.ReportingPeriod, channel string, before *string, after *string, first *int, last *int) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) HomepageEvents(ctx context.Context, before *string, after *string, first *int, last *int) (*gqlmodel.OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TaxTypes(ctx context.Context) ([]*gqlmodel.TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns graphql1.MutationResolver implementation.
func (r *Resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

// Query returns graphql1.QueryResolver implementation.
func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
