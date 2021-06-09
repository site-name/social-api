package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) StaffNotificationRecipientCreate(ctx context.Context, input StaffNotificationRecipientInput) (*StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientUpdate(ctx context.Context, id string, input StaffNotificationRecipientInput) (*StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientDelete(ctx context.Context, id string) (*StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffCreate(ctx context.Context, input StaffCreateInput) (*StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffUpdate(ctx context.Context, id string, input StaffUpdateInput) (*StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffDelete(ctx context.Context, id string) (*StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffBulkDelete(ctx context.Context, ids []*string) (*StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) StaffUsers(ctx context.Context, filter *StaffUserInput, sortBy *UserSortingInput, before *string, after *string, first *int, last *int) (*UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
