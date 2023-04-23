/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package shop

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
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
func (a *ServiceShop) ShopById(shopID string) (*model.Shop, *model.AppError) {
	return a.ShopByOptions(&model.ShopFilterOptions{
		Id: squirrel.Eq{store.ShopTableName + ".Id": shopID},
	})
}

func (a *ServiceShop) ShopByOptions(options *model.ShopFilterOptions) (*model.Shop, *model.AppError) {
	shop, err := a.srv.Store.Shop().GetByOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model.NewAppError("ShopByOptions", "app.shop.error_finding_shop_by_options.app_error", nil, err.Error(), statusCode)
	}

	return shop, nil
}
