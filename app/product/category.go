package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// CategoriesByOption returns all categories that satisfy given option
func (a *ServiceProduct) CategoriesByOption(option *product_and_discount.CategoryFilterOption) ([]*product_and_discount.Category, *model.AppError) {
	categories, err := a.srv.Store.Category().FilterByOption(option)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	}
	if len(categories) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("CategoriesByOption", "app.product.error_finding_categories_by_option.app_error", nil, errMsg, statusCode)
	}

	return categories, nil
}

// CategoryByOption returns 1 category that satisfies given option
func (a *ServiceProduct) CategoryByOption(option *product_and_discount.CategoryFilterOption) (*product_and_discount.Category, *model.AppError) {
	category, err := a.srv.Store.Category().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CategoryByOption", "app.product.error_finding_category_by_option.app_error", err)
	}

	return category, nil
}
