package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlPreorderAllocationStore struct {
	store.Store
}

func NewSqlPreorderAllocationStore(s store.Store) store.PreorderAllocationStore {
	return &SqlPreorderAllocationStore{s}
}

func (ws *SqlPreorderAllocationStore) Upsert(transaction boil.ContextTransactor, preorderAllocations model.PreorderAllocationSlice) (model.PreorderAllocationSlice, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	for _, allocation := range preorderAllocations {
		if allocation == nil {
			continue
		}

		isSaving := allocation.ID == ""
		if isSaving {
			model_helper.PreorderAllocationPreSave(allocation)
		}

		if err := model_helper.PreorderAllocationIsValid(*allocation); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = allocation.Insert(transaction, boil.Infer())
		} else {
			_, err = allocation.Update(transaction, boil.Infer())
		}

		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"preorder_allocations_order_line_id_product_variant_channel_listing_id_key"}) {
				return nil, store.NewErrInvalidInput(model.TableNames.PreorderAllocations, "OrderLineID/ProductVariantChannelListingID", "duplicate")
			}
			return nil, err
		}
	}

	return preorderAllocations, nil
}

func (ws *SqlPreorderAllocationStore) FilterByOption(options model_helper.PreorderAllocationFilterOption) (model.PreorderAllocationSlice, error) {
	conds := options.Conditions
	for _, load := range options.Preloads {
		conds = append(conds, qm.Load(load))
	}

	return model.PreorderAllocations(conds...).All(ws.GetReplica())
}

func (ws *SqlPreorderAllocationStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	_, err := model.PreorderAllocations(model.PreorderAllocationWhere.ID.IN(ids)).DeleteAll(transaction)
	return err
}
