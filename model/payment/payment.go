package payment

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

const (
	// Max lengths for some fields
	MAX_LENGTH_PAYMENT_GATEWAY       = 255
	MAX_LENGTH_PAYMENT_CHARGE_STATUS = 20
	MAX_LENGTH_PAYMENT_TOKEN         = 512
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

var validChargeStatues = model.StringArray([]string{
	NOT_CHARGED,
	PENDING,
	PARTIALLY_CHARGED,
	FULLY_CHARGED,
	PARTIALLY_REFUNDED,
	FULLY_REFUNDED,
	REFUSED,
	CANCELLED,
})

// Payment represents payment from user to shop
type Payment struct {
	Id                 string           `json:"id"`
	GateWay            string           `json:"gate_way"`
	IsActive           bool             `json:"is_active"`
	ToConfirm          bool             `json:"to_confirm"`
	CreateAt           int64            `json:"create_at"`
	UpdateAt           int64            `json:"update_at"`
	ChargeStatus       string           `json:"charge_status"`
	Token              string           `json:"token"`
	Total              *decimal.Decimal `json:"total"`
	CapturedAmount     *decimal.Decimal `json:"captured_amount"`
	Currency           string           `json:"currency"`
	CheckoutID         string           `json:"checkout_id"`
	OrderID            string           `json:"order_id"`
	BillingEmail       string           `json:"billing_email"`
	BillingFirstName   string           `json:"billing_first_name"`
	BillingLastName    string           `json:"billing_last_name"`
	BillingCompanyName string           `json:"billing_company_name"`
	BillingAddress1    string           `json:"billing_address_1"`
	BillingAddress2    string           `json:"billing_address_2"`
	BillingCity        string           `json:"billing_city"`
	BillingCityArea    string           `json:"billing_city_area"`
	BillingPostalCode  string           `json:"billing_postal_code"`
	BillingCountryCode string           `json:"billing_country_code"`
	BillingCountryArea string           `json:"billing_country_area"`
	CcFirstDigits      string           `json:"cc_first_digits"`
	CcLastDigits       string           `json:"cc_last_digits"`
	CcBranh            string           `json:"cc_brand"`
	CcExpMonth         *uint8           `json:"cc_exp_month"`
	CcExpYear          *uint16          `json:"cc_exp_year"`
	PaymentMethodType  string           `json:"payment_method_type"`
	CustomerIpAddress  *string          `json:"customer_ip_address"`
	ExtraData          string           `json:"extra_data"`
	ReturnUrl          *string          `json:"return_url_url"`
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
	res := p.Total.Sub(*p.CapturedAmount)
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
		model.StringArray(canRefundChargeStatuses).Contains(p.ChargeStatus)
}

func (p *Payment) CanConfirm() bool {
	return p.IsActive && p.IsNotCharged()
}

func (p *Payment) IsManual() bool {
	return p.GateWay == GATE_WAY_MANUAL
}

// Common method to create app error for payment
func (p *Payment) InvalidPaymentError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.payment.is_valid.%s.app_error", fieldName)
	details := ""
	if !strings.EqualFold(fieldName, "id") {
		details = "payment_id=" + p.Id
	}

	return model.NewAppError("Payment.IsValid", id, nil, details, http.StatusBadRequest)
}

