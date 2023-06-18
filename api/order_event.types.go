package api

import (
	"context"
	"strings"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

func getOrderDiscountEvent(discountObj model.StringInterface) *OrderEventDiscountObject {
	currency := discountObj.Get("currency", "").(string)
	var amount, oldAmount *goprices.Money

	amountValue, ok := discountObj.Get("amount_value", 0.0).(float64)
	if ok {
		amount, _ = goprices.NewMoney(amountValue, currency)
	}

	oldAmountValue, ok := discountObj.Get("old_amount_value", 0.0).(float64)
	if ok {
		oldAmount, _ = goprices.NewMoney(oldAmountValue, currency)
	}

	var (
		resValue        PositiveDecimal
		resValueType    DiscountValueTypeEnum
		resReason       *string
		resOldValue     *PositiveDecimal
		resOldValueType *DiscountValueTypeEnum
	)

	value, ok := discountObj.Get("value", 0.0).(float64)
	if ok {
		resValue = PositiveDecimal(decimal.NewFromFloat(value))
	}

	valueType, ok := discountObj.Get("value_type", "").(string)
	if ok {
		if vlType := DiscountValueTypeEnum(valueType); vlType.IsValid() {
			resValueType = vlType
		}
	}

	reason, ok := discountObj.Get("reason", "").(string)
	if ok && reason != "" {
		resReason = &reason
	}

	oldValue, ok := discountObj.Get("old_value", 0.0).(float64)
	if ok {
		decimalValue := PositiveDecimal(decimal.NewFromFloat(oldValue))
		resOldValue = &decimalValue
	}

	oldValueType, ok := discountObj.Get("old_value_type", "").(string)
	if ok {
		oldValueTypeEnum := DiscountValueTypeEnum(oldValueType)
		if oldValueTypeEnum.IsValid() {
			resOldValueType = &oldValueTypeEnum
		}
	}

	return &OrderEventDiscountObject{
		Amount:       SystemMoneyToGraphqlMoney(amount),
		OldAmount:    SystemMoneyToGraphqlMoney(oldAmount),
		Value:        resValue,
		ValueType:    resValueType,
		Reason:       resReason,
		OldValueType: resOldValueType,
		OldValue:     resOldValue,
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

	var emailType *OrderEventsEmailsEnum
	if et, ok := o.Parameters["email_type"]; ok && et != nil {
		mailType := OrderEventsEmailsEnum(strings.ToUpper(et.(string)))
		if mailType.IsValid() {
			emailType = &mailType
		}
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

	res := &OrderEvent{
		ID:                    o.Id,
		Email:                 email,
		EmailType:             emailType,
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
		Type:                  &o.Type,
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

	realDiscountObj, ok := discountObj.(model.StringInterface)
	if !ok {
		return nil, nil
	}

	return getOrderDiscountEvent(realDiscountObj), nil
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

		return systemRecordsToGraphql(lines, SystemFulfillmentLineToGraphqlFulfillmentLine), nil
	}

	return nil, nil
}

// Requester must be staff of shop to see.
//
// NOTE: Refer to ./schemas/order.graphqls for details on directive used.
func (o *OrderEvent) User(ctx context.Context) (*User, error) {
	if o.event.UserID == nil {
		return nil, nil
	}

	user, err := UserByUserIdLoader.Load(ctx, *o.event.UserID)()
	if err != nil {
		return nil, err
	}

	return SystemUserToGraphqlUser(user), nil
}

func (o *OrderEvent) Lines(ctx context.Context) ([]*OrderEventOrderLineObject, error) {
	rawLines := o.event.Parameters.Get("lines")
	if rawLines == nil {
		return nil, nil
	}
	lines, ok := rawLines.([]model.StringInterface)
	if !ok || len(lines) == 0 {
		return nil, nil
	}

	linePKs := []string{}
	for _, entry := range lines {
		linePK := entry.Get("line_pk", "")
		if linePK != nil {
			linePKs = append(linePKs, linePK.(string))
		}
	}

	orderLines, errs := OrderLineByIdLoader.LoadMany(ctx, linePKs)()
	if len(errs) > 0 && errs[0] != nil {
		return nil, errs[0]
	}

	orderLinesMap := lo.SliceToMap(orderLines, func(line *model.OrderLine) (string, *model.OrderLine) { return line.Id, line })

	res := []*OrderEventOrderLineObject{}

	for _, line := range lines {
		linePk := line.Get("line_pk", "").(string)
		discount, ok := line.Get("discount", model.StringInterface{}).(model.StringInterface)
		lineObject := orderLinesMap[linePk]

		if ok && discount != nil && len(discount) > 0 {
			discountObj := getOrderDiscountEvent(discount)
			quantity := line.Get("quantity", 0).(int)
			itemName := line.Get("item", "").(string)

			res = append(res, &OrderEventOrderLineObject{
				Quantity:  (*int32)(unsafe.Pointer(&quantity)),
				OrderLine: SystemOrderLineToGraphqlOrderLine(lineObject),
				ItemName:  &itemName,
				Discount:  discountObj,
			})
		}
	}

	return res, nil
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
