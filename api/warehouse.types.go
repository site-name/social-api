package api

import (
	"context"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

type Warehouse struct {
	ID                    string                               `json:"id"`
	Name                  string                               `json:"name"`
	Slug                  string                               `json:"slug"`
	Email                 string                               `json:"email"`
	IsPrivate             bool                                 `json:"isPrivate"`
	PrivateMetadata       []*MetadataItem                      `json:"privateMetadata"`
	Metadata              []*MetadataItem                      `json:"metadata"`
	ClickAndCollectOption model.WarehouseClickAndCollectOption `json:"clickAndCollectOption"`

	addressID *string
	// ShippingZones         *ShippingZoneCountableConnection   `json:"shippingZones"`
	// Address               *Address                           `json:"address"`
}

func SystemWarehouseToGraphqlWarehouse(wh *model.WareHouse) *Warehouse {
	if wh == nil {
		return nil
	}

	return &Warehouse{
		ID:                    wh.Id,
		Name:                  wh.Name,
		Slug:                  wh.Slug,
		Email:                 wh.Email,
		IsPrivate:             *wh.IsPrivate,
		Metadata:              MetadataToSlice(wh.Metadata),
		PrivateMetadata:       MetadataToSlice(wh.Metadata),
		ClickAndCollectOption: wh.ClickAndCollectOption,

		addressID: wh.AddressID,
	}
}

func (w *Warehouse) ShippingZones(ctx context.Context, args GraphqlParams) (*ShippingZoneCountableConnection, error) {
	shippingZones, err := ShippingZonesByWarehouseIDLoader.Load(ctx, w.ID)()
	if err != nil {
		return nil, err
	}

	keyFunc := func(sz *model.ShippingZone) int64 { return sz.CreateAt }
	res, appErr := newGraphqlPaginator(shippingZones, keyFunc, SystemShippingZoneToGraphqlShippingZone, args).parse("Warehouse.ShippingZones")
	if appErr != nil {
		return nil, appErr
	}

	return (*ShippingZoneCountableConnection)(unsafe.Pointer(res)), nil
}

func (w *Warehouse) Address(ctx context.Context) (*Address, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	if w.addressID == nil {
		return nil, nil
	}

	address, appErr := embedCtx.App.Srv().AccountService().AddressById(*w.addressID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemAddressToGraphqlAddress(address), nil
}

func warehouseByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.WareHouse] {
	var (
		res          = make([]*dataloader.Result[*model.WareHouse], len(ids))
		warehouseMap = map[string]*model.WareHouse{} // keys are warehouse ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	warehouses, appErr := embedCtx.App.Srv().
		WarehouseService().
		WarehousesByOption(&model.WarehouseFilterOption{
			Conditions: squirrel.Eq{model.WarehouseTableName + ".Id": ids},
		})
	if appErr != nil {
		for i := range ids {
			res[i] = &dataloader.Result[*model.WareHouse]{Error: appErr}
		}
		return res
	}

	for _, wh := range warehouses {
		warehouseMap[wh.Id] = wh
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.WareHouse]{Data: warehouseMap[id]}
	}
	return res
}

func warehousesByShippingZoneIDLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[model.Warehouses] {
	var res = make([]*dataloader.Result[model.Warehouses], len(shippingZoneIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	var shippingZones model.ShippingZones
	err := embedCtx.App.Srv().Store.GetReplica().Preload("Warehouses").Find(&shippingZones, "Id IN ?", shippingZoneIDs).Error
	if err != nil {
		appErr := model.NewAppError("warehousesByShippingZoneIDLoader", "api.warehouse.shipping_zones_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
		for i := range shippingZoneIDs {
			res[i] = &dataloader.Result[model.Warehouses]{Error: appErr}
		}
		return res
	}

	shippingZoneMap := map[string]*model.ShippingZone{}
	for _, zone := range shippingZones {
		shippingZoneMap[zone.Id] = zone
	}

	for idx, id := range shippingZoneIDs {
		zone := shippingZoneMap[id]
		if zone == nil {
			zone = new(model.ShippingZone)
		}
		res[idx] = &dataloader.Result[model.Warehouses]{Data: zone.Warehouses}
	}
	return res
}

// ---------------------- stock --------------------

type Stock struct {
	ID    string `json:"id"`
	stock *model.Stock
}

func SystemStockToGraphqlStock(s *model.Stock) *Stock {
	if s == nil {
		return nil
	}

	return &Stock{
		ID:    s.Id,
		stock: s,
	}
}

func (s *Stock) Warehouse(ctx context.Context) (*Warehouse, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	warehouse, appErr := embedCtx.App.Srv().WarehouseService().WarehouseByStockID(s.ID)
	if appErr != nil {
		return nil, appErr
	}

	return SystemWarehouseToGraphqlWarehouse(warehouse), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (s *Stock) Quantity(ctx context.Context) (int32, error) {
	return int32(s.stock.Quantity), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (s *Stock) QuantityAllocated(ctx context.Context) (int32, error) {
	allocations, err := AllocationsByStockIDLoader.Load(ctx, s.ID)()
	if err != nil {
		return 0, err
	}

	var sum int
	for _, allocation := range allocations {
		sum += allocation.QuantityAllocated
	}
	if sum < 0 {
		sum = 0
	}
	return int32(sum), nil

}

func (s *Stock) ProductVariant(ctx context.Context) (*ProductVariant, error) {
	variant, err := ProductVariantByIdLoader.Load(ctx, s.stock.ProductVariantID)()
	if err != nil {
		return nil, err
	}
	return SystemProductVariantToGraphqlProductVariant(variant), nil
}

func allocationsByStockIDLoader(ctx context.Context, stockIDs []string) []*dataloader.Result[[]*model.Allocation] {
	var (
		res            = make([]*dataloader.Result[[]*model.Allocation], len(stockIDs))
		allocationsMap = map[string]model.Allocations{} // keys are stock ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	allocations, appErr := embedCtx.App.Srv().WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
		Conditions: squirrel.Eq{model.AllocationTableName + ".StockID": stockIDs},
	})
	if appErr != nil {
		for idx := range stockIDs {
			res[idx] = &dataloader.Result[[]*model.Allocation]{Error: appErr}
		}
		return res
	}

	for _, allocation := range allocations {
		allocationsMap[allocation.StockID] = append(allocationsMap[allocation.StockID], allocation)
	}

	for idx, id := range stockIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Data: allocationsMap[id]}
	}
	return res
}

func stocksByIDLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Stock] {
	var (
		res      = make([]*dataloader.Result[*model.Stock], len(ids))
		stockMap = map[string]*model.Stock{} // keys are stock ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	stocks, appErr := embedCtx.App.Srv().
		WarehouseService().
		StocksByOption(&model.StockFilterOption{
			Conditions:             squirrel.Eq{model.StockTableName + ".Id": ids},
			SelectRelatedWarehouse: true,
		})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.Stock]{Error: appErr}
		}
		return res
	}

	for _, st := range stocks {
		stockMap[st.Id] = st
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Stock]{Data: stockMap[id]}
	}
	return res
}

