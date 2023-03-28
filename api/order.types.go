package api

import (
	"context"
	"encoding/base64"
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
	discountType := DiscountValueTypeEnum(line.UnitDiscountType)
	res.UnitDiscountType = &discountType

	if line.TaxRate != nil {
		res.TaxRate, _ = line.TaxRate.Float64()
	}

	return res
}

func (o *OrderLine) Thumbnail(ctx context.Context, args struct{ Size *int32 }) (*Image, error) {
	panic("not implemented")
}

func (o *OrderLine) DigitalContentURL(ctx context.Context) (*DigitalContentURL, error) {
	url, err := DigitalContentUrlByOrderLineID.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}
	return systemDigitalContentURLToGraphqlDigitalContentURL(url), nil
}

func (o *OrderLine) Allocations(ctx context.Context) ([]*Allocation, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()

	if embedCtx.App.Srv().AccountService().SessionHasPermissionToAny(currentSession, model.PermissionManageProducts, model.PermissionManageOrders) {
		allocations, err := AllocationsByOrderLineIdLoader.Load(ctx, o.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(allocations, systemAllocationToGraphqlAllocation), nil
	}

	return nil, model.NewAppError("OrderLine.Allocations", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func (o *OrderLine) Variant(ctx context.Context) (*ProductVariant, error) {
	if o.variantID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}
	currentSession := embedCtx.AppContext.Session()

	variant, err := ProductVariantByIdLoader.Load(ctx, *o.variantID)()
	if err != nil {
		return nil, err
	}

	if embedCtx.App.Srv().
		AccountService().
		SessionHasPermissionToAny(currentSession, model.PermissionManageOrders, model.PermissionManageDiscounts, model.PermissionManageProducts) {
		return SystemProductVariantToGraphqlProductVariant(variant), nil
	}

	channel, err := ChannelByOrderLineIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	productChannelListing, err := ProductChannelListingByProductIdAndChannelSlugLoader.Load(ctx, variant.ProductID+"__"+channel.Id)()
	if err != nil {
		return nil, err
	}

	if productChannelListing.IsVisible() {
		return SystemProductVariantToGraphqlProductVariant(variant), nil
	}

	return nil, nil
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
	ID                  string            `json:"id"`
	Created             DateTime          `json:"created"`
	Status              model.OrderStatus `json:"status"`
	TrackingClientID    string            `json:"trackingClientId"`
	ShippingMethodName  *string           `json:"shippingMethodName"`
	CollectionPointName *string           `json:"collectionPointName"`
	ShippingPrice       *TaxedMoney       `json:"shippingPrice"`
	ShippingTaxRate     float64           `json:"shippingTaxRate"`
	Token               string            `json:"token"`
	DisplayGrossPrices  bool              `json:"displayGrossPrices"`
	CustomerNote        string            `json:"customerNote"`
	Weight              *Weight           `json:"weight"`
	RedirectURL         *string           `json:"redirectUrl"`
	PrivateMetadata     []*MetadataItem   `json:"privateMetadata"`
	Metadata            []*MetadataItem   `json:"metadata"`
	Number              *string           `json:"number"`
	Origin              model.OrderOrigin `json:"origin"`
	Total               *TaxedMoney       `json:"total"`
	UndiscountedTotal   *TaxedMoney       `json:"undiscountedTotal"`
	TotalCaptured       *Money            `json:"totalCaptured"`
	TotalBalance        *Money            `json:"totalBalance"`
	LanguageCodeEnum    LanguageCodeEnum  `json:"languageCodeEnum"`

	order *model.Order // real order

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
		Status:              o.Status,
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
		Origin:              o.Origin,
		Total:               SystemTaxedMoneyToGraphqlTaxedMoney(o.Total),
		UndiscountedTotal:   SystemTaxedMoneyToGraphqlTaxedMoney(o.UnDiscountedTotal),
		TotalCaptured:       SystemMoneyToGraphqlMoney(o.TotalPaid),
		TotalBalance:        SystemMoneyToGraphqlMoney(o.TotalBalance()),
		LanguageCodeEnum:    o.LanguageCode,
		Weight: &Weight{
			Value: float64(o.WeightAmount),
			Unit:  WeightUnitsEnum(o.WeightUnit),
		},

		order: o,
	}

	if o.ShippingTaxRate != nil {
		res.ShippingTaxRate, _ = o.ShippingTaxRate.Float64()
	}

	return res
}

func (o *Order) Discounts(ctx context.Context) ([]*OrderDiscount, error) {
	rels, err := OrderDiscountsByOrderIDLoader.Load(ctx, o.ID)()
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
	if o.order.BillingAddressID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	address, err := AddressByIdLoader.Load(ctx, *o.order.BillingAddressID)()
	if err != nil {
		return nil, err
	}

	var currentSession = embedCtx.AppContext.Session()

	if (o.order.UserID != nil && *o.order.UserID == currentSession.UserId) ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		return SystemAddressToGraphqlAddress(address), nil
	}

	return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
}

