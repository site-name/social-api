package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/store"
)

// AttributePageByOption returns an attribute page filtered using given option
func (a *ServiceAttribute) AttributePageByOption(option *attribute.AttributePageFilterOption) (*attribute.AttributePage, *model.AppError) {
	attributePage, err := a.srv.Store.AttributePage().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("AttributePageByOption", "app.attribute.error_finding_attribute_page_by_option.app_error", err)
	}

	return attributePage, nil
}
