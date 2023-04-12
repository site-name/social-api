package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type GiftCardEvent struct {
	ID            string                `json:"id"`
	Date          *DateTime             `json:"date"`
	Type          *GiftCardEventsEnum   `json:"type"`
	Message       *string               `json:"message"`
	Email         *string               `json:"email"`
	OrderID       *string               `json:"orderId"`
	OrderNumber   *string               `json:"orderNumber"`
	Tag           *string               `json:"tag"`
	OldTag        *string               `json:"oldTag"`
	Balance       *GiftCardEventBalance `json:"balance"`
	ExpiryDate    *Date                 `json:"expiryDate"`
	OldExpiryDate *Date                 `json:"oldExpiryDate"`

	e *model.GiftCardEvent

	// User          *User                 `json:"user"`
	// App           *App                  `json:"app"`
}

func SystemGiftcardEventToGraphqlGiftcardEvent(evt *model.GiftCardEvent) *GiftCardEvent {
	if evt == nil {
		return nil
	}

	res := new(GiftCardEvent)
	res.ID = evt.Id
	res.e = evt
	if evt.Date != 0 {
		res.Date = &DateTime{util.TimeFromMillis(evt.Date)}
	}
	res.Type = &evt.Type

	msg, ok := evt.Parameters["message"]
	if ok && msg != nil {
		res.Message = model.NewPrimitive(msg.(string))
	}

	email, ok := evt.Parameters["email"]
	if ok && email != nil {
		res.Email = model.NewPrimitive(email.(string))
	}

	orderID, ok := evt.Parameters["order_id"]
	if ok && orderID != nil {
		res.OrderID = model.NewPrimitive(orderID.(string))
	}

	tag, ok := evt.Parameters["tag"]
	if ok && tag != nil {
		res.Tag = model.NewPrimitive(tag.(string))
	}

	oldTag, ok := evt.Parameters["old_tag"]
	if ok && oldTag != nil {
		res.OldTag = model.NewPrimitive(oldTag.(string))
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

func (e *GiftCardEvent) User(ctx context.Context) (*User, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	currentSession := embedCtx.AppContext.Session()

	if e.e.UserID != nil && currentSession.UserId == *e.e.UserID ||
		embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(currentSession, model.PermissionManageUsers, model.PermissionManageStaff) {

		if e.e.UserID == nil {
			return nil, nil
		}
		user, err := UserByUserIdLoader.Load(ctx, *e.e.UserID)()
		if err != nil {
			return nil, err
		}
		return SystemUserToGraphqlUser(user), nil
	}

	return nil, model.NewAppError("GiftCardEvent.User", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (e *GiftCardEvent) App(ctx context.Context) (*App, error) {
	panic("not implemented")
}

type GiftCard struct {
	IsActive        bool            `json:"isActive"`
	ExpiryDate      *Date           `json:"expiryDate"`
	Tag             *string         `json:"tag"`
	Created         DateTime        `json:"created"`
	LastUsedOn      *DateTime       `json:"lastUsedOn"`
	InitialBalance  *Money          `json:"initialBalance"`
	CurrentBalance  *Money          `json:"currentBalance"`
	ID              string          `json:"id"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`
	DisplayCode     string          `json:"displayCode"`

	giftcard *model.GiftCard

	// Code            string           `json:"code"`
	// CreatedByEmail  *string          `json:"createdByEmail"`
	// UsedByEmail     *string          `json:"usedByEmail"`
	// App             *App             `json:"app"`
	// Product         *Product         `json:"product"`
	// Events          []*GiftCardEvent `json:"events"`
	// BoughtInChannel *string          `json:"boughtInChannel"`
	// CreatedBy       *User            `json:"createdBy"`
	// UsedBy          *User            `json:"usedBy"`
}

func SystemGiftcardToGraphqlGiftcard(gc *model.GiftCard) *GiftCard {
	if gc == nil {
		return nil
	}

	res := new(GiftCard)

	gc.PopulateNonDbFields()

	res.ID = gc.Id
	res.IsActive = *gc.IsActive
	res.Tag = gc.Tag
	res.DisplayCode = gc.DisplayCode()
	res.giftcard = gc

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
		res.CurrentBalance = &Money{
			Amount:   gc.CurrentBalance.Amount.InexactFloat64(),
			Currency: gc.CurrentBalance.Currency,
		}
	}
	if gc.InitialBalance != nil { // these fields may got populated above
		res.InitialBalance = &Money{
			Amount:   gc.InitialBalance.Amount.InexactFloat64(),
			Currency: gc.InitialBalance.Currency,
		}
	}
	res.Metadata = MetadataToSlice(gc.Metadata)
	res.PrivateMetadata = MetadataToSlice(gc.PrivateMetadata)

	return res
}

func (g *GiftCard) App(ctx context.Context) (*App, error) {
	panic("not implemented")
}

func (g *GiftCard) Product(ctx context.Context) (*Product, error) {
	if g.giftcard.ProductID == nil {
		return nil, nil
	}

	product, err := ProductByIdLoader.Load(ctx, *g.giftcard.ProductID)()
	if err != nil {
		return nil, err
	}

	return SystemProductToGraphqlProduct(product), nil
}

func (g *GiftCard) Events(ctx context.Context) ([]*GiftCardEvent, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	// check if current user has permission to manage this giftcard
	if embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageGiftcard) {

		events, err := GiftCardEventsByGiftCardIdLoader.Load(ctx, g.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(events, SystemGiftcardEventToGraphqlGiftcardEvent), nil
	}

	return nil, model.NewAppError("giftcard.Events", ErrorUnauthorized, nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

func (g *GiftCard) CreatedByEmail(ctx context.Context) (*string, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	resolveCreatedByEmail := func(u *model.User) *string {
		if (u != nil && u.Id == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageGiftcard) {

			if u != nil {
				return &u.Email
			}

			return g.giftcard.CreatedByEmail
		}

		var email string
		if u != nil {
			email = u.Email
		} else if g.giftcard.CreatedByEmail != nil {
			email = *g.giftcard.CreatedByEmail
		}

		return model.NewPrimitive(util.ObfuscateEmail(email))
	}

	if g.giftcard.CreatedByID == nil {
		return resolveCreatedByEmail(nil), nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.CreatedByID)()
	if err != nil {
		return nil, err
	}

	return resolveCreatedByEmail(user), nil
}

func (g *GiftCard) UsedByEmail(ctx context.Context) (*string, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	resolveUsedByEmail := func(u *model.User) *string {
		if (u != nil && u.Id == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageGiftcard) {

			if u != nil {
				return &u.Email
			}

			return g.giftcard.UsedByEmail
		}

		var email string
		if u != nil {
			email = u.Email
		} else if g.giftcard.UsedByEmail != nil {
			email = *g.giftcard.UsedByEmail
		}

		return model.NewPrimitive(util.ObfuscateEmail(email))
	}

	if g.giftcard.UsedByID == nil {
		return resolveUsedByEmail(nil), nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.UsedByID)()
	if err != nil {
		return nil, err
	}

	return resolveUsedByEmail(user), nil
}

func (g *GiftCard) UsedBy(ctx context.Context) (*User, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	resolveUsedBy := func(u *model.User) (*User, *model.AppError) {
		if (u != nil && u.Id == embedCtx.AppContext.Session().UserId) ||
			embedCtx.
				App.
				Srv().
				AccountService().
				SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageUsers) {
			if u != nil {
				return SystemUserToGraphqlUser(u), nil
			}
			return nil, nil
		}

		return nil, model.NewAppError("GiftCard.UsedBy", ErrorUnauthorized, nil, "You are not authorized to perform this", http.StatusUnauthorized)
	}

	if g.giftcard.UsedByID == nil {
		return resolveUsedBy(nil)
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.UsedByID)()
	if err != nil {
		return nil, err
	}

	return resolveUsedBy(user)
}

func (g *GiftCard) CreatedBy(ctx context.Context) (*User, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	resolveCreatedBy := func(u *model.User) (*User, error) {
		if (u != nil && u.Id == embedCtx.AppContext.Session().UserId) ||
			embedCtx.App.Srv().
				AccountService().
				HasPermissionTo(embedCtx.AppContext.Session().UserId, model.PermissionManageUsers) {

			if u != nil {
				return SystemUserToGraphqlUser(u), nil
			}

			user, appErr := embedCtx.App.Srv().
				AccountService().
				UserById(ctx, embedCtx.AppContext.Session().UserId)
			if appErr != nil {
				return nil, appErr
			}

			return SystemUserToGraphqlUser(user), nil
		}

		return nil, model.NewAppError("GiftCard.CreatedBy", ErrorUnauthorized, nil, "you are not authorized to perform this", http.StatusUnauthorized)
	}

	if g.giftcard.CreatedByID == nil {
		return resolveCreatedBy(nil)
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.CreatedByID)()
	if err != nil {
		return nil, err
	}

	return resolveCreatedBy(user)
}

func (g *GiftCard) Code(ctx context.Context) (string, error) {
	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	resolveCode := func(u *model.User) (string, error) {
		if (g.giftcard.UsedByEmail == nil && embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageGiftcard)) ||
			(u != nil && u.Id == embedCtx.AppContext.Session().UserId) {
			return g.giftcard.Code, nil
		}

		return "", model.NewAppError("GiftCard.Code", ErrorUnauthorized, nil, "You are not authorized to perform this action", http.StatusUnauthorized)
	}

	if g.giftcard.UsedByID == nil {
		return resolveCode(nil)
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.UsedByID)()
	if err != nil {
		return "", err
	}

	return resolveCode(user)
}

func (g *GiftCard) BoughtInChannel(ctx context.Context) (*string, error) {
	events, _ := GiftCardEventsByGiftCardIdLoader.Load(ctx, g.ID)()

	var boughtEvent *model.GiftCardEvent
	for _, evt := range events {
		if evt.Type == model.BOUGHT {
			boughtEvent = evt
			break
		}
	}

	if boughtEvent == nil {
		return nil, nil
	}

	orderID := boughtEvent.Parameters.Get("order_id", "").(string)
	if orderID == "" {
		return nil, errors.New("bought event's parameters field has no 'order_id' key")
	}

	order, err := OrderByIdLoader.Load(ctx, orderID)()
	if err != nil {
		return nil, err
	}

	if order == nil || !model.IsValidId(order.ChannelID) {
		return nil, nil
	}

	channel, err := ChannelByIdLoader.Load(ctx, order.ChannelID)()
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, nil
	}

	return &channel.Slug, nil
}

func giftCardsByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*model.GiftCard] {
	var (
		res         = make([]*dataloader.Result[[]*model.GiftCard], len(userIDs))
		giftcardMap = map[string][]*model.GiftCard{} // keys are user ids
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	giftcards, appErr := embedCtx.
		App.
		Srv().
		GiftcardService().
		GiftcardsByOption(nil, &model.GiftCardFilterOption{
			UsedByID: squirrel.Eq{store.GiftcardTableName + ".UsedByID": userIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	for _, gc := range giftcards {
		if gc.UsedByID != nil {
			giftcardMap[*gc.UsedByID] = append(giftcardMap[*gc.UsedByID], gc)
		}
	}

	for idx, id := range userIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Data: giftcardMap[id]}
	}
	return res

errorLabel:
	for idx := range userIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Error: appErr}
	}
	return res
}

func giftCardEventsByGiftCardIdLoader(ctx context.Context, giftcardIDs []string) []*dataloader.Result[[]*model.GiftCardEvent] {
	var (
		res              = make([]*dataloader.Result[[]*model.GiftCardEvent], len(giftcardIDs))
		giftcardEventMap = map[string][]*model.GiftCardEvent{} // keys are giftcard ids
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	events, appErr := embedCtx.App.Srv().
		GiftcardService().
		GiftcardEventsByOptions(&model.GiftCardEventFilterOption{
			GiftcardID: squirrel.Eq{store.GiftcardTableName + ".GiftcardID": giftcardIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	for _, event := range events {
		giftcardEventMap[event.GiftcardID] = append(giftcardEventMap[event.GiftcardID], event)
	}

	for idx, id := range giftcardIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCardEvent]{Data: giftcardEventMap[id]}
	}
	return res

errorLabel:
	for idx := range giftcardIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCardEvent]{Error: appErr}
	}
	return res
}

func giftcardsByOrderIDsLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.GiftCard] {
	var (
		res         = make([]*dataloader.Result[[]*model.GiftCard], len(orderIDs))
		giftcards   []*model.GiftCard
		appErr      *model.AppError
		giftcardMap = map[string]*model.GiftCard{} // keys are giftcard ids

		giftcardIDs []string
		cardMap     = map[string][]*model.GiftCard{} // keys are order ids
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)

	orderGiftcardRelations, err := embedCtx.App.Srv().Store.GiftCardOrder().FilterByOptions(&model.OrderGiftCardFilterOptions{
		OrderID: squirrel.Eq{store.OrderGiftCardTableName + ".OrderID": orderIDs},
	})
	if err != nil {
		err = model.NewAppError("giftcardsByOrderIDsLoader", "app.giftcard.giftcard-order-relations_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		goto errorLabel
	}

	giftcardIDs = lo.Map(orderGiftcardRelations, func(r *model.OrderGiftCard, _ int) string { return r.GiftCardID })

	giftcards, appErr = embedCtx.App.
		Srv().
		GiftcardService().
		GiftcardsByOption(nil, &model.GiftCardFilterOption{
			Id: squirrel.Eq{store.OrderGiftCardTableName + ".Id": giftcardIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, gc := range giftcards {
		giftcardMap[gc.Id] = gc
	}

	for _, rel := range orderGiftcardRelations {
		cardMap[rel.OrderID] = append(cardMap[rel.OrderID], giftcardMap[rel.GiftCardID])
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Data: cardMap[id]}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Error: err}
	}
	return res
}
