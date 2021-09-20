package checkout

import (
	"net/http"

	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
)

// PrepareInsufficientStockCheckoutValidationAppError
func PrepareInsufficientStockCheckoutValidationAppError(where string, err *exception.InsufficientStock) *model.AppError {
	return model.NewAppError(where, "app.checkout.insufficient_stock.app_error", map[string]interface{}{"variants": err.VariantIDs()}, "", http.StatusNotAcceptable)
}
