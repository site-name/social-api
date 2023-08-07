package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// max lengths for Invoice
const (
	INVOICE_NUMBER_MAX_LENGTH       = 255
	INVOICE_EXTERNAL_URL_MAX_LENGTH = 2048
	INVOICE_STATUS_MAX_LENGTH       = 50
	INVOICE_MESSAGE_MAX_LENGTH      = 255
)

type Invoice struct {
	Id          string  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	OrderID     *string `json:"order_id" gorm:"type:uuid;column:OrderID"`
	Number      string  `json:"number" gorm:"type:varchar(255);column:Number"`
	CreateAt    int64   `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	ExternalUrl string  `json:"external_url" gorm:"type:varchar(2048);column:ExternalUrl"`
	Status      string  `json:"status" gorm:"type:varchar(50);column:Status"`
	Message     string  `json:"message" gorm:"type:varchar(255);column:Message"`
	UpdateAt    int64   `json:"update_at" gorm:"type:bigint;column:UpdateAt;autoCreateTime:milli;autoUpdateTime:milli"`
	ModelMetadata

	order *Order `gorm:"-"`
}

func (c *Invoice) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *Invoice) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *Invoice) TableName() string             { return InvoiceTableName }
func (i *Invoice) GetOrder() *Order              { return i.order }
func (i *Invoice) SetOrder(o *Order)             { i.order = o }

type InvoiceFilterOptions struct {
	Conditions squirrel.Sqlizer

	Limit              uint64
	SelectRelatedOrder bool
}

func (i *Invoice) IsValid() *AppError {
	if i.OrderID != nil && !IsValidId(*i.OrderID) {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.order_id.app_error", nil, "please provide valid order id", http.StatusBadRequest)
	}
	if !ALL_JOB_STATUSES.Contains(i.Status) {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.status.app_error", nil, "please provide valid status id", http.StatusBadRequest)
	}

	return nil
}

func (i *Invoice) DeepCopy() *Invoice {
	if i == nil {
		return nil
	}

	res := *i
	if i.OrderID != nil {
		res.OrderID = NewPrimitive(*i.OrderID)
	}
	if i.order != nil {
		res.order = i.order.DeepCopy()
	}
	res.ModelMetadata = i.ModelMetadata.DeepCopy()
	return &res
}
