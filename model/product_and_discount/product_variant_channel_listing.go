package product_and_discount

import (
	"io"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type ProductVariantChannelListing struct {
	Id              string           `json:"id"`
	VariantID       string           `json:"variant_id"` // not null
	ChannelID       string           `json:"channel_id"` // not null
	Currency        string           `json:"currency"`
	PriceAmount     *decimal.Decimal `json:"price_amount,omitempty"`
	Price           *goprices.Money  `json:"price,omitempty" db:"-"`
	CostPriceAmount *decimal.Decimal `json:"cost_price_amount"`
	CostPrice       *goprices.Money  `json:"cost_price,omitempty" db:"-"`
}

func (p *ProductVariantChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_variant_channel_listing.is_valid.%s.app_error",
		"product_variant_channel_listing_id=",
		"ProductVariantChannelListing.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.VariantID) {
		return outer("variant_id", &p.Id)
	}
	if !model.IsValidId(p.ChannelID) {
		return outer("channel_id", &p.Id)
	}
	if unit, err := currency.ParseISO(p.Currency); err != nil || !strings.EqualFold(unit.String(), p.Currency) {
		return outer("currency", &p.Id)
	}

	return nil
}

func (p *ProductVariantChannelListing) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
}

func (p *ProductVariantChannelListing) PopulateNonDbFields() {
	if p.Price == nil && p.PriceAmount != nil {
		p.Price = &goprices.Money{
			Amount:   p.PriceAmount,
			Currency: p.Currency,
		}
	}
	if p.CostPrice == nil && p.CostPriceAmount != nil {
		p.CostPrice = &goprices.Money{
			Amount:   p.CostPriceAmount,
			Currency: p.Currency,
		}
	}
}

func (p *ProductVariantChannelListing) ToJson() string {
	p.PopulateNonDbFields()
	return model.ModelToJson(p)
}

func ProductVariantChannelListingFromJson(data io.Reader) *ProductVariantChannelListing {
	var p ProductVariantChannelListing
	model.ModelFromJson(&p, data)
	return &p
}
