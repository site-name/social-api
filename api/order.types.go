package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// --------------------------- Order line -----------------------------

func SystemOrderLineToGraphqlOrderLine(line *model.OrderLine) *OrderLine {
	if line == nil {
		return nil
	}

	res := &OrderLine{
		ID:                    line.Id,
		ProductName:           line.ProductName,
		VariantName:           line.VariantName,
		ProductSku:            line.ProductSku,
		ProductVariantID:      line.ProductVariantID,
		IsShippingRequired:    line.IsShippingRequired,
		TranslatedProductName: line.TranslatedProductName,
		TranslatedVariantName: line.TranslatedVariantName,
		Quantity:              int32(line.Quantity),
		QuantityFulfilled:     int32(line.QuantityFulfilled),
		UnitDiscountReason:    line.UnitDiscountReason,
		UnitPrice:             SystemTaxedMoneyToGraphqlTaxedMoney(line.UnitPrice),
		UndiscountedUnitPrice: SystemTaxedMoneyToGraphqlTaxedMoney(line.UnDiscountedUnitPrice),
		UnitDiscount:          SystemMoneyToGraphqlMoney(line.UnitDiscount),
		UnitDiscountValue:     PositiveDecimal(*line.UnitDiscountValue),
		TotalPrice:            SystemTaxedMoneyToGraphqlTaxedMoney(line.TotalPrice),
		QuantityToFulfill:     int32(line.QuantityUnFulfilled()),
		variantID:             line.VariantID,
	}
	discountType := DiscountValueTypeEnum(strings.ToUpper(line.UnitDiscountType))
	res.UnitDiscountType = &discountType

	if line.TaxRate != nil {
		res.TaxRate, _ = line.TaxRate.Float64()
	}

	return res
}

func (o *OrderLine) Variant(ctx context.Context) (*ProductVariant, error) {
	if o.variantID == nil {
		return nil, nil
	}

	panic("not implemented")
}

func graphqlOrderLinesByIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*OrderLine] {
	var (
		res        []*dataloader.Result[*OrderLine]
		appErr     *model.AppError
		orderLines []*model.OrderLine
	)
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	orderLines, appErr = embedCtx.App.
		Srv().
		OrderService().
		OrderLinesByOption(&model.OrderLineFilterOption{
			Id: squirrel.Eq{store.OrderLineTableName + ".Id": orderLineIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, orderLine := range orderLines {
		res = append(res, &dataloader.Result[*OrderLine]{Data: SystemOrderLineToGraphqlOrderLine(orderLine)})
	}
	return res

errorLabel:
	for range orderLineIDs {
		res = append(res, &dataloader.Result[*OrderLine]{Error: err})
	}
	return res
}
