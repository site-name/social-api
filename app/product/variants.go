package product

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

// GenerateAndSetVariantName Generate ProductVariant's name based on its attributes
func (a *ServiceProduct) GenerateAndSetVariantName(variant *model.ProductVariant, sku string) *model.AppError {
	_, _ = a.srv.AttributeService().AssignedVariantAttributesByOption(&model.AssignedVariantAttributeFilterOption{
		AssignmentAttributeInputType: squirrel.Eq{model.AttributeTableName + ".InputType": model.ALLOWED_IN_VARIANT_SELECTION},
		AssignmentAttributeType:      squirrel.Eq{model.AttributeTableName + ".Type": model.PRODUCT_TYPE},
	})
	panic("not implt")
}

// GetVariantSelectionAttributes Return attributes that can be used in variant selection.
//
// Attribute must be product attribute and attribute input type must be
// in ALLOWED_IN_VARIANT_SELECTION list.
func (a *ServiceProduct) GetVariantSelectionAttributes(attributes []*model.Attribute) []*model.Attribute {

	for i, attr := range attributes {
		if !model.ALLOWED_IN_VARIANT_SELECTION.Contains(attr.InputType) ||
			attr.Type != model.PRODUCT_TYPE {
			attributes = append(attributes[:i], attributes[i+1:]...)
		}
	}

	return attributes
}
