package account

import (
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
)

type AppAccount struct {
	app.AppIface
	sessionPool sync.Pool
}

func init() {
	app.RegisterAccountApp(func(a app.AppIface) sub_app_iface.AccountApp {
		return &AppAccount{
			AppIface: a,
			sessionPool: sync.Pool{
				New: func() interface{} {
					return &model.Session{}
				},
			},
		}
	})
}
