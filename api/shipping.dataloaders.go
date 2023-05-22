package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func shippingMethodByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingMethod] {
	res := make([]*dataloader.Result[*model.ShippingMethod], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	methods, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			Id: squirrel.Eq{store.ShippingMethodTableName + ".Id": ids},
		})
	if appErr != nil {
		for i := range ids {
			res[i] = &dataloader.Result[*model.ShippingMethod]{Error: appErr}
		}
		return res
	}

	var methodMap = map[string]*model.ShippingMethod{} // keys are shipping method ids
	for _, method := range methods {
		methodMap[method.Id] = method
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ShippingMethod]{Data: methodMap[id]}
	}
	return res
}

// idPairs are slices of shippingMethodID__channelID string formats
func shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ShippingMethodChannelListing] {
	var (
		res               = make([]*dataloader.Result[*model.ShippingMethodChannelListing], len(idPairs))
		shippingMethodIDs = make([]string, len(idPairs))
		channelIDs        = make([]string, len(idPairs))
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			shippingMethodIDs[idx] = pair[:index]
			channelIDs[idx] = pair[index+2:]
		}
	}

	shippingMethodChannelListings, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
			ChannelID:        squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		for idx := range idPairs {
			res[idx] = &dataloader.Result[*model.ShippingMethodChannelListing]{Error: appErr}
		}
		return res
	}

	shippingMethodChannelListingMap := lo.SliceToMap(shippingMethodChannelListings, func(r *model.ShippingMethodChannelListing) (string, *model.ShippingMethodChannelListing) {
		return r.ShippingMethodID + "__" + r.ChannelID, r
	})
	for idx, pair := range idPairs {
		res[idx] = &dataloader.Result[*model.ShippingMethodChannelListing]{Data: shippingMethodChannelListingMap[pair]}
	}
	return res
}

func shippingZonesByChannelIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.ShippingZone] {
	res := make([]*dataloader.Result[[]*model.ShippingZone], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	relations, err := embedCtx.App.Srv().Store.
		ShippingZoneChannel().
		FilterByOptions(&model.ShippingZoneChannelFilterOptions{
			ChannelID: squirrel.Eq{store.ShippingZoneChannelTableName + ".ChannelID": ids},
		})
	if err != nil {
		err = model.NewAppError("ShippingZoneChanenlByOptions", "app.shipping.shippingzone-channel-relations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.ShippingZone]{Error: err}
		}
		return res
	}

	var (
		shippingZoneMap        = map[string]*model.ShippingZone{}
		shippingZoneIDs        []string
		channelShippingZoneMap = map[string][]string{} // keys are channel ids, values are slices of shipping zone ids
	)
	for _, rel := range relations {
		shippingZoneIDs = append(shippingZoneIDs, rel.ShippingZoneID)
		channelShippingZoneMap[rel.ChannelID] = append(channelShippingZoneMap[rel.ChannelID], rel.ShippingZoneID)
	}

	shippingZones, errs := ShippingZoneByIdLoader.LoadMany(ctx, shippingZoneIDs)()
	if len(errs) > 0 && errs[0] != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.ShippingZone]{Error: errs[0]}
		}
		return res
	}

	for _, zone := range shippingZones {
		shippingZoneMap[zone.Id] = zone
	}

	for idx, channelID := range ids {
		data := []*model.ShippingZone{}
		for _, zoneID := range channelShippingZoneMap[channelID] {
			data = append(data, shippingZoneMap[zoneID])
		}

		res[idx] = &dataloader.Result[[]*model.ShippingZone]{Data: data}
	}

	return res
}

func shippingZoneByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingZone] {
	res := make([]*dataloader.Result[*model.ShippingZone], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id: squirrel.Eq{store.ShippingZoneTableName + ".Id": ids},
	})
	if appErr != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[*model.ShippingZone]{Error: appErr}
		}
		return res
	}

	shippingZoneMap := lo.SliceToMap(shippingZones, func(s *model.ShippingZone) (string, *model.ShippingZone) { return s.Id, s })
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ShippingZone]{Data: shippingZoneMap[id]}
	}
	return res
}

func shippingMethodsByShippingZoneIdLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[[]*model.ShippingMethod] {
	res := make([]*dataloader.Result[[]*model.ShippingMethod], len(shippingZoneIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingMethods, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			ShippingZoneID: squirrel.Eq{store.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
		})
	if appErr != nil {
		for idx := range shippingZoneIDs {
			res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Error: appErr}
		}
		return res
	}

	var shippingMethodMap = map[string][]*model.ShippingMethod{} // keys are shipping zone ids
	for _, method := range shippingMethods {
		shippingMethodMap[method.ShippingZoneID] = append(shippingMethodMap[method.ShippingZoneID], method)
	}

	for idx, zoneID := range shippingZoneIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Data: shippingMethodMap[zoneID]}
	}
	return res
}

func postalCodeRulesByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodPostalCodeRule] {
	res := make([]*dataloader.Result[[]*model.ShippingMethodPostalCodeRule], len(shippingMethodIDs))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	rules, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodPostalCodeRulesByOptions(&model.ShippingMethodPostalCodeRuleFilterOptions{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodPostalCodeRuleTableName + ".ShippingMethodID": shippingMethodIDs},
		})
	if appErr != nil {
		for idx := range shippingMethodIDs {
			res[idx] = &dataloader.Result[[]*model.ShippingMethodPostalCodeRule]{Error: appErr}
		}
		return res
	}

	var ruleMap = map[string][]*model.ShippingMethodPostalCodeRule{}
	for _, rule := range rules {
		ruleMap[rule.ShippingMethodID] = append(ruleMap[rule.ShippingMethodID], rule)
	}

	for idx, id := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodPostalCodeRule]{Data: ruleMap[id]}
	}
	return res
}

func excludedProductByShippingMethodIDLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.Product] {
	res := make([]*dataloader.Result[[]*model.Product], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingMethodExcludedProducts, err := embedCtx.
		App.
		Srv().
		Store.
		ShippingMethodExcludedProduct().
		FilterByOptions(&model.ShippingMethodExcludedProductFilterOptions{
			SelectRelatedProduct: true,
			ShippingMethodID:     squirrel.Eq{store.ShippingMethodExcludedProductTableName + ".ShippingMethodID": ids},
		})
	if err != nil {
		err = model.NewAppError("excludedProductByShippingMethodIDLoader", "app.shipping.shipping_method_excluded_product_relations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
		}
		return res
	}

	var shippingMethodExcludedProductMap = map[string]model.Products{} // keys are shipping method ids
	for _, rel := range shippingMethodExcludedProducts {
		shippingMethodExcludedProductMap[rel.ShippingMethodID] = append(shippingMethodExcludedProductMap[rel.ShippingMethodID], rel.GetProduct())
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[[]*model.Product]{Data: shippingMethodExcludedProductMap[id]}
	}
	return res
}

func channelsByShippingZoneIdLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[[]*model.Channel] {
	var (
		res                           = make([]*dataloader.Result[[]*model.Channel], len(shippingZoneIDs))
		channelAndShippingZoneIDPairs = map[string][]string{} // keys are shipping zone ids, values are slices of channel ids
		channelIDs                    []string
		channels                      []*model.Channel
		channelMap                    = map[string]*model.Channel{} // keys are channel ids
		errs                          []error
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingZoneChannels, err := embedCtx.
		App.
		Srv().
		Store.
		ShippingZoneChannel().
		FilterByOptions(&model.ShippingZoneChannelFilterOptions{
			ShippingZoneID: squirrel.Eq{store.ShippingZoneChannelTableName + ".ShippingZoneID": shippingZoneIDs},
		})
	if err != nil {
		err = model.NewAppError("channelsByShippingZoneIdLoader", "app.shipping.shipping_zone_channels_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}
	channelIDs = lo.Map(shippingZoneChannels, func(r *model.ShippingZoneChannel, _ int) string { return r.ChannelID })

	for _, rel := range shippingZoneChannels {
		channelAndShippingZoneIDPairs[rel.ShippingZoneID] = append(channelAndShippingZoneIDPairs[rel.ShippingZoneID], rel.ChannelID)
	}

	channels, errs = ChannelByIdLoader.LoadMany(ctx, channelIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
	}
	for _, channel := range channels {
		channelMap[channel.Id] = channel
	}

	for idx, id := range shippingZoneIDs {
		channels := lo.Map(channelAndShippingZoneIDPairs[id], func(channelID string, _ int) *model.Channel { return channelMap[channelID] })
		channels = lo.Filter(channels, func(c *model.Channel, _ int) bool { return c != nil })

		res[idx] = &dataloader.Result[[]*model.Channel]{Data: channels}
	}
	return res

errorLabel:
	for idx := range shippingZoneIDs {
		res[idx] = &dataloader.Result[[]*model.Channel]{Error: err}
	}
	return res
}

func shippingMethodsByShippingZoneIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ShippingMethod] {
	var (
		res                                        = make([]*dataloader.Result[[]*model.ShippingMethod], len(idPairs))
		shippingZoneIDs                            = make([]string, len(idPairs))
		channelIDs                                 = make([]string, len(idPairs))
		shippingMethodIDs                          []string
		shippingMethodChannelListings              []*model.ShippingMethodChannelListing
		shippingMethodMap                          = map[string]*model.ShippingMethod{}   // keys are shipping method ids
		shippingMethodsByShippingZoneAndChannelMap = map[string][]*model.ShippingMethod{} // keys have format of shippingZoneID__channelID
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			shippingZoneIDs[idx] = pair[:index]
			channelIDs[idx] = pair[index+2:]
		}
	}

	shippingMethods, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			ShippingZoneID: squirrel.Eq{store.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	shippingMethodIDs = lo.Map(shippingMethods, func(s *model.ShippingMethod, _ int) string {
		shippingMethodMap[s.Id] = s // trick to reduce number of for loop
		return s.Id
	})
	shippingMethodChannelListings, appErr = embedCtx.App.Srv().ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	for _, listing := range shippingMethodChannelListings {
		shippingMethod, ok := shippingMethodMap[listing.ShippingMethodID]
		if ok {
			key := shippingMethod.ShippingZoneID + "__" + listing.ChannelID
			shippingMethodsByShippingZoneAndChannelMap[key] = append(shippingMethodsByShippingZoneAndChannelMap[key], shippingMethod)
		}
	}
	for idx, id := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Data: shippingMethodsByShippingZoneAndChannelMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Error: appErr}
	}
	return res
}

func shippingMethodChannelListingByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodChannelListing] {
	res := make([]*dataloader.Result[[]*model.ShippingMethodChannelListing], len(shippingMethodIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingMethodChannelListings, appErr := embedCtx.App.Srv().ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
		})
	if appErr != nil {
		for idx := range shippingMethodIDs {
			res[idx] = &dataloader.Result[[]*model.ShippingMethodChannelListing]{Error: appErr}
		}
		return res
	}

	var shippingMethodChannelListingsMap = map[string][]*model.ShippingMethodChannelListing{}
	for _, listing := range shippingMethodChannelListings {
		shippingMethodChannelListingsMap[listing.ShippingMethodID] = append(shippingMethodChannelListingsMap[listing.ShippingMethodID], listing)
	}
	for idx, id := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodChannelListing]{Data: shippingMethodChannelListingsMap[id]}
	}
	return res
}

