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
	Thumbnail             *Image                 `json:"thumbnail"`
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

func (o *OrderLine) Variant(ctx context.Context) (*ProductVariant, error) {
	if o.variantID == nil {
		return nil, nil
	}

	panic("not implemented")
}

func orderLineByIdLoader(ctx context.Context, orderLineIDs []string) []*dataloader.Result[*OrderLine] {
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

// ------------------------------- ORDER

type Order struct {
	ID                   string                  `json:"id"`
	Created              DateTime                `json:"created"`
	Status               OrderStatus             `json:"status"`
	User                 *User                   `json:"user"`
	TrackingClientID     string                  `json:"trackingClientId"`
	ShippingMethodName   *string                 `json:"shippingMethodName"`
	CollectionPointName  *string                 `json:"collectionPointName"`
	ShippingPrice        *TaxedMoney             `json:"shippingPrice"`
	ShippingTaxRate      float64                 `json:"shippingTaxRate"`
	Token                string                  `json:"token"`
	DisplayGrossPrices   bool                    `json:"displayGrossPrices"`
	CustomerNote         string                  `json:"customerNote"`
	Weight               *Weight                 `json:"weight"`
	RedirectURL          *string                 `json:"redirectUrl"`
	PrivateMetadata      []*MetadataItem         `json:"privateMetadata"`
	Metadata             []*MetadataItem         `json:"metadata"`
	Number               *string                 `json:"number"`
	Original             *string                 `json:"original"`
	Origin               OrderOriginEnum         `json:"origin"`
	IsPaid               bool                    `json:"isPaid"`
	PaymentStatus        PaymentChargeStatusEnum `json:"paymentStatus"`
	PaymentStatusDisplay string                  `json:"paymentStatusDisplay"`
	Total                *TaxedMoney             `json:"total"`
	UndiscountedTotal    *TaxedMoney             `json:"undiscountedTotal"`
	Subtotal             *TaxedMoney             `json:"subtotal"`
	StatusDisplay        *string                 `json:"statusDisplay"`
	CanFinalize          bool                    `json:"canFinalize"`
	TotalAuthorized      *Money                  `json:"totalAuthorized"`
	TotalCaptured        *Money                  `json:"totalCaptured"`
	TotalBalance         *Money                  `json:"totalBalance"`
	UserEmail            *string                 `json:"userEmail"`
	IsShippingRequired   bool                    `json:"isShippingRequired"`
	LanguageCodeEnum     LanguageCodeEnum        `json:"languageCodeEnum"`

	channelID string

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

func orderByIdLoader(ctx context.Context, ids []string) []*dataloader.Result[*Order] {
	panic("not implemented")
}
