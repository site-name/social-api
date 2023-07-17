package order

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
	"gorm.io/gorm"
)

type SqlOrderEventStore struct {
	store.Store
}

func NewSqlOrderEventStore(s store.Store) store.OrderEventStore {
	return &SqlOrderEventStore{s}
}

func (s *SqlOrderEventStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CreateAt",
		"Type",
		"OrderID",
		"Parameters",
		"UserID",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (oes *SqlOrderEventStore) Save(transaction store_iface.SqlxExecutor, orderEvent *model.OrderEvent) (*model.OrderEvent, error) {
	var executor store_iface.SqlxExecutor = oes.GetMasterX()
	if transaction != nil {
		executor = transaction
	}

	orderEvent.PreSave()
	if err := orderEvent.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + model.OrderEventTableName + "(" + oes.ModelFields("").Join(",") + ") VALUES (" + oes.ModelFields(":").Join(",") + ")"
	if _, err := executor.NamedExec(query, orderEvent); err != nil {
		return nil, errors.Wrapf(err, "failed to save order event with id=%s", orderEvent.Id)
	}

	return orderEvent, nil
}

func (oes *SqlOrderEventStore) Get(orderEventID string) (*model.OrderEvent, error) {
	var res model.OrderEvent
	err := oes.GetReplicaX().Get(&res, "SELECT * FROM "+model.OrderEventTableName+" WHERE Id = ?", orderEventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.OrderEventTableName, orderEventID)
		}
		return nil, errors.Wrapf(err, "failed to find order event iwth id=%s", orderEventID)
	}

	return &res, nil
}

func (s *SqlOrderEventStore) FilterByOptions(options *model.OrderEventFilterOptions) ([]*model.OrderEvent, error) {
	query := s.GetQueryBuilder().
		Select(s.ModelFields(model.OrderEventTableName + ".")...).
		From(model.OrderEventTableName)

	// parse options
	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.OrderID != nil {
		query = query.Where(options.OrderID)
	}
	if options.Type != nil {
		query = query.Where(options.Type)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.OrderEvent
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find order events by given options")
	}

	return res, nil
}
