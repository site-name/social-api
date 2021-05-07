package product_and_discount

import (
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"golang.org/x/text/currency"
)

type ProductChannelListing struct {
	Id                    string           `json:"id"`
	ProductID             string           `json:"product_id"`
	ChannelID             string           `json:"channel_id"`
	VisibleInListings     bool             `json:"visible_in_listings"`
	AvailableForPurchase  *time.Time       `json:"available_for_purchase"`
	Currency              string           `json:"currency"`
	DiscountedPriceAmount *decimal.Decimal `json:"discounted_price_amount"`
	DiscountedPrice       *model.Money     `json:"discounted_price,omitempty" db:"-"`
	*model.Publishable    `db:"-"`
}

// Check if product is can be bought now
func (p *ProductChannelListing) IsAvailableForPurchase() bool {
	return p.AvailableForPurchase != nil && (*p.AvailableForPurchase).Before(time.Now())
}

func (p *ProductChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_channel_listing.is_valid.%s.app_error",
		"product_channel_listing_id=",
		"ProductChannelListing.IsValid")

	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if !model.IsValidId(p.ChannelID) {
		return outer("channel_id", &p.Id)
	}
	if un, err := currency.ParseISO(p.Currency); !strings.EqualFold(un.String(), p.Currency) || err != nil {
		return outer("currency", &p.Id)
	}
	return nil
}

func (p *ProductChannelListing) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
}

func (p *ProductChannelListing) ToJson() string {
	p.DiscountedPrice = &model.Money{
		Amount:   p.DiscountedPriceAmount,
		Currency: p.Currency,
	}
	return model.ModelToJson(p)
}

func ProductChannelListingFromJson(data io.Reader) *ProductChannelListing {
	var p ProductChannelListing
	model.ModelFromJson(&p, data)
	return &p
}
