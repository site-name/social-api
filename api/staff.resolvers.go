package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) StaffNotificationRecipientCreate(ctx context.Context, args struct {
	Input StaffNotificationRecipientInput
}) (*StaffNotificationRecipientCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffNotificationRecipientUpdate(ctx context.Context, args struct {
	Id    string
	Input StaffNotificationRecipientInput
}) (*StaffNotificationRecipientUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffNotificationRecipientDelete(ctx context.Context, args struct{ Id string }) (*StaffNotificationRecipientDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffCreate(ctx context.Context, args struct{ Input StaffCreateInput }) (*StaffCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffUpdate(ctx context.Context, args struct {
	Id    string
	Input StaffUpdateInput
}) (*StaffUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffDelete(ctx context.Context, args struct{ Id string }) (*StaffDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffBulkDelete(ctx context.Context, args struct{ Ids []string }) (*StaffBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) StaffUsers(ctx context.Context, args struct {
	Filter *StaffUserInput
	SortBy *UserSortingInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*UserCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
