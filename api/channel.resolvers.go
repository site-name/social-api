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
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
	"gorm.io/gorm"
)

// NOTE: Refer to ./schemas/channel.graphqls for directive used
func (r *Resolver) ChannelCreate(ctx context.Context, args struct{ Input ChannelCreateInput }) (*ChannelCreate, error) {
	// validate input
	if !lo.EveryBy(args.Input.AddShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "addShippingZones"}, "please provide valid addShippingZones", http.StatusBadRequest)
	}
	if !args.Input.DefaultCountry.IsValid() {
		return nil, model.NewAppError("ChannelCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "defaultCountry"}, fmt.Sprintf("%s is not valid country code", args.Input.DefaultCountry), http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	channel := &model.Channel{
		Name:     args.Input.Name,
		Currency: args.Input.CurrencyCode,
	}
	if val := args.Input.IsActive; val != nil {
		channel.IsActive = *val
	}
	if val := args.Input.Slug; slug.IsSlug(val) {
		channel.Slug = val
	}

	// save new channel to db
	channel, appErr := embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
	if appErr != nil {
		return nil, appErr
	}

	// save channel shipping zone relations:
	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("ChannelCreate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	shippingZones := lo.Map(args.Input.AddShippingZones, func(id string, _ int) *model.ShippingZone {
		return &model.ShippingZone{Id: id}
	})
	appErr = embedCtx.App.Srv().ShippingService().ToggleShippingZoneRelations(transaction, shippingZones, nil, []string{channel.Id}, false)
	if appErr != nil {
		return nil, appErr
	}

	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ChannelCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
	// validate inputs
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.AddShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "AddShippingZones"}, "please provide valid ids", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Input.RemoveShippingZones, model.IsValidId) {
		return nil, model.NewAppError("ChannelUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "RemoveShippingZones"}, "please provide valid ids", http.StatusBadRequest)
	}
	intersectIds := lo.Intersect(args.Input.RemoveShippingZones, args.Input.AddShippingZones)
	if len(intersectIds) > 0 {
		return nil, model.NewAppError("ChannelUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "RemoveShippingZones/AddShippingZones"}, "remove shipping zone ids and add shipping zone ids can not have same ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// validate if channe does exist
	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
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
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("ChannelUpdate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer (transaction).Rollback()

	// create relations between shipping zones and channel
	if len(args.Input.AddShippingZones) > 0 {
		shippingZonesToAdd := lo.Map(args.Input.AddShippingZones, func(id string, _ int) *model.ShippingZone { return &model.ShippingZone{Id: id} })
		appErr = embedCtx.App.Srv().ShippingService().ToggleShippingZoneRelations(transaction, shippingZonesToAdd, []string{}, []string{channel.Id}, false)
		if appErr != nil {
			return nil, appErr
		}
	}

	if len(args.Input.RemoveShippingZones) > 0 {
		// delete relations between shipping zones and channel
		shippingZonesToRemove := lo.Map(args.Input.RemoveShippingZones, func(id string, _ int) *model.ShippingZone { return &model.ShippingZone{Id: id} })
		appErr = embedCtx.App.Srv().ShippingService().ToggleShippingZoneRelations(transaction, shippingZonesToRemove, []string{}, []string{channel.Id}, true)
		if appErr != nil {
			return nil, appErr
		}

		// delete shipping methods channel listings of shipping methods of deleted shipping zones
		shippingMethodChannelListings, appErr := embedCtx.App.Srv().
			ShippingService().
			ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
				Conditions:                          squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ChannelID": channel.Id},
				ShippingMethod_ShippingZoneID_Inner: squirrel.Eq{model.ShippingZoneTableName + ".Id": args.Input.RemoveShippingZones},
			})
		if appErr != nil {
			return nil, appErr
		}

		appErr = embedCtx.App.Srv().ShippingService().DeleteShippingMethodChannelListings(transaction, &model.ShippingMethodChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ShippingMethodChannelListingTableName + ".Id": shippingMethodChannelListings.IDs()},
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
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("ChannelUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("ChannelDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide a valid channel id", http.StatusBadRequest)
	}
	if args.Input != nil && !model.IsValidId(args.Input.ChannelID) {
		return nil, model.NewAppError("ChannelDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channelID"}, "please provide a valid channel id", http.StatusBadRequest)
	}
	if args.Input != nil && args.Input.ChannelID == args.Id {
		return nil, model.NewAppError("ChannelDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channelID"}, "target channel cannot be the channel to be deleted", http.StatusBadRequest)
	}

	deleteCheckoutsByChannelID := func(channelID string, transaction *gorm.DB) *model.AppError {
		return embedCtx.App.Srv().
			CheckoutService().
			DeleteCheckoutsByOption(transaction, &model.CheckoutFilterOption{
				Conditions: squirrel.Eq{model.CheckoutTableName + ".ChannelID": args.Id},
			})
	}

	orders, appErr := embedCtx.App.Srv().
		OrderService().
		FilterOrdersByOptions(&model.OrderFilterOption{
			Conditions:      squirrel.Eq{model.OrderTableName + ".ChannelID": args.Id},
			SelectForUpdate: true,
		})
	if appErr != nil {
		return nil, appErr
	}

	// target channel does exist
	if args.Input != nil && model.IsValidId(args.Input.ChannelID) {
		transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
		if transaction.Error != nil {
			return nil, model.NewAppError("ChannelDelete", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
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

		err := transaction.Commit().Error
		if err != nil {
			return nil, model.NewAppError("ChannelDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
		(transaction).Rollback()

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
		return nil, model.NewAppError("ChannelActivate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
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
		return nil, model.NewAppError("ChannelActivate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := embedCtx.App.Srv().ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": args.Id},
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
		return nil, model.NewAppError("Channel", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is not a valid channel id", args.Id), http.StatusBadRequest)
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
