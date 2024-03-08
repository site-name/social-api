package types

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"gorm.io/gorm"
)

// RecalculateOrderPricesFunc
type RecalculateOrderPricesFunc func(tx *gorm.DB, order *model.Order, kwargs map[string]any) *model_helper.AppError
