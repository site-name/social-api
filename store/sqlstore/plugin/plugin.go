package plugin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	defaultPluginKeyFetchLimit = 10
)

type SqlPluginStore struct {
	store.Store
}

func NewSqlPluginStore(s store.Store) store.PluginStore {
	return &SqlPluginStore{s}
}

func (ps *SqlPluginStore) Upsert(kv model.PluginKeyValue) (*model.PluginKeyValue, error) {
	if err := model_helper.PluginKeyValueIsValid(kv); err != nil {
		return nil, err
	}

	if !kv.Pvalue.Valid {
		// Setting a key to nil is the same as removing it
		err := ps.Delete(kv.PluginID, kv.Pkey)
		if err != nil {
			return nil, err
		}

		return &kv, nil
	}

	_, err := queries.Raw(
		fmt.Sprintf(
			`INSERT INTO %[1]s (
				%[2]s, 
				%[3]s, 
				%[4]s, 
				%[5]s
			) VALUES (
				$1, $2, $3, $4
			) ON CONFLICT (
				%[2]s,
				%[3]s
			) DO UPDATE SET %[4]s = $3, %[5]s = $4`,
			model.TableNames.PluginKeyValueStore, // 1
			model.PluginKeyValueColumns.PluginID, // 2
			model.PluginKeyValueColumns.Pkey,     // 3
			model.PluginKeyValueColumns.Pvalue,   // 4
			model.PluginKeyValueColumns.ExpireAt, // 5
		),
		kv.PluginID,
		kv.Pkey,
		kv.Pvalue,
		kv.ExpireAt,
	).Exec(ps.GetMaster())
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert PluginKeyValue")
	}

	return &kv, nil
}

func (ps *SqlPluginStore) CompareAndSet(kv model.PluginKeyValue, oldValue []byte) (bool, error) {
	if err := model_helper.PluginKeyValueIsValid(kv); err != nil {
		return false, nil
	}

	if kv.Pvalue.Valid == false {
		// Setting a key to nil is the same as removing it
		return ps.CompareAndDelete(kv, oldValue)
	}

	if oldValue == nil {
		// Delete any existing, expired value.
		_, err := model.PluginKeyValues(
			model.PluginKeyValueWhere.PluginID.EQ(kv.PluginID),
			model.PluginKeyValueWhere.Pkey.EQ(kv.Pkey),
			model.PluginKeyValueWhere.ExpireAt.NEQ(model_types.NewNullInt64(0)),
			model.PluginKeyValueWhere.ExpireAt.LT(model_types.NewNullInt64(model_helper.GetMillis())),
		).DeleteAll(ps.GetMaster())
		if err != nil {
			return false, err
		}

		err = kv.Insert(ps.GetMaster(), boil.Infer())
		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"PRIMARY", model.PluginKeyValueColumns.PluginID, model.PluginKeyValueColumns.Pkey}) {
				return false, nil
			}
			return false, err
		}
	} else {
		currentTime := model_helper.GetMillis()

		_, err := model.PluginKeyValues(
			model.PluginKeyValueWhere.PluginID.EQ(kv.PluginID),
			model.PluginKeyValueWhere.Pkey.EQ(kv.Pkey),
			model.PluginKeyValueWhere.Pvalue.EQ(null.Bytes{Bytes: oldValue, Valid: true}),
			model_helper.Or{
				squirrel.Eq{model.PluginKeyValueColumns.ExpireAt: 0},
				squirrel.Gt{model.PluginKeyValueColumns.ExpireAt: currentTime},
			},
		).UpdateAll(ps.GetMaster(), model.M{
			model.PluginKeyValueColumns.Pvalue:   kv.Pvalue,
			model.PluginKeyValueColumns.ExpireAt: kv.ExpireAt,
		})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (ps *SqlPluginStore) CompareAndDelete(kv model.PluginKeyValue, oldValue []byte) (bool, error) {
	if err := model_helper.PluginKeyValueIsValid(kv); err != nil {
		return false, nil
	}

	if oldValue == nil {
		// nil can't be stored. Return showing that we didn't do anything
		return false, nil
	}
	_, err := model.PluginKeyValues(
		model.PluginKeyValueWhere.PluginID.EQ(kv.PluginID),
		model.PluginKeyValueWhere.Pkey.EQ(kv.Pkey),
		model.PluginKeyValueWhere.Pvalue.EQ(null.Bytes{Bytes: oldValue, Valid: true}),
		model_helper.Or{
			squirrel.Eq{model.PluginKeyValueColumns.ExpireAt: 0},
			squirrel.Gt{model.PluginKeyValueColumns.ExpireAt: model_helper.GetMillis()},
		},
	).DeleteAll(ps.GetMaster())
	if err != nil {
		return false, errors.Wrap(err, "failed to delete PluginKeyValue")
	}

	return true, nil
}

func (ps *SqlPluginStore) SetWithOptions(pluginId string, key string, value []byte, opt model_helper.PluginKVSetOptions) (bool, error) {
	if err := opt.IsValid(); err != nil {
		return false, err
	}

	kv := model_helper.NewPluginKeyValueFromOptions(pluginId, key, value, opt)

	if opt.Atomic {
		return ps.CompareAndSet(*kv, opt.OldValue)
	}

	savedKv, nErr := ps.Upsert(*kv)
	if nErr != nil {
		return false, nErr
	}

	return savedKv != nil, nil
}

func (ps *SqlPluginStore) Get(pluginId, key string) (*model.PluginKeyValue, error) {
	record, err := model.PluginKeyValues(
		model.PluginKeyValueWhere.PluginID.EQ(pluginId),
		model.PluginKeyValueWhere.Pkey.EQ(key),
		model_helper.Or{
			squirrel.Eq{model.PluginKeyValueColumns.ExpireAt: 0},
			squirrel.Gt{model.PluginKeyValueColumns.ExpireAt: model_helper.GetMillis()},
		},
	).One(ps.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.PluginKeyValueStore, "conds")
		}
		return nil, err
	}

	return record, nil
}

