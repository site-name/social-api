package checkout

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceCheckout) CheckVariantInStock(checkout *model.Checkout, variant *model.ProductVariant, channelSlug string, quantity int, replace, checkQuantity bool) (int, *model.CheckoutLine, *model_helper.InsufficientStock, *model_helper.AppError) {
	// quantity param is default to 1

	checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(checkout.Token)
	if appErr != nil {
		return 0, nil, nil, appErr
	}

	var (
		lineWithVariant *model.CheckoutLine = nil      // checkoutLine that has variantID of given `variantID`
		lineQuantity    int                 = 0        // quantity of lineWithVariant checkout line
		newQuantity     int                 = quantity //
	)
	if !replace {
		newQuantity = quantity + lineQuantity
	}

	for _, checkoutLine := range checkoutLines {
		if checkoutLine.VariantID == variant.ID {
			lineWithVariant = checkoutLine
			break
		}
	}

	if lineWithVariant != nil {
		lineQuantity = lineWithVariant.Quantity
	}

	if newQuantity < 0 {
		return 0, nil, nil, model_helper.NewAppError(
			"CheckVariantInStock",
			"app.checkout.quantity_invalid.app_error",
			map[string]any{
				"Quantity":    quantity,
				"NewQuantity": newQuantity,
			},
			"", http.StatusBadRequest,
		)
	}

	if newQuantity > 0 && checkQuantity {
		insufficientStockErr, appErr := a.srv.Warehouse.CheckStockAndPreorderQuantity(variant, checkout.Country, channelSlug, newQuantity)
		if insufficientStockErr != nil || appErr != nil {
			return 0, nil, insufficientStockErr, appErr
		}
	}

	return newQuantity, lineWithVariant, nil, nil
}

// AddVariantToCheckout adds a product variant to checkout
//
// `quantity` default to 1, `replace` default to false, `checkQuantity` default to true
func (a *ServiceCheckout) AddVariantToCheckout(checkoutInfo *model_helper.CheckoutInfo, variant *model.ProductVariant, quantity int, replace bool, checkQuantity bool) (*model.Checkout, *model_helper.InsufficientStock, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout
	productChannelListings, appErr := a.srv.Product.ProductChannelListingsByOption(model_helper.ProductChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ProductChannelListingWhere.ChannelID.EQ(checkout.ChannelID),
			model.ProductChannelListingWhere.ProductID.EQ(variant.ProductID),
		),
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	if len(productChannelListings) == 0 || !productChannelListings[0].IsPublished {
		return nil, nil, model_helper.NewAppError("AddVariantToCheckout", model_helper.ProductNotPublishedAppErrID, nil, "Please publish the product first.", http.StatusNotAcceptable)
	}

	newQuantity, line, insufficientErr, appErr := a.CheckVariantInStock(&checkoutInfo.Checkout, variant, checkoutInfo.Channel.Slug, quantity, replace, checkQuantity)
	if appErr != nil || insufficientErr != nil {
		return nil, insufficientErr, appErr
	}

	if line == nil {
		checkoutLines, appErr := a.CheckoutLinesByOption(model_helper.CheckoutLineFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.CheckoutLineWhere.CheckoutID.EQ(checkout.Token),
				model.CheckoutLineWhere.VariantID.EQ(variant.ID),
			),
		})
		if appErr != nil && appErr.StatusCode != http.StatusNotFound { // ignore not found error
			return nil, nil, appErr
		}
		line = checkoutLines[0]
	}

	if newQuantity == 0 {
		if line != nil {
			if appErr = a.DeleteCheckoutLines(nil, []string{line.ID}); appErr != nil {
				return nil, nil, appErr
			}
		}
	} else if line == nil {
		if _, appErr = a.UpsertCheckoutLine(model.CheckoutLine{
			CheckoutID: checkoutInfo.Checkout.Token,
			VariantID:  variant.ID,
			Quantity:   newQuantity,
		}); appErr != nil {
			return nil, nil, appErr
		}
	} else if newQuantity > 0 {
		line.Quantity = newQuantity
		if _, appErr = a.UpsertCheckoutLine(*line); appErr != nil {
			return nil, nil, appErr
		}
	}

	return &checkoutInfo.Checkout, nil, nil
}

