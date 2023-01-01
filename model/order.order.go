package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/measurement"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// max lengths for some fields of order model
const (
	ORDER_STATUS_MAX_LENGTH                = 32
	ORDER_TRACKING_CLIENT_ID_MAX_LENGTH    = 36
	ORDER_ORIGIN_MAX_LENGTH                = 32
	ORDER_SHIPPING_METHOD_NAME_MAX_LENGTH  = 255
	ORDER_TOKEN_MAX_LENGTH                 = 36
	ORDER_CHECKOUT_TOKEN_MAX_LENGTH        = 36
	ORDER_COLLECTION_POINT_NAME_MAX_LENGTH = 255
)

type OrderOrigin string

// order origin valid values
const (
	CHECKOUT OrderOrigin = "checkout" // order created from checkout
	DRAFT    OrderOrigin = "draft"    // order created from draft order
	REISSUE  OrderOrigin = "reissue"  // order created from reissue existing one
)

type OrderStatus string

// ORDER STATUS VALID VALUES
const (
	STATUS_DRAFT        OrderStatus = "draft"               // fully editable, not finalized order created by staff users
	UNCONFIRMED         OrderStatus = "unconfirmed"         // order created by customers when confirmation is required
	UNFULFILLED         OrderStatus = "unfulfilled"         // order with no items marked as fulfilled
	PARTIALLY_FULFILLED OrderStatus = "partially_fulfilled" // order with some items marked as fulfilled
	FULFILLED           OrderStatus = "fulfilled"           // order with all items marked as fulfilled
	PARTIALLY_RETURNED  OrderStatus = "partially_returned"  // order with some items marked as returned
	RETURNED            OrderStatus = "returned"            // order with all items marked as returned
	CANCELED            OrderStatus = "canceled"            // permanently canceled order
)

var OrderStatusStrings = map[OrderStatus]string{
	STATUS_DRAFT:        "Draft",
	UNCONFIRMED:         "Unconfirmed",
	UNFULFILLED:         "Unfulfilled",
	PARTIALLY_FULFILLED: "Partially fulfilled",
	PARTIALLY_RETURNED:  "Partially returned",
	RETURNED:            "Returned",
	FULFILLED:           "Fulfilled",
	CANCELED:            "Canceled",
}

var OrderOriginStrings = map[OrderOrigin]string{
	CHECKOUT: "Checkout",
	DRAFT:    "Draft",
	REISSUE:  "Reissue",
}

type Order struct {
	Id                           string                 `json:"id"`
	CreateAt                     int64                  `json:"create_at"` // NOT editable
	Status                       OrderStatus            `json:"status"`    // default: UNFULFILLED
	UserID                       *string                `json:"user_id"`   //
	ShopID                       string                 `json:"shop_id"`
	LanguageCode                 string                 `json:"language_code"`       // default: "en"
	TrackingClientID             string                 `json:"tracking_client_id"`  // NOT editable
	BillingAddressID             *string                `json:"billing_address_id"`  // NOT editable
	ShippingAddressID            *string                `json:"shipping_address_id"` // NOT editable
	UserEmail                    string                 `json:"user_email"`          //
	OriginalID                   *string                `json:"original_id"`         // original order id
	Origin                       OrderOrigin            `json:"origin"`
	Currency                     string                 `json:"currency"`
	ShippingMethodID             *string                `json:"shipping_method_id"`
	CollectionPointID            *string                `json:"collection_point_id"`         // foreign key warehosue
	ShippingMethodName           *string                `json:"shipping_method_name"`        // NUL, NOT editable
	CollectionPointName          *string                `json:"collection_point_name"`       // NUL, NOTE editable
	ChannelID                    string                 `json:"channel_id"`                  //
	ShippingPriceNetAmount       *decimal.Decimal       `json:"shipping_price_net_amount"`   // NOT editable, default Zero
	ShippingPriceNet             *goprices.Money        `json:"shipping_price_net" db:"-"`   //
	ShippingPriceGrossAmount     *decimal.Decimal       `json:"shipping_price_gross_amount"` // NOT editable
	ShippingPriceGross           *goprices.Money        `json:"shipping_price_gross" db:"-"`
	ShippingPrice                *goprices.TaxedMoney   `json:"shipping_price" db:"-"`
	ShippingTaxRate              *decimal.Decimal       `json:"shipping_tax_rate"` // default: Decimal(0)
	Token                        string                 `json:"token"`             // unique
	CheckoutToken                string                 `json:"checkout_token"`    //
	TotalNetAmount               *decimal.Decimal       `json:"total_net_amount"`  // default 0
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
	DisplayGrossPrices           *bool                  `json:"display_gross_prices"` // default *true
	CustomerNote                 string                 `json:"customer_note"`
	WeightAmount                 float32                `json:"weight_amount"`
	WeightUnit                   measurement.WeightUnit `json:"weight_unit"`   // default 'kg'
	Weight                       *measurement.Weight    `json:"weight" db:"-"` // default 0
	RedirectUrl                  *string                `json:"redirect_url"`
	ModelMetadata

	populatedNonDBFields bool `json:"-" db:"-"`
}

