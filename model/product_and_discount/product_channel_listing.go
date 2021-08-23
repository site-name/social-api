package product_and_discount

import (
	"strings"
	"time"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
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
	model.Publishable

	Channel *channel.Channel `json:"-" db:"-"` // this field may be populated when store performs prefetching
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
	PrefetchChannel      bool // this tell store to prefetch channel instances also
}

type ProductChannelListings []*ProductChannelListing

type WhichID uint8

const (
	Ids WhichID = iota
	ChannelIDs
	ProductIDs
)

func (p ProductChannelListings) GetIDs(whichID WhichID) []string {
	var res []string
	for _, item := range p {
		if item != nil {
			switch whichID {
			case Ids:
				res = append(res, item.Id)
			case ProductIDs:
				res = append(res, item.ProductID)
			case ChannelIDs:
				res = append(res, item.ChannelID)
			}
		}
	}

	return res
}

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
	if p.DiscountedPriceAmount != nil && p.Currency != "" {
		p.DiscountedPrice, _ = goprices.NewMoney(p.DiscountedPriceAmount, p.Currency)
	}
}

func (p *ProductChannelListing) commonPre() {
	if p.DiscountedPrice != nil {
		p.DiscountedPriceAmount = p.DiscountedPrice.Amount
	}

	if p.Currency != "" {
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
