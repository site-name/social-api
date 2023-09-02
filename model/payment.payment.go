package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	MIN_CC_EXP_MONTH = 1
	MAX_CC_EXP_MONTH = 12
	MIN_CC_EXP_YEAR  = 1000

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

type PaymentChargeStatus string

func (p PaymentChargeStatus) IsValid() bool {
	return ChargeStatuString[p] != ""
}

// Choices for charge status
const (
	NOT_CHARGED        PaymentChargeStatus = "not_charged"
	PENDING            PaymentChargeStatus = "pending"
	PARTIALLY_CHARGED  PaymentChargeStatus = "partially_charged"
	FULLY_CHARGED      PaymentChargeStatus = "fully_charged"
	PARTIALLY_REFUNDED PaymentChargeStatus = "partially_refunded"
	FULLY_REFUNDED     PaymentChargeStatus = "fully_refunded"
	REFUSED            PaymentChargeStatus = "refused"
	CANCELLED          PaymentChargeStatus = "cancelled"
)

var ChargeStatuString = map[PaymentChargeStatus]string{
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
	Id                 UUID                `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	GateWay            string              `json:"gate_way" gorm:"type:varchar(255);column:GateWay"`
	IsActive           *bool               `json:"is_active" gorm:"default:true;column:IsActive;index:isactive_key"` // default true
	ToConfirm          bool                `json:"to_confirm" gorm:"column:ToConfirm"`
	CreateAt           int64               `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	UpdateAt           int64               `json:"update_at" gorm:"type:bigint;column:UpdateAt;autoCreateTime:milli;autoUpdateTime:milli"`
	ChargeStatus       PaymentChargeStatus `json:"charge_status" gorm:"type:varchar(20);column:ChargeStatus;index:chargestatus_key"` // default 'not_charged'
	Token              string              `json:"token" gorm:"type:varchar(512);column:Token"`
	Total              *decimal.Decimal    `json:"total" gorm:"default:0;column:Total;type:decimal(12,3)"`                    // DEFAULT decimal(0)
	CapturedAmount     *decimal.Decimal    `json:"captured_amount" gorm:"default:0;column:CapturedAmount;type:decimal(12,3)"` // DEFAULT decimal(0)
	Currency           string              `json:"currency" gorm:"type:varchar(5);column:Currency"`                           // default 'USD'
	CheckoutID         *UUID               `json:"checkout_id" gorm:"type:uuid;column:CheckoutID"`                            // foreign key to checkout
	OrderID            *UUID               `json:"order_id" gorm:"type:uuid;column:OrderID;index:orderid_key"`                // foreign key to order
	BillingEmail       string              `json:"billing_email" gorm:"type:varchar(128);column:BillingEmail"`
	BillingFirstName   string              `json:"billing_first_name" gorm:"type:varchar(256);column:BillingFirstName"`
	BillingLastName    string              `json:"billing_last_name" gorm:"type:varchar(256);column:BillingLastName"`
	BillingCompanyName string              `json:"billing_company_name" gorm:"type:varchar(256);column:BillingCompanyName"`
	BillingAddress1    string              `json:"billing_address_1" gorm:"type:varchar(256);column:BillingAddress1"`
	BillingAddress2    string              `json:"billing_address_2" gorm:"type:varchar(256);column:BillingAddress2"`
	BillingCity        string              `json:"billing_city" gorm:"type:varchar(256);column:BillingCity"`
	BillingCityArea    string              `json:"billing_city_area" gorm:"type:varchar(128);column:BillingCityArea"`
	BillingPostalCode  string              `json:"billing_postal_code" gorm:"type:varchar(256);column:BillingPostalCode"`
	BillingCountryCode CountryCode         `json:"billing_country_code" gorm:"type:varchar(3);column:BillingCountryCode"`
	BillingCountryArea string              `json:"billing_country_area" gorm:"type:varchar(256);column:BillingCountryArea"`
	CcFirstDigits      string              `json:"cc_first_digits" gorm:"type:varchar(6);column:CcFirstDigits"`
	CcLastDigits       string              `json:"cc_last_digits" gorm:"type:varchar(4);column:CcLastDigits"`
	CcBrand            string              `json:"cc_brand" gorm:"type:varchar(40);column:CcBrand"`
	CcExpMonth         *int32              `json:"cc_exp_month" gorm:"column:CcExpMonth"`
	CcExpYear          *int32              `json:"cc_exp_year" gorm:"column:CcExpYear"`
	PaymentMethodType  string              `json:"payment_method_type" gorm:"type:varchar(256);column:PaymentMethodType"`
	CustomerIpAddress  *string             `json:"customer_ip_address" gorm:"type:varchar(40);column:CustomerIpAddress"`
	ExtraData          string              `json:"extra_data" gorm:"column:ExtraData"`
	ReturnUrl          *string             `json:"return_url_url" gorm:"type:varchar(500);column:ReturnUrl"`
	PspReference       *string             `json:"psp_reference" gorm:"type:varchar(512);column:PspReference;index:pspreference_key"` // db index
	StorePaymentMethod StorePaymentMethod  `json:"store_payment_method" gorm:"type:varchar(11);column:StorePaymentMethod"`            // default to "none"
	ModelMetadata
}

