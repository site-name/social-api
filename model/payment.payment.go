package model

import (
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
	Id                 string              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
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
	CheckoutID         *string             `json:"checkout_id" gorm:"type:uuid;column:CheckoutID"`                            // foreign key to checkout
	OrderID            *string             `json:"order_id" gorm:"type:uuid;column:OrderID;index:orderid_key"`                // foreign key to order
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
		return GetPointerOfValue(decimal.Zero)
	}
	return GetPointerOfValue(p.Total.Sub(*p.CapturedAmount))
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
	outer := CreateAppErrorForModel(
		"model.payment.is_valid.%s.app_error",
		"payment_id=",
		"Payment.IsValid",
	)

	if p.OrderID != nil && !IsValidId(*p.OrderID) {
		return outer("order_id", &p.Id)
	}
	if p.CheckoutID != nil && !IsValidId(*p.CheckoutID) {
		return outer("checkout_id", &p.Id)
	}
	if !p.ChargeStatus.IsValid() {
		return outer("charge_status", &p.Id)
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
	if !p.BillingCountryCode.IsValid() {
		return outer("billing_country_code", &p.Id)
	}

	// make sure country code and currency code are match:
	region, _ := language.ParseRegion(string(p.BillingCountryCode))
	if un, ok := currency.FromRegion(region); !ok || !strings.EqualFold(un.String(), p.Currency) {
		return outer("currency", &p.Id)
	}
	if *p.CcExpMonth < MIN_CC_EXP_MONTH || *p.CcExpMonth > MAX_CC_EXP_MONTH {
		return outer("cc_exp_month", &p.Id)
	}
	if *p.CcExpYear < MIN_CC_EXP_YEAR {
		return outer("cc_exp_year", &p.Id)
	}
	if !p.StorePaymentMethod.IsValid() {
		return outer("store_payment_method", &p.Id)
	}

	return nil
}

func (p *Payment) commonPre() {
	p.BillingEmail = NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = SanitizeUnicode(CleanNamePart(p.BillingFirstName, FirstName))
	p.BillingLastName = SanitizeUnicode(CleanNamePart(p.BillingLastName, LastName))
	if p.Total == nil || p.Total.LessThanOrEqual(decimal.Zero) {
		p.Total = GetPointerOfValue(decimal.Zero)
	}
	if p.CapturedAmount == nil || p.CapturedAmount.LessThanOrEqual(decimal.Zero) {
		p.CapturedAmount = GetPointerOfValue(decimal.Zero)
	}
	if !p.ChargeStatus.IsValid() {
		p.ChargeStatus = NOT_CHARGED
	}
	if p.IsActive == nil {
		p.IsActive = GetPointerOfValue(true)
	}
	if p.Currency == "" {
		p.Currency = DEFAULT_CURRENCY
	}
	if !p.StorePaymentMethod.IsValid() {
		p.StorePaymentMethod = NONE
	}
}
