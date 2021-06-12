package account

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlTermsOfServiceStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
}

func NewSqlTermsOfServiceStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.TermsOfServiceStore {
	s := &SqlTermsOfServiceStore{sqlStore, metrics}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.TermsOfService{}, "TermsOfServices").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("UserId").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Text").SetMaxSize(model.POST_MESSAGE_MAX_BYTES_V2)
	}

	return s
}

func (s *SqlTermsOfServiceStore) CreateIndexesIfNotExists() {
}

func (s *SqlTermsOfServiceStore) Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error) {
	if termsOfService.Id != "" {
		return nil, store.NewErrInvalidInput("TermsOfService", "Id", termsOfService.Id)
	}

	termsOfService.PreSave()

	if err := termsOfService.IsValid(); err != nil {
		return nil, err
	}

	if err := s.GetMaster().Insert(termsOfService); err != nil {
		return nil, errors.Wrapf(err, "could not save a new TermsOfService")
	}

	return termsOfService, nil
}

func (s *SqlTermsOfServiceStore) GetLatest(allowFromCache bool) (*model.TermsOfService, error) {
	var termsOfService *model.TermsOfService

	query := s.GetQueryBuilder().
		Select("*").
		From("TermsOfService").
		OrderBy("CreateAt DESC").
		Limit(uint64(1))

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build sql query to get latest TOS")
	}

	if err := s.GetReplica().SelectOne(&termsOfService, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("TermsOfService", "CreateAt=latest")
		}
		return nil, errors.Wrap(err, "could not find latest TermsOfService")
	}

	return termsOfService, nil
}

func (s *SqlTermsOfServiceStore) Get(id string, allowFromCache bool) (*model.TermsOfService, error) {
	obj, err := s.GetReplica().Get(model.TermsOfService{}, id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find TermsOfService with id=%s", id)
	}
	if obj == nil {
		return nil, store.NewErrNotFound("TermsOfService", id)
	}
	return obj.(*model.TermsOfService), nil
}
