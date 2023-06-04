package commands

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
)

func getUsersFromUserArgs(a *app.App, userArgs []string) []*model.User {
	users := make([]*model.User, 0, len(userArgs))
	for _, userArg := range userArgs {
		user := getUserFromUserArg(a, userArg)
		users = append(users, user)
	}
	return users
}

func getUserFromUserArg(a *app.App, userArg string) *model.User {
	user, _ := a.Srv().Store.User().GetByOptions(context.Background(), &model.UserFilterOptions{
		Extra: squirrel.Expr("Users.Email = lower(?) OR Users.Username = ? OR Users.Id = ?", userArg, userArg, userArg),
	})

	return user
}
