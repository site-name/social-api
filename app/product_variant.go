package app

import (
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// ProductVariantById get a product variant with given id if exist
func (a *App) ProductVariantById(id string) (*product_and_discount.ProductVariant, *model.AppError) {
	variant, err := a.srv.Store.ProductVariant().Get(id)
	if err != nil {
		var nfErr *store.ErrNotFound
		statusCode := http.StatusInternalServerError
		if errors.As(err, &nfErr) {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ProductVariantById", "app.product_variant.missing_variant.app_error", nil, nfErr.Error(), statusCode)
	}

	return variant, nil
}