// OrderFilterOption is used to buils sql queries for filtering orders
type OrderFilterOption struct {
	Id            squirrel.Sqlizer // filter by order's id
	Status        squirrel.Sqlizer // for filtering order's Status
	CheckoutToken squirrel.Sqlizer // for filtering order's CheckoutToken
	ChannelSlug   squirrel.Sqlizer // for comparing the channel of this order's slug
	UserEmail     squirrel.Sqlizer // for filtering order's UserEmail
	UserID        squirrel.Sqlizer // for filtering order's UserID
}

// PopulateNonDbFields must be called after fetching order(s) from database or before perform json serialization.
func (o *Order) PopulateNonDbFields() {
	if o.populatedNonDBFields {
		return
	}
	defer func() {
		o.populatedNonDBFields = true
	}()

	// errors can be ignored since orders's Currencies were checked before saving into database
	o.ShippingPriceNet = &goprices.Money{
		Amount:   *o.ShippingPriceNetAmount,
		Currency: o.Currency,
	}
	o.ShippingPriceGross = &goprices.Money{
		Amount:   *o.ShippingPriceGrossAmount,
		Currency: o.Currency,
	}
	o.ShippingPrice, _ = goprices.NewTaxedMoney(o.ShippingPriceNet, o.ShippingPriceGross)
	o.TotalNet = &goprices.Money{
		Amount:   *o.TotalNetAmount,
		Currency: o.Currency,
	}
	o.UnDiscountedTotalNet = &goprices.Money{
		Amount:   *o.UnDiscountedTotalNetAmount,
		Currency: o.Currency,
	}
	o.TotalGross = &goprices.Money{
		Amount:   *o.TotalGrossAmount,
		Currency: o.Currency,
	}
	o.UnDiscountedTotalGross = &goprices.Money{
		Amount:   *o.UnDiscountedTotalGrossAmount,
		Currency: o.Currency,
	}
	o.UnDiscountedTotalGross = &goprices.Money{
		Amount:   *o.UnDiscountedTotalGrossAmount,
		Currency: o.Currency,
	}
	o.Total, _ = goprices.NewTaxedMoney(o.TotalNet, o.TotalGross)                                     // ignore error since arguments are trusted
	o.UnDiscountedTotal, _ = goprices.NewTaxedMoney(o.UnDiscountedTotalNet, o.UnDiscountedTotalGross) // ignore error since arguments are trusted
	o.TotalPaid = &goprices.Money{
		Amount:   *o.TotalPaidAmount,
		Currency: o.Currency,
	}
	o.Weight = &measurement.Weight{
		Amount: o.WeightAmount,
		Unit:   o.WeightUnit,
	}
}

// Orders is slice contains order(s)
type Orders []*Order

func (os Orders) PopulateNonDbFields() {
	for _, o := range os {
		o.PopulateNonDbFields()
	}
}

func (orders Orders) ChannelIDs() []string {
	return lo.Map(orders, func(o *Order, _ int) string { return o.ChannelID })
}

