package shipping

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlShippingMethodExcludedProductStore struct {
	store.Store
}

func NewSqlShippingMethodExcludedProductStore(s store.Store) store.ShippingMethodExcludedProductStore {
	return &SqlShippingMethodExcludedProductStore{s}
}

func (s *SqlShippingMethodExcludedProductStore) ScanFields(rel *model.ShippingMethodExcludedProduct) []any {
	return []any{
		&rel.Id,
		&rel.ShippingMethodID,
		&rel.ProductID,
	}
}

func (s *SqlShippingMethodExcludedProductStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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

// this constraint is defined in db/migrations/postgres/000075_create_shippingmethodexcludedproducts.up.sql
const shippingMethodExcludedProductUniqueConstraint = "shippingmethodexcludedproducts_shippingmethodid_productid_key"

// Save inserts given ShippingMethodExcludedProduct into database then returns it
func (ss *SqlShippingMethodExcludedProductStore) Save(instance *model.ShippingMethodExcludedProduct) (*model.ShippingMethodExcludedProduct, error) {
	instance.PreSave()
	if err := instance.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.ShippingMethodExcludedProductTableName + "(" + ss.ModelFields("").Join(",") + ") VALUES (" + ss.ModelFields(":").Join(",") + ") ON CONFLICT ON CONSTRAINT " + shippingMethodExcludedProductUniqueConstraint + " DO NOTHING"
	_, err := ss.GetMasterX().NamedExec(query, instance)
	if err != nil {
		if ss.IsUniqueConstraintError(err, []string{"ShippingMethodID", "ProductID", shippingMethodExcludedProductUniqueConstraint}) {
			return nil, store.NewErrInvalidInput(model.ShippingMethodExcludedProductTableName, "ShippingMethodID/ProductID", "duplicate")
		}
		return nil, errors.Wrapf(err, "failed to save shipping method excluded product with id=%s", instance.Id)
	}

	return instance, nil
}

func (s *SqlShippingMethodExcludedProductStore) FilterByOptions(options *model.ShippingMethodExcludedProductFilterOptions) ([]*model.ShippingMethodExcludedProduct, error) {
	selectFields := s.ModelFields(model.ShippingMethodExcludedProductTableName + ".")
	if options.SelectRelatedProduct {
		selectFields = append(selectFields, s.Product().ModelFields(model.ProductTableName+".")...)
	}

	query := s.GetQueryBuilder().Select(selectFields...).From(model.ShippingMethodExcludedProductTableName)

	for _, opt := range []squirrel.Sqlizer{
		options.Id,
		options.ProductID,
		options.ShippingMethodID,
	} {
		if opt != nil {
			query = query.Where(opt)
		}
	}
	if options.SelectRelatedProduct {
		query = query.InnerJoin(model.ProductTableName + " ON Products.Id = ShippingMethodExcludedProducts.ProductID")
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	rows, err := s.GetReplicaX().QueryX(queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method excluded products by given options")
	}
	defer rows.Close()

	var res []*model.ShippingMethodExcludedProduct

	for rows.Next() {
		var (
			rel        model.ShippingMethodExcludedProduct
			prd        model.Product
			scanFields = s.ScanFields(&rel)
		)
		if options.SelectRelatedProduct {
			scanFields = append(scanFields, s.Product().ScanFields(&prd)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of shipping method excluded product")
		}

		if options.SelectRelatedProduct {
			rel.SetProduct(&prd)
		}

		res = append(res, &rel)
	}

	return res, nil
}

func (s *SqlShippingMethodExcludedProductStore) Delete(transaction store_iface.SqlxExecutor, options *model.ShippingMethodExcludedProductFilterOptions) error {
	query := s.GetQueryBuilder().Delete(model.ShippingMethodExcludedProductTableName)

	for _, opt := range []squirrel.Sqlizer{
		options.Id,
		options.ProductID,
		options.ShippingMethodID,
	} {
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
		return errors.Wrap(err, "failed to delete shipping method excluded products")
	}

	return nil
}
