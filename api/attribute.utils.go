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
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

// attributeUpsertInputIface represents AttributeUpdateInput | AttributeCreateInput
type attributeUpsertInputIface interface {
	getFieldValueByString(name string) any
	// getInputType. NOTE: make sure to check result by call its .IsValid() method
	getInputType() model.AttributeInputType
}

type attributeValueUpsertInputIface interface {
	getName() string
	getFileURL() *string
	getContentType() *string
	getValue() *string
	getJsonString() JSONString
}

var (
	_ attributeValueUpsertInputIface = (*AttributeValueCreateInput)(nil)
	_ attributeValueUpsertInputIface = (*AttributeValueUpdateInput)(nil)
	_ attributeUpsertInputIface      = (*AttributeUpdateInput)(nil)
	_ attributeUpsertInputIface      = (*AttributeCreateInput)(nil)
)

type AttributeUpdateInput struct {
	Name                     *string                      `json:"name"`
	Slug                     *string                      `json:"slug"`
	Unit                     *MeasurementUnitsEnum        `json:"unit"`
	RemoveValues             []UUID                       `json:"removeValues"`
	AddValues                []*AttributeValueUpdateInput `json:"addValues"`
	ValueRequired            *bool                        `json:"valueRequired"`
	IsVariantOnly            *bool                        `json:"isVariantOnly"`
	VisibleInStorefront      *bool                        `json:"visibleInStorefront"`
	FilterableInStorefront   *bool                        `json:"filterableInStorefront"`
	FilterableInDashboard    *bool                        `json:"filterableInDashboard"`
	StorefrontSearchPosition *int32                       `json:"storefrontSearchPosition"`
	AvailableInGrid          *bool                        `json:"availableInGrid"`
}

func (a *AttributeUpdateInput) getInputType() model.AttributeInputType {
	var res model.AttributeInputType
	return res
}

// current returns value only when name == "add_values"
func (i *AttributeUpdateInput) getFieldValueByString(name string) any {
	if name == "add_values" {
		return i.AddValues
	}
	return nil
}

type AttributeCreateInput struct {
	InputType                *AttributeInputTypeEnum      `json:"inputType"`
	EntityType               *AttributeEntityTypeEnum     `json:"entityType"`
	Name                     string                       `json:"name"`
	Slug                     *string                      `json:"slug"`
	Type                     AttributeTypeEnum            `json:"type"`
	Unit                     *MeasurementUnitsEnum        `json:"unit"`
	Values                   []*AttributeValueCreateInput `json:"values"`
	ValueRequired            *bool                        `json:"valueRequired"`
	IsVariantOnly            *bool                        `json:"isVariantOnly"`
	VisibleInStorefront      *bool                        `json:"visibleInStorefront"`
	FilterableInStorefront   *bool                        `json:"filterableInStorefront"`
	FilterableInDashboard    *bool                        `json:"filterableInDashboard"`
	StorefrontSearchPosition *int32                       `json:"storefrontSearchPosition"`
	AvailableInGrid          *bool                        `json:"availableInGrid"`
}

func (a *AttributeCreateInput) getInputType() model.AttributeInputType {
	if a.InputType != nil {
		return *a.InputType
	}

	var res model.AttributeInputType
	return res
}

func (a *AttributeCreateInput) getFieldValueByString(field string) any {
	switch field {
	case "filterable_in_storefront":
		return a.FilterableInStorefront
	case "filterable_in_dashboard":
		return a.FilterableInDashboard
	case "available_in_grid":
		return a.AvailableInGrid
	case "storefront_search_position":
		return a.StorefrontSearchPosition
	case "values":
		return a.Values
	case "input_type":
		return a.InputType
	default:
		return nil
	}
}

type AttributeValueCreateInput struct {
	Name        string     `json:"name"`
	Value       *string    `json:"value"`
	RichText    JSONString `json:"richText"`
	FileURL     *string    `json:"fileUrl"`
	ContentType *string    `json:"contentType"`
}

func (a *AttributeValueCreateInput) getName() string           { return a.Name }
func (a *AttributeValueCreateInput) getFileURL() *string       { return a.FileURL }
func (a *AttributeValueCreateInput) getContentType() *string   { return a.ContentType }
func (a *AttributeValueCreateInput) getValue() *string         { return a.Value }
func (a *AttributeValueCreateInput) getJsonString() JSONString { return a.RichText }

type AttributeValueUpdateInput struct {
	AttributeValueCreateInput
}

type AttributeMixin[T attributeUpsertInputIface] struct {
	ATTRIBUTE_VALUES_FIELD string // must be either "values" or "add_values"
	srv                    *app.Server
}

func newAttributeMixin[T attributeUpsertInputIface](srv *app.Server, ATTRIBUTE_VALUES_FIELD string) *AttributeMixin[T] {
	return &AttributeMixin[T]{ATTRIBUTE_VALUES_FIELD, srv}
}

