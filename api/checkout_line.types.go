package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type CheckoutLine struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`

	// RequiresShipping *bool           `json:"requiresShipping"`
	// TotalPrice       *TaxedMoney     `json:"totalPrice"`
	// Variant          *ProductVariant `json:"variant"`
}

func SystemCheckoutLineToGraphqlCheckoutLine(line *model.CheckoutLine) *CheckoutLine {
	if line == nil {
		return nil
	}

	res := &CheckoutLine{
		ID:       line.Id,
		Quantity: int32(line.Quantity),
	}
	return res
}

func (line *CheckoutLine) Variant(ctx context.Context) *ProductVariant {
	panic("not implemented")
}

func (line *CheckoutLine) TotalPrice(ctx context.Context) *TaxedMoney {
	panic("not implemented")
}

func (line *CheckoutLine) RequiresShipping(ctx context.Context) *bool {
	panic("not implemented")
}

func graphqlCheckoutLinesByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*CheckoutLine] {
	var (
		res           []*dataloader.Result[[]*CheckoutLine]
		appErr        *model.AppError
		checkoutLines model.CheckoutLines

		// checkoutLinesMap has keys are checkout tokens.
		// values are checkout lines belong to the checkout parent
		checkoutLinesMap = map[string][]*CheckoutLine{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkoutLines, appErr = embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			CheckoutID: squirrel.Eq{store.CheckoutLineTableName + ".": tokens},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range checkoutLines {
		if line != nil {
			checkoutLinesMap[line.CheckoutID] = append(
				checkoutLinesMap[line.CheckoutID],
				SystemCheckoutLineToGraphqlCheckoutLine(line))
		}
	}

	for _, token := range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLine]{Data: checkoutLinesMap[token]})
	}
	return res

errorLabel:
	for range tokens {
		res = append(res, &dataloader.Result[[]*CheckoutLine]{Error: err})
	}
	return res
}
