package attribute

import (
	"net/http"
	"sort"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// AssociateAttributeValuesToInstance Assign given attribute values to a product or variant.
// Note: be award this function invokes the ``set`` method on the instance's
// attribute association. Meaning any values already assigned or concurrently
// assigned will be overridden by this call.
//
// `instance` must be either *Product or *ProductVariant or *Page.
// `attributeID` must be ID of processing `Attribute`
//
// Returned interface{} must be either: `*AssignedProductAttribute` or `*AssignedVariantAttribute` or `*AssignedPageAttribute`
func (a *ServiceAttribute) AssociateAttributeValuesToInstance(instance interface{}, attributeID string, values []*attribute.AttributeValue) (interface{}, *model.AppError) {

	// validate if valid `instance` was provided
	switch instance.(type) {
	case *product_and_discount.Product, *product_and_discount.ProductVariant, *page.Page:
	default:
		return nil, model.NewAppError("AssociateAttributeValuesToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "", http.StatusBadRequest)
	}

	valueIDs := attribute.AttributeValues(values).IDs()

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
	case *attribute.AssignedProductAttribute:
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

	case *attribute.AssignedVariantAttribute:
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

	case *attribute.AssignedPageAttribute:
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
func (a *ServiceAttribute) validateAttributeOwnsValues(attributeID string, valueIDs []string) *model.AppError {
	attributeValues, appErr := a.AttributeValuesOfAttribute(attributeID)
	if appErr != nil {
		return appErr
	}
	attributeActualValueIDs := attribute.AttributeValues(attributeValues).IDs()
	foundAssociatedIDs := util.StringArrayIntersection(attributeActualValueIDs, valueIDs)

	for _, associatedID := range foundAssociatedIDs {
		if !util.StringInSlice(associatedID, valueIDs) {
			return model.NewAppError("validateAttributeOwnsValues", "app.attribute.attribute_missing_some_values", nil, "", http.StatusNotFound)
		}
	}

	return nil
}

// associateAttributeToInstance associates given attribute to given instance
//
// NOTE:
//
// `instance` must be either `*product.Product` or `*product.ProductVariant` or `*page.Page`
//
// returned interface{} is either:
//  +) *AssignedProductAttribute
//  +) *AssignedVariantAttribute
//  +) *AssignedPageAttribute
func (a *ServiceAttribute) associateAttributeToInstance(instance interface{}, attributeID string) (interface{}, *model.AppError) {

	switch v := instance.(type) {
	case *product_and_discount.Product:
		attributeProduct, appErr := a.AttributeProductByOption(&attribute.AttributeProductFilterOption{
			ProductTypeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: v.ProductTypeID,
				},
			},
			AttributeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: attributeID,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedProductAttribute(&attribute.AssignedProductAttribute{
			ProductID:    v.Id,
			AssignmentID: attributeProduct.Id,
		})

	case *product_and_discount.ProductVariant:
		attrVariant, appErr := a.AttributeVariantByOption(&attribute.AttributeVariantFilterOption{
			ProductTypeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: v.ProductID,
				},
			},
			AttributeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: attributeID,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedVariantAttribute(&attribute.AssignedVariantAttribute{
			VariantID:    v.Id,
			AssignmentID: attrVariant.Id,
		})

	case *page.Page:
		attributePage, appErr := a.AttributePageByOption(&attribute.AttributePageFilterOption{
			AttributeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: attributeID,
				},
			},
			PageTypeID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: v.PageTypeID,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		return a.GetOrCreateAssignedPageAttribute(&attribute.AssignedPageAttribute{
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
//  +) *Product        - *AssignedProductAttribute
//
//  +) *ProductVariant - *AssignedVariantAttribute
//
//  +) *Page           - *AssignedPageAttribute
func (a *ServiceAttribute) sortAssignedAttributeValues(instance interface{}, assignment interface{}, valueIDs []string) *model.AppError {
	// validate if `instance` and `assignment` are provided accordingly:
	invalidArgumentErrorHandler := func(field string) *model.AppError {
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": field}, "Please read doc for this method", http.StatusBadRequest)
	}

	if instance == nil || assignment == nil || len(valueIDs) == 0 {
		return invalidArgumentErrorHandler("assignment or instance or valuesIDs")
	}

	switch instance.(type) {
	case *product_and_discount.Product:
		if assignmentValue, ok := assignment.(*attribute.AssignedProductAttribute); !ok {
			assignedProductAttrValues, attrValues, err := a.srv.Store.AssignedProductAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return store.AppErrorFromDatabaseLookupError("sortAssignedAttributeValues", "app.attribute.select_assigned_product_attribute_values_for_sort.app_error", err)
			}
			// NOTE: this sort can be done since len(assignedProductAttrValues) == len(attrValues)
			sort.Slice(assignedProductAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedProductAttrValues {
				value.SortOrder = i
			}
			// update if database:
			if err = a.srv.Store.AssignedProductAttributeValue().UpdateInBulk(assignedProductAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_product_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return invalidArgumentErrorHandler("assignment")
	case *product_and_discount.ProductVariant:
		if assignmentValue, ok := assignment.(*attribute.AssignedVariantAttribute); ok {
			assignedVariantAttrValues, attrValues, err := a.srv.Store.AssignedVariantAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return store.AppErrorFromDatabaseLookupError("sortAssignedAttributeValues", "app.attribute.select_assigned_variant_attribute_values_for_sort.app_error", err)
			}
			// NOTE: this sort can be done since len(assignedVariantAttrValues) == len(attrValues)
			sort.Slice(assignedVariantAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedVariantAttrValues {
				value.SortOrder = i
			}
			// update if database:
			if err = a.srv.Store.AssignedVariantAttributeValue().UpdateInBulk(assignedVariantAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_variant_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return invalidArgumentErrorHandler("assignment")
	case *page.Page:
		if assignmentValue, ok := assignment.(*attribute.AssignedPageAttribute); ok {
			assignedPageAttrValues, attrValues, err := a.srv.Store.AssignedPageAttributeValue().SelectForSort(assignmentValue.Id)
			// err can be *store.ErrNotFound or system error
			if err != nil {
				return store.AppErrorFromDatabaseLookupError("sortAssignedAttributeValues", "app.attribute.select_assigned_page_attribute_values_for_sort.app_error", err)
			}
			// NOTE: this sort can be done since len(assignedPageAttrValues) == len(attrValues)
			sort.Slice(assignedPageAttrValues, func(i, j int) bool {
				return sort.SearchStrings(valueIDs, attrValues[i].Id) <= sort.SearchStrings(valueIDs, attrValues[j].Id)
			})
			for i, value := range assignedPageAttrValues {
				value.SortOrder = i
			}
			// update if database:
			if err = a.srv.Store.AssignedPageAttributeValue().UpdateInBulk(assignedPageAttrValues); err != nil {
				return model.NewAppError("sortAssignedAttributeValues", "app.attribute.error_updating_assigned_page_attribute_values.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}
		// other types are not accepted and returns an error:
		return invalidArgumentErrorHandler("assignment")

	default:
		return invalidArgumentErrorHandler("instance")
	}
}