func (a *AttributeMixin[T]) cleanValues(cleanedInput attributeUpsertInputIface, attribute *model.Attribute) (model.AttributeValues, *model_helper.AppError) {
	value := cleanedInput.getFieldValueByString(a.ATTRIBUTE_VALUES_FIELD)
	valuesInput := []attributeValueUpsertInputIface{}
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
		(attributeInputType == model.AttributeInputTypeFile ||
			attributeInputType == model.AttributeInputTypeReference) {
		return nil, model_helper.NewAppError("cleanValues", "api.attribute.invalid_input_type_values.app_error", nil, fmt.Sprintf("Values cannot be used with input type %s", attributeInputType), http.StatusBadRequest)
	}

	validatedAttributeValues := model.AttributeValues{}
	for _, valueData := range valuesInput {
		value, appErr := validateValue(attribute, valueData, attributeInputType == model.AttributeInputTypeNumeric, attributeInputType == model.AttributeInputTypeSwatch)
		if appErr != nil {
			return nil, appErr
		}

		validatedAttributeValues = append(validatedAttributeValues, value)
	}

	return validatedAttributeValues, nil
}

func (a *AttributeMixin[T]) checkValuesAreUnique(valuesInput []attributeValueUpsertInputIface, attribute *model.Attribute) *model_helper.AppError {
	attributeValues, appErr := a.srv.AttributeService().FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
		Conditions: squirrel.Eq{model.AttributeValueTableName + ".AttributeID": attribute.Id},
	})
	if appErr != nil {
		return appErr
	}
	oldAttributeValueSlugs := lo.SliceToMap(attributeValues, func(v *model.AttributeValue) (string, struct{}) { return v.Slug, struct{}{} })
	uniqueSlugMap := map[string]struct{}{}

	for _, valueData := range valuesInput {
		slugg := slug.Make(valueData.getName())

		_, exist := oldAttributeValueSlugs[slugg]
		if exist {
			return model_helper.NewAppError("checkValuesAreUnique", "api.attribute.value_exist_within_attribute.app_error", nil, fmt.Sprintf("attribute value %s already exist within attribute", valueData.getName()), http.StatusBadRequest)
		}

		_, found := uniqueSlugMap[slugg]
		if !found {
			uniqueSlugMap[slugg] = struct{}{}
		}
	}

	if len(uniqueSlugMap) != len(valuesInput) {
		return model_helper.NewAppError("checkValuesAreUnique", "api.attribute.values_are_not_unique.app_error", nil, "provided values are not unique", http.StatusBadRequest)
	}
	return nil
}

