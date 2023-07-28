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
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardActivate(ctx context.Context, args struct{ Id string }) (*GiftCardActivate, error) {

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftCardActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if giftcard.IsActive != nil && !*giftcard.IsActive {
		giftcard.IsActive = model.NewPrimitive(true)
		giftcards, appErr := embedCtx.App.Srv().GiftcardService().UpsertGiftcards(nil, giftcard)
		if appErr != nil {
			return nil, appErr
		}
		giftcard = giftcards[0]

		// create giftcard event
		_, appErr = embedCtx.App.Srv().GiftcardService().BulkUpsertGiftcardEvents(nil, &model.GiftCardEvent{
			GiftcardID: giftcard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.GIFT_CARD_EVENT_TYPE_ACTIVATED,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	return &GiftCardActivate{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcard),
	}, nil
}

func (r *Resolver) GiftCardCreate(ctx context.Context, args struct{ Input GiftCardCreateInput }) (*GiftCardCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardDelete(ctx context.Context, args struct{ Id string }) (*GiftCardDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardDeactivate(ctx context.Context, args struct{ Id string }) (*GiftCardDeactivate, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftcardDeactivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if giftcard.IsActive != nil && *giftcard.IsActive {
		giftcard.IsActive = model.NewPrimitive(false)
		giftcards, appErr := embedCtx.App.Srv().GiftcardService().UpsertGiftcards(nil, giftcard)
		if appErr != nil {
			return nil, appErr
		}
		giftcard = giftcards[0]

		// giftcard event
		_, appErr = embedCtx.App.Srv().GiftcardService().BulkUpsertGiftcardEvents(nil, &model.GiftCardEvent{
			GiftcardID: giftcard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.GIFT_CARD_EVENT_TYPE_DEACTIVATED,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	return &GiftCardDeactivate{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcard),
	}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardUpdate(ctx context.Context, args struct {
	Id    string
	Input GiftCardUpdateInput
}) (*GiftCardUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// valudate input
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftcardUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	// check if giftcard relly belong to current shop
	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	oldGiftCard := giftcard.DeepCopy()

	if v := args.Input.Tag; v != nil {
		giftcard.Tag = v
	}
	if v := args.Input.ExpiryDate; v != nil {
		giftcard.ExpiryDate = &v.Time
	}
	if v := args.Input.BalanceAmount; v != nil {
		precision, _ := goprices.GetCurrencyPrecision(giftcard.Currency)
		rounded := decimal.Decimal(*v).Round(int32(precision))
		giftcard.CurrentBalanceAmount = &rounded
		giftcard.InitialBalanceAmount = &rounded
	}

	// update
	giftcards, appErr := embedCtx.App.Srv().GiftcardService().UpsertGiftcards(nil, giftcard)
	if appErr != nil {
		return nil, appErr
	}

	updatedGiftCard := giftcards[0]

	// add gift card events if needed
	eventsToAdd := []*model.GiftCardEvent{}
	if args.Input.BalanceAmount != nil {
		// embedCtx.App.Srv().GiftcardService().GiftCardEve
		eventsToAdd = append(eventsToAdd, &model.GiftCardEvent{
			GiftcardID: updatedGiftCard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.GIFT_CARD_EVENT_TYPE_BALANCE_RESET,
			Parameters: model.StringInterface{
				"currency":            updatedGiftCard.Currency,
				"initial_balance":     updatedGiftCard.InitialBalanceAmount,
				"current_balance":     updatedGiftCard.CurrentBalanceAmount,
				"old_currency":        updatedGiftCard.Currency,
				"old_initial_balance": oldGiftCard.InitialBalanceAmount,
				"old_current_balance": oldGiftCard.CurrentBalanceAmount,
			},
		})
	}
	if args.Input.ExpiryDate != nil {
		eventsToAdd = append(eventsToAdd, &model.GiftCardEvent{
			GiftcardID: updatedGiftCard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.GIFT_CARD_EVENT_TYPE_EXPIRY_DATE_UPDATED,
			Parameters: model.StringInterface{
				"expiry_date":     updatedGiftCard.ExpiryDate,
				"old_expiry_date": oldGiftCard.ExpiryDate,
			},
		})
	}
	if args.Input.Tag != nil {
		eventsToAdd = append(eventsToAdd, &model.GiftCardEvent{
			GiftcardID: updatedGiftCard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.GIFT_CARD_EVENT_TYPE_TAG_UPDATED,
			Parameters: model.StringInterface{
				"tag":     updatedGiftCard.Tag,
				"old_tag": oldGiftCard.Tag,
			},
		})
	}

	_, appErr = embedCtx.App.Srv().GiftcardService().BulkUpsertGiftcardEvents(nil, eventsToAdd...)
	if appErr != nil {
		return nil, appErr
	}

	return &GiftCardUpdate{
		GiftCard: SystemGiftcardToGraphqlGiftcard(updatedGiftCard),
	}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardResend(ctx context.Context, args struct{ Input GiftCardResendInput }) (*GiftCardResend, error) {
	// validate params
	if !model.IsValidId(args.Input.ID) {
		return nil, model.NewAppError("GiftcardResend", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, args.Input.ID+" is not a valid id", http.StatusBadRequest)
	}
	if args.Input.Email != nil && !model.IsValidEmail(*args.Input.Email) {
		return nil, model.NewAppError("GiftcardResend", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "email"}, *args.Input.Email+" is not a valid email", http.StatusBadRequest)
	}
	if !model.IsValidId(args.Input.Channel) {
		return nil, model.NewAppError("GiftcardResend", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel"}, args.Input.Channel+" is not a valid channel id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Input.ID)
	if appErr != nil {
		return nil, appErr
	}

	var targetEmail string
	switch {
	case args.Input.Email != nil:
		targetEmail = *args.Input.Email
	case giftcard.UsedByEmail != nil:
		targetEmail = *giftcard.UsedByEmail
	case giftcard.CreatedByEmail != nil:
		targetEmail = *giftcard.CreatedByEmail
	}

	receiver, appErr := embedCtx.App.Srv().AccountService().GetUserByOptions(ctx, &model.UserFilterOptions{
		Conditions: squirrel.Eq{model.UserTableName + ".Email": targetEmail},
	})
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()

	appErr = embedCtx.App.Srv().GiftcardService().SendGiftcardNotification(
		&model.User{Id: embedCtx.AppContext.Session().UserId},
		nil,
		receiver,
		targetEmail,
		*giftcard,
		pluginMng,
		args.Input.Channel,
		true,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &GiftCardResend{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcard),
	}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardAddNote(ctx context.Context, args struct {
	Id    string
	Input GiftCardAddNoteInput
}) (*GiftCardAddNote, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftcardAddNote", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, args.Id+" is not a valid id", http.StatusBadRequest)
	}
	if args.Input.Message == "" {
		return nil, model.NewAppError("GiftcardAddNote", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "message"}, "message can not be empty", http.StatusBadRequest)
	}

	// validate if giftcard really does exist
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	giftcardEvent := &model.GiftCardEvent{
		GiftcardID: args.Id,
		UserID:     &embedCtx.AppContext.Session().UserId,
		Type:       model.GIFT_CARD_EVENT_TYPE_NOTE_ADDED,
		Parameters: model.StringInterface{
			"message": args.Input.Message,
		},
	}

	events, appErr := embedCtx.App.Srv().GiftcardService().BulkUpsertGiftcardEvents(nil, giftcardEvent)
	if appErr != nil {
		return nil, appErr
	}

	return &GiftCardAddNote{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcard),
		Event:    SystemGiftcardEventToGraphqlGiftcardEvent(events[0]),
	}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardBulkDelete(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("GiftCardBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid gift card ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	appErr := embedCtx.App.Srv().GiftcardService().DeleteGiftcards(nil, args.Ids)
	if appErr != nil {
		return nil, appErr
	}

	return &GiftCardBulkDelete{Count: int32(len(args.Ids))}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardBulkActivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkActivate, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("GiftCardBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid gift card ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcards, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.Eq{
			model.GiftcardTableName + ".IsActive": false,
			model.GiftcardTableName + ".Id":       args.Ids,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("GiftCardBulkActivate", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// update
	for _, gc := range giftcards {
		gc.IsActive = model.NewPrimitive(true)
	}
	_, appErr = embedCtx.App.Srv().GiftcardService().UpsertGiftcards(transaction, giftcards...)
	if appErr != nil {
		return nil, appErr
	}

	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("GiftCardBulkActivate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &GiftCardBulkActivate{
		Count: int32(len(args.Ids)),
	}, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardBulkDeactivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDeactivate, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("GiftCardBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid gift card ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcards, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{
		Conditions: squirrel.Eq{
			model.GiftcardTableName + ".IsActive": true,
			model.GiftcardTableName + ".Id":       args.Ids,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("GiftCardBulkActivate", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// update
	for _, gc := range giftcards {
		gc.IsActive = model.NewPrimitive(false)
	}
	_, appErr = embedCtx.App.Srv().GiftcardService().UpsertGiftcards(transaction, giftcards...)
	if appErr != nil {
		return nil, appErr
	}

	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("GiftCardBulkActivate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &GiftCardBulkDeactivate{
		Count: int32(len(args.Ids)),
	}, nil
}

func (r *Resolver) GiftCard(ctx context.Context, args struct{ Id string }) (*GiftCard, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftCard", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("$s is invalid id", args.Id), http.StatusBadRequest)
	}

	giftcard, appErr := embedCtx.App.Srv().GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	return SystemGiftcardToGraphqlGiftcard(giftcard), nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCards(ctx context.Context, args struct {
	SortBy *GiftCardSortingInput
	Filter *GiftCardFilterInput
	GraphqlParams
}) (*GiftCardCountableConnection, error) {
	// validate params
	if args.Filter != nil {
		if err := args.Filter.validate(); err != nil {
			return nil, err
		}
	}
	err := args.GraphqlParams.Validate("GiftCards")
	if err != nil {
		return nil, err
	}

	giftcardFilter := args.Filter.ToSystemGiftcardFilter()
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	giftcards, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(giftcardFilter)
	if appErr != nil {
		return nil, appErr
	}

	var connection *CountableConnection[*GiftCard]

	switch args.SortBy.Field {
	case GiftCardSortFieldTag:
		keyFunc := func(g *model.GiftCard) string {
			if g.Tag != nil {
				return *g.Tag
			}
			return ""
		}
		connection, appErr = newGraphqlPaginator(giftcards, keyFunc, SystemGiftcardToGraphqlGiftcard, args.GraphqlParams).parse("GiftCards")

	case GiftCardSortFieldCurrentBalance:
		keyFunc := func(g *model.GiftCard) decimal.Decimal {
			if g.CurrentBalanceAmount == nil {
				g.CurrentBalanceAmount = &decimal.Zero
			}
			return *g.CurrentBalanceAmount
		}
		connection, appErr = newGraphqlPaginator(giftcards, keyFunc, SystemGiftcardToGraphqlGiftcard, args.GraphqlParams).parse("GiftCards")
	}
	if appErr != nil {
		return nil, appErr
	}

	return (*GiftCardCountableConnection)(unsafe.Pointer(connection)), nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (r *Resolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	giftcards, appErr := embedCtx.App.Srv().GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{})
	if appErr != nil {
		return nil, appErr
	}

	return lo.Map(giftcards, func(gc *model.GiftCard, _ int) string { return gc.Currency }), nil
}
