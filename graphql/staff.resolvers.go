package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
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

func (r *staffNotificationRecipientResolver) User(ctx context.Context, obj *gqlmodel.StaffNotificationRecipient) (*gqlmodel.User, error) {
	session, appErr := CheckUserAuthenticated("User", ctx)
	if appErr != nil {
		return nil, appErr
	}

	if (obj.UserID != nil && session.UserId == *obj.UserID) || r.Srv().AccountService().SessionHasPermissionTo(session, model.PermissionManageStaff) {
		user, appErr := r.Srv().AccountService().UserById(ctx, *obj.UserID)
		if appErr != nil {
			return nil, appErr
		}
		return gqlmodel.SystemUserToGraphqlUser(user), nil
	}

	return nil, nil
}

func (r *staffNotificationRecipientResolver) Email(ctx context.Context, obj *gqlmodel.StaffNotificationRecipient) (*string, error) {
	if obj.UserID != nil {
		user, appErr := r.Srv().AccountService().UserById(ctx, *obj.UserID)
		if appErr != nil {
			return nil, appErr
		}

		return &user.Email, nil
	}

	return obj.Email(), nil
}

// StaffNotificationRecipient returns generated.StaffNotificationRecipientResolver implementation.
func (r *Resolver) StaffNotificationRecipient() generated.StaffNotificationRecipientResolver {
	return &staffNotificationRecipientResolver{r}
}

type staffNotificationRecipientResolver struct{ *Resolver }
