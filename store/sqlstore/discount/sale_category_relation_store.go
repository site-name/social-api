package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlSaleCategoryRelationStore struct {
	store.Store
}

func NewSqlSaleCategoryRelationStore(s store.Store) store.SaleCategoryRelationStore {
	return &SqlSaleCategoryRelationStore{s}
}

func (s *SqlSaleCategoryRelationStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id", "SaleID", "CategoryID", "CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given sale-category relation into database
func (ss *SqlSaleCategoryRelationStore) Save(relation *model.SaleCategoryRelation) (*model.SaleCategoryRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.SaleCategoryRelationTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"

	if _, err := ss.GetMasterX().NamedExec(query, relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "CategoryID", "salecategories_saleid_categoryid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleCategoryRelationTableName, "SaleID/CategoryID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-category relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get returns 1 sale-category relation with given id
func (ss *SqlSaleCategoryRelationStore) Get(relationID string) (*model.SaleCategoryRelation, error) {
	var res model.SaleCategoryRelation

	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.SaleCategoryRelationTableName+" WHERE Id = ?", relationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleCategoryRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-category relation with id=%s", relationID)
	}

	return &res, nil
}

// SaleCategoriesByOption returns a slice of sale-category relations with given option
func (ss *SqlSaleCategoryRelationStore) SaleCategoriesByOption(option *model.SaleCategoryRelationFilterOption) ([]*model.SaleCategoryRelation, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleCategoryRelationTableName).
		OrderBy(store.TableOrderingMap[store.SaleCategoryRelationTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.SaleID != nil {
		query = query.Where(option.SaleID)
	}
	if option.CategoryID != nil {
		query = query.Where(option.CategoryID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SaleCategoriesByOption_ToSql")
	}

	var res []*model.SaleCategoryRelation
	err = ss.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-category relations with given option")
	}

	return res, nil
}
