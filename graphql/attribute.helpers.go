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

func (r *mutationResolver) validateValue(valueData gqlmodel.AttributeValueCreateUpdateInterface, isNumericAttr, isSwatchAttr bool) *model.AppError {
	var (
		fileURL     = valueData.GetFileURL()
		contentType = valueData.GetContentType()
		value       = valueData.GetValue()
		name        = valueData.GetName()
		richText    = valueData.GetRichText()
	)

	if !isSwatchAttr &&
		((fileURL != nil && *fileURL != "") ||
			(contentType != nil && *contentType != "") ||
			(value != nil && *value != "")) {

		return model.NewAppError("validateValue", "graphql.attribute.not_swatch_attribute.app_error", nil, "Cannot define value, file and content type fields for not swatch attribute", http.StatusBadRequest)
	}

	if isNumericAttr {
		if _, err := strconv.ParseFloat(*valueData.GetName(), 64); err != nil {
			return model.NewAppError("validateValue", "graphql.attribute.invalid_numeric_value.app_error", nil, "Value of numeric attribute must be numeric", http.StatusBadRequest)
		}
	} else if isSwatchAttr {
		if (value != nil && *value != "") &&
			(fileURL != nil && *fileURL != "") {

			return model.NewAppError("validateValue", "graphql.attribute.redundant_values.app_error", nil, "Cannot specify both value and file for swatch attribute", http.StatusBadRequest)
		}
	}

	var slugValue = *name

	if isNumericAttr {
		slugValue = strings.ReplaceAll(slugValue, ".", "_")
	}

	attributeValue := &attribute.AttributeValue{
		AttributeID: model.NewId(), // fake for validation purpose
		Name:        *name,
		RichText:    richText,
		FileUrl:     fileURL,
		ContentType: contentType,
		Slug:        slug.Make(slugValue),
	}
	if value != nil {
		attributeValue.Value = *value
	}
	if appErr := attributeValue.IsValid(); appErr != nil {
		return appErr
	}

	return nil
}
