package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func SystemGiftcardEventToGraphqlGiftcardEvent(evt *model.GiftCardEvent) *GiftCardEvent {
	if evt == nil {
		return nil
	}

	res := new(GiftCardEvent)
	res.ID = evt.Id
	if evt.Date != 0 {
		res.Date = &DateTime{util.TimeFromMillis(evt.Date)}
	}
	typ := GiftCardEventsEnum(strings.ToUpper(string(evt.Type)))
	res.Type = &typ

	msg, ok := evt.Parameters["message"]
	if str, strOk := msg.(string); ok && strOk {
		res.Message = &str
	}

	email, ok := evt.Parameters["email"]
	if str, strOk := email.(string); ok && strOk {
		res.Email = &str
	}

	orderID, ok := evt.Parameters["order_id"]
	if str, strOk := orderID.(string); ok && strOk {
		res.OrderID = &str
		res.OrderNumber = &str
	}

	tag, ok := evt.Parameters["tag"]
	if str, strOk := tag.(string); ok && strOk {
		res.Tag = &str
	}

	oldTag, ok := evt.Parameters["old_tag"]
	if str, strOk := oldTag.(string); ok && strOk {
		res.OldTag = &str
	}

	balance, ok := evt.Parameters["balance"]
	if ok && balance != nil {
		balanceMap, ok1 := balance.(map[string]any)
		if ok1 {
			currency, ok2 := balanceMap["currency"]
			if ok2 {
				strCurrency := currency.(string)

				for index, field := range []string{"initial_balance", "old_initial_balance", "current_balance", "old_current_balance"} {
					amount, ok3 := balanceMap[field]

					if ok3 && amount != nil {

						var floatValue float64

						switch t := amount.(type) {
						case int:
							floatValue = float64(t)
						case float64:
							floatValue = t
						case decimal.Decimal:
							floatValue, _ = t.Float64()
						case *decimal.Decimal:
							floatValue, _ = t.Float64()
						}

						res.Balance = new(GiftCardEventBalance)
						switch index {
						case 0:
							res.Balance.InitialBalance = &Money{strCurrency, floatValue}
						case 1:
							res.Balance.OldInitialBalance = &Money{strCurrency, floatValue}
						case 2:
							res.Balance.CurrentBalance = &Money{strCurrency, floatValue}
						case 4:
							res.Balance.OldCurrentBalance = &Money{strCurrency, floatValue}
						}
					}
				}
			}
		}
	}

	for _, key := range []string{"expiry_date", "old_expiry_date"} {
		expiryDate, ok := evt.Parameters[key]

		if ok && expiryDate != nil {
			switch t := expiryDate.(type) {
			case string:
				tim, err := time.Parse("2006-01-02", t)
				if err == nil {
					res.ExpiryDate = &Date{DateTime{tim}}
				}
			case time.Time:
				res.ExpiryDate = &Date{DateTime{util.StartOfDay(t)}}
			case *time.Time:
				res.ExpiryDate = &Date{DateTime{util.StartOfDay(*t)}}
			}
		}
	}

	return res
}

func SystemGiftcardToGraphqlGiftcard(gc *model.GiftCard) *GiftCard {
	res := new(GiftCard)

	gc.PopulateNonDbFields()

	res.ID = gc.Id
	res.IsActive = *gc.IsActive
	res.Tag = gc.Tag
	res.DisplayCode = gc.DisplayCode()

	res.createdByEmail = gc.CreatedByEmail
	res.usedByEmail = gc.UsedByEmail
	res.code = gc.Code
	res.usedByID = gc.UsedByID
	res.createdByID = gc.CreatedByID
	res.productID = gc.ProductID

	if gc.ExpiryDate != nil {
		res.ExpiryDate = &Date{
			DateTime{util.StartOfDay(*gc.ExpiryDate)},
		}
	}
	res.Created = DateTime{util.TimeFromMillis(gc.CreateAt)}
	res.Tag = gc.Tag
	if gc.LastUsedOn != nil {
		res.LastUsedOn = &DateTime{util.TimeFromMillis(*gc.LastUsedOn)}
	}
	if gc.CurrentBalance != nil { // these fields may got populated above
		flt, _ := gc.CurrentBalance.Amount.Float64()
		res.CurrentBalance = &Money{
			Amount:   flt,
			Currency: gc.CurrentBalance.Currency,
		}
	}
	if gc.InitialBalance != nil { // these fields may got populated above
		flt, _ := gc.InitialBalance.Amount.Float64()
		res.InitialBalance = &Money{
			Amount:   flt,
			Currency: gc.InitialBalance.Currency,
		}
	}
	res.Metadata = MetadataToSlice(gc.Metadata)
	res.PrivateMetadata = MetadataToSlice(gc.PrivateMetadata)

	return res
}

func (g *GiftCard) Product(ctx context.Context) (*Product, error) {
	if g.productID == nil {
		return nil, nil
	}

	return dataloaders.productsByIDs.Load(ctx, *g.productID)()
}

