package shop

// import (
// 	"database/sql"

// 	"github.com/pkg/errors"
// 	"github.com/sitename/sitename/model"
// 	"github.com/sitename/sitename/modules/util"
// 	"github.com/sitename/sitename/store"
// 	"gorm.io/gorm"
// )

// type SqlShopStore struct {
// 	store.Store
// }

// func NewSqlShopStore(s store.Store) store.ShopStore {
// 	return &SqlShopStore{s}
// }

// var shopModelFields = util.AnyArray[string]{
// 	"Id",
// 	"CreateAt",
// 	"UpdateAt",
// 	"Name",
// 	"HeaderText",
// 	"Description",
// 	"TopMenuID",
// 	"BottomMenuID",
// 	"IncludeTaxesInPrice",
// 	"DisplayGrossPrices",
// 	"ChargeTaxesOnShipping",
// 	"TrackInventoryByDefault",
// 	"DefaultWeightUnit",
// 	"AutomaticFulfillmentDigitalProducts",
// 	"DefaultDigitalMaxDownloads",
// 	"DefaultDigitalUrlValidDays",
// 	"AddressID",
// 	"CompanyAddressID",
// 	"DefaultMailSenderName",
// 	"DefaultMailSenderAddress",
// 	"CustomerSetPasswordUrl",
// 	"AutomaticallyConfirmAllNewOrders",
// 	"FulfillmentAutoApprove",
// 	"FulfillmentAllowUnPaid",
// 	"GiftcardExpiryType",
// 	"GiftcardExpiryPeriodType",
// 	"GiftcardExpiryPeriod",
// 	"AutomaticallyFulfillNonShippableGiftcard",
// }

// func (s *SqlShopStore) ModelFields(prefix string) util.AnyArray[string] {
// 	if prefix == "" {
// 		return shopModelFields
// 	}

// 	return shopModelFields.Map(func(_ int, s string) string {
// 		return prefix + s
// 	})
// }

// func (s *SqlShopStore) ScanFields(shop *model.Shop) []interface{} {
// 	return []interface{}{
// 		&shop.Id,
// 		&shop.CreateAt,
// 		&shop.UpdateAt,
// 		&shop.Name,
// 		&shop.HeaderText,
// 		&shop.Description,
// 		&shop.TopMenuID,
// 		&shop.BottomMenuID,
// 		&shop.IncludeTaxesInPrice,
// 		&shop.DisplayGrossPrices,
// 		&shop.ChargeTaxesOnShipping,
// 		&shop.TrackInventoryByDefault,
// 		&shop.DefaultWeightUnit,
// 		&shop.AutomaticFulfillmentDigitalProducts,
// 		&shop.DefaultDigitalMaxDownloads,
// 		&shop.DefaultDigitalUrlValidDays,
// 		&shop.AddressID,
// 		&shop.CompanyAddressID,
// 		&shop.DefaultMailSenderName,
// 		&shop.DefaultMailSenderAddress,
// 		&shop.CustomerSetPasswordUrl,
// 		&shop.AutomaticallyConfirmAllNewOrders,
// 		&shop.FulfillmentAutoApprove,
// 		&shop.FulfillmentAllowUnPaid,
// 		&shop.GiftcardExpiryType,
// 		&shop.GiftcardExpiryPeriodType,
// 		&shop.GiftcardExpiryPeriod,
// 		&shop.AutomaticallyFulfillNonShippableGiftcard,
// 	}
// }

// // Upsert depends on shop's Id to decide to update/insert the given shop.
// func (ss *SqlShopStore) Upsert(shopInstance *model.Shop) (*model.Shop, error) {
// 	var saving bool
// 	if shopInstance.Id == "" {
// 		saving = true
// 		shopInstance.PreSave()
// 	} else {
// 		shopInstance.PreUpdate()
// 	}

// 	if err := shopInstance.IsValid(); err != nil {
// 		return nil, err
// 	}

// 	var (
// 		err        error
// 		numUpdated int64
// 	)
// 	if saving {
// 		query := "INSERT INTO " + model.ShopTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
// 		_, err = ss.GetMaster().NamedExec(query, shopInstance)

