package checkout

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	return &SqlCheckoutLineStore{sqlStore}
}

func (cls *SqlCheckoutLineStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"CreateAt",
		"CheckoutID",
		"VariantID",
		"Quantity",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (cls *SqlCheckoutLineStore) ScanFields(line checkout.CheckoutLine) []interface{} {
	return []interface{}{
		&line.Id,
		&line.CreateAt,
		&line.CheckoutID,
		&line.VariantID,
		&line.Quantity,
	}
}

func (cls *SqlCheckoutLineStore) Upsert(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, error) {
	var isSave bool

	if !model.IsValidId(checkoutLine.Id) {
		checkoutLine.Id = ""
		isSave = true
	}

	checkoutLine.PreSave()
	if err := checkoutLine.IsValid(); err != nil {
		return nil, err
	}

	var (
		numUpdated int64
		err        error
	)
	if isSave {
		query := "INSERT INTO " + store.CheckoutLineTableName + "(" + cls.ModelFields("").Join(",") + ") VALUES (" + cls.ModelFields(":").Join(",") + ")"
		_, err = cls.GetMasterX().NamedExec(query, checkoutLine)

	} else {
		query := "UPDATE " + store.CheckoutLineTableName + " SET " + cls.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = cls.GetMasterX().NamedExec(query, checkoutLine)
		if err == nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert checkout line")
	}
	if numUpdated > 1 {
		return nil, errors.Errorf("%d checkout lines were updated instead of 1", numUpdated)
	}

	return checkoutLine, nil
}

func (cls *SqlCheckoutLineStore) Get(id string) (*checkout.CheckoutLine, error) {
	var res checkout.CheckoutLine

	err := cls.GetReplicaX().Get(&res, "SELECT * FROM "+store.CheckoutLineTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.CheckoutLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to to find checkout line with id=%s", id)
	}

	return &res, nil
}

func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutID(checkoutToken string) ([]*checkout.CheckoutLine, error) {
	var res []*checkout.CheckoutLine

	err := cls.GetReplicaX().Select(
		&res,
		`SELECT * FROM `+store.CheckoutLineTableName+`
		INNER JOIN `+store.CheckoutTableName+` ON	Checkouts.Id = CheckoutLines.CheckoutID
		WHERE	CheckoutLines.CheckoutID = ?
		ORDER BY Checkouts.CreateAt ASC`,
		checkoutToken,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get checkout lines belong to checkout with id=%s", checkoutToken)
	}

	return res, nil
}

func (cls *SqlCheckoutLineStore) DeleteLines(transaction store_iface.SqlxTxExecutor, ids []string) error {
	var executor store_iface.SqlxExecutor = cls.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	result, err := executor.Exec("DELETE FROM "+store.CheckoutLineTableName+" WHERE Id IN ?", ids)
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout lines")
	}
	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "failed to count number of checkout lines deleted")
	} else if rows != int64(len(ids)) {
		return errors.Errorf("expect %d checkout lines to be deleted but got %d", len(ids), rows)
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkUpdate(lines []*checkout.CheckoutLine) error {
	tx, err := cls.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	updateQuery := "UPDATE " + store.CheckoutLineTableName + " SET " + cls.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	for _, line := range lines {
		line.PreSave()
		if err := line.IsValid(); err != nil {
			return err
		}

		result, err := tx.NamedExec(updateQuery, line)
		if err != nil {
			return errors.Wrapf(err, "failed to update checkout line with id=%s", line.Id)
		}
		if numUpdated, _ := result.RowsAffected(); numUpdated > 1 {
			return errors.Errorf("multiple checkout lines updated: %d instead of 1", numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkCreate(lines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, error) {
	tx, err := cls.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	query := "INSERT INTO " + store.CheckoutLineTableName + "(" + cls.ModelFields("").Join(",") + ") VALUES (" + cls.ModelFields(":").Join(",") + ")"

	for _, line := range lines {
		if line != nil {
			line.PreSave()
			if appErr := line.IsValid(); appErr != nil {
				return nil, appErr
			}

			_, err = tx.NamedExec(query, line)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to save checkout line with id=%s", line.Id)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit_transaction")
	}

	return lines, nil
}

// CheckoutLinesByCheckoutWithPrefetch finds all checkout lines belong to given checkout
//
// and prefetch all related product variants, products
//
// this borrows the idea from Django's prefetch_related() method
func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutWithPrefetch(checkoutToken string) ([]*checkout.CheckoutLine, []*product_and_discount.ProductVariant, []*product_and_discount.Product, error) {
	selectFields := append(
		cls.ModelFields(store.CheckoutLineTableName+"."),
		append(
			cls.ProductVariant().ModelFields(),
			cls.Product().ModelFields()...,
		)...,
	)

	query, args, err := cls.
		GetQueryBuilder().
		Select(selectFields...).
		From(store.CheckoutLineTableName).
		InnerJoin(store.ProductVariantTableName + " ON CheckoutLines.VariantID = ProductVariants.Id").
		InnerJoin(store.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
		Where(squirrel.Eq{"CheckoutLines.CheckoutID": checkoutToken}).
		ToSql()

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "CheckoutLinesByCheckoutWithPrefetch_ToSql")
	}

	rows, err := cls.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to find checkout lines and prefetch related values, with checkoutToken=%s", checkoutToken)
	}

	var (
		checkoutLines   []*checkout.CheckoutLine
		productVariants []*product_and_discount.ProductVariant
		products        []*product_and_discount.Product
		checkoutLine    checkout.CheckoutLine
		productVariant  product_and_discount.ProductVariant
		product         product_and_discount.Product

		scanFields = append(
			cls.ScanFields(checkoutLine),
			append(
				cls.ProductVariant().ScanFields(productVariant),
				cls.Product().ScanFields(product)...,
			)...,
		)
	)

	for rows.Next() {
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to scan a row")
		}

		checkoutLines = append(checkoutLines, checkoutLine.DeepCopy())
		productVariants = append(productVariants, productVariant.DeepCopy())
		products = append(products, product.DeepCopy())
	}

	if err = rows.Close(); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to close rows")
	}

	return checkoutLines, productVariants, products, nil
}

// TotalWeightForCheckoutLines calculate total weight for given checkout lines
func (cls *SqlCheckoutLineStore) TotalWeightForCheckoutLines(checkoutLineIDs []string) (*measurement.Weight, error) {
	query, args, err := cls.
		GetQueryBuilder().
		Select(
			"CheckoutLines.Quantity",
			"ProductVariants.Weight",
			"ProductVariants.WeightUnit",
			"Products.Weight",
			"Products.WeightUnit",
			"ProductTypes.Weight",
			"ProductTypes.WeightUnit",
		).
		From(store.CheckoutLineTableName).
		InnerJoin(store.ProductVariantTableName + " ON (CheckoutLines.VariantID = ProductVariants.Id)").
		InnerJoin(store.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
		InnerJoin(store.ProductTypeTableName + " ON (ProductTypes.Id = Products.ProductTypeID)").
		Where(squirrel.Eq{"CheckoutLines.Id": checkoutLineIDs}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "TotalWeightForCheckoutLines_ToSql")
	}

	rows, err := cls.GetReplicaX().QueryX(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find values")
	}

	var (
		totalWeight  = measurement.ZeroWeight
		lineQuantity int
		weight       measurement.Weight

		variantWeightAmount *float32
		variantWeightUnit   measurement.WeightUnit

		productWeightAmount *float32
		productWeightUnit   measurement.WeightUnit

		productTypeWeightAmount *float32
		productTypeWeightUnit   measurement.WeightUnit
	)

	for rows.Next() {
		err = rows.Scan(
			&lineQuantity,
			&variantWeightAmount,
			&variantWeightUnit,
			&productWeightAmount,
			&productWeightUnit,
			&productTypeWeightAmount,
			&productTypeWeightUnit,
		)
		if err != nil {
			// return immediately if an error occured
			return nil, errors.Wrap(err, "failed to scan a row")
		}

		if variantWeightAmount != nil {
			weight = measurement.Weight{Amount: *variantWeightAmount, Unit: variantWeightUnit}
		} else if productWeightAmount != nil {
			weight = measurement.Weight{Amount: *variantWeightAmount, Unit: productWeightUnit}
		} else if productTypeWeightAmount != nil {
			weight = measurement.Weight{Amount: *productTypeWeightAmount, Unit: productTypeWeightUnit}
		}

		if weight.Amount != 0 && weight.Unit != "" {
			totalWeight, err = totalWeight.Add(weight.Mul(float32(lineQuantity)))
			if err != nil {
				return nil, err
			}
		}
	}

	if err = rows.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing rows")
	}

	return totalWeight, nil
}

// CheckoutLinesByOption finds and returns checkout lines filtered using given option
func (cls *SqlCheckoutLineStore) CheckoutLinesByOption(option *checkout.CheckoutLineFilterOption) ([]*checkout.CheckoutLine, error) {
	query := cls.GetQueryBuilder().
		Select("*").
		From(store.CheckoutLineTableName).
		OrderBy(store.TableOrderingMap[store.CheckoutLineTableName])

	// parse options
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.CheckoutID != nil {
		query = query.Where(option.CheckoutID)
	}
	if option.VariantID != nil {
		query = query.Where(option.VariantID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "CheckoutLinesByOption_ToSql")
	}

	var res []*checkout.CheckoutLine
	err = cls.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find checkout lines by given options")
	}

	return res, nil
}
