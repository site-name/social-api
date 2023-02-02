package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Channel struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	IsActive       bool            `json:"isActive"`
	Slug           string          `json:"slug"`
	CurrencyCode   string          `json:"currencyCode"`
	DefaultCountry *CountryDisplay `json:"defaultCountry"`

	// HasOrders      bool            `json:"hasOrders"`
}

func (c *Channel) HasOrders(ctx context.Context) (bool, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return false, err
	}

	// check if current user has channel management
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageChannels) {
		return false, model.NewAppError("Channel.HasOrders", ErrorUnauthorized, nil, "you are not allowed to perform this", http.StatusUnauthorized)
	}

	channel, err := ChannelWithHasOrdersByIdLoader.Load(ctx, c.ID)()
	if err != nil {
		return false, err
	}

	return channel.GetHasOrders(), nil
}

func SystemChannelToGraphqlChannel(ch *model.Channel) *Channel {
	if ch == nil {
		return nil
	}

	return &Channel{
		ID:           ch.Id,
		Name:         ch.Name,
		IsActive:     ch.IsActive,
		Slug:         ch.Slug,
		CurrencyCode: ch.Currency,
		DefaultCountry: &CountryDisplay{
			Code:    ch.DefaultCountry,
			Country: model.Countries[strings.ToUpper(ch.DefaultCountry)],
		},
	}
}

func channelByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Channel] {
	var (
		res        = make([]*dataloader.Result[*model.Channel], len(ids))
		appErr     *model.AppError
		channels   model.Channels
		channelMap = map[string]*model.Channel{} // keys are channel ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	channels, appErr = embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			Id: squirrel.Eq{store.ChannelTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	channelMap = lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Id, c })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Channel]{Error: err}
	}
	return res
}

func channelBySlugLoader(ctx context.Context, slugs []string) []*dataloader.Result[*model.Channel] {
	var (
		res        = make([]*dataloader.Result[*model.Channel], len(slugs))
		appErr     *model.AppError
		channels   model.Channels
		channelMap = map[string]*model.Channel{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	channels, appErr = embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			Slug: squirrel.Eq{store.ChannelTableName + ".Slug": slugs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	channelMap = lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Slug, c })

	for idx, slug := range slugs {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[slug]}
	}
	return res

errorLabel:
	for idx := range slugs {
		res[idx] = &dataloader.Result[*model.Channel]{Error: err}
	}
	return res
}

func channelByCheckoutLineIDLoader(ctx context.Context, checkoutLineIDs []string) []*dataloader.Result[*model.Channel] {
	var (
		res            []*dataloader.Result[*model.Channel]
		errs           []error
		checkouts      []*model.Checkout
		channels       []*model.Channel
		checkoutTokens []string
		channelIDs     []string
	)

	// find checkout lines
	checkoutLines, errs := CheckoutLineByIdLoader.LoadMany(ctx, checkoutLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	checkoutTokens = lo.Map(checkoutLines, func(item *model.CheckoutLine, _ int) string { return item.CheckoutID })

	// find checkouts
	checkouts, errs = CheckoutByTokenLoader.LoadMany(ctx, checkoutTokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	channelIDs = lo.Map(checkouts, func(item *model.Checkout, _ int) string { return item.ChannelID })

	// find channels
	channels, errs = ChannelByIdLoader.LoadMany(ctx, channelIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(ch *model.Channel, _ int) *dataloader.Result[*model.Channel] {
		return &dataloader.Result[*model.Channel]{Data: ch}
	})

errorLabel:
	for range checkoutLineIDs {
		res = append(res, &dataloader.Result[*model.Channel]{Error: errs[0]})
	}
	return res
}

func channelByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*model.Channel] {
	var (
		res        []*dataloader.Result[*model.Channel]
		orders     model.Orders
		channels   []*model.Channel
		orderLines model.OrderLines
		errs       []error
	)

	// find order lines
	orderLines, errs = OrderLineByIdLoader.LoadMany(ctx, orderLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	// find orders
	orders, errs = OrderByIdLoader.LoadMany(ctx, orderLines.OrderIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	// find channels
	channels, errs = ChannelByIdLoader.LoadMany(ctx, orders.ChannelIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(i *model.Channel, _ int) *dataloader.Result[*model.Channel] {
		return &dataloader.Result[*model.Channel]{Data: i}
	})

errorLabel:
	for range orderLineIDs {
		res = append(res, &dataloader.Result[*model.Channel]{Error: errs[0]})
	}
	return res
}

func channelWithHasOrdersByIdLoader(ctx context.Context, channelIDs []string) []*dataloader.Result[*model.Channel] {
	var (
		res        = make([]*dataloader.Result[*model.Channel], len(channelIDs))
		channels   model.Channels
		appErr     *model.AppError
		channelMap = map[string]*model.Channel{}
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	// find all channels that have orders
	channels, appErr = embedCtx.App.
		Srv().
		ChannelService().
		ChannelsByOption(&model.ChannelFilterOption{
			AnnotateHasOrders: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	channelMap = lo.SliceToMap(channels, func(c *model.Channel) (string, *model.Channel) { return c.Id, c })

	for idx, id := range channelIDs {
		res[idx] = &dataloader.Result[*model.Channel]{Data: channelMap[id]}
	}
	return res

errorLabel:
	for range channelIDs {
		res = append(res, &dataloader.Result[*model.Channel]{Error: err})
	}
	return res
}
