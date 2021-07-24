package attribute

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// AssociateAttributeValuesToInstance assigns given attribute values to a product or variant.
//
// `instance` must be either `*product.Product` or `*product.ProductVariant` or `*page.Page`
//
// `attributeID` must be ID of processing `Attribute`
//
// Returned interface{} must be either: `*AssignedProductAttribute` or `*AssignedVariantAttribute` or `*AssignedPageAttribute`
func (a *AppAttribute) AssociateAttributeValuesToInstance(instance interface{}, attributeID string, values []*attribute.AttributeValue) (interface{}, *model.AppError) {

	// validate if valid `instance` provided`
	switch instance.(type) {
	case *product_and_discount.Product, *product_and_discount.ProductVariant, *page.Page:
		// do nothing since these cases are valid.
	default:
		return nil, model.NewAppError("AssociateAttributeValuesToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "Please check doc for this method", http.StatusBadRequest)
	}

	valueIDs := make([]string, len(values))
	for i, value := range values {
		valueIDs[i] = value.Id
	}

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
		_, err := a.app.Srv().Store.AssignedProductAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			return nil, commonErrHandler(err, "AssociateAttributeToInstance", "AssignmentID", "ValueIDs")
		}
	case *attribute.AssignedVariantAttribute:
		_, err := a.app.Srv().Store.AssignedVariantAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			return nil, commonErrHandler(err, "AssociateAttributeToInstance", "AssignmentID", "ValueIDs")
		}
	case *attribute.AssignedPageAttribute:
		_, err := a.app.Srv().Store.AssignedPageAttributeValue().SaveInBulk(v.Id, valueIDs)
		if err != nil {
			return nil, commonErrHandler(err, "AssociateAttributeToInstance", "AssignmentID", "ValueIDs")
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

// commonErrHandler
func commonErrHandler(err error, where string, fields ...string) *model.AppError {
	if err == nil {
		return nil
	}

	switch t := err.(type) {
	case *model.AppError:
		return t
	case *store.ErrInvalidInput:
		return model.NewAppError(where, app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": strings.Join(fields, ", ")}, err.Error(), http.StatusBadRequest)
	default:
		return model.NewAppError(where, app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
}

// validateAttributeOwnsValues Checks given value IDs are belonging to the given attribute.
func (a *AppAttribute) validateAttributeOwnsValues(attributeID string, valueIDs []string) *model.AppError {
	attributeValues, appErr := a.AttributeValuesOfAttribute(attributeID)
	if appErr != nil {
		return appErr
	}

	sameLength := len(attributeValues) == len(valueIDs)
	sameIds := true
	for _, av := range attributeValues {
		if !util.StringInSlice(av.Id, valueIDs) {
			sameIds = false
			break
		}
	}

	if !sameLength || !sameIds {
		return model.NewAppError(
			"validateAttributeOwnsValues", app.InvalidArgumentAppErrorID,
			map[string]interface{}{
				"Fields": "valueIDs",
			},
			"",
			http.StatusBadRequest,
		)
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
//  +) *attribute.AssignedProductAttribute
//  +) *attribute.AssignedVariantAttribute
//  +) *attribute.AssignedPageAttribute
func (a *AppAttribute) associateAttributeToInstance(instance interface{}, attributeID string) (interface{}, *model.AppError) {

	switch v := instance.(type) {
	case *product_and_discount.Product:
		attributeProduct, err := a.
			app.Srv().Store.AttributeProduct().GetByOption(&attribute.AttributeProductGetOption{
			AttributeID:   attributeID,
			ProductTypeID: v.ProductTypeID,
		})
		if err != nil {
			if invlIp, ok := err.(*store.ErrInvalidInput); ok {
				return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option"}, invlIp.Error(), http.StatusBadRequest)
			}
			return nil, store.AppErrorFromDatabaseLookupError("associateAttributeToInstance", "app.attribute.attribute_product_by_option", err)
		}

		assignedProductAttribute, err := a.app.Srv().Store.AssignedProductAttribute().GetWithOption(&attribute.AssignedProductAttributeFilterOption{
			ProductID:    v.Id,
			AssignmentID: attributeProduct.Id,
		})
		if err != nil { // this error can be either: `*AppError` or `*store.ErrInvalidInput` or `system error`
			return nil, commonErrHandler(err, "associateAttributeToInstance", "option")
		}

		return assignedProductAttribute, nil
	case *product_and_discount.ProductVariant:
		attributeVariant, err := a.app.Srv().Store.AttributeVariant().GetByOption(&attribute.AttributeVariantFilterOption{
			ProductID:   v.ProductID,
			AttributeID: attributeID,
		})
		if err != nil {
			// error input is handled manually:
			if invlErr, ok := err.(*store.ErrInvalidInput); ok {
				return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option"}, invlErr.Error(), http.StatusBadRequest)
			}
			// system error, not found error:
			return nil, store.AppErrorFromDatabaseLookupError("associateAttributeToInstance", "app.attribute.error_finding_attribute_variant.app_error", err)
		}
		assignedVariantAttribute, err := a.app.Srv().Store.AssignedVariantAttribute().GetWithOption(&attribute.AssignedVariantAttributeFilterOption{
			VariantID:    v.Id,
			AssignmentID: attributeVariant.Id,
		})
		if err != nil { // this error can be either: `*AppError` or `*store.ErrInvalidInput` or `system error`
			return nil, commonErrHandler(err, "associateAttributeToInstance", "option")
		}

		return assignedVariantAttribute, nil
	case *page.Page:
		attributePage, err := a.app.Srv().Store.AttributePage().GetByOption(&attribute.AttributePageFilterOption{
			AttributeID: attributeID,
			PageTypeID:  v.PageTypeID,
		})
		if err != nil {
			// error input is handled manually:
			if invlErr, ok := err.(*store.ErrInvalidInput); ok {
				return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option"}, invlErr.Error(), http.StatusBadRequest)
			}
			// system error, not found error:
			return nil, store.AppErrorFromDatabaseLookupError("associateAttributeToInstance", "app.attribute.error_finding_attribute_page.app_error", err)
		}
		assignedPageAttribute, err := a.app.Srv().Store.AssignedPageAttribute().GetByOption(&attribute.AssignedPageAttributeFilterOption{
			PageID:       v.Id,
			AssignmentID: attributePage.Id,
		})
		if err != nil { // this error can be either: `*AppError` or `*store.ErrInvalidInput` or `system error`
			return nil, commonErrHandler(err, "associateAttributeToInstance", "option")
		}

		return assignedPageAttribute, nil
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
func (a *AppAttribute) sortAssignedAttributeValues(instance interface{}, assignment interface{}, valueIDs []string) *model.AppError {
	// validate if `instance` and `assignment` are provided accordingly:
	invalidArgumentErrorHandler := func(fields ...string) *model.AppError {
		return model.NewAppError("sortAssignedAttributeValues", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": strings.Join(fields, ", ")}, "Please read doc for this method", http.StatusBadRequest)
	}

	if instance == nil || assignment == nil {
		return invalidArgumentErrorHandler("assignment", "instance")
	}

	switch instanceType := instance.(type) {
	case *product_and_discount.Product:
		if _, ok := assignment.(*attribute.AssignedProductAttribute); !ok {
			return invalidArgumentErrorHandler("assignment")
		}
	case *product_and_discount.ProductVariant:
		if _, ok := assignment.(*attribute.AssignedVariantAttribute); !ok {
			return invalidArgumentErrorHandler("assignment")
		}
	case *page.Page:
		if _, ok := assignment.(*attribute.AssignedPageAttribute); !ok {
			return invalidArgumentErrorHandler("assignment")
		}

	default:
		return invalidArgumentErrorHandler("instance")
	}
}
