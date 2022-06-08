package warehouse

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations) not exported
func (a *ServiceWarehouse) getAvailableQuantity(stocks warehouse.Stocks) (int, *model.AppError) {
	if len(stocks) == 0 {
		return 0, nil
	}

	// reference: https://github.com/mirumee/saleor/blob/master/saleor/warehouse/availability.py, function (_get_available_quantity)
	// not sure yet why using SUM(DISTINCT 'quantity') on stocks
	var (
		totalQuantity = 0
		meetMap       = make(map[int]bool)
		stockIDs      = stocks.IDs()
	)

	for _, stock := range stocks {
		if _, met := meetMap[stock.Quantity]; !met {
			totalQuantity += stock.Quantity
			meetMap[stock.Quantity] = true // this replicates SUM(DISTINCT ...)
		}
	}

	allocations, appErr := a.AllocationsByOption(nil, &warehouse.AllocationFilterOption{
		StockID: squirrel.Eq{store.AllocationTableName + ".StockID": stockIDs},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return 0, appErr
		}
		allocations = []*warehouse.Allocation{}
	}

	var quantityAllocated int
	for _, allocation := range allocations {
		quantityAllocated += allocation.QuantityAllocated
	}

	return util.Max(totalQuantity-quantityAllocated, 0), nil
}

// CheckStockAndPreorderQuantity Validate if there is stock/preorder available for given variant.
// :raises InsufficientStock: when there is not enough items in stock for a variant
// or there is not enough available preorder items for a variant.
func (s *ServiceWarehouse) CheckStockAndPreorderQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*exception.InsufficientStock, *model.AppError) {
	var (
		insufficientStockErr *exception.InsufficientStock
		appErr               *model.AppError
	)
	if variant.IsPreorderActive() {
		insufficientStockErr, appErr = s.CheckPreorderThresholdBulk([]*product_and_discount.ProductVariant{variant}, []int{quantity}, channelSlug)
	} else {
		insufficientStockErr, appErr = s.CheckStockQuantity(variant, countryCode, channelSlug, quantity)
	}

	return insufficientStockErr, appErr
}

// Validate if there is stock available for given variant in given country.
//
// If so - returns None. If there is less stock then required raise InsufficientStock
// exception.
func (a *ServiceWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*exception.InsufficientStock, *model.AppError) {
	if *variant.TrackInventory {
		stocks, appErr := a.GetVariantStocksForCountry(nil, countryCode, channelSlug, variant.Id)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			stocks = []*warehouse.Stock{} // in case stocks is nil
		}

		if len(stocks) == 0 {
			return &exception.InsufficientStock{
				Items: []*exception.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}

		availableQuantity, appErr := a.getAvailableQuantity(stocks)
		if appErr != nil {
			return nil, appErr
		}
		if quantity > availableQuantity {
			return &exception.InsufficientStock{
				Items: []*exception.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}
	}

	return nil, nil
}

// CheckStockAndPreorderQuantityBulk Validate if products are available for stocks/preorder.
// :raises InsufficientStock: when there is not enough items in stock for a variant
// or there is not enough available preorder items for a variant.
//
// `additionalFilterBoolup`, `existingLines` can be nil, replace default to false
func (s *ServiceWarehouse) CheckStockAndPreorderQuantityBulk(variants []*product_and_discount.ProductVariant, countryCode string, quantities []int, channelSlug string, additionalFilterBoolup model.StringInterface, existingLines []*checkout.CheckoutLineInfo, replace bool) (*exception.InsufficientStock, *model.AppError) {
	stockVariants, stockQuantities, preorderVariants, preorderQuantities := s.splitLinesForTrackableAndPreorder(variants, quantities)

	var (
		insufficientStockErr *exception.InsufficientStock
		appErr               *model.AppError
	)
	if len(stockVariants) > 0 {
		insufficientStockErr, appErr = s.CheckStockQuantityBulk(stockVariants, countryCode, stockQuantities, channelSlug, additionalFilterBoolup, existingLines, replace)
	}

	if len(preorderVariants) > 0 {
		insufficientStockErr, appErr = s.CheckPreorderThresholdBulk(preorderVariants, preorderQuantities, channelSlug)
	}

	return insufficientStockErr, appErr
}

// splitLinesForTrackableAndPreorder Return variants and quantities splitted by "is_preorder_active
func (s *ServiceWarehouse) splitLinesForTrackableAndPreorder(variants []*product_and_discount.ProductVariant, quantities []int) ([]*product_and_discount.ProductVariant, []int, []*product_and_discount.ProductVariant, []int) {
	var (
		preorderVariants   []*product_and_discount.ProductVariant
		preorderQuantities []int
		stockVariants      []*product_and_discount.ProductVariant
		stockQuantities    []int
	)

	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		variant := variants[i]
		quantity := quantities[i]

		if variant.IsPreorderActive() {
			preorderVariants = append(preorderVariants, variant)
			preorderQuantities = append(preorderQuantities, quantity)
		} else {
			stockVariants = append(stockVariants, variant)
			stockQuantities = append(stockQuantities, quantity)
		}
	}

	return stockVariants, stockQuantities, preorderVariants, preorderQuantities
}

