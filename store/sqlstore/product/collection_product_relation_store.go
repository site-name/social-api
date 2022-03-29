package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlCollectionProductStore struct {
	store.Store
}

func NewSqlCollectionProductStore(s store.Store) store.CollectionProductStore {
	cps := &SqlCollectionProductStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.CollectionProduct{}, store.CollectionProductRelationTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CollectionID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("CollectionID", "ProductID")
	}
	return cps
}

func (ps *SqlCollectionProductStore) ModelFields() []string {
	return []string{
		"ProductCollections.Id",
		"ProductCollections.CollectionID",
		"ProductCollections.ProductID",
	}
}

func (ps *SqlCollectionProductStore) ScanFields(rel product_and_discount.CollectionProduct) []interface{} {
	return []interface{}{
		&rel.Id,
		&rel.CollectionID,
		&rel.ProductID,
	}
}

func (ps *SqlCollectionProductStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.CollectionProductRelationTableName, "CollectionID", store.ProductCollectionTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.CollectionProductRelationTableName, "ProductID", store.ProductTableName, "Id", true)
}

func (ps *SqlCollectionProductStore) FilterByOptions(options *product_and_discount.CollectionProductFilterOptions) ([]*product_and_discount.CollectionProduct, error) {
	selectFields := ps.ModelFields()
	if options.SelectRelatedCollection {
		selectFields = append(selectFields, ps.Collection().ModelFields()...)
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
		query = query.InnerJoin(store.ProductCollectionTableName + " ON Collections.Id = ProductCollections.CollectionID")
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

	rows, err := query.RunWith(ps.GetReplica()).Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find product-collection relations")
	}

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of product-collection relation")
		}

		if options.SelectRelatedCollection {
			collectionProduct.Collection = &collection
		}

		res = append(res, collectionProduct.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close rows")
	}

	return res, nil
}
