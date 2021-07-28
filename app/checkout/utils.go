/*
	NOTE: There are many methods or functions that are not implemented due to uncomplishment
	of the plugin system. Remember to implement them as soon as possible
*/
package checkout

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
)

func (a *AppCheckout) CheckVariantInStock(ckout *checkout.Checkout, variant *product_and_discount.ProductVariant, channelSlug string, quantity uint, replace, checkQuantity bool) (uint, *checkout.CheckoutLine, *model.AppError) {
	// quantity param is default to 1

	checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(ckout.Token)
	if appErr != nil {
		return 0, nil, appErr
	}

	var (
		lineWithVariant *checkout.CheckoutLine = nil      // checkoutLine that has variantID of given `variantID`
		lineQuantity    uint                   = 0        // quantity of lineWithVariant checkout line
		newQuantity     uint                   = quantity //
	)
	if !replace {
		newQuantity = quantity + lineQuantity
	}

	for _, checkoutLine := range checkoutLines {
		if checkoutLine.VariantID == variant.Id {
			lineWithVariant = checkoutLine
			break
		}
	}

	if lineWithVariant != nil {
		lineQuantity = lineWithVariant.Quantity
	}

	if newQuantity < 0 {
		return 0, nil, model.NewAppError(
			"CheckVariantInStock",
			"app.checkout.quantity_invalid",
			map[string]interface{}{
				"Quantity":    quantity,
				"NewQuantity": newQuantity,
			},
			"", http.StatusBadRequest,
		)
	}

	if newQuantity > 0 && checkQuantity {
		// NOTE: have a look at ckout.Country below
		_, appErr = a.app.WarehouseApp().CheckStockQuantity(variant, ckout.Country, channelSlug, newQuantity)
		if appErr != nil {
			return 0, nil, appErr
		}
	}

	return newQuantity, lineWithVariant, nil
}

