package discount

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
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

func (a *AppDiscount) AddVoucherUsageByCustomer(voucher *product_and_discount.Voucher, customerEmail string) *model.AppError {
	// validate email argument
	if !model.IsValidEmail(customerEmail) {
		return model.NewAppError("AddVoucherUsageByCustomer", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "customer email"}, "", http.StatusBadRequest)
	}

	appErr := a.ValidateOncePerCustomer(voucher, customerEmail)
	if appErr != nil {
		return appErr
	}

	_, appErr = a.CreateNewVoucherCustomer(voucher.Id, customerEmail)
	return appErr
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

	var appErr *model.AppError
	discountFuncs := []DiscountCalculator{}

	for _, discountInfo := range discountInfos {
		go func(info *product_and_discount.DiscountInfo) {
			cal, err := a.GetProductDiscountOnSale(product, uniqueCollectionIDs, discountInfo, channeL)

			a.mutex.Lock()
			if err != nil && appErr == nil {
				appErr = err
			} else {
				discountFuncs = append(discountFuncs, cal)
			}
			a.mutex.Unlock()
			a.wg.Done()

		}(discountInfo)
	}

	a.wg.Wait()

	if appErr != nil {
		return nil, appErr
	}

	return discountFuncs, nil
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
				return nil, model.NewAppError("CalculateDiscountedPrice", "app.discount.error_comparing_money.app_errir", nil, err.Error(), http.StatusBadRequest)
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

func (a *AppDiscount) ValidateVoucherInOrder(ord *order.Order) *model.AppError {
	if ord.VoucherID == nil {
		return nil // returns immediately if order has no voucher
	}

	orderSubTotal, appErr := a.OrderApp().OrderSubTotal(ord.Id, ord.Currency)
	if appErr != nil {
		return appErr
	}
	orderTotalQuantity, appErr := a.OrderApp().OrderTotalQuantity(ord.Id)
	if appErr != nil {
		return appErr
	}
	orderCustomerEmail, appErr := a.OrderApp().CustomerEmail(ord)
	if appErr != nil {
		return appErr
	}

	voucher, appErr := a.VoucherById(*ord.VoucherID)
	if appErr != nil {
		return appErr
	}

	// NOTE: orders should have owner when being created
	var orderOwnerId string
	if ord.UserID != nil {
		orderOwnerId = *ord.UserID
	}

	return a.ValidateVoucher(voucher, orderSubTotal, orderTotalQuantity, orderCustomerEmail, ord.ChannelID, orderOwnerId)
}

func (a *AppDiscount) ValidateVoucher(voucher *product_and_discount.Voucher, totalPrice *goprices.TaxedMoney, quantity uint, customerEmail string, channelID string, customerID string) *model.AppError {
	appErr := a.ValidateMinSpent(voucher, totalPrice, channelID)
	if appErr != nil {
		return appErr
	}
	appErr = voucher.ValidateMinCheckoutItemsQuantity(quantity)
	if appErr != nil {
		return appErr
	}
	if voucher.ApplyOncePerCustomer {
		appErr = a.ValidateOncePerCustomer(voucher, customerEmail)
		if appErr != nil {
			return appErr
		}
	}
	if *voucher.OnlyForStaff {
		appErr = a.ValidateVoucherOnlyForStaff(voucher, customerID)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// GetProductsVoucherDiscount Calculate discount value for a voucher of product or category type
func (a *AppDiscount) GetProductsVoucherDiscount(voucher *product_and_discount.Voucher, prices []*goprices.Money, channeL *channel.Channel) (*goprices.Money, *model.AppError) {
	// validate given prices are valid
	var invalidArg bool
	var minPrice *goprices.Money

	if len(prices) == 0 {
		invalidArg = true
	}
	for _, price := range prices {
		// check if prices's currencies is supported by system, are the same and euqal to given channel's currency
		if _, err := goprices.GetCurrencyPrecision(price.Currency); err != nil || price.Currency != prices[0].Currency || price.Currency != channeL.Currency {
			invalidArg = true
			break
		}
		if minPrice == nil || minPrice.Amount.GreaterThan(*price.Amount) {
			minPrice = price
		}
	}
	if invalidArg {
		return nil, model.NewAppError("GetProductsVoucherDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "prices"}, "", http.StatusBadRequest)
	}

	if voucher.ApplyOncePerOrder {
		return a.GetDiscountAmountFor(voucher, minPrice, channeL.Id)
	}

	totalAmount, _ := util.ZeroMoney(channeL.Currency) // ignore error since channels's Currencies are validated before saving
	resultChan := make(chan *goprices.Money)
	errorChan := make(chan *model.AppError)
	defer func() {
		close(resultChan)
		close(errorChan)
	}()

	for _, price := range prices {
		go func(pr *goprices.Money) {

			money, appErr := a.GetDiscountAmountFor(voucher, pr, channeL.Id)
			if appErr != nil {
				errorChan <- appErr
			} else {
				resultChan <- money
			}

		}(price)
	}

loop:
	for {
		select {
		case res, notClosed := <-resultChan:
			totalAmount, _ = totalAmount.Add(res)
			if !notClosed {
				break loop
			}
		case res := <-errorChan:
			return nil, res
		}
	}

	return totalAmount, nil
}

func (a *AppDiscount) FetchCategories(saleIDs []string) (interface{}, *model.AppError) {
	saleCategories, appErr := a.SaleCategoriesByOption(&product_and_discount.SaleCategoryRelationFilterOption{
		SaleID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: saleIDs,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	categoryMap := map[string][]string{}
	for _, relation := range saleCategories {
		if !util.StringInSlice(relation.CategoryID, categoryMap[relation.SaleID]) {
			categoryMap[relation.SaleID] = append(categoryMap[relation.SaleID], relation.CategoryID)
		}
	}

	subCategoryMap := map[string][]string{}

}
