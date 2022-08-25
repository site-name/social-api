package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionProductStore struct {
	store.Store
}

func NewSqlCollectionProductStore(s store.Store) store.CollectionProductStore {
	return &SqlCollectionProductStore{s}
}

func (ps *SqlCollectionProductStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
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

func (ps *SqlCollectionProductStore) ScanFields(rel product_and_discount.CollectionProduct) []interface{} {
	return []interface{}{
		&rel.Id,
		&rel.CollectionID,
		&rel.ProductID,
	}
}

func (ps *SqlCollectionProductStore) FilterByOptions(options *product_and_discount.CollectionProductFilterOptions) ([]*product_and_discount.CollectionProduct, error) {
	selectFields := ps.ModelFields(store.CollectionProductRelationTableName + ".")
	if options.SelectRelatedCollection {
		selectFields = append(selectFields, ps.Collection().ModelFields(store.CollectionTableName+".")...)
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

	var (
		res               []*product_and_discount.CollectionProduct
		collectionProduct product_and_discount.CollectionProduct
		collection        product_and_discount.Collection
		scanFields        = ps.ScanFields(collectionProduct)
	)
	if options.SelectRelatedCollection {
		scanFields = append(scanFields, ps.Collection().ScanFields(collection)...)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}
	rows, err := ps.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product-collection relations")
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product-collection relation")
		}

		if options.SelectRelatedCollection {
			collectionProduct.Collection = &collection // no need to deep copy collection here, See Collection_Product relation DeepCopy for detail
		}

		res = append(res, collectionProduct.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return res, nil
}
