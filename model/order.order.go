package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/measurement"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

// order's Origin field type. Can be either "checkout", "draft" or "reissue"
type OrderOrigin string

// order origin valid values
const (
	ORDER_ORIGIN_CHECKOUT OrderOrigin = "checkout" // order created from checkout
	ORDER_ORIGIN_DRAFT    OrderOrigin = "draft"    // order created from draft order
	ORDER_ORIGIN_REISSUE  OrderOrigin = "reissue"  // order created from reissue existing one
)

func (e OrderOrigin) IsValid() bool {
	switch e {
	case ORDER_ORIGIN_CHECKOUT, ORDER_ORIGIN_DRAFT, ORDER_ORIGIN_REISSUE:
		return true
	}
	return false
}

type OrderStatus string

// ORDER STATUS VALID VALUES
const (
	ORDER_STATUS_DRAFT               OrderStatus = "draft"               // fully editable, not finalized order created by staff users
	ORDER_STATUS_UNCONFIRMED         OrderStatus = "unconfirmed"         // order created by customers when confirmation is required
	ORDER_STATUS_UNFULFILLED         OrderStatus = "unfulfilled"         // order with no items marked as fulfilled
	ORDER_STATUS_PARTIALLY_FULFILLED OrderStatus = "partially_fulfilled" // order with some items marked as fulfilled
	ORDER_STATUS_FULFILLED           OrderStatus = "fulfilled"           // order with all items marked as fulfilled
	ORDER_STATUS_PARTIALLY_RETURNED  OrderStatus = "partially_returned"  // order with some items marked as returned
	ORDER_STATUS_RETURNED            OrderStatus = "returned"            // order with all items marked as returned
	ORDER_STATUS_CANCELED            OrderStatus = "canceled"            // permanently canceled order
)

func (e OrderStatus) IsValid() bool {
	switch e {
	case ORDER_STATUS_DRAFT,
		ORDER_STATUS_UNCONFIRMED,
		ORDER_STATUS_UNFULFILLED,
		ORDER_STATUS_PARTIALLY_FULFILLED,
		ORDER_STATUS_FULFILLED,
		ORDER_STATUS_PARTIALLY_RETURNED,
		ORDER_STATUS_RETURNED,
		ORDER_STATUS_CANCELED:
		return true
	}
	return false
}

