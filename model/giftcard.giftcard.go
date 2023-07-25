package model

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

// valid values for giftcard's ExpiryType
const (
	NeverExpire  = "never_expire"
	ExpiryPeriod = "expiry_period"
	ExpiryDate   = "expiry_date"
)

type GiftCard struct {
	Id                   string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Code                 string           `json:"code" gorm:"type:varchar(16);column:Code"`          // unique, db_index, looks like ABCD-EFGH-IJKL
	CreatedByID          *string          `json:"created_by_id" gorm:"type:uuid;column:CreatedByID"` // foreign key User, ON DELETE SET NULL
	UsedByID             *string          `json:"used_by_id" gorm:"type:uuid;column:UsedByID"`
	CreatedByEmail       *string          `json:"created_by_email" gorm:"type:varchar(128);column:CreatedByEmail"`
	UsedByEmail          *string          `json:"used_by_email" gorm:"type:varchar(128);column:UsedByEmail"`
	CreateAt             int64            `json:"created_at" gorm:"type:bigint;column:CreateAt"`
	StartDate            *time.Time       `json:"start_date" gorm:"column:StartDate"`
	ExpiryDate           *time.Time       `json:"expiry_date" gorm:"column:ExpiryDate"`
	Tag                  *string          `json:"tag" gorm:"type:varchar(255);column:Tag"`
	ProductID            *string          `json:"product_id" gorm:"type:uuid;column:ProductID"` // foreign key to Product
	LastUsedOn           *int64           `json:"last_used_on" gorm:"type:bigint;column:LastUsedOn"`
	IsActive             *bool            `json:"is_active" gorm:"column:IsActive"`                                    // default true
	Currency             string           `json:"currency" gorm:"type:varchar(3);column:Currency"`                     // UPPER cased
	InitialBalanceAmount *decimal.Decimal `json:"initial_balance_amount" gorm:"default:0;column:InitialBalanceAmount"` // default 0
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount" gorm:"default:0;column:CurrentBalanceAmount"` // default 0
	ModelMetadata

	CurrentBalance *goprices.Money `json:"current_balance,omitempty" gorm:"-"`
	InitialBalance *goprices.Money `json:"initial_balance,omitempty" gorm:"-"`

	populatedNonDBFields bool `db:"-"`

	Checkouts []*Checkout `json:"-" gorm:"many2many:GiftcardCheckouts"`
	Orders    Orders      `json:"-" gorm:"many2many:OrderGiftCards"`
}

func (c *GiftCard) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *GiftCard) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *GiftCard) TableName() string             { return GiftcardTableName }

// GiftCardFilterOption is used to buil sql queries
type GiftCardFilterOption struct {
	Conditions squirrel.Sqlizer

	CheckoutToken squirrel.Sqlizer // SELECT * FROM 'Giftcards' WHERE 'Id' IN (SELECT 'GiftcardID' FROM 'GiftCardCheckouts' WHERE 'GiftCardCheckouts.CheckoutID' ...)
	OrderID       squirrel.Sqlizer // INNER JOIN OrderGiftCards ON OrderGiftCards.GiftcardID = Giftcards.Id WHERE OrderGiftCards.OrderID ...

	Distinct        bool // if true, SELECT DISTINCT
	OrderBy         string
	SelectForUpdate bool // if true, concat `FOR UPDATE` to the end of SQL queries
}

type Giftcards []*GiftCard

func (gs Giftcards) IDs() []string {
	return lo.Map(gs, func(g *GiftCard, _ int) string { return g.Id })
}

func (gc *GiftCard) DisplayCode() string {
	if len(gc.Code) <= 4 {
		return gc.Code
	}
	return "****" + gc.Code[len(gc.Code)-4:]
}

