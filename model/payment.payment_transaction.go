package model

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type TransactionKind string

// valid values for payment transaction's kind
const (
	EXTERNAL          TransactionKind = "external"
	AUTH              TransactionKind = "auth"
	CAPTURE           TransactionKind = "capture"
	CAPTURE_FAILED    TransactionKind = "capture_failed" // ?
	ACTION_TO_CONFIRM TransactionKind = "action_to_confirm"
	VOID              TransactionKind = "void"
	REFUND            TransactionKind = "refund"
	REFUND_ONGOING    TransactionKind = "refund_ongoing"
	REFUND_FAILED     TransactionKind = "refund_failed"   // ?
	REFUND_REVERSED   TransactionKind = "refund_reversed" // ?
	CONFIRM           TransactionKind = "confirm"
	CANCEL            TransactionKind = "cancel"
	PENDING_          TransactionKind = "pending"
)

func (t TransactionKind) String() string {
	return string(t)
}

func (t TransactionKind) IsValid() bool {
	return TransactionKindString[t] != ""
}

var TransactionKindString = map[TransactionKind]string{
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

// Represents a single payment operation.
// Transaction is an attempt to transfer money between your store
// and your customers, with a chosen payment method.
type PaymentTransaction struct {
	Id                 string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt           int64            `json:"create_at" gorm:"type:bigint;column:CreateAt"` // NOT editable
	PaymentID          string           `json:"payment_id" gorm:"type:uuid;column:PaymentID"`
	Token              string           `json:"token" gorm:"type:varchar(512);column:Token"`
	Kind               TransactionKind  `json:"kind" gorm:"type:varchar(25);column:Kind"`
	IsSuccess          bool             `json:"is_success" gorm:"column:IsSuccess"`
	ActionRequired     bool             `json:"action_required" gorm:"column:ActionRequired"`
	ActionRequiredData StringInterface  `json:"action_required_data" gorm:"type:jsonb;column:ActionRequiredData"`
	Currency           string           `json:"currency" gorm:"type:varchar(5);column:Currency"`
	Amount             *decimal.Decimal `json:"amount" gorm:"default:0;column:Amount;type:decimal(12,3)"` // DEFAULT decimal(0)
	Error              *string          `json:"error" gorm:"type:varchar(256);column:Error"`
	CustomerID         *string          `json:"customer_id" gorm:"type:varchar(256);column:CustomerID"`
	GatewayResponse    StringInterface  `json:"gateway_response" gorm:"type:jsonb;column:GatewayResponse"`
	AlreadyProcessed   bool             `json:"already_processed" gorm:"column:AlreadyProcessed"`
}

func (c *PaymentTransaction) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PaymentTransaction) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PaymentTransaction) TableName() string             { return TransactionTableName }

// PaymentTransactionFilterOpts contains options for filter payment's transactions
type PaymentTransactionFilterOpts struct {
	Conditions squirrel.Sqlizer
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
	if !IsValidId(p.PaymentID) {
		return NewAppError("Transaction.IsValid", "model.transaction.is_valid.payment_id.app_error", nil, "please provide valid payment id", http.StatusBadRequest)
	}
	if un, err := currency.ParseISO(p.Currency); err != nil || !strings.EqualFold(un.String(), p.Currency) {
		return NewAppError("Transaction.IsValid", "model.transaction.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if !p.Kind.IsValid() {
		return NewAppError("Transaction.IsValid", "model.transaction.is_valid.kind.app_error", nil, "please provide valid kind", http.StatusBadRequest)
	}

	return nil
}

func (p *PaymentTransaction) commonPre() {
	if p.Error != nil {
		*p.Error = SanitizeUnicode(*p.Error)
	}
	if p.Amount == nil || p.Amount.LessThanOrEqual(decimal.Zero) {
		p.Amount = &decimal.Zero
	}
	if p.ActionRequiredData == nil {
		p.ActionRequiredData = make(StringInterface)
	}
	if p.GatewayResponse == nil {
		p.GatewayResponse = make(StringInterface)
	}
}
