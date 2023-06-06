package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentCreate(ctx context.Context, args struct {
	Input     DigitalContentUploadInput
	VariantID string
}) (*DigitalContentCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentDelete(ctx context.Context, args struct{ VariantID string }) (*DigitalContentDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentUpdate(ctx context.Context, args struct {
	Input     DigitalContentInput
	VariantID string
}) (*DigitalContentUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentURLCreate(ctx context.Context, args struct {
	Input DigitalContentURLCreateInput
}) (*DigitalContentURLCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContent(ctx context.Context, args struct{ Id string }) (*DigitalContent, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) DigitalContents(ctx context.Context, args GraphqlParams) (*DigitalContentCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
