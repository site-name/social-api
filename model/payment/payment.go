package payment

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util"
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

	// Payment Gateways
	GATE_WAY_MANUAL = "manual"
)

const (
	// some fields may have max length of 256 in common
	MAX_LENGTH_PAYMENT_COMMON_256 = 256
)

// Choices for charge status
const (
	NOT_CHARGED        = "not-charged"
	PENDING            = "pending"
	PARTIALLY_CHARGED  = "partially-charged"
	FULLY_CHARGED      = "fully-charged"
	PARTIALLY_REFUNDED = "partially-refunded"
	FULLY_REFUNDED     = "fully-refunded"
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
	Id                 string                `json:"id"`
	GateWay            string                `json:"gate_way"`
	IsActive           bool                  `json:"is_active"`
	ToConfirm          bool                  `json:"to_confirm"`
	CreateAt           int64                 `json:"create_at"`
	UpdateAt           int64                 `json:"update_at"`
	ChargeStatus       string                `json:"charge_status"`
	Token              string                `json:"token"`
	Total              *decimal.Decimal      `json:"total"`
	CapturedAmount     *decimal.Decimal      `json:"captured_amount"`
	Currency           string                `json:"currency"`
	CheckoutID         string                `json:"checkout_id"`
	OrderID            string                `json:"order_id"`
	BillingEmail       string                `json:"billing_email"`
	BillingFirstName   string                `json:"billing_first_name"`
	BillingLastName    string                `json:"billing_last_name"`
	BillingCompanyName string                `json:"billing_company_name"`
	BillingAddress1    string                `json:"billing_address_1"`
	BillingAddress2    string                `json:"billing_address_2"`
	BillingCity        string                `json:"billing_city"`
	BillingCityArea    string                `json:"billing_city_area"`
	BillingPostalCode  string                `json:"billing_postal_code"`
	BillingCountryCode string                `json:"billing_country_code"`
	BillingCountryArea string                `json:"billing_country_area"`
	CcFirstDigits      string                `json:"cc_first_digits"`
	CcLastDigits       string                `json:"cc_last_digits"`
	CcBranh            string                `json:"cc_brand"`
	CcExpMonth         *uint8                `json:"cc_exp_month"`
	CcExpYear          *uint16               `json:"cc_exp_year"`
	PaymentMethodType  string                `json:"payment_method_type"`
	CustomerIpAddress  *string               `json:"customer_ip_address"`
	ExtraData          string                `json:"extra_data"`
	ReturnUrl          *string               `json:"return_url_url"`
	Transactions       []*PaymentTransaction `json:"transactions" db:"-"`
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

// Retrieve the maximum capture possible.
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

func (p *Payment) IsAuthorized() bool {
	for _, tx := range p.Transactions {
		if tx.Kind == AUTH && tx.IsSuccess && !tx.ActionRequired {
			return true
		}
	}

	return false
}

func (p *Payment) CanVoid() bool {
	return p.IsActive && p.IsNotCharged() && p.IsAuthorized()
}

func (p *Payment) CanRefund() bool {
	canRefundChargeStatuses := []string{
		PARTIALLY_CHARGED,
		FULLY_CHARGED,
		PARTIALLY_REFUNDED,
	}

	return p.IsActive && util.StringInSlice(p.ChargeStatus, canRefundChargeStatuses)
}

func (p *Payment) CanConfirm() bool {
	return p.IsActive && p.IsNotCharged()
}

func (p *Payment) IsManual() bool {
	return p.GateWay == GATE_WAY_MANUAL
}

func (p *Payment) GetTotal() *model.Money {
	return &model.Money{
		Amount:   p.Total,
		Currency: p.Currency,
	}
}

// get most recent transaction by comparing their created time
func (p *Payment) GetLastTransaction() *PaymentTransaction {
	var maxTime int64 = 0
	var tran *PaymentTransaction

	for _, tx := range p.Transactions {
		if tx.CreateAt > maxTime {
			maxTime = tx.CreateAt
			tran = tx
		}
	}

	return tran
}

func (p *Payment) GetCapturedAmount() *model.Money {
	return &model.Money{
		Amount:   p.CapturedAmount,
		Currency: p.Currency,
	}
}

func (p *Payment) GetAuthorizedAmount() *model.Money {
	money := &model.Money{
		Amount:   &decimal.Zero,
		Currency: p.Currency,
	}

	for _, tx := range p.Transactions {
		if tx.Kind == CAPTURE && tx.IsSuccess {
			return money
		}
	}

	for _, tx := range p.Transactions {
		if tx.Kind == AUTH && tx.IsSuccess && !tx.ActionRequired {
			addedAmount := money.Amount.Add(*tx.Amount)
			money.Amount = &addedAmount
		}
	}

	return money
}

// Check if input from user is valid or not
func (p *Payment) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.payment.is_valid.%s.app_error",
		"payment_id=",
		"Payment.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.OrderID) {
		return outer("order_id", &p.Id)
	}
	if !model.IsValidId(p.CheckoutID) {
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
	if !model.IsValidEmail(p.BillingEmail) {
		return outer("billing_email", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingFirstName) > account.ADDRESS_FIRST_NAME_MAX_LENGTH || !account.IsValidNamePart(p.BillingFirstName, model.FirstName) {
		return outer("billing_first_name", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingLastName) > account.ADDRESS_LAST_NAME_MAX_LENGTH || !account.IsValidNamePart(p.BillingLastName, model.LastName) {
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
	if utf8.RuneCountInString(p.BillingCityArea) > account.ADDRESS_CITY_AREA_MAX_LENGTH {
		return outer("billing_city_area", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingPostalCode) > account.ADDRESS_POSTAL_CODE_MAX_LENGTH {
		return outer("billing_postal_code", &p.Id)
	}
	if utf8.RuneCountInString(p.BillingCountryCode) > model.MAX_LENGTH_COUNTRY_CODE {
		return outer("billing_country_code", &p.Id)
	}

	// make sure country code and currency code are match:
	region, err := language.ParseRegion(p.BillingCountryCode)
	if err != nil || !strings.EqualFold(region.String(), p.BillingCountryCode) {
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

	if ChargeStatuString[strings.ToLower(p.ChargeStatus)] == "" {
		p.ChargeStatus = NOT_CHARGED
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
	return model.ModelToJson(p)
}

func PaymentFromJson(data io.Reader) *Payment {
	var payment Payment
	model.ModelFromJson(&payment, data)
	return &payment
}