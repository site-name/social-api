package api

import (
	"context"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func getOrderDiscountEvent(discountObj model.StringInterface) *OrderEventDiscountObject {
	currency := discountObj.Get("currency", "")
	if currency == nil || currency == "" {
		slog.Error("getOrderDiscountEvent: missing value", slog.String("field", "currency"))
		return nil
	}

	var amount, oldAmount *goprices.Money

	amountValue := discountObj.Get("amount_value")
	if amountValue == nil {
		slog.Error("getOrderDiscountEvent: missing value", slog.String("field", "amount_value"))
		return nil
	}
	switch t := amountValue.(type) {
	case float64:
		amount = &goprices.Money{decimal.NewFromFloat(t), currency.(string)}
	case decimal.Decimal:
		amount = &goprices.Money{t, currency.(string)}
	}

	oldAmountValue := discountObj.Get("old_amount_value")
	if oldAmountValue != nil {
		switch t := oldAmountValue.(type) {
		case float64:
			oldAmount = &goprices.Money{decimal.NewFromFloat(t), currency.(string)}
		case decimal.Decimal:
			oldAmount = &goprices.Money{t, currency.(string)}
		}
	}

	var (
		resValue        PositiveDecimal
		resValueType    DiscountValueTypeEnum
		resReason       *string
		resOldValue     *PositiveDecimal
		resOldValueType *DiscountValueTypeEnum
	)

	value := discountObj.Get("value")
	if value == nil {
		slog.Error("getOrderDiscountEvent: missing value", slog.String("field", "value"))
		return nil
	}
	switch t := value.(type) {
	case float64:
		resValue = PositiveDecimal(decimal.NewFromFloat(t))
	case decimal.Decimal:
		resValue = PositiveDecimal(t)
	}

	valueType := discountObj.Get("value_type")
	if valueType == nil {
		slog.Error("getOrderDiscountEvent: missing value", slog.String("field", "value_type"))
		return nil
	}
	if strValueType, ok := valueType.(string); ok {
		resValueType = DiscountValueTypeEnum(strValueType)
	}

	reason := discountObj.Get("reason")
	if reason != nil && reason != "" {
		resReason = (*string)(unsafe.Pointer(&reason))
	}

	return &OrderEventDiscountObject{
		Amount: amount,
	}
}

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
	// requester must be staff of shop which has an order contains this event
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.SessionRequired()
	if embedCtx.Err != nil {
		return nil, embedCtx.Err
	}

	if embedCtx.AppContext.Session().GetUserRoles().Contains(model.ShopStaffRoleId) {
		if o.event.UserID == nil {
			return nil, nil
		}

		user, err := UserByUserIdLoader.Load(ctx, *o.event.UserID)()
		if err != nil {
			return nil, err
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, MakeUnauthorizedError("OrderEvent.User")
}

func (o *OrderEvent) Lines(ctx context.Context) ([]*OrderEventOrderLineObject, error) {
	panic("not implemented")
	// rawLines := o.event.Parameters.Get("lines", []map[string]any{})
	// if rawLines == nil {
	// 	return nil, nil
	// }
	// lines, ok := rawLines.([]map[string]any)
	// if ok && len(lines) == 0 {
	// 	return nil, nil
	// }

	// linePKs := []string{}
	// for _, entry := range lines {
	// 	linePK := entry["line_pk"]
	// 	if linePK != nil {
	// 		strLinePk, ok := linePK.(string)
	// 		if ok {
	// 			linePKs = append(linePKs, strLinePk)
	// 		}
	// 	}
	// }

	// orderLines, errs := OrderLineByIdLoader.LoadMany(ctx, linePKs)()
	// if len(errs) > 0 && errs[0] != nil {
	// 	return nil, errs[0]
	// }

	// orderLinesMap := lo.SliceToMap(orderLines, func(line *model.OrderLine) (string, *model.OrderLine) { return line.Id, line })
}

func orderEventsByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.OrderEvent] {
	var (
		res      = make([]*dataloader.Result[[]*model.OrderEvent], len(orderIDs))
		eventMap = map[string][]*model.OrderEvent{}
	)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	events, appErr := embedCtx.App.Srv().OrderService().FilterOrderEventsByOptions(&model.OrderEventFilterOptions{
		OrderID: squirrel.Eq{store.OrderEventTableName + ".OrderID": orderIDs},
	})
	if appErr != nil {
		for idx := range orderIDs {
			res[idx] = &dataloader.Result[[]*model.OrderEvent]{Error: appErr}
		}
		return res
	}

	for _, event := range events {
		eventMap[event.OrderID] = append(eventMap[event.OrderID], event)
	}
	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderEvent]{Data: eventMap[id]}
	}
	return res
}
