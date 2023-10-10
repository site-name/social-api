package product

import (
	"sort"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
)

// GenerateAndSetVariantName Generate ProductVariant's name based on its attributes
func (a *ServiceProduct) GenerateAndSetVariantName(variant *model.ProductVariant, sku string) *model.AppError {
	assignedVariantAttributes, appErr := a.srv.
		AttributeService().
		AssignedVariantAttributesByOption(&model.AssignedVariantAttributeFilterOption{
			Preloads:   []string{"Values"},
			Conditions: squirrel.Eq{model.AssignedVariantAttributeTableName + "." + model.AssignedVariantAttributeColumnVariantID: variant.Id},
			Assignment_Attribute_Conditions: squirrel.Eq{
				model.AttributeTableName + "." + model.AttributeColumnInputType: model.ALLOWED_IN_VARIANT_SELECTION,
				model.AttributeTableName + "." + model.AttributeColumnType:      model.PRODUCT_TYPE,
			},
		})

	if appErr != nil {
		return appErr
	}

	attributesDisplay := make([]string, len(assignedVariantAttributes))

	for idx, item := range assignedVariantAttributes {
		valueNames := lo.Map(item.Values, func(value *model.AttributeValue, _ int) string { return value.Name })

		attributesDisplay[idx] = strings.Join(valueNames, ", ")
	}

	sort.Sort(sort.StringSlice(attributesDisplay))
	name := strings.Join(attributesDisplay, " / ")

	if name == "" {
		if sku != "" {
			name = sku
		} else {
			name = variant.Id
		}
	}

	variant.Name = name
	_, appErr = a.UpsertProductVariant(nil, variant)
	return appErr
}

// GetVariantSelectionAttributes Return attributes that can be used in variant selection.
//
// Attribute must be product attribute and attribute input type must be
// in ALLOWED_IN_VARIANT_SELECTION list.
func (a *ServiceProduct) GetVariantSelectionAttributes(attributes []*model.Attribute) []*model.Attribute {
	return lo.Filter(attributes, func(item *model.Attribute, _ int) bool {
		return model.ALLOWED_IN_VARIANT_SELECTION.Contains(item.InputType) &&
			item.Type == model.PRODUCT_TYPE
	})
}