// 	} else {
// 		query := "UPDATE " + model.ShopTableName + " SET " + ss.
// 			ModelFields("").
// 			Map(func(_ int, s string) string {
// 				return s + "=:" + s
// 			}).
// 			Join(",") + " WHERE Id=:Id"

// 		var result sql.Result
// 		result, err = ss.GetMaster().NamedExec(query, shopInstance)
// 		if err == nil && result != nil {
// 			numUpdated, _ = result.RowsAffected()
// 		}
// 	}

// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to upsert shop with id=%s", shopInstance.Id)
// 	}

// 	if numUpdated > 1 {
// 		return nil, errors.Errorf("multiple shops updated: %d instead of 1", numUpdated)
// 	}

// 	return shopInstance, nil
// }

// // Get finds a shop with given id and returns it
// func (ss *SqlShopStore) Get(shopID string) (*model.Shop, error) {
// 	var res model.Shop
// 	err := ss.GetReplica().First(&res, "Id = ?", shopID).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.ShopTableName, shopID)
// 		}
// 		return nil, errors.Wrapf(err, "failed to find shop with id=%s", shopID)
// 	}

// 	return &res, nil
// }

// func (ss *SqlShopStore) commonQueryBuilder(options *model.ShopFilterOptions) (string, []interface{}, error) {
// 	selectFields := ss.ModelFields("Shops.")
// 	if options.SelectRelatedCompanyAddress {
// 		selectFields = append(selectFields, ss.Address().ModelFields("Addresses.")...)
// 	}
// 	query := ss.GetQueryBuilder().Select(selectFields...).From(model.ShopTableName)

// 	if options.Id != nil {
// 		query = query.Where(options.Id)
// 	}
// 	if options.OwnerID != nil {
// 		query = query.Where(options.OwnerID)
// 	}
// 	if options.Name != nil {
// 		query = query.Where(options.Name)
// 	}
// 	if options.SelectRelatedCompanyAddress {
// 		query = query.InnerJoin(model.AddressTableName + " ON Addresses.Id = Shops.CompanyAddressID")
// 	}

// 	return query.ToSql()
// }

// // FilterByOptions finds and returns shops with given options
// func (ss *SqlShopStore) FilterByOptions(options *model.ShopFilterOptions) ([]*model.Shop, error) {
// 	queryString, args, err := ss.commonQueryBuilder(options)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
// 	}

// 	rows, err := ss.GetReplica().Query(queryString, args...)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to find shops with given options")
// 	}
// 	defer rows.Close()

// 	var res []*model.Shop

// 	for rows.Next() {
// 		var shop model.Shop
// 		var address model.Address
// 		var scanFields = ss.ScanFields(&shop)
// 		if options.SelectRelatedCompanyAddress {
// 			scanFields = append(scanFields, ss.Address().ScanFields(&address)...)
// 		}

// 		err = rows.Scan(scanFields...)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "error scanning shops")
// 		}

// 		if options.SelectRelatedCompanyAddress {
// 			shop.SetCompanyAddress(&address)
// 		}
// 		res = append(res, &shop)
// 	}

// 	return res, nil
// }

// // GetByOptions finds and returns 1 shop with given options
// func (ss *SqlShopStore) GetByOptions(options *model.ShopFilterOptions) (*model.Shop, error) {
// 	queryString, args, err := ss.commonQueryBuilder(options)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
// 	}

// 	row := ss.GetReplica().QueryRow(queryString, args...)

// 	var res model.Shop
// 	var address model.Address
// 	var scanFields = ss.ScanFields(&res)
// 	if options.SelectRelatedCompanyAddress {
// 		scanFields = append(scanFields, ss.Address().ScanFields(&address)...)
// 	}

// 	err = row.Scan(scanFields...)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, store.NewErrNotFound(model.ShopTableName, "options")
// 		}
// 		return nil, errors.Wrap(err, "failed to find shop with given options")
// 	}
// 	if options.SelectRelatedCompanyAddress {
// 		res.SetCompanyAddress(&address)
// 	}

// 	return &res, nil
// }
