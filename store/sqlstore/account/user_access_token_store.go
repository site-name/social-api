package account

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlUserAccessTokenStore struct {
	store.Store
}

func NewSqlUserAccessTokenStore(sqlStore store.Store) store.UserAccessTokenStore {
	return &SqlUserAccessTokenStore{sqlStore}
}

func (s *SqlUserAccessTokenStore) Save(token *model.UserAccessToken) (*model.UserAccessToken, error) {
	err := s.GetMaster().Create(token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *SqlUserAccessTokenStore) Delete(tokenId string) error {
	transaction := s.GetMaster().Begin()
	defer transaction.Rollback()

	if err := s.deleteSessionsAndTokensById(transaction, tokenId); err == nil {
		if err := transaction.Commit().Error; err != nil {
			// don't need to rollback here since the transaction is already closed
			return errors.Wrap(err, "commit_transaction")
		}
	}

	return nil
}

func (s *SqlUserAccessTokenStore) deleteSessionsAndTokensById(transaction *gorm.DB, tokenId string) error {
	err := transaction.Raw("DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.Id = ?", tokenId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken id=%s", tokenId)
	}

	return s.deleteTokensById(transaction, tokenId)
}

func (s *SqlUserAccessTokenStore) deleteTokensById(transaction *gorm.DB, tokenId string) error {
	err := transaction.Raw("DELETE FROM UserAccessTokens WHERE Id = ?", tokenId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete UserAccessToken id=%s", tokenId)
	}

	return nil
}

func (s *SqlUserAccessTokenStore) DeleteAllForUser(userId string) error {
	transaction := s.GetMaster().Begin()
	defer transaction.Rollback()

	if err := s.deleteSessionsandTokensByUser(transaction, userId); err != nil {
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlUserAccessTokenStore) deleteSessionsandTokensByUser(transaction *gorm.DB, userId string) error {
	err := transaction.Raw("DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.UserId = ?", userId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken userId=%s", userId)
	}

	return s.deleteTokensByUser(transaction, userId)
}

func (s *SqlUserAccessTokenStore) deleteTokensByUser(transaction *gorm.DB, userId string) error {
	err := transaction.Raw("DELETE FROM UserAccessTokens WHERE UserId = ?", userId).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete UserAccessToken userId=%s", userId)
	}

	return nil
}

func (s *SqlUserAccessTokenStore) Get(tokenId string) (*model.UserAccessToken, error) {
	var token model.UserAccessToken

	if err := s.GetReplica().First(&token, "Id = ?", tokenId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("UserAccessToken", tokenId)
		}
		return nil, errors.Wrapf(err, "failed to get UserAccessToken with id=%s", tokenId)
	}

	return &token, nil
}

func (s *SqlUserAccessTokenStore) GetAll(offset, limit int) ([]*model.UserAccessToken, error) {
	tokens := []*model.UserAccessToken{}

	if err := s.GetReplica().Raw("SELECT * FROM "+model.UserAccessTokenTableName+" OFFSET ? LIMIT ?", offset, limit).Scan(&tokens).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find UserAccessTokens")
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) GetByToken(tokenString string) (*model.UserAccessToken, error) {
	var token model.UserAccessToken

	if err := s.GetReplica().First(&token, "Token = ?", tokenString).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("UserAccessToken", fmt.Sprintf("token=%s", tokenString))
		}
		return nil, errors.Wrapf(err, "failed to get UserAccessToken with token=%s", tokenString)
	}

	return &token, nil
}

func (s *SqlUserAccessTokenStore) GetByUser(userId string, offset, limit int) ([]*model.UserAccessToken, error) {
	tokens := []*model.UserAccessToken{}

	if err := s.GetReplica().Offset(offset).Limit(limit).Find(&tokens, "UserId = ?", userId).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find UserAccessTokens with userId=%s", userId)
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) Search(term string) ([]*model.UserAccessToken, error) {
	term = store.SanitizeSearchTerm(term, "\\")
	tokens := []*model.UserAccessToken{}

	if err := s.GetReplica().
		Table(model.UserAccessTokenTableName).
		Joins("INNER JOIN "+model.UserTableName+" ON Users.Id = UserAccessTokens.UserId").
		Find(&tokens, "UserAccessTokens.Id LIKE ? OR UserAccessTokens.UserId LIKE OR Users.Username LIKE ?", term, term, term).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find UserAccessTokens by term with value '%s'", term)
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) UpdateTokenEnable(tokenId string) error {
	if err := s.GetMaster().Exec("UPDATE UserAccessTokens SET IsActive = TRUE WHERE Id = ?", tokenId).Error; err != nil {
		return errors.Wrapf(err, "failed to update UserAccessTokens with id=%s", tokenId)
	}
	return nil
}

func (s *SqlUserAccessTokenStore) UpdateTokenDisable(tokenId string) error {
	transaction := s.GetMaster().Begin()
	defer transaction.Rollback()

	if err := s.deleteSessionsAndDisableToken(transaction, tokenId); err != nil {
		return err
	}
	if err := transaction.Commit().Error; err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlUserAccessTokenStore) deleteSessionsAndDisableToken(transaction *gorm.DB, tokenId string) error {
	if err := transaction.Exec("DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.Id = ?", tokenId).Error; err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken id=%s", tokenId)
	}

	return s.updateTokenDisable(transaction, tokenId)
}

func (s *SqlUserAccessTokenStore) updateTokenDisable(transaction *gorm.DB, tokenId string) error {
	if err := transaction.Exec("UPDATE UserAccessTokens SET IsActive = FALSE WHERE Id = ?", tokenId).Error; err != nil {
		return errors.Wrapf(err, "failed to update UserAccessToken with id=%s", tokenId)
	}

	return nil
}
