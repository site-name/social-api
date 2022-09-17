package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
)

// valid values for payment transaction's kind
const (
	EXTERNAL          = "external"
	AUTH              = "auth"
	CAPTURE           = "capture"
	CAPTURE_FAILED    = "capture_failed" // ?
	ACTION_TO_CONFIRM = "action_to_confirm"
	VOID              = "void"
	REFUND            = "refund"
	REFUND_ONGOING    = "refund_ongoing"
	REFUND_FAILED     = "refund_failed"   // ?
	REFUND_REVERSED   = "refund_reversed" // ?
	CONFIRM           = "confirm"
	CANCEL            = "cancel"
	PENDING_          = "pending"
)

var TransactionKindString = map[string]string{
	EXTERNAL:          "External reference",
	AUTH:              "Authorization",
	PENDING_:          "Pending", // transaction and payment share this value
	ACTION_TO_CONFIRM: "Action to confirm",
	REFUND:            "Refund",
	REFUND_ONGOING:    "Refund in progress",
	CAPTURE:           "Capture",
	VOID:              "Void",
	CONFIRM:           "Confirm",
	CANCEL:            "Cancel",
}

// max lengths for some of payment transaction's fields
const (
	TRANSACTION_KIND_MAX_LENGTH        = 25
	TRANSACTION_ERROR_MAX_LENGTH       = 256
	TRANSACTION_CUSTOMER_ID_MAX_LENGTH = 256
)

// Represents a single payment operation.
// Transaction is an attempt to transfer money between your store
// and your customers, with a chosen payment method.
type PaymentTransaction struct {
	Id                 string           `json:"id"`
	CreateAt           int64            `json:"create_at"` // NOT editable
	PaymentID          string           `json:"payment_id"`
	Token              string           `json:"token"`
	Kind               string           `json:"kind"`
	IsSuccess          bool             `json:"is_success"`
	ActionRequired     bool             `json:"action_required"`
	ActionRequiredData StringMap        `json:"action_required_data"`
	Currency           string           `json:"currency"`
	Amount             *decimal.Decimal `json:"amount"` // DEFAULT decimal(0)
	Error              *string          `json:"error"`
	CustomerID         *string          `json:"customer_id"`
	GatewayResponse    StringMap        `json:"gateway_response"`
	AlreadyProcessed   bool             `json:"already_processed"`
}

// PaymentTransactionFilterOpts contains options for filter payment's transactions
type PaymentTransactionFilterOpts struct {
	Id             squirrel.Sqlizer
	PaymentID      squirrel.Sqlizer
	Kind           squirrel.Sqlizer
	ActionRequired *bool
	IsSuccess      *bool
}

func (p *PaymentTransaction) String() string {
	return fmt.Sprintf(
		"Transaction(type=%s, is_success=%t, created=%d)",
		p.Kind,
		p.IsSuccess,
		p.CreateAt,
	)
}

func (p *PaymentTransaction) GetAmount() *goprices.Money {
	return &goprices.Money{
		Amount:   *p.Amount,
		Currency: p.Currency,
	}
}

func (p *PaymentTransaction) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"payment_transaction.is_valid.%s.app_error",
		"transaction_id=",
		"PaymentTransaction.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.PaymentID) {
		return outer("payment_id", &p.Id)
	}
	// NOTE: not sure CustomerID is uuid or not
	if p.CustomerID != nil && len(*p.CustomerID) > TRANSACTION_CUSTOMER_ID_MAX_LENGTH {
		return outer("customer_id", &p.Id)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if len(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return outer("token", &p.Id)
	}
	if len(p.Kind) > TRANSACTION_KIND_MAX_LENGTH {
		return outer("kind", &p.Id)
	}
	if p.Error != nil && utf8.RuneCountInString(*p.Error) > TRANSACTION_ERROR_MAX_LENGTH {
		return outer("error", &p.Id)
	}
	if un, err := currency.ParseISO(p.Currency); err != nil || !strings.EqualFold(un.String(), p.Currency) {
		return outer("currency", &p.Id)
	}
	if p.Amount == nil {
		return outer("amount", &p.Id)
	}

	return nil
}

func (p *PaymentTransaction) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.CreateAt = GetMillis()

	if p.ActionRequiredData == nil {
		p.ActionRequiredData = make(StringMap)
	}
	if p.Error != nil {
		*p.Error = SanitizeUnicode(*p.Error)
	}
	p.commonPre()
}

func (p *PaymentTransaction) commonPre() {
	if p.Amount == nil || p.Amount.LessThanOrEqual(decimal.Zero) {
		p.Amount = &decimal.Zero
	}
	if p.ActionRequiredData == nil {
		p.ActionRequiredData = make(StringMap)
	}
	if p.GatewayResponse == nil {
		p.GatewayResponse = make(StringMap)
	}
}

func (p *PaymentTransaction) PreUpdate() {
	p.commonPre()
}

func (p *PaymentTransaction) ToJSON() string {
	return ModelToJson(p)
}
