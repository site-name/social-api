package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AttributeVariantByOption returns an attribute variant filtered using given option
func (a *ServiceAttribute) AttributeVariantByOption(option *attribute.AttributeVariantFilterOption) (*attribute.AttributeVariant, *model.AppError) {
	attributeVariant, err := a.srv.Store.AttributeVariant().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeVariantByOption", "app.attribute.error_finding_attribute_variant_by_option.app_error", err)
	}

	return attributeVariant, nil
}
