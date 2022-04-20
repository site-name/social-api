package gqlmodel

import (
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
)

// OLD ATTRIBUTE GRAPHQL TYPE for REFERENCE
// type Attribute struct {
// 	ID                       string                             `json:"id"`
// 	ProductTypes             *ProductTypeCountableConnection    `json:"productTypes"`
// 	ProductVariantTypes      *ProductTypeCountableConnection    `json:"productVariantTypes"`
// 	PrivateMetadata          []*MetadataItem                    `json:"privateMetadata"`
// 	Metadata                 []*MetadataItem                    `json:"metadata"`
// 	InputType                *AttributeInputTypeEnum            `json:"inputType"`
// 	EntityType               *AttributeEntityTypeEnum           `json:"entityType"`
// 	Name                     *string                            `json:"name"`
// 	Slug                     *string                            `json:"slug"`
// 	Type                     *AttributeTypeEnum                 `json:"type"`
// 	Unit                     *MeasurementUnitsEnum              `json:"unit"`
// 	Choices                  *AttributeValueCountableConnection `json:"choices"`
// 	ValueRequired            bool                               `json:"valueRequired"`
// 	VisibleInStorefront      bool                               `json:"visibleInStorefront"`
// 	FilterableInStorefront   bool                               `json:"filterableInStorefront"`
// 	FilterableInDashboard    bool                               `json:"filterableInDashboard"`
// 	AvailableInGrid          bool                               `json:"availableInGrid"`
// 	Translation              *AttributeTranslation              `json:"translation"`
// 	StorefrontSearchPosition int                                `json:"storefrontSearchPosition"`
// }

type Attribute struct {
	ID                       string                                    `json:"id"`
	ProductTypes             func() *ProductTypeCountableConnection    `json:"productTypes"`
	ProductVariantTypes      func() *ProductTypeCountableConnection    `json:"productVariantTypes"`
	PrivateMetadata          []*MetadataItem                           `json:"privateMetadata"`
	Metadata                 []*MetadataItem                           `json:"metadata"`
	InputType                *AttributeInputTypeEnum                   `json:"inputType"`
	EntityType               *AttributeEntityTypeEnum                  `json:"entityType"`
	Name                     *string                                   `json:"name"`
	Slug                     *string                                   `json:"slug"`
	Type                     *AttributeTypeEnum                        `json:"type"`
	Unit                     *MeasurementUnitsEnum                     `json:"unit"`
	Choices                  func() *AttributeValueCountableConnection `json:"choices"`
	ValueRequired            bool                                      `json:"valueRequired"`            // permission required to view
	VisibleInStorefront      bool                                      `json:"visibleInStorefront"`      // permission required to view
	FilterableInStorefront   bool                                      `json:"filterableInStorefront"`   // permission required to view
	FilterableInDashboard    bool                                      `json:"filterableInDashboard"`    // permission required to view
	AvailableInGrid          bool                                      `json:"availableInGrid"`          // permission required to view
	StorefrontSearchPosition int                                       `json:"storefrontSearchPosition"` // permission required to view
	TranslationID            *string                                   `json:"translation"`
}

func (Attribute) IsNode()               {}
func (Attribute) IsObjectWithMetadata() {}

// ModelAttributeToGraphqlAttribute converts database *Attribute to graphql *Attribute
func ModelAttributeToGraphqlAttribute(a *attribute.Attribute) *Attribute {
	if a == nil {
		return nil
	}

	inputType := AttributeInputTypeEnum(strings.ToUpper(string(a.InputType)))
	attrType := AttributeTypeEnum(strings.ToUpper(a.Type))

	var unit *MeasurementUnitsEnum
	if a.Unit != nil {
		u := MeasurementUnitsEnum(strings.ToUpper(*a.Unit))
		unit = &u
	}

	var entityType *AttributeEntityTypeEnum
	if a.EntityType != nil {
		t := AttributeEntityTypeEnum(strings.ToUpper(*a.EntityType))
		entityType = &t
	}

	// NOTE: only return these fields since some other fields require specific permission to see
	return &Attribute{
		ID:              a.Id,
		Name:            &a.Name,
		Slug:            &a.Slug,
		Metadata:        MapToGraphqlMetaDataItems(a.Metadata),
		PrivateMetadata: MapToGraphqlMetaDataItems(a.PrivateMetadata),
		InputType:       &inputType,
		EntityType:      entityType,
		Type:            &attrType,
		Unit:            unit,
	}
}

//---------------------------------------------
// initial implementation
// type AttributeValue struct {
// 	ID          string                     `json:"id"`
// 	Name        *string                    `json:"name"`
// 	Slug        *string                    `json:"slug"`
// 	Value       *string                    `json:"value"`
// 	Translation *AttributeValueTranslation `json:"translation"`
// 	InputType   *AttributeInputTypeEnum    `json:"inputType"`
// 	Reference   *string                    `json:"reference"`
// 	File        *File                      `json:"file"`
// 	RichText    model.StringInterface      `json:"richText"`
// 	Boolean     *bool                      `json:"boolean"`
// 	Date        *time.Time                 `json:"date"`
// 	DateTime    *time.Time                 `json:"dateTime"`
// }

type AttributeValue struct {
	ID          string                            `json:"id"`
	Name        *string                           `json:"name"`
	Slug        *string                           `json:"slug"`
	Value       *string                           `json:"value"`
	Translation func() *AttributeValueTranslation `json:"translation"`
	InputType   func() *AttributeInputTypeEnum    `json:"inputType"`
	Reference   func() *string                    `json:"reference"`
	File        *File                             `json:"file"`
	RichText    model.StringInterface             `json:"richText"`
	Boolean     *bool                             `json:"boolean"`
	Date        *time.Time                        `json:"date"`
	DateTime    *time.Time                        `json:"dateTime"`

	AttributeID string `json:"-"`
}

