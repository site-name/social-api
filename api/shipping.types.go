package api

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type ShippingMethod struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Description         JSONString              `json:"description"`
	MinimumOrderWeight  *Weight                 `json:"minimumOrderWeight"`
	MaximumOrderWeight  *Weight                 `json:"maximumOrderWeight"`
	MaximumDeliveryDays *int32                  `json:"maximumDeliveryDays"`
	MinimumDeliveryDays *int32                  `json:"minimumDeliveryDays"`
	PrivateMetadata     []*MetadataItem         `json:"privateMetadata"`
	Metadata            []*MetadataItem         `json:"metadata"`
	Type                *ShippingMethodTypeEnum `json:"type"`

	shippingZoneID string

	// Translation         *ShippingMethodTranslation `json:"translation"`
	// ChannelListings     []*ShippingMethodChannelListing `json:"channelListings"`
	// Price               *Money                          `json:"price"`
	// MaximumOrderPrice   *Money                          `json:"maximumOrderPrice"`
	// MinimumOrderPrice   *Money                          `json:"minimumOrderPrice"`
	// PostalCodeRules     []*ShippingMethodPostalCodeRule `json:"postalCodeRules"`
	// ExcludedProducts    *ProductCountableConnection     `json:"excludedProducts"`
}

func SystemShippingMethodToGraphqlShippingMethod(m *model.ShippingMethod) *ShippingMethod {
	if m == nil {
		return nil
	}

	res := &ShippingMethod{
		ID:             m.Id,
		Name:           m.Name,
		Description:    JSONString(m.Description),
		shippingZoneID: m.ShippingZoneID,
		MinimumOrderWeight: &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(m.MinimumOrderWeight),
		},
		PrivateMetadata: MetadataToSlice(m.PrivateMetadata),
		Metadata:        MetadataToSlice(m.Metadata),
		Type:            (*ShippingMethodTypeEnum)(&m.Type),
	}

	if m.MaximumOrderWeight != nil {
		res.MaximumOrderWeight = &Weight{
			Unit:  WeightUnitsEnum(m.WeightUnit),
			Value: float64(*m.MaximumOrderWeight),
		}
	}

	if m.MaximumDeliveryDays != nil {
		res.MaximumDeliveryDays = model.NewInt32(int32(*m.MaximumDeliveryDays))
	}
	if m.MinimumDeliveryDays != nil {
		res.MinimumDeliveryDays = model.NewInt32(int32(*m.MinimumDeliveryDays))
	}

	return res
}

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

func shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ShippingMethodChannelListing] {
	panic("not implemented")
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

	shippingZones, errs = dataloaders.ShippingZoneByIdLoader.LoadMany(ctx, shippingZoneIDs)()
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
