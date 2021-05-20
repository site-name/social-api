package order

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/slog"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// max lengths for some fields of order model
const (
	ORDER_STATUS_MAX_LENGTH               = 32
	ORDER_TRACKING_CLIENT_ID_MAX_LENGTH   = 36
	ORDER_ORIGIN_MAX_LENGTH               = 32
	ORDER_SHIPPING_METHOD_NAME_MAX_LENGTH = 255
	ORDER_TOKEN_MAX_LENGTH                = 36
	ORDER_CHECKOUT_TOKEN_MAX_LENGTH       = 36
)

// fulfillment statuses
const (
	FULFILLMENT_FULFILLED             = "fulfilled"             // group of products in an order marked as fulfilled
	FULFILLMENT_REFUNDED              = "refunded"              // group of refunded products
	FULFILLMENT_RETURNED              = "returned"              // group of returned products
	FULFILLMENT_REFUNDED_AND_RETURNED = "refunded_and_returned" // group of returned and replaced products
	FULFILLMENT_REPLACED              = "replaced"              // group of replaced products
	FULFILLMENT_CANCELED              = "canceled"              // fulfilled group of products in an order marked as canceled
)

var FulfillmentStrings = map[string]string{
	FULFILLMENT_FULFILLED:             "Fulfilled",
	FULFILLMENT_REFUNDED:              "Refunded",
	FULFILLMENT_RETURNED:              "Returned",
	FULFILLMENT_REPLACED:              "Replaced",
	FULFILLMENT_REFUNDED_AND_RETURNED: "Refunded and returned",
	FULFILLMENT_CANCELED:              "Canceled",
}

// order origin valid values
const (
	CHECKOUT = "checkout" // order created from checkout
	DRAFT    = "draft"    // order created from draft order
	REISSUE  = "reissue"  // order created from reissue existing one
)

// ORDER STATUS VALID VALUES
const (
	STATUS_DRAFT        = "draft"               // fully editable, not finalized order created by staff users
	UNCONFIRMED         = "unconfirmed"         // order created by customers when confirmation is required
	UNFULFILLED         = "unfulfilled"         // order with no items marked as fulfilled
	PARTIALLY_FULFILLED = "partially fulfilled" // order with some items marked as fulfilled
	FULFILLED           = "fulfilled"           // order with all items marked as fulfilled
	PARTIALLY_RETURNED  = "partially_returned"  // order with some items marked as returned
	RETURNED            = "returned"            // order with all items marked as returned
	CANCELED            = "canceled"            // permanently canceled order
)

var OrderStatusStrings = map[string]string{
	STATUS_DRAFT:        "Draft",
	UNCONFIRMED:         "Unconfirmed",
	UNFULFILLED:         "Unfulfilled",
	PARTIALLY_FULFILLED: "Partially fulfilled",
	PARTIALLY_RETURNED:  "Partially returned",
	RETURNED:            "Returned",
	FULFILLED:           "Fulfilled",
	CANCELED:            "Canceled",
}

var OrderOriginStrings = map[string]string{
	CHECKOUT: "Checkout",
	DRAFT:    "Draft",
	REISSUE:  "Reissue",
}

