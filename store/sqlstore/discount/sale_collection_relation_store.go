package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleCollectionRelationStore struct {
	store.Store
}

func NewSqlSaleCollectionRelationStore(s store.Store) store.SaleCollectionRelationStore {
	ss := &SqlSaleCollectionRelationStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleCollectionRelation{}, store.SaleCollectionRelationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("SaleID", "CollectionID")
	}

	return ss
}

func (ss *SqlSaleCollectionRelationStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.SaleCollectionRelationTableName, "SaleID", store.SaleTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.SaleCollectionRelationTableName, "CollectionID", store.ProductCollectionTableName, "Id", false)
}

// Save insert given sale-collection relation into database
func (ss *SqlSaleCollectionRelationStore) Save(relation *product_and_discount.SaleCollectionRelation) (*product_and_discount.SaleCollectionRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	if err := ss.GetMaster().Insert(relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "CollectionID", "salecollections_saleid_collectionid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleCollectionRelationTableName, "SaleID/CollectionID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-collection relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get finds and returns a sale-collection relation with given id
func (ss *SqlSaleCollectionRelationStore) Get(relationID string) (*product_and_discount.SaleCollectionRelation, error) {
	var res product_and_discount.SaleCollectionRelation
	err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.SaleCollectionRelationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": relationID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleCollectionRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-collection relation with id=%s", relationID)
	}

	return &res, nil
}

// FilterByOption returns a list of collections filtered based on given option
func (ss *SqlSaleCollectionRelationStore) FilterByOption(option *product_and_discount.SaleCollectionRelationFilterOption) ([]*product_and_discount.SaleCollectionRelation, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleCollectionRelationTableName).
		OrderBy(store.TableOrderingMap[store.SaleCollectionRelationTableName])

	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}

	if option.SaleID != nil {
		query = query.Where(option.SaleID.ToSquirrel("SaleID"))
	}

	if option.CollectionID != nil {
		query = query.Where(option.CollectionID.ToSquirrel("CollectionID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*product_and_discount.SaleCollectionRelation
	_, err = ss.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-collection relations with given option")
	}

	return res, nil
}
