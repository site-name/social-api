package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AssignedProductAttributeByOption(option model_helper.AssignedProductAttributeFilterOption) (*model.AssignedProductAttribute, *model_helper.AppError) {
	assignedProductAttr, err := a.srv.Store.AssignedProductAttribute().GetWithOption(option)
	if err == nil {
		return assignedProductAttr, nil
	}

	statusCode := http.StatusInternalServerError
	if _, ok := err.(*store.ErrNotFound); ok {
		statusCode = http.StatusNotFound
	}
	return nil, model_helper.NewAppError("AssignedProductAttributeByOption", "app.attribute.error_finding_assigned_product_attribute_by_option.app_error", nil, err.Error(), statusCode)
}

func (a *ServiceAttribute) GetOrCreateAssignedProductAttribute(assignedProductAttribute model.AssignedProductAttribute) (*model.AssignedProductAttribute, *model_helper.AppError) {
	assignedProductAttr, appErr := a.AssignedProductAttributeByOption(model_helper.AssignedProductAttributeFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.AssignedProductAttributeWhere.ProductID.EQ(assignedProductAttribute.ProductID),
			model.AssignedProductAttributeWhere.AssignmentID.EQ(assignedProductAttribute.AssignmentID),
		),
	})
	if appErr == nil {
		return assignedProductAttr, nil
	}

	if appErr.StatusCode == http.StatusInternalServerError { // return immediately if error was caused by system
		return nil, appErr
	}
	// create new one
	assignedProductAttr, err := a.srv.Store.AssignedProductAttribute().Save(assignedProductAttribute)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("GetOrCreateAssignedProductAttribute", "app.attribute.error_creating_assigned_product_attribute.app_error", nil, err.Error(), statusCode)
	}
	return assignedProductAttr, nil
}
