package types

import (
	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(tx *gorm.DB, order *model.Order, kwargs map[string]interface{}) *model.AppError
