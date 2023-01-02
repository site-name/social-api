package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// --------------------------- Order line -----------------------------

type OrderLine struct {
	ID                    string                 `json:"id"`
	ProductName           string                 `json:"productName"`
	VariantName           string                 `json:"variantName"`
	ProductSku            *string                `json:"productSku"`
	ProductVariantID      *string                `json:"ProductVariantId"`
	IsShippingRequired    bool                   `json:"isShippingRequired"`
	Quantity              int32                  `json:"quantity"`
	QuantityFulfilled     int32                  `json:"quantityFulfilled"`
	UnitDiscountReason    *string                `json:"unitDiscountReason"`
	TaxRate               float64                `json:"taxRate"`
	UnitPrice             *TaxedMoney            `json:"unitPrice"`
	UndiscountedUnitPrice *TaxedMoney            `json:"undiscountedUnitPrice"`
	UnitDiscount          *Money                 `json:"unitDiscount"`
	UnitDiscountValue     PositiveDecimal        `json:"unitDiscountValue"`
	TotalPrice            *TaxedMoney            `json:"totalPrice"`
	TranslatedProductName string                 `json:"translatedProductName"`
	TranslatedVariantName string                 `json:"translatedVariantName"`
	QuantityToFulfill     int32                  `json:"quantityToFulfill"`
	UnitDiscountType      *DiscountValueTypeEnum `json:"unitDiscountType"`

	variantID *string
	orderID   string

	// Thumbnail             *Image                 `json:"thumbnail"`
	// Allocations           []*Allocation          `json:"allocations"`
	// DigitalContentURL     *DigitalContentURL     `json:"digitalContentUrl"`
	// Variant               *ProductVariant        `json:"variant"`
}

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

		variantID: line.VariantID,
		orderID:   line.OrderID,
	}
	discountType := DiscountValueTypeEnum(strings.ToUpper(line.UnitDiscountType))
	res.UnitDiscountType = &discountType

	if line.TaxRate != nil {
		res.TaxRate, _ = line.TaxRate.Float64()
	}

	return res
}

func (o *OrderLine) Thumbnail(ctx context.Context) (*Image, error) {
	panic("not implemented")
}

func (o *OrderLine) Allocation(ctx context.Context) (*Allocation, error) {
	panic("not implemented")
}

func (o *OrderLine) Variant(ctx context.Context) (*ProductVariant, error) {
	if o.variantID == nil {
		return nil, nil
	}

	panic("not implemented")
}

func orderLineByIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*model.OrderLine] {
	var (
		res          = make([]*dataloader.Result[*model.OrderLine], len(orderLineIDs))
		appErr       *model.AppError
		orderLines   []*model.OrderLine
		orderLineMap = map[string]*model.OrderLine{} // keys are order line ids
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

	orderLineMap = lo.SliceToMap(orderLines, func(o *model.OrderLine) (string, *model.OrderLine) { return o.Id, o })

	for idx, id := range orderLineIDs {
		res[idx] = &dataloader.Result[*model.OrderLine]{Data: orderLineMap[id]}
	}
	return res

errorLabel:
	for idx := range orderLineIDs {
		res[idx] = &dataloader.Result[*model.OrderLine]{Error: err}
	}
	return res
}

func orderLinesByOrderIdLoader(ctx context.Context, orderIDs []string) []*dataloader.Result[[]*model.OrderLine] {
	var (
		res     = make([]*dataloader.Result[[]*model.OrderLine], len(orderIDs))
		lines   model.OrderLines
		appErr  *model.AppError
		lineMap = map[string][]*model.OrderLine{} // keys are order ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	lines, appErr = embedCtx.App.Srv().
		OrderService().
		OrderLinesByOption(&model.OrderLineFilterOption{
			OrderID: squirrel.Eq{store.OrderLineTableName + ".OrderID": orderIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range lines {
		lineMap[line.OrderID] = append(lineMap[line.OrderID], line)
	}

	for idx, id := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderLine]{Data: lineMap[id]}
	}
	return res

errorLabel:
	for idx := range orderIDs {
		res[idx] = &dataloader.Result[[]*model.OrderLine]{Error: err}
	}
	return res
}

// idPairs are strings with format variantID__channelID
func orderLinesByVariantIdAndChannelIdLoader(ctx context.Context, idPairs []string) []*dataloader.Result[[]*model.OrderLine] {
	var (
		res     = make([]*dataloader.Result[[]*model.OrderLine], len(idPairs))
		lines   model.OrderLines
		appErr  *model.AppError
		lineMap = map[string]model.OrderLines{} // keys have format variantID__channelID

		variantIDs []string
		channelIDs []string
	)

	for _, pair := range idPairs {
		index := strings.Index(pair, "__")
		if index < 0 {
			continue
		}

		variantIDs = append(variantIDs, pair[:index])
		channelIDs = append(channelIDs, pair[index+2:])
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	lines, appErr = embedCtx.App.Srv().
		OrderService().
		OrderLinesByOption(&model.OrderLineFilterOption{
			VariantID:          squirrel.Eq{store.OrderLineTableName + ".VariantID": variantIDs},
			OrderChannelID:     squirrel.Eq{store.OrderTableName + ".ChannelID": channelIDs},
			SelectRelatedOrder: true,
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, line := range lines {
		if line.VariantID == nil {
			continue
		}

		key := *line.VariantID + "__" + line.GetOrder().ChannelID
		lineMap[key] = append(lineMap[key], line)
	}

	for idx, key := range idPairs {
		res[idx] = &dataloader.Result[[]*model.OrderLine]{Data: lineMap[key]}
	}
	return res

errorLabel:
	for idx := range idPairs {
		res[idx] = &dataloader.Result[[]*model.OrderLine]{Error: err}
	}
	return res
}

// ------------------------------- ORDER ---------------------------------

type Order struct {
	ID                  string           `json:"id"`
	Created             DateTime         `json:"created"`
	Status              OrderStatus      `json:"status"`
	TrackingClientID    string           `json:"trackingClientId"`
	ShippingMethodName  *string          `json:"shippingMethodName"`
	CollectionPointName *string          `json:"collectionPointName"`
	ShippingPrice       *TaxedMoney      `json:"shippingPrice"`
	ShippingTaxRate     float64          `json:"shippingTaxRate"`
	Token               string           `json:"token"`
	DisplayGrossPrices  bool             `json:"displayGrossPrices"`
	CustomerNote        string           `json:"customerNote"`
	Weight              *Weight          `json:"weight"`
	RedirectURL         *string          `json:"redirectUrl"`
	PrivateMetadata     []*MetadataItem  `json:"privateMetadata"`
	Metadata            []*MetadataItem  `json:"metadata"`
	Number              *string          `json:"number"`
	Origin              OrderOriginEnum  `json:"origin"`
	Total               *TaxedMoney      `json:"total"`
	UndiscountedTotal   *TaxedMoney      `json:"undiscountedTotal"`
	TotalCaptured       *Money           `json:"totalCaptured"`
	TotalBalance        *Money           `json:"totalBalance"`
	LanguageCodeEnum    LanguageCodeEnum `json:"languageCodeEnum"`

	channelID         string
	userID            *string
	billingAddressID  *string
	shippingAddressID *string
	order             *model.Order // parent order

	// StatusDisplay       *string          `json:"statusDisplay"`
	// IsPaid              bool             `json:"isPaid"`
	// Original             *string                 `json:"original"`
	// IsShippingRequired   bool                    `json:"isShippingRequired"`
	// User                 *User                   `json:"user"`
	// UserEmail            *string                 `json:"userEmail"`
	// CanFinalize          bool                    `json:"canFinalize"`
	// PaymentStatusDisplay string                  `json:"paymentStatusDisplay"`
	// PaymentStatus        PaymentChargeStatusEnum `json:"paymentStatus"`
	// TotalAuthorized      *Money                  `json:"totalAuthorized"`
	// Subtotal             *TaxedMoney             `json:"subtotal"`
	// BillingAddress            *Address                `json:"billingAddress"`
	// ShippingAddress           *Address                `json:"shippingAddress"`
	// Channel                   *Channel                `json:"channel"`
	// Voucher                   *Voucher                `json:"voucher"`
	// GiftCards                 []*GiftCard             `json:"giftCards"`
	// Fulfillments              []*Fulfillment          `json:"fulfillments"`
	// Lines                     []*OrderLine            `json:"lines"`
	// Actions                   []*OrderAction          `json:"actions"`
	// AvailableShippingMethods  []*ShippingMethod       `json:"availableShippingMethods"`
	// AvailableCollectionPoints []*Warehouse            `json:"availableCollectionPoints"`
	// Invoices                  []*Invoice              `json:"invoices"`
	// Payments                  []*Payment              `json:"payments"`
	// Events                    []*OrderEvent           `json:"events"`
	// DeliveryMethod            DeliveryMethod          `json:"deliveryMethod"`
	// Discounts                 []*OrderDiscount        `json:"discounts"`
}

func SystemOrderToGraphqlOrder(o *model.Order) *Order {
	if o == nil {
		return nil
	}

	o.PopulateNonDbFields()

	res := &Order{
		ID:                  o.Id,
		Created:             DateTime{util.TimeFromMillis(o.CreateAt)},
		Status:              OrderStatus(o.Status),
		TrackingClientID:    o.TrackingClientID,
		ShippingMethodName:  o.ShippingMethodName,
		CollectionPointName: o.CollectionPointName,
		ShippingPrice:       SystemTaxedMoneyToGraphqlTaxedMoney(o.ShippingPrice),
		Token:               o.Token,
		DisplayGrossPrices:  *o.DisplayGrossPrices,
		CustomerNote:        o.CustomerNote,
		RedirectURL:         o.RedirectUrl,
		PrivateMetadata:     MetadataToSlice(o.PrivateMetadata),
		Metadata:            MetadataToSlice(o.Metadata),
		Number:              &o.Id,
		Origin:              OrderOriginEnum(o.Origin),
		Total:               SystemTaxedMoneyToGraphqlTaxedMoney(o.Total),
		UndiscountedTotal:   SystemTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedTotal),
		TotalCaptured:       SystemMoneyToGraphqlMoney(o.TotalPaid),
		TotalBalance:        SystemMoneyToGraphqlMoney(o.TotalBalance()),
		LanguageCodeEnum:    SystemLanguageToGraphqlLanguageCodeEnum(o.LanguageCode),
		Weight: &Weight{
			Value: float64(o.WeightAmount),
			Unit:  WeightUnitsEnum(o.WeightUnit),
		},

		channelID:         o.ChannelID,
		userID:            o.UserID,
		billingAddressID:  o.BillingAddressID,
		shippingAddressID: o.ShippingAddressID,
		order:             o,
	}

	if o.ShippingTaxRate != nil {
		fl64, _ := o.ShippingTaxRate.Float64()
		res.ShippingTaxRate = fl64
	}

	return res
}

func (o *Order) Discounts(ctx context.Context) ([]*OrderDiscount, error) {
	rels, err := dataloaders.OrderDiscountsByOrderIDLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(rels, SystemOrderDiscountToGraphqlOrderDiscount), nil
}

func (o *Order) IsPaid(ctx context.Context) (bool, error) {
	return o.order.IsFullyPaid(), nil
}

func (o *Order) StatusDisplay(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (o *Order) BillingAddress(ctx context.Context) (*Address, error) {
	if o.billingAddressID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	var (
		currentSession = embedCtx.AppContext.Session()
		accountSrv     = embedCtx.App.Srv().AccountService()
	)

	if o.userID != nil {
		user, err := dataloaders.UserByUserIdLoader.Load(ctx, *o.userID)()
		if err != nil {
			return nil, err
		}

		address, err := dataloaders.AddressByIdLoader.Load(ctx, *o.billingAddressID)()
		if err != nil {
			return nil, err
		}

		if currentSession.UserId == user.Id || accountSrv.SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
			return SystemAddressToGraphqlAddress(address), nil
		}

		return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
	}

	// else case

	address, err := dataloaders.AddressByIdLoader.Load(ctx, *o.billingAddressID)()
	if err != nil {
		return nil, err
	}

	if accountSrv.SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		return SystemAddressToGraphqlAddress(address), nil
	}

	return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
}

func (o *Order) ShippingAddress(ctx context.Context) (*Address, error) {
	if o.shippingAddressID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	var (
		currentSession = embedCtx.AppContext.Session()
		accountSrv     = embedCtx.App.Srv().AccountService()
	)

	if o.userID != nil {
		user, err := dataloaders.UserByUserIdLoader.Load(ctx, *o.userID)()
		if err != nil {
			return nil, err
		}

		address, err := dataloaders.AddressByIdLoader.Load(ctx, *o.shippingAddressID)()
		if err != nil {
			return nil, err
		}

		if currentSession.UserId == user.Id || accountSrv.SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
			return SystemAddressToGraphqlAddress(address), nil
		}

		return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
	}

	// else case

	address, err := dataloaders.AddressByIdLoader.Load(ctx, *o.shippingAddressID)()
	if err != nil {
		return nil, err
	}

	if accountSrv.SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		return SystemAddressToGraphqlAddress(address), nil
	}

	return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
}

func (o *Order) Actions(ctx context.Context) ([]*OrderAction, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	orderSrv := embedCtx.App.Srv().OrderService()

	payments, err := dataloaders.PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	actions := []*OrderAction{}
	lastPayment := embedCtx.App.Srv().PaymentService().GetLastpayment(payments)

	ok, appErr := orderSrv.OrderCanCapture(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		ptr := OrderActionCapture
		actions = append(actions, &ptr)
	}

	ok, appErr = orderSrv.CanMarkOrderAsPaid(o.order, payments)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		ptr := OrderActionMarkAsPaid
		actions = append(actions, &ptr)
	}

	ok, appErr = orderSrv.OrderCanRefund(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		ptr := OrderActionRefund
		actions = append(actions, &ptr)
	}

	ok, appErr = orderSrv.OrderCanVoid(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		ptr := OrderActionVoid
		actions = append(actions, &ptr)
	}

	return actions, nil
}

func (o *Order) Subtotal(ctx context.Context) (*TaxedMoney, error) {
	lines, err := dataloaders.OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	subTotal, appErr := embedCtx.App.Srv().PaymentService().GetSubTotal(lines, o.order.Currency)
	if appErr != nil {
		return nil, appErr
	}

	return SystemTaxedMoneyToGraphqlTaxedMoney(subTotal), nil
}

func (o *Order) Payments(ctx context.Context) ([]*Payment, error) {
	payments, err := dataloaders.PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(payments, SystemPaymentToGraphqlPayment), nil
}

func (o *Order) TotalAuthorized(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (o *Order) Fulfillments(ctx context.Context) ([]*Fulfillment, error) {
	panic("not implemented")
}

func (o *Order) Lines(ctx context.Context) ([]*OrderLine, error) {
	lines, err := dataloaders.OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(lines, SystemOrderLineToGraphqlOrderLine), nil
}

func (o *Order) Events(ctx context.Context) ([]*OrderEvent, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	// check if current user has manage order permission to see order events:
	if !embedCtx.App.Srv().AccountService().SessionHasPermissionTo(embedCtx.AppContext.Session(), model.PermissionManageOrders) {
		return nil, model.NewAppError("Order,Events", ErrorUnauthorized, nil, "you are not authorized to see order events", http.StatusUnauthorized)
	}

	events, err := dataloaders.OrderEventsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(events, SystemOrderEventToGraphqlOrderEvent), nil
}

func (o *Order) PaymentStatus(ctx context.Context) (*PaymentChargeStatusEnum, error) {
	payments, err := dataloaders.PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		notCharged := PaymentChargeStatusEnumNotCharged
		return &notCharged, nil
	}

	if len(payments) == 1 {
		status := PaymentChargeStatusEnum(payments[0].ChargeStatus)
		return &status, nil
	}

	// find latest payment
	lastPayment := payments[0]
	for _, pm := range payments {
		if pm != nil && pm.CreateAt > lastPayment.CreateAt {
			lastPayment = pm
		}
	}

	status := PaymentChargeStatusEnum(lastPayment.ChargeStatus)
	return &status, nil
}

func (o *Order) PaymentStatusDisplay(ctx context.Context) (string, error) {
	payments, err := dataloaders.PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return "", err
	}

	if len(payments) == 0 {
		return model.ChargeStatuString[model.NOT_CHARGED], nil
	}

	if len(payments) == 1 {
		return model.ChargeStatuString[payments[0].ChargeStatus], nil
	}

	// find latest payment
	lastPayment := payments[0]
	for _, pm := range payments {
		if pm != nil && pm.CreateAt > lastPayment.CreateAt {
			lastPayment = pm
		}
	}

	return model.ChargeStatuString[lastPayment.ChargeStatus], nil
}

