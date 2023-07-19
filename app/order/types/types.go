package types

import (
	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(*gorm.DB, *model.Order, map[string]interface{}) *model.AppError
