package discount

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *ServiceDiscount) AlterVoucherUsage(voucher model.Voucher, usageDelta int) (*model.Voucher, *model_helper.AppError) {
	voucher.Used += usageDelta
	return a.UpsertVoucher(voucher)
}

func (a *ServiceDiscount) AddVoucherUsageByCustomer(voucher model.Voucher, customerEmail string) (*model_helper.NotApplicable, *model_helper.AppError) {
	_, appErr := a.VoucherCustomerByOptions(model_helper.VoucherCustomerFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherCustomerWhere.VoucherID.EQ(voucher.ID),
			model.VoucherCustomerWhere.CustomerEmail.EQ(customerEmail),
		),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// create new voucher customer
		_, appErr = a.CreateNewVoucherCustomer(voucher.ID, customerEmail)
		if appErr != nil {
			return model_helper.NewNotApplicable("AddVoucherUsageByCustomer", "Offer only valid once per customer", nil, 0), nil
		}
	}

	return nil, nil
}

func (a *ServiceDiscount) RemoveVoucherUsageByCustomer(voucher model.Voucher, customerEmail string) *model_helper.AppError {
	voucherCustomers, appErr := a.VoucherCustomersByOption(model_helper.VoucherCustomerFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherCustomerWhere.VoucherID.EQ(voucher.ID),
			model.VoucherCustomerWhere.CustomerEmail.EQ(customerEmail),
		),
	})
	if appErr != nil {
		return appErr
	}

	ids := lo.Map(voucherCustomers, func(vc *model.VoucherCustomer, _ int) string { return vc.ID })

	err := a.srv.Store.VoucherCustomer().Delete(ids)
	if err != nil {
		return model_helper.NewAppError("RemoveVoucherUsageByCustomer", "app.discount.error_delating_voucher_customer_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// GetProductDiscountOnSale Return discount value if product is on sale or raise NotApplicable
func (a *ServiceDiscount) GetProductDiscountOnSale(product model.Product, productCollectionIDs []string, discountInfo model_helper.DiscountInfo, channeL model.Channel, variantID string) (types.DiscountCalculator, *model_helper.AppError) {
	// this checks whether the given product is on sale
	isProductOnSale := discountInfo.ProductIDs.Contains(product.ID) ||
		discountInfo.CategoryIDs.Contains(product.CategoryID) ||
		discountInfo.CollectionIDs.InterSection(productCollectionIDs).Len() > 0

	isVariantOnSale := model_helper.IsValidId(variantID) && discountInfo.VariantsIDs.Contains(variantID)

	if isProductOnSale || isVariantOnSale {
		switch t := discountInfo.Sale.(type) {
		case model.Sale:
			return a.GetSaleDiscount(t, *discountInfo.ChannelListings[channeL.Slug])
		case model.Voucher:
			return a.GetVoucherDiscount(t, channeL.ID)
		}
	}

	return nil, model_helper.NewAppError("GetProductDiscountOnSale", "app.discount.discount_not_applicable_for_product.app_error", nil, "", http.StatusNotAcceptable)
}

// GetProductDiscounts Return discount values for all discounts applicable to a product.
func (a *ServiceDiscount) GetProductDiscounts(product model.Product, collections model.CollectionSlice, discountInfos []*model_helper.DiscountInfo, channeL model.Channel, variantID string) ([]types.DiscountCalculator, *model_helper.AppError) {
	var (
		atomicValue                 atomic.Int32
		appErrChan                  = make(chan *model_helper.AppError)
		valueChan                   = make(chan types.DiscountCalculator)
		discountCalculatorFunctions []types.DiscountCalculator
	)
	defer close(appErrChan)
	defer close(valueChan)
	atomicValue.Add(int32(len(discountInfos)))

	// filter duplicate collections
	uniqueCollectionIDs := lo.Uniq(lo.Map(collections, func(c *model.Collection, _ int) string { return c.ID }))

	for _, discountInfo := range discountInfos {
		go func(info *model_helper.DiscountInfo) {
			defer atomicValue.Add(-1)

			discountCalFunc, appErr := a.GetProductDiscountOnSale(product, uniqueCollectionIDs, *info, channeL, variantID)
			if appErr != nil {
				appErrChan <- appErr
				return
			}

			valueChan <- discountCalFunc
		}(discountInfo)
	}

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		case value := <-valueChan:
			discountCalculatorFunctions = append(discountCalculatorFunctions, value)
		default:
		}
	}

	return discountCalculatorFunctions, nil
}

