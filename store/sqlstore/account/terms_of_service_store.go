package account

import (
	"database/sql"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlTermsOfServiceStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
}

func NewSqlTermsOfServiceStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.TermsOfServiceStore {
	return &SqlTermsOfServiceStore{sqlStore, metrics}
}

func (s *SqlTermsOfServiceStore) Save(termsOfService model.TermsOfService) (*model.TermsOfService, error) {
	model_helper.TermsOfServicePreSave(&termsOfService)
	if err := model_helper.TermsOfServiceIsValid(termsOfService); err != nil {
		return nil, err
	}

	err := termsOfService.Insert(s.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}
	return &termsOfService, nil
}

func (s *SqlTermsOfServiceStore) GetLatest(_ bool) (*model.TermsOfService, error) {
	term, err := model.TermsOfServices(qm.OrderBy(model.TermsOfServiceColumns.CreatedAt + " DESC")).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.TermsOfServices, "latest")
		}
		return nil, err
	}

	return term, nil
}

func (s *SqlTermsOfServiceStore) Get(id string, _ bool) (*model.TermsOfService, error) {
	term, err := model.FindTermsOfService(s.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.TermsOfServices, id)
		}
		return nil, err
	}

	return term, nil
}
