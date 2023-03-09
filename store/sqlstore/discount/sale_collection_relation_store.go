package discount

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlSaleCollectionRelationStore struct {
	store.Store
}

func NewSqlSaleCollectionRelationStore(s store.Store) store.SaleCollectionRelationStore {
	return &SqlSaleCollectionRelationStore{s}
}

func (s *SqlSaleCollectionRelationStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "SaleID", "CollectionID", "CreateAt",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save insert given sale-collection relation into database
func (ss *SqlSaleCollectionRelationStore) Save(relation *model.SaleCollectionRelation) (*model.SaleCollectionRelation, error) {
	relation.PreSave()
	if err := relation.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.SaleCollectionRelationTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"

	if _, err := ss.GetMasterX().NamedExec(query, relation); err != nil {
		if ss.IsUniqueConstraintError(err, []string{"SaleID", "CollectionID", "salecollections_saleid_collectionid_key"}) {
			return nil, store.NewErrInvalidInput(store.SaleCollectionRelationTableName, "SaleID/CollectionID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save sale-collection relation with id=%s", relation.Id)
	}

	return relation, nil
}

// Get finds and returns a sale-collection relation with given id
func (ss *SqlSaleCollectionRelationStore) Get(relationID string) (*model.SaleCollectionRelation, error) {
	var res model.SaleCollectionRelation

	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.SaleCollectionRelationTableName+" WHERE Id = ?", relationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.SaleCollectionRelationTableName, relationID)
		}
		return nil, errors.Wrapf(err, "failed to find sale-collection relation with id=%s", relationID)
	}

	return &res, nil
}

// FilterByOption returns a list of collections filtered based on given option
func (ss *SqlSaleCollectionRelationStore) FilterByOption(option *model.SaleCollectionRelationFilterOption) ([]*model.SaleCollectionRelation, error) {
	query := ss.GetQueryBuilder().
		Select("*").
		From(store.SaleCollectionRelationTableName).
		OrderBy(store.TableOrderingMap[store.SaleCollectionRelationTableName])

	if option.Id != nil {
		query = query.Where(option.Id)
	}

	if option.SaleID != nil {
		query = query.Where(option.SaleID)
	}

	if option.CollectionID != nil {
		query = query.Where(option.CollectionID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.SaleCollectionRelation
	err = ss.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sale-collection relations with given option")
	}

	return res, nil
}
