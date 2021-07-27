package gqlmodel

import (
	"strings"

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

	inputType := AttributeInputTypeEnum(strings.ToUpper(a.InputType))
	attrType := AttributeTypeEnum(strings.ToUpper(a.Type))

	var unit *MeasurementUnitsEnum = nil
	if a.Unit != nil {
		u := MeasurementUnitsEnum(strings.ToUpper(*a.Unit))
		unit = &u
	}

	var entityType *AttributeEntityTypeEnum = nil
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

// ModelAttributeTranslationToGraphqlAttributeTranslation converts database attribute translation to graphql attribute translation
// func ModelAttributeTranslationToGraphqlAttributeTranslation(t *attribute.AttributeTranslation) *AttributeTranslation {
// 	if t == nil {
// 		return nil
// 	}

// 	return &AttributeTranslation{
// 		ID: t.Id,
// 		Name: t.Name,
// 		Language: &LanguageDisplay{
// 			Code: LanguageCodeEnum(stringToGraphqlEnumString(t.LanguageCode)),
// 		},
// 	}
// }