func (a *ServiceCheckout) CalculateCheckoutQuantity(lineInfos model_helper.CheckoutLineInfos) (int, *model_helper.AppError) {
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
//
//	skipStockCheck and replace are default to false
func (a *ServiceCheckout) AddVariantsToCheckout(checkout model.Checkout, variants model.ProductVariantSlice, quantities []int, channelSlug string, skipStockCheck, replace bool) (*model.Checkout, *model_helper.InsufficientStock, *model_helper.AppError) {
	// check quantities
	countryCode, appErr := a.CheckoutCountry(checkout)
	if appErr != nil {
		return nil, nil, appErr
	}
	if !skipStockCheck {
		insfStock, appErr := a.srv.Warehouse.CheckStockAndPreorderQuantityBulk(variants, countryCode, quantities, channelSlug, nil, nil, false)
		if appErr != nil || insfStock != nil {
			return nil, insfStock, appErr
		}
	}

	productIDs := lo.Map(variants, func(variant *model.ProductVariant, _ int) string { return variant.ProductID })

	channelListings, appErr := a.srv.Product.
		ProductChannelListingsByOption(model_helper.ProductChannelListingFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.ProductChannelListingWhere.ChannelID.EQ(checkout.ChannelID),
				model.ProductChannelListingWhere.ProductID.IN(productIDs),
			),
		})
	if appErr != nil {
		return nil, nil, appErr
	}

	// keys are product ids
	var listingMap = make(map[string]*model.ProductChannelListing)
	for _, listing := range channelListings {
		listingMap[listing.ProductID] = listing
	}

	for _, variant := range variants {
		productChannelListing := listingMap[variant.ProductID]
		if productChannelListing == nil || !productChannelListing.IsPublished {
			return nil, nil, model_helper.NewAppError("AddVariantsToCheckout", model_helper.ProductNotPublishedAppErrID, nil, "", http.StatusNotAcceptable)
		}
	}

	linesOfCheckout, appErr := a.CheckoutLinesByCheckoutToken(checkout.Token)
	if appErr != nil {
		return nil, nil, appErr
	}

	// keys are variant ids
	var variantIDsInLines = make(map[string]*model.CheckoutLine)
	for _, line := range linesOfCheckout {
		variantIDsInLines[line.VariantID] = line
	}

	var (
		toCreateCheckoutLines   = model.CheckoutLineSlice{}
		toUpdateCheckoutLines   = model.CheckoutLineSlice{}
		toDeleteCheckoutLineIDs = []string{}
	)
	// use Min() since two slices may not have same length
	for i := 0; i < min(len(variants), len(quantities)); i++ {
		variant := variants[i]
		quantity := quantities[i]

		if checkoutLine, exist := variantIDsInLines[variant.ID]; exist {
			if quantity > 0 {
				if replace {
					checkoutLine.Quantity = quantity
				} else {
					checkoutLine.Quantity += quantity
				}
				toUpdateCheckoutLines = append(toUpdateCheckoutLines, checkoutLine)
			} else {
				toDeleteCheckoutLineIDs = append(toDeleteCheckoutLineIDs, checkoutLine.ID)
			}
		} else if quantity > 0 {
			toCreateCheckoutLines = append(toCreateCheckoutLines, &model.CheckoutLine{
				CheckoutID: checkout.Token,
				VariantID:  variant.ID,
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
		_, appErr = a.UpsertCheckoutLines(toUpdateCheckoutLines)
		if appErr != nil {
			return nil, nil, appErr
		}
	}
	if len(toCreateCheckoutLines) > 0 {
		_, appErr = a.UpsertCheckoutLines(toCreateCheckoutLines)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	return nil, nil, nil
}

// checkNewCheckoutAddress Check if and address in checkout has changed and if to remove old one
// `address` can be nil
func (a *ServiceCheckout) checkNewCheckoutAddress(checkout model.Checkout, address *model.Address, addressType model_helper.AddressTypeEnum) (bool, bool, *model_helper.AppError) {
	oldAddressId := checkout.ShippingAddressID
	if addressType == model_helper.ADDRESS_TYPE_BILLING {
		oldAddressId = checkout.BillingAddressID
	}

	hasAddressChanged := (address == nil && !oldAddressId.IsNil()) ||
		(address != nil && oldAddressId.IsNil()) ||
		(address != nil && !oldAddressId.IsNil() && address.ID != *oldAddressId.String)

	removeOldAddress := hasAddressChanged && !oldAddressId.IsNil()
	if checkout.UserID.IsNil() {
		return hasAddressChanged, removeOldAddress, nil
	}

	addresses, appErr := a.srv.Account.AddressesByUserId(*checkout.UserID.String)
	if appErr != nil {
		return false, false, appErr
	}

	removeOldAddress = removeOldAddress && !lo.SomeBy(addresses, func(addr *model.Address) bool { return addr.ID == *oldAddressId.String })
	return hasAddressChanged, removeOldAddress, nil
}

func (a *ServiceCheckout) ChangeBillingAddressInCheckout(transaction boil.ContextTransactor, checkout model.Checkout, address *model.Address) *model_helper.AppError {
	changed, remove, appErr := a.checkNewCheckoutAddress(checkout, address, model_helper.ADDRESS_TYPE_BILLING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove && !checkout.BillingAddressID.IsNil() {
			appErr = a.srv.Account.DeleteAddresses(transaction, []string{*checkout.BillingAddressID.String})
			if appErr != nil {
				return appErr
			}
		}
		checkout.BillingAddressID.String = &address.ID
		_, appErr = a.UpsertCheckouts(transaction, model.CheckoutSlice{&checkout})
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// Save shipping address in checkout if changed.
//
// Remove previously saved address if not connected to any user.
func (a *ServiceCheckout) ChangeShippingAddressInCheckout(transaction boil.ContextTransactor, checkoutInfo model_helper.CheckoutInfo, address *model.Address, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo, manager interfaces.PluginManagerInterface) *model_helper.AppError {
	checkout := checkoutInfo.Checkout
	changed, remove, appErr := a.checkNewCheckoutAddress(checkout, address, model_helper.ADDRESS_TYPE_SHIPPING)
	if appErr != nil {
		return appErr
	}

	if changed {
		if remove && !checkout.ShippingAddressID.IsNil() {
			appErr = a.srv.Account.DeleteAddresses(transaction, []string{*checkout.ShippingAddressID.String})
			if appErr != nil {
				return appErr
			}
		}

		checkout.ShippingAddressID.String = &address.ID
		appErr = a.updateCheckoutInfoShippingAddress(checkoutInfo, address, lines, discounts, manager)
		if appErr != nil {
			return appErr
		}
		_, appErr = a.UpsertCheckouts(transaction, model.CheckoutSlice{&checkout})
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// getShippingVoucherDiscountForCheckout Calculate discount value for a voucher of shipping type
func (s *ServiceCheckout) getShippingVoucherDiscountForCheckout(manager interfaces.PluginManagerInterface, voucher model.Voucher, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.Money, *model_helper.NotApplicable, *model_helper.AppError) {
	productIDs := lo.Map(lines.Products(), func(prd *model.Product, _ int) string { return prd.ID })
	shippingRequired, appErr := s.srv.Product.ProductsRequireShipping(productIDs)
	if appErr != nil {
		return nil, nil, appErr
	}
	if !shippingRequired {
		return nil, model_helper.NewNotApplicable("getShippingVoucherDiscountForCheckout", "Your order does not require shipping.", nil, 0), nil
	}

	if checkoutInfo.DeliveryMethodInfo.GetDeliveryMethod() == nil {
		return nil, model_helper.NewNotApplicable("getShippingVoucherDiscountForCheckout", "Please select a delivery method first.", nil, 0), nil
	}

	// check if voucher is limited to specified countries
	if address != nil {
		if voucher.Countries != "" && !strings.Contains(voucher.Countries, string(address.Country)) {
			return nil, model_helper.NewNotApplicable("getShippingVoucherDiscountForCheckout", "This offer is not valid in your country.", nil, 0), nil
		}
	}

	checkoutShippingPrice, appErr := s.CheckoutShippingPrice(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, appErr
	}

	money, appErr := s.srv.Discount.GetDiscountAmountFor(voucher, checkoutShippingPrice.Gross, checkoutInfo.Channel.ID)
	return money.(*goprices.Money), nil, appErr
}

// getProductsVoucherDiscount Calculate products discount value for a voucher, depending on its type
func (s *ServiceCheckout) getProductsVoucherDiscount(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, voucher model.Voucher, discounts []*model_helper.DiscountInfo) (*goprices.Money, *model_helper.NotApplicable, *model_helper.AppError) {
	var prices []*goprices.Money

	if voucher.Type == model.VoucherTypeSpecificProduct {
		moneys, appErr := s.GetPricesOfDiscountedSpecificProduct(manager, checkoutInfo, lines, voucher, discounts)
		if appErr != nil {
			return nil, nil, appErr
		}
		prices = moneys
	}

	if len(prices) == 0 {
		return nil, model_helper.NewNotApplicable("getProductsVoucherDiscount", "This offer is only valid for selected items.", nil, 0), nil
	}

	money, appErr := s.srv.Discount.GetProductsVoucherDiscount(voucher, prices, checkoutInfo.Channel.ID)
	return money, nil, appErr
}

// GetPricesOfDiscountedSpecificProduct Get prices of variants belonging to the discounted specific products.
// Specific products are products, collections and categories.
// Product must be assigned directly to the discounted category, assigning
// product to child category won't work.
func (s *ServiceCheckout) GetPricesOfDiscountedSpecificProduct(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, voucher model.Voucher, discounts []*model_helper.DiscountInfo) ([]*goprices.Money, *model_helper.AppError) {
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
		discounts = []*model_helper.DiscountInfo{}
	}

	for _, lineInfo := range discountedLines {
		if lineInfo == nil {
			continue
		}
		lineTotal, appErr := s.CheckoutLineTotal(manager, checkoutInfo, lines, *lineInfo, discounts)
		if appErr != nil {
			return nil, appErr
		}
		taxedMoney, appErr := manager.CalculateCheckoutLineUnitPrice(*lineTotal, lineInfo.Line.Quantity, checkoutInfo, lines, *lineInfo, addresses, discounts)
		if appErr != nil {
			return nil, appErr
		}

		for i := 0; i < lineInfo.Line.Quantity; i++ {
			linePrices = append(linePrices, &taxedMoney.Gross)
		}
	}

	return linePrices, nil
}

// GetVoucherDiscountForCheckout Calculate discount value depending on voucher and discount types.
// Raise NotApplicable if voucher of given type cannot be applied.
func (s *ServiceCheckout) GetVoucherDiscountForCheckout(manager interfaces.PluginManagerInterface, voucher model.Voucher, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.Money, *model_helper.NotApplicable, *model_helper.AppError) {
	notApplicable, appErr := s.srv.Discount.ValidateVoucherForCheckout(manager, voucher, checkoutInfo, lines, discounts)
	if notApplicable != nil || appErr != nil {
		return nil, notApplicable, appErr
	}
	if voucher.Type == model.VoucherTypeEntireOrder {
		checkoutSubTotal, appErr := s.CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return nil, nil, appErr
		}
		money, appErr := s.srv.Discount.GetDiscountAmountFor(voucher, checkoutSubTotal.Gross, checkoutInfo.Channel.ID)
		if appErr != nil {
			return nil, nil, appErr
		}
		return money.(*goprices.Money), nil, nil
	}
	if voucher.Type == model.VoucherTypeShipping {
		return s.getShippingVoucherDiscountForCheckout(manager, voucher, checkoutInfo, lines, address, discounts)
	}
	if voucher.Type == model.VoucherTypeSpecificProduct {
		return s.getProductsVoucherDiscount(manager, checkoutInfo, lines, voucher, discounts)
	}

	s.srv.Log.Warn("Unknown discount type", slog.String("discount_type", string(voucher.Type)))
	return nil, nil, model_helper.NewAppError("GetVoucherDiscountForCheckout", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "voucher.Type"}, "", http.StatusBadRequest)
}

func (a *ServiceCheckout) GetDiscountedLines(checkoutLineInfos model_helper.CheckoutLineInfos, voucher model.Voucher) (model_helper.CheckoutLineInfos, *model_helper.AppError) {
	var (
		discountedProducts    []*model.Product
		discountedCategories  []*model.Category
		discountedCollections []*model.Collection
	)

	var (
		atomicValue atomic.Int32
		appErrChan  = make(chan *model_helper.AppError)
	)
	defer close(appErrChan)
	atomicValue.Add(3) //

	go func() {
		defer atomicValue.Add(-1)
		products, appErr := a.srv.Product.ProductsByVoucherID(voucher.ID)
		if appErr != nil {
			appErrChan <- appErr
			return
		}
		discountedProducts = products
	}()

	go func() {
		defer atomicValue.Add(-1)
		categories, appErr := a.srv.Product.CategoriesByOption(model_helper.CategoryFilterOption{
			VoucherID: model.VoucherCategoryWhere.VoucherID.EQ(voucher.ID),
		})
		if appErr != nil {
			appErrChan <- appErr
			return
		}
		discountedCategories = categories
	}()

	go func() {
		defer atomicValue.Add(-1)
		collections, appErr := a.srv.Product.CollectionsByVoucherID(voucher.ID)
		if appErr != nil {
			appErrChan <- appErr
			return
		}
		discountedCollections = collections
	}()

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		default:
		}
	}

	var discountedProductIDs, discountedCategoryIDs, discountedCollectionIDs util.AnyArray[string]

	for _, prd := range discountedProducts {
		discountedProductIDs = append(discountedProductIDs, prd.ID)
	}

	// filter duplicates from discountedCategories:
	meetMap := map[string]bool{}
	for _, category := range discountedCategories {
		if _, ok := meetMap[category.ID]; !ok {
			discountedCategoryIDs = append(discountedCategoryIDs, category.ID)
			meetMap[category.ID] = true
		}
	}
	// filter duplicates from discountedCollections:
	// NOTE: reuse meetMap here since UUIDs are unique
	for _, collection := range discountedCollections {
		if _, ok := meetMap[collection.ID]; !ok {
			discountedCollectionIDs = append(discountedCollectionIDs, collection.ID)
			meetMap[collection.ID] = true
		}
	}

	var discountedLines model_helper.CheckoutLineInfos
	if len(discountedProductIDs) > 0 || len(discountedCategoryIDs) > 0 || len(discountedCollectionIDs) > 0 {
		for _, lineInfo := range checkoutLineInfos {

			var lineInfoCollections_have_common_with_discountedCollections bool
			for _, collection := range lineInfo.Collections {
				if yes, exist := meetMap[collection.ID]; yes && exist {
					lineInfoCollections_have_common_with_discountedCollections = true
					break
				}
			}

			if discountedProductIDs.Contains(lineInfo.Product.ID) ||
				discountedCategoryIDs.Contains(lineInfo.Product.CategoryID) ||
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
func (a *ServiceCheckout) GetVoucherForCheckout(checkoutInfo model_helper.CheckoutInfo, vouchers model.VoucherSlice, withLock bool) (*model.Voucher, *model_helper.AppError) {
	now := model_helper.GetPointerOfValue(time.Now().UTC()) // NOTE: not sure to use UTC or system time
	checkout := checkoutInfo.Checkout

	if !checkout.VoucherCode.IsNil() {
		if len(vouchers) == 0 {
			// finds vouchers that are active in a channel
			var appErr *model_helper.AppError
			vouchers, appErr = a.srv.Discount.VouchersByOption(model_helper.VoucherFilterOption{
				VoucherChannelListing_ChannelSlug:     squirrel.Eq{model.ChannelTableName + ".Slug": checkoutInfo.Channel.Slug},
				VoucherChannelListing_ChannelIsActive: squirrel.Eq{model.ChannelTableName + ".IsActive": true},
				// Conditions: squirrel.And{
				// 	squirrel.Or{
				// 		squirrel.Eq{model.VoucherTableName + ".UsageLimit": nil},
				// 		squirrel.GtOrEq{model.VoucherTableName + ".UsageLimit": model.VoucherTableName + ".Used"},
				// 	},
				// 	squirrel.Or{
				// 		squirrel.Eq{model.VoucherTableName + ".EndDate": nil},
				// 		squirrel.GtOrEq{model.VoucherTableName + ".EndDate": now},
				// 	},
				// 	squirrel.LtOrEq{model.VoucherTableName + ".StartDate": now},
				// },
				CommonQueryOptions: model_helper.NewCommonQueryOptions(
					model_helper.Or{
						squirrel.Eq{model.VoucherTableColumns.UsageLimit: nil},
						squirrel.GtOrEq{model.VoucherTableColumns.UsageLimit: model.VoucherTableColumns.Used},
					},
					model_helper.Or{
						squirrel.Eq{model.VoucherTableColumns.EndDate: nil},
						squirrel.GtOrEq{model.VoucherTableColumns.EndDate: now},
					},
				),
			})
			if appErr != nil {
				return nil, appErr
			}
			if len(vouchers) == 0 {
				return nil, nil
			}
		}

		voucher, found := lo.Find(vouchers, func(v *model.Voucher) bool { return v != nil && !v.UsageLimit.IsNil() })
		if found && withLock {
			voucher, appErr := a.srv.Discount.VoucherByOption(model_helper.VoucherFilterOption{
				CommonQueryOptions: model_helper.NewCommonQueryOptions(
					model.VoucherWhere.ID.EQ(voucher.ID),
				),
			})
			if appErr != nil {
				return nil, appErr
			}

			return voucher, nil
		}

		return nil, nil
	}

	return nil, nil
}

// RecalculateCheckoutDiscount Recalculate `checkout.discount` based on the voucher.
// Will clear both voucher and discount if the discount is no longer applicable.
func (s *ServiceCheckout) RecalculateCheckoutDiscount(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo) *model_helper.AppError {
	checkout := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, nil, false)
	if appErr != nil {
		return appErr
	}

	if voucher != nil {
		address := checkoutInfo.ShippingAddress
		if address == nil {
			address = checkoutInfo.BillingAddress
		}

		discount, notApplicable, appErr := s.GetVoucherDiscountForCheckout(manager, *voucher, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return appErr
		}
		if notApplicable != nil {
			appErr = s.RemoveVoucherFromCheckout(&checkout)
			if appErr != nil {
				return appErr
			}
		}

		checkoutSubTotal, appErr := s.CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
		if appErr != nil {
			return appErr
		}

		if voucher.Type != model.VoucherTypeShipping {
			if checkoutSubTotal.Gross.LessThan(*discount) {
				model_helper.CheckoutSetDiscountAmount(&checkout, checkoutSubTotal.Gross)
			} else {
				model_helper.CheckoutSetDiscountAmount(&checkout, *discount)
			}
		} else {
			model_helper.CheckoutSetDiscountAmount(&checkout, *discount)
		}
		checkout.DiscountName = voucher.Name

		// check if the owner of this checkout has ther primary language:
		if checkoutInfo.User != nil {
			voucherTranslation, appErr := s.srv.Discount.GetVoucherTranslationByOption(model_helper.VoucherTranslationFilterOption{
				CommonQueryOptions: model_helper.NewCommonQueryOptions(
					model.VoucherTranslationWhere.LanguageCode.EQ(checkoutInfo.User.Locale),
					model.VoucherTranslationWhere.VoucherID.EQ(voucher.ID),
				),
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					return appErr
				}
				// ignore not found error
			} else {
				if model_helper.GetValueOfPointerOrZero(voucher.Name.String) != voucherTranslation.Name {
					checkout.TranslatedDiscountName.String = &voucherTranslation.Name
				} else {
					checkout.TranslatedDiscountName.String = model_helper.GetPointerOfValue("")
				}
			}
		}
		_, appErr = s.UpsertCheckouts(nil, model.CheckoutSlice{&checkout})
		if appErr != nil {
			return appErr
		}

		return nil
	}

	return s.RemoveVoucherFromCheckout(&checkout)
}

// AddPromoCodeToCheckout Add gift card or voucher data to checkout.
// Raise InvalidPromoCode if promo code does not match to any voucher or gift card.
func (s *ServiceCheckout) AddPromoCodeToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, promoCode string, discounts []*model_helper.DiscountInfo) (*model_helper.InvalidPromoCode, *model_helper.AppError) {
	codeIsVoucher, appErr := s.srv.Discount.PromoCodeIsVoucher(promoCode)
	if appErr != nil {
		return nil, appErr
	}

	if codeIsVoucher {
		return s.AddVoucherCodeToCheckout(manager, checkoutInfo, lines, promoCode, discounts)
	}

	codeIsGiftcard, appErr := s.srv.Giftcard.PromoCodeIsGiftCard(promoCode)
	if appErr != nil {
		return nil, appErr
	}

	if codeIsGiftcard {
		return s.srv.Giftcard.AddGiftcardCodeToCheckout(&checkoutInfo.Checkout, checkoutInfo.GetCustomerEmail(), promoCode, checkoutInfo.Channel.Currency.String())
	}

	return model_helper.NewInvalidPromoCode("AddPromoCodeToCheckout", "Promo code is invalid"), nil
}

// AddVoucherCodeToCheckout Add voucher data to checkout by code.
// Raise InvalidPromoCode() if voucher of given type cannot be applied.
func (s *ServiceCheckout) AddVoucherCodeToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, voucherCode string, discounts []*model_helper.DiscountInfo) (*model_helper.InvalidPromoCode, *model_helper.AppError) {
	vouchers, appErr := s.srv.Discount.FilterActiveVouchers(time.Now().UTC(), checkoutInfo.Channel.Slug)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return &model_helper.InvalidPromoCode{}, nil
	}

	for _, voucher := range vouchers {
		if voucher == nil {
			continue
		}
		if voucher.Code == voucherCode {
			notAplicable, appErr := s.AddVoucherToCheckout(manager, checkoutInfo, lines, *voucher, discounts)
			if appErr != nil {
				return nil, appErr
			}
			if notAplicable != nil {
				return nil, model_helper.NewAppError("AddVoucherCodeToCheckout", "app.model.voucher_not_applicabale_to_checkout.app_error", map[string]any{"code": model_helper.VOUCHER_NOT_APPLICABLE}, "", http.StatusNotAcceptable)
			}
		}
	}

	return &model_helper.InvalidPromoCode{}, nil
}

// AddVoucherToCheckout Add voucher data to checkout.
// Raise NotApplicable if voucher of given type cannot be applied.
func (s *ServiceCheckout) AddVoucherToCheckout(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, voucher model.Voucher, discounts []*model_helper.DiscountInfo) (*model_helper.NotApplicable, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout

	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}
	discountMoney, notApplicable, appErr := s.GetVoucherDiscountForCheckout(manager, voucher, checkoutInfo, lines, address, discounts)
	if appErr != nil || notApplicable != nil {
		return notApplicable, appErr
	}
	checkout.VoucherCode.String = &voucher.Code
	checkout.DiscountName = voucher.Name

	if user := checkoutInfo.User; user != nil {
		voucherTranslation, appErr := s.srv.Discount.GetVoucherTranslationByOption(model_helper.VoucherTranslationFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.VoucherTranslationWhere.LanguageCode.EQ(user.Locale),
			),
		})
		if appErr != nil {
			return nil, appErr
		}
		if voucherTranslation.Name != model_helper.GetValueOfPointerOrZero(voucher.Name.String) {
			checkout.TranslatedDiscountName.String = &voucherTranslation.Name
		} else {
			checkout.TranslatedDiscountName.String = model_helper.GetPointerOfValue("")
		}
	}
	model_helper.CheckoutSetDiscountAmount(&checkout, *discountMoney)

	_, appErr = s.UpsertCheckouts(nil, model.CheckoutSlice{&checkout})
	return nil, appErr
}

