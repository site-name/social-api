package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlFulfillmentLineStore struct {
	store.Store
}

func NewSqlFulfillmentLineStore(s store.Store) store.FulfillmentLineStore {
	return &SqlFulfillmentLineStore{s}
}

func (fls *SqlFulfillmentLineStore) Save(ffml *model.FulfillmentLine) (*model.FulfillmentLine, error) {
	if err := fls.GetMaster().Create(ffml).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save fulfillment line")
	}
	return ffml, nil
}

func (fls *SqlFulfillmentLineStore) Get(id string) (*model.FulfillmentLine, error) {
	var res model.FulfillmentLine
	if err := fls.GetReplica().First(&res, "Id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.FulfillmentLineTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find fulfillment line with id=%s", id)
	}
	return &res, nil
}

// BulkUpsert upsert given fulfillment lines
func (fls *SqlFulfillmentLineStore) BulkUpsert(transaction *gorm.DB, fulfillmentLines []*model.FulfillmentLine) ([]*model.FulfillmentLine, error) {
	if transaction == nil {
		transaction = fls.GetMaster()
	}

	var err error
	for _, line := range fulfillmentLines {
		if line.Id == "" {
			err = transaction.Create(line).Error
		} else {
			line.FulfillmentID = "" // prevent update
			err = transaction.Model(line).Updates(line).Error
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert fulfillment line")
	}

	return fulfillmentLines, nil
}

// FilterbyOption finds and returns a list of fulfillment lines by given option
func (fls *SqlFulfillmentLineStore) FilterbyOption(option *model.FulfillmentLineFilterOption) ([]*model.FulfillmentLine, error) {
	query := fls.GetQueryBuilder().
		Select(model.FulfillmentLineTableName + ".*").
		From(model.FulfillmentLineTableName).
		Where(option.Conditions)

	// this variable helps preventing the query from joining `Fulfillments` table multiple times.
	// var joinedFulfillmentTable bool

	if option.FulfillmentOrderID != nil ||
		option.FulfillmentStatus != nil {
		query = query.InnerJoin(model.FulfillmentTableName + " ON (FulfillmentLines.FulfillmentID = Fulfillments.Id)")

		query = query.
			Where(option.FulfillmentOrderID).
			Where(option.FulfillmentStatus)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	db := fls.GetReplica()
	if len(option.Preloads) > 0 {
		for _, preload := range option.Preloads {
			db = db.Preload(preload)
		}
	}

	var fulfillmentLines model.FulfillmentLines
	err = db.Raw(queryString, args...).Scan(&fulfillmentLines).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find fulfillment lines by given options")
	}

	return fulfillmentLines, nil
}

// DeleteFulfillmentLinesByOption filters fulfillment lines by given option, then deletes them
func (fls *SqlFulfillmentLineStore) DeleteFulfillmentLinesByOption(transaction *gorm.DB, option *model.FulfillmentLineFilterOption) error {
	if transaction == nil {
		transaction = fls.GetMaster()
	}

	args, err := store.BuildSqlizer(option.Conditions, "FulfillmentLine_Delete")
	if err != nil {
		return err
	}

	err = transaction.Delete(&model.FulfillmentLine{}, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete fulfillment lines by given option")
	}

	return nil
}
