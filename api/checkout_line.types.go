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

	variantID  string
	checkoutID string

	// RequiresShipping *bool           `json:"requiresShipping"`
	// TotalPrice       *TaxedMoney     `json:"totalPrice"`
	// Variant          *ProductVariant `json:"variant"`
}

func SystemCheckoutLineToGraphqlCheckoutLine(line *model.CheckoutLine) *CheckoutLine {
	if line == nil {
		return nil
	}

	res := &CheckoutLine{
		ID:         line.Id,
		Quantity:   int32(line.Quantity),
		variantID:  line.VariantID,
		checkoutID: line.CheckoutID,
	}
	return res
}

func (line *CheckoutLine) Variant(ctx context.Context) (*ProductVariant, error) {
	return dataloaders.ProductVariantByIdLoader.Load(ctx, line.variantID)()
}

func (line *CheckoutLine) TotalPrice(ctx context.Context) (*TaxedMoney, error) {
	// checkout, err := dataloaders.CheckoutByTokenLoader.Load(ctx, line.checkoutID)()
	// if err != nil {
	// 	return nil, err
	// }
	panic("not implemented")
}

func (line *CheckoutLine) RequiresShipping(ctx context.Context) (*bool, error) {
	panic("not implemented")
}

func checkoutLinesByCheckoutTokenLoader(ctx context.Context, tokens []string) []*dataloader.Result[[]*CheckoutLine] {
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
			CheckoutID: squirrel.Eq{store.CheckoutLineTableName + ".CheckoutID": tokens},
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

func checkoutLineByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*CheckoutLine] {
	var (
		res             = make([]*dataloader.Result[*CheckoutLine], len(ids))
		checkoutLines   model.CheckoutLines
		appErr          *model.AppError
		checkoutLineMap = map[string]*CheckoutLine{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	checkoutLines, appErr = embedCtx.App.Srv().
		CheckoutService().
		CheckoutLinesByOption(&model.CheckoutLineFilterOption{
			Id: squirrel.Eq{store.CheckoutLineTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range checkoutLines {
		checkoutLineMap[line.Id] = SystemCheckoutLineToGraphqlCheckoutLine(line)
	}
	for idx, id := range ids {
		res[idx] = &dataloader.Result[*CheckoutLine]{Data: checkoutLineMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*CheckoutLine]{Error: err}
	}
	return res
}
