/*
	NOTE: There are many methods or functions that are not implemented due to uncomplishment
	of the plugin system. Remember to implement them as soon as possible
*/
package checkout

import (
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (a *ServiceCheckout) CheckVariantInStock(ckout *checkout.Checkout, variant *product_and_discount.ProductVariant, channelSlug string, quantity int, replace, checkQuantity bool) (int, *checkout.CheckoutLine, *model.AppError) {
	// quantity param is default to 1

	checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(ckout.Token)
	if appErr != nil {
		return 0, nil, appErr
	}

	var (
		lineWithVariant *checkout.CheckoutLine = nil      // checkoutLine that has variantID of given `variantID`
		lineQuantity    int                    = 0        // quantity of lineWithVariant checkout line
		newQuantity     int                    = quantity //
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
		_, appErr = a.srv.WarehouseService().CheckStockQuantity(variant, ckout.Country, channelSlug, newQuantity)
		if appErr != nil {
			return 0, nil, appErr
		}
	}

	return newQuantity, lineWithVariant, nil
}

// AddVariantToCheckout adds a product variant to checkout
//
// `quantity` default to 1, `replace` default to false, `checkQuantity` default to true
func (a *ServiceCheckout) AddVariantToCheckout(checkoutInfo *checkout.CheckoutInfo, variant *product_and_discount.ProductVariant, quantity int, replace bool, checkQuantity bool) (*checkout.Checkout, *model.AppError) {
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

	prdChannelListings, appErr := a.srv.ProductService().
		ProductChannelListingsByOption(&product_and_discount.ProductChannelListingFilterOption{
			ChannelID: &model.StringFilter{
				StringOption: (&model.StringOption{
					Eq: checkoutInfo.Checkout.ChannelID,
				}).WithFilter(model.IsValidId),
			},
			ProductID: &model.StringFilter{
				StringOption: (&model.StringOption{
					Eq: variant.ProductID,
				}).WithFilter(model.IsValidId),
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

func (a *ServiceCheckout) CalculateCheckoutQuantity(lineInfos []*checkout.CheckoutLineInfo) (int, *model.AppError) {
	var sum int
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
func (a *ServiceCheckout) AddVariantsToCheckout(ckout *checkout.Checkout, variants []*product_and_discount.ProductVariant, quantities []int, channelSlug string, skipStockCheck, replace bool) (*checkout.Checkout, *warehouse.InsufficientStock, *model.AppError) {
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
		insfStock, appErr := a.srv.WarehouseService().CheckStockQuantityBulk(variants, countryCode, quantities, channelSlug)
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
	channelListings, appErr := a.srv.ProductService().
		ProductChannelListingsByOption(&product_and_discount.ProductChannelListingFilterOption{
			ChannelID: &model.StringFilter{
				StringOption: (&model.StringOption{
					Eq: ckout.ChannelID,
				}).WithFilter(model.IsValidId),
			},
			ProductID: &model.StringFilter{
				And: (&model.StringOption{
					In: productIDs,
				}).WithFilter(model.IsValidId),
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
func (a *ServiceCheckout) checkNewCheckoutAddress(ckout *checkout.Checkout, address *account.Address, addressType string) (bool, bool, *model.AppError) {
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
			addressesOfCheckoutUser, appErr := a.srv.AccountService().AddressesByUserId(*ckout.UserID)
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

func (a *ServiceCheckout) ChangeBillingAddressInCheckout(ckout *checkout.Checkout, address *account.Address) *model.AppError {
	changed, remove, appErr := a.checkNewCheckoutAddress(ckout, address, account.ADDRESS_TYPE_BILLING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove {
			appErr = a.srv.AccountService().DeleteAddresses([]string{*ckout.BillingAddressID})
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

// Save shipping address in checkout if changed.
//
// Remove previously saved address if not connected to any user.
func (a *ServiceCheckout) ChangeShippingAddressInCheckout(checkoutInfo *checkout.CheckoutInfo, address *account.Address, lineInfos []*checkout.CheckoutInfo) *model.AppError {
	panic("not implemented")
}

func (a *ServiceCheckout) GetDiscountedLines(checkoutLineInfos []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	var (
		discountedProducts    []*product_and_discount.Product
		discountedCategories  []*product_and_discount.Category
		discountedCollections []*product_and_discount.Collection
		appError              *model.AppError
	)

	setErr := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil {
			appError = err
		}
		a.mutex.Unlock()
	}

	// starting 3 go routines
	a.wg.Add(3)

	go func() {
		products, appErr := a.srv.ProductService().ProductsByVoucherID(voucher.Id)
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedProducts = products
		}

		a.wg.Done()
	}()

	go func() {
		categories, appErr := a.srv.ProductService().CategoriesByOption(&product_and_discount.CategoryFilterOption{
			VoucherIDs: []string{voucher.Id},
		})
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedCategories = categories
		}

		a.wg.Done()
	}()

	go func() {
		collections, appErr := a.srv.ProductService().CollectionsByVoucherID(voucher.Id)
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedCollections = collections
		}

		a.wg.Done()
	}()

	a.wg.Done()

	if appError != nil {
		appError.Where = "GetDiscountedLines"
		return nil, appError
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
//
// `withLock` default to false
func (a *ServiceCheckout) GetVoucherForCheckout(checkoutInfo *checkout.CheckoutInfo, withLock bool) (*product_and_discount.Voucher, *model.AppError) {

	now := model.NewTime(time.Now()) // NOTE: not sure to use UTC or system time

	if checkoutInfo.Checkout.VoucherCode != nil {
		// finds vouchers that are active in a channel
		activeInChannelVouchers, appErr := a.srv.DiscountService().VouchersByOption(&product_and_discount.VoucherFilterOption{
			UsageLimit: &model.NumberFilter{
				Or: &model.NumberOption{
					NULL: model.NewBool(true),
					ExtraExpr: []squirrel.Sqlizer{
						squirrel.Expr("?.UsageLimit > ?.Used", store.VoucherTableName, store.VoucherTableName),
					},
				},
			},
			EndDate: &model.TimeFilter{
				Or: &model.TimeOption{
					NULL: model.NewBool(true),
					GtE:  now,
				},
			},
			StartDate: &model.TimeFilter{
				TimeOption: &model.TimeOption{
					LtE: now,
				},
			},
			ChannelListingSlug: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkoutInfo.Channel.Slug,
				},
			},
			ChannelListingActive: model.NewBool(true),
			WithLook:             withLock, // this add `FOR UPDATE` to SQL query
		})

		if appErr != nil || len(activeInChannelVouchers) == 0 {
			appErr.Where = "GetVoucherForCheckout"
			return nil, appErr
		}

		// find voucher with code
		for _, voucher := range activeInChannelVouchers {
			if voucher.Code == *checkoutInfo.Checkout.VoucherCode {
				return voucher, nil
			}
		}
	}

	return nil, nil
}

// RemovePromoCodeFromCheckout Remove gift card or voucher data from checkout.
func (a *ServiceCheckout) RemovePromoCodeFromCheckout(checkoutInfo *checkout.CheckoutInfo, promoCode string) *model.AppError {
	// check if promoCode is voucher:
	promoCodeIsVoucher, appErr := a.srv.DiscountService().PromoCodeIsVoucher(promoCode)
	if appErr != nil { // this error is system error
		return appErr
	}
	if promoCodeIsVoucher {
		return a.RemoveVoucherCodeFromCheckout(checkoutInfo, promoCode)
	}

	// check promoCode is giftcard
	promoCodeIsGiftCard, appErr := a.srv.GiftcardService().PromoCodeIsGiftCard(promoCode)
	if appErr != nil {
		return appErr
	}
	if promoCodeIsGiftCard {
		return a.srv.GiftcardService().RemoveGiftcardCodeFromCheckout(&checkoutInfo.Checkout, promoCode)
	}

	return nil
}

// RemoveVoucherCodeFromCheckout Remove voucher data from checkout by code.
func (a *ServiceCheckout) RemoveVoucherCodeFromCheckout(checkoutInfo *checkout.CheckoutInfo, voucherCode string) *model.AppError {
	existingVoucher, appErr := a.GetVoucherForCheckout(checkoutInfo, false)
	if appErr != nil {
		return appErr
	}
	if existingVoucher != nil && existingVoucher.Code == voucherCode {
		return a.RemoveVoucherFromCheckout(&checkoutInfo.Checkout)
	}

	return nil
}

// RemoveVoucherFromCheckout removes voucher data from checkout
func (a *ServiceCheckout) RemoveVoucherFromCheckout(ckout *checkout.Checkout) *model.AppError {
	if ckout == nil {
		return nil
	}

	ckout.VoucherCode = nil
	ckout.DiscountName = nil
	ckout.TranslatedDiscountName = nil
	ckout.DiscountAmount = &decimal.Zero

	_, appErr := a.UpsertCheckout(ckout)

	return appErr
}

// GetValidShippingMethodsForCheckout finds all valid shipping methods for given checkout
func (a *ServiceCheckout) GetValidShippingMethodsForCheckout(checkoutInfo *checkout.CheckoutInfo, lineInfos []*checkout.CheckoutLineInfo, subTotal *goprices.TaxedMoney, countryCode string) ([]*shipping.ShippingMethod, *model.AppError) {
	var productIDs []string
	for _, line := range lineInfos {
		productIDs = append(productIDs, line.Product.Id)
	}

	// check if any product in given lineInfos requires shipping:
	requireShipping, appErr := a.srv.ProductService().ProductsRequireShipping(productIDs)
	if appErr != nil || requireShipping {
		return nil, appErr
	}

	// check if checkoutInfo
	if checkoutInfo.ShippingAddress == nil {
		return nil, nil
	}

	return a.srv.ShippingService().ApplicableShippingMethodsForCheckout(
		&checkoutInfo.Checkout,
		checkoutInfo.Checkout.ChannelID,
		subTotal.Gross,
		countryCode,
		lineInfos,
	)
}

// IsValidShippingMethod Check if shipping method is valid and remove (if not).
func (a *ServiceCheckout) IsValidShippingMethod(checkoutInfo *checkout.CheckoutInfo) (bool, *model.AppError) {
	if checkoutInfo.ShippingMethod == nil || checkoutInfo.ShippingAddress == nil {
		return false, nil
	}

	var validShippingMethodIDs []string
	if len(checkoutInfo.ValidShippingMethods) != 0 {
		for _, method := range checkoutInfo.ValidShippingMethods {
			validShippingMethodIDs = append(validShippingMethodIDs, method.Id)
		}
	}

	if len(validShippingMethodIDs) == 0 || !util.StringInSlice(checkoutInfo.ShippingMethod.Id, validShippingMethodIDs) {
		appErr := a.ClearShippingMethod(checkoutInfo)
		return false, appErr
	}

	return true, nil
}

func (a *ServiceCheckout) ClearShippingMethod(checkoutInfo *checkout.CheckoutInfo) *model.AppError {
	ckout := checkoutInfo.Checkout
	ckout.ShippingMethodID = nil

	appErr := a.UpdateCheckoutInfoShippingMethod(checkoutInfo, nil)
	if appErr != nil {
		return nil
	}

	_, appErr = a.UpsertCheckout(&ckout)
	return appErr
}

// CancelActivePayments set all active payments belong to given checkout
func (a *ServiceCheckout) CancelActivePayments(ckout *checkout.Checkout) *model.AppError {
	err := a.srv.Store.Payment().CancelActivePaymentsOfCheckout(ckout.Token)
	if err != nil {
		return model.NewAppError("CancelActivePayments", "app.checkout.cancel_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) ValidateVariantsInCheckoutLines(lines []*checkout.CheckoutLineInfo) *model.AppError {

	var notAvailableVariants []string
	for _, line := range lines {
		if line.ChannelListing == nil || line.ChannelListing.Price == nil {
			notAvailableVariants = append(notAvailableVariants, line.Variant.Id)
		}
	}

	if len(notAvailableVariants) > 0 {
		// return error indicate there are some product variants that have no channel listing or channel listing price is null
		return model.NewAppError("ValidateVariantsInCheckoutLines", "app.checkout.cannot_add_lines_with_unavailable_variants.app_error", map[string]interface{}{"variantIDs": strings.Join(notAvailableVariants, ", ")}, "", http.StatusNotAcceptable)
	}

	return nil
}
