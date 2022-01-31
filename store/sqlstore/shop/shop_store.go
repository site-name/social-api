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
	ss := &SqlShopStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shop.Shop{}, store.ShopTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("OwnerID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("TopMenuID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("AddressID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(shop.SHOP_NAME_MAX_LENGTH)
		table.ColMap("Description").SetMaxSize(shop.SHOP_DESCRIPTION_MAX_LENGTH)
		table.ColMap("DefaultWeightUnit").SetMaxSize(shop.SHOP_DEFAULT_WEIGHT_UNIT_MAX_LENGTH)
		table.ColMap("DefaultMailSenderName").SetMaxSize(shop.SHOP_DEFAULT_MAX_EMAIL_DISPLAY_NAME_LENGTH)
		table.ColMap("GiftcardExpiryType").SetMaxSize(shop.SHOP_GIFTCARD_EXPIRY_TYPE_MAX_LENGTH)
		table.ColMap("GiftcardExpiryPeriodType").SetMaxSize(shop.SHOP_GIFTCARD_EXPIRY_PERIOD_TYPE_MAX_LENGTH)
	}
	return ss
}

func (ss *SqlShopStore) CreateIndexesIfNotExists() {
	ss.CreateIndexIfNotExists("idx_shops_name", store.ShopTableName, "Name")
	ss.CreateIndexIfNotExists("idx_shops_name_lower_textpattern", store.ShopTableName, "lower(Name) text_pattern_ops")
	ss.CreateIndexIfNotExists("idx_shops_description", store.ShopTableName, "Description")
	ss.CreateIndexIfNotExists("idx_shops_description_lower_textpattern", store.ShopTableName, "lower(Description) text_pattern_ops")

	ss.CreateFullTextIndexIfNotExists("idx_shops_description", store.ShopTableName, "Description")
	ss.CreateForeignKeyIfNotExists(store.ShopTableName, "TopMenuID", store.MenuTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.ShopTableName, "OwnerID", store.UserTableName, "Id", true)
	ss.CreateForeignKeyIfNotExists(store.ShopTableName, "AddressID", store.AddressTableName, "Id", false)
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
		oldShop    *shop.Shop
	)
	if saving {
		err = ss.GetMaster().Insert(shopInstance)
	} else {
		// validate there is a shop with this id exists
		oldShop, err = ss.Get(shopInstance.Id)
		if err != nil {
			return nil, err
		}
		shopInstance.CreateAt = oldShop.CreateAt
		shopInstance.UpdateAt = model.GetMillis()
		numUpdated, err = ss.GetMaster().Update(shopInstance)
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
	err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ShopTableName+" WHERE Id = :ID", map[string]interface{}{"ID": shopID})
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

	_, err = ss.GetReplica().Select(&res, queryString, args...)
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
	err = ss.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShopTableName, "options")
		}
		return nil, errors.Wrap(err, "failed to find shop with given options")
	}

	return &res, nil
}