func (AttributeValue) IsNode() {}

func ModelAttributeValueToGraphqlAttributeValue(a *attribute.AttributeValue) *AttributeValue {
	if a == nil {
		return nil
	}

	res := &AttributeValue{
		ID:       a.Id,
		Name:     &a.Name,
		Slug:     &a.Slug,
		Value:    &a.Value,
		RichText: a.RichText,
		Boolean:  a.Boolean,
		File: &File{
			ContentType: a.ContentType,
		},
	}

	if a.FileUrl != nil {
		res.File.URL = *a.FileUrl
	}

	if a.Attribute != nil {
		// set DateTime
		if a.Attribute.InputType == attribute.DATE_TIME {
			res.DateTime = a.Datetime
		}

		// set Date
		if a.Attribute.InputType == attribute.DATE {
			res.Date = a.Datetime
		}
	}

	return res
}

func (i *AttributeCreateInput) GetValueByField(fieldName string) interface{} {
	switch fieldName {
	case "FilterableInStorefront":
		return i.FilterableInStorefront
	case "FilterableInDashboard":
		return i.FilterableInDashboard
	case "StorefrontSearchPosition":
		return i.StorefrontSearchPosition
	case "AvailableInGrid":
		return i.AvailableInGrid

	default:
		return nil
	}
}

func (i *AttributeUpdateInput) GetValueByField(fieldName string) interface{} {
	switch fieldName {
	case "FilterableInStorefront":
		return i.FilterableInStorefront
	case "FilterableInDashboard":
		return i.FilterableInDashboard
	case "StorefrontSearchPosition":
		return i.StorefrontSearchPosition
	case "AvailableInGrid":
		return i.AvailableInGrid

	default:
		return nil
	}
}

var ATTRIBUTE_PROPERTIES_CONFIGURATION = map[string][]AttributeInputTypeEnum{
	"FilterableInStorefront": {
		AttributeInputTypeEnumDropdown,
		AttributeInputTypeEnumMultiselect,
		AttributeInputTypeEnumNumeric,
		AttributeInputTypeEnumSwatch,
		AttributeInputTypeEnumBoolean,
		AttributeInputTypeEnumDate,
		AttributeInputTypeEnumDateTime,
	},
	"FilterableInDashboard": {
		AttributeInputTypeEnumDropdown,
		AttributeInputTypeEnumMultiselect,
		AttributeInputTypeEnumNumeric,
		AttributeInputTypeEnumSwatch,
		AttributeInputTypeEnumBoolean,
		AttributeInputTypeEnumDate,
		AttributeInputTypeEnumDateTime,
	},
	"AvailableInGrid": {
		AttributeInputTypeEnumDropdown,
		AttributeInputTypeEnumMultiselect,
		AttributeInputTypeEnumNumeric,
		AttributeInputTypeEnumSwatch,
		AttributeInputTypeEnumBoolean,
		AttributeInputTypeEnumDate,
		AttributeInputTypeEnumDateTime,
	},
	"StorefrontSearchPosition": {
		AttributeInputTypeEnumDropdown,
		AttributeInputTypeEnumMultiselect,
		AttributeInputTypeEnumBoolean,
		AttributeInputTypeEnumDate,
		AttributeInputTypeEnumDateTime,
	},
}

func AttributeInputTypeEnumInSlice(value AttributeInputTypeEnum, slice ...AttributeInputTypeEnum) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

func ModelAttributeValueTranslationToGraphqlAttributeValueTranslation(t *attribute.AttributeValueTranslation) *AttributeValueTranslation {
	if t == nil {
		return nil
	}

	res := &AttributeValueTranslation{
		ID:       t.Id,
		Name:     t.Name,
		RichText: t.RichText,
		Language: &LanguageDisplay{
			Code:     LanguageCodeEnum(strings.ToUpper(t.LanguageCode)),
			Language: model.Languages[strings.ToLower(t.LanguageCode)],
		},
	}

	return res
}

type AttributeValueCreateUpdateInterface interface {
	GetName() *string // should not be nil
	GetValue() *string
	GetRichText() model.StringInterface
	GetFileURL() *string
	GetContentType() *string
}

var (
	_ AttributeValueCreateUpdateInterface = (*AttributeValueCreateInput)(nil)
	_ AttributeValueCreateUpdateInterface = (*AttributeValueUpdateInput)(nil)
)

func (a *AttributeValueCreateInput) GetName() *string {
	return &a.Name
}
func (a *AttributeValueCreateInput) GetValue() *string {
	return a.Value
}
func (a *AttributeValueCreateInput) GetFileURL() *string {
	return a.FileURL
}
func (a *AttributeValueCreateInput) GetContentType() *string {
	return a.ContentType
}
func (a *AttributeValueCreateInput) GetRichText() model.StringInterface {
	return a.RichText
}

//----------------------------
func (a *AttributeValueUpdateInput) GetName() *string {
	if a.Name == nil {
		return model.NewString("")
	}
	return a.Name
}
func (a *AttributeValueUpdateInput) GetValue() *string {
	return a.Value
}
func (a *AttributeValueUpdateInput) GetFileURL() *string {
	return a.FileURL
}
func (a *AttributeValueUpdateInput) GetContentType() *string {
	return a.ContentType
}
func (a *AttributeValueUpdateInput) GetRichText() model.StringInterface {
	return a.RichText
}