func (o *Order) CanFinalize(ctx context.Context) (bool, error) {
	// if o.Status == OrderStatusDraft {
	// 	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	// 	if err != nil {
	// 		return false, err
	// 	}

	// 	country, appErr := embedCtx.App.Srv().OrderService().GetOrderCountry(o.order)
	// 	if appErr != nil {
	// 		return false, appErr
	// 	}
	// }

	// return true, nil
	panic("not implemented")
}

func (o *Order) UserEmail(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (o *Order) User(ctx context.Context) (*User, error) {
	panic("not implemented")
}

func (o *Order) DeliveryMethod(ctx context.Context) (DeliveryMethod, error) {
	panic("not implemented")
}

func (o *Order) AvailableShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	panic("not implemented")
}

func (o *Order) AvailableCollectionPoints(ctx context.Context) ([]*Warehouse, error) {
	panic("not implemented")
}

func (o *Order) Invoices(ctx context.Context) ([]*Invoice, error) {
	panic("not implemented")
}

func (o *Order) IsShippingRequired(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (o *Order) GiftCards(ctx context.Context) ([]*GiftCard, error) {
	panic("not implemented")
}

func (o *Order) Voucher(ctx context.Context) (*Voucher, error) {
	panic("not implemented")
}

func (o *Order) Original(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func orderByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*model.Order] {
	var (
		res      = make([]*dataloader.Result[*model.Order], len(ids))
		orders   model.Orders
		appErr   *model.AppError
		orderMap = map[string]*model.Order{}
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	orders, appErr = embedCtx.App.Srv().
		OrderService().
		FilterOrdersByOptions(&model.OrderFilterOption{
			Id: squirrel.Eq{store.OrderTableName + ".Id": ids},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	orderMap = lo.SliceToMap(orders, func(o *model.Order) (string, *model.Order) { return o.Id, o })

	for idx, id := range ids {
		res[idx] = &dataloader.Result[*model.Order]{Data: orderMap[id]}
	}
	return res

errorLabel:
	for idx := range ids {
		res[idx] = &dataloader.Result[*model.Order]{Error: err}
	}
	return res
}

func ordersByUserLoader(ctx context.Context, userIDs []string) []*dataloader.Result[[]*model.Order] {
	var (
		res      = make([]*dataloader.Result[[]*model.Order], len(userIDs))
		appErr   *model.AppError
		orders   model.Orders
		orderMap = map[string]model.Orders{} // keys are user ids
	)

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		goto errorLabel
	}

	orders, appErr = embedCtx.App.Srv().
		OrderService().
		FilterOrdersByOptions(&model.OrderFilterOption{
			UserID: squirrel.Eq{store.OrderTableName + ".UserID": userIDs},
		})
	if appErr != nil {
		err = appErr
		goto errorLabel
	}

	for _, ord := range orders {
		orderMap[*ord.UserID] = append(orderMap[*ord.UserID], ord)
	}

	for idx, id := range userIDs {
		res[idx] = &dataloader.Result[[]*model.Order]{Data: orderMap[id]}
	}
	return res

errorLabel:
	for idx := range userIDs {
		res[idx] = &dataloader.Result[[]*model.Order]{Error: err}
	}
	return res
}
