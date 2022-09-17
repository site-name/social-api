package model

import (
	"io"

	"github.com/Masterminds/squirrel"
)

// AttributeID unique together with ProductTypeID
type AttributeVariant struct {
	Id               string `json:"id"`
	AttributeID      string `json:"attribute_id"`
	ProductTypeID    string `json:"product_type_id"`
	VariantSelection bool   `json:"variant_selection"`
	Sortable
}

// AttributeVariantFilterOption is used to find `AttributeVariant`.
//
// properties can be provided partially or fully
type AttributeVariantFilterOption struct {
	Id            squirrel.Sqlizer
	AttributeID   squirrel.Sqlizer
	ProductTypeID squirrel.Sqlizer
	ProductIDs    []string
}

func (a *AttributeVariant) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"attribute_variant.is_valid.%s.app_error",
		"attribute_variant_id=",
		"AttributeVariant.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AttributeID) {
		return outer("attribute_id", &a.Id)
	}
	if !IsValidId(a.ProductTypeID) {
		return outer("product_type_id", &a.Id)
	}

	return nil
}

func (a *AttributeVariant) ToJSON() string {
	return ModelToJson(a)
}

func AttributeVariantFromJson(data io.Reader) *AttributeVariant {
	var a AttributeVariant
	ModelFromJson(&a, data)
	return &a
}

func (a *AttributeVariant) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
}
