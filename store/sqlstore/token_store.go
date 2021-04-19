package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/store"
)

type SqlTokenStore struct {
	*SqlStore
}

func newSqlTokenStore(sqlStore *SqlStore) store.TokenStore {
	s := &SqlTokenStore{SqlStore: sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Token{}, "Tokens").SetKeys(false, "Token")
		table.ColMap("Token").SetMaxSize(model.TOKEN_SIZE)
		table.ColMap("Type").SetMaxSize(model.TOKEN_SIZE)
		table.ColMap("Extra").SetMaxSize(model.MAX_EXTRA)
	}

	return s
}

func (s SqlTokenStore) createIndexesIfNotExists() {
}

func (s *SqlTokenStore) Save(token *model.Token) error {
	if err := token.IsValid(); err != nil {
		return err
	}

	if err := s.GetMaster().Insert(token); err != nil {
		return errors.Wrap(err, "failed to save token")
	}

	return nil
}

func (s *SqlTokenStore) Delete(token string) error {
	if _, err := s.GetMaster().Exec("DELETE FROM Tokens WHERE Token = :Token", map[string]interface{}{"Token": token}); err != nil {
		return errors.Wrapf(err, "failed to delete Token with value %s", token)
	}
	return nil
}

func (s *SqlTokenStore) GetByToken(tokenString string) (*model.Token, error) {
	token := new(model.Token)

	if err := s.GetReplica().SelectOne(token, "SELECT * FROM Tokens WHERE Token = :Token", map[string]interface{}{"Token": tokenString}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Token", fmt.Sprintf("Token=%s", tokenString))
		}
		return nil, errors.Wrapf(err, "failed to get Token with value %s", tokenString)
	}

	return token, nil
}

func (s *SqlTokenStore) Cleanup() {
	log.Debug("Cleaning up token store.")
	deltime := model.GetMillis() - model.MAX_TOKEN_EXIPRY_TIME
	if _, err := s.GetMaster().Exec("DELETE FROM Tokens WHERE CreateAt < :DelTime", map[string]interface{}{"DelTime": deltime}); err != nil {
		log.Error("Unable to cleanup token store.")
	}
}

func (s *SqlTokenStore) RemoveAllTokensByType(tokenType string) error {
	if _, err := s.GetMaster().Exec("DELETE FROM Tokens WHERE Type = :TokenType", map[string]interface{}{"TokenType": tokenType}); err != nil {
		return errors.Wrapf(err, "failed to remove all Tokens with type=%s", tokenType)
	}

	return nil
}
