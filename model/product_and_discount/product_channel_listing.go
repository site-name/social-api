package product_and_discount

import (
	"strings"
	"time"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/text/currency"
)

type ProductChannelListing struct {
	Id                    string           `json:"id"`
	ProductID             string           `json:"product_id"`
	ChannelID             string           `json:"channel_id"`
	VisibleInListings     bool             `json:"visible_in_listings"`
	AvailableForPurchase  *time.Time       `json:"available_for_purchase"`  // UTC time
	Currency              string           `json:"currency"`                // default "USD"
	DiscountedPriceAmount *decimal.Decimal `json:"discounted_price_amount"` // default decimal(0)
	DiscountedPrice       *goprices.Money  `json:"discounted_price,omitempty" db:"-"`
	CreateAt              uint64           `json:"create_at"`
	model.Publishable
}

// ProductChannelListingFilterOption is option for filtering product channel listing
type ProductChannelListingFilterOption struct {
	ProductID            *model.StringFilter
	ChannelID            *model.StringFilter
	ChannelSlug          *string // inner join Channel
	VisibleInListings    *bool
	AvailableForPurchase *model.TimeFilter
	Currency             *model.StringFilter
	ProductVariantsId    *model.StringFilter // inner join product, product variant
	PublicationDate      *model.TimeFilter
	IsPublished          *bool
}

// Check if product
func (p *ProductChannelListing) IsAvailableForPurchase() bool {
	return p.AvailableForPurchase != nil && (p.AvailableForPurchase).Before(util.StartOfDay(time.Now().UTC()))
}

func (p *ProductChannelListing) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_channel_listing.is_valid.%s.app_error",
		"product_channel_listing_id=",
		"ProductChannelListing.IsValid",
	)

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

func (p *ProductChannelListing) PopulateNonDbFields() {
	p.DiscountedPrice, _ = goprices.NewMoney(p.DiscountedPriceAmount, p.Currency)
}

func (p *ProductChannelListing) commonPre() {
	if p.DiscountedPrice != nil {
		p.DiscountedPriceAmount = p.DiscountedPrice.Amount
	} else {
		p.DiscountedPriceAmount = &decimal.Zero
	}

	if p.Currency == "" {
		p.Currency = model.DEFAULT_CURRENCY
	} else {
		p.Currency = strings.ToUpper(p.Currency)
	}
}

func (p *ProductChannelListing) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	if p.CreateAt == 0 {
		p.CreateAt = uint64(model.GetMillis())
	}
	p.commonPre()
}
