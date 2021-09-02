/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package attribute

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type ServiceAttribute struct {
	srv *app.Server
}

const (
	AttributeMissingErrID = "app.attribute.attribute_missing.app_error"
)

func init() {
	app.RegisterAttributeApp(func(s *app.Server) (sub_app_iface.AttributeService, error) {
		return &ServiceAttribute{
			srv: s,
		}, nil
	})
}

// AttributeByID returns an attribute with given id
func (a *ServiceAttribute) AttributeByID(id string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.srv.Store.Attribute().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeByID", AttributeMissingErrID, err)
	}

	return attr, nil
}

// AttributeBySlug returns an attribute with given slug
func (a *ServiceAttribute) AttributeBySlug(slug string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.srv.Store.Attribute().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeBySlug", AttributeMissingErrID, err)
	}
	return attr, nil
}

// AttributesByOption returns a list of attributes filtered using given options
func (a *ServiceAttribute) AttributesByOption(option *attribute.AttributeFilterOption) ([]*attribute.Attribute, *model.AppError) {
	attributes, err := a.srv.Store.Attribute().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributesByOption", "app.attribute.error_finding_attributes_by_option.app_error", err)
	}

	return attributes, nil
}
