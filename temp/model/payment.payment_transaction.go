package model

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type TransactionKind string

// valid values for payment transaction's kind
const (
	TRANSACTION_KIND_EXTERNAL          TransactionKind = "external"
	TRANSACTION_KIND_AUTH              TransactionKind = "auth"
	TRANSACTION_KIND_CAPTURE           TransactionKind = "capture"
	TRANSACTION_KIND_CAPTURE_FAILED    TransactionKind = "capture_failed" // ?
	TRANSACTION_KIND_ACTION_TO_CONFIRM TransactionKind = "action_to_confirm"
	TRANSACTION_KIND_VOID              TransactionKind = "void"
	TRANSACTION_KIND_REFUND            TransactionKind = "refund"
	TRANSACTION_KIND_REFUND_ONGOING    TransactionKind = "refund_ongoing"
	TRANSACTION_KIND_REFUND_FAILED     TransactionKind = "refund_failed"   // ?
	TRANSACTION_KIND_REFUND_REVERSED   TransactionKind = "refund_reversed" // ?
	TRANSACTION_KIND_CONFIRM           TransactionKind = "confirm"
	TRANSACTION_KIND_CANCEL            TransactionKind = "cancel"
	TRANSACTION_KIND_PENDING           TransactionKind = "pending"
)

func (t TransactionKind) String() string {
	return string(t)
}

func (t TransactionKind) IsValid() bool {
	return TransactionKindString[t] != ""
}

var TransactionKindString = map[TransactionKind]string{
	TRANSACTION_KIND_EXTERNAL:          "External reference",
	TRANSACTION_KIND_AUTH:              "Authorization",
	TRANSACTION_KIND_PENDING:           "Pending", // transaction and payment share this value
	TRANSACTION_KIND_ACTION_TO_CONFIRM: "Action to confirm",
	TRANSACTION_KIND_REFUND:            "Refund",
	TRANSACTION_KIND_REFUND_ONGOING:    "Refund in progress",
	TRANSACTION_KIND_CAPTURE:           "Capture",
	TRANSACTION_KIND_VOID:              "Void",
	TRANSACTION_KIND_CONFIRM:           "Confirm",
	TRANSACTION_KIND_CANCEL:            "Cancel",
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

// coumn names for payment transaction table
const (
	TransactionColumnId                 = "Id"
	TransactionColumnCreateAt           = "CreateAt"
	TransactionColumnPaymentID          = "PaymentID"
	TransactionColumnToken              = "Token"
	TransactionColumnKind               = "Kind"
	TransactionColumnIsSuccess          = "IsSuccess"
	TransactionColumnActionRequired     = "ActionRequired"
	TransactionColumnActionRequiredData = "ActionRequiredData"
	TransactionColumnCurrency           = "Currency"
	TransactionColumnAmount             = "Amount"
	TransactionColumnError              = "Error"
	TransactionColumnCustomerID         = "CustomerID"
	TransactionColumnGatewayResponse    = "GatewayResponse"
	TransactionColumnAlreadyProcessed   = "AlreadyProcessed"
)

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
		p.Amount = GetPointerOfValue(decimal.Zero)
	}
	if p.ActionRequiredData == nil {
		p.ActionRequiredData = make(StringInterface)
	}
	if p.GatewayResponse == nil {
		p.GatewayResponse = make(StringInterface)
	}
}
