package giftcard

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppGiftcard struct {
	app.AppIface
}

func init() {
	app.RegisterGiftcardApp(func(a app.AppIface) sub_app_iface.GiftcardApp {
		return &AppGiftcard{a}
	})
}

func (a *AppGiftcard) Save(id string) error {
	return nil
}
