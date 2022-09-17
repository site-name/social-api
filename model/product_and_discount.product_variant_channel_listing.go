package model

import (
	"strings"

	"github.com/Masterminds/squirrel"
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

	preorderQuantityAllocated int      `json:"-" db:"-"` // this field got populated in some db queries
	availablePreorderQuantity int      `json:"-" db:"-"`
	Channel                   *Channel `json:"-" db:"-"` // this field got populated in some db queries
}

// ProductVariantChannelListingFilterOption is used to build sql queries
type ProductVariantChannelListingFilterOption struct {
	Id          squirrel.Sqlizer
	VariantID   squirrel.Sqlizer
	ChannelID   squirrel.Sqlizer
	PriceAmount squirrel.Sqlizer

	VariantProductID squirrel.Sqlizer // INNER JOIN ProductVariants WHERE ProductVariants.ProductID ...

	SelectRelatedChannel bool   // tell store to select related Channel(s)
	SelectForUpdate      bool   // if true, add `FOR UPDATE` to the end of query
	SelectForUpdateOf    string // if provided, tell database system to lock on specific row(s)

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

type ProductVariantChannelListings []*ProductVariantChannelListing

func (p ProductVariantChannelListings) IDs() []string {
	var res []string
	for _, item := range p {
		if item != nil {
			res = append(res, item.Id)
		}
	}

	return res
}

func (p *ProductVariantChannelListing) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"product_variant_channel_listing.is_valid.%s.app_error",
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
	if p.Channel != nil {
		res.Channel = p.Channel.DeepCopy()
	}
	return &res
}