// ----------------- allocation ----------------

type Allocation struct {
	ID string `json:"id"`

	// Quantity  int32      `json:"quantity"`
	// Warehouse *Warehouse `json:"warehouse"`
	a *model.Allocation
}

func systemAllocationToGraphqlAllocation(a *model.Allocation) *Allocation {
	if a == nil {
		return nil
	}

	return &Allocation{
		ID: a.Id,
		a:  a,
	}
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (a *Allocation) Quantity(ctx context.Context) (int32, error) {
	return int32(a.a.QuantityAllocated), nil
}

// NOTE: Refer to ./schemas/warehouse.graphqls for details on directives used.
func (a *Allocation) Warehouse(ctx context.Context) (*Warehouse, error) {
	stock, err := StocksByIDLoader.Load(ctx, a.a.StockID)()
	if err != nil {
		return nil, err
	}

	warehouse, err := WarehouseByIdLoader.Load(ctx, stock.WarehouseID)()
	if err != nil {
		return nil, err
	}

	return SystemWarehouseToGraphqlWarehouse(warehouse), nil
}

func allocationsByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[[]*model.Allocation] {
	var (
		res           = make([]*dataloader.Result[[]*model.Allocation], len(orderLineIDs))
		allocationMap = map[string]model.Allocations{}
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	allocations, appErr := embedCtx.App.Srv().WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
		Conditions: squirrel.Eq{model.AllocationTableName + ".OrderLineID": orderLineIDs},
	})
	if appErr != nil {
		for idx := range orderLineIDs {
			res[idx] = &dataloader.Result[[]*model.Allocation]{Error: appErr}
		}
		return res
	}

	for _, all := range allocations {
		allocationMap[all.OrderLineID] = append(allocationMap[all.OrderLineID], all)
	}
	for idx, id := range orderLineIDs {
		res[idx] = &dataloader.Result[[]*model.Allocation]{Data: allocationMap[id]}
	}
	return res
}

