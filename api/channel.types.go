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
	hasOrders bool
}

func (c *Channel) HasOrders(ctx context.Context) (*bool, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// check if current user has channel management
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageChannels) {
		return nil, model.NewAppError("Channel.HasOrders", ErrorUnauthorized, nil, "you are not allowed to perform this", http.StatusUnauthorized)
	}

	channel, err := dataloaders.ChannelWithHasOrdersByIdLoader.Load(ctx, c.ID)()
	if err != nil {
		return nil, err
	}

	return model.NewBool(channel.hasOrders), nil
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
		hasOrders: ch.GetHasOrders(),
	}
}

func channelByIdLoader_systemResult(ctx context.Context, ids []string) []*dataloader.Result[*model.Channel] {
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

func channelByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Channel] {
	results := channelByIdLoader_systemResult(ctx, ids)

	return lo.Map(results, func(r *dataloader.Result[*model.Channel], _ int) *dataloader.Result[*Channel] {
		return &dataloader.Result[*Channel]{
			Data:  SystemChannelToGraphqlChannel(r.Data),
			Error: r.Error,
		}
	})
}

func channelBySlugLoader(ctx context.Context, slugs []string) []*dataloader.Result[*Channel] {
	var (
		res      []*dataloader.Result[*Channel]
		appErr   *model.AppError
		channels model.Channels
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

	for _, channel := range channels {
		res = append(res, &dataloader.Result[*Channel]{Data: SystemChannelToGraphqlChannel(channel)})
	}
	return res

errorLabel:
	for range slugs {
		res = append(res, &dataloader.Result[*Channel]{Error: err})
	}
	return res
}

func channelByCheckoutLineIDLoader(ctx context.Context, checkoutLineIDs []string) []*dataloader.Result[*Channel] {
	var (
		res            []*dataloader.Result[*Channel]
		errs           []error
		checkouts      []*Checkout
		channels       []*Channel
		checkoutTokens []string
		channelIDs     []string
	)

	// find checkout lines
	checkoutLines, errs := dataloaders.CheckoutLineByIdLoader.LoadMany(ctx, checkoutLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	checkoutTokens = lo.Map(checkoutLines, func(item *CheckoutLine, _ int) string { return item.checkoutID })

	// find checkouts
	checkouts, errs = dataloaders.CheckoutByTokenLoader.LoadMany(ctx, checkoutTokens)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	channelIDs = lo.Map(checkouts, func(item *Checkout, _ int) string { return item.channelID })

	// find channels
	channels, errs = dataloaders.ChannelByIdLoader.LoadMany(ctx, channelIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(ch *Channel, _ int) *dataloader.Result[*Channel] {
		return &dataloader.Result[*Channel]{Data: ch}
	})

errorLabel:
	for range checkoutLineIDs {
		res = append(res, &dataloader.Result[*Channel]{Error: errs[0]})
	}
	return res
}

func channelByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*Channel] {
	var (
		res        []*dataloader.Result[*Channel]
		orders     model.Orders
		channels   []*Channel
		orderLines model.OrderLines
		errs       []error
	)

	// find order lines
	orderLines, errs = dataloaders.OrderLineByIdLoader.LoadMany(ctx, orderLineIDs)()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	// find orders
	orders, errs = dataloaders.OrderByIdLoader.LoadMany(ctx, orderLines.OrderIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	// find channels
	channels, errs = dataloaders.ChannelByIdLoader.LoadMany(ctx, orders.ChannelIDs())()
	if len(errs) > 0 && errs[0] != nil {
		goto errorLabel
	}

	return lo.Map(channels, func(i *Channel, _ int) *dataloader.Result[*Channel] {
		return &dataloader.Result[*Channel]{Data: i}
	})

errorLabel:
	for range orderLineIDs {
		res = append(res, &dataloader.Result[*Channel]{Error: errs[0]})
	}
	return res
}

func channelWithHasOrdersByIdLoader(ctx context.Context, channelIDs []string) []*dataloader.Result[*Channel] {
	var (
		res      []*dataloader.Result[*Channel]
		channels model.Channels
		appErr   *model.AppError
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

	return lo.Map(channels, func(c *model.Channel, _ int) *dataloader.Result[*Channel] {
		return &dataloader.Result[*Channel]{Data: SystemChannelToGraphqlChannel(c)}
	})

errorLabel:
	for range channelIDs {
		res = append(res, &dataloader.Result[*Channel]{Error: err})
	}
	return res
}
