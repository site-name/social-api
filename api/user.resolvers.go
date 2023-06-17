package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) Login(ctx context.Context, args struct{ Input LoginInput }) (*LoginResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarUpdate(ctx context.Context, args struct{ Image graphql.Upload }) (*UserAvatarUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserAvatarDelete(ctx context.Context) (*UserAvatarDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) UserBulkSetActive(ctx context.Context, args struct {
	Ids      []string
	IsActive bool
}) (*UserBulkSetActive, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Me(ctx context.Context) (*User, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	return SystemUserToGraphqlUser(user), nil
}

func (r *Resolver) User(ctx context.Context, args struct {
	Id    *string
	Email *string
}) (*User, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if args.Id == nil && args.Email == nil {
		embedCtx.SetInvalidUrlParam("id, email")
		return nil, embedCtx.Err
	}
	if args.Id != nil && !model.IsValidId(*args.Id) {
		embedCtx.SetInvalidUrlParam("args.Id")
		return nil, embedCtx.Err
	}
	if args.Email != nil && !model.IsValidEmail(*args.Email) {
		embedCtx.SetInvalidUrlParam("args.Email")
		return nil, embedCtx.Err
	}

	var user *model.User
	var appErr *model.AppError
	if args.Id != nil {
		user, appErr = embedCtx.App.Srv().AccountService().UserById(ctx, *args.Id)
	} else {
		user, appErr = embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
			Email: squirrel.Eq{store.UserTableName + ".Email": *args.Email},
		})
	}
	if appErr != nil {
		return nil, appErr
	}

	return SystemUserToGraphqlUser(user), nil
}
