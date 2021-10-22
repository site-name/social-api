package shipping

import (
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type ShippingMethodChannelListing struct {
	Id                      string           `json:"id"`
	ShippingMethodID        string           `json:"shipping_method_id"`
	ChannelID               string           `json:"channel_id"`
	MinimumOrderPriceAmount *decimal.Decimal `json:"minimum_order_price_amount"`
	MinimumOrderPrice       *goprices.Money  `json:"minimum_order_price" db:"-"`
	Currency                string           `json:"currency"`
	MaximumOrderPriceAmount *decimal.Decimal `json:"maximum_order_price_amount"`
	MaximumOrderPrice       *goprices.Money  `json:"maximum_order_price" db:"-"`
	Price                   *goprices.Money  `json:"price" db:"-"`
	PriceAmount             *decimal.Decimal `json:"price_amount"`
	CreateAt                int64            `json:"create_at"`
}

// ShippingMethodChannelListingFilterOption is used to build sql queries
type ShippingMethodChannelListingFilterOption struct {
	ShippingMethodID *model.StringFilter
	ChannelID        *model.StringFilter
}

func (s *ShippingMethodChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shipping_method_channel_listing.is_valid.%s.app_error",
		"shipping_method_channel_listing_id=",
		"ShippingMethodChannelListing.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !model.IsValidId(s.ShippingMethodID) {
		return outer("shipping_method_id", &s.Id)
	}
	if !model.IsValidId(s.ChannelID) {
		return outer("channel_id", &s.Id)
	}
	if unit, err := currency.ParseISO(s.Currency); err != nil || !strings.EqualFold(unit.String(), s.Currency) {
		return outer("currency", &s.Id)
	}

	return nil
}

// PopulateNonDbFields populates non db fields of shipping method channel listing
func (s *ShippingMethodChannelListing) PopulateNonDbFields() {
	if s.MinimumOrderPriceAmount != nil {
		s.MinimumOrderPrice, _ = goprices.NewMoney(s.MinimumOrderPriceAmount, s.Currency)
	}
	if s.MaximumOrderPriceAmount != nil {
		s.MaximumOrderPrice, _ = goprices.NewMoney(s.MaximumOrderPriceAmount, s.Currency)
	}
	if s.PriceAmount != nil {
		s.Price, _ = goprices.NewMoney(s.PriceAmount, s.Currency)
	}
}

func (s *ShippingMethodChannelListing) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
	s.commonPre()
}

func (s *ShippingMethodChannelListing) commonPre() {
	if s.MinimumOrderPriceAmount == nil {
		s.MinimumOrderPriceAmount = &decimal.Zero
	}
	if s.PriceAmount == nil {
		s.PriceAmount = &decimal.Zero
	}
}

func (s *ShippingMethodChannelListing) PreUpdate() {
	s.commonPre()
}

// GetTotal retuns current ShippingMethodChannelListing's Price fields
func (s *ShippingMethodChannelListing) GetTotal() *goprices.Money {
	s.PopulateNonDbFields()
	return s.Price
}

func (s *ShippingMethodChannelListing) ToJSON() string {
	s.PopulateNonDbFields()
	return model.ModelToJson(s)
}
