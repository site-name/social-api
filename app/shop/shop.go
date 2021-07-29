package site

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

type AppShop struct {
	app app.AppIface
}

func init() {
	app.RegisterShopApp(func(a app.AppIface) sub_app_iface.ShopApp {
		return &AppShop{a}
	})
}

// ShopById finds shop by given id
func (a *AppShop) ShopById(shopID string) (*shop.Shop, *model.AppError) {
	shop, err := a.app.Srv().Store.Shop().Get(shopID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShopById", "app.shop.shop_by_id.app_error", err)
	}

	return shop, nil
}
