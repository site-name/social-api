package account

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
)

type SqlUserTermOfServiceStore struct {
	store.Store
}

const (
	UserTermOfServiceTableName = "UserTermOfServices"
)

func NewSqlUserTermOfServiceStore(s store.Store) store.UserTermOfServiceStore {
	uts := &SqlUserTermOfServiceStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(account.UserTermsOfService{}, UserTermOfServiceTableName)
		table.ColMap("UserId").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("TermsOfServiceId").SetMaxSize(store.UUID_MAX_LENGTH)
	}
	return uts
}

func (uts *SqlUserTermOfServiceStore) CreateIndexesIfNotExists() {

}

func (s *SqlUserTermOfServiceStore) GetByUser(userId string) (*account.UserTermsOfService, error) {
	var userTermsOfService *account.UserTermsOfService

	err := s.GetReplica().SelectOne(&userTermsOfService, "SELECT * FROM UserTermsOfService WHERE UserId = :userId", map[string]interface{}{"userId": userId})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("UserTermsOfService", "userId="+userId)
		}
		return nil, errors.Wrapf(err, "failed to get UserTermsOfService with userId=%s", userId)
	}
	return userTermsOfService, nil
}

func (s *SqlUserTermOfServiceStore) Save(userTermsOfService *account.UserTermsOfService) (*account.UserTermsOfService, error) {
	userTermsOfService.PreSave()
	if err := userTermsOfService.IsValid(); err != nil {
		return nil, err
	}

	c, err := s.GetMaster().Update(userTermsOfService)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update UserTermsOfService with userId=%s and termsOfServiceId=%s", userTermsOfService.UserId, userTermsOfService.TermsOfServiceId)
	}

	if c == 0 {
		if err := s.GetMaster().Insert(userTermsOfService); err != nil {
			return nil, errors.Wrapf(err, "failed to save UserTermsOfService with userId=%s and termsOfServiceId=%s", userTermsOfService.UserId, userTermsOfService.TermsOfServiceId)
		}
	}

	return userTermsOfService, nil
}

func (s *SqlUserTermOfServiceStore) Delete(userId, termsOfServiceId string) error {
	if _, err := s.GetMaster().Exec("DELETE FROM UserTermsOfService WHERE UserId = :UserId AND TermsOfServiceId = :TermsOfServiceId", map[string]interface{}{"UserId": userId, "TermsOfServiceId": termsOfServiceId}); err != nil {
		return errors.Wrapf(err, "failed to delete UserTermsOfService with userId=%s and termsOfServiceId=%s", userId, termsOfServiceId)
	}
	return nil
}
