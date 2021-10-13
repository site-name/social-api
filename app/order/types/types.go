package types

import (
	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(*gorp.Transaction, *order.Order, map[string]interface{}) *model.AppError
