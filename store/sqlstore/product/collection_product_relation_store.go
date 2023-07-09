package product

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlCollectionProductStore struct {
	store.Store
}

func NewSqlCollectionProductStore(s store.Store) store.CollectionProductStore {
	return &SqlCollectionProductStore{s}
}

func (ps *SqlCollectionProductStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CollectionID",
		"ProductID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ps *SqlCollectionProductStore) ScanFields(rel *model.CollectionProduct) []interface{} {
	return []interface{}{
		&rel.Id,
		&rel.CollectionID,
		&rel.ProductID,
	}
}

func (ps *SqlCollectionProductStore) BulkSave(transaction store_iface.SqlxTxExecutor, relations []*model.CollectionProduct) ([]*model.CollectionProduct, error) {
	runner := ps.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	for _, rel := range relations {
		if !model.IsValidId(rel.Id) {
			rel.Id = ""
			rel.PreSave()
		}

		if err := rel.IsValid(); err != nil {
			return nil, err
		}

		result, err := runner.Exec("INSERT INTO "+store.CollectionProductRelationTableName+" (Id, CollectionID, ProductID) VALUES (:Id, :CollectionID, :ProductID)", rel)
		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "CollectionID", "productcollections_collectionid_productid_key"}) {
				return nil, store.NewErrInvalidInput(store.CollectionProductRelationTableName, "CollectionID/ProductID", nil)
			}
			return nil, errors.Wrap(err, "failed to insert a collection product relation")
		}

		rowsAdded, _ := result.RowsAffected()
		if rowsAdded != 1 {
			return nil, errors.Errorf("%d relation(s) was/were added instead of 1", rowsAdded)
		}
	}

	return relations, nil
}

func (ps *SqlCollectionProductStore) FilterByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, error) {
	selectFields := ps.ModelFields(store.CollectionProductRelationTableName + ".")
	if options.SelectRelatedCollection {
		selectFields = append(selectFields, ps.Collection().ModelFields(store.CollectionTableName+".")...)
	}
	if options.SelectRelatedProduct {
		selectFields = append(selectFields, ps.Product().ModelFields(store.ProductTableName+".")...)
	}

	query := ps.GetQueryBuilder().
		Select(selectFields...).
		From(store.CollectionProductRelationTableName)

	if options.ProductID != nil {
		query = query.Where(options.ProductID)
	}
	if options.CollectionID != nil {
		query = query.Where(options.CollectionID)
	}
	if options.SelectRelatedCollection {
		query = query.InnerJoin(store.CollectionTableName + " ON Collections.Id = ProductCollections.CollectionID")
	}
	if options.SelectRelatedProduct {
		query = query.InnerJoin(store.ProductTableName + " ON Products.Id = ProductCollections.ProductID")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}
	rows, err := ps.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product-collection relations")
	}
	defer rows.Close()

	var res []*model.CollectionProduct

	for rows.Next() {
		var (
			collectionProduct model.CollectionProduct
			collection        model.Collection
			product           model.Product
			scanFields        = ps.ScanFields(&collectionProduct)
		)
		if options.SelectRelatedCollection {
			scanFields = append(scanFields, ps.Collection().ScanFields(&collection)...)
		}
		if options.SelectRelatedProduct {
			scanFields = append(scanFields, ps.Product().ScanFields(&product)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product-collection relation")
		}

		if options.SelectRelatedCollection {
			collectionProduct.SetCollection(&collection)
		}
		if options.SelectRelatedProduct {
			collectionProduct.SetProduct(&product)
		}

		res = append(res, &collectionProduct)
	}

	return res, nil
}

func (s *SqlCollectionProductStore) Delete(transaction store_iface.SqlxTxExecutor, options *model.CollectionProductFilterOptions) error {
	query := s.GetQueryBuilder().Delete(store.CollectionProductRelationTableName)

	for _, opt := range []squirrel.Sqlizer{options.CollectionID, options.ProductID} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	_, err = runner.Exec(queryString, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete collection product relations")
	}

	return nil
}
