package attribute

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceAttribute) AttributePageByOption(option model_helper.AttributePageFilterOption) (*model.AttributePage, *model_helper.AppError) {
	attributePage, err := a.srv.Store.AttributePage().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("AttributePageByOption", "app.attribute.error_finding_attribute_page_by_option.app_error", nil, err.Error(), statusCode)
	}

	return attributePage, nil
}
