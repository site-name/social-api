package warehouse

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
)

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations) not exported
func (a *ServiceWarehouse) getAvailableQuantity(stocks model.Stocks) (int, *model_helper.AppError) {
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

	allocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		Conditions: squirrel.Eq{model.AllocationTableName + ".StockID": stockIDs},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return 0, appErr
		}
		allocations = []*model.Allocation{}
	}

	var quantityAllocated int
	for _, allocation := range allocations {
		quantityAllocated += allocation.QuantityAllocated
	}

	return max(totalQuantity-quantityAllocated, 0), nil
}

// CheckStockAndPreorderQuantity Validate if there is stock/preorder available for given variant.
// :raises InsufficientStock: when there is not enough items in stock for a variant
// or there is not enough available preorder items for a variant.
func (s *ServiceWarehouse) CheckStockAndPreorderQuantity(variant *model.ProductVariant, countryCode model.CountryCode, channelSlug string, quantity int) (*model_helper.InsufficientStock, *model_helper.AppError) {
	var (
		insufficientStockErr *model_helper.InsufficientStock
		appErr               *model_helper.AppError
	)
	if variant.IsPreorderActive() {
		insufficientStockErr, appErr = s.CheckPreorderThresholdBulk([]*model.ProductVariant{variant}, []int{quantity}, channelSlug)
	} else {
		insufficientStockErr, appErr = s.CheckStockQuantity(variant, countryCode, channelSlug, quantity)
	}

	return insufficientStockErr, appErr
}

