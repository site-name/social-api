package account

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) ScanFields(addr *model.Address) []any {
	return []any{
		&addr.ID,
		&addr.FirstName,
		&addr.LastName,
		&addr.CompanyName,
		&addr.StreetAddress1,
		&addr.StreetAddress2,
		&addr.City,
		&addr.CityArea,
		&addr.PostalCode,
		&addr.Country,
		&addr.CountryArea,
		&addr.Phone,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	}
}

func (as *SqlAddressStore) Upsert(transaction store.ContextRunner, address *model.Address) (*model.Address, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	model_helper.AddressCommonPre(address)
	if err := model_helper.AddressIsValid(address); err != nil {
		return nil, err
	}

	err := address.Upsert(as.Context(), transaction, true, []string{}, boil.Infer(), boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert address")
	}

	return address, nil
}

func (as *SqlAddressStore) Get(id string) (*model.Address, error) {
	address, err := model.FindAddress(as.Context(), as.GetReplica(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.TableNames.Addresses, id)
		}
	}
	return address, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option *model_helper.AddressFilterOptions) (model.AddressSlice, error) {
	var (
		res      []*model.Address
		db       = as.GetReplica()
		andConds = squirrel.And{}
	)
	if option.Conditions != nil {
		andConds = append(andConds, option.Conditions)
	}
	if option.UserID != nil {
		andConds = append(andConds, option.UserID)
		db = db.
			Joins(
				fmt.Sprintf(
					"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.UserAddressTableName, // 1
					model.AddressTableName,     // 2
					"address_id",               // 3
					model.AddressColumnId,      // 4
				),
			)
	}

	args, err := store.BuildSqlizer(andConds, "FilterByOption")
	if err != nil {
		return nil, err
	}
	err = db.Find(&res, args...).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to find addresses by given options")
	}
	return res, nil
}

func (as *SqlAddressStore) DeleteAddresses(transaction boil.ContextTransactor, addressIDs []string) *model.AppError {
	if transaction == nil {
		transaction = as.GetMaster()
	}
	err := transaction.Delete(&model.Address{}, "Id IN ?", addressIDs).Error
	if err != nil {
		return model.NewAppError("store.DeleteAddresses", "app.account.delete_addresses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// a := models.Address{ID: }
	addrs := models.AddressSlice{}
	qm.WhereIn()

	models.Addresses(qm.WhereIn()).DeleteAll(context.Background())

	return nil
}
