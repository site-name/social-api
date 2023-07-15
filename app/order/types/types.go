package types

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store/store_iface"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(store_iface.SqlxExecutor, *model.Order, map[string]interface{}) *model.AppError
