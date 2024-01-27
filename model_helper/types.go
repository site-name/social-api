package model_helper

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

type AddressFilterOptions struct {
	UserID     qm.QueryMod // must be model.UserAddressWhere.UserID...
	Conditions []qm.QueryMod
}

type ChannelFilterOptions struct {
	Conds          []qm.QueryMod
	ShippingZoneID qm.QueryMod // INNER JOIN shipping_zone_channels szc ON ... WHERE szc.shipping_zone_id = ?
	VoucherID      qm.QueryMod // INNER JOIN voucher_channel_listings vcl ON ... WHERE vcl.voucher_id ...
}