type Order struct {
	Id                  string                 `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt            int64                  `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`         // NOT editable
	Status              OrderStatus            `json:"status" gorm:"type:varchar(32);column:Status"`                              // default ORDER_STATUS_UNFULFILLED
	UserID              *string                `json:"user_id" gorm:"type:uuid;column:UserID"`                                    //
	LanguageCode        LanguageCodeEnum       `json:"language_code" gorm:"type:varchar(35);column:LanguageCode"`                 // default: "en"
	TrackingClientID    string                 `json:"tracking_client_id" gorm:"type:varchar(36);column:TrackingClientID"`        // NOT editable
	BillingAddressID    *string                `json:"billing_address_id" gorm:"type:uuid;column:BillingAddressID"`               // NOT editable
	ShippingAddressID   *string                `json:"shipping_address_id" gorm:"type:uuid;column:ShippingAddressID"`             // NOT editable
	UserEmail           string                 `json:"user_email" gorm:"type:varchar(128);column:UserEmail;index:user_email_key"` //
	OriginalID          *string                `json:"original_id" gorm:"type:uuid;column:OriginalID"`                            // original order id
	Origin              OrderOrigin            `json:"origin" gorm:"type:varchar(32);column:Origin"`
	Currency            string                 `json:"currency" gorm:"type:varchar(3);column:Currency"`
	ShippingMethodID    *string                `json:"shipping_method_id" gorm:"type:uuid;column:ShippingMethodID"`
	CollectionPointID   *string                `json:"collection_point_id" gorm:"type:uuid;column:CollectionPointID"`             // foreign key warehosue
	ShippingMethodName  *string                `json:"shipping_method_name" gorm:"type:varchar(255);column:ShippingMethodName"`   // NULL, NOT editable
	CollectionPointName *string                `json:"collection_point_name" gorm:"type:varchar(255);column:CollectionPointName"` // NULL, NOT editable
	ChannelID           string                 `json:"channel_id" gorm:"type:uuid;column:ChannelID"`                              //
	Token               string                 `json:"token" gorm:"type:varchar(36);column:Token;uniqueIndex:token_unique_key"`   // unique
	CheckoutToken       string                 `json:"checkout_token" gorm:"type:varchar(36);column:CheckoutToken"`               //
	VoucherID           *string                `json:"voucher_id" gorm:"type:uuid;column:VoucherID"`
	DisplayGrossPrices  *bool                  `json:"display_gross_prices" gorm:"column:DisplayGrossPrices"` // default *true
	CustomerNote        string                 `json:"customer_note" gorm:"column:CustomerNote"`
	WeightAmount        float32                `json:"weight_amount" gorm:"default:0;column:WeightAmount"`
	WeightUnit          measurement.WeightUnit `json:"weight_unit" gorm:"column:WeightUnit;type:varchar(5)"` // default 'kg'
	RedirectUrl         *string                `json:"redirect_url" gorm:"type:varchar(200);column:RedirectUrl"`
	ShippingTaxRate     *decimal.Decimal       `json:"shipping_tax_rate" gorm:"default:0;column:ShippingTaxRate;type:decimal(5, 4)"` // default Decimal(0)

	TotalPaidAmount *decimal.Decimal `json:"total_paid_amount" gorm:"default:0;column:TotalPaidAmount;type:decimal(12,3)"`
	TotalPaid       *goprices.Money  `json:"total_paid" gorm:"-"`

	TotalNetAmount   *decimal.Decimal     `json:"total_net_amount" gorm:"default:0;column:TotalNetAmount;type:decimal(12,3)"` // default Decimal(0)
	TotalNet         *goprices.Money      `json:"total_net" gorm:"-"`
	TotalGrossAmount *decimal.Decimal     `json:"total_gross_amount" gorm:"default:0;column:TotalGrossAmount;type:decimal(12,3)"` // default decimal(0)
	TotalGross       *goprices.Money      `json:"total_gross" gorm:"-"`
	Total            *goprices.TaxedMoney `json:"total" gorm:"-"` // from TotalNet, TotalGross

	UnDiscountedTotalNetAmount   *decimal.Decimal     `json:"undiscounted_total_net_amount" gorm:"default:0;column:UnDiscountedTotalNetAmount;type:decimal(12,3)"`
	UnDiscountedTotalNet         *goprices.Money      `json:"undiscounted_total_net" gorm:"-"`
	UnDiscountedTotalGrossAmount *decimal.Decimal     `json:"undiscounted_total_gross_amount" gorm:"default:0;column:UnDiscountedTotalGrossAmount;type:decimal(12,3)"` // default 0
	UnDiscountedTotalGross       *goprices.Money      `json:"undiscounted_total_gross" gorm:"-"`
	UnDiscountedTotal            *goprices.TaxedMoney `json:"undiscounted_total" gorm:"-"` // from UnDiscountedTotalNet, UnDiscountedTotalGross

	ShippingPriceNetAmount   *decimal.Decimal     `json:"shipping_price_net_amount" gorm:"default:0;column:ShippingPriceNetAmount;type:decimal(12,3)"` // NOT editable, default Zero
	ShippingPriceNet         *goprices.Money      `json:"shipping_price_net" gorm:"-"`
	ShippingPriceGrossAmount *decimal.Decimal     `json:"shipping_price_gross_amount" gorm:"default:0;column:ShippingPriceGrossAmount;type:decimal(12,3)"` // NOT editable
	ShippingPriceGross       *goprices.Money      `json:"shipping_price_gross" gorm:"-"`
	ShippingPrice            *goprices.TaxedMoney `json:"shipping_price" gorm:"-"` // from ShippingPriceNet, ShippingPriceGross

	Weight *measurement.Weight `json:"weight" gorm:"-"` // default 0

	ModelMetadata
	GiftCards            []*GiftCard `json:"-" gorm:"many2many:OrderGiftCards"`
	OrderLines           OrderLines  `json:"-" gorm:"foreignKey:OrderID"`
	populatedNonDBFields bool        `gorm:"-"`
	Channel              *Channel    `json:"-"`
}

func (c *Order) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Order) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Order) TableName() string             { return OrderTableName }

// OrderFilterOption is used to build sql queries for filtering orders
type OrderFilterOption struct {
	Conditions squirrel.Sqlizer // filter by order's id

	ChannelSlug squirrel.Sqlizer // for comparing the channel of this order's slug

	SelectForUpdate bool // if true, add FOR UPDATE to the end of sql queries. NOTE: Only applies if Transaction is set
	Transaction     *gorm.DB

	Preload []string
}

// PopulateNonDbFields must be called after fetching order(s) from database or before perform json serialization.
func (o *Order) PopulateNonDbFields() {
	if o.populatedNonDBFields {
		return
	}
	o.populatedNonDBFields = true

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
		"model.order.is_valid.%s.app_error",
		"order_id=",
		"Order.IsValid",
	)
	if o.UserID != nil && !IsValidId(*o.UserID) {
		return outer("user_id", &o.Id)
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
	if !IsValidId(o.ChannelID) {
		return outer("channel_id", &o.Id)
	}
	if o.VoucherID != nil && !IsValidId(*o.VoucherID) {
		return outer("voucher_id", &o.Id)
	}
	if !o.LanguageCode.IsValid() {
		return outer("language_code", &o.Id)
	}
	if !IsValidEmail(o.UserEmail) {
		return outer("user_email", &o.Id)
	}
	if unit, err := currency.ParseISO(o.Currency); err != nil || !strings.EqualFold(unit.String(), o.Currency) {
		return outer("currency", &o.Id)
	}
	if o.RedirectUrl != nil && !IsValidHTTPURL(*o.RedirectUrl) {
		return outer("redirect_url", &o.Id)
	}

	return nil
}

func (o *Order) commonPre() {
	// shipping price
	if o.ShippingPrice != nil {
		o.ShippingPriceNet = o.ShippingPrice.Net
		o.ShippingPriceGross = o.ShippingPrice.Gross

		o.ShippingPriceNetAmount = &o.ShippingPriceNet.Amount
		o.ShippingPriceGrossAmount = &o.ShippingPriceGross.Amount
	}
	if o.ShippingPriceNetAmount == nil {
		o.ShippingPriceNetAmount = &decimal.Zero
	}
	if o.ShippingPriceGrossAmount == nil {
		o.ShippingPriceGrossAmount = &decimal.Zero
	}

	// total
	if o.Total != nil {
		o.TotalNet = o.Total.Net
		o.TotalGross = o.Total.Gross

		o.TotalNetAmount = &o.TotalNet.Amount
		o.TotalGrossAmount = &o.TotalGross.Amount
	}
	if o.TotalNetAmount == nil {
		o.TotalNetAmount = &decimal.Zero
	}
	if o.TotalGrossAmount == nil {
		o.TotalGrossAmount = &decimal.Zero
	}

	// total paid
	if o.TotalPaid != nil {
		o.TotalPaidAmount = &o.TotalPaid.Amount
	} else {
		o.TotalPaidAmount = &decimal.Zero
	}

	// un-discounted total
	if o.UnDiscountedTotal != nil {
		o.UnDiscountedTotalNet = o.UnDiscountedTotal.Net
		o.UnDiscountedTotalGross = o.UnDiscountedTotal.Gross

		o.UnDiscountedTotalNetAmount = &o.UnDiscountedTotalNet.Amount
		o.UnDiscountedTotalGrossAmount = &o.UnDiscountedTotalGross.Amount
	}

	if o.DisplayGrossPrices == nil {
		o.DisplayGrossPrices = NewPrimitive(true)
	}
	if o.WeightUnit == "" {
		o.WeightUnit = measurement.KG
	}
	if o.ShippingMethodName != nil {
		o.ShippingMethodName = NewPrimitive(SanitizeUnicode(*o.ShippingMethodName))
	}
	o.CustomerNote = SanitizeUnicode(o.CustomerNote)

	if o.Status == "" {
		o.Status = ORDER_STATUS_UNFULFILLED
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
		o.CollectionPointName = NewPrimitive(SanitizeUnicode(*o.CollectionPointName))
	}
	if o.Weight != nil {
		o.WeightAmount = o.Weight.Amount
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
	return o.Status == ORDER_STATUS_DRAFT
}

// IsUnconfirmed checks if current order's Status is "unconfirmed"
func (o *Order) IsUnconfirmed() bool {
	return o.Status == ORDER_STATUS_UNCONFIRMED
}

// IsOpen checks if current order's Status if "draft" OR "partially_fulfilled"
func (o *Order) IsOpen() bool {
	return o.Status == ORDER_STATUS_DRAFT || o.Status == ORDER_STATUS_PARTIALLY_FULFILLED
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

	s.UserID = CopyPointer(s.UserID)
	s.VoucherID = CopyPointer(s.VoucherID)
	s.OriginalID = CopyPointer(s.OriginalID)
	s.RedirectUrl = CopyPointer(s.RedirectUrl)
	s.TotalNetAmount = CopyPointer(s.TotalNetAmount)
	s.ShippingTaxRate = CopyPointer(s.ShippingTaxRate)
	s.TotalPaidAmount = CopyPointer(s.TotalPaidAmount)
	s.BillingAddressID = CopyPointer(s.BillingAddressID)
	s.ShippingMethodID = CopyPointer(s.ShippingMethodID)
	s.TotalGrossAmount = CopyPointer(s.TotalGrossAmount)
	s.ShippingAddressID = CopyPointer(s.ShippingAddressID)
	s.CollectionPointID = CopyPointer(s.CollectionPointID)
	s.ShippingMethodName = CopyPointer(s.ShippingMethodName)
	s.CollectionPointName = CopyPointer(s.CollectionPointName)
	s.ShippingPriceNetAmount = CopyPointer(s.ShippingPriceNetAmount)
	s.ShippingPriceGrossAmount = CopyPointer(s.ShippingPriceGrossAmount)
	s.UnDiscountedTotalNetAmount = CopyPointer(s.UnDiscountedTotalNetAmount)
	s.UnDiscountedTotalGrossAmount = CopyPointer(s.UnDiscountedTotalGrossAmount)

	return &order
}
