package checkout

import (
	"github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	cls := &SqlCheckoutLineStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(checkout.CheckoutLine{}, store.CheckoutLineTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CheckoutID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("VariantID").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return cls
}

func (cls *SqlCheckoutLineStore) ModelFields() []string {
	return []string{
		"CheckoutLines.Id",
		"CheckoutLines.CreateAt",
		"CheckoutLines.CheckoutID",
		"CheckoutLines.VariantID",
		"CheckoutLines.Quantity",
	}
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

func (cls *SqlCheckoutLineStore) CreateIndexesIfNotExists() {
	cls.CreateIndexIfNotExists("idx_checkoutlines_checkout_id", store.CheckoutLineTableName, "CheckoutID")
	cls.CreateIndexIfNotExists("idx_checkoutlines_variant_id", store.CheckoutLineTableName, "VariantID")

	// foreign keys:
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "CheckoutID", store.CheckoutTableName, "Token", true)
	cls.CreateForeignKeyIfNotExists(store.CheckoutLineTableName, "VariantID", store.ProductVariantTableName, "Id", true)
}

func (cls *SqlCheckoutLineStore) Upsert(checkoutLine *checkout.CheckoutLine) (*checkout.CheckoutLine, error) {
	var isSave bool

	if checkoutLine.Id == "" {
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
		err = cls.GetMaster().Insert(checkoutLine)
	} else {
		numUpdated, err = cls.GetMaster().Update(checkoutLine)
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
	err := cls.GetReplica().SelectOne(&res, "SELECT * FROM "+store.CheckoutLineTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to to find checkout line with id=%s", id)
	}

	return &res, nil
}

func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutID(checkoutToken string) ([]*checkout.CheckoutLine, error) {
	var res []*checkout.CheckoutLine
	_, err := cls.GetReplica().Select(
		&res,
		`SELECT * FROM `+store.CheckoutLineTableName+` AS CkL 
		INNER JOIN `+store.CheckoutTableName+` AS Ck ON (
			Ck.Id = CkL.CheckoutID
		)
		WHERE (
			CkL.CheckoutID = :CheckoutID
		) 
		ORDER BY Ck.CreateAt ASC`,
		map[string]interface{}{"CheckoutID": checkoutToken},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get checkout lines belong to checkout with id=%s", checkoutToken)
	}

	return res, nil
}

func (cls *SqlCheckoutLineStore) DeleteLines(transaction *gorp.Transaction, ids []string) error {
	var executor squirrel.Execer = cls.GetMaster()
	if transaction != nil {
		executor = transaction
	}

	result, err := executor.Exec("DELETE FROM "+store.CheckoutLineTableName+" WHERE Id IN :IDs", map[string]interface{}{"IDs": ids})
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
	for _, line := range lines {
		if line == nil || line.IsValid() != nil {
			return store.NewErrInvalidInput(store.CheckoutLineTableName, "lines", "nil value")
		}
	}

	tx, err := cls.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, line := range lines {
		numUpdated, err := tx.Update(line)
		if err != nil {
			return errors.Wrapf(err, "failed to update checkout line with id=%s", line.Id)
		}
		if numUpdated > 1 {
			return errors.Errorf("multiple checkout lines updated: %d instead of 1", numUpdated)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkCreate(lines []*checkout.CheckoutLine) ([]*checkout.CheckoutLine, error) {
	tx, err := cls.GetMaster().Begin()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(tx)

	for _, line := range lines {
		if line != nil {
			line.PreSave()
			if appErr := line.IsValid(); appErr != nil {
				return nil, appErr
			}
			err = tx.Insert(line)
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
		cls.ModelFields(),
		append(
			cls.ProductVariant().ModelFields(),
			cls.Product().ModelFields()...,
		)...,
	)

	rows, err := cls.
		GetQueryBuilder().
		Select(selectFields...).
		From(store.CheckoutLineTableName).
		InnerJoin(store.ProductVariantTableName + " ON CheckoutLines.VariantID = ProductVariants.Id").
		InnerJoin(store.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
		Where(squirrel.Eq{"CheckoutLines.CheckoutID": checkoutToken}).
		RunWith(cls.GetReplica()).
		Query()

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
		scanFields      = append(
			cls.ScanFields(checkoutLine),
			append(cls.ProductVariant().ScanFields(productVariant), cls.Product().ScanFields(product)...)...,
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

	rows, err := cls.
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
		RunWith(cls.GetReplica()).
		Query()

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
	_, err = cls.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find checkout lines by given options")
	}

	return res, nil
}
