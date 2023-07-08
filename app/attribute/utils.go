package attribute

import (
	"net/http"
	"sort"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// AssociateAttributeValuesToInstance Assign given attribute values to a product or variant.
// Note: be award this function invokes the “set“ method on the instance's
// attribute association. Meaning any values already assigned or concurrently
// assigned will be overridden by this call.
//
// `instance` must be either *Product or *ProductVariant or *Page.
// `attributeID` must be ID of processing `Attribute`
//
// Returned interface{} must be either: `*AssignedProductAttribute` or `*AssignedVariantAttribute` or `*AssignedPageAttribute`
func (a *ServiceAttribute) AssociateAttributeValuesToInstance(instance interface{}, attributeID string, values model.AttributeValues) (interface{}, *model.AppError) {

	// validate if valid `instance` was provided
	switch instance.(type) {
	case *model.Product, *model.ProductVariant, *model.Page:
	default:
		return nil, model.NewAppError("AssociateAttributeValuesToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "", http.StatusBadRequest)
	}

	valueIDs := values.IDs()

	// Ensure the values are actually form the given attribute:
	if appErr := a.validateAttributeOwnsValues(attributeID, valueIDs); appErr != nil {
		return nil, appErr
	}

	// Associate the attribute and the passed values
	assignment, appErr := a.associateAttributeToInstance(instance, attributeID)
	if appErr != nil {
		return nil, appErr
	}

	// save in bulk value relationships:
	switch v := assignment.(type) {
	case *model.AssignedProductAttribute:
		_, err := a.srv.Store.AssignedProductAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusNotFound
			}
			return nil, model.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.error_creating_assigned_product_attribute_values.app_error", nil, err.Error(), statusCode)
		}

	case *model.AssignedVariantAttribute:
		_, err := a.srv.Store.AssignedVariantAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusNotFound
			}
			return nil, model.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.error_creating_assigned_variants_attribute_values.app_error", nil, err.Error(), statusCode)
		}

	case *model.AssignedPageAttribute:
		_, err := a.srv.Store.AssignedPageAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusNotFound
			}
			return nil, model.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.error_creating_assigned_page_attribute_values.app_error", nil, err.Error(), statusCode)
		}

	default:
		return nil, model.NewAppError("AssociateAttributeValuesToInstance", "app.attribute.unknown_returned_assignment.app_error", nil, "", http.StatusNotAcceptable)
	}

	// sort assigned attribute values:
	appErr = a.sortAssignedAttributeValues(instance, assignment, valueIDs)
	if appErr != nil {
		return nil, appErr
	}

	return assignment, nil
}

// validateAttributeOwnsValues Checks given value IDs are belonging to the given attribute.
func (a *ServiceAttribute) validateAttributeOwnsValues(attributeID string, valueIDs util.AnyArray[string]) *model.AppError {
	attributeValues, appErr := a.AttributeValuesOfAttribute(attributeID)
	if appErr != nil {
		return appErr
	}
	attributeActualValueIDs := attributeValues.IDs()
	foundAssociatedIDs := valueIDs.InterSection(attributeActualValueIDs...)

	for _, associatedID := range foundAssociatedIDs {
		if !valueIDs.Contains(associatedID) {
			return model.NewAppError("validateAttributeOwnsValues", "app.attribute.attribute_missing_some_values", nil, "", http.StatusNotFound)
		}
	}

	return nil
}

