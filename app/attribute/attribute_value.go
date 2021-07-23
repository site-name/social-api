package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

func (a *AppAttribute) AttributeValuesOfAttribute(attributeID string) ([]*attribute.AttributeValue, *model.AppError) {
	attrValues, appErr := a.app.Srv().Store.AttributeValue().GetAllByAttributeID(attributeID)
	if appErr != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeValuesOfAttribute", "app.attribute.attribute_values_of_attribute_lookup_error.app_error", appErr)
	}

	return attrValues, nil
}