// CalculateDiscountedPrice Return minimum product's price of all prices with discounts applied
//
// `discounts` is optional
func (a *ServiceDiscount) CalculateDiscountedPrice(product model.Product, price goprices.Money, collections []*model.Collection, discounts []*model_helper.DiscountInfo, channeL model.Channel, variantID string) (*goprices.Money, *model_helper.AppError) {
	if len(discounts) > 0 {

		discountCalFuncs, appErr := a.GetProductDiscounts(product, collections, discounts, channeL, variantID)
		if appErr != nil {
			return nil, appErr
		}

		for _, discountFunc := range discountCalFuncs {
			discountedIface, err := discountFunc(price, nil)
			if err != nil {
				return nil, model_helper.NewAppError("CalculateDiscountedPrice", "app.discount.calculate_discount_error.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			discountedPrice := discountedIface.(goprices.Money)
			if discountedPrice.LessThan(price) {
				price = discountedPrice
			}
		}
	}

	return &price, nil
}

// ValidateVoucherForCheckout validates given voucher
func (a *ServiceDiscount) ValidateVoucherForCheckout(manager interfaces.PluginManagerInterface, voucher model.Voucher, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo) (*model_helper.NotApplicable, *model_helper.AppError) {
	quantity, appErr := a.srv.Checkout.CalculateCheckoutQuantity(lines)
	if appErr != nil {
		return nil, appErr
	}
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}
	checkoutSubTotal, appErr := a.srv.Checkout.CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}
	customerEmail := checkoutInfo.GetCustomerEmail()
	return a.ValidateVoucher(voucher, *checkoutSubTotal, quantity, customerEmail, checkoutInfo.Channel.ID, checkoutInfo.User.ID)
}

func (a *ServiceDiscount) ValidateVoucherInOrder(order *model.Order) (notApplicableErr *model_helper.NotApplicable, appErr *model_helper.AppError) {
	if order.VoucherID.IsNil() {
		return // returns immediately if order has no voucher
	}

	orderSubTotal, appErr := a.srv.Order.OrderSubTotal(order)
	if appErr != nil {
		return
	}
	orderTotalQuantity, appErr := a.srv.Order.OrderTotalQuantity(order.ID)
	if appErr != nil {
		return
	}
	orderCustomerEmail, appErr := a.srv.Order.CustomerEmail(order)
	if appErr != nil {
		return
	}

	voucher, appErr := a.VoucherById(*order.VoucherID.String)
	if appErr != nil {
		return
	}

	// NOTE: orders should have owner when being created
	var orderOwnerId string
	if !order.UserID.IsNil() {
		orderOwnerId = *order.UserID.String
	}

	return a.ValidateVoucher(*voucher, *orderSubTotal, orderTotalQuantity, orderCustomerEmail, order.ChannelID, orderOwnerId)
}

func (a *ServiceDiscount) ValidateVoucher(voucher model.Voucher, totalPrice goprices.TaxedMoney, quantity int, customerEmail string, channelID string, customerID string) (notApplicableErr *model_helper.NotApplicable, appErr *model_helper.AppError) {
	notApplicableErr, appErr = a.ValidateMinSpent(voucher, totalPrice, channelID)
	if appErr != nil || notApplicableErr != nil {
		return
	}

	notApplicableErr = model_helper.VoucherValidateMinCheckoutItemsQuantity(voucher, quantity)
	if appErr != nil || notApplicableErr != nil {
		return
	}
	if voucher.ApplyOncePerCustomer {
		notApplicableErr, appErr = a.ValidateOncePerCustomer(voucher, customerEmail)
		if appErr != nil || notApplicableErr != nil {
			return
		}
	}
	if voucher.OnlyForStaff {
		notApplicableErr, appErr = a.ValidateOnlyForStaff(voucher, customerID)
		if appErr != nil || notApplicableErr != nil {
			return
		}
	}

	return
}