// PopulateNonDbFields populates money fields for giftcard
func (gc *GiftCard) PopulateNonDbFields() {
	if gc.populatedNonDBFields {
		return
	}
	gc.populatedNonDBFields = true

	if gc.InitialBalanceAmount == nil {
		gc.InitialBalanceAmount = &decimal.Zero
	}
	gc.InitialBalance = &goprices.Money{
		Amount:   *gc.InitialBalanceAmount,
		Currency: gc.Currency,
	}

	if gc.CurrentBalanceAmount == nil || gc.CurrentBalanceAmount.LessThan(decimal.Zero) {
		gc.CurrentBalanceAmount = &decimal.Zero
	}
	gc.CurrentBalance = &goprices.Money{
		Amount:   *gc.CurrentBalanceAmount,
		Currency: gc.Currency,
	}
}

func (gc *GiftCard) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.gift_card.is_valid.%s.app_error",
		"gift_card_id=",
		"GiftCard.IsValid",
	)

	if gc.CreatedByID != nil && !IsValidId(*gc.CreatedByID) {
		return outer("created_by_id", &gc.Id)
	}
	if gc.UsedByID != nil && !IsValidId(*gc.UsedByID) {
		return outer("used_by_id", &gc.Id)
	}
	if gc.CreatedByEmail != nil && !IsValidEmail(*gc.CreatedByEmail) {
		return outer("created_by_email", &gc.Id)
	}
	if gc.UsedByEmail != nil && !IsValidEmail(*gc.UsedByEmail) {
		return outer("used_by_email", &gc.Id)
	}
	if gc.ProductID != nil && !IsValidId(*gc.ProductID) {
		return outer("product_id", &gc.Id)
	}
	if gc.LastUsedOn != nil && *gc.LastUsedOn <= 0 {
		return outer("last_used_on", &gc.Id)
	}
	if _, err := currency.ParseISO(gc.Currency); err != nil {
		return outer("currency", &gc.Id)
	}

	return nil
}

func (gc *GiftCard) PreSave() {
	if gc.Code == "" {
		rawString := NewRandomString(16)
		gc.Code = fmt.Sprintf("%s-%s-%s-%s", rawString[:4], rawString[4:8], rawString[8:12], rawString[12:])
	}
	gc.commonPre()
}

func (gc *GiftCard) commonPre() {
	if gc.CurrentBalanceAmount == nil {
		gc.CurrentBalanceAmount = &decimal.Zero
	}

	if gc.InitialBalanceAmount == nil {
		gc.InitialBalanceAmount = &decimal.Zero
	}

	if gc.IsActive == nil {
		gc.IsActive = NewPrimitive(true)
	}

	if gc.Currency == "" {
		gc.Currency = DEFAULT_CURRENCY
	}
	if gc.StartDate == nil {
		today := util.StartOfDay(time.Now())
		gc.StartDate = &today
	}
}

func (gc *GiftCard) PreUpdate() {
	gc.commonPre()
}

// NOTE: If you want co use *InitialBalance* or *CurrentBalance* after DeepCopy,
// you have to call PopulateNonDbFields() again.
func (gc *GiftCard) DeepCopy() *GiftCard {
	return &GiftCard{
		Id:                   gc.Id,
		Code:                 gc.Code,
		CreatedByID:          CopyPointer(gc.CreatedByID),
		UsedByID:             CopyPointer(gc.UsedByID),
		CreatedByEmail:       CopyPointer(gc.CreatedByEmail),
		UsedByEmail:          CopyPointer(gc.UsedByEmail),
		CreateAt:             gc.CreateAt,
		StartDate:            CopyPointer(gc.StartDate),
		ExpiryDate:           CopyPointer(gc.ExpiryDate),
		Tag:                  CopyPointer(gc.Tag),
		ProductID:            CopyPointer(gc.ProductID),
		LastUsedOn:           CopyPointer(gc.LastUsedOn),
		IsActive:             CopyPointer(gc.IsActive),
		InitialBalanceAmount: CopyPointer(gc.InitialBalanceAmount),
		CurrentBalanceAmount: CopyPointer(gc.CurrentBalanceAmount),
		ModelMetadata:        gc.ModelMetadata.DeepCopy(),
	}
}
