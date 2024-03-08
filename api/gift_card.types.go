package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
	"golang.org/x/text/currency"
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

	e *model.GiftcardEvent

	// User          *User                 `json:"user"`
	// App           *App                  `json:"app"`
}

func SystemGiftcardEventToGraphqlGiftcardEvent(evt *model.GiftcardEvent) *GiftCardEvent {
	if evt == nil {
		return nil
	}

	res := new(GiftCardEvent)
	res.ID = evt.ID
	res.e = evt
	if evt.Date != 0 {
		res.Date = &DateTime{util.TimeFromMillis(evt.Date)}
	}
	res.Type = &evt.Type

	msg, ok := evt.Parameters["message"]
	if ok && msg != nil {
		res.Message = model_helper.GetPointerOfValue(msg.(string))
	}

	email, ok := evt.Parameters["email"]
	if ok && email != nil {
		res.Email = model_helper.GetPointerOfValue(email.(string))
	}

	orderID, ok := evt.Parameters["order_id"]
	if ok && orderID != nil {
		res.OrderID = model_helper.GetPointerOfValue(orderID.(string))
	}

	tag, ok := evt.Parameters["tag"]
	if ok && tag != nil {
		res.Tag = model_helper.GetPointerOfValue(tag.(string))
	}

	oldTag, ok := evt.Parameters["old_tag"]
	if ok && oldTag != nil {
		res.OldTag = model_helper.GetPointerOfValue(oldTag.(string))
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
				tim, err := time.Parse(time.DateOnly, t)
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

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (e *GiftCardEvent) User(ctx context.Context) (*User, error) {
	if e.e.UserID == nil {
		return nil, nil
	}
	user, err := UserByUserIdLoader.Load(ctx, *e.e.UserID)()
	if err != nil {
		return nil, err
	}
	return SystemUserToGraphqlUser(user), nil
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
	Code            string          `json:"code"`

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
	res.Code = gc.Code
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

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (g *GiftCard) Events(ctx context.Context) ([]*GiftCardEvent, error) {
	events, err := GiftCardEventsByGiftCardIdLoader.Load(ctx, g.ID)()
	if err != nil {
		return nil, err
	}

	return systemRecordsToGraphql(events, SystemGiftcardEventToGraphqlGiftcardEvent), nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (g *GiftCard) CreatedByEmail(ctx context.Context) (*string, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if g.giftcard.CreatedByID != nil {
		user, appErr := embedCtx.App.Srv().AccountService().UserById(ctx, *g.giftcard.CreatedByID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			return g.giftcard.CreatedByEmail, nil
		}

		return &user.Email, nil
	}

	return g.giftcard.CreatedByEmail, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (g *GiftCard) UsedByEmail(ctx context.Context) (*string, error) {
	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.UsedByID)()
	if err != nil {
		return nil, err
	}

	return &user.Email, nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (g *GiftCard) UsedBy(ctx context.Context) (*User, error) {
	if g.giftcard.UsedByID == nil {
		return nil, nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.UsedByID)()
	if err != nil {
		return nil, err
	}
	return SystemUserToGraphqlUser(user), nil
}

// NOTE: Refer to ./schemas/gift_card.graphqls for details on directive used.
func (g *GiftCard) CreatedBy(ctx context.Context) (*User, error) {
	if g.giftcard.CreatedByID == nil {
		return nil, nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *g.giftcard.CreatedByID)()
	if err != nil {
		return nil, err
	}

	return SystemUserToGraphqlUser(user), nil
}

func (g *GiftCard) BoughtInChannel(ctx context.Context) (*string, error) {
	events, err := GiftCardEventsByGiftCardIdLoader.Load(ctx, g.ID)()
	if err != nil {
		return nil, err
	}

	var boughtEvent *model.GiftCardEvent
	for _, evt := range events {
		if evt.Type == model.GIFT_CARD_EVENT_TYPE_BOUGHT {
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

	if order == nil || !model_helper.IsValidId(order.ChannelID) {
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

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	_, giftcards, appErr := embedCtx.
		App.
		Srv().
		GiftcardService().
		GiftcardsByOption(&model.GiftCardFilterOption{
			Conditions: squirrel.Eq{model.GiftcardTableName + ".UsedByID": userIDs},
		})
	if appErr != nil {
		for idx := range userIDs {
			res[idx] = &dataloader.Result[[]*model.GiftCard]{Error: appErr}
		}
		return res
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
}

func giftCardEventsByGiftCardIdLoader(ctx context.Context, giftcardIDs []string) []*dataloader.Result[[]*model.GiftCardEvent] {
	var (
		res              = make([]*dataloader.Result[[]*model.GiftCardEvent], len(giftcardIDs))
		giftcardEventMap = map[string][]*model.GiftCardEvent{} // keys are giftcard ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	events, appErr := embedCtx.App.Srv().
		GiftcardService().
		GiftcardEventsByOptions(&model.GiftCardEventFilterOption{
			Conditions: squirrel.Eq{model.GiftcardTableName + ".GiftcardID": giftcardIDs},
		})
	if appErr != nil {
		for idx := range giftcardIDs {
			res[idx] = &dataloader.Result[[]*model.GiftCardEvent]{Error: appErr}
		}
		return res
	}

	for _, event := range events {
		giftcardEventMap[event.GiftcardID] = append(giftcardEventMap[event.GiftcardID], event)
	}

	for idx, id := range giftcardIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCardEvent]{Data: giftcardEventMap[id]}
	}
	return res
}

func giftcardsByOrderIDsLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.GiftCard] {
	var (
		res      = make([]*dataloader.Result[[]*model.GiftCard], len(orderIDs))
		orderMap = map[string]*model.Order{} // keys are order ids
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	var orders model.Orders
	err := embedCtx.App.Srv().Store.GetReplica().Preload("GiftCards").Find(&orders, "Id IN ?", orderIDs).Error
	if err != nil {
		goto errorLabel
	}

	for _, order := range orders {
		orderMap[order.Id] = order
	}
	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Data: orderMap[id].GiftCards}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.GiftCard]{Error: err}
	}
	return res
}

type GiftCardFilterInput struct {
	IsActive       *bool            `json:"isActive"`
	Tag            *string          `json:"tag"`
	Tags           []string         `json:"tags"`
	Products       []string         `json:"products"` // product ids
	UsedBy         []string         `json:"usedBy"`   // user ids
	Currency       *string          `json:"currency"` // should be upper-cased
	CurrentBalance *PriceRangeInput `json:"currentBalance"`
	InitialBalance *PriceRangeInput `json:"initialBalance"`
}

func (g *GiftCardFilterInput) validate() *model_helper.AppError {
	if g.Tag != nil && *g.Tag == "" {
		return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "tag"}, "tag must not be empty", http.StatusBadRequest)
	}
	if len(g.Products) > 0 && !lo.EveryBy(g.Products, model_helper.IsValidId) {
		return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "products"}, "please provide valid product ids", http.StatusBadRequest)
	}
	if len(g.UsedBy) > 0 && !lo.EveryBy(g.UsedBy, model_helper.IsValidId) {
		return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "usedBy"}, "please provide valid user ids", http.StatusBadRequest)
	}
	if g.Currency != nil {
		_, err := currency.ParseISO(*g.Currency)
		if err != nil {
			return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "currency"}, *g.Currency+" is not a valid currency", http.StatusBadRequest)
		}
	}
	if g.Currency == nil && (g.CurrentBalance != nil || g.InitialBalance != nil) {
		return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "currency"}, "please provide a currency", http.StatusBadRequest)
	}

	for name, priceRange := range map[string]*PriceRangeInput{
		"CurrentBalance": g.CurrentBalance,
		"InitialBalance": g.InitialBalance,
	} {
		if priceRange != nil &&
			priceRange.Gte != nil &&
			priceRange.Lte != nil &&
			*priceRange.Gte >= *priceRange.Lte {
			return model_helper.NewAppError("GiftCardFilterInput.Validate", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": name}, "Lte must be greater than Gte", http.StatusBadRequest)
		}
	}

	return nil
}

