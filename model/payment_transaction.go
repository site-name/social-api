package model

import (
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
)

// default choices for payment transaction kind
const (
	EXTERNAL          = "external"
	AUTH              = "auth"
	CAPTURE           = "capture"
	CAPTURE_FAILED    = "capture_failed"
	ACTION_TO_CONFIRM = "action_to_confirm"
	VOID              = "void"
	PENDING_          = "pending"
	REFUND            = "refund"
	REFUND_ONGOING    = "refund_ongoing"
	REFUND_FAILED     = "refund_failed"
	REFUND_REVERSED   = "refund_reversed"
	CONFIRM           = "confirm"
	CANCEL            = "cancel"
)

var validTransactionKinds = StringArray([]string{
	EXTERNAL,
	AUTH,
	CAPTURE,
	CAPTURE_FAILED,
	ACTION_TO_CONFIRM,
	VOID,
	PENDING_,
	REFUND,
	REFUND_ONGOING,
	REFUND_FAILED,
	REFUND_REVERSED,
	CONFIRM,
	CANCEL,
})

const (
	TRANSACTION_KIND_MAX_LENGTH  = 25
	SEARCHABLE_KEY_MAX_LENGTH    = 512
	TRANSACTION_ERROR_MAX_LENGTH = 256
)

type PaymentTransaction struct {
	Id                 string               `json:"id"`
	CreateAt           int64                `json:"create_at"`
	PaymentID          string               `json:"payment_id"`
	Token              string               `json:"token"`
	Kind               string               `json:"kind"`
	IsSuccess          bool                 `json:"is_success"`
	ActionRequired     bool                 `json:"action_required"`
	ActionRequiredData StringMap            `json:"action_required_data"`
	Currency           string               `json:"currency"`
	Amount             *decimal.NullDecimal `json:"amount"`
	Error              *string              `json:"error"`
	CustomerID         *string              `json:"customer_id"`
	GatewayResponse    StringMap            `json:"gateway_response"`
	AlreadyProcessed   bool                 `json:"already_processed"`
	SearchableKey      *string              `json:"searchable_key"`
}

func (p *PaymentTransaction) String() string {
	return fmt.Sprintf(
		"Transaction(type=%s, is_success=%t, created=%d)",
		p.Kind,
		p.IsSuccess,
		p.CreateAt,
	)
}

func (p *PaymentTransaction) GetAmount() {
	panic("not implemented")
}

// Common method for creating app error for payment transaction
func InvalidPaymentTransactionErr(fieldName string, paymentTransactionID string) *AppError {
	id := fmt.Sprintf("model.payment_transaction.is_valid.%s.app_error", fieldName)
	var details string
	if paymentTransactionID != "" {
		details = "transaction_id=" + paymentTransactionID
	}

	return NewAppError("PaymentTransaction.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *PaymentTransaction) IsValid() *AppError {
	if !IsValidId(p.Id) {
		return InvalidPaymentTransactionErr("id", "")
	}
	if !IsValidId(p.PaymentID) {
		return InvalidPaymentTransactionErr("payment_id", p.Id)
	}
	if p.CustomerID != nil && !IsValidId(*p.CustomerID) {
		return InvalidPaymentTransactionErr("customer_id", p.Id)
	}
	if p.CreateAt == 0 {
		return InvalidPaymentTransactionErr("create_at", p.Id)
	}
	if len(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return InvalidPaymentTransactionErr("token", p.Id)
	}
	if len(p.Kind) > TRANSACTION_KIND_MAX_LENGTH {
		return InvalidPaymentTransactionErr("kind", p.Id)
	}
	if !validTransactionKinds.Contains(p.Kind) {
		return InvalidPaymentTransactionErr("kind", p.Id)
	}
	if len(p.Currency) > MAX_LENGTH_PAYMENT_CURRENCY_CODE {
		return InvalidPaymentTransactionErr("currency", p.Id)
	}
	if p.Error != nil && len(*p.Error) > TRANSACTION_ERROR_MAX_LENGTH {
		return InvalidPaymentTransactionErr("error", p.Id)
	}
	if p.SearchableKey != nil && utf8.RuneCountInString(*p.SearchableKey) > SEARCHABLE_KEY_MAX_LENGTH {
		return InvalidPaymentTransactionErr("searchable_key", p.Id)
	}
	if un, err := currency.ParseISO(p.Currency); err != nil || un.String() != p.Currency {
		return InvalidPaymentError("currency", p.Id)
	}
	if p.Amount == nil || !p.Amount.Valid {
		return InvalidPaymentError("amount", p.Id)
	}

	return nil
}

func (p *PaymentTransaction) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}

	if p.Amount == nil || !p.Amount.Valid {
		p.Amount = DEFAULT_DECIMAL_VALUE
	}

	p.CreateAt = GetMillis()
}

func (p *PaymentTransaction) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func PaymentTransactionFromJson(data io.Reader) *PaymentTransaction {
	var pmtr *PaymentTransaction
	json.JSON.NewDecoder(data).Decode(pmtr)
	return pmtr
}
