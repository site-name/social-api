package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type AttributeUpsertInputIface interface {
	getFieldValueByString(name string) any
	// getInputType. NOTE: make sure to check result by call its .IsValid() method
	getInputType() AttributeInputTypeEnum
}

type attributeValueInputIface interface {
	getName() string
	getFileURL() *string
	getContentType() *string
	getValue() *string
	getJsonString() JSONString
}

type AttributeMixin[T AttributeUpsertInputIface] struct {
	ATTRIBUTE_VALUES_FIELD string // must be either "values" or "add_values"
	srv                    *app.Server
}

func (a *AttributeMixin[T]) cleanValues(cleanedInput AttributeUpsertInputIface, attribute *model.Attribute) (model.AttributeValues, *model.AppError) {
	value := cleanedInput.getFieldValueByString(a.ATTRIBUTE_VALUES_FIELD)
	valuesInput := []attributeValueInputIface{}
	if value != nil {
		switch t := value.(type) {
		case []*AttributeValueUpdateInput:
			for _, item := range t {
				valuesInput = append(valuesInput, item)
			}
		case []*AttributeValueCreateInput:
			for _, item := range t {
				valuesInput = append(valuesInput, item)
			}
		}
	}

	attributeInputType := cleanedInput.getFieldValueByString("input_type")
	if attributeInputType == nil {
		attributeInputType = attribute.InputType
	}
	if len(valuesInput) == 0 {
		return nil, nil
	}

	if attributeInputType != nil &&
		(attributeInputType == model.FILE_ ||
			attributeInputType == model.REFERENCE) {
		return nil, model.NewAppError("cleanValues", "api.attribute.invalid_input_type_values.app_error", nil, fmt.Sprintf("Values cannot be used with input type %s", attributeInputType), http.StatusBadRequest)
	}

	validatedAttributeValues := model.AttributeValues{}
	for _, valueData := range valuesInput {
		value, appErr := validateValue(attribute, valueData, attributeInputType == model.NUMERIC, attributeInputType == model.SWATCH)
		if appErr != nil {
			return nil, appErr
		}

		validatedAttributeValues = append(validatedAttributeValues, value)
	}

	return validatedAttributeValues, nil
}

func (a *AttributeMixin[T]) checkValuesAreUnique(valuesInput []attributeValueInputIface, attribute *model.Attribute) *model.AppError {
	attributeValues, appErr := a.srv.AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		AttributeID: squirrel.Eq{store.AttributeValueTableName + ".AttributeID": attribute.Id},
	})
	if appErr != nil {
		return appErr
	}
	oldAttributeValueSlugs := lo.SliceToMap(attributeValues, func(v *model.AttributeValue) (string, struct{}) { return v.Slug, struct{}{} })
	uniqueSlugMap := map[string]struct{}{}

	for _, valueData := range valuesInput {
		slug := slug.Make(valueData.getName())

		_, exist := oldAttributeValueSlugs[slug]
		if exist {
			return model.NewAppError("checkValuesAreUnique", "api.attribute.value_exist_within_attribute.app_error", nil, fmt.Sprintf("attribute value %s already exist within attribute", valueData.getName()), http.StatusBadRequest)
		}

		_, found := uniqueSlugMap[slug]
		if !found {
			uniqueSlugMap[slug] = struct{}{}
		}
	}

	if len(uniqueSlugMap) != len(valuesInput) {
		return model.NewAppError("checkValuesAreUnique", "api.attribute.values_are_not_unique.app_error", nil, "provided values are not unique", http.StatusBadRequest)
	}
	return nil
}

func validateValue(attribute *model.Attribute, valueData attributeValueInputIface, isNumericAttr, isSwatchAttr bool) (*model.AttributeValue, *model.AppError) {
	name := valueData.getName()
	appErr := cleanValueInputData(valueData, isSwatchAttr)
	if appErr != nil {
		return nil, appErr
	}

	if isNumericAttr {
		appErr = validateNumericValue(name)
	} else if isSwatchAttr {
		appErr = validateSwatchAttributeValue(valueData)
	}
	if appErr != nil {
		return nil, appErr
	}

	var slugValue string
	if !isNumericAttr {
		slugValue = name
	} else {
		slugValue = strings.ReplaceAll(name, ".", "_")
	}
	slugValue = slug.Make(slugValue)

	attributeValue := &model.AttributeValue{
		Name:        name,
		Slug:        slugValue,
		FileUrl:     valueData.getFileURL(),
		ContentType: valueData.getContentType(),
		RichText:    model.StringInterface(valueData.getJsonString()),
		AttributeID: attribute.Id,
	}
	value := valueData.getValue()
	if value != nil {
		attributeValue.Value = *value
	}

	attributeValue.PreSave() // call this to make uuid
	appErr = attributeValue.IsValid()
	if appErr != nil {
		return nil, appErr
	}
	return attributeValue, nil
}

func validateSwatchAttributeValue(valueData attributeValueInputIface) *model.AppError {
	var (
		value   = valueData.getValue()
		fileUrl = valueData.getFileURL()
	)

	if (value != nil && *value != "") &&
		(fileUrl != nil && *fileUrl != "") {
		return model.NewAppError("validateSwatchAttributeValue", "api.attribute.swatch_attribute_cannot_has_file_and_value.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func validateNumericValue(value string) *model.AppError {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return model.NewAppError("validateNumericValue", "api.attribute.numeric_value_invalid.app_error", nil, "value of numeric attribute must be float number", http.StatusBadRequest)
	}
	return nil
}

func cleanValueInputData(valueData attributeValueInputIface, isSwatchAttr bool) *model.AppError {
	var (
		contentType = valueData.getContentType()
		fileUrl     = valueData.getFileURL()
		value       = valueData.getValue()
	)
	if !isSwatchAttr &&
		((contentType != nil && *contentType != "") ||
			(fileUrl != nil && *fileUrl != "") ||
			(value != nil && *value != "")) {

		return model.NewAppError("cleanValueInputData", "api.attribute.not_swatch_attribute_has_no_file_nor_value_nor_contentType.app_error", nil, "cannot define value, file and content type fields for non-swatch attribute", http.StatusBadRequest)
	}

	return nil
}

func cleanAttributeSettings(instance *model.Attribute, cleanedInput AttributeUpsertInputIface) *model.AppError {
	inputType := cleanedInput.getInputType()
	if !inputType.IsValid() {
		inputType = instance.InputType
	}

	for field, allowedInputTypes := range model.ATTRIBUTE_PROPERTIES_CONFIGURATION {
		if !lo.Contains(allowedInputTypes, inputType) &&
			cleanedInput.getFieldValueByString(field) != nil {
			return model.NewAppError("cleanAttributeSettings", "api.attribute.field_cannot_be_set.app_error", nil, fmt.Sprintf("cannot set %s on a %s attribute", field, inputType), http.StatusBadRequest)
		}
	}

	return nil
}
