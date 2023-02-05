package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// CategoriesByOption returns all categories that satisfy given option
func (a *ServiceProduct) CategoriesByOption(option *model.CategoryFilterOption) ([]*model.Category, *model.AppError) {
	categories, err := a.srv.Store.Category().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("CategoriesByOption", "app.product.error_finding_categories_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return categories, nil
}

// CategoryByOption returns 1 category that satisfies given option
func (a *ServiceProduct) CategoryByOption(option *model.CategoryFilterOption) (*model.Category, *model.AppError) {
	category, err := a.srv.Store.Category().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("CategoryByOption", "app.product.error_finding_category_by_option.app_error", nil, err.Error(), statusCode)
	}

	return category, nil
}
