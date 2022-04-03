package attribute

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AttributeValuesOfAttribute(attributeID string) ([]*attribute.AttributeValue, *model.AppError) {
	attrValues, err := a.srv.Store.AttributeValue().FilterByOptions(attribute.AttributeValueFilterOptions{
		AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attributeID},
	})
	var (
		statusCode = 0
		errMsg     string
	)
	if err != nil {
		errMsg = err.Error()
		statusCode = http.StatusInternalServerError
	}
	if len(attrValues) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AttributeValuesOfAttribute", "app.attribute.error_finding_attribute_values_by_attribute_id.app_error", nil, errMsg, statusCode)
	}

	return attrValues, nil
}

// UpsertAttributeValue insderts or updates given attribute value then returns it
func (a *ServiceAttribute) UpsertAttributeValue(attrValue *attribute.AttributeValue) (*attribute.AttributeValue, *model.AppError) {
	attrValue, err := a.srv.Store.AttributeValue().Upsert(attrValue)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertAttributeValue", "app.attribute.error_upserting_attribute_value.app_error", nil, err.Error(), statusCode)
	}

	return attrValue, nil
}
