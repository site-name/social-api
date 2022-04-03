/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package attribute

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type ServiceAttribute struct {
	srv *app.Server
}

func init() {
	app.RegisterAttributeService(func(s *app.Server) (sub_app_iface.AttributeService, error) {
		return &ServiceAttribute{
			srv: s,
		}, nil
	})
}

// AttributeByID returns an attribute with given id
func (a *ServiceAttribute) AttributeByID(id string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.srv.Store.Attribute().Get(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AttributeByID", "app.attriute.error_finding_attribute_by_id.app_error", nil, err.Error(), statusCode)
	}

	return attr, nil
}

// AttributeBySlug returns an attribute with given slug
func (a *ServiceAttribute) AttributeBySlug(slug string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.srv.Store.Attribute().GetBySlug(slug)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AttributeBySlug", "app.attribute.error_finding_attribute_by_slug.app_error", nil, err.Error(), statusCode)
	}

	return attr, nil
}

// AttributesByOption returns a list of attributes filtered using given options
func (a *ServiceAttribute) AttributesByOption(option *attribute.AttributeFilterOption) ([]*attribute.Attribute, *model.AppError) {
	attributes, err := a.srv.Store.Attribute().FilterbyOption(option)

	var (
		statusCode int = 0
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(attributes) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("AttributesByOption", "app.attribute.error_finding_attributes_by_options.app_error", nil, errMsg, statusCode)
	}

	return attributes, nil
}

// UpsertAttribute inserts or updates given attribute and returns it
func (s *ServiceAttribute) UpsertAttribute(attr *attribute.Attribute) (*attribute.Attribute, *model.AppError) {
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
