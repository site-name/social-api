package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

type OrderEvent struct {
	ID                    string                 `json:"id"`
	Date                  *DateTime              `json:"date"`
	Type                  *OrderEventsEnum       `json:"type"`
	Message               *string                `json:"message"`
	Email                 *string                `json:"email"`
	EmailType             *OrderEventsEmailsEnum `json:"emailType"`
	Amount                *float64               `json:"amount"`
	PaymentID             *string                `json:"paymentId"`
	PaymentGateway        *string                `json:"paymentGateway"`
	Quantity              *int32                 `json:"quantity"`
	ComposedID            *string                `json:"composedId"`
	OrderNumber           *string                `json:"orderNumber"`
	InvoiceNumber         *string                `json:"invoiceNumber"`
	OversoldItems         []string               `json:"oversoldItems"`
	TransactionReference  *string                `json:"transactionReference"`
	ShippingCostsIncluded *bool                  `json:"shippingCostsIncluded"`

	event *model.OrderEvent

	// Lines                 []*OrderEventOrderLineObject `json:"lines"`
	// FulfilledItems        []*FulfillmentLine           `json:"fulfilledItems"`
	// Warehouse             *Warehouse                   `json:"warehouse"`
	// RelatedOrder          *Order                       `json:"relatedOrder"`
	// Discount              *OrderEventDiscountObject    `json:"discount"`
	// User                  *User                        `json:"user"`
}

func SystemOrderEventToGraphqlOrderEvent(o *model.OrderEvent) *OrderEvent {
	if o == nil {
		return nil
	}

	var email *string
	if em, ok := o.Parameters["email"]; ok && em != nil {
		email = model.NewPrimitive(em.(string))
	}

	var emailType OrderEventsEmailsEnum
	if et, ok := o.Parameters["email_type"]; ok && et != nil {
		emailType = OrderEventsEmailsEnum(strings.ToUpper(et.(string)))
	}

	var amount *float64
	if am, ok := o.Parameters["amount"]; ok && am != nil {
		amount = model.NewPrimitive(am.(float64))
	}

	var paymentID *string
	if pi, ok := o.Parameters["payment_id"]; ok && pi != nil {
		paymentID = model.NewPrimitive(pi.(string))
	}

	var paymentGateway *string
	if pg, ok := o.Parameters["payment_gateway"]; ok && pg != nil {
		paymentGateway = model.NewPrimitive(pg.(string))
	}

	var quantity *int32
	if qt, ok := o.Parameters["quantity"]; ok && qt != nil {
		quantity = model.NewPrimitive(int32(qt.(int)))
	}

	var message *string
	if msg, ok := o.Parameters["message"]; ok && msg != nil {
		message = model.NewPrimitive(msg.(string))
	}

	var composedID *string
	if cpID, ok := o.Parameters["composed_id"]; ok && cpID != nil {
		composedID = model.NewPrimitive(cpID.(string))
	}

	var overSoldItems []string
	item, ok := o.Parameters["oversold_items"]
	if ok && item != nil {
		overSoldItems = item.([]string)
	}

	var invoiceNumber *string
	if in, ok := o.Parameters["invoice_number"]; ok && in != nil {
		invoiceNumber = model.NewPrimitive(in.(string))
	}

	var transactionReference *string
	if tr, ok := o.Parameters["transaction_reference"]; ok && tr != nil {
		transactionReference = model.NewPrimitive(tr.(string))
	}

	var shippingCostsIncluded *bool
	if si, ok := o.Parameters["shipping_costs_included"]; ok && si != nil {
		shippingCostsIncluded = model.NewPrimitive(si.(bool))
	}

	var orderEventType = OrderEventsEnum(o.Type)

	res := &OrderEvent{
		ID:                    o.Id,
		Email:                 email,
		EmailType:             &emailType,
		Amount:                amount,
		PaymentID:             paymentID,
		PaymentGateway:        paymentGateway,
		Quantity:              quantity,
		Message:               message,
		ComposedID:            composedID,
		OversoldItems:         overSoldItems,
		OrderNumber:           &o.OrderID,
		InvoiceNumber:         invoiceNumber,
		TransactionReference:  transactionReference,
		ShippingCostsIncluded: shippingCostsIncluded,
		Type:                  &orderEventType,
		Date:                  &DateTime{util.TimeFromMillis(o.CreateAt)},

		event: o,
	}

	return res
}

func (o *OrderEvent) Discount(ctx context.Context) (*OrderEventDiscountObject, error) {
	discountObj := o.event.Parameters.Get("discount")
	if discountObj == nil {
		return nil, nil
	}

	obj, ok := discountObj.(model.StringInterface)
	if !ok {
		return nil, nil
	}

	currency := obj.Get("currency")
	if currency == nil {
		return nil, nil
	}

	panic("not implemented")
}

func (o *OrderEvent) RelatedOrder(ctx context.Context) (*Order, error) {
	orderID, ok := o.event.Parameters["related_order_pk"]
	if ok && orderID != nil {
		order, err := OrderByIdLoader.Load(ctx, orderID.(string))()
		if err != nil {
			return nil, err
		}

		return SystemOrderToGraphqlOrder(order), nil
	}

	return nil, nil
}

func (o *OrderEvent) Warehouse(ctx context.Context) (*Warehouse, error) {
	warehouseID, ok := o.event.Parameters["warehouse"]
	if ok && warehouseID != nil {
		warehouse, err := WarehouseByIdLoader.Load(ctx, warehouseID.(string))()
		if err != nil {
			return nil, err
		}

		return SystemWarehouseToGraphqlWarehouse(warehouse), nil
	}

	return nil, nil
}

func (o *OrderEvent) FulfilledItems(ctx context.Context) ([]*FulfillmentLine, error) {
	fulfillmentLineIDs, ok := o.event.Parameters["fulfilled_items"]
	if ok && fulfillmentLineIDs != nil {
		lines, errs := FulfillmentLinesByIdLoader.LoadMany(ctx, fulfillmentLineIDs.([]string))()
		if errs != nil && errs[0] != nil {
			return nil, errs[0]
		}

		return DataloaderResultMap(lines, SystemFulfillmentLineToGraphqlFulfillmentLine), nil
	}

	return nil, nil
}

func (o *OrderEvent) User(ctx context.Context) (*User, error) {
	if o.event.UserID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()

	if currentSession.UserId == *o.event.UserID ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionToAny(currentSession, model.PermissionManageUsers, model.PermissionManageStaff) {

		user, err := UserByUserIdLoader.Load(ctx, *o.event.UserID)()
		if err != nil {
			return nil, err
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, nil
}

func (o *OrderEvent) Lines(ctx context.Context) ([]*OrderEventOrderLineObject, error) {
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
