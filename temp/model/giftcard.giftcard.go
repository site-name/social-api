package model

import (
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
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
	CreateAt             int64            `json:"created_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
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

	Checkouts []*Checkout `json:"-" gorm:"many2many:GiftcardCheckouts"`
	Orders    Orders      `json:"-" gorm:"many2many:OrderGiftCards"`

	// NOTE: fields below are used for sorting purpose
	RelatedProductSlug     string `json:"-" gorm:"-"`
	RelatedProductName     string `json:"-" gorm:"-"`
	RelatedUsedByFirstName string `json:"-" gorm:"-"`
	RelatedUsedByLastName  string `json:"-" gorm:"-"`
}

func (c *GiftCard) BeforeCreate(_ *gorm.DB) error { c.PreSave(); return c.IsValid() }
func (c *GiftCard) BeforeUpdate(_ *gorm.DB) error { c.PreUpdate(); return c.IsValid() }
func (c *GiftCard) TableName() string             { return GiftcardTableName }

// GiftCardFilterOption is used to buil sql queries
type GiftCardFilterOption struct {
	Conditions squirrel.Sqlizer

	CheckoutToken squirrel.Sqlizer // Id IN (SELECT 'GiftcardID' FROM 'GiftcardCheckouts' WHERE 'GiftcardCheckouts.CheckoutID' ...)
	OrderID       squirrel.Sqlizer // INNER JOIN OrderGiftCards ON OrderGiftCards.GiftcardID = Giftcards.Id WHERE OrderGiftCards.OrderID ...

	Distinct        bool // if true, SELECT DISTINCT
	SelectForUpdate bool // if true, concat `FOR UPDATE` to the end of SQL queries. NOTE: only apply when Transaction is set
	Transaction     *gorm.DB

	AnnotateRelatedProductName bool
	AnnotateRelatedProductSlug bool
	AnnotateUsedByFirstName    bool
	AnnotateUsedByLastName     bool

	CountTotal              bool
	GraphqlPaginationValues GraphqlPaginationValues
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
	if gc.InitialBalanceAmount == nil {
		gc.InitialBalanceAmount = GetPointerOfValue(decimal.Zero)
	}
	gc.InitialBalance = &goprices.Money{
		Amount:   *gc.InitialBalanceAmount,
		Currency: gc.Currency,
	}

	if gc.CurrentBalanceAmount == nil || gc.CurrentBalanceAmount.LessThan(decimal.Zero) {
		gc.CurrentBalanceAmount = GetPointerOfValue(decimal.Zero)
	}
	gc.CurrentBalance = &goprices.Money{
		Amount:   *gc.CurrentBalanceAmount,
		Currency: gc.Currency,
	}
}

func (gc *GiftCard) IsValid() *AppError {
	if gc.CreatedByID != nil && !IsValidId(*gc.CreatedByID) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.created_by.app_error", nil, "please provide valid created by id", http.StatusBadRequest)
	}
	if gc.UsedByID != nil && !IsValidId(*gc.UsedByID) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.usedby_id.app_error", nil, "please provide valid used by by id", http.StatusBadRequest)
	}
	if gc.CreatedByEmail != nil && !IsValidEmail(*gc.CreatedByEmail) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.created_by_email.app_error", nil, "please provide valid created by email", http.StatusBadRequest)
	}
	if gc.UsedByEmail != nil && !IsValidEmail(*gc.UsedByEmail) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.used_by_email.app_error", nil, "please provide valid used by email", http.StatusBadRequest)
	}
	if gc.ProductID != nil && !IsValidId(*gc.ProductID) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if gc.LastUsedOn != nil && *gc.LastUsedOn <= 0 {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.last_used_on.app_error", nil, "please provide valid last used on", http.StatusBadRequest)
	}
	if _, err := currency.ParseISO(gc.Currency); err != nil {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if !PromoCodeRegex.MatchString(gc.Code) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.code.app_error", nil, "Code must look like DF6F-HGYG-78TU", http.StatusBadRequest)
	}
	if gc.StartDate != nil && gc.StartDate.IsZero() {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.start_date.app_error", nil, "please provide valid start date", http.StatusBadRequest)
	}
	if gc.ExpiryDate != nil && gc.ExpiryDate.IsZero() {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.expiry_date.app_error", nil, "please provide valid expiry date", http.StatusBadRequest)
	}
	if gc.StartDate != nil && gc.ExpiryDate != nil && gc.StartDate.After(*gc.ExpiryDate) {
		return NewAppError("Giftcard.IsValid", "model.gift_card.is_valid.dates.app_error", nil, "start date must be before expiry date", http.StatusBadRequest)
	}

	return nil
}

func (gc *GiftCard) PreSave() {
	gc.commonPre()
}

func (gc *GiftCard) commonPre() {
	if gc.CurrentBalanceAmount == nil {
		gc.CurrentBalanceAmount = GetPointerOfValue(decimal.Zero)
	}

	if gc.InitialBalanceAmount == nil {
		gc.InitialBalanceAmount = GetPointerOfValue(decimal.Zero)
	}

	if gc.IsActive == nil {
		gc.IsActive = GetPointerOfValue(true)
	}

	if gc.Currency == "" {
		gc.Currency = DEFAULT_CURRENCY
	}
	if gc.StartDate == nil {
		today := util.StartOfDay(time.Now())
		gc.StartDate = &today
	}
	if gc.Code == "" {
		gc.Code = NewPromoCode()
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
