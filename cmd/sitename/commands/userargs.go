package commands

import (
	"context"

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
	user, _ := a.Srv().Store.User().GetByEmail(userArg)

	if user == nil {
		var err error
		if user, err = a.Srv().Store.User().GetByUsername(userArg); err == nil {
			return user
		}
	}

	if user == nil {
		user, _ = a.Srv().Store.User().Get(context.Background(), userArg)
	}

	return user
}