func shippingZonesByWarehouseIDLoader(ctx context.Context, warehouseIDs []string) []*dataloader.Result[model.ShippingZones] {
	var (
		res                       = make([]*dataloader.Result[model.ShippingZones], len(warehouseIDs))
		shippingZoneIDs           []string
		shippingZones             model.ShippingZones
		appErr                    *model.AppError
		warehouseShippingZonesMap = map[string]model.ShippingZones{} // keys are warehouse ids
		shippingZoneMap           = map[string]*model.ShippingZone{} // keys are shipping zone ids
	)
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	warehouseShippingZones, err := embedCtx.
		App.
		Srv().
		Store.
		WarehouseShippingZone().
		FilterByOptions(&model.WarehouseShippingZoneFilterOption{
			WarehouseID: squirrel.Eq{store.WarehouseShippingZoneTableName + ".WarehouseID": warehouseIDs},
		})
	if err != nil {
		appErr = model.NewAppError("shippingZonesByWarehouseIDLoader", "app.shipping.warehouse_shipping_zones_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	shippingZoneIDs = lo.Map(warehouseShippingZones, func(r *model.WarehouseShippingZone, _ int) string { return r.ShippingZoneID })
	shippingZones, appErr = embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id: squirrel.Eq{store.ShippingZoneTableName + ".Id": shippingZoneIDs},
	})
	if appErr != nil {
		goto errorLabel
	}
	shippingZoneMap = lo.SliceToMap(shippingZones, func(s *model.ShippingZone) (string, *model.ShippingZone) { return s.Id, s })

	for _, rel := range warehouseShippingZones {
		zone, exist := shippingZoneMap[rel.ShippingZoneID]
		if exist {
			warehouseShippingZonesMap[rel.WarehouseID] = append(warehouseShippingZonesMap[rel.WarehouseID], zone)
		}
	}

	for idx, id := range warehouseIDs {
		res[idx] = &dataloader.Result[model.ShippingZones]{Data: warehouseShippingZonesMap[id]}
	}
	return res

errorLabel:
	for idx := range warehouseIDs {
		res[idx] = &dataloader.Result[model.ShippingZones]{Error: appErr}
	}
	return res
}
