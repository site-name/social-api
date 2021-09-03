package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AttributeProductByOption returns an attribute product filtered using given option
func (a *ServiceAttribute) AttributeProductByOption(option *attribute.AttributeProductFilterOption) (*attribute.AttributeProduct, *model.AppError) {
	attributeProduct, err := a.srv.Store.AttributeProduct().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeProductByOption", "app.attribute.error_finding_attribute_product_by_option.app_error", err)
	}

	return attributeProduct, nil
}
