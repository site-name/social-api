package discount

import (
	"fmt"
	"net/http"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// IncreaseVoucherUsage increase voucher's uses by 1
func (a *AppDiscount) IncreaseVoucherUsage(voucher *product_and_discount.Voucher) *model.AppError {
	voucher.Used++
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

// DecreaseVoucherUsage decreases voucher's uses by 1
func (a *AppDiscount) DecreaseVoucherUsage(voucher *product_and_discount.Voucher) *model.AppError {
	voucher.Used--
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

func (a *AppDiscount) AddVoucherUsageByCustomer(voucher *product_and_discount.Voucher, customerEmail string) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "AddVoucherUsageByCustomer"
		}
		if notApplicableErr != nil {
			notApplicableErr.Where = "AddVoucherUsageByCustomer"
		}
	}()

	// validate email argument
	if !model.IsValidEmail(customerEmail) {
		appErr = model.NewAppError("AddVoucherUsageByCustomer", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "customer email"}, "", http.StatusBadRequest)
		return
	}

	notApplicableErr, appErr = a.ValidateOncePerCustomer(voucher, customerEmail)
	if appErr != nil || notApplicableErr != nil {
		return
	}

	_, appErr = a.CreateNewVoucherCustomer(voucher.Id, customerEmail)
	return
}

