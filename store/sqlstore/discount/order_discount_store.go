package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderDiscountStore struct {
	store.Store
}

func NewSqlOrderDiscountStore(sqlStore store.Store) store.OrderDiscountStore {
	return &SqlOrderDiscountStore{sqlStore}
}

// Upsert depends on given order discount's Id property to decide to update/insert it
func (ods *SqlOrderDiscountStore) Upsert(transaction *gorm.DB, orderDiscount *model.OrderDiscount) (*model.OrderDiscount, error) {
	if transaction == nil {
		transaction = ods.GetMaster()
	}

	err := transaction.Save(orderDiscount).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert given order discount")
	}
	return orderDiscount, nil
}

// Get finds and returns an order discount with given id
func (ods *SqlOrderDiscountStore) Get(orderDiscountID string) (*model.OrderDiscount, error) {
	var res model.OrderDiscount

	err := ods.GetReplica().First(&res, "Id = ?", orderDiscountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderDiscountTableName, orderDiscountID)
		}
		return nil, errors.Wrapf(err, "failed to save order discount with id=%s", orderDiscountID)
	}

	return &res, nil
}

// FilterbyOption filters order discounts that satisfy given option, then returns them
func (ods *SqlOrderDiscountStore) FilterbyOption(option *model.OrderDiscountFilterOption) ([]*model.OrderDiscount, error) {
	db := ods.GetReplica()
	if option.PreloadOrder {
		db = db.Preload("Order")
	}

	args, err := store.BuildSqlizer(option.Conditions, "FilterByOptions")
	if err != nil {
		return nil, err
	}

	var res []*model.OrderDiscount
	err = db.Find(&res, args...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order discounts with given option")
	}

	return res, nil
}

// BulkDelete perform bulk delete all given order discount ids
func (ods *SqlOrderDiscountStore) BulkDelete(orderDiscountIDs []string) error {
	err := ods.GetMaster().Table(model.OrderDiscountTableName).Delete("Id IN ?", orderDiscountIDs).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete order discounts")
	}
	return nil
}
