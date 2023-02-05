package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// AttributeProductByOption returns an attribute product filtered using given option
func (a *ServiceAttribute) AttributeProductByOption(option *model.AttributeProductFilterOption) (*model.AttributeProduct, *model.AppError) {
	attributeProduct, err := a.srv.Store.AttributeProduct().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AttributeProductByOption", "app.attribute.attribute_product_by_options.app_error", nil, err.Error(), statusCode)
	}

	return attributeProduct, nil
}

func (s *ServiceAttribute) AttributeProductsByOption(option *model.AttributeProductFilterOption) ([]*model.AttributeProduct, *model.AppError) {
	attributeProducts, err := s.srv.Store.AttributeProduct().FilterByOptions(option)
	if err != nil {
		return nil, model.NewAppError("AttributeProductsByOption", "app.attribute.attribute_products_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return attributeProducts, nil
}
