package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

func (r *Resolver) GiftCardActivate(ctx context.Context, args struct{ Id string }) (*GiftCardActivate, error) {
	// only staffs of shop which issued the giftcard can activate them
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateGiftcard)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftCardActivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	giftcard, appErr := r.srv.GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if giftcard.IsActive != nil && !*giftcard.IsActive {
		giftcard.IsActive = model.NewPrimitive(true)
		giftcards, appErr := r.srv.GiftcardService().UpsertGiftcards(nil, giftcard)
		if appErr != nil {
			return nil, appErr
		}
		giftcard = giftcards[0]

		// create giftcard event
		_, appErr = r.srv.GiftcardService().BulkUpsertGiftcardEvents(nil, &model.GiftCardEvent{
			GiftcardID: giftcard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.ACTIVATED,
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

func (r *Resolver) GiftCardDeactivate(ctx context.Context, args struct{ Id string }) (*GiftCardDeactivate, error) {
	// only staffs of shop which issued the giftcard can deactivate them
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateGiftcard)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftcardDeactivate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	giftcard, appErr := r.srv.GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	if giftcard.IsActive != nil && *giftcard.IsActive {
		giftcard.IsActive = model.NewPrimitive(false)
		giftcards, appErr := r.srv.GiftcardService().UpsertGiftcards(nil, giftcard)
		if appErr != nil {
			return nil, appErr
		}
		giftcard = giftcards[0]

		// giftcard event
		_, appErr = r.srv.GiftcardService().BulkUpsertGiftcardEvents(nil, &model.GiftCardEvent{
			GiftcardID: giftcard.Id,
			UserID:     &embedCtx.AppContext.Session().UserId,
			Type:       model.DEACTIVATED,
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	return &GiftCardDeactivate{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcard),
	}, nil
}

func (r *Resolver) GiftCardUpdate(ctx context.Context, args struct {
	Id    string
	Input GiftCardUpdateInput
}) (*GiftCardUpdate, error) {
	// requester must be staff at the shop which issued given giftcard to update it
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasPermissionToAll(model.PermissionUpdateGiftcard)
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	// valudate input
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftcardUpdate", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("%s is invalid id", args.Id), http.StatusBadRequest)
	}

	// check if giftcard relly belong to current shop
	giftcard, appErr := r.srv.GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

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
	giftcards, appErr := r.srv.GiftcardService().UpsertGiftcards(nil, giftcard)
	if appErr != nil {
		return nil, appErr
	}

	// TODO: check if we need create some giftcard event

	return &GiftCardUpdate{
		GiftCard: SystemGiftcardToGraphqlGiftcard(giftcards[0]),
	}, nil
}

func (r *Resolver) GiftCardResend(ctx context.Context, args struct{ Input GiftCardResendInput }) (*GiftCardResend, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardAddNote(ctx context.Context, args struct {
	Id    string
	Input GiftCardAddNoteInput
}) (*GiftCardAddNote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDelete(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkActivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkActivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardBulkDeactivate(ctx context.Context, args struct{ Ids []string }) (*GiftCardBulkDeactivate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCard(ctx context.Context, args struct{ Id string }) (*GiftCard, error) {
	// requester must be authenticated to see giftcard
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("GiftCard", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, fmt.Sprintf("$s is invalid id", args.Id), http.StatusBadRequest)
	}

	giftcard, appErr := r.srv.GiftcardService().GetGiftCard(args.Id)
	if appErr != nil {
		return nil, appErr
	}

	return SystemGiftcardToGraphqlGiftcard(giftcard), nil
}

func (r *Resolver) GiftCardSettings(ctx context.Context) (*GiftCardSettings, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCards(ctx context.Context, args struct {
	SortBy *GiftCardSortingInput
	Filter *GiftCardFilterInput
	GraphqlParams
}) (*GiftCardCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) GiftCardCurrencies(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}
