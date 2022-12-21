package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func shippingZonesByChannelIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*ShippingZone] {
	panic("not implemented")
}

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

func shippingMethodByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*ShippingMethod] {
	results := shippingMethodByIdLoader_SystemResult(ctx, ids)
	res := make([]*dataloader.Result[*ShippingMethod], len(results))

	for idx, result := range results {
		res[idx] = &dataloader.Result[*ShippingMethod]{
			Data:  SystemShippingMethodToGraphqlShippingMethod(result.Data),
			Error: result.Error,
		}
	}

	return res
}

// shippingMethodByIdLoaderSystemResult returns slice of dataloader results that have data fields are system shipping methods
func shippingMethodByIdLoader_SystemResult(ctx context.Context, ids []string) []*dataloader.Result[*model.ShippingMethod] {
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

func shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader(ctx context.Context, idPairs []string) []*dataloader.Result[*ShippingMethodChannelListing] {
	panic("not implemented")
}

func shippingMethodChannelListingByShippingMethodIdAndChannelSlugLoader_systemResult(ctx context.Context, idPairs []string) []*dataloader.Result[*model.ShippingMethodChannelListing] {
	panic("not implemented")
}
