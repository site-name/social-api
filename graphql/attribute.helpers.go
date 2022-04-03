package graphql

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
)

func (r *mutationResolver) validateValue(valueData gqlmodel.AttributeValueCreateInput, isNumericAttr, isSwatchAttr bool) *model.AppError {
	if !isSwatchAttr &&
		((valueData.FileURL != nil && *valueData.FileURL != "") ||
			(valueData.ContentType != nil && *valueData.ContentType != "") ||
			(valueData.Value != nil && *valueData.Value != "")) {

		return model.NewAppError("validateValue", "graphql.attribute.not_swatch_attribute.app_error", nil, "Cannot define value, file and content type fields for not swatch attribute", http.StatusBadRequest)
	}

	if isNumericAttr {
		if _, err := strconv.ParseFloat(valueData.Name, 64); err != nil {
			return model.NewAppError("validateValue", "graphql.attribute.invalid_numeric_value.app_error", nil, "Value of numeric attribute must be numeric", http.StatusBadRequest)
		}
	} else if isSwatchAttr {
		if (valueData.Value != nil && *valueData.Value != "") &&
			(valueData.FileURL != nil && *valueData.FileURL != "") {

			return model.NewAppError("validateValue", "graphql.attribute.redundant_values.app_error", nil, "Cannot specify both value and file for swatch attribute", http.StatusBadRequest)
		}
	}

	var slugValue = valueData.Name

	if isNumericAttr {
		slugValue = strings.ReplaceAll(slugValue, ".", "_")
	}

	attributeValue := &attribute.AttributeValue{
		AttributeID: model.NewId(), // fake for validation purpose
		Name:        valueData.Name,
		RichText:    valueData.RichText,
		FileUrl:     valueData.FileURL,
		ContentType: valueData.ContentType,
		Slug:        slug.Make(slugValue),
	}
	if valueData.Value != nil {
		attributeValue.Value = *valueData.Value
	}
	if appErr := attributeValue.IsValid(); appErr != nil {
		return appErr
	}

	return nil
}
