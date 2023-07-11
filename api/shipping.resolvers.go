package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/shipping.graphqls for details on directive used.
func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingMethodChannelListingInput
}) (*ShippingMethodChannelListingUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method channel listing id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Input.RemoveChannels, model.IsValidId) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "remove channels"}, "please provide valid remove channel ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channelIdsRequesterWantToAdd := lo.Map(args.Input.AddChannels, func(item *ShippingMethodChannelListingAddInput, _ int) string { return item.ChannelID })
	if !lo.EveryBy(channelIdsRequesterWantToAdd, model.IsValidId) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide valid add channel ids", http.StatusBadRequest)
	}

	// find channels that requester want to add
	channels, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": channelIdsRequesterWantToAdd},
	})
	channelsAboutToAddMap := lo.SliceToMap(channels, func(ch *model.Channel) (string, *model.Channel) { return ch.Id, ch })

	// find shipping method given by user
	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Id: squirrel.Eq{store.ShippingMethodTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if channels to add are assigned to given shipping method:
	channelsOfShippingMethod, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{
		ShippingZoneChannels_ShippingZoneID: squirrel.Eq{store.ShippingZoneChannelTableName + ".ShippingZoneID": shippingMethod.ShippingZoneID},
	})
	if appErr != nil {
		return nil, appErr
	}

	channelIDsNotAssignedToShippingMethod, _ := lo.Difference(channelIdsRequesterWantToAdd, channelsOfShippingMethod.IDs())
	if len(channelIDsNotAssignedToShippingMethod) > 0 {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide channels that are assigned to the shipping method", http.StatusBadRequest)
	}

	// keep only add channels that have ShippingMethodChannelListing relations with shipping method
	shippingMethodChannelListings, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": args.Id},
		ChannelID:        squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": channelIdsRequesterWantToAdd},
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		shippingMethodChannelListingsMap         = make(map[string]string, len(shippingMethodChannelListings)) // keys have form ofchannelID+shippingMethodID, values are shipping method channel listing ids
		channelIDsOfShippingMethodChannelListing = make(util.AnyArray[string], len(shippingMethodChannelListings))
	)
	for idx, listing := range shippingMethodChannelListings {
		shippingMethodChannelListingsMap[listing.ChannelID+listing.ShippingMethodID] = listing.Id
		channelIDsOfShippingMethodChannelListing[idx] = listing.ChannelID
	}

	channelIDsCanAdd := channelIDsOfShippingMethodChannelListing.InterSection(channelIdsRequesterWantToAdd...)
	actualChannelIdsCanAddMap := lo.SliceToMap(channelIDsCanAdd, func(id string) (string, bool) { return id, true })

	for _, channelInput := range args.Input.AddChannels {
		if channelInput.Price == nil && !actualChannelIdsCanAddMap[channelInput.ChannelID] {
			return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channel"}, "price field is required", http.StatusBadRequest)
		}

		// validate given max order price > min order price
		if channelInput.MinimumOrderPrice != nil &&
			channelInput.MaximumOrderPrice != nil &&
			channelInput.MaximumOrderPrice.LessThanOrEqual(channelInput.MinimumOrderPrice) {
			return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "max prices must be greater than min prices", http.StatusBadRequest)
		}
	}

	// update
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	// add channels:
	shippingMethodChannelListingsToAdd := lo.Map(args.Input.AddChannels, func(item *ShippingMethodChannelListingAddInput, _ int) *model.ShippingMethodChannelListing {
		channel := channelsAboutToAddMap[item.ChannelID]

		for _, price := range [...]*PositiveDecimal{
			item.Price,
			item.MinimumOrderPrice,
			item.MaximumOrderPrice,
		} {
			if price != nil {
				money := &goprices.Money{Amount: decimal.Decimal(*price), Currency: channel.Currency}
				// NOTE: we don't validate price amount precisions for currencies here.
				// call Quantize() to round prices for currencies properly
				money, _ = money.Quantize(nil, goprices.Up)
				*price = *(*PositiveDecimal)(unsafe.Pointer(&money.Amount))
			}
		}

		return &model.ShippingMethodChannelListing{
			Id: shippingMethodChannelListingsMap[item.ChannelID+args.Id], // this will be an uuid or "", helps store to decide weather to update or insert

			Currency:                channelsAboutToAddMap[item.ChannelID].Currency,
			MinimumOrderPriceAmount: (*decimal.Decimal)(unsafe.Pointer(item.MinimumOrderPrice)),
			MaximumOrderPriceAmount: (*decimal.Decimal)(unsafe.Pointer(item.MaximumOrderPrice)),
			PriceAmount:             (*decimal.Decimal)(unsafe.Pointer(item.Price)),
			ShippingMethodID:        args.Id,
			ChannelID:               item.ChannelID,
		}
	})
	_, appErr = embedCtx.App.Srv().ShippingService().UpsertShippingMethodChannelListings(transaction, shippingMethodChannelListingsToAdd)
	if appErr != nil {
		return nil, appErr
	}

	// remove:
	appErr = embedCtx.App.Srv().ShippingService().DeleteShippingMethodChannelListings(transaction, &model.ShippingMethodChannelListingFilterOption{
		ChannelID:        squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": args.Input.RemoveChannels},
		ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingMethodChannelListingUpdate{
		ShippingMethod: SystemShippingMethodToGraphqlShippingMethod(shippingMethod),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directive used.
func (r *Resolver) ShippingPriceCreate(ctx context.Context, args struct{ Input ShippingPriceInput }) (*ShippingPriceCreate, error) {
	appErr := args.Input.Validate("ShippingPriceCreate")
	if appErr != nil {
		return nil, appErr
	}

	shippingMethod := new(model.ShippingMethod)
	args.Input.Patch(shippingMethod)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// start transaction:
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("ShippingPriceCreate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	shippingMethod, appErr = embedCtx.App.Srv().ShippingService().UpsertShippingMethod(transaction, shippingMethod)
	if appErr != nil {
		return nil, appErr
	}

	shippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Id: squirrel.Eq{store.ShippingZoneTableName + ".Id": shippingMethod.ShippingZoneID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("ShippingPriceCreate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceCreate{
		ShippingMethod: SystemShippingMethodToGraphqlShippingMethod(shippingMethod),
		ShippingZone:   SystemShippingZoneToGraphqlShippingZone(shippingZones[0]),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceDelete(ctx context.Context, args struct{ Id string }) (*ShippingPriceDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingPriceDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Id:                        squirrel.Eq{store.ShippingMethodTableName + ".Id": args.Id},
		SelectRelatedShippingZone: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	err := embedCtx.App.Srv().Store.ShippingMethod().Delete(nil, args.Id)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceDelete", "app.shipping.delete_shipping_method.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceDelete{
		ShippingMethod: SystemShippingMethodToGraphqlShippingMethod(shippingMethod),
		ShippingZone:   SystemShippingZoneToGraphqlShippingZone(shippingMethod.GetShippingZone()),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ShippingPriceBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.ShippingMethod().Delete(nil, args.Ids...)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceDelete", "app.shipping.delete_shipping_methods.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceBulkDelete{
		Count: int32(len(args.Ids)),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceInput
}) (*ShippingPriceUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingPriceUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	appErr := args.Input.Validate("ShippingPriceUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// start transaction:
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("ShippingPriceUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Id:                        squirrel.Eq{store.ShippingMethodTableName + ".Id": args.Id},
		SelectRelatedShippingZone: true, //
	})
	if appErr != nil {
		return nil, appErr
	}
	if args.Input.Patch(shippingMethod) {
		shippingMethod, appErr = embedCtx.App.Srv().Shipping.UpsertShippingMethod(transaction, shippingMethod)
		if err != nil {
			return nil, appErr
		}
	}

	if len(args.Input.DeletePostalCodeRules) > 0 {
		err = embedCtx.App.Srv().Store.ShippingMethodPostalCodeRule().Delete(transaction, args.Input.DeletePostalCodeRules...)
		if err != nil {
			return nil, model.NewAppError("ShippingPriceUpdate", "app.shipping.error_delete_shipping_method_postal_code_rules.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if len(args.Input.AddPostalCodeRules) > 0 {
		rules := lo.Map(args.Input.AddPostalCodeRules, func(item *ShippingPostalCodeRulesCreateInputRange, _ int) *model.ShippingMethodPostalCodeRule {
			rule := &model.ShippingMethodPostalCodeRule{
				ShippingMethodID: args.Id,
				Start:            item.Start,
			}
			if item.End != nil {
				rule.End = *item.End
			}
			if args.Input.InclusionType != nil {
				rule.InclusionType = *args.Input.InclusionType
			}

			return rule
		})

		_, appErr = embedCtx.App.Srv().ShippingService().CreateShippingMethodPostalCodeRules(transaction, rules)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("ShippingPriceUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// NOTE: The returning shipping zone below is old version of parent shipping zone of shipping method

	return &ShippingPriceUpdate{
		ShippingMethod: SystemShippingMethodToGraphqlShippingMethod(shippingMethod),
		ShippingZone:   SystemShippingZoneToGraphqlShippingZone(shippingMethod.GetShippingZone()),
	}, nil
}

func (r *Resolver) ShippingPriceTranslate(ctx context.Context, args struct {
	Id           string
	Input        ShippingPriceTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ShippingPriceTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceExcludeProducts(ctx context.Context, args struct {
	Id    string
	Input ShippingPriceExcludeProductsInput
}) (*ShippingPriceExcludeProducts, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingPriceExcludeProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Input.Products, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceExcludeProducts", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "product ids"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	for _, productID := range args.Input.Products {
		// NOTE: no need to worry unique constraint violation since it's handled in store layer.
		_, err := embedCtx.App.Srv().Store.ShippingMethodExcludedProduct().Save(&model.ShippingMethodExcludedProduct{
			ShippingMethodID: args.Id,
			ProductID:        productID,
		})
		if err != nil {
			return nil, model.NewAppError("ShippingPriceExcludeProducts", "app.shipping.insert_shipping_method_excluded_product.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return &ShippingPriceExcludeProducts{
		ShippingMethod: &ShippingMethod{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingPriceRemoveProductFromExclude(ctx context.Context, args struct {
	Id       string
	Products []string
}) (*ShippingPriceRemoveProductFromExclude, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingPriceRemoveProductFromExclude", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Products, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceRemoveProductFromExclude", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "product ids"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.ShippingMethodExcludedProduct().Delete(nil, &model.ShippingMethodExcludedProductFilterOptions{
		ShippingMethodID: squirrel.Eq{store.ShippingMethodExcludedProductTableName + ".ShippingMethodID": args.Id},
		ProductID:        squirrel.Eq{store.ShippingMethodExcludedProductTableName + ".ProductID": args.Products},
	})
	if err != nil {
		return nil, model.NewAppError("ShippingPriceRemoveProductFromExclude", "app.shipping.delete_shipping_method_excluded_products.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceRemoveProductFromExclude{
		ShippingMethod: &ShippingMethod{ID: args.Id},
	}, nil
}

func (r *Resolver) ShippingZoneCreate(ctx context.Context, args struct {
	Input ShippingZoneCreateInput
}) (*ShippingZoneCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZoneDelete(ctx context.Context, args struct{ Id string }) (*ShippingZoneDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingZoneDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

}

func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ShippingZoneBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZoneUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingZoneUpdateInput
}) (*ShippingZoneUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZone(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*ShippingZone, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ShippingZones(ctx context.Context, args struct {
	Filter  *ShippingZoneFilterInput
	Channel *string
	GraphqlParams
}) (*ShippingZoneCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
