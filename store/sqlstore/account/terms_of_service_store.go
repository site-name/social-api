package account

import (
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlTermsOfServiceStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
}

func (s *SqlTermsOfServiceStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
	err := s.GetMaster().Create(termsOfService).Error
	if err != nil {
		return nil, err
	}
	return termsOfService, nil
}

func (s *SqlTermsOfServiceStore) GetLatest(allowFromCache bool) (*model.TermsOfService, error) {
	var termsOfService model.TermsOfService

	err := s.GetReplica().Order("CreateAt DESC").First(&termsOfService).Error
	if err != nil {
		return nil, err
	}
	return &termsOfService, nil
}

func (s *SqlTermsOfServiceStore) Get(id string, allowFromCache bool) (*model.TermsOfService, error) {
	var res model.TermsOfService
	err := s.GetReplica().First(&res, "Id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}
