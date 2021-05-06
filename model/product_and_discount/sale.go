package product_and_discount

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
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

func (s *Sale) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.sale.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "sale_id=" + s.Id
	}

	return model.NewAppError("Sale.IsValid", id, nil, details, http.StatusBadRequest)
}

func (s *Sale) IsValid() *model.AppError {
	if s.Id == "" {
		return s.createAppError("id")
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return s.createAppError("name")
	}
	if len(s.Type) > SALE_TYPE_MAX_LENGTH || !SALE_TYPES.Contains(s.Type) {
		return s.createAppError("type")
	}
	if s.StartDate == 0 {
		return s.createAppError("start_date")
	}

	return nil
}

func (s *Sale) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SaleFromJson(data io.Reader) *Sale {
	var s Sale
	err := json.JSON.NewDecoder(data).Decode(&s)
	if err != nil {
		return nil
	}
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

func (s *SaleTranslation) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.sale_translation.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "sale_translation_id=" + s.Id
	}
	return model.NewAppError("SaleTranslation.IsValid", id, nil, details, http.StatusBadRequest)
}

func (s *SaleTranslation) IsValid() *model.AppError {
	if s.Id == "" {
		return s.createAppError("id")
	}
	if s.SaleID == "" {
		return s.createAppError("sale_id")
	}
	if tag, err := language.Parse(s.LanguageCode); err != nil || !strings.EqualFold(tag.String(), s.LanguageCode) {
		return s.createAppError("language_code")
	}
	if utf8.RuneCountInString(s.Name) > SALE_NAME_MAX_LENGTH {
		return s.createAppError("name")
	}

	return nil
}

func (s *SaleTranslation) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SaleTranslationFromJson(data io.Reader) *SaleTranslation {
	var st SaleTranslation
	err := json.JSON.NewDecoder(data).Decode(&st)
	if err != nil {
		return nil
	}
	return &st
}
