package api

import (
	"github.com/sitename/sitename/model"
)

type Sale struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Type            SaleType        `json:"type"`
	StartDate       DateTime        `json:"startDate"`
	EndDate         *DateTime       `json:"endDate"`
	PrivateMetadata []*MetadataItem `json:"privateMetadata"`
	Metadata        []*MetadataItem `json:"metadata"`

	// DiscountValue   *float64        `json:"discountValue"`
	// Currency        *string         `json:"currency"`
	// Categories      *CategoryCountableConnection       `json:"categories"`
	// Collections     *CollectionCountableConnection     `json:"collections"`
	// Products        *ProductCountableConnection        `json:"products"`
	// Variants        *ProductVariantCountableConnection `json:"variants"`
	// Translation     *SaleTranslation                   `json:"translation"`
	// ChannelListings []*SaleChannelListing              `json:"channelListings"`
}

func systemSaleToGraphqlSale(s *model.Sale) *Sale {
	if s == nil {
		return nil
	}

	res := &Sale{
		ID:              s.Id,
		Name:            s.Name,
		Type:            SaleType(s.Type),
		StartDate:       DateTime{s.StartDate},
		Metadata:        MetadataToSlice(s.Metadata),
		PrivateMetadata: MetadataToSlice(s.PrivateMetadata),
	}
	if s.EndDate != nil {
		res.EndDate = &DateTime{*s.EndDate}
	}

	return res
}
