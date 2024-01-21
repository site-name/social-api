package account

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type SqlUserAccessTokenStore struct {
	store.Store
}

func NewSqlUserAccessTokenStore(sqlStore store.Store) store.UserAccessTokenStore {
	return &SqlUserAccessTokenStore{sqlStore}
}

func (s *SqlUserAccessTokenStore) Save(token model.UserAccessToken) (*model.UserAccessToken, error) {
	model_helper.UserAccessTokenCommonPre(&token)
	if err := model_helper.UserAccessTokenIsValid(token); err != nil {
		return nil, err
	}

	err := token.Insert(s.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *SqlUserAccessTokenStore) Delete(tokenId string) error {
	tx, err := s.GetMaster().BeginTx(s.Context(), &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	// delete related session
	_, err = tx.ExecContext(
		s.Context(),
		fmt.Sprintf(
			"DELETE FROM %s USING %s WHERE %s = %s AND %s = $1",
			model.TableNames.Sessions,
			model.TableNames.UserAccessTokens,
			model.SessionTableColumns.Token,
			model.UserAccessTokenTableColumns.Token,
			model.UserAccessTokenTableColumns.ID,
		),
		tokenId)
	if err != nil {
		return errors.Wrap(err, "failed to delete related session of given user access token")
	}

	// delete user access token
	_, err = model.UserAccessTokens(model.UserAccessTokenWhere.ID.EQ(tokenId)).DeleteAll(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}

	return nil
}

func (s *SqlUserAccessTokenStore) DeleteAllForUser(userId string) error {
	tx, err := s.GetMaster().BeginTx(s.Context(), &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	// delete related user session
	_, err = tx.ExecContext(
		s.Context(),
		fmt.Sprintf(
			"DELETE FROM %s USING %s WHERE %s = %s AND %s = $1",
			model.TableNames.Sessions,
			model.TableNames.UserAccessTokens,
			model.SessionTableColumns.Token,
			model.UserAccessTokenTableColumns.Token,
			model.UserAccessTokenTableColumns.UserID,
		),
		userId)

	// delete user access token
	_, err = model.UserAccessTokens(model.UserAccessTokenWhere.UserID.EQ(userId)).DeleteAll(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlUserAccessTokenStore) Get(tokenId string) (*model.UserAccessToken, error) {
	token, err := model.FindUserAccessToken(s.GetReplica(), tokenId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.UserAccessTokens, tokenId)
		}
		return nil, err
	}

	return token, nil
}

func (s *SqlUserAccessTokenStore) GetAll(conds ...qm.QueryMod) (model.UserAccessTokenSlice, error) {
	return model.UserAccessTokens(conds...).All(s.GetReplica())
}

func (s *SqlUserAccessTokenStore) GetByToken(tokenString string) (*model.UserAccessToken, error) {
	token, err := model.UserAccessTokens(model.UserAccessTokenWhere.Token.EQ(tokenString)).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.UserAccessTokens, "token")
		}
		return nil, err
	}

	return token, nil
}

func (s *SqlUserAccessTokenStore) Search(term string) (model.UserAccessTokenSlice, error) {
	term = store.SanitizeSearchTerm(term, "\\")

	return model.
		UserAccessTokens(
			qm.InnerJoin(fmt.Sprintf(
				"%[1]s ON %[2]s = %[3]s",
				model.TableNames.Users,                   // 1
				model.UserTableColumns.ID,                // 2
				model.UserAccessTokenTableColumns.UserID, // 3
			)),
			model.UserAccessTokenWhere.ID.LIKE(term),
			qm.Or(model.UserAccessTokenTableColumns.UserID+" LIKE ?", term),
			qm.Or(model.UserTableColumns.Username+" LIKE ?", term),
		).
		All(s.GetReplica())
}

func (s *SqlUserAccessTokenStore) UpdateTokenEnable(tokenId string) error {
	_, err := model.UserAccessTokens(model.UserAccessTokenWhere.ID.EQ(tokenId)).
		UpdateAll(s.GetMaster(), model.M{
			model.UserAccessTokenColumns.IsActive: true,
		})
	return err
}

func (s *SqlUserAccessTokenStore) UpdateTokenDisable(tokenId string) error {
	tx, err := s.GetMaster().BeginTx(s.Context(), &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	// delete related session
	_, err = tx.ExecContext(
		s.Context(),
		fmt.Sprintf(
			"DELETE FROM %s USING %s WHERE %s = %s AND %s = $1",
			model.TableNames.Sessions,
			model.TableNames.UserAccessTokens,
			model.SessionTableColumns.Token,
			model.UserAccessTokenTableColumns.Token,
			model.UserAccessTokenTableColumns.ID,
		),
		tokenId)
	if err != nil {
		return errors.Wrap(err, "failed to delete related session of given user access token")
	}

	_, err = model.UserAccessTokens(model.UserAccessTokenWhere.ID.EQ(tokenId)).
		UpdateAll(s.GetMaster(), model.M{
			model.UserAccessTokenColumns.IsActive: false,
		})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}