func (o *Order) ShippingAddress(ctx context.Context) (*Address, error) {
	if o.order.ShippingAddressID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	address, err := AddressByIdLoader.Load(ctx, *o.order.ShippingAddressID)()
	if err != nil {
		return nil, err
	}

	var currentSession = embedCtx.AppContext.Session()

	if (o.order.UserID != nil && *o.order.UserID == currentSession.UserId) ||
		embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		return SystemAddressToGraphqlAddress(address), nil
	}

	return SystemAddressToGraphqlAddress(address.Obfuscate()), nil
}

func (o *Order) Actions(ctx context.Context) ([]OrderAction, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	orderSrv := embedCtx.App.Srv().OrderService()

	payments, err := PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	actions := []OrderAction{}
	lastPayment := embedCtx.App.Srv().PaymentService().GetLastpayment(payments)

	ok, appErr := orderSrv.OrderCanCapture(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		actions = append(actions, OrderActionCapture)
	}

	ok, appErr = orderSrv.CanMarkOrderAsPaid(o.order, payments)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		actions = append(actions, OrderActionMarkAsPaid)
	}

	ok, appErr = orderSrv.OrderCanRefund(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		actions = append(actions, OrderActionRefund)
	}

	ok, appErr = orderSrv.OrderCanVoid(o.order, lastPayment)
	if appErr != nil {
		return nil, appErr
	}
	if ok {
		actions = append(actions, OrderActionVoid)
	}

	return actions, nil
}

func (o *Order) Subtotal(ctx context.Context) (*TaxedMoney, error) {
	lines, err := OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
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
	payments, err := PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(payments, SystemPaymentToGraphqlPayment), nil
}

func (o *Order) TotalAuthorized(ctx context.Context) (*Money, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	payments, err := PaymentsByOrderIdLoader.Load(ctx, o.order.Id)()
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return &Money{
			Amount:   0,
			Currency: o.order.Currency,
		}, nil
	}

	// find most recent payment:
	var mostRecentPayment = payments[0]
	if len(payments) > 1 {
		for _, pm := range payments {
			if pm != nil && pm.CreateAt > mostRecentPayment.CreateAt {
				mostRecentPayment = pm
			}
		}
	}
	if !*mostRecentPayment.IsActive {
		return &Money{
			Amount:   0,
			Currency: o.order.Currency,
		}, nil
	}

	money, appErr := embedCtx.App.Srv().PaymentService().PaymentGetAuthorizedAmount(mostRecentPayment)
	if appErr != nil {
		return nil, appErr
	}
	return SystemMoneyToGraphqlMoney(money), nil
}

func (o *Order) Fulfillments(ctx context.Context) ([]*Fulfillment, error) {
	fulfillments, err := FulfillmentsByOrderIdLoader.Load(ctx, o.order.Id)()
	if err != nil {
		return nil, err
	}
	/*
		TODO: https://github.com/site-name/social-api/issues/11
	*/

	return DataloaderResultMap(fulfillments, SystemFulfillmentToGraphqlFulfillment), nil
}

func (o *Order) Lines(ctx context.Context) ([]*OrderLine, error) {
	lines, err := OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
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
		return nil, model.NewAppError("Order.Events", ErrorUnauthorized, nil, "you are not authorized to see order events", http.StatusUnauthorized)
	}

	events, err := OrderEventsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(events, SystemOrderEventToGraphqlOrderEvent), nil
}

func (o *Order) PaymentStatus(ctx context.Context) (*PaymentChargeStatusEnum, error) {
	payments, err := PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		notCharged := PaymentChargeStatusEnumNotCharged
		return &notCharged, nil
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
	payments, err := PaymentsByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return "", err
	}

	if len(payments) == 0 {
		return model.ChargeStatuString[model.NOT_CHARGED], nil
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
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	var currentSession = embedCtx.AppContext.Session()

	if (o.order.UserID != nil && *o.order.UserID == currentSession.UserId) ||
		embedCtx.App.Srv().
			AccountService().
			SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {

		return &o.order.UserEmail, nil
	}

	return model.NewPrimitive(util.ObfuscateEmail(o.order.UserEmail)), nil
}

