package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AssignedVariantAttributeByOption returns an assigned variant attribute filtered by given option
func (a *AppAttribute) AssignedVariantAttributeByOption(option *attribute.AssignedVariantAttributeFilterOption) (*attribute.AssignedVariantAttribute, *model.AppError) {
	assignedVariantAttr, err := a.app.Srv().Store.AssignedVariantAttribute().GetWithOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AssignedVariantAttributeByOption", "app.attribute.error_finding_assigned_variant_attribute_by_option.app_error", err)
	}

	return assignedVariantAttr, nil
}

// GetOrCreateAssignedVariantAttribute get or create new assigned variant attribute with given option then returns it
func (a *AppAttribute) GetOrCreateAssignedVariantAttribute(assignedVariantAttr *attribute.AssignedVariantAttribute) (*attribute.AssignedVariantAttribute, *model.AppError) {

	option := new(attribute.AssignedVariantAttributeFilterOption)
	if assignedVariantAttr.VariantID != "" {
		option.VariantID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: assignedVariantAttr.VariantID,
			},
		}
	}
	if assignedVariantAttr.AssignmentID != "" {
		option.AssignmentID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: assignedVariantAttr.AssignmentID,
			},
		}
	}
	assignedVariantAttribute, appErr := a.AssignedVariantAttributeByOption(option)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // returns immediately if error is system caused
		}

		// create new instance:
		assignedVariantAttribute, err := a.app.Srv().Store.AssignedVariantAttribute().Save(assignedVariantAttr)
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
func (a *AppAttribute) AssignedVariantAttributesByOption(option *attribute.AssignedVariantAttributeFilterOption) ([]*attribute.AssignedVariantAttribute, error) {
	assignedVariantAttrs, err := a.app.Srv().Store.AssignedVariantAttribute().FilterByOption(option)

	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(assignedVariantAttrs) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AssignedVariantAttributesByOption", "app.attribute.error_finding_assigned_variant_attributes_by_option.app_error", nil, errMsg, statusCode)
	}

	return assignedVariantAttrs, nil
}
