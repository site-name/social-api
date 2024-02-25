package order

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderLineStore struct {
	store.Store
}

func NewSqlOrderLineStore(sqlStore store.Store) store.OrderLineStore {
	return &SqlOrderLineStore{sqlStore}
}

// Upsert depends on given orderLine's Id to decide to update or save it
func (ols *SqlOrderLineStore) Upsert(transaction *gorm.DB, orderLine *model.OrderLine) (*model.OrderLine, error) {
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

// BulkUpsert performs upsert multiple order lines in once
func (ols *SqlOrderLineStore) BulkUpsert(transaction *gorm.DB, orderLines []*model.OrderLine) ([]*model.OrderLine, error) {
	for _, orderLine := range orderLines {
		_, err := ols.Upsert(transaction, orderLine)
		if err != nil {
			return nil, err
		}
	}

	return orderLines, nil
}

func (ols *SqlOrderLineStore) Get(id string) (*model.OrderLine, error) {
	var odl model.OrderLine
	err := ols.GetReplica().First(&odl, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find order line with id=%s", id)
	}

	return &odl, nil
}

// BulkDelete delete all given order lines. NOTE: validate given ids are valid uuids before calling me
func (ols *SqlOrderLineStore) BulkDelete(tx *gorm.DB, orderLineIDs []string) error {
	if tx == nil {
		tx = ols.GetMaster()
	}
	err := tx.Raw("DELETE FROM "+model.OrderLineTableName+" WHERE Id IN ?", orderLineIDs).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete order lines with given ids")
	}

	return nil
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
