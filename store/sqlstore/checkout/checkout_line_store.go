package checkout

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlCheckoutLineStore struct {
	store.Store
}

func NewSqlCheckoutLineStore(sqlStore store.Store) store.CheckoutLineStore {
	return &SqlCheckoutLineStore{sqlStore}
}

func (cls *SqlCheckoutLineStore) Upsert(checkoutLines model.CheckoutLineSlice) (model.CheckoutLineSlice, error) {
	for _, line := range checkoutLines {
		if line == nil {
			continue
		}

		isSaving := false
		if line.ID == "" {
			isSaving = true
			model_helper.CheckoutLinePreSave(line)
		}
		if err := model_helper.CheckoutLineIsValid(*line); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = line.Insert(cls.GetMaster(), boil.Infer())
		} else {
			_, err = line.Update(cls.GetMaster(), boil.Blacklist(model.CheckoutLineColumns.CreatedAt, model.CheckoutLineColumns.CheckoutID))
		}

		if err != nil {
			return nil, err
		}
	}

	return checkoutLines, nil
}

func (cls *SqlCheckoutLineStore) Get(id string) (*model.CheckoutLine, error) {
	checkoutLine, err := model.FindCheckoutLine(cls.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.CheckoutLines, id)
		}
		return nil, err
	}
	return checkoutLine, nil
}

func (cls *SqlCheckoutLineStore) DeleteLines(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = cls.GetMaster()
	}

	_, err := model.CheckoutLines(model.CheckoutLineWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}

// CheckoutLinesByCheckoutWithPrefetch finds all checkout lines belong to given checkout
//
// and prefetch all related product variants, products
//
// this borrows the idea from Django's prefetch_related() method
// func (cls *SqlCheckoutLineStore) CheckoutLinesByCheckoutWithPrefetch(checkoutToken string) ([]*model.CheckoutLine, []*model.ProductVariant, []*model.Product, error) {
// 	selectFields := []string{
// 		model.CheckoutLineTableName + ".*",
// 		model.ProductVariantTableName + ".*",
// 		model.ProductTableName + ".*",
// 	}

// 	query, args, err := cls.
// 		GetQueryBuilder().
// 		Select(selectFields...).
// 		From(model.CheckoutLineTableName).
// 		InnerJoin(model.ProductVariantTableName + " ON CheckoutLines.VariantID = ProductVariants.Id").
// 		InnerJoin(model.ProductTableName + " ON ProductVariants.ProductID = Products.Id").
// 		Where(squirrel.Eq{"CheckoutLines.CheckoutID": checkoutToken}).
// 		ToSql()

// 	if err != nil {
// 		return nil, nil, nil, errors.Wrap(err, "CheckoutLinesByCheckoutWithPrefetch_ToSql")
// 	}

// 	rows, err := cls.GetReplica().Raw(query, args...).Rows()
// 	if err != nil {
// 		return nil, nil, nil, errors.Wrapf(err, "failed to find checkout lines and prefetch related values, with checkoutToken=%s", checkoutToken)
// 	}
// 	defer rows.Close()

// 	var (
// 		checkoutLines   model.CheckoutLines
// 		productVariants model.ProductVariants
// 		products        model.Products
// 	)

// 	for rows.Next() {
// 		var (
// 			checkoutLine   model.CheckoutLine
// 			productVariant model.ProductVariant
// 			product        model.Product
// 			scanFields     = append(
// 				cls.ScanFields(&checkoutLine),
// 				append(
// 					cls.ProductVariant().ScanFields(&productVariant),
// 					cls.Product().ScanFields(&product)...,
// 				)...,
// 			)
// 		)
// 		err = rows.Scan(scanFields...)
// 		if err != nil {
// 			return nil, nil, nil, errors.Wrap(err, "failed to scan a row")
// 		}

// 		checkoutLines = append(checkoutLines, &checkoutLine)
// 		productVariants = append(productVariants, &productVariant)
// 		products = append(products, &product)
// 	}

// 	return checkoutLines, productVariants, products, nil
// }

