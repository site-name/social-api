package attribute

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AssignedPageAttributeByOption returns 1 assigned page attribute
func (a *ServiceAttribute) AssignedPageAttributeByOption(option *attribute.AssignedPageAttributeFilterOption) (*attribute.AssignedPageAttribute, *model.AppError) {
	assignedPageAttr, err := a.srv.Store.AssignedPageAttribute().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AssignedPageAttributeByOption", "app.attribute.error_finding_assigned_page_attribute_by_option.app_error", err)
	}

	return assignedPageAttr, nil
}

// GetOrCreateAssignedPageAttribute gets or create an assigned page attribute, then returns it
func (a *ServiceAttribute) GetOrCreateAssignedPageAttribute(assignedPageAttribute *attribute.AssignedPageAttribute) (*attribute.AssignedPageAttribute, *model.AppError) {
	option := new(attribute.AssignedPageAttributeFilterOption)
	if assignedPageAttribute.PageID != "" {
		option.PageID = squirrel.Eq{store.AssignedPageAttributeTableName + ".PageID": assignedPageAttribute.PageID}
	}
	if assignedPageAttribute.AssignmentID != "" {
		option.AssignmentID = squirrel.Eq{store.AssignedPageAttributeTableName + ".AssignmentID": assignedPageAttribute.AssignmentID}
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
