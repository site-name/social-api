package commands

import (
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	user, _ := a.Srv().Store.User().Get(
		model.UserWhere.Email.EQ(strings.ToLower(userArg)),
		qm.Or(model.UserColumns.Username+" = ?", userArg),
		qm.Or(model.UserColumns.ID+" = ?", userArg),
	)

	return user
}
