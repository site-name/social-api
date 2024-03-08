package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AssignedPageAttributeByOption(option model_helper.AssignedPageAttributeFilterOption) (*model.AssignedPageAttribute, *model_helper.AppError) {
	assignedPageAttrs, err := a.srv.Store.AssignedPageAttribute().FilterByOptions(option)

	var statusCode int
	var errMsg string

	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(assignedPageAttrs) == 0 {
		statusCode = http.StatusNotFound
	}
	if statusCode > 0 {
		return nil, model_helper.NewAppError("AssignedPageAttributeByOption", "app.attribute.assigned_page_attribute_by_options.app_error", nil, errMsg, statusCode)
	}

	return assignedPageAttrs[0], nil
}

func (a *ServiceAttribute) GetOrCreateAssignedPageAttribute(assignedPageAttribute model.AssignedPageAttribute) (*model.AssignedPageAttribute, *model_helper.AppError) {
	assignedPageAttr, appErr := a.AssignedPageAttributeByOption(model_helper.AssignedPageAttributeFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.AssignedPageAttributeWhere.PageID.EQ(assignedPageAttribute.PageID),
			model.AssignedPageAttributeWhere.AssignmentID.EQ(assignedPageAttribute.AssignmentID),
		),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// create new
		assignedPageAttr, err := a.srv.Store.AssignedPageAttribute().Upsert(assignedPageAttribute)
		if err != nil {
			if appErr, ok := err.(*model_helper.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusBadRequest
			}

			return nil, model_helper.NewAppError("GetOrCreateAssignedPageAttribute", "app.attribute.error_creating_assigned_page_attribute.app_error", nil, err.Error(), statusCode)
		}

		return assignedPageAttr, nil
	}

	return assignedPageAttr, nil
}
