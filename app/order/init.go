package order

import (
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
)

type AppOrder struct {
	app.AppIface
	wg    sync.WaitGroup
	mutex sync.Mutex

	RecalculateOrderPrices RecalculateOrderPricesFunc // decorated function, initialized along with this app
}

func init() {
	app.RegisterOrderApp(func(a app.AppIface) sub_app_iface.OrderApp {
		orderApp := &AppOrder{
			AppIface: a,
		}

		orderApp.RecalculateOrderPrices = orderApp.UpdateVoucherDiscount(func() *model.AppError {

		})

		return orderApp
	})
}

type RecalculateOrderPricesFunc func(*order.Order, map[string]interface{}) *model.AppError

func (a *AppOrder) UpdateVoucherDiscount(fun RecalculateOrderPricesFunc) RecalculateOrderPricesFunc {

	return func(ord *order.Order, kwargs map[string]interface{}) *model.AppError {

		return fun(ord, kwargs)
	}
}
