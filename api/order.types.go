package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
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

// ------------------------------- ORDER

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
	IsPaid              bool             `json:"isPaid"`
	Total               *TaxedMoney      `json:"total"`
	UndiscountedTotal   *TaxedMoney      `json:"undiscountedTotal"`
	StatusDisplay       *string          `json:"statusDisplay"`
	TotalCaptured       *Money           `json:"totalCaptured"`
	TotalBalance        *Money           `json:"totalBalance"`
	LanguageCodeEnum    LanguageCodeEnum `json:"languageCodeEnum"`

	channelID string

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

	res := &Order{
		ID: o.Id,

		channelID: o.ChannelID,
	}
	panic("not implemented")
	return res
}

func (o *Order) Discounts(ctx context.Context) ([]*OrderDiscount, error) {
	panic("not implemented")
}

func (o *Order) BillingAddress(ctx context.Context) (*Address, error) {
	panic("not implemented")
}

func (o *Order) ShippingAddress(ctx context.Context) (*Address, error) {
	panic("not implemented")
}

func (o *Order) Actions(ctx context.Context) ([]*OrderAction, error) {
	panic("not implemented")
}

func (o *Order) Subtotal(ctx context.Context) (*TaxedMoney, error) {
	panic("not implemented")
}

func (o *Order) TotalAuthorized(ctx context.Context) (*Money, error) {
	panic("not implemented")
}

func (o *Order) Fulfillments(ctx context.Context) ([]*Fulfillment, error) {
	panic("not implemented")
}

func (o *Order) Lines(ctx context.Context) ([]*OrderLine, error) {
	panic("not implemented")
}

func (o *Order) Events(ctx context.Context) ([]*OrderEvent, error) {
	panic("not implemented")
}

func (o *Order) PaymentStatus(ctx context.Context) (PaymentChargeStatusEnum, error) {
	panic("not implemented")
}

func (o *Order) PaymentStatusDisplay(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (o *Order) CanFinalize(ctx context.Context) (bool, error) {
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

// func (o *Order) Errors(ctx context.Context) (*string, error) {
// 	panic("not implemented")
// }

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
