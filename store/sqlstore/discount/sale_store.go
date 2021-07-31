package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDiscountSaleStore struct {
	store.Store
}

func NewSqlDiscountSaleStore(sqlStore store.Store) store.DiscountSaleStore {
	ss := &SqlDiscountSaleStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.Sale{}, store.SaleTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShopID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(product_and_discount.SALE_NAME_MAX_LENGTH)
		table.ColMap("Type").SetMaxSize(10)
	}
	return ss
}

func (ss *SqlDiscountSaleStore) CreateIndexesIfNotExists() {
	ss.CreateIndexIfNotExists("idx_sales_name", store.SaleTableName, "Name")
	ss.CreateIndexIfNotExists("idx_sales_type", store.SaleTableName, "Type")
	ss.CreateForeignKeyIfNotExists(store.SaleTableName, "ShopID", store.ShopTableName, "Id", false)
}

// Upsert bases on sale's Id to decide to update or insert given sale
func (ss *SqlDiscountSaleStore) Upsert(sale *product_and_discount.Sale) (*product_and_discount.Sale, error) {
	var saving bool
	if sale.Id == "" {
		saving = true
		sale.PreSave()
	} else {
		sale.PreUpdate()
	}

	if err := sale.IsValid(); err != nil {
		return nil, err
	}

	var (
		err           error
		numberUpdated int64
		oldSale       *product_and_discount.Sale
	)
	if saving {
		err = ss.GetMaster().Insert(sale)
	} else {
		oldSale, err = ss.Get(sale.Id)
		if err != nil {
			return nil, err
		}
		sale.ShopID = oldSale.ShopID // shop id CANNOT be edited

		numberUpdated, err = ss.GetMaster().Update(sale)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert sale with id=%s", sale.Id)
	}
	if numberUpdated > 1 {
		return nil, errors.Errorf("multiple sales were updated: %d instead of 1", numberUpdated)
	}

	return sale, nil
}

// Get finds and returns a sale with given saleID
func (ss *SqlDiscountSaleStore) Get(saleID string) (*product_and_discount.Sale, error) {
	result, err := ss.GetReplica().Get(product_and_discount.Sale{}, saleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleTableName, saleID)
		}
		return nil, errors.Wrapf(err, "failed to finds sale with id=%s", saleID)
	}

	return result.(*product_and_discount.Sale), nil
}

// FilterSalesByOption filter sales by option
func (ss *SqlDiscountSaleStore) FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, error) {
	query := ss.
		GetQueryBuilder().
		Select("*").
		From(store.SaleTableName).
		OrderBy("CreateAt ASC")

	// check shop id
	query = query.Where(option.ShopID.ToSquirrel("ShopID"))

	// check sale start date
	if option.StartDate != nil {
		query = query.Where(option.StartDate.ToSquirrel("StartDate"))
	}

	// check sale end date
	if option.EndDate != nil {
		query = query.Where(option.EndDate.ToSquirrel("EndDate"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query_to_sql")
	}

	var sales []*product_and_discount.Sale
	_, err = ss.GetReplica().Select(&sales, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleTableName, "")
		}
		return nil, errors.Wrap(err, "failed to find sales with given condition.")
	}

	return sales, nil
}
