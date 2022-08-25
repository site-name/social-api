package types

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/store/store_iface"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(store_iface.SqlxTxExecutor, *order.Order, map[string]interface{}) *model.AppError
