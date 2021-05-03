package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Id                 string           `json:"id"`
	CreateAt           int64            `json:"create_at"`
	PaymentID          string           `json:"payment_id"`
	Token              string           `json:"token"`
	Kind               string           `json:"kind"`
	IsSuccess          bool             `json:"is_success"`
	ActionRequired     bool             `json:"action_required"`
	ActionRequiredData StringMap        `json:"action_required_data"`
	Currency           string           `json:"currency"`
	Amount             *decimal.Decimal `json:"amount"`
	Error              *string          `json:"error"`
	CustomerID         *string          `json:"customer_id"`
	GatewayResponse    StringMap        `json:"gateway_response"`
	AlreadyProcessed   bool             `json:"already_processed"`
	SearchableKey      *string          `json:"searchable_key"`
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
func (p *PaymentTransaction) InvalidPaymentTransactionErr(fieldName string) *AppError {
	id := fmt.Sprintf("model.payment_transaction.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "transaction_id=" + p.Id
	}

	return NewAppError("PaymentTransaction.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *PaymentTransaction) IsValid() *AppError {
	if !IsValidId(p.Id) {
		return p.InvalidPaymentTransactionErr("id")
	}
	if !IsValidId(p.PaymentID) {
		return p.InvalidPaymentTransactionErr("payment_id")
	}
	if p.CustomerID != nil && !IsValidId(*p.CustomerID) {
		return p.InvalidPaymentTransactionErr("customer_id")
	}
	if p.CreateAt == 0 {
		return p.InvalidPaymentTransactionErr("create_at")
	}
	if len(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return p.InvalidPaymentTransactionErr("token")
	}
	if len(p.Kind) > TRANSACTION_KIND_MAX_LENGTH {
		return p.InvalidPaymentTransactionErr("kind")
	}
	if !validTransactionKinds.Contains(p.Kind) {
		return p.InvalidPaymentTransactionErr("kind")
	}
	if len(p.Currency) > MAX_LENGTH_CURRENCY_CODE {
		return p.InvalidPaymentTransactionErr("currency")
	}
	if p.Error != nil && len(*p.Error) > TRANSACTION_ERROR_MAX_LENGTH {
		return p.InvalidPaymentTransactionErr("error")
	}
	if p.SearchableKey != nil && utf8.RuneCountInString(*p.SearchableKey) > SEARCHABLE_KEY_MAX_LENGTH {
		return p.InvalidPaymentTransactionErr("searchable_key")
	}
	if un, err := currency.ParseISO(p.Currency); err != nil || un.String() != p.Currency {
		return p.InvalidPaymentTransactionErr("currency")
	}
	if p.Amount == nil {
		return p.InvalidPaymentTransactionErr("amount")
	}

	return nil
}

func (p *PaymentTransaction) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}

	if p.Amount == nil {
		p.Amount = &decimal.Zero
	}

	p.CreateAt = GetMillis()
}

func (p *PaymentTransaction) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func PaymentTransactionFromJson(data io.Reader) *PaymentTransaction {
	var pmtr PaymentTransaction
	err := json.JSON.NewDecoder(data).Decode(&pmtr)
	if err != nil {
		return nil
	}
	return &pmtr
}
