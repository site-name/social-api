package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type Fulfillment struct {
	ID               string            `json:"id"`
	FulfillmentOrder int32             `json:"fulfillmentOrder"`
	Status           FulfillmentStatus `json:"status"`
	TrackingNumber   string            `json:"trackingNumber"`
	Created          DateTime          `json:"created"`
	PrivateMetadata  []*MetadataItem   `json:"privateMetadata"`
	Metadata         []*MetadataItem   `json:"metadata"`

	fulfillment *model.Fulfillment

	// Lines            []*FulfillmentLine `json:"lines"`
	// StatusDisplay    *string            `json:"statusDisplay"`
	// Warehouse        *Warehouse         `json:"warehouse"`
}

func SystemFulfillmentToGraphqlFulfillment(fulfillment *model.Fulfillment) *Fulfillment {
	if fulfillment == nil {
		return nil
	}

	return &Fulfillment{
		ID:               fulfillment.Id,
		FulfillmentOrder: int32(fulfillment.FulfillmentOrder),
		Status:           FulfillmentStatus(fulfillment.Status),
		TrackingNumber:   fulfillment.TrackingNumber,
		Created:          DateTime{util.TimeFromMillis(fulfillment.CreateAt)},
		Metadata:         MetadataToSlice(fulfillment.Metadata),
		PrivateMetadata:  MetadataToSlice(fulfillment.PrivateMetadata),

		fulfillment: fulfillment,
	}
}

func (f *Fulfillment) Lines(ctx context.Context) ([]*FulfillmentLine, error) {
	lines, err := FulfillmentLinesByFulfillmentIDLoader.Load(ctx, f.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(lines, SystemFulfillmentLineToGraphqlFulfillmentLine), nil
}

func (f *Fulfillment) StatusDisplay(ctx context.Context) (*string, error) {
	res := model.FulfillmentStrings[f.fulfillment.Status]
	return &res, nil
}

func (f *Fulfillment) Warehouse(ctx context.Context) (*Warehouse, error) {
	fulfillmentLines, err := FulfillmentLinesByFulfillmentIDLoader.Load(ctx, f.ID)()
	if err != nil {
		return nil, err
	}

	if len(fulfillmentLines) > 0 && fulfillmentLines[0].StockID != nil {
		stock, err := StocksByIDLoader.Load(ctx, *fulfillmentLines[0].StockID)()
		if err != nil {
			return nil, err
		}

		if stock.GetWarehouse() != nil {
			return SystemWarehouseToGraphqlWarehouse(stock.GetWarehouse()), nil
		}

		return nil, nil
	}

	return nil, nil
}

func fulfillmentsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.Fulfillment] {
	var (
		res            = make([]*dataloader.Result[[]*model.Fulfillment], len(orderIDs))
		fulfillmentMap = map[string]model.Fulfillments{}
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	fulfillments, appErr := embedCtx.App.Srv().OrderService().FulfillmentsByOption(nil, &model.FulfillmentFilterOption{
		OrderID: squirrel.Eq{store.FulfillmentTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		goto errorLabel
	}

	for _, f := range fulfillments {
		fulfillmentMap[f.OrderID] = append(fulfillmentMap[f.OrderID], f)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Fulfillment]{Data: fulfillmentMap[id]}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.Fulfillment]{Error: appErr}
	}
	return res
}

// ------------ fulfillment line -------------------

type FulfillmentLine struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`

	// OrderLine *OrderLine `json:"orderLine"`
	fml *model.FulfillmentLine
}

func SystemFulfillmentLineToGraphqlFulfillmentLine(fml *model.FulfillmentLine) *FulfillmentLine {
	if fml == nil {
		return nil
	}

	return &FulfillmentLine{
		ID:       fml.Id,
		Quantity: int32(fml.Quantity),
		fml:      fml,
	}
}

func (f *FulfillmentLine) OrderLine(ctx context.Context) (*OrderLine, error) {
	orderLine, err := OrderLineByIdLoader.Load(ctx, f.fml.OrderLineID)()
	if err != nil {
		return nil, err
	}

	return SystemOrderLineToGraphqlOrderLine(orderLine), nil
}

func fulfillmentLinesByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.FulfillmentLine] {
	var (
		res     = make([]*dataloader.Result[*model.FulfillmentLine], len(ids))
		LineMap = map[string]*model.FulfillmentLine{}
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	lines, appErr := embedCtx.App.Srv().OrderService().FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Id: squirrel.Eq{store.FulfillmentLineTableName + ".Id": ids},
	})
	if appErr != nil {
		goto errorLabel
	}

	LineMap = lo.SliceToMap(lines, func(l *model.FulfillmentLine) (string, *model.FulfillmentLine) { return l.Id, l })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.FulfillmentLine]{Data: LineMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.FulfillmentLine]{Error: appErr}
	}
	return res
}

func fulfillmentLinesByFulfillmentIDLoader(ctx context.Context, fulfillmentIDs []string) []*dataloader.Result[[]*model.FulfillmentLine] {
	var (
		res                = make([]*dataloader.Result[[]*model.FulfillmentLine], len(fulfillmentIDs))
		fulfillmentLineMap = map[string]model.FulfillmentLines{} // keys are fulfillment ids
	)

	embedCtx, _ := GetContextValue[*web.Context](ctx, WebCtx)
	fulfillmentLines, appErr := embedCtx.App.
		Srv().
		OrderService().
		FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
			FulfillmentID: squirrel.Eq{store.FulfillmentLineTableName + ".FulfillmentID": fulfillmentIDs},
		})
	if appErr != nil {
		goto errorLabel
	}

	for _, line := range fulfillmentLines {
		fulfillmentLineMap[line.FulfillmentID] = append(fulfillmentLineMap[line.FulfillmentID], line)
	}

	for idx, id := range fulfillmentIDs {
		res[idx] = &dataloader.Result[[]*model.FulfillmentLine]{Data: fulfillmentLineMap[id]}
	}
	return res

errorLabel:
	for idx := range fulfillmentIDs {
		res[idx] = &dataloader.Result[[]*model.FulfillmentLine]{Error: appErr}
	}
	return res
}
