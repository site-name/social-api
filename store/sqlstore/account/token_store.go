package account

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
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
	if err := token.IsValid(); err != nil {
		return err
	}

	query := "INSERT INTO " + model.TokenTableName + " (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
	if _, err := s.GetMasterX().NamedExec(query, token); err != nil {
		return errors.Wrap(err, "failed to save token")
	}

	return nil
}

func (s *SqlTokenStore) Delete(token string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM "+model.TokenTableName+" WHERE Token = ?", token); err != nil {
		return errors.Wrapf(err, "failed to delete Token with value %s", token)
	}
	return nil
}

func (s *SqlTokenStore) GetByToken(tokenString string) (*model.Token, error) {
	var token model.Token

	if err := s.GetReplicaX().Get(token, "SELECT * FROM "+model.TokenTableName+" WHERE Token = ?", token); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Token", fmt.Sprintf("Token=%s", tokenString))
		}
		return nil, errors.Wrapf(err, "failed to get Token with value %s", tokenString)
	}

	return &token, nil
}

func (s *SqlTokenStore) Cleanup() {
	slog.Debug("Cleaning up token store.")

	deltime := model.GetMillis() - model.MAX_TOKEN_EXIPRY_TIME
	if _, err := s.GetMasterX().Exec("DELETE FROM "+model.TokenTableName+" WHERE CreateAt < ?", deltime); err != nil {
		slog.Error("Unable to cleanup token store.")
	}
}

func (s *SqlTokenStore) RemoveAllTokensByType(tokenType string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM Tokens WHERE Type = ?", tokenType); err != nil {
		return errors.Wrapf(err, "failed to remove all Tokens with type=%s", tokenType)
	}

	return nil
}

func (s *SqlTokenStore) GetAllTokensByType(tokenType string) ([]*model.Token, error) {
	var tokens []*model.Token
	if err := s.GetReplicaX().Select(&tokens, "SELECT * FROM "+model.TokenTableName+" WHERE Type = ?", tokenType); err != nil {
		return nil, errors.Wrapf(err, "failed to find tokens with type=%s", tokenType)
	}

	return tokens, nil
}
