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
