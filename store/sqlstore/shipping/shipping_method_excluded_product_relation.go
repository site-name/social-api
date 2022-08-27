package shipping

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodExcludedProductStore struct {
	store.Store
}

func NewSqlShippingMethodExcludedProductStore(s store.Store) store.ShippingMethodExcludedProductStore {
	return &SqlShippingMethodExcludedProductStore{s}
}

func (s *SqlShippingMethodExcludedProductStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"ShippingMethodID",
		"ProductID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

// Save inserts given ShippingMethodExcludedProduct into database then returns it
func (ss *SqlShippingMethodExcludedProductStore) Save(instance *shipping.ShippingMethodExcludedProduct) (*shipping.ShippingMethodExcludedProduct, error) {
	instance.PreSave()
	if err := instance.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ShippingMethodExcludedProductTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ")"
	_, err := ss.GetMasterX().NamedExec(query, instance)
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
	err := ss.GetReplicaX().Get(&res, "SELECT * FROM "+store.ShippingMethodExcludedProductTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ShippingMethodExcludedProductTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find shipping method excluded product with id=%s", id)
	}

	return &res, nil
}