type Order struct {
	Id                           string                 `json:"id"`
	CreateAt                     int64                  `json:"create_at"` // default: now()
	Status                       string                 `json:"status"`    // default: UNFULFILLED
	UserID                       *string                `json:"user_id"`   // default: en
	LanguageCode                 string                 `json:"language_code"`
	TrackingClientID             string                 `json:"tracking_client_id"`
	BillingAddressID             *string                `json:"billing_address_id"`
	ShippingAddressID            *string                `json:"shipping_address_id"`
	UserEmail                    string                 `json:"user_email"`  // default: ""
	OriginalID                   *string                `json:"original_id"` // original order
	Origin                       string                 `json:"origin"`
	Currency                     string                 `json:"currency"`
	ShippingMethodID             *string                `json:"shipping_method_id"`
	ShippingMethodName           *string                `json:"shipping_method_name"`
	ChannelID                    string                 `json:"channel_id"`
	ShippingPriceNetAmount       *decimal.Decimal       `json:"shipping_price_net_amount"`
	ShippingPriceNet             *goprices.Money        `json:"shipping_price_net" db:"-"`
	ShippingPriceGrossAmount     *decimal.Decimal       `json:"shipping_price_gross_amount"`
	ShippingPriceGross           *goprices.Money        `json:"shipping_price_gross" db:"-"`
	ShippingPrice                *goprices.TaxedMoney   `json:"shipping_price" db:"-"`
	ShippingTaxRate              *decimal.Decimal       `json:"shipping_tax_rate"` // default: Decimal(0)
	Token                        string                 `json:"token"`             // unique
	CheckoutToken                string                 `json:"checkout_token"`
	TotalNetAmount               *decimal.Decimal       `json:"total_net_amount"` // default 0
	UnDiscountedTotalNetAmount   *decimal.Decimal       `json:"undiscounted_total_net_amount"`
	TotalNet                     *goprices.Money        `json:"total_net" db:"-"`
	UnDiscountedTotalNet         *goprices.Money        `json:"undiscounted_total_net" db:"-"`
	TotalGrossAmount             *decimal.Decimal       `json:"total_gross_amount"`
	UnDiscountedTotalGrossAmount *decimal.Decimal       `json:"undiscounted_total_gross_amount"`
	TotalGross                   *goprices.Money        `json:"total_gross" db:"-"`
	UnDiscountedTotalGross       *goprices.Money        `json:"undiscounted_total_gross" db:"-"`
	Total                        *goprices.TaxedMoney   `json:"total" db:"-"`
	UnDiscountedTotal            *goprices.TaxedMoney   `json:"undiscounted_total" db:"-"`
	TotalPaidAmount              *decimal.Decimal       `json:"total_paid_amount"`
	TotalPaid                    *goprices.Money        `json:"total_paid" db:"-"`
	VoucherID                    *string                `json:"voucher_id"`
	GiftCards                    []*giftcard.GiftCard   `json:"gift_cards" db:"-"`
	DisplayGrossPrices           *bool                  `json:"display_gross_prices"`
	CustomerNote                 string                 `json:"customer_note"`
	WeightAmount                 float32                `json:"weight_amount"`
	WeightUnit                   measurement.WeightUnit `json:"weight_unit"`
	Weight                       *measurement.Weight    `json:"weight" db:"-"`
	RedirectUrl                  *string                `json:"redirect_url"`
	model.ModelMetadata
}

func (o *Order) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order.is_valid.%s.app_error",
		"order_id=",
		"Order.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.UserID != nil && !model.IsValidId(*o.UserID) {
		return outer("user_id", &o.Id)
	}
	if o.BillingAddressID != nil && !model.IsValidId(*o.BillingAddressID) {
		return outer("billing_address_id", &o.Id)
	}
	if o.ShippingAddressID != nil && !model.IsValidId(*o.ShippingAddressID) {
		return outer("shipping_address_id", &o.Id)
	}
	if o.OriginalID != nil && !model.IsValidId(*o.OriginalID) {
		return outer("original_id", &o.Id)
	}
	if o.ShippingMethodID != nil && !model.IsValidId(*o.ShippingMethodID) {
		return outer("shipping_method_id", &o.Id)
	}
	if !model.IsValidId(o.ChannelID) {
		return outer("channel_id", &o.Id)
	}
	if o.VoucherID != nil && !model.IsValidId(*o.VoucherID) {
		return outer("voucher_id", &o.Id)
	}
	if len(o.Status) > ORDER_STATUS_MAX_LENGTH {
		return outer("status", &o.Id)
	}
	if len(o.TrackingClientID) > ORDER_TRACKING_CLIENT_ID_MAX_LENGTH {
		return outer("tracking_client_id", &o.Id)
	}
	if len(o.Origin) > ORDER_ORIGIN_MAX_LENGTH {
		return outer("origin", &o.Id)
	}
	if o.ShippingMethodName != nil && utf8.RuneCountInString(*o.ShippingMethodName) > ORDER_SHIPPING_METHOD_NAME_MAX_LENGTH {
		return outer("shipping_method_name", &o.Id)
	}
	if len(o.Token) > ORDER_TOKEN_MAX_LENGTH {
		return outer("token", &o.Id)
	}
	if len(o.CheckoutToken) > ORDER_CHECKOUT_TOKEN_MAX_LENGTH {
		return outer("checkout_token", &o.Id)
	}
	if tag, err := language.Parse(o.LanguageCode); err != nil || !strings.EqualFold(tag.String(), o.LanguageCode) {
		return outer("language_code", &o.Id)
	}
	if !model.IsValidEmail(o.UserEmail) || len(o.UserEmail) > model.USER_EMAIL_MAX_LENGTH {
		return outer("user_email", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}
	if o.CreateAt == 0 {
		return outer("create_at", &o.Id)
	}
	if o.RedirectUrl != nil && !model.IsValidHttpUrl(*o.RedirectUrl) {
		return outer("redirect_url", &o.Id)
	}

	return nil
}

