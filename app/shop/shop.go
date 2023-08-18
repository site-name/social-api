/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package shop

import (
	"github.com/sitename/sitename/app"
)

type ServiceShop struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Shop = &ServiceShop{srv: s}
		return nil
	})
}

// ShopById finds shop by given id
// func (a *ServiceShop) ShopById(shopID string) (model.ShopSettings, *model.AppError) {
// 	return a.ShopByOptions(&model.ShopFilterOptions{
// 		Id: squirrel.Eq{model.ShopTableName + ".Id": shopID},
// 	})
// }

// func (a *ServiceShop) ShopByOptions(options *model.ShopFilterOptions) (model.ShopSettings, *model.AppError) {
// 	shop, err := a.srv.Store.Shop().GetByOptions(options)
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		if _, ok := err.(*store.ErrNotFound); ok {
// 			statusCode = http.StatusNotFound
// 		}

// 		return nil, model.NewAppError("ShopByOptions", "app.shop.error_finding_shop_by_options.app_error", nil, err.Error(), statusCode)
// 	}

// 	return shop, nil
// }