func validateValue(attribute *model.Attribute, valueData attributeValueUpsertInputIface, isNumericAttr, isSwatchAttr bool) (*model.AttributeValue, *model_helper.AppError) {
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

func validateSwatchAttributeValue(valueData attributeValueUpsertInputIface) *model_helper.AppError {
	var (
		value   = valueData.getValue()
		fileUrl = valueData.getFileURL()
	)

	if (value != nil && *value != "") &&
		(fileUrl != nil && *fileUrl != "") {
		return model_helper.NewAppError("validateSwatchAttributeValue", "api.attribute.swatch_attribute_cannot_has_file_and_value.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func validateNumericValue(value string) *model_helper.AppError {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return model_helper.NewAppError("validateNumericValue", "api.attribute.numeric_value_invalid.app_error", nil, "value of numeric attribute must be float number", http.StatusBadRequest)
	}
	return nil
}

func cleanValueInputData(valueData attributeValueUpsertInputIface, isSwatchAttr bool) *model_helper.AppError {
	var (
		contentType = valueData.getContentType()
		fileUrl     = valueData.getFileURL()
		value       = valueData.getValue()
	)
	if !isSwatchAttr &&
		((contentType != nil && *contentType != "") ||
			(fileUrl != nil && *fileUrl != "") ||
			(value != nil && *value != "")) {

		return model_helper.NewAppError("cleanValueInputData", "api.attribute.not_swatch_attribute_has_no_file_nor_value_nor_contentType.app_error", nil, "cannot define value, file and content type fields for non-swatch attribute", http.StatusBadRequest)
	}

	return nil
}

func cleanAttributeSettings(instance *model.Attribute, cleanedInput attributeUpsertInputIface) *model_helper.AppError {
	inputType := cleanedInput.getInputType()
	if !inputType.IsValid() {
		inputType = instance.InputType
	}

	for field, allowedInputTypes := range model.ATTRIBUTE_PROPERTIES_CONFIGURATION {
		if !lo.Contains(allowedInputTypes, inputType) &&
			cleanedInput.getFieldValueByString(field) != nil {
			return model_helper.NewAppError("cleanAttributeSettings", "api.attribute.field_cannot_be_set.app_error", nil, fmt.Sprintf("cannot set %s on a %s attribute", field, inputType), http.StatusBadRequest)
		}
	}

	return nil
}

type AttrValuesInput struct {
	GlobalID    string
	Values      []string
	References  []string
	FileUrl     *string
	ContentType *string
	RichText    model.StringInterface
	Boolean     *bool
	Date        *string
	DateTime    *string
}

// type attributeValidator func(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError

func isValueRequired(attribute *model.Attribute, variantValidation bool) bool {
	return attribute.ValueRequired ||
		(variantValidation && model.ALLOWED_IN_VARIANT_SELECTION.Contains(attribute.InputType))
}

func validateFileAttributesInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError {
	if (attributeValues.FileUrl == nil || *attributeValues.FileUrl == "") &&
		isValueRequired(attribute, variantValidation) {
		return model_helper.NewAppError("validateFileAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "FileUrl"}, "please provide file", http.StatusBadRequest)
	}
	return nil
}

func validateReferenceAttributesInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValiation bool) *model_helper.AppError {
	if len(attributeValues.References) == 0 && isValueRequired(attribute, variantValiation) {
		return model_helper.NewAppError("validateReferenceAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "References"}, "please provide references", http.StatusBadRequest)
	}
	return nil
}

func validateBooleanInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError {
	if attribute.ValueRequired && attributeValues.Boolean == nil {
		return model_helper.NewAppError("validateBooleanInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Boolean"}, "please provide boolean value", http.StatusBadRequest)
	}
	return nil
}

func validateRichTextAttributesInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError {
	if (attributeValues.RichText == nil || len(attributeValues.RichText) == 0) && attribute.ValueRequired {
		return model_helper.NewAppError("validateRichTextAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "RichText"}, "please provide rich text", http.StatusBadRequest)
	}
	return nil
}

func validateStandardAttributesInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError {
	if len(attributeValues.Values) == 0 && isValueRequired(attribute, variantValidation) {
		return model_helper.NewAppError("validateStandardAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Values"}, "please provide values", http.StatusBadRequest)
	}
	if attribute.InputType != model.AttributeInputTypeMultiSelect && len(attributeValues.Values) != 1 {
		return model_helper.NewAppError("validateStandardAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Values"}, "only 1 value allowed", http.StatusBadRequest)
	}

	isNumeric := attribute.InputType == model.AttributeInputTypeNumeric

	for _, value := range attributeValues.Values {
		if !isNumeric && strings.TrimSpace(value) == "" {
			return model_helper.NewAppError("validateStandardAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Values"}, "please provide non empty values", http.StatusBadRequest)
		}
		if !isNumeric && len(value) > model.AttributeValueNameMaxLength {
			return model_helper.NewAppError("validateStandardAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Values"}, fmt.Sprintf("some value has length exceeds allowed (%d)", model.AttributeValueNameMaxLength), http.StatusBadRequest)
		}

		if isNumeric {
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return model_helper.NewAppError("validateStandardAttributesInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Values"}, "please provide numeric values. err: "+err.Error(), http.StatusBadRequest)
			}
		}
	}

	return nil
}

func validateDatetimeInput(attribute *model.Attribute, attributeValues AttrValuesInput, variantValidation bool) *model_helper.AppError {
	isBlankDate := attribute.InputType == model.AttributeInputTypeDate && attributeValues.Date == nil
	isBlankDatetime := attribute.InputType == model.AttributeInputTypeDateTime && attributeValues.DateTime == nil

	if attribute.ValueRequired && (isBlankDate || isBlankDatetime) {
		return model_helper.NewAppError("validateDatetimeInput", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "date/datetime"}, "please provide date or datetime value", http.StatusBadRequest)
	}
	return nil
}

type attributeInput struct {
	attribute            model.Attribute
	attributeValuesInput AttrValuesInput
}

