package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
)

func (a *AppCheckout) CheckVariantInStock(variant *product_and_discount.ProductVariant, checkoutID, channelSlug string, quantity *uint, replace, checkQuantity bool) (interface{}, *checkout.CheckoutLine, *model.AppError) {
	// quantity param is default to 1
	if quantity == nil {
		quantity = model.NewUint(1)
	}

	lines, appErr := a.CheckoutLinesByCheckoutID(checkoutID)
	if appErr != nil {
		return nil, nil, appErr
	}

	var (
		lineWithVariant *checkout.CheckoutLine             // checkoutLine that has variantID of given `variantID`
		lineQuantity    uint                               // quantity of lineWithVariant checkout line
		newQuantity     uint                   = *quantity //
	)

	for _, line := range lines {
		if line.VariantID == variant.Id {
			lineWithVariant = line
			break
		}
	}

	if lineWithVariant != nil {
		lineQuantity = lineWithVariant.Quantity
	}

	if !replace {
		newQuantity = *quantity + lineQuantity
	}

	if newQuantity < 0 {
		return nil, nil, model.NewAppError("CheckVariantInStock", "app.checkout.quantity_not_negative",
			map[string]interface{}{"Quantity": *quantity, "NewQuantity": newQuantity}, "", http.StatusBadRequest)
	}

	if newQuantity > 0 && checkQuantity {
		a.WarehouseApp().CheckStockQuantity()
	}

	return newQuantity, lineWithVariant, nil
}

func (a *AppCheckout) AddVariantToCheckout() {

}
