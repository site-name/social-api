package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
	"gorm.io/gorm"
)

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelCreate(ctx context.Context, args struct{ Input ChannelCreateInput }) (*ChannelCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate input
	if !lo.EveryBy(args.Input.AddShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "addShippingZones"}, "please provide valid addShippingZones", http.StatusBadRequest)
	}
	if !args.Input.DefaultCountry.IsValid() {
		return nil, model.NewAppError("ChannelCreate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "defaultCountry"}, fmt.Sprintf("%s is not valid country code", args.Input.DefaultCountry), http.StatusBadRequest)
	}

	channel := &model.Channel{
		Name:     args.Input.Name,
		Currency: args.Input.CurrencyCode,
	}
	if val := args.Input.IsActive; val != nil {
		channel.IsActive = *val
	}
	if val := args.Input.Slug; val != "" && slug.IsSlug(val) {
		channel.Slug = val
	}

	// save new channel to db
	channel, appErr := embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
	if appErr != nil {
		return nil, appErr
	}

	// save channel shipping zone relations:
	// begin transaction
	transaction, err := embedCtx.App.Srv().Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("ChannelCreate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	shippingZoneChannelRelations := lo.Map(args.Input.AddShippingZones, func(id string, _ int) *model.ShippingZoneChannel {
		return &model.ShippingZoneChannel{
			ShippingZoneID: id,
			ChannelID:      channel.Id,
		}
	})
	_, appErr = embedCtx.App.Srv().ChannelService().BulkUpsertShippingZoneChannels(transaction, shippingZoneChannelRelations)
	if appErr != nil {
		return nil, appErr
	}

	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("ChannelCreate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ChannelCreate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelUpdate(ctx context.Context, args struct {
	Id    string
	Input ChannelUpdateInput
}) (*ChannelUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate inputs
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.AddShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "AddShippingZones"}, "please provide valid ids", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.RemoveShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "RemoveShippingZones"}, "please provide valid ids", http.StatusBadRequest)
	}
	intersectIds := lo.Intersect(args.Input.RemoveShippingZones, args.Input.AddShippingZones)
	if len(intersectIds) > 0 {
		return nil, model.NewAppError("ChannelUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "RemoveShippingZones/AddShippingZones"}, "remove shipping zone ids and add shipping zone ids can not have same ids", http.StatusBadRequest)
	}

	// validate if channe does exist
	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	// update field values
	if val := args.Input.Name; val != nil && *val != "" {
		channel.Name = *val
	}
	if val := args.Input.Slug; val != nil && *val != "" && slug.IsSlug(*val) {
		channel.Slug = *val
	}
	if val := args.Input.IsActive; val != nil {
		channel.IsActive = *val
	}
	if val := args.Input.DefaultCountry; val != nil && val.IsValid() {
		channel.DefaultCountry = *val
	}

	// update channel in db
	channel, appErr = embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
	if appErr != nil {
		return nil, appErr
	}

	// update m2m relations
	// begin transaction
	transaction, err := embedCtx.App.Srv().Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("ChannelUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	// create relations between shipping zones and channel
	if len(args.Input.AddShippingZones) > 0 {
		var insertRelations = make([]*model.ShippingZoneChannel, len(args.Input.AddShippingZones))
		for idx, shipZoneID := range args.Input.AddShippingZones {
			insertRelations[idx] = &model.ShippingZoneChannel{ShippingZoneID: shipZoneID, ChannelID: channel.Id}
		}

		_, appErr = embedCtx.App.Srv().ChannelService().BulkUpsertShippingZoneChannels(transaction, insertRelations)
		if appErr != nil {
			return nil, appErr
		}
	}

	if len(args.Input.RemoveShippingZones) > 0 {
		// delete relations between shipping zones and channel
		var deleteRelations = make([]*model.ShippingZoneChannel, len(args.Input.RemoveShippingZones))
		for idx, shipZoneID := range args.Input.RemoveShippingZones {
			deleteRelations[idx] = &model.ShippingZoneChannel{ShippingZoneID: shipZoneID, ChannelID: channel.Id}
		}

		appErr = embedCtx.App.Srv().ChannelService().BulkDeleteShippingZoneChannels(transaction, &model.ShippingZoneChannelFilterOptions{
			Conditions: squirrel.And{
				squirrel.Eq{model.ShippingZoneChannelTableName + ".ShippingZoneID": args.Input.AddShippingZones},
				squirrel.Eq{model.ShippingZoneChannelTableName + ".ChannelID": args.Id},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		// delete shipping methods channel listings of shipping methods of deleted shipping zones
		shippingMethodChannelListings, appErr := embedCtx.App.Srv().ShippingService().
			ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
				ChannelID:                           squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ChannelID": channel.Id},
				ShippingMethod_ShippingZoneID_Inner: squirrel.Eq{model.ShippingZoneTableName + ".Id": args.Input.RemoveShippingZones},
			})
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().ShippingService().DeleteShippingMethodChannelListings(transaction, &model.ShippingMethodChannelListingFilterOption{
			Id: squirrel.Eq{model.ShippingMethodChannelListingTableName + ".Id": shippingMethodChannelListings.IDs()},
		})
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().ShippingService().DropInvalidShippingMethodsRelationsForGivenChannels(transaction, shippingMethodChannelListings.ShippingMethodIDs(), []string{channel.Id})
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("ChannelUpdate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ChannelUpdate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelDelete(ctx context.Context, args struct {
	Id    string
	Input *ChannelDeleteInput
}) (*ChannelDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate input
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide a valid channel id", http.StatusBadRequest)
	}
	if args.Input != nil && !model.IsValidId(args.Input.ChannelID) {
		return nil, model.NewAppError("ChannelDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channelID"}, "please provide a valid channel id", http.StatusBadRequest)
	}
	if args.Input != nil && args.Input.ChannelID == args.Id {
		return nil, model.NewAppError("ChannelDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channelID"}, "target channel cannot be the channel to be deleted", http.StatusBadRequest)
	}

	deleteCheckoutsByChannelID := func(channelID string, transaction *gorm.DB) *model.AppError {
		return embedCtx.App.Srv().
			CheckoutService().
			DeleteCheckoutsByOption(transaction, &model.CheckoutFilterOption{
				ChannelID: squirrel.Eq{model.CheckoutTableName + ".ChannelID": args.Id},
			})
	}

	orders, appErr := embedCtx.App.Srv().
		OrderService().
		FilterOrdersByOptions(&model.OrderFilterOption{
			ChannelID:       squirrel.Eq{model.OrderTableName + ".ChannelID": args.Id},
			SelectForUpdate: true,
		})
	if appErr != nil {
		return nil, appErr
	}

	// target channel does exist
	if args.Input != nil && model.IsValidId(args.Input.ChannelID) {
		transaction, err := embedCtx.App.Srv().Store.GetMaster().Begin()
		if err != nil {
			return nil, model.NewAppError("ChannelDelete", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		// delete checkouts of origin channel
		appErr := deleteCheckoutsByChannelID(args.Id, transaction)
		if appErr != nil {
			return nil, appErr
		}

		// migrate orders to target channel
		for _, order := range orders {
			order.ChannelID = args.Input.ChannelID
		}

		_, appErr = embedCtx.App.Srv().OrderService().BulkUpsertOrders(transaction, orders)
		if appErr != nil {
			return nil, appErr
		}

		err = transaction.Commit()
		if err != nil {
			return nil, model.NewAppError("ChannelDelete", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
		store.FinalizeTransaction(transaction)

	} else {
		if len(orders) > 0 {
			return nil, model.NewAppError("ChannelDelete", "api.channel.delete_channel_with_orders.app_error", nil, "you must specify a target channel to migrate orders of given channel to", http.StatusNotAcceptable)
		}

		appErr := deleteCheckoutsByChannelID(args.Id, nil)
		if appErr != nil {
			return nil, appErr
		}
	}

	// delete channel
	appErr = embedCtx.App.Srv().ChannelService().DeleteChannels(nil, args.Id)
	if appErr != nil {
		return nil, appErr
	}

	return &ChannelDelete{
		Channel: &Channel{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelActivate(ctx context.Context, args struct{ Id string }) (*ChannelActivate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate channel
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	if !channel.IsActive {
		channel.IsActive = true
		channel, appErr = embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ChannelActivate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelDeactivate(ctx context.Context, args struct{ Id string }) (*ChannelDeactivate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// validate channel
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	if channel.IsActive {
		channel.IsActive = false
		channel, appErr = embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ChannelDeactivate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

func (r *Resolver) Channel(ctx context.Context, args struct{ Id string }) (*Channel, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Channel", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is not a valid channel id", args.Id), http.StatusBadRequest)
	}

	channel, err := ChannelByIdLoader.Load(ctx, args.Id)()
	if err != nil {
		return nil, err
	}

	return SystemChannelToGraphqlChannel(channel), nil
}

func (r *Resolver) Channels(ctx context.Context) ([]*Channel, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channels, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{})
	if appErr != nil {
		return nil, appErr
	}

	return systemRecordsToGraphql(channels, SystemChannelToGraphqlChannel), nil
}
