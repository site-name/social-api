package preference

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlPreferenceStore struct {
	store.Store
}

func NewSqlPreferenceStore(sqlStore store.Store) store.PreferenceStore {
	return &SqlPreferenceStore{sqlStore}
}

func (s SqlPreferenceStore) DeleteUnusedFeatures() {
	sql, args, err := s.GetQueryBuilder().
		Delete(model.PreferenceTableName).
		Where(sq.Eq{"Category": model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS}).
		Where(sq.Eq{"Value": "false"}).
		Where(sq.Like{"Name": store.FeatureTogglePrefix + "%"}).ToSql()
	if err != nil {
		slog.Warn(errors.Wrap(err, "could not build sql query to delete unused features!").Error())
	}
	if err = s.GetMaster().Raw(sql, args...).Error; err != nil {
		slog.Warn("Failed to delete unused features", slog.Err(err))
	}
}

func (s SqlPreferenceStore) Save(preferences model.Preferences) error {
	// wrap in a transaction so that if one fails, everything fails
	tx := s.GetMaster().Begin()
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	for _, preference := range preferences {
		preference := preference
		if upsertErr := s.save(tx, &preference); upsertErr != nil {
			return upsertErr
		}
	}

	if err := tx.Commit().Error; err != nil {
		// don't need to rollback here since the transaction is already closed
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s SqlPreferenceStore) save(transaction *gorm.DB, preference *model.Preference) error {
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

	if err = transaction.Raw(queryString, args...).Error; err != nil {
		return errors.Wrap(err, "failed to count Preferences")
	}

	return nil
}

func (s SqlPreferenceStore) insert(transaction *gorm.DB, preference *model.Preference) error {
	if err := transaction.Create(preference).Error; err != nil {
		if s.IsUniqueConstraintError(err, []string{"UserId", "preferences_pkey"}) {
			return store.NewErrInvalidInput("Preference", "<userId, category, name>", fmt.Sprintf("<%s, %s, %s>", preference.UserId, preference.Category, preference.Name))
		}
		return errors.Wrapf(err, "failed to save Preference with userId=%s, category=%s, name=%s", preference.UserId, preference.Category, preference.Name)
	}

	return nil
}

func (s SqlPreferenceStore) update(transaction *gorm.DB, preference *model.Preference) error {
	if err := transaction.Save(preference).Error; err != nil {
		return errors.Wrapf(err, "failed to update Preference with userId=%s, category=%s, name=%s", preference.UserId, preference.Category, preference.Name)
	}

	return nil
}

func (s SqlPreferenceStore) Get(userId string, category string, name string) (*model.Preference, error) {
	var preference model.Preference
	if err := s.GetReplica().First(&preference, "UserId = ? AND Category = ? AND Name = ?", userId, category, name).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s, name=%s", userId, category, name)
	}

	return &preference, nil
}

func (s SqlPreferenceStore) GetCategory(userId string, category string) (model.Preferences, error) {
	var preferences model.Preferences
	if err := s.GetReplica().Find(&preferences, "UserId = ? AND Category = ?", userId, category).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s, category=%s", userId, category)
	}
	return preferences, nil

}

func (s SqlPreferenceStore) GetAll(userId string) (model.Preferences, error) {
	var preferences model.Preferences
	if err := s.GetReplica().Find(&preferences, "UserId = ?", userId).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find Preference with userId=%s", userId)
	}
	return preferences, nil
}

func (s SqlPreferenceStore) PermanentDeleteByUser(userId string) error {
	if err := s.GetMaster().Raw("DELETE FROM "+model.PreferenceTableName+" WHERE UserId = ?", userId).Error; err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s", userId)
	}
	return nil
}

func (s SqlPreferenceStore) Delete(userId, category, name string) error {
	if err := s.GetMaster().Raw("DELETE FROM "+model.PreferenceTableName+" WHERE UserId = ? AND Category = ? AND Name = ?", userId, category, name).Error; err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s, category=%s and name=%s", userId, category, name)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategory(userId string, category string) error {
	if err := s.GetMaster().Exec("DELETE FROM "+model.PreferenceTableName+" WHERE UserId = ? AND Category = ?", userId, category).Error; err != nil {
		return errors.Wrapf(err, "failed to delete Preference with userId=%s and category=%s", userId, category)
	}

	return nil
}

func (s SqlPreferenceStore) DeleteCategoryAndName(category string, name string) error {
	if err := s.GetMaster().Raw("DELETE FROM "+model.PreferenceTableName+" WHERE Name = ? AND Category = ?", name, category).Error; err != nil {
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

	result := s.GetMaster().Raw(query, args...)
	if result.Error != nil {
		return int64(0), errors.Wrap(result.Error, "failed to delete Preference")
	}

	return result.RowsAffected, nil
}
