package model

import (
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
	"gorm.io/gorm"
)

type ProductChannelListing struct {
	Id                    string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	ProductID             string           `json:"product_id" gorm:"type:uuid;column:ProductID;index:productid_channelid_key"`
	ChannelID             string           `json:"channel_id" gorm:"type:uuid;column:ChannelID;index:productid_channelid_key"`
	VisibleInListings     bool             `json:"visible_in_listings" gorm:"column:VisibleInListings"`
	AvailableForPurchase  *time.Time       `json:"available_for_purchase" gorm:"column:AvailableForPurchase"` // precision to date. E.g 2021-09-08
	Currency              string           `json:"currency" gorm:"type:varchar(5);column:Currency"`
	DiscountedPriceAmount *decimal.Decimal `json:"discounted_price_amount" gorm:"column:DiscountedPriceAmount"` // can be NULL
	CreateAt              uint64           `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	Publishable

	DiscountedPrice *goprices.Money `json:"discounted_price,omitempty" gorm:"-"` // can be NULL

	channel *Channel `gorm:"-"` // this field may be populated when store performs prefetching
}

func (c *ProductChannelListing) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductChannelListing) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent updating
	return c.IsValid()
}
func (c *ProductChannelListing) TableName() string     { return ProductChannelListingTableName }
func (p *ProductChannelListing) GetChannel() *Channel  { return p.channel }
func (p *ProductChannelListing) SetChannel(c *Channel) { p.channel = c }

// ProductChannelListingFilterOption is option for filtering product channel listing
type ProductChannelListingFilterOption struct {
	Conditions squirrel.Sqlizer

	ProductVariantsId squirrel.Sqlizer // INNER JOIN Products ON ... INNER JOIN ProductVariants ON ... WHERE ProductVariants.Id ...
	ChannelSlug       squirrel.Sqlizer // INNER JOIN Channels ON ... WHERE Channels.Slug ...
	PrefetchChannel   bool             // this tell store to prefetch channel instances also
}

func (p *ProductChannelListing) IsAvailableForPurchase() bool {
	return p.AvailableForPurchase != nil && (p.AvailableForPurchase).Before(util.StartOfDay(time.Now().UTC()))
}

func (p *ProductChannelListing) IsValid() *AppError {
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if !IsValidId(p.ChannelID) {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.channel_id.app_error", nil, "please provide valid channel id", http.StatusBadRequest)
	}
	if un, err := currency.ParseISO(p.Currency); !strings.EqualFold(un.String(), p.Currency) || err != nil {
		return NewAppError("ProductChannelListing.IsValid", "model.product_channel_listing.is_valid.currency.app_error", nil, "please provide valid currency", http.StatusBadRequest)
	}
	if err := ValidateDecimal("ProductChannelListing.IsValid", p.DiscountedPriceAmount, DECIMAL_TOTAL_DIGITS_ALLOWED, DECIMAL_MAX_DECIMAL_PLACES_ALLOWED); err != nil {
		return err
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

func (p *ProductChannelListing) DeepCopy() *ProductChannelListing {
	if p == nil {
		return nil
	}

	res := *p
	if p.channel != nil {
		res.channel = p.channel.DeepCopy()
	}
	res.AvailableForPurchase = CopyPointer(p.AvailableForPurchase)
	res.Publishable = *p.Publishable.DeepCopy()
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
