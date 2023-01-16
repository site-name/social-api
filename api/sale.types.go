package api

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
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

func (s *Sale) Translation(ctx context.Context, args struct{ LanguageCode LanguageCodeEnum }) (*SaleTranslation, error) {
	panic("not implemented")
}

func (s *Sale) Categories(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CategoryCountableConnection, error) {
	panic("not implemented")
}

func (s *Sale) Collections(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*CollectionCountableConnection, error) {
	panic("not implemented")
}

func (s *Sale) Products(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductCountableConnection, error) {
	panic("not implemented")
}

func (s *Sale) Variants(ctx context.Context, args struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}) (*ProductVariantCountableConnection, error) {
	panic("not implemented")
}

func (v *Sale) DiscountValue(ctx context.Context) (*float64, error) {
	// VoucherChannelListingByVoucherIdLoader.Load(ctx, v.ID)()
	panic("not implemented")
}

func (v *Sale) Currency(ctx context.Context) (*string, error) {
	panic("not implemented")
}

func (v *Sale) ChannelListings(ctx context.Context) ([]*SaleChannelListing, error) {
	embedCtx, err := GetContextValue[*web.Context](ctx, WebCtx)
	if err != nil {
		return nil, err
	}

	currentSession := embedCtx.AppContext.Session()
	if embedCtx.App.Srv().AccountService().SessionHasPermissionTo(currentSession, model.PermissionManageDiscounts) {
		listings, err := SaleChannelListingBySaleIdLoader.Load(ctx, v.ID)()
		if err != nil {
			return nil, err
		}

		return DataloaderResultMap(listings, systemSaleChannelListingToGraphqlSaleChannelListing), nil
	}

	return nil, model.NewAppError("Voucher.ChannelListings", ErrorUnauthorized, nil, "you are not authorized to perform this action", http.StatusUnauthorized)
}
