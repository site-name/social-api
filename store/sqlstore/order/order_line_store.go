package order

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	return &SqlOrderLineStore{sqlStore}
}

func (ols *SqlOrderLineStore) Upsert(transaction boil.ContextTransactor, orderLine model.OrderLine) (*model.OrderLine, error) {
	if transaction == nil {
		transaction = ols.GetMaster()
	}

	var err error

	if orderLine.Id == "" {
		err = transaction.Create(orderLine).Error
	} else {
		orderLine.CreateAt = 0
		err = transaction.Model(orderLine).Updates(orderLine).Error
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert order line")
	}

	return orderLine, nil
}

func (ols *SqlOrderLineStore) Get(id string) (*model.OrderLine, error) {
	orderLine, err := model.FindOrderLine(ols.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.OrderLines, id)
		}
		return nil, err
	}

	return orderLine, nil

}

func (ols *SqlOrderLineStore) Delete(tx boil.ContextTransactor, orderLineIDs []string) error {
	if tx == nil {
		tx = ols.GetMaster()
	}

	_, err := model.OrderLines(model.OrderLineWhere.ID.IN(orderLineIDs)).DeleteAll(tx)
	return err
}

func (ols *SqlOrderLineStore) FilterbyOption(option *model.OrderLineFilterOption) (model.OrderLineSlice, error) {
	query := ols.GetReplica()
	if len(option.Preload) > 0 {
		for _, rel := range option.Preload {
			query = query.Preload(rel)
		}
	}

	conditions := squirrel.And{
		option.Conditions,
	}

	if option.RelatedOrderConditions != nil {
		query = query.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.OrderTableName,         // 1
				model.OrderLineTableName,     // 2
				model.OrderColumnId,          // 3
				model.OrderLineColumnOrderID, // 4
			),
		)

		conditions = append(conditions, option.RelatedOrderConditions)
	}

	if option.VariantDigitalContentID != nil || option.VariantProductID != nil {
		query = query.Joins(
			fmt.Sprintf(
				"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
				model.ProductVariantTableName,         // 1
				model.OrderLineTableName,              // 2
				model.ProductVariantColumnId,          // 3
				model.OrderLineColumnProductVariantID, // 4
			),
		)
		conditions = append(conditions, option.VariantProductID)

		if option.VariantDigitalContentID != nil {
			query = query.Joins(
				fmt.Sprintf(
					"INNER JOIN %[1]s ON %[1]s.%[3]s = %[2]s.%[4]s",
					model.DigitalContentTableName,              // 1
					model.ProductVariantTableName,              // 2
					model.DigitalContentColumnProductVariantID, // 3
					model.ProductVariantColumnId,               // 4
				),
			)
			conditions = append(conditions, option.VariantDigitalContentID)
		}
	}

	conds, args, err := conditions.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "please provide valid lookup conditions")
	}

	var orderLines model.OrderLineSlice
	err = query.Find(&orderLines, []any{conds, args}...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order lines by given options")
	}

	return orderLines, nil
}
