package attribute

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// AssignedPageAttributeByOption returns 1 assigned page attribute
func (a *ServiceAttribute) AssignedPageAttributeByOption(option *model.AssignedPageAttributeFilterOption) (*model.AssignedPageAttribute, *model.AppError) {
	assignedPageAttr, err := a.srv.Store.AssignedPageAttribute().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AssignedPageAttributeByOption", "app.attribute.assigned_page_attribute_by_options.app_error", nil, err.Error(), statusCode)
	}

	return assignedPageAttr, nil
}

// GetOrCreateAssignedPageAttribute gets or create an assigned page attribute, then returns it
func (a *ServiceAttribute) GetOrCreateAssignedPageAttribute(assignedPageAttribute *model.AssignedPageAttribute) (*model.AssignedPageAttribute, *model.AppError) {
	option := new(model.AssignedPageAttributeFilterOption)
	if assignedPageAttribute.PageID != "" {
		option.PageID = squirrel.Eq{model.AssignedPageAttributeTableName + ".PageID": assignedPageAttribute.PageID}
	}
	if assignedPageAttribute.AssignmentID != "" {
		option.AssignmentID = squirrel.Eq{model.AssignedPageAttributeTableName + ".AssignmentID": assignedPageAttribute.AssignmentID}
	}

	assignedPageAttr, appErr := a.AssignedPageAttributeByOption(option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// create new
		assignedPageAttr, err := a.srv.Store.AssignedPageAttribute().Save(assignedPageAttribute)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusBadRequest
			}

			return nil, model.NewAppError("GetOrCreateAssignedPageAttribute", "app.attribute.error_creating_assigned_page_attribute.app_error", nil, err.Error(), statusCode)
		}

		return assignedPageAttr, nil
	}

	return assignedPageAttr, nil
}
