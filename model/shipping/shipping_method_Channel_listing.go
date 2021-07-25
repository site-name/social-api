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

func (s *ShippingMethodChannelListing) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.MinimumOrderPriceAmount == nil {
		s.MinimumOrderPriceAmount = &decimal.Zero
	}
	if s.PriceAmount == nil {
		s.PriceAmount = &decimal.Zero
	}
}

func (s *ShippingMethodChannelListing) GetTotal() *goprices.Money {
	return s.Price
}

func (s *ShippingMethodChannelListing) ToJson() string {
	if s.MinimumOrderPrice == nil {
		s.MinimumOrderPrice = &goprices.Money{
			Amount:   s.MinimumOrderPriceAmount,
			Currency: s.Currency,
		}
	}
	if s.MaximumOrderPrice == nil {
		s.MaximumOrderPrice = &goprices.Money{
			Amount:   s.MaximumOrderPriceAmount,
			Currency: s.Currency,
		}
	}
	if s.Price == nil {
		s.Price = &goprices.Money{
			Amount:   s.PriceAmount,
			Currency: s.Currency,
		}
	}

	return model.ModelToJson(s)
}
