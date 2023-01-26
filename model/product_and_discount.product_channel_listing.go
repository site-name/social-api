package model

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
)

type ProductChannelListing struct {
	Id                    string           `json:"id"`
	ProductID             string           `json:"product_id"`
	ChannelID             string           `json:"channel_id"`
	VisibleInListings     bool             `json:"visible_in_listings"`
	AvailableForPurchase  *time.Time       `json:"available_for_purchase"` // UTC time
	Currency              string           `json:"currency"`
	DiscountedPriceAmount *decimal.Decimal `json:"discounted_price_amount"`           // can be NULL
	DiscountedPrice       *goprices.Money  `json:"discounted_price,omitempty" db:"-"` // can be NULL
	CreateAt              uint64           `json:"create_at"`
	Publishable

	channel *Channel `db:"-"` // this field may be populated when store performs prefetching
}

func (p *ProductChannelListing) GetChannel() *Channel {
	return p.channel
}

func (p *ProductChannelListing) SetChannel(c *Channel) {
	p.channel = c
}

// ProductChannelListingFilterOption is option for filtering product channel listing
type ProductChannelListingFilterOption struct {
	Id                   squirrel.Sqlizer
	ProductID            squirrel.Sqlizer
	ChannelID            squirrel.Sqlizer
	AvailableForPurchase squirrel.Sqlizer
	Currency             squirrel.Sqlizer
	ProductVariantsId    squirrel.Sqlizer // INNER JOIN Products ON ... INNER JOIN ProductVariants ON ... WHERE ProductVariants.Id ...
	PublicationDate      squirrel.Sqlizer //
	ChannelSlug          *string          // INNER JOIN Channels ON ... WHERE Channels.Slug ...
	VisibleInListings    *bool
	IsPublished          *bool
	PrefetchChannel      bool // this tell store to prefetch channel instances also
}

func (p *ProductChannelListing) IsAvailableForPurchase() bool {
	return p.AvailableForPurchase != nil && (p.AvailableForPurchase).Before(util.StartOfDay(time.Now().UTC()))
}

func (p *ProductChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"product_channel_listing.is_valid.%s.app_error",
		"product_channel_listing_id=",
		"ProductChannelListing.IsValid",
	)

	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if !IsValidId(p.ChannelID) {
		return outer("channel_id", &p.Id)
	}
	if un, err := currency.ParseISO(p.Currency); !strings.EqualFold(un.String(), p.Currency) || err != nil {
		return outer("currency", &p.Id)
	}
	return nil
}

func (p *ProductChannelListing) PopulateNonDbFields() {
	if p.DiscountedPriceAmount != nil && p.Currency != "" {
		p.DiscountedPrice = &goprices.Money{
			Amount:   *p.DiscountedPriceAmount,
			Currency: p.Currency,
		}
	}
}

func (p *ProductChannelListing) commonPre() {
	if p.DiscountedPrice != nil {
		p.DiscountedPriceAmount = &p.DiscountedPrice.Amount
	}

	if p.Currency != "" {
		p.Currency = strings.ToUpper(p.Currency)
	}
}

func (p *ProductChannelListing) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	if p.CreateAt == 0 {
		p.CreateAt = uint64(GetMillis())
	}
	p.commonPre()
}

func (p *ProductChannelListing) PreUpdate() {
	p.commonPre()
}

func (p *ProductChannelListing) DeepCopy() *ProductChannelListing {
	if p == nil {
		return nil
	}

	res := *p
	if p.channel != nil {
		res.channel = p.channel.DeepCopy()
	}
	if p.AvailableForPurchase != nil {
		res.AvailableForPurchase = NewPrimitive(*p.AvailableForPurchase)
	}
	return &res
}

type ProductChannelListings []*ProductChannelListing

func (p ProductChannelListings) IDs() []string {
	return lo.Map(p, func(r *ProductChannelListing, _ int) string { return r.Id })
}

func (p ProductChannelListings) ChannelIDs() []string {
	return lo.Map(p, func(r *ProductChannelListing, _ int) string { return r.ChannelID })
}

func (p ProductChannelListings) ProductIDs() []string {
	return lo.Map(p, func(r *ProductChannelListing, _ int) string { return r.ProductID })
}

func (p ProductChannelListings) DeepCopy() ProductChannelListings {
	return lo.Map(p, func(r *ProductChannelListing, _ int) *ProductChannelListing { return r.DeepCopy() })
}