// GetProductsVoucherDiscount Calculate discount value for a voucher of product or category type
func (a *ServiceDiscount) GetProductsVoucherDiscount(voucher model.Voucher, prices []*goprices.Money, channelID string) (*goprices.Money, *model_helper.AppError) {
	// validate params
	if len(prices) == 0 {
		return nil, model_helper.NewAppError("GetProductsVoucherDiscount", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "prices"}, "please provide prices list", http.StatusBadRequest)
	}

	minPrice, _ := util.MinMaxMoneyInMoneySlice(prices)

	if voucher.ApplyOncePerOrder {
		price, appErr := a.GetDiscountAmountFor(voucher, minPrice, channelID)
		if appErr != nil {
			return nil, appErr
		}
		return price.(*goprices.Money), nil
	}

	totalAmount, _ := util.ZeroMoney(prices[0].GetCurrency()) // ignore error since channels's Currencies are validated before saving

	var (
		atomicValue atomic.Int32
		appErrChan  = make(chan *model_helper.AppError)
		valueChan   = make(chan *goprices.Money)
	)
	defer func() {
		close(appErrChan)
		close(valueChan)
	}()

	atomicValue.Add(int32(len(prices)))

	for _, price := range prices {
		go func(aPrice *goprices.Money) {
			defer atomicValue.Add(-1)

			money, appErr := a.GetDiscountAmountFor(voucher, aPrice, channelID)
			if appErr != nil {
				appErrChan <- appErr
				return
			}

			valueChan <- money.(*goprices.Money)
		}(price)
	}

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		case money := <-valueChan:
			addedMoney, err := totalAmount.Add(*money)
			if err != nil {
				return nil, model_helper.NewAppError("GetProductsVoucherDiscount", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
			}
			totalAmount = addedMoney
		default:
		}
	}

	return totalAmount, nil
}

// FetchCategories returns a map with keys are sale ids, values are slices of category ids
func (a *ServiceDiscount) FetchCategories(saleIDs []string) (map[string][]string, *model_helper.AppError) {
	saleCategories, appErr := a.SaleCategoriesByOption(squirrel.Eq{"sale_id": saleIDs})
	if appErr != nil {
		return nil, appErr
	}

	// categoryMap has keys are sale ids, values are slices of category ids
	var categoryMap = map[string]*util.AnySet[string]{}
	for _, relation := range saleCategories {
		if categoryMap[relation.SaleID] == nil {
			categoryMap[relation.SaleID] = util.NewSet[string]()
		}
		categoryMap[relation.SaleID].Add(relation.CategoryID)
	}

	var subCategoryMap = map[string][]string{}
	for saleID, categoryIDs := range categoryMap {
		categories, appErr := a.srv.Product.CategoryByIds(categoryIDs.Values(), true)
		if appErr != nil {
			return nil, appErr
		}

		categoryIDs := lo.Map(categories, func(c *model.Category, _ int) string { return c.ID })
		subCategoryMap[saleID] = util.NewSet[string](categoryIDs...).Values()
	}

	return subCategoryMap, nil
}

// FetchCollections returns a map with keys are sale ids, values are slices of UNIQUE collection ids
func (a *ServiceDiscount) FetchCollections(saleIDs []string) (map[string][]string, *model_helper.AppError) {
	saleCollections, appErr := a.SaleCollectionsByOptions(squirrel.Eq{"sale_id": saleIDs})
	if appErr != nil {
		return nil, appErr
	}

	var (
		collectionMap = map[string][]string{}
		meetMap       = map[string]map[string]bool{}
		saleID        string
		collectionID  string
	)

	for _, saleCollection := range saleCollections {
		saleID = saleCollection.SaleID
		collectionID = saleCollection.CollectionID

		if !meetMap[saleID][collectionID] {
			collectionMap[saleID] = append(collectionMap[saleID], collectionID)
			meetMap[saleID][collectionID] = true
		}
	}

	return collectionMap, nil
}

