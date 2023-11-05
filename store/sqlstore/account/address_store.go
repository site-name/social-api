package account

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAddressStore struct {
	store.Store
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) ScanFields(addr *model.Address) []interface{} {
	return []interface{}{
		&addr.Id,
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
		&addr.CreateAt,
		&addr.UpdateAt,
	}
}

func (as *SqlAddressStore) Upsert(transaction *gorm.DB, address *model.Address) (*model.Address, error) {
	if transaction == nil {
		transaction = as.GetMaster()
	}

	var err = transaction.Save(address).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert address")
	}

	return address, nil
}

func (as *SqlAddressStore) Get(id string) (*model.Address, error) {
	var res model.Address
	err := as.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AddressTableName, id)
		}
		return nil, errors.Wrap(err, "failed to find address with id="+id)
	}
	return &res, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option *model.AddressFilterOption) ([]*model.Address, error) {
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

func (as *SqlAddressStore) DeleteAddresses(transaction *gorm.DB, addressIDs []string) *model.AppError {
	if transaction == nil {
		transaction = as.GetMaster()
	}
	err := transaction.Delete(&model.Address{}, "Id IN ?", addressIDs).Error
	if err != nil {
		return model.NewAppError("store.DeleteAddresses", "app.account.delete_addresses.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
