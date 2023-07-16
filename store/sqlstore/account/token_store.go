package account

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlTokenStore struct {
	store.Store
}

func NewSqlTokenStore(sqlStore store.Store) store.TokenStore {
	return &SqlTokenStore{sqlStore}
}

func (s *SqlTokenStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Token",
		"CreateAt",
		"Type",
		"Extra",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (s *SqlTokenStore) Save(token *model.Token) error {
	return s.GetMaster().Create(token).Error
}

func (s *SqlTokenStore) Delete(token string) error {
	return s.GetMaster().Delete(&model.Token{}, "Token = ?", token).Error
}

func (s *SqlTokenStore) GetByToken(tokenString string) (*model.Token, error) {
	var token model.Token
	err := s.GetReplica().First(&token, "Token = ?", tokenString).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("Token", tokenString)
		}
		return nil, errors.Wrapf(err, "failed to get Token with value %s", tokenString)
	}

	return &token, nil
}

func (s *SqlTokenStore) Cleanup() {
	slog.Debug("Cleaning up token store.")

	deltime := model.GetMillis() - model.MAX_TOKEN_EXIPRY_TIME
	err := s.GetMaster().Delete(&model.Token{}, "CreateAt < ?", deltime).Error
	if err != nil {
		slog.Error("failed to delete tokens", slog.Err(err))
	}
}

func (s *SqlTokenStore) RemoveAllTokensByType(tokenType string) error {
	return s.GetMaster().Delete(&model.Token{}, "Type = ?", tokenType).Error
}

func (s *SqlTokenStore) GetAllTokensByType(tokenType model.TokenType) ([]*model.Token, error) {
	var tokens []*model.Token
	err := s.GetReplica().Find(&tokens, "Type = ?", tokenType).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
