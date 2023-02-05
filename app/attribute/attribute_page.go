package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// AttributePageByOption returns an attribute page filtered using given option
func (a *ServiceAttribute) AttributePageByOption(option *model.AttributePageFilterOption) (*model.AttributePage, *model.AppError) {
	attributePage, err := a.srv.Store.AttributePage().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("AttributePageByOption", "app.attribute.error_finding_attribute_page_by_option.app_error", nil, err.Error(), statusCode)
	}

	return attributePage, nil
}