// TotalWeightForCheckoutLines calculate total weight for given checkout lines
func (cls *SqlCheckoutLineStore) TotalWeightForCheckoutLines(checkoutLineIDs []string) (*measurement.Weight, error) {
	// query, args, err := cls.
	// 	GetQueryBuilder().
	// 	Select(
	// 		"CheckoutLines.Quantity",
	// 		"ProductVariants.Weight",
	// 		"ProductVariants.WeightUnit",
	// 		"Products.Weight",
	// 		"Products.WeightUnit",
	// 		"ProductTypes.Weight",
	// 		"ProductTypes.WeightUnit",
	// 	).
	// 	From(model.CheckoutLineTableName).
	// 	InnerJoin(model.ProductVariantTableName + " ON (CheckoutLines.VariantID = ProductVariants.Id)").
	// 	InnerJoin(model.ProductTableName + " ON (Products.Id = ProductVariants.ProductID)").
	// 	InnerJoin(model.ProductTypeTableName + " ON (ProductTypes.Id = Products.ProductTypeID)").
	// 	Where(squirrel.Eq{"CheckoutLines.Id": checkoutLineIDs}).
	// 	ToSql()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "TotalWeightForCheckoutLines_ToSql")
	// }

	// rows, err := cls.GetReplica().Raw(query, args...).Rows()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to find values")
	// }
	// defer rows.Close()

	// var (
	// 	totalWeight  = measurement.ZeroWeight
	// 	lineQuantity int
	// 	weight       measurement.Weight

	// 	variantWeightAmount *float32
	// 	variantWeightUnit   measurement.WeightUnit

	// 	productWeightAmount *float32
	// 	productWeightUnit   measurement.WeightUnit

	// 	productTypeWeightAmount *float32
	// 	productTypeWeightUnit   measurement.WeightUnit
	// )

	// for rows.Next() {
	// 	err = rows.Scan(
	// 		&lineQuantity,
	// 		&variantWeightAmount,
	// 		&variantWeightUnit,
	// 		&productWeightAmount,
	// 		&productWeightUnit,
	// 		&productTypeWeightAmount,
	// 		&productTypeWeightUnit,
	// 	)
	// 	if err != nil {
	// 		// return immediately if an error occured
	// 		return nil, errors.Wrap(err, "failed to scan a row")
	// 	}

	// 	if variantWeightAmount != nil {
	// 		weight = measurement.Weight{Amount: *variantWeightAmount, Unit: variantWeightUnit}
	// 	} else if productWeightAmount != nil {
	// 		weight = measurement.Weight{Amount: *variantWeightAmount, Unit: productWeightUnit}
	// 	} else if productTypeWeightAmount != nil {
	// 		weight = measurement.Weight{Amount: *productTypeWeightAmount, Unit: productTypeWeightUnit}
	// 	}

	// 	if weight.Amount != 0 && weight.Unit != "" {
	// 		totalWeight, err = totalWeight.Add(weight.Mul(float32(lineQuantity)))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	// return totalWeight, nil

	var queryResult []*struct {
		LineQuantity        int
		VariantWeightAmount *float64
		VariantWeightUnit   measurement.WeightUnit
		ProductWeightAmount *float64
		ProductWeightUnit   measurement.WeightUnit
	}

	err := model.CheckoutLines(
		qm.Select(
			model.CheckoutLineTableColumns.Quantity,
			model.ProductVariantTableColumns.Weight,
			model.ProductVariantTableColumns.WeightUnit,
			model.ProductTableColumns.Weight,
			model.ProductTableColumns.WeightUnit,
		),
		qm.InnerJoin("%s ON %s = %s", model.TableNames.ProductVariants, model.CheckoutLineTableColumns.VariantID, model.ProductVariantTableColumns.ID),
		qm.InnerJoin("%s ON %s = %s", model.TableNames.Products, model.ProductVariantTableColumns.ProductID, model.ProductTableColumns.ID),
		model.CheckoutLineWhere.ID.IN(checkoutLineIDs),
	).Bind(nil, cls.GetReplica(), &queryResult)

	if err != nil {
		return nil, err
	}

	var totalWeight = measurement.ZeroWeight

	for _, result := range queryResult {
		var weight measurement.Weight

		if result.VariantWeightAmount != nil {
			weight = measurement.Weight{Amount: *result.VariantWeightAmount, Unit: result.VariantWeightUnit}
		} else if result.ProductWeightAmount != nil {
			weight = measurement.Weight{Amount: *result.ProductWeightAmount, Unit: result.ProductWeightUnit}
		}

		if weight.Amount != 0 && weight.Unit != "" {
			addedWeight, err := totalWeight.Add(weight.Mul(result.LineQuantity))
			if err != nil {
				return nil, err
			}
			totalWeight = *addedWeight
		}
	}

	return &totalWeight, nil
}

// CheckoutLinesByOption finds and returns checkout lines filtered using given option
func (cls *SqlCheckoutLineStore) CheckoutLinesByOption(option model_helper.CheckoutLineFilterOptions) (model.CheckoutLineSlice, error) {
	return model.CheckoutLines(option.Conditions...).All(cls.GetReplica())
}