func (a *AppDiscount) RemoveVoucherUsageByCustomer(voucher *product_and_discount.Voucher, customerEmail string) *model.AppError {
	// validate email argument
	if !model.IsValidEmail(customerEmail) {
		return model.NewAppError("RemoveVoucherUsageByCustomer", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "customer email"}, "", http.StatusBadRequest)
	}

	voucherCustomers, appErr := a.VoucherCustomerByCustomerEmailAndVoucherID(voucher.Id, customerEmail)
	if appErr != nil {
		return appErr
	}

	if len(voucherCustomers) > 0 {
		err := a.Srv().Store.VoucherCustomer().DeleteInBulk(voucherCustomers)
		if err != nil {
			return model.NewAppError("RemoveVoucherUsageByCustomer", "app.discount.error_delating_voucher_customer_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}

// GetProductDiscountOnSale Return discount value if product is on sale or raise NotApplicable
func (a *AppDiscount) GetProductDiscountOnSale(product *product_and_discount.Product, productCollectionIDs []string, discountInfo *product_and_discount.DiscountInfo, channeL *channel.Channel) (DiscountCalculator, *model.AppError) {
	// this checks whether the given product is on sale
	if util.StringInSlice(product.Id, discountInfo.ProductIDs) ||
		(product.CategoryID != nil && util.StringInSlice(*product.CategoryID, discountInfo.CategoryIDs)) ||
		len(util.StringArrayIntersection(productCollectionIDs, discountInfo.CollectionIDs)) > 0 {

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
func (a *AppDiscount) GetProductDiscounts(product *product_and_discount.Product, collections []*product_and_discount.Collection, discountInfos []*product_and_discount.DiscountInfo, channeL *channel.Channel) ([]DiscountCalculator, *model.AppError) {
	// filter duplicate collections
	uniqueCollectionIDs := []string{}
	meetMap := map[string]bool{}

	for _, collection := range collections {
		if _, met := meetMap[collection.Id]; !met {
			uniqueCollectionIDs = append(uniqueCollectionIDs, collection.Id)
			meetMap[collection.Id] = true
		}
	}

	a.wg.Add(len(uniqueCollectionIDs))

	var (
		appError                    *model.AppError
		discountCalculatorFunctions []DiscountCalculator
	)

	for _, discountInfo := range discountInfos {
		go func(info *product_and_discount.DiscountInfo) {
			discountCalFunc, appErr := a.GetProductDiscountOnSale(product, uniqueCollectionIDs, info, channeL)

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
func (a *AppDiscount) CalculateDiscountedPrice(product *product_and_discount.Product, price *goprices.Money, collections []*product_and_discount.Collection, discounts []*product_and_discount.DiscountInfo, channeL *channel.Channel) (*goprices.Money, *model.AppError) {
	if len(discounts) > 0 {

		discountCalFuncs, appErr := a.GetProductDiscounts(product, collections, discounts, channeL)
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

func (a *AppDiscount) ValidateVoucherForCheckout() {
	panic("not implemented")
}

func (a *AppDiscount) ValidateVoucherInOrder(ord *order.Order) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
	defer func() {
		if notApplicableErr != nil {
			notApplicableErr.Where = "ValidateVoucherInOrder"
		}
		if appErr != nil {
			appErr.Where = "ValidateVoucherInOrder"
		}
	}()

	if ord.VoucherID == nil {
		return // returns immediately if order has no voucher
	}

	orderSubTotal, appErr := a.OrderApp().OrderSubTotal(ord)
	if appErr != nil {
		return
	}
	orderTotalQuantity, appErr := a.OrderApp().OrderTotalQuantity(ord.Id)
	if appErr != nil {
		return
	}
	orderCustomerEmail, appErr := a.OrderApp().CustomerEmail(ord)
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

func (a *AppDiscount) ValidateVoucher(voucher *product_and_discount.Voucher, totalPrice *goprices.TaxedMoney, quantity int, customerEmail string, channelID string, customerID string) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "ValidateVoucher"
		}
		if notApplicableErr != nil {
			notApplicableErr.Where = "ValidateVoucher"
		}
	}()

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
func (a *AppDiscount) GetProductsVoucherDiscount(voucher *product_and_discount.Voucher, prices []*goprices.Money, channelID string) (*goprices.Money, *model.AppError) {
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

func (a *AppDiscount) FetchCategories(saleIDs []string) (map[string][]string, *model.AppError) {
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

func (a *AppDiscount) FetchCollections(saleIDs []string) (map[string][]string, *model.AppError) {
	saleCollections, err := a.Srv().Store.SaleCollectionRelation().FilterByOption(&product_and_discount.SaleCollectionRelationFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FetchCollections", "app.discount.error_finding_sale_collections_by_option.app_error", err)
	}

	collectionMap := map[string][]string{}
	meetMap := map[string]bool{}

	for _, saleCollection := range saleCollections {
		if _, met := meetMap[saleCollection.CollectionID]; !met {
			collectionMap[saleCollection.SaleID] = append(collectionMap[saleCollection.SaleID], saleCollection.CollectionID)
			meetMap[saleCollection.CollectionID] = true
		}
	}

	return collectionMap, nil
}

func (a *AppDiscount) FetchProducts(saleIDs []string) (map[string][]string, *model.AppError) {
	saleProducts, err := a.Srv().Store.SaleProductRelation().SaleProductsByOption(&product_and_discount.SaleProductRelationFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FetchProducts", "app.discount,error_finding_sale_products_by_option.app_error", err)
	}

	productMap := map[string][]string{}
	meetMap := map[string]bool{}

	for _, saleProduct := range saleProducts {
		if _, met := meetMap[saleProduct.ProductID]; !met {
			productMap[saleProduct.SaleID] = append(productMap[saleProduct.SaleID], saleProduct.ProductID)
			meetMap[saleProduct.ProductID] = true
		}
	}

	return productMap, nil
}

func (a *AppDiscount) FetchSaleChannelListings(saleIDs []string) (map[string]map[string]*product_and_discount.SaleChannelListing, *model.AppError) {
	channelListings, err := a.Srv().Store.DiscountSaleChannelListing().SaleChannelListingsWithOption(&product_and_discount.SaleChannelListingFilterOption{
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

func (a *AppDiscount) FetchDiscounts(date *time.Time) ([]*product_and_discount.DiscountInfo, *model.AppError) {
	// finds active sales
	activeSales, apErr := a.ActiveSales(date)
	if apErr != nil {
		return nil, apErr
	}

	activeSaleIDs := []string{}
	for _, sale := range activeSales {
		activeSaleIDs = append(activeSaleIDs, sale.Id)
	}

	var (
		collections         map[string][]string
		products            map[string][]string
		categories          map[string][]string
		appError            *model.AppError
		saleChannelListings map[string]map[string]*product_and_discount.SaleChannelListing
	)

	safelySetAppError := func(err *model.AppError) {
		if err != nil {
			a.mutex.Lock()
			if appError == nil {
				appError = err
			}
			a.mutex.Unlock()
		}
	}

	a.wg.Add(4)

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
		})
	}

	return discountInfos, nil
}

func (a *AppDiscount) FetchActiveDiscounts() ([]*product_and_discount.DiscountInfo, *model.AppError) {
	return a.FetchDiscounts(util.NewTime(time.Now().UTC()))
}
