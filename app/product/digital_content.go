package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// DigitalContentbyOption returns 1 digital content filtered using given option
func (a *AppProduct) DigitalContentbyOption(option *product_and_discount.DigitalContenetFilterOption) (*product_and_discount.DigitalContent, *model.AppError) {
	digitalContent, err := a.Srv().Store.DigitalContent().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("DigitalContent", "app.product.error_finding_digital_content_by_option,app_error", err)
	}

	return digitalContent, nil
}