// FetchProducts returns a map with keys are sale ids, values are slices of UNIQUE product ids
func (a *ServiceDiscount) FetchProducts(saleIDs []string) (map[string][]string, *model_helper.AppError) {
	saleProducts, appErr := a.SaleProductsByOptions(squirrel.Eq{"sale_id": saleIDs})
	if appErr != nil {
		return nil, appErr
	}

	var (
		productMap = map[string][]string{}
		meetMap    = map[string]map[string]bool{}
		saleID     string
		productID  string
	)

	for _, saleProduct := range saleProducts {
		saleID = saleProduct.SaleID
		productID = saleProduct.ProductID

		if !meetMap[saleID][productID] {
			productMap[saleID] = append(productMap[saleID], productID)
			meetMap[saleID][productID] = true
		}
	}

	return productMap, nil
}

// FetchVariants returns a map with keys are sale ids and values are slice of UNIQUE product variant ids
func (s *ServiceDiscount) FetchVariants(saleIDs []string) (map[string][]string, *model_helper.AppError) {
	saleProductVariants, appErr := s.SaleProductVariantsByOptions(squirrel.Eq{"sale_id": saleIDs})
	if appErr != nil {
		return nil, appErr
	}

	var (
		meetMap   = make(map[string]map[string]bool) // keys are sale ids, value maps have keys are product variant ids
		res       = make(map[string][]string)        // keys are sale ids, values are slice of UNIQUE product variant ids
		saleID    string
		variantID string
	)
	for _, relation := range saleProductVariants {
		saleID = relation.SaleID
		variantID = relation.ProductVariantID

		if !meetMap[saleID][variantID] {
			meetMap[saleID][variantID] = true

			res[saleID] = append(res[saleID], variantID)
		}
	}

	return res, nil
}

// FetchSaleChannelListings returns a map with keys are sale ids, values are maps with keys are channel slugs
func (a *ServiceDiscount) FetchSaleChannelListings(saleIDs []string) (map[string]map[string]*model.SaleChannelListing, *model_helper.AppError) {
	channelListings, err := a.srv.Store.
		DiscountSaleChannelListing().
		FilterByOptions(model_helper.SaleChannelListingFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.SaleChannelListingWhere.SaleID.IN(saleIDs),
				qm.Load(model.SaleChannelListingRels.Channel),
			),
		})
	if err != nil {
		return nil, model_helper.NewAppError("FetchSaleChannelListings", "app.discount.sale_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	channelListingMap := map[string]map[string]*model.SaleChannelListing{}

	for _, listing := range channelListings {
		channelListingMap[listing.SaleID][listing.R.Channel.Slug] = listing
	}

	return channelListingMap, nil
}

func (a *ServiceDiscount) FetchDiscounts(date time.Time) ([]*model_helper.DiscountInfo, *model_helper.AppError) {
	// finds active sales
	activeSales, appErr := a.ActiveSales(&date)
	if appErr != nil {
		return nil, appErr
	}

	activeSaleIDs := lo.Map(activeSales, func(s *model.Sale, _ int) string { return s.ID })

	var (
		collections         map[string][]string
		products            map[string][]string
		categories          map[string][]string
		variants            map[string][]string
		saleChannelListings map[string]map[string]*model.SaleChannelListing
		atomicValue         atomic.Int32
		appErrChan          = make(chan *model_helper.AppError)
	)
	defer close(appErrChan)
	atomicValue.Add(5) //

	go func() {
		defer atomicValue.Add(-1)
		clts, apErr := a.FetchCollections(activeSaleIDs)
		if apErr != nil {
			appErrChan <- apErr
		}

		collections = clts
	}()

	go func() {
		defer atomicValue.Add(-1)
		scls, apErr := a.FetchSaleChannelListings(activeSaleIDs)
		if apErr != nil {
			appErrChan <- apErr
		}

		saleChannelListings = scls
	}()

	go func() {
		defer atomicValue.Add(-1)
		prds, apErr := a.FetchProducts(activeSaleIDs)
		if apErr != nil {
			appErrChan <- apErr
		}

		products = prds
	}()

	go func() {
		defer atomicValue.Add(-1)
		ctgrs, apErr := a.FetchCategories(activeSaleIDs)
		if apErr != nil {
			appErrChan <- apErr
		}

		categories = ctgrs
	}()

	go func() {
		defer atomicValue.Add(-1)
		vras, apErr := a.FetchVariants(activeSaleIDs)
		if apErr != nil {
			appErrChan <- apErr
		}

		variants = vras
	}()

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		default:
		}
	}

	var discountInfos []*model_helper.DiscountInfo

	for _, sale := range activeSales {
		discountInfos = append(discountInfos, &model_helper.DiscountInfo{
			Sale:            sale,
			CategoryIDs:     categories[sale.ID],
			ChannelListings: saleChannelListings[sale.ID],
			CollectionIDs:   collections[sale.ID],
			ProductIDs:      products[sale.ID],
			VariantsIDs:     variants[sale.ID],
		})
	}

	return discountInfos, nil
}

