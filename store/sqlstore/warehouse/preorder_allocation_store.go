package warehouse

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlPreorderAllocationStore struct {
	store.Store
}

func NewSqlPreorderAllocationStore(s store.Store) store.PreorderAllocationStore {
	return &SqlPreorderAllocationStore{s}
}

func (ws *SqlPreorderAllocationStore) ScanFields(preorderAllocation *model.PreorderAllocation) []interface{} {
	return []interface{}{
		&preorderAllocation.Id,
		&preorderAllocation.OrderLineID,
		&preorderAllocation.Quantity,
		&preorderAllocation.ProductVariantChannelListingID,
	}
}

// BulkCreate bulk inserts given preorderAllocations and returns them
func (ws *SqlPreorderAllocationStore) BulkCreate(transaction *gorm.DB, preorderAllocations []*model.PreorderAllocation) ([]*model.PreorderAllocation, error) {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	for _, allocation := range preorderAllocations {
		err := transaction.Save(allocation).Error
		if err != nil {
			if ws.IsUniqueConstraintError(err, []string{"OrderLineID", "ProductVariantChannelListingID", "orderlineid_productvariantchannellistingid_key"}) {
				return nil, store.NewErrInvalidInput(model.PreOrderAllocationTableName, "OrderLineID/ProductVariantChannelListingID", "duplicate")
			}
			return nil, errors.Wrap(err, "failed to upsert preorder allocation")
		}
	}

	return preorderAllocations, nil
}

// FilterByOption finds and returns a list of preorder allocations filtered using given options
func (ws *SqlPreorderAllocationStore) FilterByOption(options *model.PreorderAllocationFilterOption) ([]*model.PreorderAllocation, error) {
	selectFields := []string{model.PreOrderAllocationTableName + ".*"}
	if options.SelectRelated_OrderLine {
		selectFields = append(selectFields, model.OrderLineTableName+".*")
	}
	if options.SelectRelated_OrderLine_Order && options.SelectRelated_OrderLine {
		selectFields = append(selectFields, model.OrderTableName+".*")
	}

	query := ws.GetQueryBuilder().Select(selectFields...).From(model.PreOrderAllocationTableName).Where(options.Conditions)

	if options.SelectRelated_OrderLine {
		query = query.InnerJoin(model.OrderLineTableName + " ON PreorderAllocations.OrderLineID = Orderlines.Id")

		if options.SelectRelated_OrderLine_Order {
			query = query.InnerJoin(model.OrderTableName + " ON Orderlines.OrderID = Orders.Id")
		}
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := ws.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find preorder allocations with given options")
	}
	defer rows.Close()

	var res model.PreorderAllocations

	for rows.Next() {
		var (
			preorderAllocation model.PreorderAllocation
			orderLine          model.OrderLine
			orDer              model.Order
			scanFields         = ws.ScanFields(&preorderAllocation)
		)
		if options.SelectRelated_OrderLine {
			scanFields = append(scanFields, ws.OrderLine().ScanFields(&orderLine)...)
		}
		if options.SelectRelated_OrderLine && options.SelectRelated_OrderLine_Order {
			scanFields = append(scanFields, ws.Order().ScanFields(&orDer)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row of preorder allocation")
		}

		// join data.
		if options.SelectRelated_OrderLine {
			preorderAllocation.SetOrderLine(&orderLine)

			if options.SelectRelated_OrderLine_Order {
				orderLine.Order = &orDer
			}
		}

		res = append(res, &preorderAllocation)
	}

	return res, nil
}

// Delete deletes preorder-allocations by given ids
func (ws *SqlPreorderAllocationStore) Delete(transaction *gorm.DB, preorderAllocationIDs ...string) error {
	if transaction == nil {
		transaction = ws.GetMaster()
	}

	err := transaction.Raw("DELETE FROM "+model.PreOrderAllocationTableName+" WHERE Id IN ?", preorderAllocationIDs).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete preorder-allocations with given ids")
	}

	return nil
}
