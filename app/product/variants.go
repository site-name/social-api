package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

// GenerateAndSetVariantName Generate ProductVariant's name based on its attributes
func (a *AppProduct) GenerateAndSetVariantName(variant *product_and_discount.ProductVariant, sku string) *model.AppError {

	// a.AttributeApp().
}

// GetVariantSelectionAttributes Return attributes that can be used in variant selection.
//
// Attribute must be product attribute and attribute input type must be
// in ALLOWED_IN_VARIANT_SELECTION list.
func (a *AppProduct) GetVariantSelectionAttributes(attributes []*attribute.Attribute) []*attribute.Attribute {

	for i, attr := range attributes {
		if !util.StringInSlice(attr.InputType, attribute.ALLOWED_IN_VARIANT_SELECTION) || attr.Type != attribute.PRODUCT_TYPE {
			attributes = append(attributes[:i], attributes[i+1:]...)
		}
	}

	return attributes
}
