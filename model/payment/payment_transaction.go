package payment

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
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

var TransactionKindString = map[string]string{
	EXTERNAL:          "External reference",
	AUTH:              "Authorization",
	PENDING:           "Pending",
	ACTION_TO_CONFIRM: "Action to confirm",
	REFUND:            "Refund",
	REFUND_ONGOING:    "Refund in progress",
	CAPTURE:           "Capture",
	VOID:              "Void",
	CONFIRM:           "Confirm",
	CANCEL:            "Cancel",
}

const (
	TRANSACTION_KIND_MAX_LENGTH  = 25
	SEARCHABLE_KEY_MAX_LENGTH    = 512
	TRANSACTION_ERROR_MAX_LENGTH = 256
)

// Represents a single payment operation.
//
// Transaction is an attempt to transfer money between your store
// and your customers, with a chosen payment method.
type PaymentTransaction struct {
	Id                 string           `json:"id"`
	CreateAt           int64            `json:"create_at"`
	PaymentID          string           `json:"payment_id"`
	Token              string           `json:"token"`
	Kind               string           `json:"kind"`
	IsSuccess          bool             `json:"is_success"`
	ActionRequired     bool             `json:"action_required"`
	ActionRequiredData model.StringMap  `json:"action_required_data"`
	Currency           string           `json:"currency"`
	Amount             *decimal.Decimal `json:"amount"`
	Error              *string          `json:"error"`
	CustomerID         *string          `json:"customer_id"`
	GatewayResponse    model.StringMap  `json:"gateway_response"`
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

func (p *PaymentTransaction) GetAmount() *model.Money {
	return &model.Money{
		Amount:   p.Amount,
		Currency: p.Currency,
	}
}

func (p *PaymentTransaction) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.payment_transaction.is_valid.%s.app_error",
		"transaction_id=",
		"PaymentTransaction.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.PaymentID) {
		return outer("payment_id", &p.Id)
	}
	if p.CustomerID != nil && !model.IsValidId(*p.CustomerID) {
		return outer("customer_id", &p.Id)
	}
	if p.CreateAt == 0 {
		return outer("create_at", &p.Id)
	}
	if len(p.Token) > MAX_LENGTH_PAYMENT_TOKEN {
		return outer("token", &p.Id)
	}
	if TransactionKindString[strings.ToLower(p.Kind)] == "" {
		return outer("kind", &p.Id)
	}
	if p.Error != nil && len(*p.Error) > TRANSACTION_ERROR_MAX_LENGTH {
		return outer("error", &p.Id)
	}
	if p.SearchableKey != nil && utf8.RuneCountInString(*p.SearchableKey) > SEARCHABLE_KEY_MAX_LENGTH {
		return outer("searchable_key", &p.Id)
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
		p.Id = model.NewId()
	}

	if p.Amount == nil {
		p.Amount = &decimal.Zero
	}

	p.CreateAt = model.GetMillis()
}

func (p *PaymentTransaction) ToJson() string {
	return model.ModelToJson(p)
}

func PaymentTransactionFromJson(data io.Reader) *PaymentTransaction {
	var pmtr PaymentTransaction
	model.ModelFromJson(&pmtr, data)
	return &pmtr
}