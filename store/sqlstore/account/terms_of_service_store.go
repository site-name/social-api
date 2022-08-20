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

func (s *SqlTermsOfServiceStore) ModelFields(prefix string) model.StringArray {
	res := model.StringArray{
		"Id",
		"CreateAt",
		"UserID",
		"Text",
	}

	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func NewSqlTermsOfServiceStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.TermsOfServiceStore {
	return &SqlTermsOfServiceStore{sqlStore, metrics}
}

func (s *SqlTermsOfServiceStore) Save(termsOfService *model.TermsOfService) (*model.TermsOfService, error) {
	termsOfService.PreSave()

	if err := termsOfService.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.TermsOfServiceTableName + " (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
	if _, err := s.GetMasterX().NamedExec(query, termsOfService); err != nil {
		return nil, errors.Wrapf(err, "failed to save save a new TermsOfService")
	}

	return termsOfService, nil
}

func (s *SqlTermsOfServiceStore) GetLatest(allowFromCache bool) (*model.TermsOfService, error) {
	var termsOfService model.TermsOfService

	query := s.GetQueryBuilder().
		Select("*").
		From("TermsOfService").
		OrderBy("CreateAt DESC").
		Limit(uint64(1))

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build sql query to get latest TOS")
	}

	if err := s.GetReplicaX().Get(&termsOfService, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("TermsOfService", "CreateAt=latest")
		}
		return nil, errors.Wrap(err, "could not find latest TermsOfService")
	}

	return &termsOfService, nil
}

func (s *SqlTermsOfServiceStore) Get(id string, allowFromCache bool) (*model.TermsOfService, error) {
	var res model.TermsOfService

	err := s.GetReplicaX().Get(&res, "SELECT * FROM "+store.TermsOfServiceTableName+" WHERE Id = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find TermsOfService with id=%s", id)
	}
	return &res, nil
}
