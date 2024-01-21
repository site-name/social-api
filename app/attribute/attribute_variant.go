package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

// AttributeVariantByOption returns an attribute variant filtered using given option
func (a *ServiceAttribute) AttributeVariantByOption(option *model.AttributeVariantFilterOption) (*model.AttributeVariant, *model_helper.AppError) {
	attributeVariant, err := a.srv.Store.AttributeVariant().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("AttributeVariantByOption", "app.attribute.attribute_variant_by_options.app_error", nil, err.Error(), statusCode)
	}

	return attributeVariant, nil
}
