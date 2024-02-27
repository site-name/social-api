package preference

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlPreferenceStore struct {
	store.Store
}

func NewSqlPreferenceStore(sqlStore store.Store) store.PreferenceStore {
	return &SqlPreferenceStore{sqlStore}
}

func (s *SqlPreferenceStore) DeleteUnusedFeatures() {
	_, err := model.Preferences(
		model.PreferenceWhere.Category.EQ(model_helper.PREFERENCE_CATEGORY_ADVANCED_SETTINGS),
		model.PreferenceWhere.Value.EQ("false"),
		model.PreferenceWhere.Name.LIKE(store.FeatureTogglePrefix+"%"),
	).DeleteAll(s.GetMaster())
	if err != nil {
		slog.Warn("Failed to delete unused features", slog.Err(err))
	}
}

func (s *SqlPreferenceStore) Save(preferences model.PreferenceSlice) error {
	tx, err := s.GetMaster().BeginTx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "begin_transaction")
	}
	defer s.FinalizeTransaction(tx)

	for _, preference := range preferences {
		if preference == nil {
			continue
		}

		if upsertErr := s.save(tx, preference); upsertErr != nil {
			return upsertErr
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit_transaction")
	}
	return nil
}

func (s *SqlPreferenceStore) save(transaction boil.ContextTransactor, preference *model.Preference) error {
	model_helper.PreferencePreUpdate(preference)

	if err := model_helper.PreferenceIsValid(*preference); err != nil {
		return err
	}
	// postgres has no way to upsert values until version 9.5 and trying inserting and then updating causes transactions to abort
	return preference.Upsert(
		transaction,
		true,
		[]string{model.PreferenceColumns.UserID, model.PreferenceColumns.Category, model.PreferenceColumns.Name},
		boil.Whitelist(model.PreferenceColumns.Value),
		boil.Infer(),
	)
}

func (s *SqlPreferenceStore) Get(userId string, category string, name string) (*model.Preference, error) {
	record, err := model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
		model.PreferenceWhere.Category.EQ(category),
		model.PreferenceWhere.Name.EQ(name),
	).One(s.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Preferences, userId+"-"+category+"-"+name)
		}
		return nil, err
	}

	return record, nil
}

func (s *SqlPreferenceStore) GetCategory(userId string, category string) (model.PreferenceSlice, error) {
	return model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
		model.PreferenceWhere.Category.EQ(category),
	).All(s.GetReplica())
}

func (s *SqlPreferenceStore) GetAll(userId string) (model.PreferenceSlice, error) {
	return model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
	).All(s.GetReplica())
}

func (s *SqlPreferenceStore) PermanentDeleteByUser(userId string) error {
	_, err := model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
	).DeleteAll(s.GetMaster())
	return err
}

func (s *SqlPreferenceStore) Delete(userId, category, name string) error {
	_, err := model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
		model.PreferenceWhere.Category.EQ(category),
		model.PreferenceWhere.Name.EQ(name),
	).DeleteAll(s.GetMaster())
	return err
}

func (s *SqlPreferenceStore) DeleteCategory(userId string, category string) error {
	_, err := model.Preferences(
		model.PreferenceWhere.UserID.EQ(userId),
		model.PreferenceWhere.Category.EQ(category),
	).DeleteAll(s.GetMaster())
	return err
}

func (s *SqlPreferenceStore) DeleteCategoryAndName(category string, name string) error {
	_, err := model.Preferences(
		model.PreferenceWhere.Name.EQ(name),
		model.PreferenceWhere.Category.EQ(category),
	).DeleteAll(s.GetMaster())
	return err
}

// func (s *SqlPreferenceStore) CleanupFlagsBatch(limit int64) (int64, error) {
// if limit < 0 {
// 	// uint64 does not throw an error, it overflows if it is negative.
// 	// it is better to manually check here, or change the function type to uint64
// 	return int64(0), errors.Errorf("Received a negative limit")
// }
// nameInQ, nameInArgs, err := sq.Select("*").
// 	FromSelect(
// 		sq.Select(model.PreferenceTableColumns.Name).
// 			From(model.TableNames.Preferences).
// 			LeftJoin("Posts ON Preferences.Name = Posts.Id").
// 			Where(sq.Eq{"Preferences.Category": model.PREFERENCE_CATEGORY_FLAGGED_POST}).
// 			Where(sq.Eq{"Posts.Id": nil}).
// 			Limit(uint64(limit)),
// 		"t").
// 	ToSql()
// if err != nil {
// 	return int64(0), errors.Wrap(err, "could not build nested sql query to delete preference")
// }
// query, args, err := s.GetQueryBuilder().Delete("Preferences").
// 	Where(sq.Eq{"Category": model.PREFERENCE_CATEGORY_FLAGGED_POST}).
// 	Where(sq.Expr("name IN ("+nameInQ+")", nameInArgs...)).
// 	ToSql()

// if err != nil {
// 	return int64(0), errors.Wrap(err, "could not build sql query to delete preference")
// }

// result := s.GetMaster().Raw(query, args...)
// if result.Error != nil {
// 	return int64(0), errors.Wrap(result.Error, "failed to delete Preference")
// }

// return result.RowsAffected, nil
// panic("not implemented")
// }