// associateAttributeToInstance associates given attribute to given instance
//
// NOTE:
//
// `instance` must be either `*product.Product` or `*product.ProductVariant` or `*model.Page`
//
// returned interface{} is either:
//
//	+) *AssignedProductAttribute
//	+) *AssignedVariantAttribute
//	+) *AssignedPageAttribute
func (a *ServiceAttribute) associateAttributeToInstance(instance interface{}, attributeID string) (interface{}, *model.AppError) {
	switch v := instance.(type) {
	case *model.Product:
		attributeProduct, appErr := a.AttributeProductByOption(&model.AttributeProductFilterOption{
			ProductTypeID: squirrel.Eq{store.AttributeProductTableName + ".ProductTypeID": v.ProductTypeID},
			AttributeID:   squirrel.Eq{store.AttributeProductTableName + ".AttributeID": attributeID},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedProductAttribute(&model.AssignedProductAttribute{
			ProductID:    v.Id,
			AssignmentID: attributeProduct.Id,
		})

	case *model.ProductVariant:
		attrVariant, appErr := a.AttributeVariantByOption(&model.AttributeVariantFilterOption{
			ProductTypeID: squirrel.Eq{store.AttributeVariantTableName + ".ProductTypeID": v.ProductID},
			AttributeID:   squirrel.Eq{store.AttributeVariantTableName + ".AttributeID": attributeID},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedVariantAttribute(&model.AssignedVariantAttribute{
			VariantID:    v.Id,
			AssignmentID: attrVariant.Id,
		})

	case *model.Page:
		attributePage, appErr := a.AttributePageByOption(&model.AttributePageFilterOption{
			AttributeID: squirrel.Eq{store.AttributePageTableName + ".AttributeID": attributeID},
			PageTypeID:  squirrel.Eq{store.AttributePageTableName + ".PageTypeID": v.PageTypeID},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedPageAttribute(&model.AssignedPageAttribute{
			PageID:       v.Id,
			AssignmentID: attributePage.Id,
		})

	default:
		return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "", http.StatusBadRequest)
	}
}

// sortAssignedAttributeValues Sorts assigned attribute values based on values list order
//
// (instance - assignment) must be provided by pair like this:
//
//	+) *Product        - *AssignedProductAttribute
//
//	+) *ProductVariant - *AssignedVariantAttribute
//
//	+) *Page           - *AssignedPageAttribute
func (a *ServiceAttribute) sortAssignedAttributeValues(instance interface{}, assignment interface{}, valueIDs []string) *model.AppError {
	if instance == nil || assignment == nil || len(valueIDs) == 0 {
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "assignment or instance or valuesIDs"}, "", http.StatusBadRequest)
	}

	switch instance.(type) {
	case *model.Product:
		if assignmentValue, ok := assignment.(*model.AssignedProductAttribute); !ok {
			assignedProductAttrValues, attrValues, err := a.srv.Store.AssignedProductAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_product_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			// NOTE: this sort can be done since len(assignedProductAttrValues) == len(attrValues)
			sort.Slice(assignedProductAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedProductAttrValues {
				value.SortOrder = &i
			}
			// update if database:
			if err = a.srv.Store.AssignedProductAttributeValue().UpdateInBulk(assignedProductAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_product_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "assignment"}, "", http.StatusBadRequest)

	case *model.ProductVariant:
		if assignmentValue, ok := assignment.(*model.AssignedVariantAttribute); ok {
			assignedVariantAttrValues, attrValues, err := a.srv.Store.AssignedVariantAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_variant_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			// NOTE: this sort can be done since len(assignedVariantAttrValues) == len(attrValues)
			sort.Slice(assignedVariantAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedVariantAttrValues {
				value.SortOrder = &i
			}
			// update if database:
			if err = a.srv.Store.AssignedVariantAttributeValue().UpdateInBulk(assignedVariantAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_variant_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "assignment"}, "", http.StatusBadRequest)

	case *model.Page:
		if assignmentValue, ok := assignment.(*model.AssignedPageAttribute); ok {
			assignedPageAttrValues, attrValues, err := a.srv.Store.AssignedPageAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.select_assigned_page_attribute_values_for_sort.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			// NOTE: this sort can be done since len(assignedPageAttrValues) == len(attrValues)
			sort.Slice(assignedPageAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedPageAttrValues {
				value.SortOrder = &i
			}
			// update if database:
			if err = a.srv.Store.AssignedPageAttributeValue().UpdateInBulk(assignedPageAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_page_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "assignment"}, "", http.StatusBadRequest)

	default:
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "", http.StatusBadRequest)
	}
}
