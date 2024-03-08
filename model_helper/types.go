package model_helper

import (
	"github.com/sitename/sitename/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AddressFilterOptions struct {
	CommonQueryOptions
}

type ChannelFilterOptions struct {
	CommonQueryOptions
	ShippingZoneID    qm.QueryMod // INNER JOIN shipping_zone_channels szc ON ... WHERE szc.shipping_zone_id = ?
	VoucherID         qm.QueryMod // INNER JOIN voucher_channel_listings vcl ON ... WHERE vcl.voucher_id ...
	AnnotateHasOrders bool        // this tells the store to annotate if the channels has order(s) attached
}

type VatFilterOptions struct {
	CommonQueryOptions
}

type ExternalAccessTokens struct {
	Token        *string
	RefreshToken *string
	CsrfToken    *string
	User         *model.User
}
