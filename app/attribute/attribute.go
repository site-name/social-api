/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package attribute

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type ServiceAttribute struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Attribute = &ServiceAttribute{s}
		return nil
	})
}

func (a *ServiceAttribute) AttributeByOption(option *model.AttributeFilterOption) (*model.Attribute, *model.AppError) {
	attributes, appErr := a.AttributesByOption(option)
	if appErr != nil {
		return nil, appErr
	}
	if len(attributes) == 0 {
		return nil, model.NewAppError("AttributeByOption", "app.attribute.error_finding_attribute_by_option.app_error", nil, "", http.StatusNotFound)
	}

	return attributes[0], nil
}

// AttributesByOption returns a list of attributes filtered using given options
func (a *ServiceAttribute) AttributesByOption(option *model.AttributeFilterOption) ([]*model.Attribute, *model.AppError) {
	attributes, err := a.srv.Store.Attribute().FilterbyOption(option)
	if err != nil {
		return nil, model.NewAppError("AttributesByOption", "app.attribute.attributes_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return attributes, nil
}

// UpsertAttribute inserts or updates given attribute and returns it
func (s *ServiceAttribute) UpsertAttribute(attr *model.Attribute) (*model.Attribute, *model.AppError) {
	attr, err := s.srv.Store.Attribute().Upsert(attr)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("UpsertAttribute", "app.attribute.error_upserting_attribute.app_error", nil, err.Error(), statusCode)
	}

	return attr, nil
}

func (s *ServiceAttribute) DeleteAttributes(ids ...string) (int64, *model.AppError) {
	numDeleted, err := s.srv.Store.Attribute().Delete(ids...)
	if err != nil {
		return 0, model.NewAppError("DeleteAttribute", "app.attribute.error_deleting_attributes.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return numDeleted, nil
}
