package giftcard

import (
	"strings"
	"time"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
)

// max lengths for some fields of giftcard
const (
	GiftcardCodeMaxLength             = 16
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
	ExpiryType           string           `json:"expiry_type"`
	ExpiryPeriodType     *string          `json:"expiry_period_type"`
	ExpiryPeriod         *int             `json:"expiry_period"`
	Tag                  *string          `json:"tag"`
	ProductID            *string          `json:"product_id"` // foreign key to Product
	LastUsedOn           *int64           `json:"last_used_on"`
	IsActive             *bool            `json:"is_active"` // default true
	Currency             string           `json:"currency"`
	InitialBalanceAmount *decimal.Decimal `json:"initial_balance_amount"` // default 0
	InitialBalance       *goprices.Money  `json:"initial_balance,omitempty" db:"-"`
	CurrentBalanceAmount *decimal.Decimal `json:"current_balance_amount"` // default 0
	CurrentBalance       *goprices.Money  `json:"current_balance,omitempty" db:"-"`
	model.ModelMetadata
}

// GiftCardFilterOption is used to buil sql queries
type GiftCardFilterOption struct {
	ExpiryDate    *model.TimeFilter
	StartDate     *model.TimeFilter
	Code          *model.StringFilter
	Currency      *model.StringFilter
	CreatedByID   *model.StringFilter
	CheckoutToken *model.StringFilter // SELECT * FROM 'Giftcards' WHERE 'Id' IN (SELECT 'GiftcardID' FROM 'GiftCardCheckouts' WHERE 'CheckoutID' ...)
	IsActive      *bool
	Distinct      bool // if true, SELECT DISTINCT

	SelectForUpdate bool // if true, concat `FOR UPDATE` to the end of SQL queries
}

func (gc *GiftCard) DisplayCode() string {
	return "****" + gc.Code[len(gc.Code)-4:]
}

// PopulateNonDbFields populates money fields for giftcard
func (gc *GiftCard) PopulateNonDbFields() {
	gc.InitialBalance, _ = goprices.NewMoney(gc.InitialBalanceAmount, gc.Currency)
	gc.CurrentBalance, _ = goprices.NewMoney(gc.CurrentBalanceAmount, gc.Currency)
}

func (gc *GiftCard) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.gift_card.is_valid.%s.app_error",
		"gift_card_id=",
		"GiftCard.IsValid",
	)

	if !model.IsValidId(gc.Id) {
		return outer("id", nil)
	}
	if gc.CreatedByID != nil && !model.IsValidId(*gc.CreatedByID) {
		return outer("created_by_id", &gc.Id)
	}
	if gc.UsedByID != nil && !model.IsValidId(*gc.UsedByID) {
		return outer("used_by_id", &gc.Id)
	}
	if gc.CreatedByEmail != nil && len(*gc.CreatedByEmail) > model.USER_EMAIL_MAX_LENGTH {
		return outer("created_by_email", &gc.Id)
	}
	if gc.UsedByEmail != nil && len(*gc.UsedByEmail) > model.USER_EMAIL_MAX_LENGTH {
		return outer("used_by_email", &gc.Id)
	}
	if len(gc.ExpiryType) > GiftcardExpiryTypeMaxLength || GiftcardExpiryTypeMap[gc.ExpiryType] == "" {
		return outer("expiry_type", &gc.Id)
	}
	if gc.ExpiryPeriodType != nil && (len(*gc.ExpiryPeriodType) > GiftcardExpiryPeriodTypeMaxLength || model.TimePeriodMap[*gc.ExpiryPeriodType] == "") {
		return outer("expiry_period_type", &gc.Id)
	}
	if gc.ExpiryPeriod != nil && *gc.ExpiryPeriod < 0 {
		return outer("expiry_period", &gc.Id)
	}
	if gc.Tag != nil && len(*gc.Tag) > GiftcardTagMaxLength {
		return outer("tag", &gc.Id)
	}
	if gc.ProductID != nil && !model.IsValidId(*gc.ProductID) {
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
		gc.Id = model.NewId()
	}
	gc.CreateAt = model.GetMillis()

	gc.commonPre()
}

func (gc *GiftCard) commonPre() {
	if gc.CurrentBalance != nil {
		gc.CurrentBalanceAmount = gc.CurrentBalance.Amount
	} else {
		gc.CurrentBalanceAmount = &decimal.Zero
	}

	if gc.InitialBalance != nil {
		gc.InitialBalanceAmount = gc.InitialBalance.Amount
	} else {
		gc.InitialBalanceAmount = &decimal.Zero
	}

	if gc.IsActive == nil {
		gc.IsActive = model.NewBool(true)
	}

	if gc.Currency == "" {
		gc.Currency = model.DEFAULT_CURRENCY
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
