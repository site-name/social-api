package api

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type OrderEvent struct {
	ID                    string                       `json:"id"`
	Date                  *DateTime                    `json:"date"`
	Type                  *OrderEventsEnum             `json:"type"`
	User                  *User                        `json:"user"`
	Message               *string                      `json:"message"`
	Email                 *string                      `json:"email"`
	EmailType             *OrderEventsEmailsEnum       `json:"emailType"`
	Amount                *float64                     `json:"amount"`
	PaymentID             *string                      `json:"paymentId"`
	PaymentGateway        *string                      `json:"paymentGateway"`
	Quantity              *int32                       `json:"quantity"`
	ComposedID            *string                      `json:"composedId"`
	OrderNumber           *string                      `json:"orderNumber"`
	InvoiceNumber         *string                      `json:"invoiceNumber"`
	OversoldItems         []string                     `json:"oversoldItems"`
	Lines                 []*OrderEventOrderLineObject `json:"lines"`
	FulfilledItems        []*FulfillmentLine           `json:"fulfilledItems"`
	Warehouse             *Warehouse                   `json:"warehouse"`
	TransactionReference  *string                      `json:"transactionReference"`
	ShippingCostsIncluded *bool                        `json:"shippingCostsIncluded"`
	RelatedOrder          *Order                       `json:"relatedOrder"`
	Discount              *OrderEventDiscountObject    `json:"discount"`
}

func SystemOrderEventToGraphqlOrderEvent(o *model.OrderEvent) *OrderEvent {
	if o == nil {
		return &OrderEvent{}
	}
	panic("not implemented")
}

func orderEventsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.OrderEvent] {
	var (
		res      = make([]*dataloader.Result[[]*model.OrderEvent], len(orderIDs))
		events   []*model.OrderEvent
		eventMap = map[string][]*model.OrderEvent{}
		appErr   *model.AppError
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	events, appErr = embedCtx.App.Srv().OrderService().FilterOrderEventsByOptions(&model.OrderEventFilterOptions{
		OrderID: squirrel.Eq{store.OrderEventTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, event := range events {
		eventMap[event.OrderID] = append(eventMap[event.OrderID], event)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderEvent]{Data: eventMap[id]}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderEvent]{Error: err}
	}
	return res
}
