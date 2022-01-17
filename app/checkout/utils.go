/*
	NOTE: There are many methods or functions that are not implemented due to uncomplishment
	of the plugin system. Remember to implement them as soon as possible
*/
package checkout

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// CheckVariantInStock
func (a *ServiceCheckout) CheckVariantInStock(checkOut *checkout.Checkout, variant *product_and_discount.ProductVariant, channelSlug string, quantity int, replace, checkQuantity bool) (int, *checkout.CheckoutLine, *exception.InsufficientStock, *model.AppError) {
	// quantity param is default to 1

	checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(checkOut.Token)
	if appErr != nil {
		return 0, nil, nil, appErr
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
		return 0, nil, nil, model.NewAppError(
			"CheckVariantInStock",
			"app.checkout.quantity_invalid.app_error",
			map[string]interface{}{
				"Quantity":    quantity,
				"NewQuantity": newQuantity,
			},
			"", http.StatusBadRequest,
		)
	}

	if newQuantity > 0 && checkQuantity {
		insufficientStockErr, appErr := a.srv.WarehouseService().CheckStockAndPreorderQuantity(variant, checkOut.Country, channelSlug, newQuantity)
		if insufficientStockErr != nil || appErr != nil {
			return 0, nil, insufficientStockErr, appErr
		}
	}

	return newQuantity, lineWithVariant, nil, nil
}

