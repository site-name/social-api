package checkout

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type AppCheckout struct {
	app.AppIface
}

func init() {
	app.RegisterCheckoutApp(func(a app.AppIface) sub_app_iface.CheckoutApp {
		return &AppCheckout{a}
	})
}

func (a *AppCheckout) CheckoutbyId(id string) (*checkout.Checkout, *model.AppError) {
	checkout, err := a.Srv().Store.Checkout().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutById", "app.checkout.missing_checkout.app_error", err)
	}

	return checkout, nil
}
