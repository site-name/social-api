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

	isSaving := orderLine.ID == ""
	if isSaving {
		model_helper.OrderLinePreSave(&orderLine)
	} else {
		model_helper.OrderLineCommonPre(&orderLine)
	}

	if err := model_helper.OrderLineIsValid(orderLine); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = orderLine.Insert(transaction, boil.Infer())
	} else {
		_, err = orderLine.Update(transaction, boil.Blacklist(
			model.OrderLineColumns.ID,
			model.OrderLineColumns.CreatedAt,
			model.OrderLineColumns.OrderID,
		))
	}

	if err != nil {
		return nil, err
	}

	return &orderLine, nil
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

func (ols *SqlOrderLineStore) FilterbyOption(option model_helper.OrderLineFilterOptions) (model.OrderLineSlice, error) {
	conds := option.Conditions
	if option.RelatedOrderConditions != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.Orders, model.OrderTableColumns.ID, model.OrderLineTableColumns.OrderID)),
			option.RelatedOrderConditions,
		)
	}

	for _, load := range option.Preload {
		conds = append(conds, qm.Load(load))
	}

	if option.VariantProductID != nil {
		conds = append(
			conds,
			qm.InnerJoin(fmt.Sprintf("%s ON %s = %s", model.TableNames.ProductVariants, model.ProductVariantTableColumns.ID, model.OrderLineTableColumns.ProductVariantID)),
			option.VariantProductID,
		)
	}

	return model.OrderLines(conds...).All(ols.GetReplica())
}
