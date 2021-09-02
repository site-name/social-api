/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package shop

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

type ServiceShop struct {
	srv *app.Server
}

func init() {
	app.RegisterShopService(func(s *app.Server) (sub_app_iface.ShopService, error) {
		return &ServiceShop{
			srv: s,
		}, nil
	})
}

// ShopById finds shop by given id
func (a *ServiceShop) ShopById(shopID string) (*shop.Shop, *model.AppError) {
	shop, err := a.srv.Store.Shop().Get(shopID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShopById", "app.shop.error_finding_shop_by_id.app_error", err)
	}

	return shop, nil
}
