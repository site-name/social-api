package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func SystemCheckoutToGraphqlCheckout(ckout *model.Checkout) *Checkout {
	if ckout == nil {
		return nil
	}

	res := &Checkout{}
	panic("not implemented")
	return res
}

// arrays inside keys have format of [userID, channelID]
func graphqlCheckoutsByUserAndChannelLoader(ctx context.Context, keys [][2]string) ([]*dataloader.Result[[]*Checkout], error) {
	var (
		appErr                  *model.AppError
		res                     []*dataloader.Result[[]*Checkout]
		userIDs                 []string
		channelIDs              []string
		userID_checkoutID_pairs []string
		checkouts               []*model.Checkout
		checkoutsMap            = map[string][]*Checkout{}
	)

	for _, item := range keys {
		if len(item) != 2 {
			continue
		}

		userIDs = append(userIDs, item[0])
		channelIDs = append(channelIDs, item[1])
		userID_checkoutID_pairs = append(userID_checkoutID_pairs, item[0]+item[1])
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkouts, appErr = embedCtx.App.Srv().
		CheckoutService().
		CheckoutsByOption(&model.CheckoutFilterOption{
			UserID:          squirrel.Eq{store.CheckoutTableName + "UserID": userIDs},
			ChannelID:       squirrel.Eq{store.CheckoutTableName + "ChannelID": channelIDs},
			ChannelIsActive: model.NewBool(true),
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, checkout := range checkouts {
		if checkout.UserID != nil {
			key := *checkout.UserID + checkout.ChannelID
			checkoutsMap[key] = append(checkoutsMap[key], SystemCheckoutToGraphqlCheckout(checkout))
		}
	}

	for _, key := range userID_checkoutID_pairs {
		res = append(res, &dataloader.Result[[]*Checkout]{Data: checkoutsMap[key]})
	}
	return res, nil

errorLabel:
	for range keys {
		res = append(res, &dataloader.Result[[]*Checkout]{Error: err})
	}
	return res, nil
}
