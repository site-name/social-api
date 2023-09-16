package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/shipping.graphqls for details on directive used.
func (r *Resolver) ShippingMethodChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingMethodChannelListingInput
}) (*ShippingMethodChannelListingUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method channel listing id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Input.RemoveChannels, model.IsValidId) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "remove channels"}, "please provide valid remove channel ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channelIdsRequesterWantToAdd := lo.Map(args.Input.AddChannels, func(item *ShippingMethodChannelListingAddInput, _ int) string { return item.ChannelID })
	if !lo.EveryBy(channelIdsRequesterWantToAdd, model.IsValidId) {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide valid add channel ids", http.StatusBadRequest)
	}

	// find channels that requester want to add
	channels, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": channelIdsRequesterWantToAdd},
	})
	channelsAboutToAddMap := lo.SliceToMap(channels, func(ch *model.Channel) (string, *model.Channel) { return ch.Id, ch })

	// find shipping method given by user
	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Conditions: squirrel.Eq{model.ShippingMethodTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if channels to add are assigned to given shipping method:
	channelsOfShippingMethod, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{
		ShippingZoneChannels_ShippingZoneID: squirrel.Eq{model.ShippingZoneChannelTableName + ".ShippingZoneID": shippingMethod.ShippingZoneID},
	})
	if appErr != nil {
		return nil, appErr
	}

	channelIDsNotAssignedToShippingMethod, _ := lo.Difference(channelIdsRequesterWantToAdd, channelsOfShippingMethod.IDs())
	if len(channelIDsNotAssignedToShippingMethod) > 0 {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide channels that are assigned to the shipping method", http.StatusBadRequest)
	}

	// keep only add channels that have ShippingMethodChannelListing relations with shipping method
	shippingMethodChannelListings, appErr := embedCtx.App.Srv().
		ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ShippingMethodChannelListingTableName + ".ShippingMethodID": args.Id,
				model.ShippingMethodChannelListingTableName + ".ChannelID":        channelIdsRequesterWantToAdd,
			},
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

	channelIDsCanAdd := channelIDsOfShippingMethodChannelListing.InterSection(channelIdsRequesterWantToAdd)
	actualChannelIdsCanAddMap := lo.SliceToMap(channelIDsCanAdd, func(id string) (string, bool) { return id, true })

	for _, channelInput := range args.Input.AddChannels {
		if channelInput.Price == nil && !actualChannelIdsCanAddMap[channelInput.ChannelID] {
			return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channel"}, "price field is required", http.StatusBadRequest)
		}

		// validate given max order price > min order price
		if channelInput.MinimumOrderPrice != nil &&
			channelInput.MaximumOrderPrice != nil &&
			channelInput.MaximumOrderPrice.ToDecimal().LessThanOrEqual(channelInput.MinimumOrderPrice.ToDecimal()) {
			return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "max prices must be greater than min prices", http.StatusBadRequest)
		}
	}

	// update
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	defer transaction.Rollback()

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
				money, _ = money.Quantize(goprices.Up, -1)
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
		Conditions: squirrel.Eq{
			model.ShippingMethodChannelListingTableName + ".ChannelID":        args.Input.RemoveChannels,
			model.ShippingMethodChannelListingTableName + ".ShippingMethodID": args.Id,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ShippingMethodChannelListingUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	defer transaction.Rollback()

	shippingMethod, appErr = embedCtx.App.Srv().ShippingService().UpsertShippingMethod(transaction, shippingMethod)
	if appErr != nil {
		return nil, appErr
	}

	shippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{
		Conditions: squirrel.Eq{model.ShippingZoneTableName + ".Id": shippingMethod.ShippingZoneID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ShippingPriceCreate", model.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("ShippingPriceDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Conditions:                squirrel.Eq{model.ShippingMethodTableName + ".Id": args.Id},
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
		return nil, model.NewAppError("ShippingPriceBulkDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid ids", http.StatusBadRequest)
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
		return nil, model.NewAppError("ShippingPriceUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid id", http.StatusBadRequest)
	}

	appErr := args.Input.Validate("ShippingPriceUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// start transaction:
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	defer transaction.Rollback()

	shippingMethod, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
		Conditions:                squirrel.Eq{model.ShippingMethodTableName + ".Id": args.Id},
		SelectRelatedShippingZone: true, //
	})
	if appErr != nil {
		return nil, appErr
	}
	if args.Input.Patch(shippingMethod) {
		shippingMethod, appErr = embedCtx.App.Srv().Shipping.UpsertShippingMethod(transaction, shippingMethod)
		if appErr != nil {
			return nil, appErr
		}
	}

	if len(args.Input.DeletePostalCodeRules) > 0 {
		err := embedCtx.App.Srv().Store.ShippingMethodPostalCodeRule().Delete(transaction, args.Input.DeletePostalCodeRules...)
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
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ShippingPriceUpdate", model.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("ShippingPriceExcludeProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Input.Products, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceExcludeProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "product ids"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productsToExclude := lo.Map(args.Input.Products, func(id string, _ int) *model.Product { return &model.Product{Id: id} })
	err := embedCtx.App.Srv().Store.GetMaster().
		Model(&model.ShippingMethod{Id: args.Id}).
		Association("ExcludedProducts").
		Append(productsToExclude)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceExcludeProducts", "app.shipping.insert_shipping_method_excluded_product.app_error", nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("ShippingPriceRemoveProductFromExclude", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping method id", http.StatusBadRequest)
	}

	if !lo.EveryBy(args.Products, model.IsValidId) {
		return nil, model.NewAppError("ShippingPriceRemoveProductFromExclude", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "product ids"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productsToRemove := lo.Map(args.Products, func(id string, _ int) *model.Product { return &model.Product{Id: id} })
	err := embedCtx.App.Srv().Store.GetMaster().
		Model(&model.ShippingMethod{Id: args.Id}).
		Association("ExcludedProducts").
		Delete(productsToRemove)
	if err != nil {
		return nil, model.NewAppError("ShippingPriceExcludeProducts", "app.shipping.insert_shipping_method_excluded_product.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingPriceRemoveProductFromExclude{
		ShippingMethod: &ShippingMethod{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZoneCreate(ctx context.Context, args struct {
	Input ShippingZoneCreateInput
}) (*ShippingZoneCreate, error) {
	// validate params
	if !lo.EveryBy(args.Input.AddWarehouses, model.IsValidId) {
		return nil, model.NewAppError("ShippingZoneCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add warehouses"}, "please provide valid warehouse ids", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.AddChannels, model.IsValidId) {
		return nil, model.NewAppError("ShippingZoneCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide valid channel ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	defer transaction.Rollback()

	shippingZone := &model.ShippingZone{
		Default:   args.Input.Default,
		Countries: strings.Join(args.Input.Countries, " "),
	}
	if val := args.Input.Name; val != nil {
		shippingZone.Name = *val
	}
	if val := args.Input.Description; val != nil {
		shippingZone.Description = *val
	}

	shippingZone, appErr := embedCtx.App.Srv().Shipping.UpsertShippingZone(transaction, shippingZone)
	if appErr != nil {
		return nil, appErr
	}

	// save m2m warehouse shipping zones
	if len(args.Input.AddWarehouses) > 0 {
		warehousesToAdd := lo.Map(args.Input.AddWarehouses, func(id string, _ int) *model.WareHouse { return &model.WareHouse{Id: id} })
		err := transaction.Model(shippingZone).Association("Warehouses").Append(warehousesToAdd)
		if err != nil {
			return nil, model.NewAppError("ShippingZoneCreate", "app.shipping.add_warehouse_shipping_zones.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// save m2m shipping zone channels
	if len(args.Input.AddChannels) > 0 {
		channelsToAdd := lo.Map(args.Input.AddChannels, func(id string, _ int) *model.Channel { return &model.Channel{Id: id} })
		err := transaction.Model(shippingZone).Association("Channels").Append(channelsToAdd)
		if err != nil {
			return nil, model.NewAppError("ShippingZoneCreate", "app.shipping.add_channel_shipping_zones.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ShippingZoneCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingZoneCreate{
		ShippingZone: SystemShippingZoneToGraphqlShippingZone(shippingZone),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZoneDelete(ctx context.Context, args struct{ Id string }) (*ShippingZoneDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingZoneDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping zone id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, err := embedCtx.App.Srv().ShippingService().DeleteShippingZones(nil, &model.ShippingZoneFilterOption{
		Conditions: squirrel.Eq{model.ShippingZoneTableName + ".Id": args.Id},
	})
	if err != nil {
		return nil, err
	}

	// NOTE: ShippingZoneChannels and WarehouseShippingZones are auto deleted thanks to ON DELETE CASCADE options

	return &ShippingZoneDelete{
		ShippingZone: &ShippingZone{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZoneBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ShippingZoneBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("ShippingZoneCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid shipping zone ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	numDeleted, err := embedCtx.App.Srv().ShippingService().DeleteShippingZones(nil, &model.ShippingZoneFilterOption{
		Conditions: squirrel.Eq{model.ShippingZoneTableName + ".Id": args.Ids},
	})
	if err != nil {
		return nil, err
	}

	// NOTE: ShippingZoneChannels and WarehouseShippingZones are auto deleted thanks to ON DELETE CASCADE options

	return &ShippingZoneBulkDelete{
		Count: int32(numDeleted),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZoneUpdate(ctx context.Context, args struct {
	Id    string
	Input ShippingZoneUpdateInput
}) (*ShippingZoneUpdate, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingZoneUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping zone id", http.StatusBadRequest)
	}

	for argsName, ids := range map[string][]string{
		"addWarehouses":    args.Input.AddWarehouses,
		"addChannels":      args.Input.AddChannels,
		"removeWarehouses": args.Input.RemoveWarehouses,
		"removeChannels":   args.Input.RemoveChannels,
	} {
		if !lo.EveryBy(ids, model.IsValidId) {
			return nil, model.NewAppError("ShippingZoneUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": argsName}, "please provide valid "+argsName, http.StatusBadRequest)
		}
	}

	if len(lo.Intersect(args.Input.AddWarehouses, args.Input.RemoveWarehouses)) > 0 {
		return nil, model.NewAppError("ShippingZoneUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "AddWarehouses / RemoveWarehouses"}, "add warehouses must differ from remove warehouses", http.StatusBadRequest)
	}
	if len(lo.Intersect(args.Input.AddChannels, args.Input.RemoveChannels)) > 0 {
		return nil, model.NewAppError("ShippingZoneUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "AddChannels / RemoveChannels"}, "add channels must differ from remove channels", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("ShippingZoneUpdate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// find all shipping zones
	allShippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(&model.ShippingZoneFilterOption{})
	if appErr != nil {
		return nil, appErr
	}
	if len(allShippingZones) == 0 {
		return nil, model.NewAppError("ShippingZoneUpdate", "app.shipping.no_shipping_zone.app_error", nil, "system has no shipping zone", http.StatusNotImplemented)
	}
	shippingZoneToUpdate, found := lo.Find(allShippingZones, func(sp *model.ShippingZone) bool { return sp.Id == args.Id })
	if !found {
		return nil, model.NewAppError("ShippingZoneUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping zone id", http.StatusBadRequest)
	}

	// clean default:
	// check if:
	// 1) requester want to set given shipping zone as default.
	// 2) there is already a shippingzone with `Default` set to true.
	// If 2 conditions are met, return error
	if val := args.Input.Default; val != nil && *val {
		if lo.SomeBy(allShippingZones, func(sp *model.ShippingZone) bool { return sp.Id != args.Id && sp.Default != nil && *sp.Default }) {
			return nil, model.NewAppError("ShippingZoneUpdate", "app.shipping.default_shipping_zone_exists.app_error", nil, "default shipping zone exists", http.StatusNotModified)
		}
		// find all countries code that are not used by any shipping zones
		usedCountriesByShippingZones := map[string]struct{}{}
		for _, spz := range allShippingZones {
			for _, country := range strings.Fields(spz.Countries) {
				usedCountriesByShippingZones[country] = struct{}{}
			}
		}

		countriesNotUsedByShippingZones := []string{}
		for country := range model.Countries {
			countryCode := string(country)
			_, used := usedCountriesByShippingZones[countryCode]
			if !used {
				countriesNotUsedByShippingZones = append(countriesNotUsedByShippingZones, countryCode)
			}
		}

		args.Input.Countries = countriesNotUsedByShippingZones
	}

	// update shipping zone:
	if val := args.Input.Name; val != nil {
		shippingZoneToUpdate.Name = *val
	}
	if val := args.Input.Description; val != nil {
		shippingZoneToUpdate.Description = *val
	}
	shippingZoneToUpdate.Default = args.Input.Default
	shippingZoneToUpdate.Countries = strings.Join(args.Input.Countries, " ")

	shippingZoneToUpdate, appErr = embedCtx.App.Srv().Shipping.UpsertShippingZone(transaction, shippingZoneToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	// add channels, warehouses to shipping zones
	appErr = embedCtx.App.Srv().ShippingService().ToggleShippingZoneRelations(transaction, model.ShippingZones{{Id: args.Id}}, args.Input.AddWarehouses, args.Input.AddChannels, false)
	if appErr != nil {
		return nil, appErr
	}

	// remove channels, warehouses from shipping zones
	appErr = embedCtx.App.Srv().ShippingService().ToggleShippingZoneRelations(transaction, model.ShippingZones{{Id: args.Id}}, args.Input.RemoveWarehouses, args.Input.RemoveChannels, true)
	if appErr != nil {
		return nil, appErr
	}

	if len(args.Input.RemoveChannels) > 0 {
		shippingChannelListings, appErr := embedCtx.App.Srv().ShippingService().ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			ShippingMethod_ShippingZoneID_Inner: squirrel.Eq{model.ShippingZoneTableName + ".Id": args.Id},
			Conditions:                          squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ChannelID": args.Input.RemoveChannels},
		})
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().ShippingService().DeleteShippingMethodChannelListings(transaction, &model.ShippingMethodChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ShippingMethodChannelListingTableName + ".Id": shippingChannelListings.IDs()},
		})
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().ShippingService().DropInvalidShippingMethodsRelationsForGivenChannels(transaction, shippingChannelListings.ShippingMethodIDs(), args.Input.RemoveChannels)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ShippingZoneUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ShippingZoneUpdate{
		ShippingZone: SystemShippingZoneToGraphqlShippingZone(shippingZoneToUpdate),
	}, nil
}

// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZone(ctx context.Context, args struct {
	Id      string
	Channel *string // TODO: Check if we need this
}) (*ShippingZone, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ShippingZone", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid shipping zone id", http.StatusBadRequest)
	}

	zone, err := ShippingZoneByIdLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}

	return SystemShippingZoneToGraphqlShippingZone(zone), nil
}

// NOTE: shipping zones order by Names
//
// NOTE: Refer to ./schemas/shipping.graphqls for details on directives used.
func (r *Resolver) ShippingZones(ctx context.Context, args struct {
	Filter    *ShippingZoneFilterInput
	ChannelID *string
	GraphqlParams
}) (*ShippingZoneCountableConnection, error) {
	channelIDs := []string{}

	shippingZoneFilterOpts := &model.ShippingZoneFilterOption{}

	if args.ChannelID != nil {
		if !model.IsValidId(*args.ChannelID) {
			return nil, model.NewAppError("ShippingZones", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel id"}, "please provide valid channel id", http.StatusBadRequest)
		}
		channelIDs = append(channelIDs, *args.ChannelID)
	}

	appErr := args.GraphqlParams.validate("ShippingZones")
	if appErr != nil {
		return nil, appErr
	}

	if args.Filter != nil {
		if !lo.EveryBy(args.Filter.Channels, model.IsValidId) {
			return nil, model.NewAppError("ShippingZones", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channels"}, "please provide valid channel id", http.StatusBadRequest)
		}
		channelIDs = append(channelIDs, args.Filter.Channels...)

		if search := args.Filter.Search; search != nil &&
			*search != "" &&
			!stringsContainSqlExpr.MatchString(*search) { // NOTE: this check is needed
			shippingZoneFilterOpts.Conditions = squirrel.ILike{model.ShippingZoneTableName + ".Name": "%" + *search + "%"}
		}
	}

	if len(channelIDs) > 0 {
		shippingZoneFilterOpts.ChannelID = squirrel.Eq{model.ShippingZoneChannelTableName + ".ChannelID": channelIDs}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	shippingZones, appErr := embedCtx.App.Srv().ShippingService().ShippingZonesByOption(shippingZoneFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := func(sz *model.ShippingZone) []any { return []any{model.ShippingZoneTableName + ".Name", sz.Name} }
	res, appErr := newGraphqlPaginator(shippingZones, keyFunc, SystemShippingZoneToGraphqlShippingZone, args.GraphqlParams).parse("ShippingZones")
	if appErr != nil {
		return nil, appErr
	}

	return (*ShippingZoneCountableConnection)(unsafe.Pointer(res)), nil
}
