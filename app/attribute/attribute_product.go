package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// AttributeProductByOption returns an attribute product filtered using given option
func (a *ServiceAttribute) AttributeProductByOption(option *model.AttributeProductFilterOption) (*model.AttributeProduct, *model_helper.AppError) {
	attributeProduct, err := a.srv.Store.AttributeProduct().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("AttributeProductByOption", "app.attribute.attribute_product_by_options.app_error", nil, err.Error(), statusCode)
	}

	return attributeProduct, nil
}

func (s *ServiceAttribute) AttributeProductsByOption(option *model.AttributeProductFilterOption) ([]*model.AttributeProduct, *model_helper.AppError) {
	attributeProducts, err := s.srv.Store.AttributeProduct().FilterByOptions(option)
	if err != nil {
		return nil, model_helper.NewAppError("AttributeProductsByOption", "app.attribute.attribute_products_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return attributeProducts, nil
}

func (s *ServiceAttribute) AssignedProductAttributesByOption(options *model.AssignedProductAttributeFilterOption) (model.AssignedProductAttributes, *model_helper.AppError) {
	assignedProductAttributes, err := s.srv.Store.AssignedProductAttribute().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("AssignedProductAttributesByOption", "app.attribute.assigned_product_attributes_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return assignedProductAttributes, nil
}
