package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// DigitalContentbyOption returns 1 digital content filtered using given option
func (a *ServiceProduct) DigitalContentbyOption(option *model.DigitalContenetFilterOption) (*model.DigitalContent, *model.AppError) {
	digitalContent, err := a.srv.Store.DigitalContent().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("DigitalContent", "app.product.error_finding_digital_content_by_option,app_error", err)
	}

	return digitalContent, nil
}