func (o *Order) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"order.is_valid.%s.app_error",
		"order_id=",
		"Order.IsValid",
	)
	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.UserID != nil && !IsValidId(*o.UserID) {
		return outer("user_id", &o.Id)
	}
	if !IsValidId(o.ShopID) {
		return outer("shop_id", &o.Id)
	}
	if o.BillingAddressID != nil && !IsValidId(*o.BillingAddressID) {
		return outer("billing_address_id", &o.Id)
	}
	if o.ShippingAddressID != nil && !IsValidId(*o.ShippingAddressID) {
		return outer("shipping_address_id", &o.Id)
	}
	if o.OriginalID != nil && !IsValidId(*o.OriginalID) {
		return outer("original_id", &o.Id)
	}
	if o.ShippingMethodID != nil && !IsValidId(*o.ShippingMethodID) {
		return outer("shipping_method_id", &o.Id)
	}
	if o.CollectionPointID != nil && !IsValidId(*o.CollectionPointID) {
		return outer("collection_point_id", &o.Id)
	}
	if o.CollectionPointName != nil && utf8.RuneCountInString(*o.CollectionPointName) > ORDER_COLLECTION_POINT_NAME_MAX_LENGTH {
		return outer("collection_point_name", &o.Id)
	}
	if !IsValidId(o.ChannelID) {
		return outer("channel_id", &o.Id)
	}
	if o.VoucherID != nil && !IsValidId(*o.VoucherID) {
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
	if !IsValidEmail(o.UserEmail) || len(o.UserEmail) > USER_EMAIL_MAX_LENGTH {
		return outer("user_email", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}
	if o.CreateAt == 0 {
		return outer("create_at", &o.Id)
	}
	if o.RedirectUrl != nil && !IsValidHTTPURL(*o.RedirectUrl) {
		return outer("redirect_url", &o.Id)
	}

	return nil
}

func (o *Order) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
	o.CreateAt = GetMillis()

	o.commonPre()
}

// PreUpdate
func (o *Order) PreUpdate() {
	o.commonPre()
}

func (o *Order) commonPre() {
	if o.ShippingPriceNet != nil {
		o.ShippingPriceNetAmount = &o.ShippingPriceNet.Amount
	} else {
		o.ShippingPriceNetAmount = &decimal.Zero
	}

	if o.ShippingPriceGross != nil {
		o.ShippingPriceGrossAmount = &o.ShippingPriceGross.Amount
	} else {
		o.ShippingPriceGrossAmount = &decimal.Zero
	}

	if o.TotalNet != nil {
		o.TotalNetAmount = &o.TotalNet.Amount
	} else {
		o.TotalNetAmount = &decimal.Zero
	}

	if o.TotalGross != nil {
		o.TotalGrossAmount = &o.TotalGross.Amount
	} else {
		o.TotalGrossAmount = &decimal.Zero
	}

	if o.TotalPaid != nil {
		o.TotalPaidAmount = &o.TotalPaid.Amount
	} else {
		o.TotalPaidAmount = &decimal.Zero
	}

	if o.DisplayGrossPrices == nil {
		o.DisplayGrossPrices = NewBool(true)
	}
	if o.WeightUnit == "" {
		o.WeightUnit = measurement.KG
	}
	if o.ShippingMethodName != nil {
		o.ShippingMethodName = NewString(SanitizeUnicode(*o.ShippingMethodName))
	}
	o.CustomerNote = SanitizeUnicode(o.CustomerNote)

	if o.Status == "" {
		o.Status = UNFULFILLED
	}
	if o.LanguageCode == "" {
		o.LanguageCode = DEFAULT_LOCALE
	}
	if o.Token == "" {
		o.NewToken()
	}
	if o.Currency != "" {
		o.Currency = strings.ToUpper(o.Currency)
	} else {
		o.Currency = DEFAULT_CURRENCY
	}
	if o.CollectionPointName != nil {
		o.CollectionPointName = NewString(SanitizeUnicode(*o.CollectionPointName))
	}
}

// NewToken generates an uuid and assign it to current order's Id
func (o *Order) NewToken() {
	o.Token = NewId()
}

