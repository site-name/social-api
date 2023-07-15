package attribute

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// AssignedVariantAttributeByOption returns an assigned variant attribute filtered by given option
func (a *ServiceAttribute) AssignedVariantAttributeByOption(option *model.AssignedVariantAttributeFilterOption) (*model.AssignedVariantAttribute, *model.AppError) {
	assignedVariantAttr, err := a.srv.Store.AssignedVariantAttribute().GetWithOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AssignedVariantAttributeByOption", "app.attribute.error_finding_assigned_variant_attribute_by_option.app_error", nil, err.Error(), statusCode)
	}

	return assignedVariantAttr, nil
}

// GetOrCreateAssignedVariantAttribute get or create new assigned variant attribute with given option then returns it
func (a *ServiceAttribute) GetOrCreateAssignedVariantAttribute(assignedVariantAttr *model.AssignedVariantAttribute) (*model.AssignedVariantAttribute, *model.AppError) {
	option := new(model.AssignedVariantAttributeFilterOption)
	if assignedVariantAttr.VariantID != "" {
		option.VariantID = squirrel.Eq{model.AssignedVariantAttributeTableName + ".VariantID": assignedVariantAttr.VariantID}
	}
	if assignedVariantAttr.AssignmentID != "" {
		option.AssignmentID = squirrel.Eq{model.AssignedVariantAttributeTableName + ".AssignmentID": assignedVariantAttr.AssignmentID}
	}
	assignedVariantAttribute, appErr := a.AssignedVariantAttributeByOption(option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // returns immediately if error is system caused
		}

		// create new instance:
		assignedVariantAttribute, err := a.srv.Store.AssignedVariantAttribute().Save(assignedVariantAttr)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				return nil, appErr
			}
			statusCode := http.StatusInternalServerError
			if _, ok := err.(*store.ErrInvalidInput); ok {
				statusCode = http.StatusBadRequest
			}

			return nil, model.NewAppError("GetOrCreateAssignedVariantAttribute", "app.attribute.error_creating_assigned_variant_attribute.app_error", nil, err.Error(), statusCode)
		}

		return assignedVariantAttribute, nil
	}

	return assignedVariantAttribute, nil
}

// AssignedVariantAttributesByOption returns a list of assigned variant attributes filtered by given options
func (a *ServiceAttribute) AssignedVariantAttributesByOption(option *model.AssignedVariantAttributeFilterOption) ([]*model.AssignedVariantAttribute, *model.AppError) {
	assignedVariantAttrs, err := a.srv.Store.AssignedVariantAttribute().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("AssignedVariantAttributesByOption", "app.attribute.error_finding_assigned_variant_attributes_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return assignedVariantAttrs, nil
}
