package account

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

const (
	userNotOwnAddress = "app.account.user_not_own_address.app_error"
)

type AppAccount struct {
	app.AppIface
}

func init() {
	app.RegisterAccountApp(func(a app.AppIface) sub_app_iface.AccountApp {
		return &AppAccount{a}
	})
}
