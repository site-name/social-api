package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"golang.org/x/text/currency"
)

type ProductVariantChannelListing struct {
	Id                        string           `json:"id"`
	VariantID                 string           `json:"variant_id"` // not null
	ChannelID                 string           `json:"channel_id"` // not null
	Currency                  string           `json:"currency"`
	PriceAmount               *decimal.Decimal `json:"price_amount,omitempty"` // can be NULL
	Price                     *goprices.Money  `json:"price,omitempty" db:"-"`
	CostPriceAmount           *decimal.Decimal `json:"cost_price_amount"` // can be NULL
	CostPrice                 *goprices.Money  `json:"cost_price,omitempty" db:"-"`
	PreorderQuantityThreshold *int             `json:"preorder_quantity_threshold"`
	CreateAt                  int64            `json:"create_at"`

	preorderQuantityAllocated int             `db:"-"` // this field got populated in some db queries
	availablePreorderQuantity int             `db:"-"`
	channel                   *Channel        `db:"-"` // this field got populated in some db queries
	variant                   *ProductVariant `db:"-"` // this field got populated in some store functions
}

// ProductVariantChannelListingFilterOption is used to build sql queries
type ProductVariantChannelListingFilterOption struct {
	Id          squirrel.Sqlizer
	VariantID   squirrel.Sqlizer
	ChannelID   squirrel.Sqlizer
	PriceAmount squirrel.Sqlizer

	VariantProductID squirrel.Sqlizer // INNER JOIN ProductVariants ON ... WHERE ProductVariants.ProductID ...

	SelectRelatedChannel        bool   // tell store to select related Channel(s)
	SelectRelatedProductVariant bool   // tell store to select related product variant
	SelectForUpdate             bool   // if true, add `FOR UPDATE` to the end of query
	SelectForUpdateOf           string // if provided, tell database system to lock on specific row(s)

	AnnotatePreorderQuantityAllocated bool // set true to populate `preorderQuantityAllocated` field of returning product variant channel listings
	AnnotateAvailablePreorderQuantity bool // set true to populate `availablePreorderQuantity` field of returning product variant channel listings
}

func (p *ProductVariantChannelListing) Set_preorderQuantityAllocated(value int) {
	p.preorderQuantityAllocated = value
}
func (p *ProductVariantChannelListing) Get_preorderQuantityAllocated() int {
	return p.preorderQuantityAllocated
}

func (p *ProductVariantChannelListing) Set_availablePreorderQuantity(value int) {
	p.availablePreorderQuantity = value
}
func (p *ProductVariantChannelListing) Get_availablePreorderQuantity() int {
	return p.availablePreorderQuantity
}

func (p *ProductVariantChannelListing) GetChannel() *Channel {
	return p.channel
}
func (p *ProductVariantChannelListing) SetChannel(c *Channel) {
	p.channel = c
}

func (p *ProductVariantChannelListing) GetVariant() *ProductVariant {
	return p.variant
}
func (p *ProductVariantChannelListing) SetVariant(c *ProductVariant) {
	p.variant = c
}

type ProductVariantChannelListings []*ProductVariantChannelListing

func (p ProductVariantChannelListings) IDs() []string {
	return lo.Map(p, func(l *ProductVariantChannelListing, _ int) string { return l.Id })
}

func (ps ProductVariantChannelListings) DeepCopy() ProductVariantChannelListings {
	return lo.Map(ps, func(p *ProductVariantChannelListing, _ int) *ProductVariantChannelListing { return p.DeepCopy() })
}

func (p *ProductVariantChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.product_variant_channel_listing.is_valid.%s.app_error",
		"product_variant_channel_listing_id=",
		"ProductVariantChannelListing.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.VariantID) {
		return outer("variant_id", &p.Id)
	}
	if !IsValidId(p.ChannelID) {
		return outer("channel_id", &p.Id)
	}
	if unit, err := currency.ParseISO(p.Currency); err != nil || !strings.EqualFold(unit.String(), p.Currency) {
		return outer("currency", &p.Id)
	}

	return nil
}

func (p *ProductVariantChannelListing) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.CreateAt = GetMillis()
	p.commonPre()
}

func (p *ProductVariantChannelListing) commonPre() {
	if p.Price != nil {
		p.PriceAmount = &p.Price.Amount
	}

	if p.CostPrice != nil {
		p.CostPriceAmount = &p.CostPrice.Amount
	}

	if p.Currency != "" {
		p.Currency = strings.ToUpper(p.Currency)
	} else {
		p.Currency = DEFAULT_CURRENCY
	}
}

func (p *ProductVariantChannelListing) PopulateNonDbFields() {
	if p.PriceAmount != nil {
		p.Price = &goprices.Money{
			Amount:   *p.PriceAmount,
			Currency: p.Currency,
		}
	}
	if p.CostPriceAmount != nil {
		p.CostPrice = &goprices.Money{
			Amount:   *p.CostPriceAmount,
			Currency: p.Currency,
		}
	}
}

func (p *ProductVariantChannelListing) DeepCopy() *ProductVariantChannelListing {
	res := *p

	if p.PriceAmount != nil {
		res.PriceAmount = NewPrimitive(*p.PriceAmount)
	}
	if p.CostPriceAmount != nil {
		res.CostPriceAmount = NewPrimitive(*p.CostPriceAmount)
	}
	if p.PreorderQuantityThreshold != nil {
		res.PreorderQuantityThreshold = NewPrimitive(*p.PreorderQuantityThreshold)
	}
	if p.channel != nil {
		res.channel = p.channel.DeepCopy()
	}
	if p.variant != nil {
		res.variant = p.variant.DeepCopy()
	}
	return &res
}