// Validate if there is stock available for given variant in given country.
//
// If so - returns None. If there is less stock then required raise InsufficientStock
// exception.
func (a *ServiceWarehouse) CheckStockQuantity(variant *model.ProductVariant, countryCode model.CountryCode, channelSlug string, quantity int) (*model_helper.InsufficientStock, *model_helper.AppError) {
	if *variant.TrackInventory {
		stocks, appErr := a.GetVariantStocksForCountry(countryCode, channelSlug, variant.Id)
		if appErr != nil {
			return nil, appErr
		}

		if len(stocks) == 0 {
			return &model_helper.InsufficientStock{
				Items: []*model.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}

		availableQuantity, appErr := a.getAvailableQuantity(stocks)
		if appErr != nil {
			return nil, appErr
		}
		if quantity > availableQuantity {
			return &model_helper.InsufficientStock{
				Items: []*model.InsufficientStockData{
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
func (s *ServiceWarehouse) CheckStockAndPreorderQuantityBulk(variants []*model.ProductVariant, countryCode model.CountryCode, quantities []int, channelSlug string, additionalFilterBoolup model_types.JSONString, existingLines model_helper.CheckoutLineInfos, replace bool) (*model_helper.InsufficientStock, *model_helper.AppError) {
	stockVariants, stockQuantities, preorderVariants, preorderQuantities := s.splitLinesForTrackableAndPreorder(variants, quantities)

	if len(stockVariants) > 0 {
		insufficientStockErr, appErr := s.CheckStockQuantityBulk(stockVariants, countryCode, stockQuantities, channelSlug, additionalFilterBoolup, existingLines, replace)
		if insufficientStockErr != nil || appErr != nil {
			return insufficientStockErr, appErr
		}
	}

	if len(preorderVariants) > 0 {
		insufficientStockErr, appErr := s.CheckPreorderThresholdBulk(preorderVariants, preorderQuantities, channelSlug)
		if insufficientStockErr != nil || appErr != nil {
			return insufficientStockErr, appErr
		}
	}

	return nil, nil
}

// splitLinesForTrackableAndPreorder Return variants and quantities splitted by "is_preorder_active
func (s *ServiceWarehouse) splitLinesForTrackableAndPreorder(variants []*model.ProductVariant, quantities []int) ([]*model.ProductVariant, []int, []*model.ProductVariant, []int) {
	var (
		preorderVariants   []*model.ProductVariant
		preorderQuantities []int
		stockVariants      []*model.ProductVariant
		stockQuantities    []int
	)

	for i := 0; i < min(len(variants), len(quantities)); i++ {
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
	variants model.ProductVariantSlice,
	countryCode model.CountryCode,
	quantities []int,
	channelSlug string,
	additionalFilterLookup model_types.JSONString, // can be nil, if non-nil then it must be map[string]any{"warehouse_id": <an UUID string>}
	existingLines model_helper.CheckoutLineInfos, // can be nil
	replace bool, // default false
) (*model_helper.InsufficientStock, *model_helper.AppError) {

	variants = variants.FilterNils()

	// build a filter option
	allVariantStockFilterOption := &model.StockFilterOptionsForCountryAndChannel{
		CountryCode:               countryCode,
		ChannelSlug:               channelSlug,
		ProductVariantIDFilter:    squirrel.Eq{model.StockTableName + ".ProductVariantID": variants.IDs()},
		AnnotateAvailableQuantity: true,
	}

	// check if `additionalFilterLookup` is not nil:
	if additionalFilterLookup != nil && additionalFilterLookup["warehouse_id"] != nil {
		allVariantStockFilterOption.WarehouseID = additionalFilterLookup["warehouse_id"].(string)
	}

	allVariantStocks, appErr := a.FilterStocksForCountryAndChannel(allVariantStockFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	// variantStocks has keys of product variant ids
	var variantStocks = map[string][]*model.Stock{}
	for _, stock := range allVariantStocks {
		if stock != nil {
			variantStocks[stock.ProductVariantID] = append(variantStocks[stock.ProductVariantID], stock)
		}
	}

	var (
		// keys are product variant ids, values are order line quantities
		variantsQuantities = map[string]int{}
		insufficientStocks = []*model.InsufficientStockData{}
	)

	for _, lineInfo := range existingLines {
		if lineInfo != nil && model_helper.IsValidId(lineInfo.Variant.Id) {
			variantsQuantities[lineInfo.Variant.Id] = lineInfo.Line.Quantity
		}
	}

	for i := 0; i < min(len(variants), len(quantities)); i++ {
		quantity := quantities[i]
		variant := variants[i]

		if !replace {
			quantity += variantsQuantities[variant.Id]
		}

		stocks, exists := variantStocks[variant.Id]
		if !exists {
			stocks = []*model.Stock{}
		}

		var availableQuantity int = 0
		for _, stock := range stocks {
			availableQuantity += stock.AvailableQuantity
		}

		if quantity > 0 {
			if len(stocks) == 0 {
				insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &availableQuantity,
				})
			} else if *variant.TrackInventory && quantity > availableQuantity {
				insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &availableQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return &model_helper.InsufficientStock{
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
func (s *ServiceWarehouse) CheckPreorderThresholdBulk(variants model.ProductVariantSlice, quantities []int, channelSlug string) (*model_helper.InsufficientStock, *model_helper.AppError) {
	allVariantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		Conditions:                        squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": variants.IDs()},
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
		variantChannels = map[string][]*model.ProductVariantChannelListing{}
		// variantsGlobalAllocations has keys are product variant ids
		variantsGlobalAllocations = map[string]int{}
	)

	for _, channelListing := range allVariantChannelListings {
		if channelListing.GetChannel() != nil && channelListing.GetChannel().Slug == channelSlug {

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

	var insufficientStocks []*model.InsufficientStockData
	for i := 0; i < min(len(variants), len(quantities)); i++ {
		var (
			variant  = variants[i]
			quantity = quantities[i]
		)

		if variantsChannelAvailability[variant.Id].PreorderQuantityThreshold != nil {
			if quantity > variantsChannelAvailability[variant.Id].AvailablePreorderQuantity {
				insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: variantsChannelAvailability[variant.Id].PreorderQuantityThreshold,
				})
			}
		}

		if variant.PreOrderGlobalThreshold != nil {
			globalQuantity := *variant.PreOrderGlobalThreshold - variantsGlobalAllocations[variant.Id]
			if quantity > globalQuantity {
				insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &globalQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return model.NewInsufficientStock(insufficientStocks), nil
	}

	return nil, nil
}
