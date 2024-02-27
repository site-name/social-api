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

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	return &SqlFulfillmentLineStore{s}
}

func (fls *SqlFulfillmentLineStore) Upsert(ffml model.FulfillmentLine) (*model.FulfillmentLine, error) {
	isSaving := ffml.ID == ""
	if isSaving {
		model_helper.FulfillmentLinePreSave(&ffml)
	}

	if err := model_helper.FulfillmentLineIsValid(ffml); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = ffml.Insert(fls.GetMaster(), boil.Infer())
	} else {
		_, err = ffml.Update(fls.GetMaster(), boil.Infer())
	}

	if err != nil {
		return nil, err
	}

	return &ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*model.FulfillmentLine, error) {
	line, err := model.FindFulfillmentLine(fls.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.FulfillmentLines, id)
		}
		return nil, err
	}

	return line, nil
}

func (fls *SqlFulfillmentLineStore) FilterByOptions(option model_helper.FulfillmentLineFilterOption) (model.FulfillmentLineSlice, error) {
	conds := option.Conditions
	if option.RelatedFulfillmentConds != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Fulfillments, model.FulfillmentTableColumns.ID, model.FulfillmentLineTableColumns.FulfillmentID)),
			option.RelatedFulfillmentConds,
		)
	}

	for _, load := range option.Preload {
		conds = append(conds, qm.Load(load))
	}

	return model.FulfillmentLines(conds...).All(fls.GetReplica())
}

func (fls *SqlFulfillmentLineStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = fls.GetMaster()
	}

	_, err := model.FulfillmentLines(model.FulfillmentLineWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
