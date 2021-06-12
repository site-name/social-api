package preference

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

type SqlPreferenceStore struct {
	store.Store
}

func NewSqlPreferenceStore(sqlStore store.Store) store.PreferenceStore {
	s := &SqlPreferenceStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Preference{}, "Preferences").SetKeys(false, "UserId", "Category", "Name")
		table.ColMap("UserId").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Category").SetMaxSize(32)
		table.ColMap("Name").SetMaxSize(32)
		table.ColMap("Value").SetMaxSize(2000)
	}

	return s
}

func (s SqlPreferenceStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_preferences_user_id", "Preferences", "UserId")
	s.CreateIndexIfNotExists("idx_preferences_category", "Preferences", "Category")
	s.CreateIndexIfNotExists("idx_preferences_name", "Preferences", "Name")
}

func (s SqlPreferenceStore) DeleteUnusedFeatures() {
	slog.Debug("Deleting any unused pre-release features")
	sql, args, err := s.GetQueryBuilder().
		Delete("Preferences").
		Where(sq.Eq{"Category": model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS}).
		Where(sq.Eq{"Value": "false"}).
		Where(sq.Like{"Name": store.FeatureTogglePrefix + "%"}).ToSql()
	if err != nil {
		slog.Warn(errors.Wrap(err, "could not build sql query to delete unused features!").Error())
	}
	if _, err = s.GetMaster().Exec(sql, args...); err != nil {
		slog.Warn("Failed to delete unused features", slog.Err(err))
	}
}

func (s SqlPreferenceStore) Save(preferences *model.Preferences) error {
	// wrap in a transaction so that if one fails, everything fails
	transaction, err := s.GetMaster().Begin()
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}

	defer store.FinalizeTransaction(transaction)
	for _, preference := range *preferences {
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

func (s SqlPreferenceStore) save(transaction *gorp.Transaction, preference *model.Preference) error {
	preference.PreUpdate()

	if err := preference.IsValid(); err != nil {
		return err
	}
	// postgres has no way to upsert values until version 9.5 and trying inserting and then updating causes transactions to abort
	queryString, args, err := s.GetQueryBuilder().
		Select("count(0)").
		From("Preferences").
		Where(sq.Eq{"UserId": preference.UserId}).
		Where(sq.Eq{"Category": preference.Category}).
		Where(sq.Eq{"Name": preference.Name}).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to generate sqlquery")
	}

	count, err := transaction.SelectInt(queryString, args...)
	if err != nil {
		return errors.Wrap(err, "failed to count Preferences")
	}

	if count == 1 {
		return s.update(transaction, preference)
	}
	return s.insert(transaction, preference)
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
	var preference *model.Preference
	query, args, err := s.GetQueryBuilder().
		Select("*").
		From("Preferences").
		Where(sq.Eq{"UserId": userId}).
		Where(sq.Eq{"Category": category}).
		Where(sq.Eq{"Name": name}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "could not build sql query to get preference")
	}
	if err = s.GetReplica().SelectOne(&preference, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s, name=%s", userId, category, name)
	}

	return preference, nil
}

func (s SqlPreferenceStore) GetCategory(userId string, category string) (model.Preferences, error) {
	var preferences model.Preferences
	query, args, err := s.GetQueryBuilder().
		Select("*").
		From("Preferences").
		Where(sq.Eq{"UserId": userId}).
		Where(sq.Eq{"Category": category}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build sql query to get preference")
	}
	if _, err = s.GetReplica().Select(&preferences, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s", userId, category)
	}
	return preferences, nil

}

func (s SqlPreferenceStore) GetAll(userId string) (model.Preferences, error) {
	var preferences model.Preferences
	query, args, err := s.GetQueryBuilder().
		Select("*").
		From("Preferences").
		Where(sq.Eq{"UserId": userId}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not build sql query to get preference")
	}
	if _, err = s.GetReplica().Select(&preferences, query, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s", userId)
	}
	return preferences, nil
}

func (s SqlPreferenceStore) PermanentDeleteByUser(userId string) error {
	sql, args, err := s.GetQueryBuilder().
		Delete("Preferences").
		Where(sq.Eq{"UserId": userId}).ToSql()
	if err != nil {
		return errors.Wrap(err, "could not build sql query to get delete preference by user")
	}
	if _, err := s.GetMaster().Exec(sql, args...); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s", userId)
	}
	return nil
}

func (s SqlPreferenceStore) Delete(userId, category, name string) error {

	sql, args, err := s.GetQueryBuilder().
		Delete("Preferences").
		Where(sq.Eq{"UserId": userId}).
		Where(sq.Eq{"Category": category}).
		Where(sq.Eq{"Name": name}).ToSql()

	if err != nil {
		return errors.Wrap(err, "could not build sql query to get delete preference")
	}

	if _, err = s.GetMaster().Exec(sql, args...); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s, category=%s and name=%s", userId, category, name)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategory(userId string, category string) error {

	sql, args, err := s.GetQueryBuilder().
		Delete("Preferences").
		Where(sq.Eq{"UserId": userId}).
		Where(sq.Eq{"Category": category}).ToSql()

	if err != nil {
		return errors.Wrap(err, "could not build sql query to get delete preference by category")
	}

	if _, err = s.GetMaster().Exec(sql, args...); err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s and category=%s", userId, category)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategoryAndName(category string, name string) error {
	sql, args, err := s.GetQueryBuilder().
		Delete("Preferences").
		Where(sq.Eq{"Name": name}).
		Where(sq.Eq{"Category": category}).ToSql()

	if err != nil {
		return errors.Wrap(err, "could not build sql query to get delete preference by category and name")
	}

	if _, err = s.GetMaster().Exec(sql, args...); err != nil {
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

	sqlResult, err := s.GetMaster().Exec(query, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to delete Preference")
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return int64(0), errors.Wrap(err, "unable to get rows affected")
	}

	return rowsAffected, nil
}
