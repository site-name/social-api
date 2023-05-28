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
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) ChannelCreate(ctx context.Context, args struct{ Input ChannelCreateInput }) (*ChannelCreate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionCreateChannel)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate ids
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

	// insert
	channel, appErr := embedCtx.App.Srv().ChannelService().UpsertChannel(nil, channel)
	if appErr != nil {
		return nil, appErr
	}

	// save channel shipping zone relations:
	// begin transaction
	transaction, err := r.srv.Store.GetMasterX().Beginx()
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
	_, appErr = r.srv.ChannelService().BulkUpsertShippingZoneChannels(transaction, shippingZoneChannelRelations)
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

func (r *Resolver) ChannelUpdate(ctx context.Context, args struct {
	Id    string
	Input ChannelUpdateInput
}) (*ChannelUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateChannel)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

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
	channel, appErr := r.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": args.Id},
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
	channel, appErr = r.srv.ChannelService().UpsertChannel(nil, channel)
	if appErr != nil {
		return nil, appErr
	}

	// update m2m relations
	// begin transaction
	transaction, err := r.srv.Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("ChannelUpdate", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	// NOTE: Don't defer transaction.RollBack() here

	// create relations between shipping zones and channel
	insertRelations := lo.Map(args.Input.AddShippingZones, func(id string, _ int) *model.ShippingZoneChannel {
		return &model.ShippingZoneChannel{ShippingZoneID: id, ChannelID: channel.Id}
	})
	_, appErr = r.srv.ChannelService().BulkUpsertShippingZoneChannels(transaction, insertRelations)
	if appErr != nil {
		return nil, appErr
	}

	if len(args.Input.RemoveShippingZones) > 0 {
		// delete relations between shipping zones and channel
		deleteRelations := lo.Map(args.Input.RemoveShippingZones, func(id string, _ int) *model.ShippingZoneChannel {
			return &model.ShippingZoneChannel{ShippingZoneID: id, ChannelID: channel.Id}
		})
		appErr = r.srv.ChannelService().BulkDeleteShippingZoneChannels(transaction, deleteRelations)
		if appErr != nil {
			return nil, appErr
		}

		// delete shipping methods channel listings of shipping methods of deleted shipping zones
		shippingMethodChannelListings, appErr := r.srv.ShippingService().
			ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
				ChannelID:                           squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": channel.Id},
				ShippingMethod_ShippingZoneID_Inner: squirrel.Eq{store.ShippingZoneTableName + ".Id": args.Input.RemoveShippingZones},
			})
		if appErr != nil {
			return nil, appErr
		}

		err = r.srv.Store.ShippingMethodChannelListing().BulkDelete(transaction, shippingMethodChannelListings.IDs())
		if err != nil {
			return nil, model.NewAppError("ChannelUpdate", "app.shippng.delete_shipping_method_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		r.srv.Go(func() {
			appErr := r.srv.ShippingService().DropInvalidShippingMethodsRelationsForGivenChannels(transaction, shippingMethodChannelListings.ShippingMethodIDs(), []string{channel.Id})
			if appErr != nil {
				slog.Error("failed to update invalid shipping method", slog.Err(appErr))
			}

			err := transaction.Commit()
			if err != nil {
				slog.Error("failed to commit transaction", slog.Err(err))
			}

			store.FinalizeTransaction(transaction)
		})

	} else {
		err := transaction.Commit()
		if err != nil {
			return nil, model.NewAppError("ChannelUpdate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		store.FinalizeTransaction(transaction)
	}

	return &ChannelUpdate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

func (r *Resolver) ChannelDelete(ctx context.Context, args struct {
	Id    string
	Input *ChannelDeleteInput
}) (*ChannelDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ChannelActivate(ctx context.Context, args struct{ Id string }) (*ChannelActivate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateChannel)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate channel
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := r.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	if !channel.IsActive {
		channel.IsActive = true
		channel, appErr = r.srv.ChannelService().UpsertChannel(nil, channel)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ChannelActivate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

func (r *Resolver) ChannelDeactivate(ctx context.Context, args struct{ Id string }) (*ChannelDeactivate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateChannel)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// validate channel
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("ChannelActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid channel id", http.StatusBadRequest)
	}

	channel, appErr := r.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	if channel.IsActive {
		channel.IsActive = false
		channel, appErr = r.srv.ChannelService().UpsertChannel(nil, channel)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ChannelDeactivate{
		Channel: SystemChannelToGraphqlChannel(channel),
	}, nil
}

func (r *Resolver) Channel(ctx context.Context, args struct{ Id *string }) (*Channel, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Channels(ctx context.Context) ([]*Channel, error) {
	panic(fmt.Errorf("not implemented"))
}
