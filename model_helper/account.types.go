package model_helper

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

type AddressFilterOptions struct {
	UserID     qm.QueryMod // Id IN (SELECT AddressID FROM UserAddresses ON ... WHERE UserAddresses.UserID ...)
	Conditions []qm.QueryMod
}
