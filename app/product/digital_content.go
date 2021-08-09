package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// DigitalContentByProductVariantID finds and returns 1 digital content that is related to given product variant
func (a *AppProduct) DigitalContentByProductVariantID(variantID string) (*product_and_discount.DigitalContent, *model.AppError) {
	digitalContent, err := a.Srv().Store.DigitalContent().GetByProductVariantID(variantID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("", "app.product.error_finding_digital_content_with_varant_id.app_error", err)
	}

	return digitalContent, nil
}
