package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) StaffNotificationRecipientCreate(ctx context.Context, input gqlmodel.StaffNotificationRecipientInput) (*gqlmodel.StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientUpdate(ctx context.Context, id string, input gqlmodel.StaffNotificationRecipientInput) (*gqlmodel.StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffNotificationRecipientDelete(ctx context.Context, id string) (*gqlmodel.StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffCreate(ctx context.Context, input gqlmodel.StaffCreateInput) (*gqlmodel.StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffUpdate(ctx context.Context, id string, input gqlmodel.StaffUpdateInput) (*gqlmodel.StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffDelete(ctx context.Context, id string) (*gqlmodel.StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StaffBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) StaffUsers(ctx context.Context, filter *gqlmodel.StaffUserInput, sortBy *gqlmodel.UserSortingInput, before *string, after *string, first *int, last *int) (*gqlmodel.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
