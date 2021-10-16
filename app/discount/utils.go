package discount

import (
	"fmt"
	"net/http"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// IncreaseVoucherUsage increase voucher's uses by 1
func (a *ServiceDiscount) IncreaseVoucherUsage(voucher *product_and_discount.Voucher) *model.AppError {
	voucher.Used++
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

// DecreaseVoucherUsage decreases voucher's uses by 1
func (a *ServiceDiscount) DecreaseVoucherUsage(voucher *product_and_discount.Voucher) *model.AppError {
	voucher.Used--
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

// AddVoucherUsageByCustomer adds an usage for given voucher, by given customer
func (a *ServiceDiscount) AddVoucherUsageByCustomer(voucher *product_and_discount.Voucher, customerEmail string) (*product_and_discount.NotApplicable, *model.AppError) {
	_, appErr := a.VoucherCustomerByOptions(&product_and_discount.VoucherCustomerFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: voucher.Id,
			},
		},
		CustomerEmail: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: customerEmail,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// create new voucher customer
		_, appErr = a.CreateNewVoucherCustomer(voucher.Id, customerEmail)
		if appErr != nil {
			return product_and_discount.NewNotApplicable("AddVoucherUsageByCustomer", "Offer only valid once per customer", nil, 0), nil
		}
	}

	return nil, nil
}

// RemoveVoucherUsageByCustomer deletes voucher customers for given voucher
func (a *ServiceDiscount) RemoveVoucherUsageByCustomer(voucher *product_and_discount.Voucher, customerEmail string) *model.AppError {
	voucherCustomers, appErr := a.VoucherCustomersByOption(&product_and_discount.VoucherCustomerFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: voucher.Id,
			},
		},
		CustomerEmail: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: customerEmail,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	if len(voucherCustomers) > 0 {
		err := a.srv.Store.VoucherCustomer().DeleteInBulk(voucherCustomers)
		if err != nil {
			return model.NewAppError("RemoveVoucherUsageByCustomer", "app.discount.error_delating_voucher_customer_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}

// GetProductDiscountOnSale Return discount value if product is on sale or raise NotApplicable
func (a *ServiceDiscount) GetProductDiscountOnSale(product *product_and_discount.Product, productCollectionIDs []string, discountInfo *product_and_discount.DiscountInfo, channeL *channel.Channel, variantID string) (types.DiscountCalculator, *model.AppError) {
	// this checks whether the given product is on sale
	isProductOnSale := util.StringInSlice(product.Id, discountInfo.ProductIDs) ||
		(product.CategoryID != nil && util.StringInSlice(*product.CategoryID, discountInfo.CategoryIDs)) ||
		len(util.StringArrayIntersection(productCollectionIDs, discountInfo.CollectionIDs)) > 0

	isVariantOnSale := model.IsValidId(variantID) && util.StringInSlice(variantID, discountInfo.VariantsIDs)

	if isProductOnSale || isVariantOnSale {
		switch t := discountInfo.Sale.(type) {
		case *product_and_discount.Sale:
			return a.GetSaleDiscount(t, discountInfo.ChannelListings[channeL.Slug])
		case *product_and_discount.Voucher:
			return a.GetVoucherDiscount(t, channeL.Id)
		}
	}

	return nil, model.NewAppError("GetProductDiscountOnSale", "app.discount.discount_not_applicable_for_product.app_error", nil, "", http.StatusNotAcceptable)
}

// GetProductDiscounts Return discount values for all discounts applicable to a product.
func (a *ServiceDiscount) GetProductDiscounts(product *product_and_discount.Product, collections []*product_and_discount.Collection, discountInfos []*product_and_discount.DiscountInfo, channeL *channel.Channel, variantID string) ([]types.DiscountCalculator, *model.AppError) {
	// filter duplicate collections
	var (
		uniqueCollectionIDs = []string{}
		meetMap             = map[string]bool{}
	)

	for _, collection := range collections {
		if _, met := meetMap[collection.Id]; !met {
			uniqueCollectionIDs = append(uniqueCollectionIDs, collection.Id)
			meetMap[collection.Id] = true
		}
	}

	a.wg.Add(len(uniqueCollectionIDs))

	var (
		appError                    *model.AppError
		discountCalculatorFunctions []types.DiscountCalculator
	)

	for _, discountInfo := range discountInfos {
		go func(info *product_and_discount.DiscountInfo) {
			discountCalFunc, appErr := a.GetProductDiscountOnSale(product, uniqueCollectionIDs, info, channeL, variantID)

			a.mutex.Lock()
			if appErr != nil && appError == nil {
				appError = appErr
			} else {
				discountCalculatorFunctions = append(discountCalculatorFunctions, discountCalFunc)
			}
			a.mutex.Unlock()

			a.wg.Done()

		}(discountInfo)
	}

	a.wg.Wait()

	if appError != nil {
		return nil, appError
	}

	return discountCalculatorFunctions, nil
}

// CalculateDiscountedPrice Return minimum product's price of all prices with discounts applied
//
// `discounts` is optional
func (a *ServiceDiscount) CalculateDiscountedPrice(product *product_and_discount.Product, price *goprices.Money, collections []*product_and_discount.Collection, discounts []*product_and_discount.DiscountInfo, channeL *channel.Channel, variantID string) (*goprices.Money, *model.AppError) {
	if len(discounts) > 0 {

		discountCalFuncs, appErr := a.GetProductDiscounts(product, collections, discounts, channeL, variantID)
		if appErr != nil {
			return nil, appErr
		}

		for _, discountFunc := range discountCalFuncs {
			discountedIface, err := discountFunc(price)
			if err != nil {
				return nil, model.NewAppError("CalculateDiscountedPrice", "app.discount.calculate_discount_error.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			discountedPrice := discountedIface.(*goprices.Money)
			less, err := discountedPrice.LessThan(price)
			if err != nil {
				return nil, model.NewAppError("CalculateDiscountedPrice", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusBadRequest)
			}
			if less {
				price = discountedPrice
			}
		}
	}

	return price, nil
}

// ValidateVoucherForCheckout validates given voucher
func (a *ServiceDiscount) ValidateVoucherForCheckout(manager interface{}, voucher *product_and_discount.Voucher, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) (*product_and_discount.NotApplicable, *model.AppError) {
	quantity, appErr := a.srv.CheckoutService().CalculateCheckoutQuantity(lines)
	if appErr != nil {
		return nil, appErr
	}
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}
	checkoutSubTotal, appErr := a.srv.CheckoutService().CheckoutSubTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}
	customerEmail := checkoutInfo.GetCustomerEmail()
	return a.ValidateVoucher(voucher, checkoutSubTotal, quantity, customerEmail, checkoutInfo.Channel.Id, checkoutInfo.User.Id)
}

func (a *ServiceDiscount) ValidateVoucherInOrder(ord *order.Order) (notApplicableErr *product_and_discount.NotApplicable, appErr *model.AppError) {

	if ord.VoucherID == nil {
		return // returns immediately if order has no voucher
	}

	orderSubTotal, appErr := a.srv.OrderService().OrderSubTotal(ord)
	if appErr != nil {
		return
	}
	orderTotalQuantity, appErr := a.srv.OrderService().OrderTotalQuantity(ord.Id)
	if appErr != nil {
		return
	}
	orderCustomerEmail, appErr := a.srv.OrderService().CustomerEmail(ord)
	if appErr != nil {
		return
	}

	voucher, appErr := a.VoucherById(*ord.VoucherID)
	if appErr != nil {
		return
	}

	// NOTE: orders should have owner when being created
	var orderOwnerId string
	if ord.UserID != nil {
		orderOwnerId = *ord.UserID
	}

	return a.ValidateVoucher(voucher, orderSubTotal, orderTotalQuantity, orderCustomerEmail, ord.ChannelID, orderOwnerId)
}

func (a *ServiceDiscount) ValidateVoucher(voucher *product_and_discount.Voucher, totalPrice *goprices.TaxedMoney, quantity int, customerEmail string, channelID string, customerID string) (notApplicableErr *product_and_discount.NotApplicable, appErr *model.AppError) {

	notApplicableErr, appErr = a.ValidateMinSpent(voucher, totalPrice, channelID)
	if appErr != nil || notApplicableErr != nil {
		return
	}
	notApplicableErr = voucher.ValidateMinCheckoutItemsQuantity(quantity)
	if appErr != nil || notApplicableErr != nil {
		return
	}
	if voucher.ApplyOncePerCustomer {
		notApplicableErr, appErr = a.ValidateOncePerCustomer(voucher, customerEmail)
		if appErr != nil || notApplicableErr != nil {
			return
		}
	}
	if *voucher.OnlyForStaff {
		notApplicableErr, appErr = a.ValidateVoucherOnlyForStaff(voucher, customerID)
		if appErr != nil || notApplicableErr != nil {
			return
		}
	}

	return
}

// GetProductsVoucherDiscount Calculate discount value for a voucher of product or category type
func (a *ServiceDiscount) GetProductsVoucherDiscount(voucher *product_and_discount.Voucher, prices []*goprices.Money, channelID string) (*goprices.Money, *model.AppError) {
	// validate given prices are valid:
	var (
		invalidArg   bool
		minPrice     *goprices.Money
		appErrDetail string
	)

	if len(prices) == 0 {
		invalidArg = true
		appErrDetail = "len(prices) == 0"
	}

	for index, price := range prices {
		// check if prices's currencies is supported by system and are the same
		if _, err := goprices.GetCurrencyPrecision(price.Currency); err != nil || price.Currency != prices[0].Currency {
			invalidArg = true
			appErrDetail = fmt.Sprintf("a price has invalid currency unit: index: %d, currency unit: %s", index+1, price.Currency)
			break
		}
		if minPrice == nil || minPrice.Amount.GreaterThan(*price.Amount) {
			minPrice = price
		}
	}
	if invalidArg {
		return nil, model.NewAppError("GetProductsVoucherDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "prices"}, appErrDetail, http.StatusBadRequest)
	}

	if voucher.ApplyOncePerOrder {
		price, appErr := a.GetDiscountAmountFor(voucher, minPrice, channelID)
		if appErr != nil {
			return nil, appErr
		}
		return price.(*goprices.Money), nil
	}

	totalAmount, _ := util.ZeroMoney(prices[0].Currency) // ignore error since channels's Currencies are validated before saving
	var appErr *model.AppError

	setAppErr := func(err *model.AppError) {
		if err != nil {
			a.mutex.Lock()
			if appErr == nil {
				appErr = err
			}
			a.mutex.Unlock()
		}
	}

	a.wg.Add(len(prices))

	for _, price := range prices {
		go func(aPrice *goprices.Money) {

			money, err := a.GetDiscountAmountFor(voucher, aPrice, channelID)
			if err != nil {
				setAppErr(err)
			} else {
				addedAmount, err := totalAmount.Add(money.(*goprices.Money))
				if err != nil {
					setAppErr(
						model.NewAppError("GetProductsVoucherDiscount", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError),
					)
				} else {
					a.mutex.Lock()
					totalAmount = addedAmount
					a.mutex.Unlock()
				}
			}

		}(price)
	}

	a.wg.Wait()

	if appErr != nil {
		return nil, appErr
	}

	return totalAmount, nil
}

func (a *ServiceDiscount) FetchCategories(saleIDs []string) (map[string][]string, *model.AppError) {
	// saleCategories, appErr := a.SaleCategoriesByOption(&product_and_discount.SaleCategoryRelationFilterOption{
	// 	SaleID: &model.StringFilter{
	// 		StringOption: &model.StringOption{
	// 			In: saleIDs,
	// 		},
	// 	},
	// })
	// if appErr != nil {
	// 	return nil, appErr
	// }

	// categoryMap := map[string][]string{}
	// for _, relation := range saleCategories {
	// 	if !util.StringInSlice(relation.CategoryID, categoryMap[relation.SaleID]) {
	// 		categoryMap[relation.SaleID] = append(categoryMap[relation.SaleID], relation.CategoryID)
	// 	}
	// }

	// subCategoryMap := map[string][]string{}

	// TODO: implement me
	panic("not implemented")
}

// FetchCollections returns a map with keys are sale ids, values are slices of UNIQUE collection ids
func (a *ServiceDiscount) FetchCollections(saleIDs []string) (map[string][]string, *model.AppError) {
	saleCollections, appErr := a.SaleCollectionsByOptions(&product_and_discount.SaleCollectionRelationFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return make(map[string][]string), nil
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
func (a *ServiceDiscount) FetchProducts(saleIDs []string) (map[string][]string, *model.AppError) {
	saleProducts, appErr := a.SaleProductsByOptions(&product_and_discount.SaleProductRelationFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return make(map[string][]string), nil
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
func (s *ServiceDiscount) FetchVariants(salePKs []string) (map[string][]string, *model.AppError) {
	saleProductVariants, appErr := s.SaleProductVariantsByOptions(&product_and_discount.SaleProductVariantFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: salePKs,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return make(map[string][]string), nil
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
func (a *ServiceDiscount) FetchSaleChannelListings(saleIDs []string) (map[string]map[string]*product_and_discount.SaleChannelListing, *model.AppError) {
	channelListings, err := a.srv.Store.DiscountSaleChannelListing().SaleChannelListingsWithOption(&product_and_discount.SaleChannelListingFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FetchSaleChannelListings", "app.discount.sale_channel_listings_by_option.app_error", err)
	}

	channelListingMap := map[string]map[string]*product_and_discount.SaleChannelListing{}

	for _, listing := range channelListings {
		channelListingMap[listing.SaleID][listing.ChannelSlug] = &listing.SaleChannelListing
	}

	return channelListingMap, nil
}

func (a *ServiceDiscount) FetchDiscounts(date *time.Time) ([]*product_and_discount.DiscountInfo, *model.AppError) {
	// finds active sales
	activeSales, apErr := a.ActiveSales(date)
	if apErr != nil {
		return nil, apErr
	}

	activeSaleIDs := activeSales.IDs()

	var (
		collections         map[string][]string
		products            map[string][]string
		categories          map[string][]string
		variants            map[string][]string
		appError            *model.AppError
		saleChannelListings map[string]map[string]*product_and_discount.SaleChannelListing
	)

	safelySetAppError := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		a.mutex.Unlock()
	}

	a.wg.Add(5)

	go func() {
		// find collections
		collections, apErr = a.FetchCollections(activeSaleIDs)
		safelySetAppError(apErr)

		a.wg.Done()
	}()

	go func() {
		saleChannelListings, apErr = a.FetchSaleChannelListings(activeSaleIDs)
		safelySetAppError(apErr)

		a.wg.Done()
	}()

	go func() {
		products, apErr = a.FetchProducts(activeSaleIDs)
		safelySetAppError(apErr)

		a.wg.Done()
	}()

	go func() {
		categories, apErr = a.FetchCategories(activeSaleIDs)
		safelySetAppError(apErr)

		a.wg.Done()
	}()

	go func() {
		variants, apErr = a.FetchVariants(activeSaleIDs)
		safelySetAppError(apErr)

		a.wg.Done()
	}()

	a.wg.Wait()

	if appError != nil {
		return nil, appError
	}

	var discountInfos []*product_and_discount.DiscountInfo

	for _, sale := range activeSales {
		discountInfos = append(discountInfos, &product_and_discount.DiscountInfo{
			Sale:            sale,
			CategoryIDs:     categories[sale.Id],
			ChannelListings: saleChannelListings[sale.Id],
			CollectionIDs:   collections[sale.Id],
			ProductIDs:      products[sale.Id],
			VariantsIDs:     variants[sale.Id],
		})
	}

	return discountInfos, nil
}

// FetchActiveDiscounts returns discounts that are activated
func (a *ServiceDiscount) FetchActiveDiscounts() ([]*product_and_discount.DiscountInfo, *model.AppError) {
	return a.FetchDiscounts(util.NewTime(time.Now().UTC()))
}

func (s *ServiceDiscount) FetchCatalogueInfo(instance *product_and_discount.Sale) (map[string][]string, *model.AppError) {
	panic("not implemented")
}

// IsValidPromoCode checks if given code is valid giftcard code or voucher code
func (s *ServiceDiscount) IsValidPromoCode(code string) bool {
	codeIsGiftcard, appErr := s.srv.GiftcardService().PromoCodeIsGiftCard(code)
	if appErr != nil {
		s.srv.Log.Error("IsValidPromoCode", slog.Err(appErr))
	}

	codeIsVoucher, appErr := s.PromoCodeIsVoucher(code)
	if appErr != nil {
		s.srv.Log.Error("IsValidPromoCode", slog.Err(appErr))
	}

	return !(codeIsGiftcard || codeIsVoucher)
}

// GeneratePromoCode randomly generate promo code
func (s *ServiceDiscount) GeneratePromoCode() string {
	code := model.NewId()
	for !s.IsValidPromoCode(code) {
		code = model.NewId()
	}

	return code
}
