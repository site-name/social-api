package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) StaffNotificationRecipientCreate(ctx context.Context, args struct {
	Input gqlmodel.StaffNotificationRecipientInput
}) (*gqlmodel.StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffNotificationRecipientUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.StaffNotificationRecipientInput
}) (*gqlmodel.StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffNotificationRecipientDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffCreate(ctx context.Context, args struct{ Input gqlmodel.StaffCreateInput }) (*gqlmodel.StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.StaffUpdateInput
}) (*gqlmodel.StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffUsers(ctx context.Context, args struct {
	Filter *gqlmodel.StaffUserInput
	SortBy *gqlmodel.UserSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
