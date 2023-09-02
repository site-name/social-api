package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type ProductVariantChannelListing struct {
	Id                        UUID             `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	VariantID                 UUID             `json:"variant_id" gorm:"type:uuid;column:VariantID"` // not null
	ChannelID                 UUID             `json:"channel_id" gorm:"type:uuid;column:ChannelID"` // not null
	Currency                  string           `json:"currency" gorm:"type:varchar(5);column:Currency"`
	PriceAmount               *decimal.Decimal `json:"price_amount,omitempty" gorm:"column:PriceAmount;type:decimal(12,3)"` // can be NULL
	CostPriceAmount           *decimal.Decimal `json:"cost_price_amount" gorm:"column:CostPriceAmount;type:decimal(12,3)"`  // can be NULL
	PreorderQuantityThreshold *int             `json:"preorder_quantity_threshold" gorm:"column:PreorderQuantityThreshold"`
	CreateAt                  int64            `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`

	Price     *goprices.Money `json:"price,omitempty" gorm:"-"`
	CostPrice *goprices.Money `json:"cost_price,omitempty" gorm:"-"`

	preorderQuantityAllocated int             `gorm:"-"` // this field got populated in some db queries
	availablePreorderQuantity int             `gorm:"-"`
	channel                   *Channel        `gorm:"-"` // this field got populated in some db queries
	variant                   *ProductVariant `gorm:"-"` // this field got populated in some store functions
}

func (c *ProductVariantChannelListing) BeforeCreate(_ *gorm.DB) error {
	c.commonPre()
	return c.IsValid()
}
func (c *ProductVariantChannelListing) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent updating
	return c.IsValid()
}
func (c *ProductVariantChannelListing) TableName() string {
	return ProductVariantChannelListingTableName
}

// ProductVariantChannelListingFilterOption is used to build sql queries
type ProductVariantChannelListingFilterOption struct {
	Conditions squirrel.Sqlizer

	VariantProductID squirrel.Sqlizer // INNER JOIN ProductVariants ON ... WHERE ProductVariants.ProductID ...

	SelectRelatedChannel        bool   // tell store to select related Channel(s)
	SelectRelatedProductVariant bool   // tell store to select related product variant
	SelectForUpdate             bool   // if true, add `FOR UPDATE` to the end of query. NOTE: only apply when Transaction is set
	SelectForUpdateOf           string // if provided, tell database system to lock on specific row(s)

	AnnotatePreorderQuantityAllocated bool // set true to populate `preorderQuantityAllocated` field of returning product variant channel listings
	AnnotateAvailablePreorderQuantity bool // set true to populate `availablePreorderQuantity` field of returning product variant channel listings

	Transaction *gorm.DB
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

func (p ProductVariantChannelListings) IDs() []UUID {
	return lo.Map(p, func(l *ProductVariantChannelListing, _ int) UUID { return l.Id })
}

func (p ProductVariantChannelListings) VariantIDs() util.AnyArray[UUID] {
	return lo.Map(p, func(l *ProductVariantChannelListing, _ int) UUID { return l.VariantID })
}

func (ps ProductVariantChannelListings) DeepCopy() ProductVariantChannelListings {
	return lo.Map(ps, func(p *ProductVariantChannelListing, _ int) *ProductVariantChannelListing { return p.DeepCopy() })
}

func (p *ProductVariantChannelListing) IsValid() *AppError {
	if !IsValidId(p.VariantID) {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.variant_id.app_error", nil, "please provide valid variant id", http.StatusBadRequest)
	}
	if !IsValidId(p.ChannelID) {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if _, err := currency.ParseISO(p.Currency); err != nil {
		return NewAppError("ProductVariantChannelListing.IsValid", "model.product_variant_channel_listing.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}

	return nil
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

	res.PriceAmount = CopyPointer(p.PriceAmount)
	res.CostPriceAmount = CopyPointer(p.CostPriceAmount)
	res.PreorderQuantityThreshold = CopyPointer(p.PreorderQuantityThreshold)

	if p.channel != nil {
		res.channel = p.channel.DeepCopy()
	}
	if p.variant != nil {
		res.variant = p.variant.DeepCopy()
	}
	return &res
}