// NOTE: Call me after calling validate()
func (g *GiftCardFilterInput) ToSystemGiftcardFilter() (*model.GiftCardFilterOption, *model_helper.AppError) {
	appErr := g.validate()
	if appErr != nil {
		return nil, appErr
	}

	andConditions := squirrel.And{}

	if g.IsActive != nil {
		andConditions = append(andConditions, squirrel.Eq{model.GiftcardTableName + ".IsActive": *g.IsActive})
	}
	var tagFilter = squirrel.And{}
	if g.Tag != nil && *g.Tag != "" {
		tagFilter = append(tagFilter, squirrel.ILike{model.GiftcardTableName + ".Tag": "%" + *g.Tag + "%"})
	}
	if len(g.Tags) > 0 {
		tagFilter = append(tagFilter, squirrel.Eq{model.GiftcardTableName + ".Tag": g.Tags})
	}
	if len(tagFilter) > 0 {
		andConditions = append(andConditions, tagFilter)
	}

	if len(g.Products) > 0 {
		andConditions = append(andConditions, squirrel.Eq{model.GiftcardTableName + ".ProductID": g.Products})
	}
	if len(g.UsedBy) > 0 {
		andConditions = append(andConditions, squirrel.Eq{model.GiftcardTableName + ".UsedByID": g.UsedBy})
	}
	if g.Currency != nil {
		andConditions = append(andConditions, squirrel.Eq{model.GiftcardTableName + ".Currency": strings.ToUpper(*g.Currency)})
	}

	for fieldName, priceRange := range map[string]*PriceRangeInput{
		".CurrentBalanceAmount": g.CurrentBalance,
		".InitialBalanceAmount": g.InitialBalance,
	} {
		if priceRange != nil {
			balanceConditions := squirrel.And{}
			if gte := priceRange.Gte; gte != nil {
				balanceConditions = append(balanceConditions, squirrel.GtOrEq{model.GiftcardTableName + fieldName: *gte})
			}
			if lte := priceRange.Lte; lte != nil {
				balanceConditions = append(balanceConditions, squirrel.LtOrEq{model.GiftcardTableName + fieldName: *lte})
			}
			andConditions = append(andConditions, balanceConditions)
		}
	}

	return &model.GiftCardFilterOption{Conditions: andConditions}, nil
}