// IsFullyPaid checks current order's total paid is greater than its total gross
func (o *Order) IsFullyPaid() bool {
	o.PopulateNonDbFields()
	return o.Total.Gross.LessThanOrEqual(o.TotalPaid)
}

// IsPartlyPaid checks if order has `TotalPaidAmount` > 0
func (o *Order) IsPartlyPaid() bool {
	o.PopulateNonDbFields()
	return o.TotalPaidAmount != nil && decimal.Zero.LessThan(*o.TotalPaidAmount)
}

// IsDraft checks if current order's Status if "draft"
func (o *Order) IsDraft() bool {
	return o.Status == STATUS_DRAFT
}

// IsUnconfirmed checks if current order's Status is "unconfirmed"
func (o *Order) IsUnconfirmed() bool {
	return o.Status == UNCONFIRMED
}

// IsOpen checks if current order's Status if "draft" OR "partially_fulfilled"
func (o *Order) IsOpen() bool {
	return o.Status == STATUS_DRAFT || o.Status == PARTIALLY_FULFILLED
}

// TotalCaptured returns current order's TotalPaid money
func (o *Order) TotalCaptured() *goprices.Money {
	o.PopulateNonDbFields()
	return o.TotalPaid
}

// TotalBalance substracts order's total paid to order's total gross
func (o *Order) TotalBalance() *goprices.Money {
	value, _ := o.TotalCaptured().Sub(o.Total.Gross)
	return value
}

// GetTotalWeight returns current order's Weight
func (o *Order) GetTotalWeight() *measurement.Weight {
	o.PopulateNonDbFields()
	return o.Weight
}

func (s *Order) DeepCopy() *Order {
	order := *s

	if s.UserID != nil {
		order.UserID = NewString(*s.UserID)
	}
	if s.BillingAddressID != nil {
		order.BillingAddressID = NewString(*s.BillingAddressID)
	}
	if s.ShippingAddressID != nil {
		order.ShippingAddressID = NewString(*s.ShippingAddressID)
	}
	if s.ShippingMethodID != nil {
		order.ShippingMethodID = NewString(*s.ShippingMethodID)
	}
	if s.OriginalID != nil {
		order.OriginalID = NewString(*s.OriginalID)
	}
	if s.CollectionPointID != nil {
		order.CollectionPointID = NewString(*s.CollectionPointID)
	}
	if s.ShippingMethodName != nil {
		order.ShippingMethodName = NewString(*s.ShippingMethodName)
	}
	if s.CollectionPointName != nil {
		order.CollectionPointName = NewString(*s.CollectionPointName)
	}
	if s.VoucherID != nil {
		order.VoucherID = NewString(*s.VoucherID)
	}
	if s.RedirectUrl != nil {
		order.RedirectUrl = NewString(*s.RedirectUrl)
	}

	if s.ShippingPriceNetAmount != nil {
		order.ShippingPriceNetAmount = NewDecimal(*s.ShippingPriceNetAmount)
	}
	if s.ShippingPriceGrossAmount != nil {
		order.ShippingPriceGrossAmount = NewDecimal(*s.ShippingPriceGrossAmount)
	}
	if s.ShippingTaxRate != nil {
		order.ShippingTaxRate = NewDecimal(*s.ShippingTaxRate)
	}
	if s.TotalNetAmount != nil {
		order.TotalNetAmount = NewDecimal(*s.TotalNetAmount)
	}
	if s.UnDiscountedTotalNetAmount != nil {
		order.UnDiscountedTotalNetAmount = NewDecimal(*s.UnDiscountedTotalNetAmount)
	}

	if s.TotalGrossAmount != nil {
		order.TotalGrossAmount = NewDecimal(*s.TotalGrossAmount)
	}
	if s.UnDiscountedTotalGrossAmount != nil {
		order.UnDiscountedTotalGrossAmount = NewDecimal(*s.UnDiscountedTotalGrossAmount)
	}
	if s.TotalPaidAmount != nil {
		order.TotalPaidAmount = NewDecimal(*s.TotalPaidAmount)
	}

	return &order
}
