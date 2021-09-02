package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AssignedProductAttributeByOption returns an assigned product attribute filtered using given option
func (a *ServiceAttribute) AssignedProductAttributeByOption(option *attribute.AssignedProductAttributeFilterOption) (*attribute.AssignedProductAttribute, *model.AppError) {
	assignedProductAttr, err := a.srv.Store.AssignedProductAttribute().GetWithOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AssignedProductAttributeByOption", "app.attribute.error_finding_assigned_product_attribute_by_option", err)
	}

	return assignedProductAttr, nil
}

// GetOrCreateAssignedProductAttribute get or create new instance from the given, then returns it
func (a *ServiceAttribute) GetOrCreateAssignedProductAttribute(assignedProductAttribute *attribute.AssignedProductAttribute) (*attribute.AssignedProductAttribute, *model.AppError) {
	// try get first:
	option := new(attribute.AssignedProductAttributeFilterOption)
	if assignedProductAttribute.ProductID != "" {
		option.ProductID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: assignedProductAttribute.ProductID,
			},
		}
	}
	if assignedProductAttribute.AssignmentID != "" {
		option.AssignmentID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: assignedProductAttribute.AssignmentID,
			},
		}
	}

	assignedProductAttr, appErr := a.AssignedProductAttributeByOption(option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // return immediately if error was caused by system
			return nil, appErr
		}
		// create new one
		assignedProductAttr, err := a.srv.Store.AssignedProductAttribute().Save(assignedProductAttribute)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}

			var statusCode = http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusBadRequest
			}

			return nil, model.NewAppError("GetOrCreateAssignedProductAttribute", "app.attribute.error_creating_assigned_product_attribute.app_error", nil, err.Error(), statusCode)
		}
		return assignedProductAttr, nil
	}

	return assignedProductAttr, nil
}
