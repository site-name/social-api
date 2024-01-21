package account

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlTokenStore struct {
	store.Store
}

func NewSqlTokenStore(sqlStore store.Store) store.TokenStore {
	return &SqlTokenStore{sqlStore}
}

func (s *SqlTokenStore) Save(token model.Token) (*model.Token, error) {
	model_helper.TokenPreSave(&token)
	if err := model_helper.TokenIsValid(token); err != nil {
		return nil, err
	}

	err := token.Insert(s.GetMaster(), boil.Infer())
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *SqlTokenStore) Delete(token string) error {
	_, err := model.Tokens(model.TokenWhere.Token.EQ(token)).DeleteAll(s.GetMaster())
	return err
}

func (s *SqlTokenStore) GetByToken(tokenString string) (*model.Token, error) {
	token, err := model.FindToken(s.GetReplica(), tokenString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Tokens, tokenString)
		}
		return nil, err
	}
	return token, nil
}

func (s *SqlTokenStore) Cleanup() error {
	deltime := model_helper.GetMillis() - model_helper.MAX_TOKEN_EXIPRY_TIME
	_, err := model.Tokens(model.TokenWhere.CreatedAt.LT(deltime)).DeleteAll(s.GetMaster())
	return err
}

func (s *SqlTokenStore) GetAllTokensByType(tokenType model_helper.TokenType) (model.TokenSlice, error) {
	return model.Tokens(model.TokenWhere.Type.EQ(string(tokenType))).All(s.GetReplica())
}
