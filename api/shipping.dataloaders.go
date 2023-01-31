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
	var (
		res       = make([]*dataloader.Result[*model.ShippingMethod], len(ids))
		appErr    *model.AppError
		methods   []*model.ShippingMethod
		methodMap = map[string]*model.ShippingMethod{} // keys are shipping method ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	methods, appErr = embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			Id: squirrel.Eq{store.ShippingMethodTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, method := range methods {
		methodMap[method.Id] = method
	}

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ShippingMethod]{Data: methodMap[id]}
	}
	return res

errorLabel:
	for i := range ids {
		res[i] = &dataloader.Result[*model.ShippingMethod]{Error: err}
	}
	return res
}

// idPairs are slices of shippingMethodID__channelID string formats
func shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ShippingMethodChannelListing] {
	var (
		res                             = make([]*dataloader.Result[*model.ShippingMethodChannelListing], len(idPairs))
		shippingMethodIDs               = make([]string, len(idPairs))
		channelIDs                      = make([]string, len(idPairs))
		shippingMethodChannelListings   []*model.ShippingMethodChannelListing
		appErr                          *model.AppError
		shippingMethodChannelListingMap = map[string]*model.ShippingMethodChannelListing{} // keys are pair of shippingMethodID__channelID formats
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	for idx, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index >= 0 {
			shippingMethodIDs[idx] = pair[:index]
			channelIDs[idx] = pair[index+2:]
		}
	}

	shippingMethodChannelListings, appErr = embedCtx.App.Srv().
		ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
			ChannelID:        squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": channelIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	shippingMethodChannelListingMap = lo.SliceToMap(shippingMethodChannelListings, func(r *model.ShippingMethodChannelListing) (string, *model.ShippingMethodChannelListing) {
		return r.ShippingMethodID + "__" + r.ChannelID, r
	})

	for idx, pair := range idPairs {
		res[idx] = &dataloader.Result[*model.ShippingMethodChannelListing]{Data: shippingMethodChannelListingMap[pair]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[*model.ShippingMethodChannelListing]{Error: err}
	}
	return res
}

func shippingZonesByChannelIdLoader(ctx context.Context, ids []string) []*dataloader.Result[[]*model.ShippingZone] {
	var (
		res                    = make([]*dataloader.Result[[]*model.ShippingZone], len(ids))
		shippingZones          []*model.ShippingZone
		shippingZoneMap        = map[string]*model.ShippingZone{}
		relations              []*model.ShippingZoneChannel
		err                    error
		errs                   []error
		shippingZoneIDs        []string
		channelShippingZoneMap = map[string][]string{} // keys are channel ids, values are slices of shipping zone ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	relations, err = embedCtx.App.Srv().Store.
		ShippingZoneChannel().
		FilterByOptions(&model.ShippingZoneChannelFilterOptions{
			ChannelID: squirrel.Eq{store.ShippingZoneChannelTableName + ".ChannelID": ids},
		})
	if err != nil {
		err = model.NewAppError("ShippingZoneChanenlByOptions", "app.shipping.shippingzone-channel-relations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	for _, rel := range relations {
		shippingZoneIDs = append(shippingZoneIDs, rel.ShippingZoneID)
		channelShippingZoneMap[rel.ChannelID] = append(channelShippingZoneMap[rel.ChannelID], rel.ShippingZoneID)
	}

	shippingZones, errs = ShippingZoneByIdLoader.LoadMany(ctx, shippingZoneIDs)()
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		goto errorLabel
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

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[[]*model.ShippingZone]{Error: err}
	}
	return res
}

func shippingZoneByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingZone] {
	var (
		res             = make([]*dataloader.Result[*model.ShippingZone], len(ids))
		shippingZones   []*model.ShippingZone
		appErr          *model.AppError
		shippingZoneMap = map[string]*model.ShippingZone{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	shippingZones, appErr = embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id: squirrel.Eq{store.ShippingZoneTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	shippingZoneMap = lo.SliceToMap(shippingZones, func(s *model.ShippingZone) (string, *model.ShippingZone) { return s.Id, s })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.ShippingZone]{Data: shippingZoneMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.ShippingZone]{Error: err}
	}
	return res
}

func shippingMethodsByShippingZoneIdLoader(ctx context.Context, shippingZoneIDs []string) []*dataloader.Result[[]*model.ShippingMethod] {
	var (
		res               = make([]*dataloader.Result[[]*model.ShippingMethod], len(shippingZoneIDs))
		shippingMethods   []*model.ShippingMethod
		appErr            *model.AppError
		shippingMethodMap = map[string][]*model.ShippingMethod{} // keys are shipping zone ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	shippingMethods, appErr = embedCtx.App.Srv().
		ShippingService().
		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
			ShippingZoneID: squirrel.Eq{store.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, method := range shippingMethods {
		shippingMethodMap[method.ShippingZoneID] = append(shippingMethodMap[method.ShippingZoneID], method)
	}

	for idx, zoneID := range shippingZoneIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Data: shippingMethodMap[zoneID]}
	}
	return res

errorLabel:
	for idx := range shippingZoneIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Error: err}
	}
	return res
}

func postalCodeRulesByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodPostalCodeRule] {
	var (
		res     = make([]*dataloader.Result[[]*model.ShippingMethodPostalCodeRule], len(shippingMethodIDs))
		rules   []*model.ShippingMethodPostalCodeRule
		appErr  *model.AppError
		ruleMap = map[string][]*model.ShippingMethodPostalCodeRule{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	rules, appErr = embedCtx.App.Srv().
		ShippingService().
		ShippingMethodPostalCodeRulesByOptions(&model.ShippingMethodPostalCodeRuleFilterOptions{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodPostalCodeRuleTableName + ".ShippingMethodID": shippingMethodIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, rule := range rules {
		ruleMap[rule.ShippingMethodID] = append(ruleMap[rule.ShippingMethodID], rule)
	}

	for idx, id := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodPostalCodeRule]{Data: ruleMap[id]}
	}
	return res

errorLabel:
	for idx := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodPostalCodeRule]{Error: err}
	}
	return res
}

func shippingMethodsByShippingZoneIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.ShippingMethod] {
	// 	var (
	// 		res             = make([]*dataloader.Result[[]*model.ShippingMethod], len(idPairs))
	// 		shippingZoneIDs = make([]string, len(idPairs))
	// 		channelIDs      = make([]string, len(idPairs))
	// 		shippingMethods []*model.ShippingMethod
	// 		appErr          *model.AppError
	// 	)

	// 	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	// 	if err != nil {
	// 		goto errorLabel
	// 	}

	// 	for idx, pair := range idPairs {
	// 		index := strings.Index(pair, "__")
	// 		if index >= 0 {
	// 			shippingZoneIDs[idx] = pair[:index]
	// 			channelIDs[idx] = pair[index+2:]
	// 		}
	// 	}

	// 	shippingMethods, appErr = embedCtx.App.Srv().
	// 		ShippingService().
	// 		ShippingMethodsByOptions(&model.ShippingMethodFilterOption{
	// 			ShippingZoneID: squirrel.Eq{store.ShippingMethodTableName + ".ShippingZoneID": shippingZoneIDs},
	// 		})
	// 	if appErr != nil {
	// 		err = appErr
	// 		goto errorLabel
	// 	}

	// errorLabel:
	//
	//	for idx := range idPairs {
	//		res[idx] = &dataloader.Result[[]*model.ShippingMethod]{Error: err}
	//	}
	//	return res
	panic("not implemented")
}

func shippingMethodChannelListingByShippingMethodIdLoader(ctx context.Context, shippingMethodIDs []string) []*dataloader.Result[[]*model.ShippingMethodChannelListing] {
	var (
		res                              = make([]*dataloader.Result[[]*model.ShippingMethodChannelListing], len(shippingMethodIDs))
		shippingMethodChannelListings    []*model.ShippingMethodChannelListing
		shippingMethodChannelListingsMap = map[string][]*model.ShippingMethodChannelListing{}
		appErr                           *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	shippingMethodChannelListings, appErr = embedCtx.App.Srv().ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethodIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, listing := range shippingMethodChannelListings {
		shippingMethodChannelListingsMap[listing.ShippingMethodID] = append(shippingMethodChannelListingsMap[listing.ShippingMethodID], listing)
	}
	for idx, id := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodChannelListing]{Data: shippingMethodChannelListingsMap[id]}
	}
	return res

errorLabel:
	for idx := range shippingMethodIDs {
		res[idx] = &dataloader.Result[[]*model.ShippingMethodChannelListing]{Error: err}
	}
	return res
}