func (o *Order) ToJson() string {
	if o.ShippingPriceNet == nil {
		o.ShippingPriceNet = &goprices.Money{
			Amount:   o.ShippingPriceNetAmount,
			Currency: o.Currency,
		}
	}
	if o.ShippingPriceGross == nil {
		o.ShippingPriceGross = &goprices.Money{
			Amount:   o.ShippingPriceGrossAmount,
			Currency: o.Currency,
		}
	}
	if o.ShippingPrice == nil {
		o.ShippingPrice = &goprices.TaxedMoney{
			Net:      o.ShippingPriceNet,
			Gross:    o.ShippingPriceGross,
			Currency: o.Currency,
		}
	}
	if o.TotalNet == nil {
		o.TotalNet = &goprices.Money{
			Amount:   o.TotalNetAmount,
			Currency: o.Currency,
		}
	}
	if o.UnDiscountedTotalNet == nil {
		o.UnDiscountedTotalNet = &goprices.Money{
			Amount:   o.UnDiscountedTotalNetAmount,
			Currency: o.Currency,
		}
	}
	if o.TotalGross == nil {
		o.TotalGross = &goprices.Money{
			Amount:   o.TotalGrossAmount,
			Currency: o.Currency,
		}
	}
	if o.UnDiscountedTotalGross == nil {
		o.UnDiscountedTotalGross = &goprices.Money{
			Amount:   o.UnDiscountedTotalGrossAmount,
			Currency: o.Currency,
		}
	}
	if o.TotalPaid == nil {
		o.TotalPaid = &goprices.Money{
			Amount:   o.TotalPaidAmount,
			Currency: o.Currency,
		}
	}
	if o.Total == nil {
		o.Total = &goprices.TaxedMoney{
			Net:      o.UnDiscountedTotalNet,
			Gross:    o.UnDiscountedTotalGross,
			Currency: o.Currency,
		}
	}
	if o.UnDiscountedTotal == nil {
		o.UnDiscountedTotal = &goprices.TaxedMoney{
			Net:      o.UnDiscountedTotalNet,
			Gross:    o.UnDiscountedTotalGross,
			Currency: o.Currency,
		}
	}
	if o.Weight == nil {
		o.Weight = &measurement.Weight{
			Amount: o.WeightAmount,
			Unit:   o.WeightUnit,
		}
	}

	return model.ModelToJson(o)
}

func OrderFromJson(data io.Reader) *Order {
	var o Order
	model.ModelFromJson(&o, data)
	return &o
}

func (o *Order) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	if o.CreateAt == 0 {
		o.CreateAt = model.GetMillis()
	}
	if o.Status == "" {
		o.Status = UNFULFILLED
	}
	if o.LanguageCode == "" {
		o.LanguageCode = model.DEFAULT_LANGUAGE_CODE
	}
	if o.ShippingPriceNetAmount == nil {
		o.ShippingPriceNetAmount = &decimal.Zero
	}
	if o.ShippingPriceGrossAmount == nil {
		o.ShippingPriceGrossAmount = &decimal.Zero
	}
	if o.ShippingTaxRate == nil {
		o.ShippingTaxRate = &decimal.Zero
	}
	if o.TotalNetAmount == nil {
		o.TotalNetAmount = &decimal.Zero
	}
	if o.UnDiscountedTotalNetAmount == nil {
		o.UnDiscountedTotalNetAmount = &decimal.Zero
	}
	if o.TotalGrossAmount == nil {
		o.TotalGrossAmount = &decimal.Zero
	}
	if o.TotalPaidAmount == nil {
		o.TotalPaidAmount = &decimal.Zero
	}
	if o.DisplayGrossPrices == nil {
		o.DisplayGrossPrices = model.NewBool(true)
	}
	if o.WeightUnit == "" {
		o.WeightUnit = measurement.KG
	}
	if o.ShippingMethodName != nil {
		o.ShippingMethodName = model.NewString(model.SanitizeUnicode(*o.ShippingMethodName))
	}
}

func (o *Order) IsFullyPaid() bool {
	ok, err := o.Total.Gross.LessThanOrEqual(o.TotalPaid)
	if err != nil {
		slog.Error("error checking order is fully paid", slog.Err(err), slog.String("order_id", o.Id))
		return false
	}
	return ok
}

func (o *Order) IsPartlyPaid() bool {
	return decimal.Zero.LessThan(*o.TotalPaidAmount)
}

func (o *Order) GetCustomerEmail() string {
	panic("not implemented") // TODO: fixme
}

func (o *Order) String() string {
	return "#" + o.Id
}

func (o *Order) IsDraft() bool {
	return o.Status == DRAFT
}

func (o *Order) IsUnconfirmed() bool {
	return o.Status == UNCONFIRMED
}

func (o *Order) IsOpen() bool {
	return o.Status == DRAFT || o.Status == PARTIALLY_FULFILLED
}

func (o *Order) TotalCaptured() *goprices.Money {
	return o.TotalPaid
}

func (o *Order) TotalBalance() (*goprices.Money, error) {
	return o.TotalCaptured().Sub(o.Total.Gross)
}

func (o *Order) GetTotalWeight() *measurement.Weight {
	if o.Weight != nil {
		return o.Weight
	}
	return &measurement.Weight{
		Amount: o.WeightAmount,
		Unit:   o.WeightUnit,
	}
}
