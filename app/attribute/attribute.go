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

func (a *AppAttribute) AttributeByID(id string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.app.Srv().Store.Attribute().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeByID", AttributeMissingErrID, err)
	}

	return attr, nil
}

func (a *AppAttribute) AttributeBySlug(slug string) (*attribute.Attribute, *model.AppError) {
	attr, err := a.app.Srv().Store.Attribute().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributeBySlug", AttributeMissingErrID, err)
	}
	return attr, nil
}
