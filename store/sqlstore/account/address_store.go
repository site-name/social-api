package account

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlAddressStore struct {
	store.Store
}

var modelFields = util.AnyArray[string]{
	"Id",
	"FirstName",
	"LastName",
	"CompanyName",
	"StreetAddress1",
	"StreetAddress2",
	"City",
	"CityArea",
	"PostalCode",
	"Country",
	"CountryArea",
	"Phone",
	"CreateAt",
	"UpdateAt",
}

// new address database store
func NewSqlAddressStore(sqlStore store.Store) store.AddressStore {
	return &SqlAddressStore{Store: sqlStore}
}

func (as *SqlAddressStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return modelFields
	}

	return modelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
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

	var result = transaction.Save(address)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to upsert address")
	}
	if result.RowsAffected != 1 {
		return nil, errors.Errorf("%d address(es) upserted instead of 1", result.RowsAffected)
	}

	return address, nil
}

func (as *SqlAddressStore) Get(addressID string) (*model.Address, error) {
	var res model.Address
	err := as.GetReplica().First(&res, "Id = ?", addressID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.AddressTableName, addressID)
		}
		return nil, errors.Wrap(err, "failed to find address with id="+addressID)
	}
	return &res, nil
}

// FilterByOption finds and returns a list of address(es) filtered by given option
func (as *SqlAddressStore) FilterByOption(option *model.AddressFilterOption) ([]*model.Address, error) {
	var res []*model.Address
	db := as.GetReplica().Table(model.AddressTableName)
	andConds := squirrel.And{}
	for _, opt := range []squirrel.Sqlizer{
		option.Id,
		option.UserID,
		option.Other,
	} {
		if opt != nil {
			andConds = append(andConds, opt)
		}
	}
	if option.OrderID != nil {
		andConds = append(andConds, option.OrderID.Id) //
		db = db.Joins("INNER JOIN "+model.OrderTableName+" ON Orders.? = Addresses.Id", option.OrderID.On)
	}
	if option.UserID != nil {
		db = db.Joins("INNER JOIN " + model.UserAddressTableName + " ON UserAddresses.address_id = Addresses.Id")
	}
	err := db.Find(&res, store.BuildSqlizer(andConds)...).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to find addresses by given options")
	}
	return res, nil
}

func (as *SqlAddressStore) DeleteAddresses(transaction *gorm.DB, addressIDs []string) error {
	if transaction == nil {
		transaction = as.GetMaster()
	}
	err := transaction.Delete(&model.Address{}, "Id IN ?", addressIDs).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete addresses with given ids")
	}

	return nil
}
