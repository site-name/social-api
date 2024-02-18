package attribute

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlCustomProductAttributeStore struct {
	store.Store
}

func NewSqlCustomProductAttributeStore(s store.Store) store.CustomProductAttributeStore {
	return &SqlCustomProductAttributeStore{s}
}

func (cpas *SqlCustomProductAttributeStore) Upsert(tx boil.ContextTransactor, record model.CustomProductAttribute) (*model.CustomProductAttribute, error) {
	if tx == nil {
		tx = cpas.GetMaster()
	}

	isSaving := record.ID == ""
	if isSaving {
		model_helper.CustomProductAttributePreSave(&record)
	} else {
		model_helper.CustomProductAttributeCommonPre(&record)
	}

	if err := model_helper.CustomProductAttributeIsValid(record); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = record.Insert(tx, boil.Infer())
	} else {
		_, err = record.Update(tx, boil.Blacklist(model.CustomProductAttributeColumns.ProductID))
	}

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (cpas *SqlCustomProductAttributeStore) FilterByOptions(options model_helper.CustomProductAttributeFilterOptions) (model.CustomProductAttributeSlice, error) {
	return model.CustomProductAttributes(options.Conditions...).All(cpas.GetReplica())
}

func (cpas *SqlCustomProductAttributeStore) Delete(tx boil.ContextTransactor, ids []string) (int64, error) {
	if tx == nil {
		tx = cpas.GetMaster()
	}

	return model.CustomProductAttributes(model.CustomProductAttributeWhere.ID.IN(ids)).DeleteAll(tx)
}
