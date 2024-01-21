package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

// ProductVariantTranslationsByOption returns a list of product variant translations
func (s *ServiceProduct) ProductVariantTranslationsByOption(option *model.ProductVariantTranslationFilterOption) ([]*model.ProductVariantTranslation, *model_helper.AppError) {
	translations, err := s.srv.Store.ProductVariantTranslation().FilterByOption(option)
	var (
		errMessage string
		statusCode int
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(translations) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model_helper.NewAppError("ProductVariantTranslationsByOption", "app.product.error_finding_product_variant_translations_by_option.app_error", nil, errMessage, statusCode)
	}

	return translations, nil
}