// AddVariantToCheckout adds a product variant to checkout
//
// `quantity` default to 1, `replace` default to false, `checkQuantity` default to true
func (a *ServiceCheckout) AddVariantToCheckout(checkoutInfo *checkout.CheckoutInfo, variant *product_and_discount.ProductVariant, quantity int, replace bool, checkQuantity bool) (*checkout.Checkout, *exception.InsufficientStock, *model.AppError) {
	// validate arguments
	var invalidArgs string
	if checkoutInfo == nil {
		invalidArgs = "checkoutInfo"
	}
	if variant == nil {
		invalidArgs += ", variant"
	}
	if invalidArgs != "" {
		return nil, nil, model.NewAppError("AddVariantToCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invalidArgs}, "", http.StatusBadRequest)
	}

	checkOut := checkoutInfo.Checkout
	productChannelListings, appErr := a.srv.ProductService().ProductChannelListingsByOption(&product_and_discount.ProductChannelListingFilterOption{
		ChannelID: &model.StringFilter{
			StringOption: (&model.StringOption{
				Eq: checkOut.ChannelID,
			}),
		},
		ProductID: &model.StringFilter{
			StringOption: (&model.StringOption{
				Eq: variant.ProductID,
			}),
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	if len(productChannelListings) == 0 || !productChannelListings[0].IsPublished {
		return nil, nil, model.NewAppError("AddVariantToCheckout", app.ProductNotPublishedAppErrID, nil, "Please publish the product first.", http.StatusNotAcceptable)
	}

	newQuantity, line, insufficientErr, appErr := a.CheckVariantInStock(&checkoutInfo.Checkout, variant, checkoutInfo.Channel.Slug, quantity, replace, checkQuantity)
	if appErr != nil || insufficientErr != nil {
		return nil, insufficientErr, appErr
	}

	if line == nil {
		checkoutLines, appErr := a.CheckoutLinesByOption(&checkout.CheckoutLineFilterOption{
			CheckoutID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkOut.Token,
				},
			},
			VariantID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: variant.Id,
				},
			},
		})
		if appErr != nil && appErr.StatusCode != http.StatusNotFound { // ignore not found error
			return nil, nil, appErr
		}
		line = checkoutLines[0]
	}

	if newQuantity == 0 {
		if line != nil {
			if appErr = a.DeleteCheckoutLines(nil, []string{line.Id}); appErr != nil {
				return nil, nil, appErr
			}
		}
	} else if line == nil {
		if _, appErr = a.UpsertCheckoutLine(&checkout.CheckoutLine{
			CheckoutID: checkoutInfo.Checkout.Token,
			VariantID:  variant.Id,
			Quantity:   newQuantity,
		}); appErr != nil {
			return nil, nil, appErr
		}
	} else if newQuantity > 0 {
		line.Quantity = newQuantity
		if _, appErr = a.UpsertCheckoutLine(line); appErr != nil {
			return nil, nil, appErr
		}
	}

	return &checkoutInfo.Checkout, nil, nil
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
func (a *ServiceCheckout) AddVariantsToCheckout(checkOut *checkout.Checkout, variants []*product_and_discount.ProductVariant, quantities []int, channelSlug string, skipStockCheck, replace bool) (*checkout.Checkout, *exception.InsufficientStock, *model.AppError) {
	// validate input arguments:
	var invlArgs string
	if checkOut == nil {
		invlArgs = "checkout"
	}
	if len(variants) == 0 {
		invlArgs += ", variants"
	}
	if invlArgs != "" {
		return nil, nil, model.NewAppError("AddVariantsToCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invlArgs}, "", http.StatusBadRequest)
	}

	// check quantities
	countryCode, appErr := a.CheckoutCountry(checkOut)
	if appErr != nil {
		return nil, nil, appErr
	}
	if !skipStockCheck {
		insfStock, appErr := a.srv.WarehouseService().CheckStockAndPreorderQuantityBulk(variants, countryCode, quantities, channelSlug, nil, nil, false)
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
					Eq: checkOut.ChannelID,
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

	linesOfCheckout, appErr := a.CheckoutLinesByCheckoutToken(checkOut.Token)
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
		variant := variants[i]
		quantity := quantities[i]

		if checkoutLine, exist := variantIDsInLines[variant.Id]; exist {
			if quantity > 0 {
				if replace {
					checkoutLine.Quantity = quantity
				} else {
					checkoutLine.Quantity += quantity
				}
				toUpdateCheckoutLines = append(toUpdateCheckoutLines, checkoutLine)
			} else {
				toDeleteCheckoutLineIDs = append(toDeleteCheckoutLineIDs, checkoutLine.Id)
			}
		} else if quantity > 0 {
			toCreateCheckoutLines = append(toCreateCheckoutLines, &checkout.CheckoutLine{
				CheckoutID: checkOut.Token,
				VariantID:  variant.Id,
				Quantity:   quantity,
			})
		}
	}

	if len(toDeleteCheckoutLineIDs) > 0 {
		appErr = a.DeleteCheckoutLines(nil, toDeleteCheckoutLineIDs)
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
func (a *ServiceCheckout) checkNewCheckoutAddress(checkOut *checkout.Checkout, address *account.Address, addressType string) (bool, bool, *model.AppError) {
	// validate if non-nill checkout was provided
	var invalidArguments string
	if checkOut == nil {
		invalidArguments = "checkOut"
	}
	if invalidArguments != "" {
		return false, false, model.NewAppError("checkNewCheckoutAddress", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": invalidArguments}, "", http.StatusBadRequest)
	}

	oldAddressId := checkOut.ShippingAddressID
	if addressType == account.ADDRESS_TYPE_BILLING {
		oldAddressId = checkOut.BillingAddressID
	}

	hasAddressChanged := (address == nil || oldAddressId != nil) ||
		(address != nil && oldAddressId == nil) ||
		(address != nil && oldAddressId != nil && address.Id != *oldAddressId)

	if oldAddressId == nil {
		return hasAddressChanged, false, nil
	} else {
		if checkOut.UserID == nil {
			return hasAddressChanged, hasAddressChanged, nil
		} else {
			var oldAddressNOTbelongToCheckoutUser bool
			addressesOfCheckoutUser, appErr := a.srv.AccountService().AddressesByUserId(*checkOut.UserID)
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

func (a *ServiceCheckout) ChangeBillingAddressInCheckout(checkOut *checkout.Checkout, address *account.Address) *model.AppError {
	changed, remove, appErr := a.checkNewCheckoutAddress(checkOut, address, account.ADDRESS_TYPE_BILLING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove {
			appErr = a.srv.AccountService().DeleteAddresses(*checkOut.BillingAddressID)
			if appErr != nil {
				return appErr
			}
		}
		checkOut.BillingAddressID = &address.Id
		_, appErr = a.UpsertCheckout(checkOut)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// Save shipping address in checkout if changed.
//
// Remove previously saved address if not connected to any user.
func (a *ServiceCheckout) ChangeShippingAddressInCheckout(checkoutInfo checkout.CheckoutInfo, address *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interfaces.PluginManagerInterface) *model.AppError {
	checkOut := checkoutInfo.Checkout
	changed, remove, appErr := a.checkNewCheckoutAddress(&checkOut, address, account.ADDRESS_TYPE_SHIPPING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove && checkOut.ShippingAddressID != nil {
			appErr = a.srv.AccountService().DeleteAddresses(*checkOut.ShippingAddressID)
			if appErr != nil {
				return appErr
			}
		}

		checkOut.ShippingAddressID = &address.Id
		appErr = a.UpdateCheckoutInfoShippingAddress(checkoutInfo, address, lines, discounts, manager)
		if appErr != nil {
			return appErr
		}
		_, appErr = a.UpsertCheckout(&checkOut)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// getShippingVoucherDiscountForCheckout Calculate discount value for a voucher of shipping type
func (s *ServiceCheckout) getShippingVoucherDiscountForCheckout(manager interfaces.PluginManagerInterface, voucher *product_and_discount.Voucher, checkoutInfo checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.Money, *product_and_discount.NotApplicable, *model.AppError) {
	shippingRequired, appErr := s.srv.ProductService().ProductsRequireShipping(lines.Products().IDs())
	if appErr != nil {
		return nil, nil, appErr
	}
	if !shippingRequired {
		return nil, product_and_discount.NewNotApplicable("getShippingVoucherDiscountForCheckout", "Your order does not require shipping.", nil, 0), nil
	}

	if checkoutInfo.DeliveryMethodInfo.GetDeliveryMethod() == nil {
		return nil, product_and_discount.NewNotApplicable("getShippingVoucherDiscountForCheckout", "Please select a delivery method first.", nil, 0), nil
	}

	// check if voucher is limited to specified countries
	if address != nil {
		if voucher.Countries != "" && !strings.Contains(voucher.Countries, address.Country) {
			return nil, product_and_discount.NewNotApplicable("getShippingVoucherDiscountForCheckout", "This offer is not valid in your country.", nil, 0), nil
		}
	}

	checkoutShippingPrice, appErr := s.CheckoutShippingPrice(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, appErr
	}

	money, appErr := s.srv.DiscountService().GetDiscountAmountFor(voucher, checkoutShippingPrice.Gross, checkoutInfo.Channel.Id)
	return money.(*goprices.Money), nil, appErr
}

// getProductsVoucherDiscount Calculate products discount value for a voucher, depending on its type
func (s *ServiceCheckout) getProductsVoucherDiscount(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher, discounts []*product_and_discount.DiscountInfo) (*goprices.Money, *product_and_discount.NotApplicable, *model.AppError) {
	var prices []*goprices.Money

	if voucher.Type == product_and_discount.SPECIFIC_PRODUCT {
		moneys, appErr := s.GetPricesOfDiscountedSpecificProduct(manager, checkoutInfo, lines, voucher, discounts)
		if appErr != nil {
			return nil, nil, appErr
		}
		prices = moneys
	}

	if prices == nil || len(prices) == 0 {
		return nil, product_and_discount.NewNotApplicable("getProductsVoucherDiscount", "This offer is only valid for selected items.", nil, 0), nil
	}

	money, appErr := s.srv.DiscountService().GetProductsVoucherDiscount(voucher, prices, checkoutInfo.Channel.Id)
	return money, nil, appErr
}

// GetPricesOfDiscountedSpecificProduct Get prices of variants belonging to the discounted specific products.
// Specific products are products, collections and categories.
// Product must be assigned directly to the discounted category, assigning
// product to child category won't work.
func (s *ServiceCheckout) GetPricesOfDiscountedSpecificProduct(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher, discounts []*product_and_discount.DiscountInfo) ([]*goprices.Money, *model.AppError) {
	var linePrices []*goprices.Money

	discountedLines, appErr := s.GetDiscountedLines(lines, voucher)
	if appErr != nil {
		return nil, appErr
	}

	addresses := checkoutInfo.ShippingAddress
	if addresses == nil {
		addresses = checkoutInfo.BillingAddress
	}
	if discounts == nil {
		discounts = []*product_and_discount.DiscountInfo{}
	}

	for _, lineInfo := range discountedLines {
		lineTotal, appErr := s.CheckoutLineTotal(manager, checkoutInfo, lines, lineInfo, discounts)
		if appErr != nil {
			return nil, appErr
		}
		taxedMoney, appErr := manager.CalculateCheckoutLineUnitPrice(*lineTotal, lineInfo.Line.Quantity, checkoutInfo, lines, *lineInfo, addresses, discounts)
		if appErr != nil {
			return nil, appErr
		}

		for i := 0; i < lineInfo.Line.Quantity; i++ {
			linePrices = append(linePrices, taxedMoney.Gross)
		}
	}

	return linePrices, nil
}

// GetVoucherDiscountForCheckout Calculate discount value depending on voucher and discount types.
// Raise NotApplicable if voucher of given type cannot be applied.
func (s *ServiceCheckout) GetVoucherDiscountForCheckout(manager interfaces.PluginManagerInterface, voucher *product_and_discount.Voucher, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.Money, *product_and_discount.NotApplicable, *model.AppError) {
	notApplicable, appErr := s.srv.DiscountService().ValidateVoucherForCheckout(manager, voucher, checkoutInfo, lines, discounts)
	if notApplicable != nil || appErr != nil {
		return nil, notApplicable, appErr
	}
	if voucher.Type == product_and_discount.ENTIRE_ORDER {
		checkoutSubTotal, appErr := s.CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return nil, nil, appErr
		}
		money, appErr := s.srv.DiscountService().GetDiscountAmountFor(voucher, checkoutSubTotal.Gross, checkoutInfo.Channel.Id)
		if appErr != nil {
			return nil, nil, appErr
		}
		return money.(*goprices.Money), nil, nil
	}
	if voucher.Type == product_and_discount.SHIPPING {
		return s.getShippingVoucherDiscountForCheckout(manager, voucher, checkoutInfo, lines, address, discounts)
	}
	if voucher.Type == product_and_discount.SPECIFIC_PRODUCT {
		return s.getProductsVoucherDiscount(manager, checkoutInfo, lines, voucher, discounts)
	}

	s.srv.Log.Warn("Unknown discount type", slog.String("discount_type", voucher.Type))
	return nil, nil, model.NewAppError("GetVoucherDiscountForCheckout", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "voucher.Type"}, "", http.StatusBadRequest)
}

func (a *ServiceCheckout) GetDiscountedLines(checkoutLineInfos []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	var (
		discountedProducts    []*product_and_discount.Product
		discountedCategories  []*product_and_discount.Category
		discountedCollections []*product_and_discount.Collection
		appError              *model.AppError
		mut                   sync.Mutex
		wg                    sync.WaitGroup
	)

	setErr := func(err *model.AppError) {
		mut.Lock()
		if err != nil {
			appError = err
		}
		mut.Unlock()
	}

	// starting 3 go routines
	wg.Add(3)

	go func() {
		products, appErr := a.srv.ProductService().ProductsByVoucherID(voucher.Id)
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedProducts = products
		}

		wg.Done()
	}()

	go func() {
		categories, appErr := a.srv.ProductService().CategoriesByOption(&product_and_discount.CategoryFilterOption{
			VoucherID: squirrel.Eq{a.srv.Store.VoucherCategory().TableName("VoucherID"): voucher.Id},
		})
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedCategories = categories
		}

		wg.Done()
	}()

	go func() {
		collections, appErr := a.srv.ProductService().CollectionsByVoucherID(voucher.Id)
		if appErr != nil {
			setErr(appErr)
		} else {
			discountedCollections = collections
		}

		wg.Done()
	}()

	wg.Done()

	if appError != nil {
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

			if util.StringInSlice(lineInfo.Product.Id, discountedProductIDs) ||
				(lineInfo.Product.CategoryID != nil && util.StringInSlice(*lineInfo.Product.CategoryID, discountedCategoryIDs)) ||
				lineInfoCollections_have_common_with_discountedCollections {
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
func (a *ServiceCheckout) GetVoucherForCheckout(checkoutInfo checkout.CheckoutInfo, withLock bool) (*product_and_discount.Voucher, *model.AppError) {

	now := model.NewTime(time.Now()) // NOTE: not sure to use UTC or system time
	checKout := checkoutInfo.Checkout

	voucherFilterOption := &product_and_discount.VoucherFilterOption{
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
	}

	if checKout.VoucherCode != nil {
		// finds vouchers that are active in a channel
		activeInChannelVouchers, appErr := a.srv.DiscountService().VouchersByOption(voucherFilterOption)

		if appErr != nil || len(activeInChannelVouchers) == 0 {
			return nil, appErr
		}

		// find voucher with code
		for _, voucher := range activeInChannelVouchers {
			if voucher.Code == *checKout.VoucherCode && voucher.UsageLimit != nil && withLock {

				voucherFilterOption.WithLook = true // this tell database to append `FOR UPDATE` to the end of query
				voucher, appErr = a.srv.DiscountService().VoucherByOption(voucherFilterOption)
				if appErr != nil {
					return nil, appErr
				}
			}
			return voucher, nil
		}
	}

	return nil, nil
}

// RecalculateCheckoutDiscount Recalculate `checkout.discount` based on the voucher.
// Will clear both voucher and discount if the discount is no longer applicable.
func (s *ServiceCheckout) RecalculateCheckoutDiscount(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) *model.AppError {
	checkOut := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, false)
	if appErr != nil {
		return appErr
	}

	if voucher != nil {
		address := checkoutInfo.ShippingAddress
		if address == nil {
			address = checkoutInfo.BillingAddress
		}

		discount, notApplicable, appErr := s.GetVoucherDiscountForCheckout(manager, voucher, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return appErr
		}
		if notApplicable != nil {
			appErr = s.RemoveVoucherFromCheckout(&checkOut)
			if appErr != nil {
				return appErr
			}
		}

		checkoutSubTotal, appErr := s.CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return appErr
		}
		if voucher.Type != product_and_discount.SHIPPING {
			if less, err := checkoutSubTotal.Gross.LessThan(discount); less && err == nil {
				checkOut.Discount = checkoutSubTotal.Gross
			} else {
				checkOut.Discount = discount
			}
		} else {
			checkOut.Discount = discount
		}
		checkOut.DiscountName = voucher.Name

		// check if the owner of this checkout has ther primary language:
		if checkoutInfo.User != nil && model.Languages[checkoutInfo.User.Locale] != "" {
			voucherTranslation, appErr := s.srv.DiscountService().GetVoucherTranslationByOption(&product_and_discount.VoucherTranslationFilterOption{
				LanguageCode: &model.StringFilter{
					StringOption: &model.StringOption{
						Eq: checkoutInfo.User.Locale,
					},
				},
				VoucherID: &model.StringFilter{
					StringOption: &model.StringOption{
						Eq: voucher.Id,
					},
				},
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					return appErr
				}
				// ignore not found error
			} else {
				if voucherTranslation.Name != *voucher.Name {
					checkOut.TranslatedDiscountName = &voucherTranslation.Name
				} else {
					checkOut.TranslatedDiscountName = model.NewString("")
				}
			}
		}
		_, appErr = s.UpsertCheckout(&checkOut)
		if appErr != nil {
			return appErr
		}

		return nil
	}

	return s.RemoveVoucherFromCheckout(&checkOut)
}

// AddPromoCodeToCheckout Add gift card or voucher data to checkout.
// Raise InvalidPromoCode if promo code does not match to any voucher or gift card.
func (s *ServiceCheckout) AddPromoCodeToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, promoCode string, discounts []*product_and_discount.DiscountInfo) (*giftcard.InvalidPromoCode, *model.AppError) {
	codeIsVoucher, appErr := s.srv.DiscountService().PromoCodeIsVoucher(promoCode)
	if appErr != nil {
		return nil, appErr
	}

	if codeIsVoucher {
		return s.AddVoucherCodeToCheckout(manager, checkoutInfo, lines, promoCode, discounts)
	}

	codeIsGiftcard, appErr := s.srv.GiftcardService().PromoCodeIsGiftCard(promoCode)
	if appErr != nil {
		return nil, appErr
	}

	if codeIsGiftcard {
		return s.srv.GiftcardService().AddGiftcardCodeToCheckout(&checkoutInfo.Checkout, checkoutInfo.GetCustomerEmail(), promoCode, checkoutInfo.Channel.Currency)
	}

	return giftcard.NewInvalidPromoCode("AddPromoCodeToCheckout", "Promo code is invalid"), nil
}

// AddVoucherCodeToCheckout Add voucher data to checkout by code.
// Raise InvalidPromoCode() if voucher of given type cannot be applied.
func (s *ServiceCheckout) AddVoucherCodeToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, voucherCode string, discounts []*product_and_discount.DiscountInfo) (*giftcard.InvalidPromoCode, *model.AppError) {
	vouchers, appErr := s.srv.DiscountService().FilterActiveVouchers(model.NewTime(time.Now().UTC()), checkoutInfo.Channel.Slug)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return &giftcard.InvalidPromoCode{}, nil
	}

	for _, voucher := range vouchers {
		if voucher.Code == voucherCode {
			notAplicable, appErr := s.AddVoucherToCheckout(manager, checkoutInfo, lines, voucher, discounts)
			if appErr != nil {
				return nil, appErr
			}
			if notAplicable != nil {
				return nil, model.NewAppError("AddVoucherCodeToCheckout", "app.checkout.voucher_not_applicabale_to_checkout.app_error", map[string]interface{}{"code": exception.VOUCHER_NOT_APPLICABLE}, "", http.StatusNotAcceptable)
			}
		}
	}

	return &giftcard.InvalidPromoCode{}, nil
}

// AddVoucherToCheckout Add voucher data to checkout.
// Raise NotApplicable if voucher of given type cannot be applied.
func (s *ServiceCheckout) AddVoucherToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, voucher *product_and_discount.Voucher, discounts []*product_and_discount.DiscountInfo) (*product_and_discount.NotApplicable, *model.AppError) {
	checkout := checkoutInfo.Checkout

	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}
	discountMoney, notApplicable, appErr := s.GetVoucherDiscountForCheckout(manager, voucher, checkoutInfo, lines, address, discounts)
	if appErr != nil || notApplicable != nil {
		return notApplicable, appErr
	}
	checkout.VoucherCode = &voucher.Code
	checkout.DiscountName = voucher.Name

	if user := checkoutInfo.User; user != nil && model.Languages[user.Locale] != "" {
		voucherTranslation, appErr := s.srv.DiscountService().GetVoucherTranslationByOption(&product_and_discount.VoucherTranslationFilterOption{
			LanguageCode: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: user.Locale,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}
		if voucherTranslation.Name != *voucher.Name {
			checkout.TranslatedDiscountName = &voucherTranslation.Name
		} else {
			checkout.TranslatedDiscountName = model.NewString("")
		}
	}
	checkout.Discount = discountMoney

	_, appErr = s.UpsertCheckout(&checkout)
	return nil, appErr
}

// RemovePromoCodeFromCheckout Remove gift card or voucher data from checkout.
func (a *ServiceCheckout) RemovePromoCodeFromCheckout(checkoutInfo checkout.CheckoutInfo, promoCode string) *model.AppError {
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
func (a *ServiceCheckout) RemoveVoucherCodeFromCheckout(checkoutInfo checkout.CheckoutInfo, voucherCode string) *model.AppError {
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
func (a *ServiceCheckout) RemoveVoucherFromCheckout(checkOut *checkout.Checkout) *model.AppError {
	if checkOut == nil {
		return nil
	}

	checkOut.VoucherCode = nil
	checkOut.DiscountName = nil
	checkOut.TranslatedDiscountName = nil
	checkOut.DiscountAmount = &decimal.Zero

	_, appErr := a.UpsertCheckout(checkOut)

	return appErr
}

// GetValidShippingMethodsForCheckout finds all valid shipping methods for given checkout
func (a *ServiceCheckout) GetValidShippingMethodsForCheckout(checkoutInfo checkout.CheckoutInfo, lineInfos []*checkout.CheckoutLineInfo, subTotal *goprices.TaxedMoney, countryCode string) ([]*shipping.ShippingMethod, *model.AppError) {
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

// GetValidCollectionPointsForCheckout Return a collection of `Warehouse`s that can be used as a collection point.
// Note that `quantity_check=False` should be used, when stocks quantity will
// be validated in further steps (checkout completion) in order to raise
// 'InsufficientProductStock' error instead of 'InvalidShippingError'.
func (s *ServiceCheckout) GetValidCollectionPointsForCheckout(lines checkout.CheckoutLineInfos, countryCode string, quantityCheck bool) ([]*warehouse.WareHouse, *model.AppError) {
	linesRequireShipping, appErr := s.srv.ProductService().ProductsRequireShipping(lines.Products().IDs())
	if appErr != nil {
		return nil, appErr
	}
	if !linesRequireShipping {
		return []*warehouse.WareHouse{}, nil
	}

	if model.Countries[strings.ToUpper(countryCode)] == "" {
		return []*warehouse.WareHouse{}, nil
	}
	checkoutLines, appErr := s.CheckoutLinesByOption(&checkout.CheckoutLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: lines.CheckoutLines().IDs(),
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return []*warehouse.WareHouse{}, nil
	}

	// TODO: implement me.
	panic("not implemented")
}

func (a *ServiceCheckout) ClearDeliveryMethod(checkoutInfo checkout.CheckoutInfo) *model.AppError {
	checkOut := checkoutInfo.Checkout
	checkOut.CollectionPointID = nil
	checkOut.ShippingMethodID = nil

	appErr := a.UpdateCheckoutInfoDeliveryMethod(checkoutInfo, nil)
	if appErr != nil {
		return nil
	}

	_, appErr = a.UpsertCheckout(&checkOut)
	return appErr
}

// IsFullyPaid Check if provided payment methods cover the checkout's total amount.
// Note that these payments may not be captured or charged at all.
func (s *ServiceCheckout) IsFullyPaid(manager interfaces.PluginManagerInterface, checkoutInfo checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) (bool, *model.AppError) {
	checkOut := checkoutInfo.Checkout
	payments, appErr := s.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		CheckoutToken: checkOut.Token,
		IsActive:      model.NewBool(true),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
		// ignore not found error
	}

	totalPaid := &decimal.Zero
	for _, payMent := range payments {
		totalPaid = model.NewDecimal(totalPaid.Add(*payMent.Total))
	}
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	checkoutTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return false, appErr
	}
	checkoutTotalGiftcardBalance, appErr := s.CheckoutTotalGiftCardsBalance(&checkOut)
	if appErr != nil {
		return false, appErr
	}

	sub, err := checkoutTotal.Sub(checkoutTotalGiftcardBalance)
	if err != nil {
		return false, model.NewAppError("IsFullyPaid", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	checkoutTotal = sub

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkOut.Currency)
	if less, err := checkoutTotal.LessThan(zeroTaxedMoney); less && err == nil {
		checkoutTotal = zeroTaxedMoney
	}

	return checkoutTotal.Gross.Amount.LessThan(*totalPaid), nil
}

// CancelActivePayments set all active payments belong to given checkout
func (a *ServiceCheckout) CancelActivePayments(checkOut *checkout.Checkout) *model.AppError {
	err := a.srv.Store.Payment().CancelActivePaymentsOfCheckout(checkOut.Token)
	if err != nil {
		return model.NewAppError("CancelActivePayments", "app.checkout.cancel_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) ValidateVariantsInCheckoutLines(lines []*checkout.CheckoutLineInfo) *model.AppError {
	var notAvailableVariantIDs []string
	for _, line := range lines {
		if line.ChannelListing.Price == nil {
			notAvailableVariantIDs = append(notAvailableVariantIDs, line.Variant.Id)
		}
	}

	if len(notAvailableVariantIDs) > 0 {
		notAvailableVariantIDs = util.RemoveDuplicatesFromStringArray(notAvailableVariantIDs)
		// return error indicate there are some product variants that have no channel listing or channel listing price is null
		return model.NewAppError("ValidateVariantsInCheckoutLines", "app.checkout.cannot_add_lines_with_unavailable_variants.app_error", map[string]interface{}{"variants": strings.Join(notAvailableVariantIDs, ", ")}, "", http.StatusNotAcceptable)
	}

	return nil
}
