package api

import (
	"context"
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

func (c *Channel) HasOrders(ctx context.Context) (*bool, error) {
	panic("not implemented")
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

func channelByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Channel] {
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
			Id: squirrel.Eq{store.ChannelTableName + ".Id": ids},
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
	for range ids {
		res = append(res, &dataloader.Result[*Channel]{Error: err})
	}
	return res
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
	for _, err := range errs {
		res = append(res, &dataloader.Result[*Channel]{Error: err})
	}
	return res
}

func channelByOrderLineIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[[]*Channel] {
	panic("not implemented")
}

func channelWithHasOrdersByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Channel] {
	panic("not implemented")
}