func (c *Payment) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Payment) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *Payment) TableName() string             { return PaymentTableName }

type PaymentFilterOption struct {
	Conditions squirrel.Sqlizer

	TransactionsKind           squirrel.Sqlizer // INNER JOIN Transactions ON ... WHERE Transactions.Kind ...
	TransactionsActionRequired squirrel.Sqlizer // INNER JOIN Transactions ON ... WHERE Transactions.ActionRequired ...
	TransactionsIsSuccess      squirrel.Sqlizer // INNER JOIN Transactions ON ... WHERE Transactions.IsSuccess ...

	DbTransaction *gorm.DB
	LockForUpdate bool
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

var canRefundChargeStatuses = util.AnyArray[PaymentChargeStatus]{
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
	if p.OrderID != nil && !IsValidId(*p.OrderID) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.order_id.app_error", nil, "please provide valid payment order id", http.StatusBadRequest)
	}
	if p.CheckoutID != nil && !IsValidId(*p.CheckoutID) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.checkout_id.app_error", nil, "please provide valid payment checkout id", http.StatusBadRequest)
	}
	if !p.ChargeStatus.IsValid() {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.charge_status.app_error", nil, "please provide valid payment charge status", http.StatusBadRequest)
	}
	if p.Total == nil {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.total.app_error", nil, "please provide valid payment total", http.StatusBadRequest)
	}
	if p.CapturedAmount == nil {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.captured_amount.app_error", nil, "please provide valid payment captured amount", http.StatusBadRequest)
	}
	if !IsValidEmail(p.BillingEmail) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_email.app_error", nil, "please provide valid payment billing email", http.StatusBadRequest)
	}
	if !IsValidNamePart(p.BillingFirstName, FirstName) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_firstname.app_error", nil, "please provide valid payment billing first name", http.StatusBadRequest)
	}
	if !IsValidNamePart(p.BillingLastName, LastName) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_lastname.app_error", nil, "please provide valid payment billing last name", http.StatusBadRequest)
	}
	if !p.BillingCountryCode.IsValid() {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.billing_amountry_code.app_error", nil, "please provide valid payment billing country code", http.StatusBadRequest)
	}

	// make sure country code and currency code are match:
	region, _ := language.ParseRegion(string(p.BillingCountryCode))
	if un, ok := currency.FromRegion(region); !ok || !strings.EqualFold(un.String(), p.Currency) {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.currency.app_error", nil, "please provide valid payment currency code", http.StatusBadRequest)
	}
	if *p.CcExpMonth < MIN_CC_EXP_MONTH || *p.CcExpMonth > MAX_CC_EXP_MONTH {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.cc_exp_month.app_error", nil, "please provide valid payment cc exp month", http.StatusBadRequest)
	}
	if *p.CcExpYear < MIN_CC_EXP_YEAR {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.cc_exp_year.app_error", nil, "please provide valid payment cc exp year", http.StatusBadRequest)
	}
	if !p.StorePaymentMethod.IsValid() {
		return NewAppError("Payment.IsValid", "model.payment.is_valid.store_payment_method.app_error", nil, "please provide valid payment store payment method", http.StatusBadRequest)
	}

	return nil
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
	if !p.ChargeStatus.IsValid() {
		p.ChargeStatus = NOT_CHARGED
	}
	if p.IsActive == nil {
		p.IsActive = NewPrimitive(true)
	}
	if p.Currency == "" {
		p.Currency = DEFAULT_CURRENCY
	}
	if !p.StorePaymentMethod.IsValid() {
		p.StorePaymentMethod = NONE
	}
}
