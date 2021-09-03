package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AttributeValuesOfAttribute(attributeID string) ([]*attribute.AttributeValue, *model.AppError) {
	attrValues, err := a.srv.Store.AttributeValue().GetAllByAttributeID(attributeID)
	if err != nil {
		// since err can be wither `store.ErrNotFound` or `system error`, so use this shortcut
		return nil, store.AppErrorFromDatabaseLookupError("AttributeValuesOfAttribute", "app.attribute.attribute_values_of_attribute_lookup_error.app_error", err)
	}

	return attrValues, nil
}
