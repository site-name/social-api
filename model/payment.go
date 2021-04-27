package model

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

const (
	// Max lengths for some fields
	MAX_LENGTH_PAYMENT_GATEWAY       = 255
	MAX_LENGTH_PAYMENT_CHARGE_STATUS = 20
	MAX_LENGTH_PAYMENT_TOKEN         = 512
	MAX_LENGTH_PAYMENT_CURRENCY_CODE = 3
	MAX_LENGTH_PAYMENT_COUNTRY_CODE  = 2
	MAX_LENGTH_CC_FIRST_DIGITS       = 6
	MAX_LENGTH_CC_LAST_DIGITS        = 4
	MAX_LENGTHCC_BRAND               = 40
	MIN_CC_EXP_MONTH                 = 1
	MAX_CC_EXP_MONTH                 = 12
	MIN_CC_EXP_YEAR                  = 1000

	// some fields may have max length of 256 in common
	MAX_LENGTH_PAYMENT_COMMON_256 = 256

	// Choices for charge status
	NOT_CHARGED           = "not-charged"
	PENDING               = "pending"
	PARTIALLY_CHARGED     = "partially-charged"
	FULLY_CHARGED         = "fully-charged"
	PARTIALLY_REFUNDED    = "partially-refunded"
	FULLY_REFUNDED        = "fully-refunded"
	REFUSED               = "refused"
	CANCELLED             = "cancelled"
	DEFAULT_CHARGE_STATUS = NOT_CHARGED

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

var validChargeStatues = StringArray([]string{
	NOT_CHARGED,
	PENDING,
	PARTIALLY_CHARGED,
	FULLY_CHARGED,
	PARTIALLY_REFUNDED,
	FULLY_REFUNDED,
	REFUSED,
	CANCELLED,
})

// Default decimal values
var (
	DEFAULT_DECIMAL_VALUE *decimal.NullDecimal
	initOnce              sync.Once
)

func init() {
	initOnce.Do(func() {
		DEFAULT_DECIMAL_VALUE = &decimal.NullDecimal{
			Decimal: decimal.Zero,
			Valid:   true,
		}
	})
}

// Payment represents payment from user to shop
type Payment struct {
	Id                 string               `json:"id"`
	GateWay            string               `json:"gate_way"`
	IsActive           bool                 `json:"is_active"`
	ToConfirm          bool                 `json:"to_confirm"`
	CreateAt           int64                `json:"create_at"`
	UpdateAt           int64                `json:"update_at"`
	ChargeStatus       string               `json:"charge_status"`
	Token              string               `json:"token"`
	Total              *decimal.NullDecimal `json:"total"`
	CapturedAmount     *decimal.NullDecimal `json:"captured_amount"`
	Currency           string               `json:"currency"`
	CheckoutID         string               `json:"checkout_id"`
	OrderID            string               `json:"order_id"`
	BillingEmail       string               `json:"billing_email"`
	BillingFirstName   string               `json:"billing_first_name"`
	BillingLastName    string               `json:"billing_last_name"`
	BillingCompanyName string               `json:"billing_company_name"`
	BillingAddress1    string               `json:"billing_address_1"`
	BillingAddress2    string               `json:"billing_address_2"`
	BillingCity        string               `json:"billing_city"`
	BillingCityArea    string               `json:"billing_city_area"`
	BillingPostalCode  string               `json:"billing_postal_code"`
	BillingCountryCode string               `json:"billing_country_code"`
	BillingCountryArea string               `json:"billing_country_area"`
	CcFirstDigits      string               `json:"cc_first_digits"`
	CcLastDigits       string               `json:"cc_last_digits"`
	CcBranh            string               `json:"cc_brand"`
	CcExpMonth         *uint8               `json:"cc_exp_month"`
	CcExpYear          *uint16              `json:"cc_exp_year"`
	PaymentMethodType  string               `json:"payment_method_type"`
	CustomerIpAddress  *string              `json:"customer_ip_address"`
	ExtraData          string               `json:"extra_data"`
	ReturnUrl          *string              `json:"return_url_url"`
}

func (p *Payment) String() string {
	return fmt.Sprintf(
		"Payment(gateway=%s, is_active=%t, created=%d, charge_status=%s)",
		p.GateWay,
		p.IsActive,
		p.CreateAt,
		p.ChargeStatus,
	)
}

func (p *Payment) GetChargeAmount() decimal.Decimal {
	res := p.Total.Decimal.Sub(p.CapturedAmount.Decimal)
	return res
}

func (p *Payment) IsNotCharged() bool {
	return p.ChargeStatus == NOT_CHARGED
}

func (p *Payment) CanAuthorize() bool {
	return p.IsActive && p.IsNotCharged()
}

func (p *Payment) CanCapture() bool {
	if !p.IsActive && !p.IsNotCharged() {
		return false
	}

	return true
}

// func (p *Payment) CanVoid() bool {
// 	return p.IsActive && p.IsNotCharged() && p.Is
// }

func (p *Payment) CanRefund() bool {
	canRefundChargeStatuses := []string{
		PARTIALLY_CHARGED,
		FULLY_CHARGED,
		PARTIALLY_REFUNDED,
	}

	return p.IsActive &&
		StringArray(canRefundChargeStatuses).Contains(p.ChargeStatus)
}

func (p *Payment) CanConfirm() bool {
	return p.IsActive && p.IsNotCharged()
}

func (p *Payment) IsManual() bool {
	return p.GateWay == GATE_WAY_MANUAL
}

// Common method to create app error for payment
func InvalidPaymentError(fieldName string, paymentID string) *AppError {
	id := fmt.Sprintf("model.payment.is_valid.%s.app_error", fieldName)
	details := ""
	if paymentID != "" {
		details = "payment_id=" + paymentID
	}

	return NewAppError("Payment.IsValid", id, nil, details, http.StatusBadRequest)
}

// Check if input from user is valid or not
func (p *Payment) IsValid() *AppError {
	if !IsValidId(p.Id) {
		return InvalidPaymentError("id", "")
	}
	if !IsValidId(p.OrderID) {
		return InvalidPaymentError("order_id", p.Id)
	}
	if !IsValidId(p.CheckoutID) {
		return InvalidPaymentError("checkout_id", p.Id)
	}
	if p.CreateAt == 0 {
		return InvalidPaymentError("create_at", p.Id)
	}
	if p.UpdateAt == 0 {
		return InvalidPaymentError("update_at", p.Id)
	}
	if utf8.RuneCountInString(p.GateWay) > MAX_LENGTH_PAYMENT_GATEWAY {
		return InvalidPaymentError("gateway", p.Id)
	}
	if p.ChargeStatus == "" ||
		utf8.RuneCountInString(p.ChargeStatus) > MAX_LENGTH_PAYMENT_CHARGE_STATUS ||
		!validChargeStatues.Contains(p.ChargeStatus) {
		return InvalidPaymentError("charge_status", p.Id)
	}
	if utf8.RuneCountInString(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return InvalidPaymentError("token", p.Id)
	}
	if p.Total == nil || !p.Total.Valid {
		return InvalidPaymentError("total", p.Id)
	}
	if p.CapturedAmount == nil || !p.CapturedAmount.Valid {
		return InvalidPaymentError("captured_amount", p.Id)
	}
	if len(p.BillingEmail) > USER_EMAIL_MAX_LENGTH || p.BillingEmail == "" || !IsValidEmail(p.BillingEmail) {
		return InvalidPaymentError("billing_email", p.Id)
	}
	if utf8.RuneCountInString(p.BillingFirstName) > FIRST_NAME_MAX_LENGTH || !IsValidNamePart(p.BillingFirstName, firstName) {
		return InvalidPaymentError("billing_first_name", p.Id)
	}
	if utf8.RuneCountInString(p.BillingLastName) > LAST_NAME_MAX_LENGTH || !IsValidNamePart(p.BillingLastName, lastName) {
		return InvalidPaymentError("billing_last_name", p.Id)
	}
	if utf8.RuneCountInString(p.BillingCompanyName) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("billing_company_name", p.Id)
	}
	if utf8.RuneCountInString(p.BillingAddress1) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("billing_address_1", p.Id)
	}
	if utf8.RuneCountInString(p.BillingAddress2) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("billing_address_2", p.Id)
	}
	if utf8.RuneCountInString(p.BillingCity) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("billing_city", p.Id)
	}
	if utf8.RuneCountInString(p.BillingCityArea) > CITY_AREA_MAX_LENGTH {
		return InvalidPaymentError("billing_city_area", p.Id)
	}
	if utf8.RuneCountInString(p.BillingPostalCode) > POSTAL_CODE_MAX_LENGTH {
		return InvalidPaymentError("billing_postal_code", p.Id)
	}
	if utf8.RuneCountInString(p.BillingCountryCode) > MAX_LENGTH_PAYMENT_COUNTRY_CODE {
		return InvalidPaymentError("billing_country_code", p.Id)
	}
	region, err := language.ParseRegion(p.BillingCountryCode)
	if err != nil || region.String() != p.BillingCountryCode {
		return InvalidPaymentError("billing_country_code", p.Id)
	}
	if utf8.RuneCountInString(p.Currency) > MAX_LENGTH_PAYMENT_CURRENCY_CODE {
		return InvalidPaymentError("currency", p.Id)
	}
	if un, ok := currency.FromRegion(region); !ok || un.String() != p.Currency {
		return InvalidPaymentError("currency", p.Id)
	}
	if utf8.RuneCountInString(p.BillingCountryArea) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("billing_country_area", p.Id)
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_FIRST_DIGITS {
		return InvalidPaymentError("cc_first_digits", p.Id)
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_LAST_DIGITS {
		return InvalidPaymentError("cc_last_digits", p.Id)
	}
	if *p.CcExpMonth < MIN_CC_EXP_MONTH || *p.CcExpMonth > MAX_CC_EXP_MONTH {
		return InvalidPaymentError("cc_exp_month", p.Id)
	}
	if *p.CcExpYear < MIN_CC_EXP_YEAR {
		return InvalidPaymentError("cc_exp_year", p.Id)
	}
	if len(p.PaymentMethodType) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return InvalidPaymentError("payment_method_type", p.Id)
	}

	return nil
}

