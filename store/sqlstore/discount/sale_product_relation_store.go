package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleProductRelationStore struct {
	store.Store
}

func NewSqlSaleProductRelationStore(s store.Store) store.SaleProductRelationStore {
	ss := &SqlSaleProductRelationStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleProductRelation{}, store.SaleProductRelationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("SaleID", "ProductID")
	}

	return ss
}

func (ss *SqlSaleProductRelationStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.SaleProductRelationTableName, "SaleID", store.SaleTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.SaleProductRelationTableName, "ProductID", store.ProductTableName, "Id", false)
}

// Save inserts given sale-product relation into database then returns it
func (ss *SqlSaleProductRelationStore) Save(relation *product_and_discount.SaleProductRelation) (*product_and_discount.SaleProductRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	if err := ss.GetMaster().Insert(relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "ProductID", "paleproducts_saleid_productid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleProductRelationTableName, "SaleID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-product relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get finds and returns a sale-product relation with given id
func (ss *SqlSaleProductRelationStore) Get(relationID string) (*product_and_discount.SaleProductRelation, error) {
	var res product_and_discount.SaleProductRelation
	err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.SaleProductRelationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": relationID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleProductRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-product relation with id=%s", relationID)
	}

	return &res, nil
}
