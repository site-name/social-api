package discount

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// IncreaseVoucherUsage increase voucher's uses by 1
func (a *ServiceDiscount) IncreaseVoucherUsage(voucher *model.Voucher) *model.AppError {
	voucher.Used++
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

// DecreaseVoucherUsage decreases voucher's uses by 1
func (a *ServiceDiscount) DecreaseVoucherUsage(voucher *model.Voucher) *model.AppError {
	voucher.Used--
	_, appErr := a.UpsertVoucher(voucher)
	return appErr
}

// AddVoucherUsageByCustomer adds an usage for given voucher, by given customer
func (a *ServiceDiscount) AddVoucherUsageByCustomer(voucher *model.Voucher, customerEmail string) (*model.NotApplicable, *model.AppError) {
	_, appErr := a.VoucherCustomerByOptions(&model.VoucherCustomerFilterOption{
		VoucherID:     squirrel.Eq{store.VoucherCustomerTableName + ".VoucherID": voucher.Id},
		CustomerEmail: squirrel.Eq{store.VoucherCustomerTableName + ".CustomerEmail": customerEmail},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		// create new voucher customer
		_, appErr = a.CreateNewVoucherCustomer(voucher.Id, customerEmail)
		if appErr != nil {
			return model.NewNotApplicable("AddVoucherUsageByCustomer", "Offer only valid once per customer", nil, 0), nil
		}
	}

	return nil, nil
}

// RemoveVoucherUsageByCustomer deletes voucher customers for given voucher
func (a *ServiceDiscount) RemoveVoucherUsageByCustomer(voucher *model.Voucher, customerEmail string) *model.AppError {
	err := a.srv.Store.VoucherCustomer().DeleteInBulk(&model.VoucherCustomerFilterOption{
		VoucherID:     squirrel.Eq{store.VoucherCustomerTableName + ".VoucherID": voucher.Id},
		CustomerEmail: squirrel.Eq{store.VoucherCustomerTableName + ".CustomerEmail": customerEmail},
	})
	if err != nil {
		return model.NewAppError("RemoveVoucherUsageByCustomer", "app.discount.error_delating_voucher_customer_relations.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// GetProductDiscountOnSale Return discount value if product is on sale or raise NotApplicable
func (a *ServiceDiscount) GetProductDiscountOnSale(product model.Product, productCollectionIDs []string, discountInfo *model.DiscountInfo, channeL model.Channel, variantID string) (types.DiscountCalculator, *model.AppError) {
	// this checks whether the given product is on sale
	isProductOnSale := util.ItemInSlice(product.Id, discountInfo.ProductIDs) ||
		(product.CategoryID != nil && util.ItemInSlice(*product.CategoryID, discountInfo.CategoryIDs)) ||
		len(util.SlicesIntersection(productCollectionIDs, discountInfo.CollectionIDs)) > 0

	isVariantOnSale := model.IsValidId(variantID) && util.ItemInSlice(variantID, discountInfo.VariantsIDs)

	if isProductOnSale || isVariantOnSale {
		switch t := discountInfo.Sale.(type) {
		case *model.Sale:
			return a.GetSaleDiscount(t, discountInfo.ChannelListings[channeL.Slug])
		case *model.Voucher:
			return a.GetVoucherDiscount(t, channeL.Id)
		}
	}

	return nil, model.NewAppError("GetProductDiscountOnSale", "app.discount.discount_not_applicable_for_product.app_error", nil, "", http.StatusNotAcceptable)
}

// GetProductDiscounts Return discount values for all discounts applicable to a product.
func (a *ServiceDiscount) GetProductDiscounts(product model.Product, collections []*model.Collection, discountInfos []*model.DiscountInfo, channeL model.Channel, variantID string) ([]types.DiscountCalculator, *model.AppError) {
	// filter duplicate collections
	var (
		uniqueCollectionIDs = []string{}
		meetMap             = map[string]bool{}
		wg                  sync.WaitGroup
		mut                 sync.Mutex
	)

	for _, collection := range collections {
		if _, met := meetMap[collection.Id]; !met {
			uniqueCollectionIDs = append(uniqueCollectionIDs, collection.Id)
			meetMap[collection.Id] = true
		}
	}

	wg.Add(len(uniqueCollectionIDs))

	var (
		appError                    *model.AppError
		discountCalculatorFunctions []types.DiscountCalculator
	)

	for _, discountInfo := range discountInfos {
		go func(info *model.DiscountInfo) {
			discountCalFunc, appErr := a.GetProductDiscountOnSale(product, uniqueCollectionIDs, info, channeL, variantID)

			mut.Lock()
			if appErr != nil && appError == nil {
				appError = appErr
			} else {
				discountCalculatorFunctions = append(discountCalculatorFunctions, discountCalFunc)
			}
			mut.Unlock()

			wg.Done()

		}(discountInfo)
	}

	wg.Wait()

	if appError != nil {
		return nil, appError
	}

	return discountCalculatorFunctions, nil
}

// CalculateDiscountedPrice Return minimum product's price of all prices with discounts applied
//
// `discounts` is optional
func (a *ServiceDiscount) CalculateDiscountedPrice(product model.Product, price *goprices.Money, collections []*model.Collection, discounts []*model.DiscountInfo, channeL model.Channel, variantID string) (*goprices.Money, *model.AppError) {
	if len(discounts) > 0 {

		discountCalFuncs, appErr := a.GetProductDiscounts(product, collections, discounts, channeL, variantID)
		if appErr != nil {
			return nil, appErr
		}

		for _, discountFunc := range discountCalFuncs {
			discountedIface, err := discountFunc(price, nil)
			if err != nil {
				return nil, model.NewAppError("CalculateDiscountedPrice", "app.discount.calculate_discount_error.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
			discountedPrice := discountedIface.(*goprices.Money)
			if discountedPrice.LessThan(price) {
				price = discountedPrice
			}
		}
	}

	return price, nil
}

// ValidateVoucherForCheckout validates given voucher
func (a *ServiceDiscount) ValidateVoucherForCheckout(manager interfaces.PluginManagerInterface, voucher *model.Voucher, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, discounts []*model.DiscountInfo) (*model.NotApplicable, *model.AppError) {
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

func (a *ServiceDiscount) ValidateVoucherInOrder(ord *model.Order) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {

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

func (a *ServiceDiscount) ValidateVoucher(voucher *model.Voucher, totalPrice *goprices.TaxedMoney, quantity int, customerEmail string, channelID string, customerID string) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
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
func (a *ServiceDiscount) GetProductsVoucherDiscount(voucher *model.Voucher, prices []*goprices.Money, channelID string) (*goprices.Money, *model.AppError) {
	// validate given prices are valid:
	var (
		invalidArg   bool
		minPrice     *goprices.Money
		appErrDetail string
		mut          sync.Mutex
		wg           sync.WaitGroup
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
		if minPrice == nil || minPrice.Amount.GreaterThan(price.Amount) {
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
			mut.Lock()
			if appErr == nil {
				appErr = err
			}
			mut.Unlock()
		}
	}

	wg.Add(len(prices))

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
					mut.Lock()
					totalAmount = addedAmount
					mut.Unlock()
				}
			}

		}(price)
	}

	wg.Wait()

	if appErr != nil {
		return nil, appErr
	}

	return totalAmount, nil
}

// FetchCategories returns a map with keys are sale ids, values are slices of category ids
func (a *ServiceDiscount) FetchCategories(saleIDs []string) (map[string][]string, *model.AppError) {
	saleCategories, appErr := a.SaleCategoriesByOption(&model.SaleCategoryRelationFilterOption{
		SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": saleIDs},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return map[string][]string{}, nil
	}

	// categoryMap has keys are sale ids, values are slices of category ids
	var categoryMap = map[string][]string{}
	for _, relation := range saleCategories {
		categoryMap[relation.SaleID] = append(categoryMap[relation.SaleID], relation.CategoryID)
	}

	allCategories, appErr := a.srv.ProductService().CategoriesByOption(&model.CategoryFilterOption{
		All: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		return map[string][]string{}, nil
	}

	categorizedCategories := model.ClassifyCategories(allCategories)

	var subCategoriesMap = map[string][]string{}
	for saleID, categoryIDs := range categoryMap {
		subCategoriesMap[saleID] = util.SlicesIntersection(categoryIDs, categorizedCategories.IDs())
	}

	return subCategoriesMap, nil
}

// FetchCollections returns a map with keys are sale ids, values are slices of UNIQUE collection ids
func (a *ServiceDiscount) FetchCollections(saleIDs []string) (map[string][]string, *model.AppError) {
	saleCollections, appErr := a.SaleCollectionsByOptions(&model.SaleCollectionRelationFilterOption{
		SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": saleIDs},
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
	saleProducts, appErr := a.SaleProductsByOptions(&model.SaleProductRelationFilterOption{
		SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": saleIDs},
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
	saleProductVariants, appErr := s.SaleProductVariantsByOptions(&model.SaleProductVariantFilterOption{
		SaleID: squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": salePKs},
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
func (a *ServiceDiscount) FetchSaleChannelListings(saleIDs []string) (map[string]map[string]*model.SaleChannelListing, *model.AppError) {
	channelListings, err := a.srv.Store.
		DiscountSaleChannelListing().
		SaleChannelListingsWithOption(&model.SaleChannelListingFilterOption{
			SaleID:               squirrel.Eq{store.SaleChannelListingTableName + ".SaleID": saleIDs},
			SelectRelatedChannel: true,
		})
	if err != nil {
		return nil, model.NewAppError("FetchSaleChannelListings", "app.discount.sale_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	channelListingMap := map[string]map[string]*model.SaleChannelListing{}

	for _, listing := range channelListings {
		channelListingMap[listing.SaleID][listing.GetChannel().Slug] = listing
	}

	return channelListingMap, nil
}

func (a *ServiceDiscount) FetchDiscounts(date time.Time) ([]*model.DiscountInfo, *model.AppError) {
	// finds active sales
	activeSales, apErr := a.ActiveSales(&date)
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
		saleChannelListings map[string]map[string]*model.SaleChannelListing
		mut                 sync.Mutex
		wg                  sync.WaitGroup
	)

	safelySetAppError := func(err *model.AppError) {
		mut.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		mut.Unlock()
	}

	wg.Add(5)

	go func() {
		collections, apErr = a.FetchCollections(activeSaleIDs)
		safelySetAppError(apErr)

		wg.Done()
	}()

	go func() {
		saleChannelListings, apErr = a.FetchSaleChannelListings(activeSaleIDs)
		safelySetAppError(apErr)

		wg.Done()
	}()

	go func() {
		products, apErr = a.FetchProducts(activeSaleIDs)
		safelySetAppError(apErr)

		wg.Done()
	}()

	go func() {
		categories, apErr = a.FetchCategories(activeSaleIDs)
		safelySetAppError(apErr)

		wg.Done()
	}()

	go func() {
		variants, apErr = a.FetchVariants(activeSaleIDs)
		safelySetAppError(apErr)

		wg.Done()
	}()

	wg.Wait()

	if appError != nil {
		return nil, appError
	}

	var discountInfos []*model.DiscountInfo

	for _, sale := range activeSales {
		discountInfos = append(discountInfos, &model.DiscountInfo{
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
func (a *ServiceDiscount) FetchActiveDiscounts() ([]*model.DiscountInfo, *model.AppError) {
	return a.FetchDiscounts(time.Now().UTC())
}

// FetchCatalogueInfo may return a map with keys are ["categories", "collections", "products", "variants"].
//
// values are slices of uuid strings
func (s *ServiceDiscount) FetchCatalogueInfo(instance model.Sale) (map[string][]string, *model.AppError) {

	var (
		wg       sync.WaitGroup
		mut      sync.Mutex
		appError *model.AppError

		categorieIDs  []string
		collectionIDs []string
		productIDs    []string
		variantIDs    []string

		setAppErr = func(err *model.AppError, setWhenCode ...int) {
			mut.Lock()
			defer mut.Unlock()

			if err != nil && appError == nil && util.ItemInSlice(err.StatusCode, setWhenCode) {
				appError = err
			}
		}
	)

	wg.Add(4)

	go func() {
		defer wg.Done()

		cates, appErr := s.srv.ProductService().CategoriesByOption(&model.CategoryFilterOption{
			SaleID: squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": instance.Id},
		})
		if appErr != nil {
			setAppErr(appErr, http.StatusInternalServerError)
			return
		}
		categorieIDs = model.Categories(cates).IDs()
	}()

	go func() {
		defer wg.Done()

		collecs, appErr := s.srv.ProductService().CollectionsByOption(&model.CollectionFilterOption{
			SaleID: squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": instance.Id},
		})
		if appErr != nil {
			setAppErr(appErr, http.StatusInternalServerError)
			return
		}
		collectionIDs = model.Collections(collecs).IDs()
	}()

	go func() {
		defer wg.Done()

		prds, appErr := s.srv.ProductService().ProductsByOption(&model.ProductFilterOption{
			SaleID: squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": instance.Id},
		})
		if appErr != nil {
			setAppErr(appErr, http.StatusInternalServerError)
			return
		}
		productIDs = model.Products(prds).IDs()
	}()

	go func() {
		defer wg.Done()

		productVariants, appErr := s.srv.ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{})
		if appErr != nil {
			setAppErr(appErr, http.StatusInternalServerError)
			return
		}
		variantIDs = model.ProductVariants(productVariants).IDs()
	}()

	wg.Wait()

	if appError != nil {
		return nil, appError
	}

	return map[string][]string{
		"categories":  categorieIDs,
		"collections": collectionIDs,
		"products":    productIDs,
		"variants":    variantIDs,
	}, nil
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
