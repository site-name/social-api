package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func (r *mutationResolver) VariantMediaAssign(ctx context.Context, mediaID string, variantID string) (*VariantMediaAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VariantMediaUnassign(ctx context.Context, mediaID string, variantID string) (*VariantMediaUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignNavigation(ctx context.Context, menu *string, navigationType NavigationType) (*AssignNavigation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*FileUpload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ReportProductSales(ctx context.Context, period ReportingPeriod, channel string, before *string, after *string, first *int, last *int) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) HomepageEvents(ctx context.Context, before *string, after *string, first *int, last *int) (*OrderEventCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TaxTypes(ctx context.Context) ([]*TaxType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
