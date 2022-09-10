package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ExportProducts(ctx context.Context, args struct{ Input ExportProductsInput }) (*ExportProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExportFile(ctx context.Context, args struct{ Id string }) (*ExportFile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ExportFiles(ctx context.Context, args struct {
	Filter *ExportFileFilterInput
	SortBy *ExportFileSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*ExportFileCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
