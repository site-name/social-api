package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ExportProducts(ctx context.Context, input ExportProductsInput) (*ExportProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFile(ctx context.Context, id string) (*ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ExportFiles(ctx context.Context, filter *ExportFileFilterInput, sortBy *ExportFileSortingInput, before *string, after *string, first *int, last *int) (*ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
