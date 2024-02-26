package order

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlFulfillmentStore struct {
	store.Store
}

func NewSqlFulfillmentStore(sqlStore store.Store) store.FulfillmentStore {
	return &SqlFulfillmentStore{sqlStore}
}

// Upsert depends on given fulfillment's Id to decide update or insert it
func (fs *SqlFulfillmentStore) Upsert(transaction *gorm.DB, fulfillment *model.Fulfillment) (*model.Fulfillment, error) {
	if transaction == nil {
		transaction = fs.GetMaster()
	}

	var err error
	if fulfillment.Id == "" {
		err = transaction.Create(fulfillment).Error
	} else {
		// prevent update:
		fulfillment.CreateAt = 0
		fulfillment.OrderID = ""
		fulfillment.FulfillmentOrder = 0

		err = transaction.Model(fulfillment).Updates(fulfillment).Error
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert fulfillment")
	}
	return fulfillment, nil
}

// Get fidns and returns a fulfillment with given id
func (fs *SqlFulfillmentStore) Get(id string) (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	if err := fs.GetReplica().First(&ffm, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.FulfillmentTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment with id=%s", id)
	}

	return &ffm, nil
}

func (fs *SqlFulfillmentStore) commonQueryBuild(option *model.FulfillmentFilterOption) squirrel.SelectBuilder {
	// decide which fiedlds to select
	selectFields := []string{model.FulfillmentTableName + ".*"}
	if option.SelectRelatedOrder {
		selectFields = append(selectFields, model.OrderTableName+".*")
	}

	// build query:
	query := fs.GetQueryBuilder().
		Select(selectFields...).
		From(model.FulfillmentTableName).
		Where(option.Conditions)

	// parse option
	if option.FulfillmentLineID != nil {
		query = query.
			InnerJoin(model.FulfillmentLineTableName + " ON FulfillmentLines.FulfillmentID = Fulfillments.Id").
			Where(option.FulfillmentLineID)
	} else if option.HaveNoFulfillmentLines {
		query = query.
			LeftJoin(model.FulfillmentLineTableName + " ON FulfillmentLines.FulfillmentID = Fulfillments.Id").
			Where("FulfillmentLines.FulfillmentID IS NULL")
	}

	if option.SelectRelatedOrder {
		query = query.InnerJoin(model.OrderTableName + " ON (Orders.Id = Fulfillments.OrderID)")
	}
	if option.SelectForUpdate && option.Transaction != nil {
		query = query.Suffix("FOR UPDATE")
	}

	return query
}

// GetByOption returns 1 fulfillment, filtered by given option
func (fs *SqlFulfillmentStore) GetByOption(option *model.FulfillmentFilterOption) (*model.Fulfillment, error) {
	queryString, args, err := fs.commonQueryBuild(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByOption_ToSql")
	}

	var (
		fulfillment model.Fulfillment
		order       model.Order
		scanFields  = fs.ScanFields(&fulfillment)
	)
	if option.SelectRelatedOrder {
		scanFields = append(scanFields, fs.Order().ScanFields(&order)...)
	}

	runner := fs.GetMaster()
	if option.Transaction != nil {
		runner = option.Transaction
	}
	err = runner.Raw(queryString, args...).Row().Scan(scanFields)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.NewErrNotFound(model.FulfillmentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find fulfillment based on given option")
	}

	if option.SelectForUpdate {
		fulfillment.SetOrder(&order)
	}

	return &fulfillment, nil
}

// FilterByOption finds and returns a slice of fulfillments by given option
func (fs *SqlFulfillmentStore) FilterByOption(option *model.FulfillmentFilterOption) ([]*model.FulfillmentSlice, error) {
	runner := fs.GetMaster()
	if option.Transaction != nil {
		runner = option.Transaction
	}

	queryString, args, err := fs.commonQueryBuild(option).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	rows, err := runner.Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillments with given option")
	}
	defer rows.Close()

	var res model.Fulfillments

	for rows.Next() {
		var (
			fulfillment model.Fulfillment
			order       model.Order
			scanFields  = fs.ScanFields(&fulfillment)
		)
		if option.SelectRelatedOrder {
			scanFields = append(scanFields, fs.Order().ScanFields(&order)...)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan a row on fulfillment and related order")
		}

		if option.SelectRelatedOrder {
			fulfillment.SetOrder(&order)
		}
		res = append(res, &fulfillment)
	}

	return res, nil
}

// BulkDeleteFulfillments deletes given fulfillments
func (fs *SqlFulfillmentStore) BulkDeleteFulfillments(transaction *gorm.DB, fulfillments model.Fulfillments) error {
	if transaction == nil {
		transaction = fs.GetMaster()
	}

	err := transaction.Raw("DELETE * FROM "+model.FulfillmentTableName+" WHERE Id IN ?", fulfillments.IDs()).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillments")
	}

	return nil
}
