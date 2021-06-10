package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) ExportProducts(ctx context.Context, input gqlmodel.ExportProductsInput) (*gqlmodel.ExportProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*gqlmodel.ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *gqlmodel.ExportFileFilterInput, sortBy *gqlmodel.ExportFileSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
