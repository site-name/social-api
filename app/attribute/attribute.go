package attribute

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

type AppAttribute struct {
	app app.AppIface
}

const (
	AttributeMissingErrID = "app.attribute.attribute_missing.app_error"
)

func init() {
	app.RegisterAttributeApp(func(a app.AppIface) sub_app_iface.AttributeApp {
		return &AppAttribute{
			app: a,
		}
	})
}

// AttributeByID returns an attribute with given id
func (a *AppAttribute) AttributeByID(id string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.app.Srv().Store.Attribute().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeByID", AttributeMissingErrID, err)
	}

	return attr, nil
}

// AttributeBySlug returns an attribute with given slug
func (a *AppAttribute) AttributeBySlug(slug string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.app.Srv().Store.Attribute().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeBySlug", AttributeMissingErrID, err)
	}
	return attr, nil
}

// AttributesByOption returns a list of attributes filtered using given options
func (a *AppAttribute) AttributesByOption(option *attribute.AttributeFilterOption) ([]*attribute.Attribute, *model.AppError) {
	attributes, err := a.app.Srv().Store.Attribute().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributesByOption", "app.attribute.error_finding_attributes_by_option.app_error", err)
	}

	return attributes, nil
}
