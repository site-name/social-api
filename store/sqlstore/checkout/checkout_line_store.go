package checkout

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	return &SqlCheckoutLineStore{sqlStore}
}

func (cls *SqlCheckoutLineStore) ScanFields(line *model.CheckoutLine) []interface{} {
	return []interface{}{
		&line.Id,
		&line.CreateAt,
		&line.CheckoutID,
		&line.VariantID,
		&line.Quantity,
	}
}

func (cls *SqlCheckoutLineStore) Upsert(checkoutLine *model.CheckoutLine) (*model.CheckoutLine, error) {
	err := cls.GetMaster().Save(checkoutLine).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert checkout line")
	}
	return checkoutLine, nil
}

func (cls *SqlCheckoutLineStore) Get(id string) (*model.CheckoutLine, error) {
	var res model.CheckoutLine

	err := cls.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.CheckoutLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to to find checkout line with id=%s", id)
	}

	return &res, nil
}

func (cls *SqlCheckoutLineStore) DeleteLines(transaction *gorm.DB, ids []string) error {
	if transaction == nil {
		transaction = cls.GetMaster()
	}

	err := transaction.Raw("DELETE FROM "+model.CheckoutLineTableName+" WHERE Id IN ?", ids).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete checkout lines")
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkUpdate(lines []*model.CheckoutLine) error {
	for _, line := range lines {
		err := cls.GetMaster().Table(model.CheckoutLineTableName).Updates(line).Error
		if err != nil {
			return errors.Wrapf(err, "failed to update checkout line with id=%s", line.Id)
		}
	}

	return nil
}

func (cls *SqlCheckoutLineStore) BulkCreate(lines []*model.CheckoutLine) ([]*model.CheckoutLine, error) {
	err := cls.GetMaster().Create(lines).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert a checkout line")
	}
	return lines, nil
}

// CheckoutLinesByCheckoutWithPrefetch finds all checkout lines belong to given checkout
//
// and prefetch all related product variants, products
//
// this borrows the idea from Django's prefetch_related() method
func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutWithPrefetch(checkoutToken string) ([]*model.CheckoutLine, []*model.ProductVariant, []*model.Product, error) {
	selectFields := []string{
		model.CheckoutLineTableName + ".*",
		model.ProductVariantTableName + ".*",
		model.ProductTableName + ".*",
	}

	query, args, err := cls.
		GetQueryBuilder().
		Select(selectFields...).
		From(model.CheckoutLineTableName).
		InnerJoin(model.ProductVariantTableName + " ON CheckoutLines.VariantID = ProductVariants.Id").
		InnerJoin(model.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
		Where(squirrel.Eq{"CheckoutLines.CheckoutID": checkoutToken}).
		ToSql()

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "CheckoutLinesByCheckoutWithPrefetch_ToSql")
	}

	rows, err := cls.GetReplica().Raw(query, args...).Rows()
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to find checkout lines and prefetch related values, with checkoutToken=%s", checkoutToken)
	}
	defer rows.Close()

	var (
		checkoutLines   model.CheckoutLines
		productVariants model.ProductVariants
		products        model.Products
	)

	for rows.Next() {
		var (
			checkoutLine   model.CheckoutLine
			productVariant model.ProductVariant
			product        model.Product
			scanFields     = append(
				cls.ScanFields(&checkoutLine),
				append(
					cls.ProductVariant().ScanFields(&productVariant),
					cls.Product().ScanFields(&product)...,
				)...,
			)
		)
		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to scan a row")
		}

		checkoutLines = append(checkoutLines, &checkoutLine)
		productVariants = append(productVariants, &productVariant)
		products = append(products, &product)
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
		From(model.CheckoutLineTableName).
		InnerJoin(model.ProductVariantTableName + " ON (CheckoutLines.VariantID = ProductVariants.Id)").
		InnerJoin(model.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
		InnerJoin(model.ProductTypeTableName + " ON (ProductTypes.Id = Products.ProductTypeID)").
		Where(squirrel.Eq{"CheckoutLines.Id": checkoutLineIDs}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "TotalWeightForCheckoutLines_ToSql")
	}

	rows, err := cls.GetReplica().Raw(query, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find values")
	}
	defer rows.Close()

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

	return totalWeight, nil
}

// CheckoutLinesByOption finds and returns checkout lines filtered using given option
func (cls *SqlCheckoutLineStore) CheckoutLinesByOption(option *model.CheckoutLineFilterOption) ([]*model.CheckoutLine, error) {
	args, err := store.BuildSqlizer(option.Conditions, "CheckoutLinesByOption")
	if err != nil {
		return nil, err
	}
	var res []*model.CheckoutLine
	err = cls.GetReplica().Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find checkout lines by given options")
	}

	return res, nil
}
