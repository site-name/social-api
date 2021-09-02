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

type ServiceShopConfig struct {
	Server *app.Server
}

func NewServiceShop(config *ServiceShopConfig) sub_app_iface.ShopService {
	return &ServiceShop{
		srv: config.Server,
	}
}

// ShopById finds shop by given id
func (a *ServiceShop) ShopById(shopID string) (*shop.Shop, *model.AppError) {
	shop, err := a.srv.Store.Shop().Get(shopID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShopById", "app.shop.error_finding_shop_by_id.app_error", err)
	}

	return shop, nil
}