func availableQuantityByProductVariantIdCountryCodeAndChannelIdLoader(ctx context.Context, idTripple []string) []*dataloader.Result[int] {
	var (
		res                          = make([]*dataloader.Result[int], len(idTripple))
		variantsByCountryAndChannels = map[[2]string][]string{} // keys have format of [2]string{countryCode, channelID}, values are variant ids
		quantityByVariantAndCountry  = map[string]int{}         // keys have format of: variantID__countryCode__channelID
		embedCtx                     = GetContextValue[*web.Context](ctx, WebCtx)
	)

	// the result map has keys are variant ids
	var batchLoadQuantitiesByCountry = func(countryCode, channelID string, variantIDs []string) (map[string]int, *model.AppError) {
		stockFilterOptions := &model.StockFilterOption{
			AnnotateAvailabeQuantity: true,
		}
		conditions := squirrel.And{
			squirrel.Eq{model.StockTableName + ".ProductVariantID": variantIDs},
		}

		warehouseShippingZones, err := embedCtx.App.Srv().Store.
			Warehouse().
			WarehouseShipingZonesByCountryCodeAndChannelID(countryCode, channelID)
		if err != nil {
			return nil, model.NewAppError("availableQuantityByProductVariantIdCountryCodeAndChannelIdLoader", "app.warehouse.warehouse_shipping_zones_by_country_code_and_channel_id.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		var warehouseShippingZonesMap = map[string][]string{} // keys are warehouse ids, values are shipping zone ids
		for _, warehouseShippingZone := range warehouseShippingZones {
			warehouseShippingZonesMap[warehouseShippingZone.WarehouseID] = append(warehouseShippingZonesMap[warehouseShippingZone.WarehouseID], warehouseShippingZone.ShippingZoneID)
		}

		if countryCode != "" || channelID != "" {
			conditions = append(conditions, squirrel.Eq{model.StockTableName + ".WarehouseID": lo.Keys(warehouseShippingZonesMap)})
		}

		stockFilterOptions.Conditions = conditions
		stocks, appErr := embedCtx.App.Srv().WarehouseService().StocksByOption(stockFilterOptions)
		if appErr != nil {
			return nil, appErr
		}

		// A single country code (or a missing country code) can return results from
		// multiple shipping zones. We want to combine all quantities within a single
		// zone and then find out which zone contains the highest total.

		// keys are product variant ids, values are maps with keys are shipping zone ids
		var quantityByShippingZoneByProductVariant = map[string]map[string]int{}

		for _, stock := range stocks {
			quantity := util.GetMinMax(0, stock.AvailableQuantity).Max
			for _, shippingZoneID := range warehouseShippingZonesMap[stock.WarehouseID] {
				quantityByShippingZoneByProductVariant[stock.ProductVariantID][shippingZoneID] += quantity
			}
		}

		var quantityMap = map[string]int{} // keys are variant ids

		for variantID, quantityByShippingZone := range quantityByShippingZoneByProductVariant {
			quantityValues := lo.Values(quantityByShippingZone)

			if countryCode != "" {
				// When country code is known, return the sum of quantities from all
				// shipping zones supporting given country.
				quantityMap[variantID] = util.AnyArray[int](quantityValues).Sum()
			} else {
				// When country code is unknown, return the highest known quantity.
				quantityMap[variantID] = util.GetMinMax(quantityValues...).Max
			}
		}

		// Return the quantities after capping them at the maximum quantity allowed in
		// checkout. This prevent users from tracking the store's precise stock levels.
		for key, value := range quantityMap {
			quantityMap[key] = util.GetMinMax(value, *embedCtx.App.Config().ShopSettings.MaxCheckoutLineQuantity).Min
		}
		return quantityMap, nil
	}

	for _, tripple := range idTripple {
		split := strings.Split(tripple, "__")
		if len(split) == 3 {
			key := [2]string{split[1], split[2]}
			variantsByCountryAndChannels[key] = append(variantsByCountryAndChannels[key], split[0])
		}
	}

	for key, variantIDs := range variantsByCountryAndChannels {
		countryCode, channelID := key[0], key[1]

		quantityMap, appErr := batchLoadQuantitiesByCountry(countryCode, channelID, variantIDs)
		if appErr != nil {
			for idx := range idTripple {
				res[idx] = &dataloader.Result[int]{Error: appErr}
			}
			return res
		}

		for variantID, quantity := range quantityMap {
			key := variantID + "__" + countryCode + "__" + channelID
			quantityByVariantAndCountry[key] = util.GetMinMax(0, quantity).Max
		}
	}

	for idx, tripple := range idTripple {
		res[idx] = &dataloader.Result[int]{Data: quantityByVariantAndCountry[tripple]}
	}
	return res
}

func stocksWithAvailableQuantityByProductVariantIdCountryCodeAndChannelLoader(ctx context.Context, idTripple []string) []*dataloader.Result[model.Stocks] {
	var (
		variantsByCountryAndChannel = map[[2]string][]string{} // keys have form of countryCode__channelID, values are variant ids
		res                         = make([]*dataloader.Result[model.Stocks], len(idTripple))
		stocksByVariantAndCountry   = map[string]model.Stocks{} // keys have format of variantID__countryCode__channelID
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	batchLoadStocksByCountry := func(countryCode, channelID string, variantIDs []string) (map[string]model.Stocks, *model.AppError) {
		countryCode = strings.ToUpper(countryCode)

		stockFilterOptions := &model.StockFilterOption{
			Conditions:               squirrel.Eq{model.StockTableName + ".ProductVariantID": variantIDs},
			AnnotateAvailabeQuantity: true,
		}
		if countryCode != "" {
			stockFilterOptions.Warehouse_ShippingZone_countries = squirrel.Like{model.ShippingZoneTableName + ".Countries::text": "%" + countryCode + "%"}
		}
		if channelID != "" {
			stockFilterOptions.Warehouse_ShippingZone_ChannelID = squirrel.Eq{model.ShippingZoneChannelTableName + ".ChannelID": channelID}
		}

		stocks, appErr := embedCtx.App.Srv().WarehouseService().StocksByOption(stockFilterOptions)
		if appErr != nil {
			return nil, appErr
		}

		var stocksByVariantIdMap = map[string]model.Stocks{} // keys are variant ids
		for _, stock := range stocks {
			stocksByVariantIdMap[stock.ProductVariantID] = append(stocksByVariantIdMap[stock.ProductVariantID], stock)
		}

		return stocksByVariantIdMap, nil
	}
	// end function

	for _, tripple := range idTripple {
		split := strings.Split(tripple, "__")
		if len(split) == 3 {
			key := [2]string{split[0], split[1]}
			variantsByCountryAndChannel[key] = append(variantsByCountryAndChannel[key], split[2])
		}
	}

	for key, variantIDs := range variantsByCountryAndChannel {
		countryCode, channelID := key[0], key[1]

		stocksByVariantIdMap, appErr := batchLoadStocksByCountry(countryCode, channelID, variantIDs)
		if appErr != nil {
			for idx := range idTripple {
				res[idx] = &dataloader.Result[model.Stocks]{Error: appErr}
			}
			return res
		}

		for _, variantID := range variantIDs {
			stocks, ok := stocksByVariantIdMap[variantID]
			if ok {
				key := variantID + "__" + countryCode + "__" + channelID
				stocksByVariantAndCountry[key] = append(stocksByVariantAndCountry[key], stocks...)
			}
		}
	}

	for idx, tripple := range idTripple {
		res[idx] = &dataloader.Result[model.Stocks]{Data: stocksByVariantAndCountry[tripple]}
	}
	return res
}