func (ps *SqlPluginStore) Delete(pluginId, key string) error {
	_, err := model.PluginKeyValues(
		model.PluginKeyValueWhere.PluginID.EQ(pluginId),
		model.PluginKeyValueWhere.Pkey.EQ(key),
	).DeleteAll(ps.GetMaster())
	return err
}

func (ps *SqlPluginStore) DeleteAllForPlugin(pluginId string) error {
	_, err := model.PluginKeyValues(
		model.PluginKeyValueWhere.PluginID.EQ(pluginId),
	).DeleteAll(ps.GetMaster())
	return err
}

func (ps *SqlPluginStore) DeleteAllExpired() error {
	currentTime := model_helper.GetMillis()
	_, err := model.PluginKeyValues(
		model.PluginKeyValueWhere.ExpireAt.NEQ(model_types.NewNullInt64(0)),
		model.PluginKeyValueWhere.ExpireAt.LT(model_types.NewNullInt64(currentTime)),
	).DeleteAll(ps.GetMaster())
	return err
}

func (ps *SqlPluginStore) List(pluginId string, offset int, limit int) ([]string, error) {
	if limit <= 0 {
		limit = defaultPluginKeyFetchLimit
	}

	if offset <= 0 {
		offset = 0
	}

	var keys []string

	err := model.PluginKeyValues(
		qm.Select(model.PluginKeyValueColumns.Pkey),
		model.PluginKeyValueWhere.PluginID.EQ(pluginId),
		model_helper.Or{
			squirrel.Eq{model.PluginKeyValueColumns.ExpireAt: 0},
			squirrel.Gt{model.PluginKeyValueColumns.ExpireAt: model_helper.GetMillis()},
		},
		qm.OrderBy(fmt.Sprintf("%s %s", model.PluginKeyValueColumns.Pkey, model_helper.ASC)),
		qm.Limit(limit),
		qm.Offset(offset),
	).Query.Bind(context.Background(), ps.GetReplica(), &keys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find plugin keys of plugins")
	}

	return keys, nil
}
