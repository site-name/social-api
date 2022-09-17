package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ProductTranslationsByOption returns a list of product translations
func (s *ServiceProduct) ProductTranslationsByOption(option *model.ProductTranslationFilterOption) ([]*model.ProductTranslation, *model.AppError) {
	translations, err := s.srv.Store.ProductTranslation().FilterByOption(option)
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
		return nil, model.NewAppError("ProductTranslationsByOption", "app.product.error_finding_product_translations_by_option.app_error", nil, errMessage, statusCode)
	}

	return translations, nil
}