// Validate if there is stock available for given variants in given country.
//
// :raises InsufficientStock: when there is not enough items in stock for a variant
func (a *ServiceWarehouse) CheckStockQuantityBulk(
	variants product_and_discount.ProductVariants,
	countryCode string,
	quantities []int,
	channelSlug string,
	additionalFilterLookup model.StringInterface, // can be nil, if non-nil then it must be map[string]interface{}{"warehouse_id": <an UUID string>}
	existingLines []*checkout.CheckoutLineInfo, // can be nil
	replace bool, // default false
) (*exception.InsufficientStock, *model.AppError) {

	variants = variants.FilterNils()

	// build a filter option
	allVariantStockFilterOption := &warehouse.StockFilterForCountryAndChannel{
		CountryCode:              countryCode,
		ChannelSlug:              channelSlug,
		ProductVariantIDFilter:   squirrel.Eq{store.StockTableName + ".ProductVariantID": variants.IDs()},
		AnnotateAvailabeQuantity: true,
	}

	// check if `additionalFilterLookup` is not nil:
	if additionalFilterLookup != nil && additionalFilterLookup["warehouse_id"] != nil {
		allVariantStockFilterOption.WarehouseID = additionalFilterLookup["warehouse_id"].(string)
	}

	allVariantStocks, appErr := a.FilterStocksForCountryAndChannel(nil, allVariantStockFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	// variantStocks has keys of product variant ids
	var variantStocks = map[string][]*warehouse.Stock{}
	for _, stock := range allVariantStocks {
		if stock != nil {
			variantStocks[stock.ProductVariantID] = append(variantStocks[stock.ProductVariantID], stock)
		}
	}

	var (
		// keys are product variant ids, values are order line quantities
		variantsQuantities = map[string]int{}
		insufficientStocks = []*exception.InsufficientStockData{}
	)

	for _, lineInfo := range existingLines {
		if lineInfo != nil && model.IsValidId(lineInfo.Variant.Id) {
			variantsQuantities[lineInfo.Variant.Id] = lineInfo.Line.Quantity
		}
	}

	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		quantity := quantities[i]
		variant := variants[i]

		if !replace {
			quantity += variantsQuantities[variant.Id]
		}

		stocks, exists := variantStocks[variant.Id]
		if !exists {
			stocks = []*warehouse.Stock{}
		}

		var availableQuantity int = 0
		for _, stock := range stocks {
			availableQuantity += stock.AvailableQuantity
		}

		if quantity > 0 {
			if len(stocks) == 0 {
				insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &availableQuantity,
				})
			} else if *variant.TrackInventory && quantity > availableQuantity {
				insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &availableQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return &exception.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	return nil, nil
}

type structObject struct {
	AvailablePreorderQuantity int
	PreorderQuantityThreshold *int
}

// CheckPreorderThresholdBulk Validate if there is enough preordered variants according to thresholds.
// :raises InsufficientStock: when there is not enough available items for a variant.
func (s *ServiceWarehouse) CheckPreorderThresholdBulk(variants product_and_discount.ProductVariants, quantities []int, channelSlug string) (*exception.InsufficientStock, *model.AppError) {
	allVariantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(nil, &product_and_discount.ProductVariantChannelListingFilterOption{
		VariantID:                         squirrel.Eq{s.srv.Store.ProductVariantChannelListing().TableName("VariantID"): variants.IDs()},
		SelectRelatedChannel:              true,
		AnnotatePreorderQuantityAllocated: true,
		AnnotateAvailablePreorderQuantity: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		// variantsChannelAvailability has keys are product variant ids
		variantsChannelAvailability = map[string]structObject{}
		// variantChannels has keys are product variant ids
		variantChannels = map[string][]*product_and_discount.ProductVariantChannelListing{}
		// variantsGlobalAllocations has keys are product variant ids
		variantsGlobalAllocations = map[string]int{}
	)

	for _, channelListing := range allVariantChannelListings {
		if channelListing.Channel != nil && channelListing.Channel.Slug == channelSlug {

			variantsChannelAvailability[channelListing.VariantID] = structObject{
				AvailablePreorderQuantity: channelListing.Get_availablePreorderQuantity(),
				PreorderQuantityThreshold: channelListing.PreorderQuantityThreshold,
			}

		}

		variantChannels[channelListing.VariantID] = append(variantChannels[channelListing.VariantID], channelListing)
	}

	for variantID, channelListings := range variantChannels {
		for _, channelListing := range channelListings {
			variantsGlobalAllocations[variantID] += channelListing.Get_preorderQuantityAllocated()
		}
	}

	var insufficientStocks []*exception.InsufficientStockData
	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		var (
			variant  = variants[i]
			quantity = quantities[i]
		)

		if variantsChannelAvailability[variant.Id].PreorderQuantityThreshold != nil {
			if quantity > variantsChannelAvailability[variant.Id].AvailablePreorderQuantity {
				insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: variantsChannelAvailability[variant.Id].PreorderQuantityThreshold,
				})
			}
		}

		if variant.PreOrderGlobalThreshold != nil {
			globalQuantity := *variant.PreOrderGlobalThreshold - variantsGlobalAllocations[variant.Id]
			if quantity > globalQuantity {
				insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &globalQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return exception.NewInsufficientStock(insufficientStocks), nil
	}

	return nil, nil
}
