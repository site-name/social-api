package shipping

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodExcludedProductStore struct {
	store.Store
}

func NewSqlShippingMethodExcludedProductStore(s store.Store) store.ShippingMethodExcludedProductStore {
	ss := &SqlShippingMethodExcludedProductStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodExcludedProduct{}, store.ShippingMethodExcludedProductTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "ProductID")
	}

	return ss
}

func (ss *SqlShippingMethodExcludedProductStore) CreateIndexesIfNotExists() {
	ss.CreateForeignKeyIfNotExists(store.ShippingMethodExcludedProductTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", false)
	ss.CreateForeignKeyIfNotExists(store.ShippingMethodExcludedProductTableName, "ProductID", store.ProductTableName, "Id", false)
}

// Save inserts given ShippingMethodExcludedProduct into database then returns it
func (ss *SqlShippingMethodExcludedProductStore) Save(instance *shipping.ShippingMethodExcludedProduct) (*shipping.ShippingMethodExcludedProduct, error) {
	instance.PreSave()
	if err := instance.IsValid(); err != nil {
		return nil, err
	}

	err := ss.GetMaster().Insert(instance)
	if err != nil {
		if ss.IsUniqueConstraintError(err, []string{"ShippingMethodID", "ProductID", "shippingmethodexcludedproducts_shippingmethodid_productid_key"}) {
			return nil, store.NewErrInvalidInput(store.ShippingMethodExcludedProductTableName, "ShippingMethodID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save shipping method excluded product with id=%s", instance.Id)
	}

	return instance, nil
}

// Get finds and returns a shipping method excluded product with given id then reutrns it
func (ss *SqlShippingMethodExcludedProductStore) Get(id string) (*shipping.ShippingMethodExcludedProduct, error) {
	var res shipping.ShippingMethodExcludedProduct
	err := ss.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ShippingMethodExcludedProductTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodExcludedProductTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method excluded product with id=%s", id)
	}

	return &res, nil
}
