package preference

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlPreferenceStore struct {
	store.Store
}

func NewSqlPreferenceStore(sqlStore store.Store) store.PreferenceStore {
	return &SqlPreferenceStore{sqlStore}
}

func (s SqlPreferenceStore) DeleteUnusedFeatures() {
	sql, args, err := s.GetQueryBuilder().
		Delete(store.PreferenceTableName).
		Where(sq.Eq{"Category": model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS}).
		Where(sq.Eq{"Value": "false"}).
		Where(sq.Like{"Name": store.FeatureTogglePrefix + "%"}).ToSql()
	if err != nil {
		slog.Warn(errors.Wrap(err, "could not build sql query to delete unused features!").Error())
	}
	if _, err = s.GetMasterX().Exec(sql, args...); err != nil {
		slog.Warn("Failed to delete unused features", slog.Err(err))
	}
}

func (s SqlPreferenceStore) Save(preferences model.Preferences) error {
	// wrap in a transaction so that if one fails, everything fails
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}

	defer store.FinalizeTransaction(transaction)
	for _, preference := range preferences {
		preference := preference
		if upsertErr := s.save(transaction, &preference); upsertErr != nil {
			return upsertErr
		}
	}

	if err := transaction.Commit(); err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s SqlPreferenceStore) save(transaction store_iface.SqlxTxExecutor, preference *model.Preference) error {
	preference.PreUpdate()

	if err := preference.IsValid(); err != nil {
		return err
	}
	// postgres has no way to upsert values until version 9.5 and trying inserting and then updating causes transactions to abort
	queryString, args, err := s.GetQueryBuilder().
		Insert("Preferences").
		Columns("UserId", "Category", "Name", "Value").
		Values(preference.UserId, preference.Category, preference.Name, preference.Value).
		SuffixExpr(squirrel.Expr("ON CONFLICT (userid, category, name) DO UPDATE SET Value = ?", preference.Value)).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to generate sqlquery")
	}

	if _, err = transaction.Exec(queryString, args...); err != nil {
		return errors.Wrap(err, "failed to count Preferences")
	}

	return nil
}

func (s SqlPreferenceStore) insert(transaction *gorp.Transaction, preference *model.Preference) error {
	if err := transaction.Insert(preference); err != nil {
		if s.IsUniqueConstraintError(err, []string{"UserId", "preferences_pkey"}) {
			return store.NewErrInvalidInput("Preference", "<userId, category, name>", fmt.Sprintf("<%s, %s, %s>", preference.UserId, preference.Category, preference.Name))
		}
		return errors.Wrapf(err, "failed to save Preference with userId=%s, category=%s, name=%s", preference.UserId, preference.Category, preference.Name)
	}

	return nil
}

func (s SqlPreferenceStore) update(transaction *gorp.Transaction, preference *model.Preference) error {
	if _, err := transaction.Update(preference); err != nil {
		return errors.Wrapf(err, "failed to update Preference with userId=%s, category=%s, name=%s", preference.UserId, preference.Category, preference.Name)
	}

	return nil
}

func (s SqlPreferenceStore) Get(userId string, category string, name string) (*model.Preference, error) {
	var preference model.Preference
	if err := s.GetReplicaX().Select(&preference, "SELECT * FROM "+store.PreferenceTableName+" WHERE UserId = ? AND Category = ? AND Name = ?", userId, category, name); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s, name=%s", userId, category, name)
	}

	return &preference, nil
}

func (s SqlPreferenceStore) GetCategory(userId string, category string) (model.Preferences, error) {
	var preferences model.Preferences
	if err := s.GetReplicaX().Select(&preferences, "SELECT * FROM "+store.PreferenceTableName+" WHERE UserId = ? AND Category = ?", userId, category); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s", userId, category)
	}
	return preferences, nil

}

func (s SqlPreferenceStore) GetAll(userId string) (model.Preferences, error) {
	var preferences model.Preferences
	if err := s.GetReplicaX().Select(&preferences, "SELECT * FROM "+store.PreferenceTableName+" WHERE UserId = ?", userId); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s", userId)
	}
	return preferences, nil
}

func (s SqlPreferenceStore) PermanentDeleteByUser(userId string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM "+store.PreferenceTableName+" WHERE UserId = ?", userId); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s", userId)
	}
	return nil
}

func (s SqlPreferenceStore) Delete(userId, category, name string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM "+store.PreferenceTableName+" WHERE UserId = ? AND Category = ? AND Name = ?", userId, category, name); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s, category=%s and name=%s", userId, category, name)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategory(userId string, category string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM "+store.PreferenceTableName+" WHERE UserId = ? AND Category = ?", userId, category); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s and category=%s", userId, category)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategoryAndName(category string, name string) error {
	if _, err := s.GetMasterX().Exec("DELETE FROM "+store.PreferenceTableName+" WHERE Name = ? AND Category = ?", name, category); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with category=%s and name=%s", category, name)
	}

	return nil
}

func (s SqlPreferenceStore) CleanupFlagsBatch(limit int64) (int64, error) {
	if limit < 0 {
		// uint64 does not throw an error, it overflows if it is negative.
		// it is better to manually check here, or change the function type to uint64
		return int64(0), errors.Errorf("Received a negative limit")
	}
	nameInQ, nameInArgs, err := sq.Select("*").
		FromSelect(
			sq.Select("Preferences.Name").
				From("Preferences").
				LeftJoin("Posts ON Preferences.Name = Posts.Id").
				Where(sq.Eq{"Preferences.Category": model.PREFERENCE_CATEGORY_FLAGGED_POST}).
				Where(sq.Eq{"Posts.Id": nil}).
				Limit(uint64(limit)),
			"t").
		ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "could not build nested sql query to delete preference")
	}
	query, args, err := s.GetQueryBuilder().Delete("Preferences").
		Where(sq.Eq{"Category": model.PREFERENCE_CATEGORY_FLAGGED_POST}).
		Where(sq.Expr("name IN ("+nameInQ+")", nameInArgs...)).
		ToSql()

	if err != nil {
		return int64(0), errors.Wrap(err, "could not build sql query to delete preference")
	}

	sqlResult, err := s.GetMasterX().Exec(query, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to delete Preference")
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return int64(0), errors.Wrap(err, "unable to get rows affected")
	}

	return rowsAffected, nil
}
