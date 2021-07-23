package attribute

import (
	"net/http"

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
// `valueIDs` must be a list of IDs of `AttributeValues`
//
// Returned interface{} must be either: `*AssignedProductAttribute` or `*AssignedVariantAttribute` or `*AssignedPageAttribute`
func (a *AppAttribute) AssociateAttributeValuesToInstance(instance interface{}, attributeID string, valueIDs []string) (interface{}, *model.AppError) {

	// Ensure the values are actually form the given attribute:
	if appErr := a.validateAttributeOwnsValues(attributeID, valueIDs); appErr != nil {
		return nil, appErr
	}

	// Associate the attribute and the passed values
	if appErr := a.validateAttributeOwnsValues(attributeID, valueIDs); appErr != nil {
		return nil, appErr
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
			switch t := err.(type) {
			case *model.AppError:
				return nil, t
			case *store.ErrInvalidInput:
				return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option"}, err.Error(), http.StatusBadRequest)
			default:
				return nil, model.NewAppError("associateAttributeToInstance", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError)
			}
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
			switch t := err.(type) {
			case *model.AppError:
				return nil, t
			case *store.ErrInvalidInput:
				return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "option"}, err.Error(), http.StatusBadRequest)
			default:
				return nil, model.NewAppError("associateAttributeToInstance", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError)
			}
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

	default:
		return nil, model.NewAppError("associateAttributeToInstance", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "instance"}, "", http.StatusBadRequest)
	}
}