// Check if input from user is valid or not
func (p *Payment) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.InvalidPaymentError("id")
	}
	if !model.IsValidId(p.OrderID) {
		return p.InvalidPaymentError("order_id")
	}
	if !model.IsValidId(p.CheckoutID) {
		return p.InvalidPaymentError("checkout_id")
	}
	if p.CreateAt == 0 {
		return p.InvalidPaymentError("create_at")
	}
	if p.UpdateAt == 0 {
		return p.InvalidPaymentError("update_at")
	}
	if utf8.RuneCountInString(p.GateWay) > MAX_LENGTH_PAYMENT_GATEWAY {
		return p.InvalidPaymentError("gateway")
	}
	if p.ChargeStatus == "" ||
		utf8.RuneCountInString(p.ChargeStatus) > MAX_LENGTH_PAYMENT_CHARGE_STATUS ||
		!validChargeStatues.Contains(p.ChargeStatus) {
		return p.InvalidPaymentError("charge_status")
	}
	if utf8.RuneCountInString(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return p.InvalidPaymentError("token")
	}
	if p.Total == nil {
		return p.InvalidPaymentError("total")
	}
	if p.CapturedAmount == nil {
		return p.InvalidPaymentError("captured_amount")
	}
	if len(p.BillingEmail) > model.USER_EMAIL_MAX_LENGTH || p.BillingEmail == "" || !model.IsValidEmail(p.BillingEmail) {
		return p.InvalidPaymentError("billing_email")
	}
	if utf8.RuneCountInString(p.BillingFirstName) > account.FIRST_NAME_MAX_LENGTH || !account.IsValidNamePart(p.BillingFirstName, model.FirstName) {
		return p.InvalidPaymentError("billing_first_name")
	}
	if utf8.RuneCountInString(p.BillingLastName) > account.LAST_NAME_MAX_LENGTH || !account.IsValidNamePart(p.BillingLastName, model.LastName) {
		return p.InvalidPaymentError("billing_last_name")
	}
	if utf8.RuneCountInString(p.BillingCompanyName) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("billing_company_name")
	}
	if utf8.RuneCountInString(p.BillingAddress1) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("billing_address_1")
	}
	if utf8.RuneCountInString(p.BillingAddress2) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("billing_address_2")
	}
	if utf8.RuneCountInString(p.BillingCity) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("billing_city")
	}
	if utf8.RuneCountInString(p.BillingCityArea) > account.CITY_AREA_MAX_LENGTH {
		return p.InvalidPaymentError("billing_city_area")
	}
	if utf8.RuneCountInString(p.BillingPostalCode) > account.POSTAL_CODE_MAX_LENGTH {
		return p.InvalidPaymentError("billing_postal_code")
	}
	if utf8.RuneCountInString(p.BillingCountryCode) > model.MAX_LENGTH_COUNTRY_CODE {
		return p.InvalidPaymentError("billing_country_code")
	}
	region, err := language.ParseRegion(p.BillingCountryCode)
	if err != nil || !strings.EqualFold(region.String(), p.BillingCountryCode) {
		return p.InvalidPaymentError("billing_country_code")
	}
	if utf8.RuneCountInString(p.Currency) > model.MAX_LENGTH_CURRENCY_CODE {
		return p.InvalidPaymentError("currency")
	}
	if un, ok := currency.FromRegion(region); !ok || !strings.EqualFold(un.String(), p.Currency) {
		return p.InvalidPaymentError("currency")
	}
	if utf8.RuneCountInString(p.BillingCountryArea) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("billing_country_area")
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_FIRST_DIGITS {
		return p.InvalidPaymentError("cc_first_digits")
	}
	if len(p.CcFirstDigits) > MAX_LENGTH_CC_LAST_DIGITS {
		return p.InvalidPaymentError("cc_last_digits")
	}
	if *p.CcExpMonth < MIN_CC_EXP_MONTH || *p.CcExpMonth > MAX_CC_EXP_MONTH {
		return p.InvalidPaymentError("cc_exp_month")
	}
	if *p.CcExpYear < MIN_CC_EXP_YEAR {
		return p.InvalidPaymentError("cc_exp_year")
	}
	if len(p.PaymentMethodType) > MAX_LENGTH_PAYMENT_COMMON_256 {
		return p.InvalidPaymentError("payment_method_type")
	}

	return nil
}

// populate some fields if empty and perform some sanitizes
func (p *Payment) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}

	p.BillingEmail = model.NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = model.SanitizeUnicode(account.CleanNamePart(p.BillingFirstName, model.FirstName))
	p.BillingLastName = model.SanitizeUnicode(account.CleanNamePart(p.BillingLastName, model.LastName))

	if p.Total == nil {
		p.Total = &decimal.Zero
	}

	if p.CapturedAmount == nil {
		p.CapturedAmount = &decimal.Zero
	}

	if p.ChargeStatus == "" || !validChargeStatues.Contains(p.ChargeStatus) {
		p.ChargeStatus = DEFAULT_CHARGE_STATUS
	}

	p.CreateAt = model.GetMillis()
	p.UpdateAt = p.CreateAt
}

func (p *Payment) PreUpdate() {
	p.BillingEmail = model.NormalizeEmail(p.BillingEmail)
	p.BillingFirstName = model.SanitizeUnicode(account.CleanNamePart(p.BillingFirstName, model.FirstName))
	p.BillingLastName = model.SanitizeUnicode(account.CleanNamePart(p.BillingLastName, model.LastName))

	p.UpdateAt = model.GetMillis()
}

func (p *Payment) ToJson() string {
	b, _ := json.JSON.Marshal(p)

	return string(b)
}

func PaymentFromJson(data io.Reader) *Payment {
	var payment Payment
	err := json.JSON.NewDecoder(data).Decode(&payment)
	if err != nil {
		return nil
	}
	return &payment
}