// populate some fields if empty and perform some sanitizes
func (p *Payment) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}

	p.BillingEmail = NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = SanitizeUnicode(CleanNamePart(p.BillingFirstName, firstName))
	p.BillingLastName = SanitizeUnicode(CleanNamePart(p.BillingLastName, lastName))

	if p.Total == nil || !p.Total.Valid {
		p.Total = DEFAULT_DECIMAL_VALUE
	}

	if p.CapturedAmount == nil || !p.CapturedAmount.Valid {
		p.CapturedAmount = DEFAULT_DECIMAL_VALUE
	}

	if p.ChargeStatus == "" || !validChargeStatues.Contains(p.ChargeStatus) {
		p.ChargeStatus = DEFAULT_CHARGE_STATUS
	}

	p.CreateAt = GetMillis()
	p.UpdateAt = p.CreateAt
}

func (p *Payment) PreUpdate() {
	p.BillingEmail = NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = SanitizeUnicode(CleanNamePart(p.BillingFirstName, firstName))
	p.BillingLastName = SanitizeUnicode(CleanNamePart(p.BillingLastName, lastName))

	p.UpdateAt = GetMillis()
}

func (p *Payment) ToJson() string {
	b, _ := json.JSON.Marshal(p)

	return string(b)
}

func PaymentFromJson(data io.Reader) *Payment {
	var payment *Payment
	json.JSON.NewDecoder(data).Decode(payment)

	return payment
}