// FetchActiveDiscounts returns discounts that are activated
func (a *ServiceDiscount) FetchActiveDiscounts() ([]*model_helper.DiscountInfo, *model_helper.AppError) {
	return a.FetchDiscounts(time.Now().UTC())
}

// FetchCatalogueInfo may return a map with keys are ["categories", "collections", "products", "variants"].
//
// values are slices of uuid strings
func (s *ServiceDiscount) FetchCatalogueInfo(instance model.Sale) (map[string][]string, *model_helper.AppError) {
	var (
		res        = map[string][]string{}
		appError   = make(chan *model_helper.AppError)
		val        = make(chan any)
		atmicValue atomic.Int32
	)
	defer func() {
		close(appError)
		close(val)
	}()

	atmicValue.Add(4)

	go func() {
		categories, appErr := s.srv.Product.CategoriesByOption(model_helper.CategoryFilterOption{
			SaleID: model.SaleCategoryWhere.SaleID.EQ(instance.ID),
		})
		if appErr != nil {
			appError <- appErr
			return
		}
		val <- categories
	}()

	go func() {
		_, collections, appErr := s.srv.Product.CollectionsByOption(model_helper.CollectionFilterOptions{
			SaleID: model.SaleCollectionWhere.SaleID.EQ(instance.ID),
		})
		if appErr != nil {
			appError <- appErr
			return
		}
		val <- collections
	}()

	go func() {
		products, appErr := s.srv.Product.ProductsByOption(model_helper.ProductFilterOption{
			SaleID: model.SaleProductWhere.SaleID.EQ(instance.ID),
		})
		if appErr != nil {
			appError <- appErr
			return
		}
		val <- products
	}()

	go func() {
		productVariants, appErr := s.srv.Product.ProductVariantsByOption(model_helper.ProductVariantFilterOptions{
			SaleID: model.SaleProductVariantWhere.SaleID.EQ(instance.ID),
		})
		if appErr != nil {
			appError <- appErr
			return
		}
		val <- productVariants
	}()

	for atmicValue.Load() != 0 {
		select {
		case err := <-appError:
			return nil, err

		case v := <-val:
			atmicValue.Add(-1)

			switch t := v.(type) {
			case model.CategorySlice:
				res["categories"] = lo.Map(t, func(c *model.Category, _ int) string { return c.ID })
			case model.ProductSlice:
				res["products"] = lo.Map(t, func(p *model.Product, _ int) string { return p.ID })
			case model.ProductVariantSlice:
				res["variants"] = lo.Map(t, func(v *model.ProductVariant, _ int) string { return v.ID })
			case model.CollectionSlice:
				res["collections"] = lo.Map(t, func(c *model.Collection, _ int) string { return c.ID })
			}
		}
	}

	return res, nil
}

// IsValidPromoCode checks if given code is valid giftcard code or voucher code
func (s *ServiceDiscount) IsValidPromoCode(code string) bool {
	codeIsGiftcard, appErr := s.srv.Giftcard.PromoCodeIsGiftCard(code)
	if appErr != nil {
		s.srv.Log.Error("IsValidPromoCode", slog.Err(appErr))
		return false
	}

	codeIsVoucher, appErr := s.PromoCodeIsVoucher(code)
	if appErr != nil {
		s.srv.Log.Error("IsValidPromoCode", slog.Err(appErr))
	}

	return !(codeIsGiftcard || codeIsVoucher)
}
