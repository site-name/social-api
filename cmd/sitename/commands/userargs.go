package commands

import (
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

func getUsersFromUserArgs(a *app.App, userArgs []string) []*model.User {
	users := make([]*model.User, 0, len(userArgs))
	for _, userArg := range userArgs {
		user := getUserFromUserArg(a, userArg)
		if user != nil {
			users = append(users, user)
		}
	}
	return users
}

func getUserFromUserArg(a *app.App, userArg string) *model.User {
	user, err := model.Users(
		model_helper.Or{
			squirrel.Eq{model.UserColumns.Email: strings.ToLower(userArg)},
			squirrel.Eq{model.UserColumns.Username: userArg},
			squirrel.Eq{model.UserColumns.ID: userArg},
		},
	).One(a.Srv().Store.GetReplica())
	if err != nil {
		return nil
	}

	return user
}
