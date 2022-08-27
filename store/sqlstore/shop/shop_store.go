package shop

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/store"
)

type SqlShopStore struct {
	store.Store
}

func NewSqlShopStore(s store.Store) store.ShopStore {
	return &SqlShopStore{s}
}

func (s *SqlShopStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"OwnerID",
		"CreateAt",
		"UpdateAt",
		"Name",
		"Description",
		"TopMenuID",
		"IncludeTaxesInPrice",
		"DisplayGrossPrices",
		"ChargeTaxesOnShipping",
		"TrackInventoryByDefault",
		"DefaultWeightUnit",
		"AutomaticFulfillmentDigitalProducts",
		"DefaultDigitalMaxDownloads",
		"DefaultDigitalUrlValidDays",
		"AddressID",
		"DefaultMailSenderName",
		"DefaultMailSenderAddress",
		"CustomerSetPasswordUrl",
		"AutomaticallyConfirmAllNewOrders",
		"FulfillmentAutoApprove",
		"FulfillmentAllowUnPaid",
		"GiftcardExpiryType",
		"GiftcardExpiryPeriodType",
		"GiftcardExpiryPeriod",
		"AutomaticallyFulfillNonShippableGiftcard",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Upsert depends on shop's Id to decide to update/insert the given shop.
func (ss *SqlShopStore) Upsert(shopInstance *shop.Shop) (*shop.Shop, error) {
	var saving bool
	if shopInstance.Id == "" {
		saving = true
		shopInstance.PreSave()
	} else {
		shopInstance.PreUpdate()
	}

	if err := shopInstance.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)
	if saving {
		query := "INSERT INTO " + store.ShopTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
		_, err = ss.GetMasterX().NamedExec(query, shopInstance)

	} else {
		query := "UPDATE " + store.ShopTableName + " SET " + ss.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = ss.GetMasterX().NamedExec(query, shopInstance)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert shop with id=%s", shopInstance.Id)
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("multiple shops updated: %d instead of 1", numUpdated)
	}

	return shopInstance, nil
}

// Get finds a shop with given id and returns it
func (ss *SqlShopStore) Get(shopID string) (*shop.Shop, error) {
	var res shop.Shop
	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.ShopTableName+" WHERE Id = ?", shopID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopTableName, shopID)
		}
		return nil, errors.Wrapf(err, "failed to find shop with id=%s", shopID)
	}

	return &res, nil
}

func (ss *SqlShopStore) commonQueryBuilder(options *shop.ShopFilterOptions) (string, []interface{}, error) {
	query := ss.GetQueryBuilder().Select("*").From(store.ShopTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.OwnerID != nil {
		query = query.Where(options.OwnerID)
	}
	if options.Name != nil {
		query = query.Where(options.Name)
	}

	return query.ToSql()
}

// FilterByOptions finds and returns shops with given options
func (ss *SqlShopStore) FilterByOptions(options *shop.ShopFilterOptions) ([]*shop.Shop, error) {
	var res []*shop.Shop

	queryString, args, err := ss.commonQueryBuilder(options)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	err = ss.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shops with given options")
	}

	return res, nil
}

// GetByOptions finds and returns 1 shop with given options
func (ss *SqlShopStore) GetByOptions(options *shop.ShopFilterOptions) (*shop.Shop, error) {
	queryString, args, err := ss.commonQueryBuilder(options)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res shop.Shop
	err = ss.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find shop with given options")
	}

	return &res, nil
}
