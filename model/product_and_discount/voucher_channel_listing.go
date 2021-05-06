package product_and_discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"golang.org/x/text/currency"
)

type VoucherChannelListing struct {
	Id            string           `json:"id"`
	VoucherID     string           `json:"voucher_id"`
	ChannelID     string           `json:"channel_id"`
	DiscountValue *decimal.Decimal `json:"discount_value"`
	Discount      *model.Money     `json:"discount,omitempty" db:"_"`
	MinSpent      *model.Money     `json:"min_spent,omitempty" db:"_"`
	Currency      string           `json:"currency"`
	MinSpenAmount *decimal.Decimal `json:"min_spent_amount"`
}

func (v *VoucherChannelListing) ToJson() string {
	v.Discount = &model.Money{
		Amount:   v.DiscountValue,
		Currency: v.Currency,
	}
	v.MinSpent = &model.Money{
		Amount:   v.MinSpenAmount,
		Currency: v.Currency,
	}
	b, _ := json.JSON.Marshal(v)
	return string(b)
}

func VoucherChannelListingFromJson(data io.Reader) *VoucherChannelListing {
	var vcl VoucherChannelListing
	err := json.JSON.NewDecoder(data).Decode(&vcl)
	if err != nil {
		return nil
	}
	return &vcl
}

func (v *VoucherChannelListing) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.voucher_channel_listing.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "voucher_channel_listing_id=" + v.Id
	}

	return model.NewAppError("VoucherChannelListing.IsValid", id, nil, details, http.StatusBadRequest)
}

func (v *VoucherChannelListing) IsValid() *model.AppError {
	if v.Id == "" {
		return v.createAppError("id")
	}
	if v.VoucherID == "" {
		return v.createAppError("voucher_id")
	}
	if v.ChannelID == "" {
		return v.createAppError("channel_id")
	}
	if unit, err := currency.ParseISO(v.Currency); err != nil || !strings.EqualFold(unit.String(), v.Currency) {
		return v.createAppError("currency")
	}

	return nil
}
