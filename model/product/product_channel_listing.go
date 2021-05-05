package product

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
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
}

func (p *ProductChannelListing) IsAvailableForPurchase() bool {
	return p.AvailableForPurchase != nil && (*p.AvailableForPurchase).Before(time.Now())
}

func (p *ProductChannelListing) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.product_channel_listing.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "product_channel_listing_id=" + p.Id
	}

	return model.NewAppError("ProductChannelListing.IsValid", id, nil, details, http.StatusBadRequest)
}

func (p *ProductChannelListing) IsValid() *model.AppError {
	if !model.IsValidId(p.Id) {
		return p.createAppError("id")
	}
	if !model.IsValidId(p.ProductID) {
		return p.createAppError("product_id")
	}
	if !model.IsValidId(p.ChannelID) {
		return p.createAppError("channel_id")
	}
	if un, err := currency.ParseISO(p.Currency); !strings.EqualFold(un.String(), p.Currency) || err != nil {
		return p.createAppError("currency")
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
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

func ProductChannelListingFromJson(data io.Reader) *ProductChannelListing {
	var p ProductChannelListing
	err := json.JSON.NewDecoder(data).Decode(&p)
	if err != nil {
		return nil
	}
	return &p
}

const (
	PRODUCT_VARIANT_NAME_MAX_LENGTH = 255
	PRODUCT_VARIANT_SKU_MAX_LENGTH  = 255
)

// type VariantMedia struct {
// 	Id        string `json:"id"`
// 	VariantID string `json:"variant_id"`
// 	MediaID   string `json:"media_id"`
// }