func (o *Order) User(ctx context.Context) (*User, error) {
	if o.order.UserID == nil {
		return nil, nil
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()

	if (o.order.UserID != nil && currentSession.UserId == *o.order.UserID) ||
		embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageUsers) {
		user, err := UserByUserIdLoader.Load(ctx, *o.order.UserID)()
		if err != nil {
			return nil, err
		}

		return SystemUserToGraphqlUser(user), nil
	}

	return nil, model.NewAppError("Order.User", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func (o *Order) DeliveryMethod(ctx context.Context) (DeliveryMethod, error) {
	if o.order.ShippingMethodID != nil {
		shippingMethod, err := ShippingMethodByIdLoader.Load(ctx, *o.order.ShippingMethodID)()
		if err != nil {
			return nil, err
		}
		return SystemShippingMethodToGraphqlShippingMethod(shippingMethod), nil
	}

	if o.order.CollectionPointID != nil {
		warehouse, err := WarehouseByIdLoader.Load(ctx, *o.order.CollectionPointID)()
		if err != nil {
			return nil, err
		}

		return SystemWarehouseToGraphqlWarehouse(warehouse), nil
	}

	return nil, nil
}

func (o *Order) AvailableShippingMethods(ctx context.Context) ([]*ShippingMethod, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	methods, appErr := embedCtx.App.Srv().OrderService().GetValidShippingMethodsForOrder(o.order)
	if appErr != nil {
		return nil, appErr
	}

	if len(methods) == 0 {
		return []*ShippingMethod{}, nil
	}

	// TODO: complete plugin manager
	panic("not implemented")
}

func (o *Order) Channel(ctx context.Context) (*Channel, error) {
	channel, err := ChannelByIdLoader.Load(ctx, o.order.ChannelID)()
	if err != nil {
		return nil, err
	}
	return SystemChannelToGraphqlChannel(channel), nil
}

func (o *Order) AvailableCollectionPoints(ctx context.Context) ([]*Warehouse, error) {
	lines, err := OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	address, err := o.ShippingAddress(ctx)
	if err != nil {
		return nil, err
	}

	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	warehouses, appErr := embedCtx.App.Srv().OrderService().GetValidCollectionPointsForOrder(lines, address.Address.Country)
	if appErr != nil {
		return nil, appErr
	}

	return DataloaderResultMap(warehouses, SystemWarehouseToGraphqlWarehouse), nil
}

func (o *Order) Invoices(ctx context.Context) ([]*Invoice, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()

	if (o.order.UserID != nil && *o.order.UserID == currentSession.UserId) ||
		embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageOrders) {
		invoices, err := InvoicesByOrderIDLoader.Load(ctx, o.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(invoices, SystemInvoiceToGraphqlInvoice), nil
	}

	return nil, model.NewAppError("Order.Invoice", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}

func (o *Order) IsShippingRequired(ctx context.Context) (bool, error) {
	lines, err := OrderLinesByOrderIdLoader.Load(ctx, o.ID)()
	if err != nil {
		return false, err
	}

	return lo.SomeBy(lines, func(o *model.OrderLine) bool { return o.IsShippingRequired }), nil
}

func (o *Order) GiftCards(ctx context.Context) ([]*GiftCard, error) {
	giftcards, err := GiftcardsByOrderIDsLoader.Load(ctx, o.ID)()
	if err != nil {
		return nil, err
	}

	return DataloaderResultMap(giftcards, SystemGiftcardToGraphqlGiftcard), nil
}

func (o *Order) Voucher(ctx context.Context) (*Voucher, error) {
	if o.order.VoucherID == nil {
		return nil, nil
	}

	voucher, err := VoucherByIDLoader.Load(ctx, *o.order.VoucherID)()
	if err != nil {
		return nil, err
	}

	return systemVoucherToGraphqlVoucher(voucher), nil
}

func (o *Order) Original(ctx context.Context) (*string, error) {
	if o.order.OriginalID != nil {
		return nil, nil
	}
	value := append([]byte("Order"), *o.order.OriginalID...)

	return model.NewPrimitive(base64.StdEncoding.EncodeToString(value)), nil
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
		if ord.UserID == nil {
			continue
		}
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
