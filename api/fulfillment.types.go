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

	// Lines            []*FulfillmentLine `json:"lines"`
	// StatusDisplay    *string            `json:"statusDisplay"`
	// Warehouse        *Warehouse         `json:"warehouse"`
}

func SystemFulfillmentToGraphqlFulfillment(f *model.Fulfillment) *Fulfillment {
	if f == nil {
		return &Fulfillment{}
	}

	return &Fulfillment{
		ID:               f.Id,
		FulfillmentOrder: int32(f.FulfillmentOrder),
		Status:           FulfillmentStatus(f.Status),
		TrackingNumber:   f.TrackingNumber,
		Created:          DateTime{util.TimeFromMillis(f.CreateAt)},
		Metadata:         MetadataToSlice(f.Metadata),
		PrivateMetadata:  MetadataToSlice(f.PrivateMetadata),
	}
}

func (f *Fulfillment) Lines(ctx context.Context) ([]*FulfillmentLine, error) {
	panic("not implemented")
}

func (f *Fulfillment) StatusDisplay(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (f *Fulfillment) Warehouse(ctx context.Context) (*Warehouse, error) {
	panic("not implemented")
}

func fulfillmentsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.Fulfillment] {
	var (
		res            = make([]*dataloader.Result[[]*model.Fulfillment], len(orderIDs))
		appErr         *model.AppError
		fulfillmentMap = map[string]model.Fulfillments{}
		fulfillments   model.Fulfillments
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	fulfillments, appErr = embedCtx.App.Srv().OrderService().FulfillmentsByOption(nil, &model.FulfillmentFilterOption{
		OrderID: squirrel.Eq{store.FulfillmentTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		err = appErr
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
		res[idx] = &dataloader.Result[[]*model.Fulfillment]{Error: err}
	}
	return res
}

// ------------

type FulfillmentLine struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
	// OrderLine *OrderLine `json:"orderLine"`
}

func SystemFulfillmentLineToGraphqlFulfillmentLine(l *model.FulfillmentLine) *FulfillmentLine {
	if l == nil {
		return &FulfillmentLine{}
	}

	return &FulfillmentLine{
		ID:       l.Id,
		Quantity: int32(l.Quantity),
	}
}

func (f *FulfillmentLine) OrderLine(ctx context.Context) (*OrderLine, error) {
	panic("not implemented")
}

func fulfillmentLinesByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.FulfillmentLine] {
	var (
		res     = make([]*dataloader.Result[*model.FulfillmentLine], len(ids))
		lines   model.FulfillmentLines
		appErr  *model.AppError
		LineMap = map[string]*model.FulfillmentLine{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	lines, appErr = embedCtx.App.Srv().OrderService().FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Id: squirrel.Eq{store.FulfillmentLineTableName + ".Id": ids},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	LineMap = lo.SliceToMap(lines, func(l *model.FulfillmentLine) (string, *model.FulfillmentLine) { return l.Id, l })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.FulfillmentLine]{Data: LineMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.FulfillmentLine]{Error: err}
	}
	return res
}
