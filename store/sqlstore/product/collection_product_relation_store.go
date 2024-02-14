package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCollectionProductStore struct {
	store.Store
}

func NewSqlCollectionProductStore(s store.Store) store.CollectionProductStore {
	return &SqlCollectionProductStore{s}
}

func (ps *SqlCollectionProductStore) BulkSave(transaction *gorm.DB, relations []*model.CollectionProduct) ([]*model.CollectionProduct, error) {
	if transaction == nil {
		transaction = ps.GetMaster()
	}

	for _, rel := range relations {
		err := transaction.Save(rel).Error

		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"ProductID", "CollectionID", "productcollections_collectionid_productid_key"}) {
				return nil, store.NewErrInvalidInput(model.CollectionProductRelationTableName, "CollectionID/ProductID", nil)
			}
			return nil, errors.Wrap(err, "failed to insert a collection product relation")
		}
	}

	return relations, nil
}

func (ps *SqlCollectionProductStore) FilterByOptions(options *model.CollectionProductFilterOptions) ([]*model.CollectionProduct, error) {
	selectFields := []string{model.CollectionProductRelationTableName + ".*"}
	if options.SelectRelatedCollection {
		selectFields = append(selectFields, model.CollectionTableName+".*")
	}
	if options.SelectRelatedProduct {
		selectFields = append(selectFields, model.ProductTableName+".*")
	}

	query := ps.GetQueryBuilder().
		Select(selectFields...).
		From(model.CollectionProductRelationTableName).
		Where(options.Conditions)

	if options.SelectRelatedCollection {
		query = query.InnerJoin(model.CollectionTableName + " ON Collections.Id = ProductCollections.CollectionID")
	}
	if options.SelectRelatedProduct {
		query = query.InnerJoin(model.ProductTableName + " ON Products.Id = ProductCollections.ProductID")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}
	rows, err := ps.GetReplica().Raw(queryString, args...).Rows()
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

func (s *SqlCollectionProductStore) Delete(transaction *gorm.DB, options *model.CollectionProductFilterOptions) error {
	query := s.GetQueryBuilder().Delete(model.CollectionProductRelationTableName).Where(options.Conditions)
	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	if transaction == nil {
		transaction = s.GetMaster()
	}

	err = transaction.Raw(queryString, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete collection product relations")
	}

	return nil
}
