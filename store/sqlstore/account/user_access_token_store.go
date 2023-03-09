package account

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlUserAccessTokenStore struct {
	store.Store
}

func NewSqlUserAccessTokenStore(sqlStore store.Store) store.UserAccessTokenStore {
	return &SqlUserAccessTokenStore{sqlStore}
}

func (s *SqlUserAccessTokenStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"Token",
		"UserID",
		"Description",
		"IsActive",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (s *SqlUserAccessTokenStore) Save(token *model.UserAccessToken) (*model.UserAccessToken, error) {
	token.PreSave()

	if err := token.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.UserAccessTokenTableName + " (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
	if _, err := s.GetMasterX().NamedExec(query, token); err != nil {
		return nil, errors.Wrap(err, "failed to save UserAccessToken")
	}
	return token, nil
}

func (s *SqlUserAccessTokenStore) Delete(tokenId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}

	defer store.FinalizeTransaction(transaction)

	if err := s.deleteSessionsAndTokensById(transaction, tokenId); err == nil {
		if err := transaction.Commit(); err != nil {
			// don't need to rollback here since the transaction is already closed
			return errors.Wrap(err, "commit_transaction")
		}
	}

	return nil

}

func (s *SqlUserAccessTokenStore) deleteSessionsAndTokensById(transaction store_iface.SqlxTxExecutor, tokenId string) error {
	if _, err := transaction.Exec("DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.Id = ?", tokenId); err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken id=%s", tokenId)
	}

	return s.deleteTokensById(transaction, tokenId)
}

func (s *SqlUserAccessTokenStore) deleteTokensById(transaction store_iface.SqlxTxExecutor, tokenId string) error {
	if _, err := transaction.Exec("DELETE FROM UserAccessTokens WHERE Id = ?", tokenId); err != nil {
		return errors.Wrapf(err, "failed to delete UserAccessToken id=%s", tokenId)
	}

	return nil
}

func (s *SqlUserAccessTokenStore) DeleteAllForUser(userId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(transaction)

	if err := s.deleteSessionsandTokensByUser(transaction, userId); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlUserAccessTokenStore) deleteSessionsandTokensByUser(transaction store_iface.SqlxTxExecutor, userId string) error {
	query := "DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.UserId = ?"

	if _, err := transaction.Exec(query, userId); err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken userId=%s", userId)
	}

	return s.deleteTokensByUser(transaction, userId)
}

func (s *SqlUserAccessTokenStore) deleteTokensByUser(transaction store_iface.SqlxTxExecutor, userId string) error {
	if _, err := transaction.Exec("DELETE FROM UserAccessTokens WHERE UserId = ?", userId); err != nil {
		return errors.Wrapf(err, "failed to delete UserAccessToken userId=%s", userId)
	}

	return nil
}

func (s *SqlUserAccessTokenStore) Get(tokenId string) (*model.UserAccessToken, error) {
	var token model.UserAccessToken

	if err := s.GetReplicaX().Get(&token, "SELECT * FROM UserAccessTokens WHERE Id = ?", tokenId); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("UserAccessToken", tokenId)
		}
		return nil, errors.Wrapf(err, "failed to get UserAccessToken with id=%s", tokenId)
	}

	return &token, nil
}

func (s *SqlUserAccessTokenStore) GetAll(offset, limit int) ([]*model.UserAccessToken, error) {
	tokens := []*model.UserAccessToken{}

	if err := s.GetReplicaX().Select(&tokens, "SELECT * FROM UserAccessTokens LIMIT ? OFFSET ?", limit, offset); err != nil {
		return nil, errors.Wrap(err, "failed to find UserAccessTokens")
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) GetByToken(tokenString string) (*model.UserAccessToken, error) {
	var token model.UserAccessToken

	if err := s.GetReplicaX().Get(&token, "SELECT * FROM UserAccessTokens WHERE Token = ?", token); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("UserAccessToken", fmt.Sprintf("token=%s", tokenString))
		}
		return nil, errors.Wrapf(err, "failed to get UserAccessToken with token=%s", tokenString)
	}

	return &token, nil
}

func (s *SqlUserAccessTokenStore) GetByUser(userId string, offset, limit int) ([]*model.UserAccessToken, error) {
	tokens := []*model.UserAccessToken{}

	if err := s.GetReplicaX().Select(&tokens, "SELECT * FROM UserAccessTokens WHERE UserId = ? LIMIT ? OFFSET ?", userId, limit, offset); err != nil {
		return nil, errors.Wrapf(err, "failed to find UserAccessTokens with userId=%s", userId)
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) Search(term string) ([]*model.UserAccessToken, error) {
	term = store.SanitizeSearchTerm(term, "\\")
	tokens := []*model.UserAccessToken{}
	params := []interface{}{term, term, term}
	query := `
		SELECT
			uat.*
		FROM UserAccessTokens uat
		INNER JOIN Users u
			ON uat.UserId = u.Id
		WHERE uat.Id LIKE ? OR uat.UserId LIKE ? OR u.Username LIKE ?`

	if err := s.GetReplicaX().Select(&tokens, query, params...); err != nil {
		return nil, errors.Wrapf(err, "failed to find UserAccessTokens by term with value '%s'", term)
	}

	return tokens, nil
}

func (s *SqlUserAccessTokenStore) UpdateTokenEnable(tokenId string) error {
	if _, err := s.GetMasterX().Exec("UPDATE UserAccessTokens SET IsActive = TRUE WHERE Id = ?", tokenId); err != nil {
		return errors.Wrapf(err, "failed to update UserAccessTokens with id=%s", tokenId)
	}
	return nil
}

func (s *SqlUserAccessTokenStore) UpdateTokenDisable(tokenId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer store.FinalizeTransaction(transaction)

	if err := s.deleteSessionsAndDisableToken(transaction, tokenId); err != nil {
		return err
	}
	if err := transaction.Commit(); err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlUserAccessTokenStore) deleteSessionsAndDisableToken(transaction store_iface.SqlxTxExecutor, tokenId string) error {
	if _, err := transaction.Exec("DELETE FROM Sessions s USING UserAccessTokens o WHERE o.Token = s.Token AND o.Id = ?", tokenId); err != nil {
		return errors.Wrapf(err, "failed to delete Sessions with UserAccessToken id=%s", tokenId)
	}

	return s.updateTokenDisable(transaction, tokenId)
}

func (s *SqlUserAccessTokenStore) updateTokenDisable(transaction store_iface.SqlxTxExecutor, tokenId string) error {
	if _, err := transaction.Exec("UPDATE UserAccessTokens SET IsActive = FALSE WHERE Id = ?", tokenId); err != nil {
		return errors.Wrapf(err, "failed to update UserAccessToken with id=%s", tokenId)
	}

	return nil
}
