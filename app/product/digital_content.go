package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// DigitalContentbyOption returns 1 digital content filtered using given option
func (a *ServiceProduct) DigitalContentbyOption(option *model.DigitalContenetFilterOption) (*model.DigitalContent, *model.AppError) {
	digitalContent, err := a.srv.Store.DigitalContent().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("DigitalContentbyOption", "app.product.error_finding_digital_content_by_option,app_error", nil, err.Error(), statusCode)
	}

	return digitalContent, nil
}
