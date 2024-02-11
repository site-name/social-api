package discount

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlOrderDiscountStore struct {
	store.Store
}

func NewSqlOrderDiscountStore(sqlStore store.Store) store.OrderDiscountStore {
	return &SqlOrderDiscountStore{sqlStore}
}

func (ods *SqlOrderDiscountStore) Upsert(transaction boil.ContextTransactor, orderDiscount model.OrderDiscount) (*model.OrderDiscount, error) {
	if transaction == nil {
		transaction = ods.GetMaster()
	}

	isSaving := false
	if orderDiscount.ID == "" {
		isSaving = true
		model_helper.OrderDiscountPreSave(&orderDiscount)
	} else {
		model_helper.OrderDiscountPreUpdate(&orderDiscount)
	}

	if err := model_helper.OrderDiscountIsValid(orderDiscount); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = orderDiscount.Insert(transaction, boil.Infer())
	} else {
		_, err = orderDiscount.Update(transaction, boil.Infer())
	}

	if err != nil {
		return nil, err
	}

	return &orderDiscount, nil
}

func (ods *SqlOrderDiscountStore) Get(orderDiscountID string) (*model.OrderDiscount, error) {
	orderDiscount, err := model.FindOrderDiscount(ods.GetReplica(), orderDiscountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.OrderDiscounts, orderDiscountID)
		}
		return nil, err
	}
	return orderDiscount, nil
}

func (ods *SqlOrderDiscountStore) FilterbyOption(option model_helper.OrderDiscountFilterOption) (model.OrderDiscountSlice, error) {
	return model.OrderDiscounts(option.Conditions...).All(ods.GetReplica())
}

func (ods *SqlOrderDiscountStore) BulkDelete(orderDiscountIDs []string) error {
	_, err := model.OrderDiscounts(model.OrderDiscountWhere.ID.IN(orderDiscountIDs)).DeleteAll(ods.GetMaster())
	if err != nil {
		return err
	}
	return nil
}
