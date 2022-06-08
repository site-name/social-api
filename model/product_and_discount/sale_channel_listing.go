package product_and_discount

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type SaleChannelListing struct {
	Id            string           `json:"id"`
	SaleID        string           `json:"sale_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"` // default decimal(0)
	Currency      string           `json:"currency"`
	CreateAt      int64            `json:"create_at"`
}

type SaleChannelListingFilterOption struct {
	Id        squirrel.Sqlizer
	SaleID    squirrel.Sqlizer
	ChannelID squirrel.Sqlizer
}

func (s *SaleChannelListing) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
	if s.DiscountValue == nil || s.DiscountValue.LessThan(decimal.Zero) {
		s.DiscountValue = &decimal.Zero
	}
	s.Currency = strings.ToUpper(s.Currency)
}

func (s *SaleChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_channel_listing.is_valid.%s.app_error",
		"sale_channel_listing_id=",
		"SaleChannelListing.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !model.IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !model.IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	if unit, err := currency.ParseISO(s.Currency); err == nil || !strings.EqualFold(unit.String(), s.Currency) {
		return outer("currency", &s.Id)
	}

	return nil
}