func (g *GiftCard) Events(ctx context.Context) ([]*GiftCardEvent, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// check if current user has permission to manage this giftcard
	if embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageGiftcard) {

		// dataloaders.giftcardEventsByGiftcardIDs.Load()()
		panic("not implemented")
	}

	return nil, model.NewAppError("giftcard.Events", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (g *GiftCard) CreatedByEmail(ctx context.Context) (*string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	resolveCreatedByEmail := func(u *User) *string {
		if (u != nil && u.ID == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageGiftcard) {

			if u != nil {
				return &u.Email
			}

			return g.createdByEmail
		}

		var email string
		if u != nil {
			email = u.Email
		} else if g.createdByEmail != nil {
			email = *g.createdByEmail
		}

		return model.NewString(util.ObfuscateEmail(email))
	}

	if g.createdByID == nil {
		return resolveCreatedByEmail(nil), nil
	}

	user, err := dataloaders.usersByIDs.Load(ctx, *g.createdByID)()
	if err != nil {
		return nil, err
	}

	return resolveCreatedByEmail(user), nil
}

func (g *GiftCard) UsedByEmail(ctx context.Context) (*string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	resolveUsedByEmail := func(u *User) *string {
		if (u != nil && u.ID == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageGiftcard) {

			if u != nil {
				return &u.Email
			}

			return g.usedByEmail
		}

		var email string
		if u != nil {
			email = u.Email
		} else if g.usedByEmail != nil {
			email = *g.usedByEmail
		}

		return model.NewString(util.ObfuscateEmail(email))
	}

	if g.usedByID == nil {
		return resolveUsedByEmail(nil), nil
	}

	user, err := dataloaders.usersByIDs.Load(ctx, *g.usedByID)()
	if err != nil {
		return nil, err
	}

	return resolveUsedByEmail(user), nil
}

func (g *GiftCard) UsedBy(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	resolveUsedBy := func(u *User) (*User, *model.AppError) {
		if (u != nil && u.ID == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().AccountService().
				SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageUsers) {
			if u != nil {
				return u, nil
			}
			return nil, nil
		}

		return nil, model.NewAppError("GiftCard.UsedBy", ErrorUnauthorized, nil, "You are not authorized to perform this", http.StatusUnauthorized)
	}

	if g.usedByID == nil {
		return resolveUsedBy(nil)
	}

	user, err := dataloaders.usersByIDs.Load(ctx, *g.usedByID)()
	if err != nil {
		return nil, err
	}

	return resolveUsedBy(user)
}

func (g *GiftCard) CreatedBy(ctx context.Context) (*User, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	resolveCreatedBy := func(u *User) (*User, error) {
		if (u != nil && u.ID == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageUsers) {

			if u != nil {
				return u, nil
			}

			user, appErr := embedCtx.App.Srv().
				AccountService().UserById(ctx, embedCtx.AppContext.Session().UserId)
			if appErr != nil {
				return nil, appErr
			}

			return SystemUserToGraphqlUser(user), nil
		}

		return nil, model.NewAppError("GiftCard.CreatedBy", ErrorUnauthorized, nil, "you are not authorized to perform this", http.StatusUnauthorized)
	}

	if g.createdByID == nil {
		return resolveCreatedBy(nil)
	}

	user, err := dataloaders.usersByIDs.Load(ctx, *g.createdByID)()
	if err != nil {
		return nil, err
	}

	return resolveCreatedBy(user)
}

func (g *GiftCard) Code(ctx context.Context) (string, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return "", err
	}

	resolveCode := func(u *User) (string, *model.AppError) {
		if (g.usedByEmail == nil && embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageGiftcard)) ||
			(u != nil && u.ID == embedCtx.AppContext.Session().UserId) {
			return g.code, nil
		}

		return "", model.NewAppError("GiftCard.Code", ErrorUnauthorized, nil, "You are not authorized to perform this action", http.StatusUnauthorized)
	}

	if g.usedByID == nil {
		return resolveCode(nil)
	}

	user, err := dataloaders.usersByIDs.Load(ctx, *g.usedByID)()
	if err != nil {
		return "", err
	}

	return resolveCode(user)
}

func (g *GiftCard) BoughtInChannel(ctx context.Context) {
	panic("not implemented")
}

func graphqlGiftcardsByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[*GiftCard] {
	var (
		res       []*dataloader.Result[*GiftCard]
		appErr    *model.AppError
		giftcards []*model.GiftCard
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	giftcards, appErr = embedCtx.App.Srv().
		GiftcardService().
		GiftcardsByOption(nil, &model.GiftCardFilterOption{
			UsedByID: squirrel.Eq{store.GiftcardTableName + ".UsedByID": userIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, gc := range giftcards {
		res = append(res, &dataloader.Result[*GiftCard]{Data: SystemGiftcardToGraphqlGiftcard(gc)})
	}
	return res

errorLabel:
	for range userIDs {
		res = append(res, &dataloader.Result[*GiftCard]{Error: err})
	}
	return res
}

func graphqlGiftcardEventsByGiftcardIDsLoader(ctx context.Context, giftcardIDs []string) []*dataloader.Result[*GiftCardEvent] {
	var (
		res    []*dataloader.Result[*GiftCardEvent]
		appErr *model.AppError
		events []*model.GiftCardEvent
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	events, appErr = embedCtx.App.Srv().
		GiftcardService().
		GiftcardEventsByOptions(&model.GiftCardEventFilterOption{
			GiftcardID: squirrel.Eq{store.GiftcardTableName + ".GiftcardID": giftcardIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, evt := range events {
		res = append(res, &dataloader.Result[*GiftCardEvent]{Data: SystemGiftcardEventToGraphqlGiftcardEvent(evt)})
	}
	return res

errorLabel:
	for range giftcardIDs {
		res = append(res, &dataloader.Result[*GiftCardEvent]{Error: err})
	}
	return res
}
