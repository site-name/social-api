package checkout

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
)

func (a *AppCheckout) CheckVariantInStock(variant *product_and_discount.ProductVariant, ckout *checkout.Checkout, channelSlug string, quantity *uint, replace, checkQuantity bool) (uint, *checkout.CheckoutLine, *model.AppError) {
	// quantity param is default to 1
	if quantity == nil {
		quantity = model.NewUint(1)
	}

	checkoutLines, appErr := a.CheckoutLinesByCheckoutID(ckout.Token)
	if appErr != nil {
		return 0, nil, appErr
	}

	var (
		lineWithVariant *checkout.CheckoutLine             // checkoutLine that has variantID of given `variantID`
		lineQuantity    uint                               // quantity of lineWithVariant checkout line
		newQuantity     uint                   = *quantity //
	)

	for _, checkoutLine := range checkoutLines {
		if checkoutLine.VariantID == variant.Id {
			lineWithVariant = checkoutLine
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
		return 0, nil, model.NewAppError(
			"CheckVariantInStock",
			"app.checkout.quantity_invalid",
			map[string]interface{}{
				"Quantity":    *quantity,
				"NewQuantity": newQuantity,
			},
			"", http.StatusBadRequest,
		)
	}

	if newQuantity > 0 && checkQuantity {
		// NOTE: have a look at ckout.Country below
		_, appErr := a.WarehouseApp().CheckStockQuantity(variant, ckout.Country, channelSlug, newQuantity)
		if appErr != nil {
			return 0, nil, appErr
		}
	}

	return newQuantity, lineWithVariant, nil
}

// AddVariantToCheckout adds a product variant to checkout
//
// `quantity` default to 1, `replace` default to false, `checkQuantity` default to true
func (a *AppCheckout) AddVariantToCheckout(checkoutInfo *checkout.CheckoutInfo, variant *product_and_discount.ProductVariant, quantity int, replace bool, checkQuantity bool) (*checkout.Checkout, *model.AppError) {
	panic("not implt")
}