// RemovePromoCodeFromCheckout Remove gift card or voucher data from checkout.
func (a *ServiceCheckout) RemovePromoCodeFromCheckout(checkoutInfo model_helper.CheckoutInfo, promoCode string) *model_helper.AppError {
	// check if promoCode is voucher:
	promoCodeIsVoucher, appErr := a.srv.Discount.PromoCodeIsVoucher(promoCode)
	if appErr != nil { // this error is system error
		return appErr
	}
	if promoCodeIsVoucher {
		return a.RemoveVoucherCodeFromCheckout(checkoutInfo, promoCode)
	}

	// check promoCode is giftcard
	promoCodeIsGiftCard, appErr := a.srv.Giftcard.PromoCodeIsGiftCard(promoCode)
	if appErr != nil {
		return appErr
	}
	if promoCodeIsGiftCard {
		return a.srv.Giftcard.RemoveGiftcardCodeFromCheckout(&checkoutInfo.Checkout, promoCode)
	}

	return nil
}

// RemoveVoucherCodeFromCheckout Remove voucher data from checkout by code.
func (a *ServiceCheckout) RemoveVoucherCodeFromCheckout(checkoutInfo model_helper.CheckoutInfo, voucherCode string) *model_helper.AppError {
	existingVoucher, appErr := a.GetVoucherForCheckout(checkoutInfo, nil, false)
	if appErr != nil {
		return appErr
	}
	if existingVoucher != nil && existingVoucher.Code == voucherCode {
		return a.RemoveVoucherFromCheckout(&checkoutInfo.Checkout)
	}

	return nil
}

