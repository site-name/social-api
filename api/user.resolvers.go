package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) Login(ctx context.Context, args struct{ Input LoginInput }) (*LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarUpdate(ctx context.Context, args struct{ Image Upload }) (*UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarDelete(ctx context.Context) (*UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/user.graphql for details on directives used.
func (r *Resolver) UserBulkSetActive(ctx context.Context, args struct {
	Ids      []string
	IsActive bool
}) (*UserBulkSetActive, error) {
	// validate given ids are valid uuids
	// if !lo.EveryBy(args.Ids, model_helper.IsValidId) {
	// 	return nil, model_helper.NewAppError("UserBulkSetActive", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
	// }

	// embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// currentSession := embedCtx.AppContext.Session()
	// requesterIsSystemUserManager := currentSession.GetUserRoles().Contains(model.SystemUserManagerRoleId)

	// usersToUpdate, appErr := embedCtx.App.Srv().AccountService().GetUsersByIds(args.Ids, &store.UserGetByIdsOpts{})
	// if appErr != nil {
	// 	return nil, appErr
	// }
	// usersToUpdateMap := lo.SliceToMap(usersToUpdate, func(u *model.User) (string, *model.User) { return u.Id, u })

	// // validate if all given ids are valid,
	// // NOTE: the rules below are applied:
	// // 1) system admin can update every users EXCEPT himself.
	// // 2) system user manager can update every users EXCEPT system admin, himself
	// for userId, user := range usersToUpdateMap {
	// 	if currentSession.UserId == userId {
	// 		return nil, model_helper.NewAppError("UserBulkSetActive", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "ids"}, "you can not update yourself", http.StatusForbidden)
	// 	}
	// 	if requesterIsSystemUserManager && user.GetRoles().Contains(model.SystemAdminRoleId) {
	// 		return nil, model_helper.NewAppError("UserBulkSetActive", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "ids"}, "you can not update system admin", http.StatusForbidden)
	// 	}

	// 	// update user
	// 	user.IsActive = args.IsActive
	// }

	// // update
	// embedCtx.App.Srv().AccountService().Up

	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/user.graphql for details on directives used.
func (r *Resolver) Me(ctx context.Context) (*User, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	user, err := UserByUserIdLoader.Load(ctx, embedCtx.AppContext.Session().UserID)()
	if err != nil {
		return nil, err
	}

	return SystemUserToGraphqlUser(user), nil
}

func (r *Resolver) User(ctx context.Context, args struct {
	Id    *string
	Email *string
}) (*User, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if (args.Id == nil && args.Email == nil) || (args.Id != nil && args.Email != nil) {
		return nil, model_helper.NewAppError("User", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "id, email"}, "please provide either email or id, not both", http.StatusBadRequest)
	}
	if args.Id != nil && !model_helper.IsValidId(*args.Id) {
		return nil, model_helper.NewAppError("User", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "id"}, "please provide a valid id", http.StatusBadRequest)
	}
	if args.Email != nil && !model_helper.IsValidEmail(*args.Email) {
		return nil, model_helper.NewAppError("User", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "email"}, "please provide a valid email", http.StatusBadRequest)
	}

	var user *model.User
	var appErr *model_helper.AppError

	if args.Id != nil {
		user, appErr = embedCtx.App.Srv().AccountService().UserById(ctx, *args.Id)
	} else {
		user, appErr = embedCtx.App.Srv().AccountService().GetUserByOptions(model_helper.UserFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.UserWhere.Email.EQ(*args.Email),
			),
		})
	}
	if appErr != nil {
		return nil, appErr
	}

	return SystemUserToGraphqlUser(user), nil
}
