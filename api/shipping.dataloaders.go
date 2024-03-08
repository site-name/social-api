package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

func shippingMethodByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingMethod] {
	res := make([]*dataloader.Result[*model.ShippingMethod], len(ids))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	methods, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			Conditions: squirrel.Eq{model.ShippingMethodTableName + ".Id": ids},
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
			Conditions: squirrel.Eq{
				model.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs,
				model.ShippingMethodChannelListingTableName + ".ChannelID":        channelIDs,
			},
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
	var (
		res      = make([]*dataloader.Result[[]*model.ShippingZone], len(ids))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		channels model.Channels
	)

	err := embedCtx.App.Srv().
		Store.
		GetReplica().
		Preload("ShippingZones").
		Find(&channels, "Id IN ?", ids).Error
	if err != nil {
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.ShippingZone]{Error: err}
		}
		return res
	}

	channelMap := map[string]*model.Channel{}
	for _, ch := range channels {
		channelMap[ch.Id] = ch
	}

	for idx, id := range ids {
		var zones model.ShippingZones

		channel := channelMap[id]
		if channel != nil {
			zones = channel.ShippingZones
		}
		res[idx] = &dataloader.Result[[]*model.ShippingZone]{Data: zones}
	}
	return res
}

func shippingZoneByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingZone] {
	res := make([]*dataloader.Result[*model.ShippingZone], len(ids))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Conditions: squirrel.Eq{model.ShippingZoneTableName + ".Id": ids},
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

func shippingMethodsByShippingZoneIdLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[model.ShippingMethodSlice] {
	res := make([]*dataloader.Result[model.ShippingMethodSlice], len(shippingZoneIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingMethods, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			Conditions: squirrel.Eq{model.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
		})
	if appErr != nil {
		for idx := range shippingZoneIDs {
			res[idx] = &dataloader.Result[model.ShippingMethodSlice]{Error: appErr}
		}
		return res
	}

	var shippingMethodMap = map[string]model.ShippingMethodSlice{} // keys are shipping zone ids
	for _, method := range shippingMethods {
		shippingMethodMap[method.ShippingZoneID] = append(shippingMethodMap[method.ShippingZoneID], method)
	}

	for idx, zoneID := range shippingZoneIDs {
		res[idx] = &dataloader.Result[model.ShippingMethodSlice]{Data: shippingMethodMap[zoneID]}
	}
	return res
}

func postalCodeRulesByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodPostalCodeRule] {
	res := make([]*dataloader.Result[[]*model.ShippingMethodPostalCodeRule], len(shippingMethodIDs))

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	rules, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodPostalCodeRulesByOptions(&model.ShippingMethodPostalCodeRuleFilterOptions{
			Conditions: squirrel.Eq{model.ShippingMethodPostalCodeRuleTableName + ".ShippingMethodID": shippingMethodIDs},
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
	var (
		res      = make([]*dataloader.Result[[]*model.Product], len(ids))
		embedCtx = GetContextValue[*web.Context](ctx, WebCtx)
		methods  model.ShippingMethodSlice
	)

	err := embedCtx.App.Srv().Store.GetReplica().Preload("").Find(&methods, "Id IN ?", ids).Error
	if err != nil {
		err = model_helper.NewAppError("excludedProductByShippingMethodIDLoader", "app.shipping.shipping_method_excluded_product_relations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		for idx := range ids {
			res[idx] = &dataloader.Result[[]*model.Product]{Error: err}
		}
		return res
	}

	var methodMap = map[string]*model.ShippingMethod{}
	for _, method := range methods {
		methodMap[method.Id] = method
	}

	for idx, id := range ids {
		var products model.Products

		method := methodMap[id]
		if method != nil {
			products = method.ExcludedProducts
		}
		res[idx] = &dataloader.Result[[]*model.Product]{Data: products}
	}
	return res
}

func channelsByShippingZoneIdLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[[]*model.Channel] {
	var (
		res           = make([]*dataloader.Result[[]*model.Channel], len(shippingZoneIDs))
		embedCtx      = GetContextValue[*web.Context](ctx, WebCtx)
		shippingZones model.ShippingZones
	)

	err := embedCtx.App.Srv().Store.
		GetReplica().
		Preload("Channels").
		Find(&shippingZones, "Id IN ?", shippingZoneIDs).
		Error
	if err != nil {
		for idx := range shippingZoneIDs {
			res[idx] = &dataloader.Result[[]*model.Channel]{Error: err}
		}
		return res
	}

	zonesMap := map[string]*model.ShippingZone{}
	for _, zone := range shippingZones {
		zonesMap[zone.Id] = zone
	}

	for idx, id := range shippingZoneIDs {
		var channels model.Channels
		zone := zonesMap[id]
		if zone != nil {
			channels = zone.Channels
		}

		res[idx] = &dataloader.Result[[]*model.Channel]{Data: channels}
	}
	return res
}

func shippingMethodsByShippingZoneIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[model.ShippingMethodSlice] {
	var (
		res                                        = make([]*dataloader.Result[model.ShippingMethodSlice], len(idPairs))
		shippingZoneIDs                            = make([]string, len(idPairs))
		channelIDs                                 = make([]string, len(idPairs))
		shippingMethodIDs                          []string
		shippingMethodChannelListings              []*model.ShippingMethodChannelListing
		shippingMethodMap                          = map[string]*model.ShippingMethod{}     // keys are shipping method ids
		shippingMethodsByShippingZoneAndChannelMap = map[string]model.ShippingMethodSlice{} // keys have format of shippingZoneID__channelID
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
			Conditions: squirrel.Eq{model.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
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
			Conditions: squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
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
		res[idx] = &dataloader.Result[model.ShippingMethodSlice]{Data: shippingMethodsByShippingZoneAndChannelMap[id]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[model.ShippingMethodSlice]{Error: appErr}
	}
	return res
}

func shippingMethodChannelListingByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodChannelListing] {
	res := make([]*dataloader.Result[[]*model.ShippingMethodChannelListing], len(shippingMethodIDs))
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingMethodChannelListings, appErr := embedCtx.App.Srv().ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
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
		res        = make([]*dataloader.Result[model.ShippingZones], len(warehouseIDs))
		embedCtx   = GetContextValue[*web.Context](ctx, WebCtx)
		warehouses model.Warehouses
	)

	err := embedCtx.App.Srv().Store.
		GetReplica().
		Preload("ShippingZones").
		Find(&warehouses, "Id IN ?", warehouseIDs).
		Error
	if err != nil {
		for idx := range warehouseIDs {
			res[idx] = &dataloader.Result[model.ShippingZones]{Error: err}
		}
		return res
	}

	warehouseMap := map[string]*model.WareHouse{}
	for _, warehouse := range warehouses {
		warehouseMap[warehouse.Id] = warehouse
	}

	for idx, id := range warehouseIDs {
		var zones model.ShippingZones
		warehouse := warehouseMap[id]
		if warehouse != nil {
			zones = warehouse.ShippingZones
		}
		res[idx] = &dataloader.Result[model.ShippingZones]{Data: zones}
	}
	return res
}