func validateRequiredAttributes(inputData []attributeInput, attributes []*model.Attribute) *model_helper.AppError {
	providedAttributeIdsMap := lo.SliceToMap(inputData, func(item attributeInput) (string, bool) { return item.attribute.Id, true })

	if lo.SomeBy(attributes, func(item *model.Attribute) bool {
		return item.ValueRequired && !providedAttributeIdsMap[item.Id]
	}) {
		return model_helper.NewAppError("validateRequiredAttributes", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attributes"}, "All attributes flagged as having a value required must be supplied.", http.StatusBadRequest)
	}

	return nil
}

func validateAttributesInput(inputData []attributeInput, attributes []*model.Attribute, isPageAttributes, variantValidation bool) *model_helper.AppError {
	for _, item := range inputData {

		var appErr *model_helper.AppError
		switch item.attribute.InputType {
		case model.AttributeInputTypeFile:
			appErr = validateFileAttributesInput(&item.attribute, item.attributeValuesInput, variantValidation)
		case model.AttributeInputTypeReference:
			appErr = validateReferenceAttributesInput(&item.attribute, item.attributeValuesInput, variantValidation)
		case model.AttributeInputTypeRichText:
			appErr = validateRichTextAttributesInput(&item.attribute, item.attributeValuesInput, variantValidation)
		case model.AttributeInputTypeBoolean:
			appErr = validateBooleanInput(&item.attribute, item.attributeValuesInput, variantValidation)
		case model.AttributeInputTypeDate,
			model.AttributeInputTypeDateTime:
			appErr = validateDatetimeInput(&item.attribute, item.attributeValuesInput, variantValidation)
		default:
			appErr = validateStandardAttributesInput(&item.attribute, item.attributeValuesInput, variantValidation)
		}

		if appErr != nil {
			return appErr
		}
	}

	if !variantValidation {
		appErr := validateRequiredAttributes(inputData, attributes)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

type attributeAssignmentMixin struct {
	ctx *web.Context
}

func (a *attributeAssignmentMixin) preSaveNumericValues(attribute *model.Attribute, attrValues AttrValuesInput) (*model.AttributeValue, *model_helper.AppError) {
	if len(attrValues.Values) == 0 {
		return nil, nil
	}

	return a.ctx.App.Srv().AttributeService().UpsertAttributeValue(&model.AttributeValue{
		AttributeID: attribute.Id,
		Name:        attrValues.Values[0],
	})
}

func (a *attributeAssignmentMixin) preSaveValues(attribute *model.Attribute, attrValues AttrValuesInput) (model.AttributeValues, *model_helper.AppError) {
	existingAttributeValuesWithAttributeIDAndName, appErr := a.ctx.App.Srv().
		AttributeService().
		FilterAttributeValuesByOptions(model.AttributeValueFilterOptions{
			Conditions: squirrel.Eq{
				model.AttributeValueTableName + "." + model.AttributeValueColumnAttributeID: attribute.Id,
				model.AttributeValueTableName + "." + model.AttributeValueColumnName:        attrValues.Values,
			},
		})
	if appErr != nil {
		return nil, appErr
	}

	attributeValuesMap := map[string]*model.AttributeValue{} // keys have format of : attributeID_attributeValueName
	for _, value := range existingAttributeValuesWithAttributeIDAndName {
		attributeValuesMap[value.AttributeID+"_"+value.Name] = value
	}

	for _, name := range attrValues.Values {
		key := attribute.Id + "_" + name

		_, existed := attributeValuesMap[key]

		if !existed {
			attrValue, appErr := a.ctx.App.Srv().
				AttributeService().
				UpsertAttributeValue(&model.AttributeValue{
					AttributeID: attribute.Id,
					Name:        name,
				})
			if appErr != nil {
				return nil, appErr
			}

			attributeValuesMap[key] = attrValue
		}
	}

	return lo.Values(attributeValuesMap), nil
}

func (a *attributeAssignmentMixin) resolveAttributeNodes(attributes []*model.Attribute, globalIds, pks, slugs []string) ([]*model.Attribute, *model_helper.AppError) {
	pksMap := lo.SliceToMap(pks, func(item string) (string, bool) { return item, true })
	slugsMap := lo.SliceToMap(slugs, func(item string) (string, bool) { return item, true })

	attributes = lo.Filter(attributes, func(item *model.Attribute, index int) bool {
		return item != nil && (pksMap[item.Id] || slugsMap[item.Slug])
	})

	if len(attributes) == 0 {
		return nil, model_helper.NewAppError("resolveAttributeNodes", "app.attribute.cannot_resolve_node.app_error", nil, "could not resolve to a note", http.StatusBadRequest)
	}

	var (
		attributePksMap   = map[string]bool{}
		attributeSlugsMap = map[string]bool{}
	)
	for _, attribute := range attributes {
		attributePksMap[attribute.Id] = true
		attributeSlugsMap[attribute.Slug] = true
	}

	for i := 0; i < min(len(pks), len(globalIds)); i++ {
		if !attributePksMap[pks[i]] {
			return nil, model_helper.NewAppError("resolveAttributeNodes", "app.attribute.cannot_resolve_node.app_error", nil, "could not resolve id "+globalIds[i]+" to attribute", http.StatusBadRequest)
		}
	}

	for _, slug := range slugs {
		if !attributeSlugsMap[slug] {
			return nil, model_helper.NewAppError("resolveAttributeNodes", "app.attribute.cannot_resolve_node.app_error", nil, "could not resolve slug "+slug+" to attribute", http.StatusBadRequest)
		}
	}

	return attributes, nil
}