// AddVariantToCheckout adds a product variant to checkout
//
// `quantity` default to 1, `replace` default to false, `checkQuantity` default to true
func (a *AppCheckout) AddVariantToCheckout(checkoutInfo *checkout.CheckoutInfo, variant *product_and_discount.ProductVariant, quantity uint, replace bool, checkQuantity bool) (*checkout.Checkout, *model.AppError) {
	// validate arguments
	var invalidArgs string
	if checkoutInfo == nil {
		invalidArgs = "checkoutInfo"
	}
	if variant == nil {
		invalidArgs += ", variant"
	}
	if invalidArgs != "" {
		return nil, model.NewAppError("AddVariantToCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invalidArgs}, "", http.StatusBadRequest)
	}

	prdChannelListings, appErr := a.app.ProductApp().ProductChannelListingsByOption(&product_and_discount.ProductChannelListingFilterOption{
		ChannelID: &product_and_discount.StringFilter{
			Eq: checkoutInfo.Checkout.ChannelID,
		},
		ProductID: &product_and_discount.StringFilter{
			Eq: variant.ProductID,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	if len(prdChannelListings) == 0 || !prdChannelListings[0].IsPublished {
		return nil, model.NewAppError("AddVariantToCheckout", app.ProductNotPublishedAppErrID, nil, "Please publish the product first.", http.StatusNotAcceptable)
	}

	newQuantity, line, appErr := a.CheckVariantInStock(&checkoutInfo.Checkout, variant, checkoutInfo.Channel.Slug, quantity, replace, checkQuantity)
	if appErr != nil {
		return nil, appErr
	}

	if line == nil {
		checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(checkoutInfo.Checkout.Token)
		if appErr != nil && appErr.StatusCode != http.StatusNotFound { // ignore not found error
			return nil, appErr
		}
		line = checkoutLines[0]
	}

	if newQuantity == 0 {
		if line != nil {
			if appErr = a.DeleteCheckoutLines([]string{line.Id}); appErr != nil {
				return nil, appErr
			}
		}
	} else if line == nil {
		if _, appErr = a.UpsertCheckoutLine(&checkout.CheckoutLine{
			CheckoutID: checkoutInfo.Checkout.Token,
			VariantID:  variant.Id,
			Quantity:   newQuantity,
		}); appErr != nil {
			return nil, appErr
		}
	} else if newQuantity > 0 {
		line.Quantity = newQuantity
		if _, appErr = a.UpsertCheckoutLine(line); appErr != nil {
			return nil, appErr
		}
	}

	return &checkoutInfo.Checkout, nil
}

func (a *AppCheckout) CalculateCheckoutQuantity(lineInfos []*checkout.CheckoutLineInfo) (uint, *model.AppError) {
	var sum uint
	for _, info := range lineInfos {
		sum += info.Line.Quantity
	}

	return sum, nil
}

// AddVariantsToCheckout Add variants to checkout.
//
// If a variant is not placed in checkout, a new checkout line will be created.
// If quantity is set to 0, checkout line will be deleted.
// Otherwise, quantity will be added or replaced (if replace argument is True).
//  skipStockCheck and replace are default to false
func (a *AppCheckout) AddVariantsToCheckout(ckout *checkout.Checkout, variants []*product_and_discount.ProductVariant, quantities []uint, channelSlug string, skipStockCheck, replace bool) (*checkout.Checkout, *warehouse.InsufficientStock, *model.AppError) {
	// validate input arguments:
	var invlArgs string
	if ckout == nil {
		invlArgs = "ckout"
	}
	if len(variants) == 0 {
		invlArgs += ", variants"
	}
	if invlArgs != "" {
		return nil, nil, model.NewAppError("AddVariantsToCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invlArgs}, "", http.StatusBadRequest)
	}

	// check quantities
	countryCode, appErr := a.CheckoutCountry(ckout)
	if appErr != nil {
		return nil, nil, appErr
	}
	if !skipStockCheck {
		insfStock, appErr := a.app.WarehouseApp().CheckStockQuantityBulk(variants, countryCode, quantities, channelSlug)
		if appErr != nil {
			return nil, nil, appErr
		}
		if insfStock != nil && len(insfStock.Items) > 0 {
			return nil, insfStock, nil
		}
	}

	productIDs := make([]string, len(variants))
	for i, variant := range variants {
		productIDs[i] = variant.ProductID
	}
	channelListings, appErr := a.app.ProductApp().
		ProductChannelListingsByOption(&product_and_discount.ProductChannelListingFilterOption{
			ChannelID: &product_and_discount.StringFilter{
				Eq: ckout.ChannelID,
			},
			ProductID: &product_and_discount.StringFilter{
				In: productIDs,
			},
		})
	if appErr != nil {
		return nil, nil, appErr
	}

	listingMap := make(map[string]*product_and_discount.ProductChannelListing)
	for _, listing := range channelListings {
		listingMap[listing.ProductID] = listing
	}

	for _, productID := range productIDs {
		if listingMap[productID] == nil || !listingMap[productID].IsPublished {
			return nil, nil, model.NewAppError("AddVariantsToCheckout", app.ProductNotPublishedAppErrID, nil, "", http.StatusNotAcceptable)
		}
	}

	linesOfCheckout, appErr := a.CheckoutLinesByCheckoutToken(ckout.Token)
	if appErr != nil {
		return nil, nil, appErr
	}

	variantIDsInLines := make(map[string]*checkout.CheckoutLine)
	for _, line := range linesOfCheckout {
		variantIDsInLines[line.VariantID] = line
	}

	var (
		toCreateCheckoutLines   = []*checkout.CheckoutLine{}
		toUpdateCheckoutLines   = []*checkout.CheckoutLine{}
		toDeleteCheckoutLineIDs = []string{}
	)
	// use Min() since two slices may not have same length
	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		if checkoutLine, exist := variantIDsInLines[variants[i].Id]; exist {
			if quantities[i] > 0 {
				if replace {
					checkoutLine.Quantity = quantities[i]
				} else {
					checkoutLine.Quantity += quantities[i]
				}
				toUpdateCheckoutLines = append(toUpdateCheckoutLines, checkoutLine)
			} else {
				toDeleteCheckoutLineIDs = append(toDeleteCheckoutLineIDs, checkoutLine.Id)
			}
		} else {
			toCreateCheckoutLines = append(toCreateCheckoutLines, &checkout.CheckoutLine{
				CheckoutID: ckout.Token,
				VariantID:  variants[i].Id,
				Quantity:   quantities[i],
			})
		}
	}

	if len(toDeleteCheckoutLineIDs) > 0 {
		appErr = a.DeleteCheckoutLines(toDeleteCheckoutLineIDs)
		if appErr != nil {
			return nil, nil, appErr
		}
	}
	if len(toUpdateCheckoutLines) > 0 {
		appErr = a.BulkUpdateCheckoutLines(toUpdateCheckoutLines)
		if appErr != nil {
			return nil, nil, appErr
		}
	}
	if len(toCreateCheckoutLines) > 0 {
		_, appErr = a.BulkCreateCheckoutLines(toCreateCheckoutLines)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	return nil, nil, nil
}

// checkNewCheckoutAddress Check if and address in checkout has changed and if to remove old one
func (a *AppCheckout) checkNewCheckoutAddress(ckout *checkout.Checkout, address *account.Address, addressType string) (bool, bool, *model.AppError) {
	// validate if non-nill checkout was provided
	var invalidArguments string
	if ckout == nil {
		invalidArguments = "ckout"
	}
	if invalidArguments != "" {
		return false, false, model.NewAppError("checkNewCheckoutAddress", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invalidArguments}, "", http.StatusBadRequest)
	}

	oldAddressId := ckout.ShippingAddressID
	if addressType == account.ADDRESS_TYPE_BILLING {
		oldAddressId = ckout.BillingAddressID
	}

	hasAddressChanged := (address == nil || oldAddressId != nil) ||
		(address != nil && oldAddressId == nil) ||
		(address != nil && oldAddressId != nil && address.Id != *oldAddressId)

	if oldAddressId == nil {
		return hasAddressChanged, false, nil
	} else {
		if ckout.UserID == nil {
			return hasAddressChanged, hasAddressChanged, nil
		} else {
			var oldAddressNOTbelongToCheckoutUser bool
			addressesOfCheckoutUser, appErr := a.app.AccountApp().AddressesByUserId(*ckout.UserID)
			if appErr != nil {
				if appErr.StatusCode == http.StatusNotFound { // user owns 0 address
					oldAddressNOTbelongToCheckoutUser = true
				}
				return false, false, appErr // must returns since this is system's error
			} else {
				oldAddressNOTbelongToCheckoutUser = true
				for _, addr := range addressesOfCheckoutUser {
					if *oldAddressId == addr.Id {
						oldAddressNOTbelongToCheckoutUser = false
						break
					}
				}
			}
			return hasAddressChanged, hasAddressChanged && oldAddressNOTbelongToCheckoutUser, nil
		}
	}
}

func (a *AppCheckout) ChangeBillingAddressInCheckout(ckout *checkout.Checkout, address *account.Address) *model.AppError {
	changed, remove, appErr := a.checkNewCheckoutAddress(ckout, address, account.ADDRESS_TYPE_BILLING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove {
			appErr = a.app.AccountApp().DeleteAddresses([]string{*ckout.BillingAddressID})
			if appErr != nil {
				return appErr
			}
		}
		ckout.BillingAddressID = &address.Id
		_, appErr = a.UpsertCheckout(ckout)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// func (a *AppCheckout) ChangeShippingAddressInCheckout(checkoutInfo *checkout.CheckoutInfo, address *account.Address, lineInfos []*checkout.CheckoutInfo) *model.AppError {

// }

func (a *AppCheckout) GetDiscountedLines(checkoutLineInfos []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	discountedProducts, appErr := a.app.ProductApp().ProductsByVoucherID(voucher.Id)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // system's error, returns immediately
			return nil, appErr
		}
	}
	discountedCategories, appErr := a.app.ProductApp().CategoriesByVoucherID(voucher.Id)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // system's error, returns immediately
			return nil, appErr
		}
	}
	discountedCollections, appErr := a.app.ProductApp().CollectionsByVoucherID(voucher.Id)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // system's error, returns immediately
			return nil, appErr
		}
	}

	var (
		discountedProductIDs    []string
		discountedCategoryIDs   []string
		discountedCollectionIDs []string
	)
	for _, prd := range discountedProducts {
		discountedProductIDs = append(discountedProductIDs, prd.Id)
	}

	// filter duplicates from discountedCategories:
	meetMap := map[string]bool{}
	for _, category := range discountedCategories {
		if _, ok := meetMap[category.Id]; !ok {
			discountedCategoryIDs = append(discountedCategoryIDs, category.Id)
			meetMap[category.Id] = true
		}
	}
	// filter duplicates from discountedCollections:
	// NOTE: reuse meetMap here since UUIDs are unique
	for _, collection := range discountedCollections {
		if _, ok := meetMap[collection.Id]; !ok {
			discountedCollectionIDs = append(discountedCollectionIDs, collection.Id)
			meetMap[collection.Id] = true
		}
	}

	var discountedLines []*checkout.CheckoutLineInfo
	if len(discountedProductIDs) > 0 || len(discountedCategoryIDs) > 0 || len(discountedCollectionIDs) > 0 {
		for _, lineInfo := range checkoutLineInfos {

			var lineInfoCollections_have_common_with_discountedCollections bool
			for _, collection := range lineInfo.Collections {
				if yes, exist := meetMap[collection.Id]; yes && exist {
					lineInfoCollections_have_common_with_discountedCollections = true
					break
				}
			}

			if lineInfo.Variant != nil && (util.StringInSlice(lineInfo.Product.Id, discountedProductIDs) ||
				(lineInfo.Product.CategoryID != nil && util.StringInSlice(*lineInfo.Product.CategoryID, discountedCategoryIDs)) ||
				lineInfoCollections_have_common_with_discountedCollections) {
				discountedLines = append(discountedLines, lineInfo)
			}
		}
		return discountedLines, nil
	} else {
		// If there's no discounted products, collections or categories,
		// it means that all products are discounted
		return checkoutLineInfos, nil
	}
}

// GetVoucherForCheckout returns voucher with voucher code saved in checkout if active or None
func (a *AppCheckout) GetVoucherForCheckout(checkoutInfo *checkout.CheckoutInfo, withLock bool) (*product_and_discount.Voucher, *model.AppError) {
	if checkoutInfo.Checkout.VoucherCode != nil {

	}
}