// RemoveVoucherFromCheckout removes voucher data from checkout
func (a *ServiceCheckout) RemoveVoucherFromCheckout(checkout *model.Checkout) *model_helper.AppError {
	if checkout == nil {
		return nil
	}

	checkout.VoucherCode.String = nil
	checkout.DiscountName.String = nil
	checkout.TranslatedDiscountName.String = nil
	checkout.DiscountAmount = decimal.NewFromInt(0)

	_, appErr := a.UpsertCheckouts(nil, model.CheckoutSlice{checkout})

	return appErr
}

// GetValidShippingMethodsForCheckout finds all valid shipping methods for given checkout
func (a *ServiceCheckout) GetValidShippingMethodsForCheckout(checkoutInfo model_helper.CheckoutInfo, lineInfos model_helper.CheckoutLineInfos, subTotal *goprices.TaxedMoney, countryCode model.CountryCode) (model.ShippingMethodSlice, *model_helper.AppError) {
	var productIDs []string
	for _, line := range lineInfos {
		productIDs = append(productIDs, line.Product.ID)
	}

	// check if any product in given lineInfos requires shipping:
	requireShipping, appErr := a.srv.Product.ProductsRequireShipping(productIDs)
	if appErr != nil || requireShipping {
		return nil, appErr
	}

	// check if checkoutInfo
	if checkoutInfo.ShippingAddress == nil {
		return nil, nil
	}

	return a.srv.Shipping.ApplicableShippingMethodsForCheckout(
		checkoutInfo.Checkout,
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
func (s *ServiceCheckout) GetValidCollectionPointsForCheckout(lines model_helper.CheckoutLineInfos, countryCode model.CountryCode, quantityCheck bool) (model.WarehouseSlice, *model_helper.AppError) {
	productIDs := lo.Map(lines.Products(), func(product *model.Product, _ int) string { return product.ID })
	linesRequireShipping, appErr := s.srv.Product.ProductsRequireShipping(productIDs)
	if appErr != nil {
		return nil, appErr
	}
	if !linesRequireShipping {
		return model.WarehouseSlice{}, nil
	}

	if countryCode.IsValid() != nil {
		return model.WarehouseSlice{}, nil
	}

	checkoutLineIDs := lo.Map(lines.CheckoutLines(), func(line *model.CheckoutLine, _ int) string { return line.ID })
	checkoutLines, appErr := s.CheckoutLinesByOption(model_helper.CheckoutLineFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.CheckoutLineWhere.ID.IN(checkoutLineIDs),
		),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return model.WarehouseSlice{}, nil
	}

	var (
		warehouses model.WarehouseSlice
		err        error
	)
	if quantityCheck {
		warehouses, err = s.srv.Store.Warehouse().ApplicableForClickAndCollectCheckoutLines(checkoutLines, countryCode)
	} else {
		warehouses, err = s.srv.Store.Warehouse().ApplicableForClickAndCollectNoQuantityCheck(checkoutLines, countryCode)
	}

	if err != nil {
		return nil, model_helper.NewAppError("GetValidCollectionPointsForCheckout", "app.warehouse.error_finding_warehouses_by_checkout_lines_and_country.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return warehouses, nil
}

func (a *ServiceCheckout) ClearDeliveryMethod(checkoutInfo model_helper.CheckoutInfo) *model_helper.AppError {
	checkout := checkoutInfo.Checkout
	checkout.CollectionPointID.String = nil
	checkout.ShippingMethodID.String = nil

	appErr := a.UpdateCheckoutInfoDeliveryMethod(checkoutInfo, nil)
	if appErr != nil {
		return nil
	}

	_, appErr = a.UpsertCheckouts(nil, model.CheckoutSlice{&checkout})
	return appErr
}

// IsFullyPaid Check if provided payment methods cover the checkout's total amount.
// Note that these payments may not be captured or charged at all.
func (s *ServiceCheckout) IsFullyPaid(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo) (bool, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout
	payments, appErr := s.srv.Payment.PaymentsByOption(model_helper.PaymentFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentWhere.CheckoutID.EQ(model_types.NewNullString(checkout.Token)),
			model.PaymentWhere.IsActive.EQ(true),
		),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return false, appErr
		}
		// ignore not found error
	}

	totalPaid := decimal.Zero
	for _, payment := range payments {
		if payment == nil {
			continue
		}
		totalPaid = totalPaid.Add(payment.Total)
	}
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	checkoutTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return false, appErr
	}
	checkoutTotalGiftcardBalance, appErr := s.CheckoutTotalGiftCardsBalance(checkout)
	if appErr != nil {
		return false, appErr
	}

	sub, err := checkoutTotal.Sub(checkoutTotalGiftcardBalance)
	if err != nil {
		return false, model_helper.NewAppError("IsFullyPaid", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	checkoutTotal = sub

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkout.Currency.String())
	if checkoutTotal.LessThan(*zeroTaxedMoney) {
		checkoutTotal = zeroTaxedMoney
	}

	return checkoutTotal.Gross.Amount.LessThan(totalPaid), nil
}

// CancelActivePayments set all active payments belong to given checkout
func (a *ServiceCheckout) CancelActivePayments(checkout *model.Checkout) *model_helper.AppError {
	err := a.srv.Store.Payment().CancelActivePaymentsOfCheckout(checkout.Token)
	if err != nil {
		return model_helper.NewAppError("CancelActivePayments", "app.checkout.cancel_payments_of_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceCheckout) ValidateVariantsInCheckoutLines(lines model_helper.CheckoutLineInfos) *model_helper.AppError {
	var notAvailableVariantIDs util.AnyArray[string]
	for _, line := range lines {
		if line.ChannelListing.PriceAmount.IsNil() {
			notAvailableVariantIDs = append(notAvailableVariantIDs, line.Variant.ID)
		}
	}

	if len(notAvailableVariantIDs) > 0 {
		notAvailableVariantIDs = notAvailableVariantIDs.Dedup()
		// return error indicate there are some product variants that have no channel listing or channel listing price is null
		return model_helper.NewAppError("ValidateVariantsInCheckoutLines", "app.checkout.cannot_add_lines_with_unavailable_variants.app_error", map[string]any{"variants": strings.Join(notAvailableVariantIDs, ", ")}, "", http.StatusNotAcceptable)
	}

	return nil
}

// PrepareInsufficientStockCheckoutValidationAppError
func (s *ServiceCheckout) PrepareInsufficientStockCheckoutValidationAppError(where string, err model_helper.InsufficientStock) *model_helper.AppError {
	return model_helper.NewAppError(where, "app.checkout.insufficient_stock.app_error", map[string]any{"variants": err.VariantIDs()}, "", http.StatusNotAcceptable)
}

type DeliveryMethod interface {
	*model.ShippingMethod | *model.Warehouse
}

// Check if current shipping method is valid
func (s *ServiceCheckout) CleanDeliveryMethod(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, method any) (bool, *model_helper.AppError) {
	if method == nil {
		// no shipping method was provided, it is valid
		return true, nil
	}

	productIDs := lo.Map(lines.Products(), func(prd *model.Product, _ int) string { return prd.ID })
	shippingRequired, appErr := s.srv.Product.ProductsRequireShipping(productIDs)
	if appErr != nil {
		return false, appErr
	}

	if !shippingRequired {
		return false, model_helper.NewAppError("CleanDeliveryMethod", "app.checkout.clean_delivery_method.shipping_not_required.app_error", nil, "shipping not required", http.StatusNotAcceptable)
	}

	switch t := any(method).(type) {
	case *model.ShippingMethod:
		if checkoutInfo.ShippingAddress == nil {
			return false, model_helper.NewAppError("CleanDeliveryMethod", "app.checkout.checkout_no_shipping_address.app_error", nil, "cannot choose a shipping method for a checkout without the shipping address", http.StatusNotAcceptable)
		}
		return lo.SomeBy(checkoutInfo.ValidShippingMethods, func(item *model.ShippingMethod) bool { return item != nil && item.ID == t.ID }), nil

	case *model.Warehouse:
		return lo.SomeBy(checkoutInfo.ValidPickupPoints, func(item *model.Warehouse) bool { return item != nil && item.ID == t.ID }), nil

	// this code will never reach
	default:
		return false, nil
	}
}

func (s *ServiceCheckout) UpdateCheckoutShippingMethodIfValid(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) *model_helper.AppError {
	quantity, appErr := s.CalculateCheckoutQuantity(lines)
	if appErr != nil {
		return appErr
	}

	// remove shipping method when empty checkout
	if quantity == 0 {
		appErr := s.ClearDeliveryMethod(checkoutInfo)
		if appErr != nil {
			return appErr
		}
	} else {
		productIDs := lo.Map(lines.Products(), func(prd *model.Product, _ int) string { return prd.ID })
		requireShipping, appErr := s.srv.Product.ProductsRequireShipping(productIDs)
		if appErr != nil {
			return appErr
		}
		if !requireShipping {
			appErr := s.ClearDeliveryMethod(checkoutInfo)
			if appErr != nil {
				return appErr
			}
		}
	}

	isValid, appErr := s.CleanDeliveryMethod(checkoutInfo, lines, checkoutInfo.DeliveryMethodInfo.GetDeliveryMethod())
	if appErr != nil {
		return appErr
	}

	if !isValid {
		appErr = s.ClearDeliveryMethod(checkoutInfo)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceCheckout) CheckLinesQuantity(variants model.ProductVariantSlice, quantities []int, country model.CountryCode, channelSlug string, allowZeroQuantity bool, existingLines model_helper.CheckoutLineInfos, replace bool) *model_helper.AppError {
	for _, quantity := range quantities {
		if (!allowZeroQuantity && quantity <= 0) || (allowZeroQuantity && quantity < 0) {
			return model_helper.NewAppError("CheckLinesQuantity", "app.checkout.zero_quantity_not_allowed.app_error", nil, "quantity must be heigher than zero", http.StatusNotAcceptable)
		}

		shopMaxCheckoutQuantity := *s.srv.Config().ShopSettings.MaxCheckoutLineQuantity
		if quantity > shopMaxCheckoutQuantity {
			return model_helper.NewAppError("CheckLinesQuantity", "app.checkout.quantity_exceed_max_allowed.app_error", nil, fmt.Sprintf("cannot add more than %d items", shopMaxCheckoutQuantity), http.StatusNotAcceptable)
		}
	}

	insufficientStockErr, appErr := s.srv.Warehouse.CheckStockAndPreorderQuantityBulk(variants, country, quantities, channelSlug, nil, existingLines, replace)
	if appErr != nil {
		return appErr
	}
	if insufficientStockErr != nil {
		errors := make([]string, len(insufficientStockErr.Items))
		for idx, item := range insufficientStockErr.Items {
			errors[idx] = fmt.Sprintf("could not add items %s. Only %d remainning in stock.", model_helper.ProductVariantString(item.Variant), max(*item.AvailableQuantity, 0))
		}
		return model_helper.NewAppError("CheckLinesQuantity", "app.checkout.insufficient_stock.app_error", map[string]any{"Quantity": errors}, insufficientStockErr.Error(), http.StatusNotAcceptable)
	}

	return nil
}
