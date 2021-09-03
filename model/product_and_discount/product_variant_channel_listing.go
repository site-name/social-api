package product_and_discount

import (
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"golang.org/x/text/currency"
)

type ProductVariantChannelListing struct {
	Id              string           `json:"id"`
	VariantID       string           `json:"variant_id"` // not null
	ChannelID       string           `json:"channel_id"` // not null
	Currency        string           `json:"currency"`
	PriceAmount     *decimal.Decimal `json:"price_amount,omitempty"` // can be NULL
	Price           *goprices.Money  `json:"price,omitempty" db:"-"`
	CostPriceAmount *decimal.Decimal `json:"cost_price_amount"` // can be NULL
	CostPrice       *goprices.Money  `json:"cost_price,omitempty" db:"-"`
	CreateAt        int64            `json:"create_at"`

	Channel *channel.Channel `json:"-" db:"-"`
}

// ProductVariantChannelListingFilterOption is used to build sql queries
type ProductVariantChannelListingFilterOption struct {
	Id          *model.StringFilter
	VariantID   *model.StringFilter
	ChannelID   *model.StringFilter
	PriceAmount *model.NumberFilter

	VariantProductID *model.StringFilter // INNER JOIN ProductVariants
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
	p.CreateAt = model.GetMillis()
	p.commonPre()
}

func (p *ProductVariantChannelListing) commonPre() {
	if p.Price != nil {
		p.PriceAmount = p.Price.Amount
	}

	if p.CostPrice != nil {
		p.CostPriceAmount = p.CostPrice.Amount
	}

	if p.Currency != "" {
		p.Currency = strings.ToUpper(p.Currency)
	}
}

func (p *ProductVariantChannelListing) PopulateNonDbFields() {
	if p.PriceAmount != nil && p.Currency != "" {
		p.Price, _ = goprices.NewMoney(p.PriceAmount, p.Currency)
	}
	if p.CostPriceAmount != nil && p.Currency != "" {
		p.CostPrice, _ = goprices.NewMoney(p.CostPriceAmount, p.Currency)
	}
}

func (p *ProductVariantChannelListing) ToJson() string {
	p.PopulateNonDbFields()
	return model.ModelToJson(p)
}
