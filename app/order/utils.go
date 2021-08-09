package order

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store"
)

// FilterOrdersByOptions is common method for filtering orders by given option
func (a *AppOrder) FilterOrdersByOptions(option *order.OrderFilterOption) ([]*order.Order, *model.AppError) {
	orders, err := a.Srv().Store.Order().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FilterOrdersbyOption", "app.order.error_finding_orders_by_option.app_error", err)
	}

	return orders, nil
}
