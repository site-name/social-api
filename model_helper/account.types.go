package model_helper

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

type AddressFilterOptions struct {
	UserID     qm.QueryMod // must bemodel.UserAddressWhere.UserID...
	Conditions []qm.QueryMod
}
