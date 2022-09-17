package model

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
)

// max lengths for some fields of giftcard
const (
	GiftcardCodeMaxLength             = 40
	GiftcardExpiryTypeMaxLength       = 32
	GiftcardExpiryPeriodTypeMaxLength = 32
	GiftcardTagMaxLength              = 255
)

// valid values for giftcard's ExpiryType
const (
	NeverExpire  = "never_expire"
	ExpiryPeriod = "expiry_period"
	ExpiryDate   = "expiry_date"
)

var GiftcardExpiryTypeMap = map[string]string{
	NeverExpire:  "Never expire",
	ExpiryPeriod: "Expiry period",
	ExpiryDate:   "Expiry date",
}

type GiftCard struct {
	Id                   string           `json:"id"`
	Code                 string           `json:"code"`          // unique, db_index
	CreatedByID          *string          `json:"created_by_id"` // foreign key User, ON DELETE SET NULL
	UsedByID             *string          `json:"used_by_id"`
	CreatedByEmail       *string          `json:"created_by_email"`
	UsedByEmail          *string          `json:"used_by_email"`
	CreateAt             int64            `json:"created_at"`
	StartDate            *time.Time       `json:"start_date"`
	ExpiryDate           *time.Time       `json:"expiry_date"`
	Tag                  *string          `json:"tag"`
	ProductID            *string          `json:"product_id"` // foreign key to Product
	LastUsedOn           *int64           `json:"last_used_on"`
	IsActive             *bool            `json:"is_active"` // default true
	Currency             string           `json:"currency"`
	InitialBalanceAmount *decimal.Decimal `json:"initial_balance_amount"` // default 0
	InitialBalance       *goprices.Money  `json:"initial_balance,omitempty" db:"-"`
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount"` // default 0
	CurrentBalance       *goprices.Money  `json:"current_balance,omitempty" db:"-"`
	ModelMetadata

	populatedNonDBFields bool `json:"-" db:"_"`
}

// GiftCardFilterOption is used to buil sql queries
type GiftCardFilterOption struct {
	ExpiryDate    squirrel.Sqlizer
	StartDate     squirrel.Sqlizer
	Code          squirrel.Sqlizer
	Currency      squirrel.Sqlizer
	CreatedByID   squirrel.Sqlizer
	CheckoutToken squirrel.Sqlizer // SELECT * FROM 'Giftcards' WHERE 'Id' IN (SELECT 'GiftcardID' FROM 'GiftCardCheckouts' WHERE 'GiftCardCheckouts.CheckoutID' ...)
	IsActive      *bool
	Distinct      bool // if true, SELECT DISTINCT

	SelectForUpdate bool // if true, concat `FOR UPDATE` to the end of SQL queries
}

type Giftcards []*GiftCard

func (g Giftcards) IDs() []string {
	res := []string{}
	for _, item := range g {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (gc *GiftCard) DisplayCode() string {
	return "****" + gc.Code[len(gc.Code)-4:]
}

// PopulateNonDbFields populates money fields for giftcard
func (gc *GiftCard) PopulateNonDbFields() {
	if gc.populatedNonDBFields {
		return
	}
	defer func() { gc.populatedNonDBFields = true }()

	if gc.InitialBalanceAmount == nil {
		gc.InitialBalanceAmount = &decimal.Zero
	}
	gc.InitialBalance = &goprices.Money{
		Amount:   *gc.InitialBalanceAmount,
		Currency: gc.Currency,
	}

	if gc.CurrentBalanceAmount == nil {
		gc.CurrentBalanceAmount = &decimal.Zero
	}
	gc.CurrentBalance = &goprices.Money{
		Amount:   *gc.CurrentBalanceAmount,
		Currency: gc.Currency,
	}
}

func (gc *GiftCard) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"gift_card.is_valid.%s.app_error",
		"gift_card_id=",
		"GiftCard.IsValid",
	)

	if !IsValidId(gc.Id) {
		return outer("id", nil)
	}
	if gc.CreatedByID != nil && !IsValidId(*gc.CreatedByID) {
		return outer("created_by_id", &gc.Id)
	}
	if gc.UsedByID != nil && !IsValidId(*gc.UsedByID) {
		return outer("used_by_id", &gc.Id)
	}
	if gc.CreatedByEmail != nil && len(*gc.CreatedByEmail) > USER_EMAIL_MAX_LENGTH {
		return outer("created_by_email", &gc.Id)
	}
	if gc.UsedByEmail != nil && len(*gc.UsedByEmail) > USER_EMAIL_MAX_LENGTH {
		return outer("used_by_email", &gc.Id)
	}
	if gc.Tag != nil && len(*gc.Tag) > GiftcardTagMaxLength {
		return outer("tag", &gc.Id)
	}
	if gc.ProductID != nil && !IsValidId(*gc.ProductID) {
		return outer("product_id", &gc.Id)
	}
	if gc.CreateAt == 0 {
		return outer("create_at", &gc.Id)
	}
	if gc.LastUsedOn != nil && *gc.LastUsedOn <= 0 {
		return outer("last_used_on", &gc.Id)
	}
	if len(gc.Code) > GiftcardCodeMaxLength {
		return outer("code", &gc.Id)
	}
	if unit, err := currency.ParseISO(gc.Currency); err != nil || !strings.EqualFold(unit.String(), gc.Currency) {
		return outer("currency", &gc.Id)
	}

	return nil
}

func (gc *GiftCard) PreSave() {
	if gc.Id == "" {
		gc.Id = NewId()
	}
	gc.CreateAt = GetMillis()

	gc.commonPre()
}

func (gc *GiftCard) commonPre() {
	if gc.CurrentBalance != nil {
		gc.CurrentBalanceAmount = &gc.CurrentBalance.Amount
	} else {
		gc.CurrentBalanceAmount = &decimal.Zero
	}

	if gc.InitialBalance != nil {
		gc.InitialBalanceAmount = &gc.InitialBalance.Amount
	} else {
		gc.InitialBalanceAmount = &decimal.Zero
	}

	if gc.IsActive == nil {
		gc.IsActive = NewBool(true)
	}

	if gc.Currency == "" {
		gc.Currency = DEFAULT_CURRENCY
	} else {
		gc.Currency = strings.ToUpper(gc.Currency)
	}
	if gc.StartDate == nil {
		today := util.StartOfDay(time.Now())
		gc.StartDate = &today
	}
}

func (gc *GiftCard) PreUpdate() {
	gc.commonPre()
}
