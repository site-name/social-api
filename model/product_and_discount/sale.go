package product_and_discount

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"golang.org/x/text/language"
)

// max lengths for some fields
const (
	SALE_NAME_MAX_LENGTH = 255
	SALE_TYPE_MAX_LENGTH = 10
)

type Sale struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Products    []*Product    `json:"products,omitempty" db:"-"`
	Categories  []*Category   `json:"categories,omitempty" db:"-"`
	Collections []*Collection `json:"collections,omitempty" db:"-"`
	StartDate   int64         `json:"start_date"`
	EndDate     *int64        `json:"end_date"`
}

func (s *Sale) String() string {
	return s.Name
}

func (s *Sale) GetDiscount(scl *SaleChannelListing) (*model.Money, error) {
	if scl == nil {
		return nil, &NotApplicable{
			Msg: "This sale if not assigned to this channel.",
		}
	}

	// if s.Type == FIXED {
	// 	discountAmount := &model.Money{
	// 		Amount:   scl.DiscountValue,
	// 		Currency: scl.Currency,
	// 	}
	// } else if s.Type == PERCENTAGE {

	// }

	// TODO: Fix me

	return nil, errors.New("unknown discount type")
}

func (s *Sale) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale.is_valid.%s.app_error",
		"sale_id=",
		"Sale.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if len(s.Type) > SALE_TYPE_MAX_LENGTH || !SALE_TYPES.Contains(s.Type) {
		return outer("type", &s.Id)
	}
	if s.StartDate == 0 {
		return outer("start_date", &s.Id)
	}

	return nil
}

func (s *Sale) ToJson() string {
	return model.ModelToJson(s)
}

func SaleFromJson(data io.Reader) *Sale {
	var s Sale
	model.ModelFromJson(&s, data)
	return &s
}

func (s *Sale) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.Type == "" || !SALE_TYPES.Contains(s.Type) {
		s.Type = FIXED
	}
	if s.StartDate == 0 {
		s.StartDate = model.GetMillis()
	}
}

type SaleTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	SaleID       string `json:"sale_id"`
}

func (s *SaleTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_translation.is_valid.%s.app_error",
		"sale_translation_id=",
		"SaleTranslation.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if tag, err := language.Parse(s.LanguageCode); err != nil || !strings.EqualFold(tag.String(), s.LanguageCode) {
		return outer("language_code", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}

	return nil
}

func (s *SaleTranslation) ToJson() string {
	return model.ModelToJson(s)
}

func SaleTranslationFromJson(data io.Reader) *SaleTranslation {
	var st SaleTranslation
	model.ModelFromJson(&st, data)
	return &st
}
