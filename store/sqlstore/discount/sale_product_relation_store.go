package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlSaleProductRelationStore struct {
	store.Store
}

func NewSqlSaleProductRelationStore(s store.Store) store.SaleProductRelationStore {
	return &SqlSaleProductRelationStore{s}
}

func (s *SqlSaleProductRelationStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id", "SaleID", "ProductID", "CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given sale-product relation into database then returns it
func (ss *SqlSaleProductRelationStore) Save(relation *model.SaleProductRelation) (*model.SaleProductRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.SaleProductRelationTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"

	if _, err := ss.GetMasterX().NamedExec(query, relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "ProductID", "paleproducts_saleid_productid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleProductRelationTableName, "SaleID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-product relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get finds and returns a sale-product relation with given id
func (ss *SqlSaleProductRelationStore) Get(relationID string) (*model.SaleProductRelation, error) {
	var res model.SaleProductRelation

	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.SaleProductRelationTableName+" WHERE Id = ?", relationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleProductRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-product relation with id=%s", relationID)
	}

	return &res, nil
}

// SaleProductsByOption returns a slice of sale-product relations, filtered by given option
func (ss *SqlSaleProductRelationStore) SaleProductsByOption(option *model.SaleProductRelationFilterOption) ([]*model.SaleProductRelation, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleProductRelationTableName).
		OrderBy(store.TableOrderingMap[store.SaleProductRelationTableName])

	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.SaleID != nil {
		query = query.Where(option.SaleID)
	}
	if option.ProductID != nil {
		query = query.Where(option.ProductID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SaleProductsByOption_ToSql")
	}

	var res []*model.SaleProductRelation
	err = ss.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-product relations with given option")
	}

	return res, nil
}
