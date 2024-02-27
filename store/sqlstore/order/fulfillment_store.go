package order

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlFulfillmentStore struct {
	store.Store
}

func NewSqlFulfillmentStore(sqlStore store.Store) store.FulfillmentStore {
	return &SqlFulfillmentStore{sqlStore}
}

func (fs *SqlFulfillmentStore) Upsert(transaction boil.ContextTransactor, fulfillment model.Fulfillment) (*model.Fulfillment, error) {
	if transaction == nil {
		transaction = fs.GetMaster()
	}

	isSaving := fulfillment.ID == ""
	if isSaving {
		model_helper.FulfillmentPreSave(&fulfillment)
	} else {
		model_helper.FulfillmentCommonPre(&fulfillment)
	}

	if err := model_helper.FulfillmentIsValid(fulfillment); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = fulfillment.Insert(transaction, boil.Infer())
	} else {
		_, err = fulfillment.Update(transaction, boil.Blacklist(
			model.FulfillmentColumns.CreatedAt,
			model.FulfillmentColumns.OrderID,
		))
	}
	if err != nil {
		return nil, err
	}

	return &fulfillment, nil
}

func (fs *SqlFulfillmentStore) Get(id string) (*model.Fulfillment, error) {
	record, err := model.FindFulfillment(fs.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Fulfillments, id)
		}
		return nil, err
	}

	return record, nil
}

func (fs *SqlFulfillmentStore) commonQueryBuild(option model_helper.FulfillmentFilterOption) []qm.QueryMod {
	conds := option.Conditions
	if option.FulfillmentLineID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.FulfillmentLines, model.FulfillmentTableColumns.ID, model.FulfillmentLineTableColumns.FulfillmentID)),
		)
	} else if option.HaveNoFulfillmentLines {
		conds = append(
			conds,
			qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.FulfillmentLines, model.FulfillmentTableColumns.ID, model.FulfillmentLineTableColumns.FulfillmentID)),
			qm.Where(fmt.Sprintf("%s IS NULL", model.FulfillmentLineTableColumns.FulfillmentID)),
		)
	}

	for _, load := range option.Preload {
		conds = append(conds, qm.Load(load))
	}

	return conds
}

func (fs *SqlFulfillmentStore) FilterByOption(option model_helper.FulfillmentFilterOption) (model.FulfillmentSlice, error) {
	conds := fs.commonQueryBuild(option)
	return model.Fulfillments(conds...).All(fs.GetReplica())
}

func (fs *SqlFulfillmentStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = fs.GetMaster()
	}

	_, err := model.Fulfillments(model.FulfillmentWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
