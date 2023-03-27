package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// Max lengths for some payment's fields
const (
	MAX_LENGTH_PAYMENT_GATEWAY              = 255
	MAX_LENGTH_PAYMENT_CHARGE_STATUS        = 20
	MAX_LENGTH_PAYMENT_TOKEN                = 512
	PAYMENT_PSP_REFERENCE_MAX_LENGTH        = 512
	MAX_LENGTH_CC_FIRST_DIGITS              = 6
	MAX_LENGTH_CC_LAST_DIGITS               = 4
	MAX_LENGTH_CC_BRAND                     = 40
	MIN_CC_EXP_MONTH                        = 1
	MAX_CC_EXP_MONTH                        = 12
	MIN_CC_EXP_YEAR                         = 1000
	MAX_LENGTH_PAYMENT_COMMON_256           = 256
	PAYMENT_STORE_PAYMENT_METHOD_MAX_LENGTH = 11

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

// Choices for charge status
const (
	NOT_CHARGED        = "not_charged"
	PENDING            = "pending"
	PARTIALLY_CHARGED  = "partially_charged"
	FULLY_CHARGED      = "fully_charged"
	PARTIALLY_REFUNDED = "partially_refunded"
	FULLY_REFUNDED     = "fully_refunded"
	REFUSED            = "refused"
	CANCELLED          = "cancelled"
)

var ChargeStatuString = map[string]string{
	NOT_CHARGED:        "Not charged",
	PENDING:            "Pending",
	PARTIALLY_CHARGED:  "Partially charged",
	FULLY_CHARGED:      "Fully charged",
	PARTIALLY_REFUNDED: "Partially refunded",
	FULLY_REFUNDED:     "Fully refunded",
	REFUSED:            "Refused",
	CANCELLED:          "Cancelled",
}

// Payment represents payment from user to shop
type Payment struct {
	Id                 string             `json:"id"`
	GateWay            string             `json:"gate_way"`
	IsActive           *bool              `json:"is_active"` // default true
	ToConfirm          bool               `json:"to_confirm"`
	CreateAt           int64              `json:"create_at"`
	UpdateAt           int64              `json:"update_at"`
	ChargeStatus       string             `json:"charge_status"`
	Token              string             `json:"token"`
	Total              *decimal.Decimal   `json:"total"`           // DEFAULT decimal(0)
	CapturedAmount     *decimal.Decimal   `json:"captured_amount"` // DEFAULT decimal(0)
	Currency           string             `json:"currency"`        // default 'USD'
	CheckoutID         *string            `json:"checkout_id"`     // foreign key to checkout
	OrderID            *string            `json:"order_id"`        // foreign key to order
	BillingEmail       string             `json:"billing_email"`
	BillingFirstName   string             `json:"billing_first_name"`
	BillingLastName    string             `json:"billing_last_name"`
	BillingCompanyName string             `json:"billing_company_name"`
	BillingAddress1    string             `json:"billing_address_1"`
	BillingAddress2    string             `json:"billing_address_2"`
	BillingCity        string             `json:"billing_city"`
	BillingCityArea    string             `json:"billing_city_area"`
	BillingPostalCode  string             `json:"billing_postal_code"`
	BillingCountryCode CountryCode        `json:"billing_country_code"`
	BillingCountryArea string             `json:"billing_country_area"`
	CcFirstDigits      string             `json:"cc_first_digits"`
	CcLastDigits       string             `json:"cc_last_digits"`
	CcBrand            string             `json:"cc_brand"`
	CcExpMonth         *uint8             `json:"cc_exp_month"`
	CcExpYear          *uint16            `json:"cc_exp_year"`
	PaymentMethodType  string             `json:"payment_method_type"`
	CustomerIpAddress  *string            `json:"customer_ip_address"`
	ExtraData          string             `json:"extra_data"`
	ReturnUrl          *string            `json:"return_url_url"`
	PspReference       *string            `json:"psp_reference"`        // db index
	StorePaymentMethod StorePaymentMethod `json:"store_payment_method"` // default to "none"
	ModelMetadata
}

// PaymentFilterOption is used to build sql queries
type PaymentFilterOption struct {
	Id                         squirrel.Sqlizer
	OrderID                    squirrel.Sqlizer
	CheckoutID                 squirrel.Sqlizer
	IsActive                   *bool
	TransactionsKind           squirrel.Sqlizer // for filtering payment's transactions's `Kind`
	TransactionsActionRequired *bool            // for filtering payment's transactions's `ActionRequired`
	TransactionsIsSuccess      *bool            // for filtering payment's transactions's `IsSuccess`
}

// PaymentPatch is used to update payments
type PaymentPatch struct {
	CheckoutID   string
	OrderID      string
	BillingEmail string
}

// Retrieve the maximum capture possible.
func (p *Payment) GetChargeAmount() *decimal.Decimal {
	if p.Total == nil || p.CapturedAmount == nil {
		return &decimal.Zero
	}
	return NewPrimitive(p.Total.Sub(*p.CapturedAmount))
}

// NotCharged checks if current payment's charge status is "not_charged"
func (p *Payment) NotCharged() bool {
	return p.ChargeStatus == NOT_CHARGED
}

// CanAuthorize checks if current payment is active and not charged
func (p *Payment) CanAuthorize() bool {
	return *p.IsActive && p.NotCharged()
}

// CanCapture checks if payment is not active and is not charged => false, else => true.
func (p *Payment) CanCapture() bool {
	return *p.IsActive && p.NotCharged()
}

var canRefundChargeStatuses = util.AnyArray[string]{
	PARTIALLY_CHARGED,
	FULLY_CHARGED,
	PARTIALLY_REFUNDED,
}

// CanRefund checks if current payment is active && (partially charged || fully charged || partially refunded)
func (p *Payment) CanRefund() bool {
	return *p.IsActive && canRefundChargeStatuses.Contains(p.ChargeStatus)
}

// CanConfirm checks if current payment is active && not charged
func (p *Payment) CanConfirm() bool {
	return *p.IsActive && p.NotCharged()
}

// IsManual checks if current payment's gateway == "manual"
func (p *Payment) IsManual() bool {
	return p.GateWay == GATE_WAY_MANUAL
}

func (p *Payment) GetTotal() *goprices.Money {
	return &goprices.Money{
		Amount:   *p.Total,
		Currency: p.Currency,
	}
}

func (p *Payment) GetCapturedAmount() *goprices.Money {
	return &goprices.Money{
		Amount:   *p.CapturedAmount,
		Currency: p.Currency,
	}
}

// Check if input from user is valid or not
func (p *Payment) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.payment.is_valid.%s.app_error",
		"payment_id=",
		"Payment.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if p.OrderID != nil && !IsValidId(*p.OrderID) {
		return outer("order_id", &p.Id)
	}
	if p.CheckoutID != nil && !IsValidId(*p.CheckoutID) {
		return outer("checkout_id", &p.Id)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if p.UpdateAt == 0 {
		return outer("update_at", &p.Id)
	}
	if utf8.RuneCountInString(p.GateWay) > MAX_LENGTH_PAYMENT_GATEWAY {
		return outer("gateway", &p.Id)
	}
	if ChargeStatuString[strings.ToLower(p.ChargeStatus)] == "" {
		return outer("charge_status", &p.Id)
	}
	if utf8.RuneCountInString(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return outer("token", &p.Id)
	}
	if p.Total == nil {
		return outer("total", &p.Id)
	}
	if p.CapturedAmount == nil {
		return outer("captured_amount", &p.Id)
	}
	if !IsValidEmail(p.BillingEmail) {
		return outer("billing_email", &p.Id)
	}
	if !IsValidNamePart(p.BillingFirstName, FirstName) {
		return outer("billing_first_name", &p.Id)
	}
	if !IsValidNamePart(p.BillingLastName, LastName) {
		return outer("billing_last_name", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingCompanyName) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("billing_company_name", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingAddress1) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("billing_address_1", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingAddress2) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("billing_address_2", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingCity) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("billing_city", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingCityArea) > ADDRESS_CITY_AREA_MAX_LENGTH {
		return outer("billing_city_area", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingPostalCode) > ADDRESS_POSTAL_CODE_MAX_LENGTH {
		return outer("billing_postal_code", &p.Id)
	}
	// make sure country code and currency code are match:
	region, err := language.ParseRegion(string(p.BillingCountryCode))
	if err != nil || !strings.EqualFold(region.String(), string(p.BillingCountryCode)) {
		return outer("billing_country_code", &p.Id)
	}
	if un, ok := currency.FromRegion(region); !ok || !strings.EqualFold(un.String(), p.Currency) {
		return outer("currency", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingCountryArea) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("billing_country_area", &p.Id)
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_FIRST_DIGITS {
		return outer("cc_first_digits", &p.Id)
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_LAST_DIGITS {
		return outer("cc_last_digits", &p.Id)
	}
	if *p.CcExpMonth < MIN_CC_EXP_MONTH || *p.CcExpMonth > MAX_CC_EXP_MONTH {
		return outer("cc_exp_month", &p.Id)
	}
	if *p.CcExpYear < MIN_CC_EXP_YEAR {
		return outer("cc_exp_year", &p.Id)
	}
	if len(p.PaymentMethodType) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return outer("payment_method_type", &p.Id)
	}
	if len(p.CcBrand) > MAX_LENGTH_CC_BRAND {
		return outer("cc_brand", &p.Id)
	}
	if p.PspReference != nil && len(*p.PspReference) > PAYMENT_PSP_REFERENCE_MAX_LENGTH {
		return outer("psp_reference", &p.Id)
	}
	if len(p.StorePaymentMethod) > PAYMENT_STORE_PAYMENT_METHOD_MAX_LENGTH || StorePaymentMethodStringValues[p.StorePaymentMethod] == "" {
		return outer("store_payment_method", &p.Id)
	}

	return nil
}

// populate some fields if empty and perform some sanitizes
func (p *Payment) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.CreateAt = GetMillis()
	p.UpdateAt = p.CreateAt
	p.commonPre()
}

func (p *Payment) commonPre() {
	p.BillingEmail = NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = SanitizeUnicode(CleanNamePart(p.BillingFirstName, FirstName))
	p.BillingLastName = SanitizeUnicode(CleanNamePart(p.BillingLastName, LastName))
	if p.Total == nil || p.Total.LessThanOrEqual(decimal.Zero) {
		p.Total = &decimal.Zero
	}
	if p.CapturedAmount == nil || p.CapturedAmount.LessThanOrEqual(decimal.Zero) {
		p.CapturedAmount = &decimal.Zero
	}
	if _, ok := ChargeStatuString[strings.ToLower(p.ChargeStatus)]; !ok {
		p.ChargeStatus = NOT_CHARGED
	}
	if p.IsActive == nil {
		p.IsActive = NewPrimitive(true)
	}
	if p.Currency == "" {
		p.Currency = DEFAULT_CURRENCY
	}
	if StorePaymentMethodStringValues[p.StorePaymentMethod] == "" {
		p.StorePaymentMethod = NONE
	}
}

func (p *Payment) PreUpdate() {
	p.UpdateAt = GetMillis()
	p.commonPre()
}

func (p *Payment) ToJSON() string {
	return ModelToJson(p)
}
