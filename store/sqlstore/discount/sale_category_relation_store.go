package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlSaleCategoryRelationStore struct {
	store.Store
}

func NewSqlSaleCategoryRelationStore(s store.Store) store.SaleCategoryRelationStore {
	ss := &SqlSaleCategoryRelationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.SaleCategoryRelation{}, store.SaleCategoryRelationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("SaleID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CategoryID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("SaleID", "CategoryID")
	}
	return ss
}

func (ss *SqlSaleCategoryRelationStore) TableName(withField string) string {
	name := "SaleCategories"
	if withField != "" {
		name += "." + withField
	}
	return name
}

func (ss *SqlSaleCategoryRelationStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.SaleCategoryRelationTableName, "SaleID", store.SaleTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.SaleCategoryRelationTableName, "CategoryID", store.ProductCategoryTableName, "Id", false)
}

// Save inserts given sale-category relation into database
func (ss *SqlSaleCategoryRelationStore) Save(relation *product_and_discount.SaleCategoryRelation) (*product_and_discount.SaleCategoryRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	if err := ss.GetMaster().Insert(relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "CategoryID", "salecategories_saleid_categoryid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleCategoryRelationTableName, "SaleID/CategoryID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-category relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get returns 1 sale-category relation with given id
func (ss *SqlSaleCategoryRelationStore) Get(relationID string) (*product_and_discount.SaleCategoryRelation, error) {
	var res product_and_discount.SaleCategoryRelation
	err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.SaleCategoryRelationTableName+" WHERE Id = :ID", map[string]interface{}{"ID": relationID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleCategoryRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-category relation with id=%s", relationID)
	}

	return &res, nil
}

// SaleCategoriesByOption returns a slice of sale-category relations with given option
func (ss *SqlSaleCategoryRelationStore) SaleCategoriesByOption(option *product_and_discount.SaleCategoryRelationFilterOption) ([]*product_and_discount.SaleCategoryRelation, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleCategoryRelationTableName).
		OrderBy(store.TableOrderingMap[store.SaleCategoryRelationTableName])

	// check id
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("Id"))
	}

	// check saleID
	if option.SaleID != nil {
		query = query.Where(option.SaleID.ToSquirrel("SaleID"))
	}

	// check categoryID
	if option.CategoryID != nil {
		query = query.Where(option.CategoryID.ToSquirrel("CategoryID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SaleCategoriesByOption_ToSql")
	}

	var res []*product_and_discount.SaleCategoryRelation
	_, err = ss.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-category relations with given option")
	}

	return res, nil
}
