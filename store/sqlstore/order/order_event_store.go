package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	return &SqlOrderEventStore{s}
}

func (oes *SqlOrderEventStore) Save(transaction *gorm.DB, orderEvent *model.OrderEvent) (*model.OrderEvent, error) {
	if transaction == nil {
		transaction = oes.GetMaster()
	}

	if err := transaction.Save(orderEvent).Error; err != nil {
		return nil, errors.Wrap(err, "failed to save order event")
	}

	return orderEvent, nil
}

func (oes *SqlOrderEventStore) Get(orderEventID string) (*model.OrderEvent, error) {
	var res model.OrderEvent
	err := oes.GetReplica().First(&res, "Id = ?", orderEventID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderEventTableName, orderEventID)
		}
		return nil, errors.Wrapf(err, "failed to find order event iwth id=%s", orderEventID)
	}

	return &res, nil
}

func (s *SqlOrderEventStore) FilterByOptions(options *model.OrderEventFilterOptions) (int64, []*model.OrderEvent, error) {
	query := s.
		GetQueryBuilder().
		Select(model.OrderEventTableName + ".*").
		Where(options.Conditions)

	var totalCount int64
	if options.CountTotal {
		countQuery, args, err := s.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOptions_CounTotal_ToSql")
		}

		err = s.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of order events by given options")
		}
	}

	// NOTE: we apply pagination conditions after count total (if required)
	options.GraphqlPaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}
	var res []*model.OrderEvent
	err = s.GetReplica().Raw(queryStr, args...).Scan(&res).Error
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find order events by given options")
	}

	return totalCount, res, nil
}
