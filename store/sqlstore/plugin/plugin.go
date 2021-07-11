package plugin

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/store"
)

const (
	defaultPluginKeyFetchLimit = 10
)

type SqlPluginStore struct {
	store.Store
}

func NewSqlPluginStore(s store.Store) store.PluginStore {
	ps := &SqlPluginStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(plugins.PluginKeyValue{}, store.PluginKeyValueStoreTableName).SetKeys(false, "PluginId", "Key")
		table.ColMap("PluginId").SetMaxSize(plugins.KEY_VALUE_PLUGIN_ID_MAX_RUNES)
		table.ColMap("Key").SetMaxSize(plugins.KEY_VALUE_KEY_MAX_RUNES)
		table.ColMap("Value").SetMaxSize(8192)
	}
	return ps
}

func (ps *SqlPluginStore) CreateIndexesIfNotExists() {
}

func (ps *SqlPluginStore) SaveOrUpdate(kv *plugins.PluginKeyValue) (*plugins.PluginKeyValue, error) {
	if err := kv.IsValid(); err != nil {
		return nil, err
	}

	if kv.Value == nil {
		// Setting a key to nil is the same as removing it
		err := ps.Delete(kv.PluginId, kv.Key)
		if err != nil {
			return nil, err
		}

		return kv, nil
	}

	// Unfortunately PostgreSQL pre-9.5 does not have an atomic upsert, so we use
	// separate update and insert queries to accomplish our upsert
	if rowsAffected, err := ps.GetMaster().Update(kv); err != nil {
		return nil, errors.Wrap(err, "failed to update PluginKeyValue")
	} else if rowsAffected == 0 {
		// No rows were affected by the update, so let's try an insert
		if err := ps.GetMaster().Insert(kv); err != nil {
			return nil, errors.Wrap(err, "failed to save PluginKeyValue")
		}
	}

	return kv, nil
}

func (ps *SqlPluginStore) CompareAndSet(kv *plugins.PluginKeyValue, oldValue []byte) (bool, error) {
	if err := kv.IsValid(); err != nil {
		return false, err
	}

	if kv.Value == nil {
		// Setting a key to nil is the same as removing it
		return ps.CompareAndDelete(kv, oldValue)
	}

	if oldValue == nil {
		// Delete any existing, expired value.
		query := ps.GetQueryBuilder().
			Delete("PluginKeyValueStore").
			Where(sq.Eq{"PluginId": kv.PluginId}).
			Where(sq.Eq{"PKey": kv.Key}).
			Where(sq.NotEq{"ExpireAt": int(0)}).
			Where(sq.Lt{"ExpireAt": model.GetMillis()})

		queryString, args, err := query.ToSql()
		if err != nil {
			return false, errors.Wrap(err, "plugin_tosql")
		}

		if _, err := ps.GetMaster().Exec(queryString, args...); err != nil {
			return false, errors.Wrap(err, "failed to delete PluginKeyValue")
		}

		// Insert if oldValue is nil
		if err := ps.GetMaster().Insert(kv); err != nil {
			// If the error is from unique constraints violation, it's the result of a
			// race condition, return false and no error. Otherwise we have a real error and
			// need to return it.
			if ps.IsUniqueConstraintError(err, []string{"PRIMARY", "PluginId", "Key", "PKey", "pkey"}) {
				return false, nil
			}
			return false, errors.Wrap(err, "failed to insert PluginKeyValue")
		}
	} else {
		currentTime := model.GetMillis()

		// Update if oldValue is not nil
		query := ps.GetQueryBuilder().
			Update("PluginKeyValueStore").
			Set("PValue", kv.Value).
			Set("ExpireAt", kv.ExpireAt).
			Where(sq.Eq{"PluginId": kv.PluginId}).
			Where(sq.Eq{"PKey": kv.Key}).
			Where(sq.Eq{"PValue": oldValue}).
			Where(sq.Or{
				sq.Eq{"ExpireAt": int(0)},
				sq.Gt{"ExpireAt": currentTime},
			})

		queryString, args, err := query.ToSql()
		if err != nil {
			return false, errors.Wrap(err, "plugin_tosql")
		}

		updateResult, err := ps.GetMaster().Exec(queryString, args...)
		if err != nil {
			return false, errors.Wrap(err, "failed to update PluginKeyValue")
		}

		if rowsAffected, err := updateResult.RowsAffected(); err != nil {
			// Failed to update
			return false, errors.Wrap(err, "unable to get rows affected")
		} else if rowsAffected == 0 {

			// No rows were affected by the update, where condition was not satisfied,
			// return false, but no error.
			return false, nil
		}
	}

	return true, nil
}

func (ps SqlPluginStore) CompareAndDelete(kv *plugins.PluginKeyValue, oldValue []byte) (bool, error) {
	if err := kv.IsValid(); err != nil {
		return false, err
	}

	if oldValue == nil {
		// nil can't be stored. Return showing that we didn't do anything
		return false, nil
	}

	query := ps.GetQueryBuilder().
		Delete("PluginKeyValueStore").
		Where(sq.Eq{"PluginId": kv.PluginId}).
		Where(sq.Eq{"PKey": kv.Key}).
		Where(sq.Eq{"PValue": oldValue}).
		Where(sq.Or{
			sq.Eq{"ExpireAt": int(0)},
			sq.Gt{"ExpireAt": model.GetMillis()},
		})

	queryString, args, err := query.ToSql()
	if err != nil {
		return false, errors.Wrap(err, "plugin_tosql")
	}

	deleteResult, err := ps.GetMaster().Exec(queryString, args...)
	if err != nil {
		return false, errors.Wrap(err, "failed to delete PluginKeyValue")
	}

	if rowsAffected, err := deleteResult.RowsAffected(); err != nil {
		return false, errors.Wrap(err, "unable to get rows affected")
	} else if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (ps SqlPluginStore) SetWithOptions(pluginId string, key string, value []byte, opt plugins.PluginKVSetOptions) (bool, error) {
	if err := opt.IsValid(); err != nil {
		return false, err
	}

	kv, err := plugins.NewPluginKeyValueFromOptions(pluginId, key, value, opt)
	if err != nil {
		return false, err
	}

	if opt.Atomic {
		return ps.CompareAndSet(kv, opt.OldValue)
	}

	savedKv, nErr := ps.SaveOrUpdate(kv)
	if nErr != nil {
		return false, nErr
	}

	return savedKv != nil, nil
}

func (ps SqlPluginStore) Get(pluginId, key string) (*plugins.PluginKeyValue, error) {
	currentTime := model.GetMillis()
	query := ps.GetQueryBuilder().Select("PluginId, PKey, PValue, ExpireAt").
		From("PluginKeyValueStore").
		Where(sq.Eq{"PluginId": pluginId}).
		Where(sq.Eq{"PKey": key}).
		Where(sq.Or{sq.Eq{"ExpireAt": 0}, sq.Gt{"ExpireAt": currentTime}})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "plugin_tosql")
	}

	row := ps.GetReplica().Db.QueryRow(queryString, args...)
	var kv plugins.PluginKeyValue
	if err := row.Scan(&kv.PluginId, &kv.Key, &kv.Value, &kv.ExpireAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("PluginKeyValue", fmt.Sprintf("pluginId=%s, key=%s", pluginId, key))
		}
		return nil, errors.Wrapf(err, "failed to get PluginKeyValue with pluginId=%s and key=%s", pluginId, key)
	}

	return &kv, nil
}

func (ps SqlPluginStore) Delete(pluginId, key string) error {
	query := ps.GetQueryBuilder().
		Delete("PluginKeyValueStore").
		Where(sq.Eq{"PluginId": pluginId}).
		Where(sq.Eq{"Pkey": key})

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "plugin_tosql")
	}

	if _, err := ps.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrapf(err, "failed to delete PluginKeyValue with pluginId=%s and key=%s", pluginId, key)
	}
	return nil
}

func (ps SqlPluginStore) DeleteAllForPlugin(pluginId string) error {
	query := ps.GetQueryBuilder().
		Delete("PluginKeyValueStore").
		Where(sq.Eq{"PluginId": pluginId})

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "plugin_tosql")
	}

	if _, err := ps.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrapf(err, "failed to get all PluginKeyValues with pluginId=%s ", pluginId)
	}
	return nil
}

func (ps SqlPluginStore) DeleteAllExpired() error {
	currentTime := model.GetMillis()
	query := ps.GetQueryBuilder().
		Delete("PluginKeyValueStore").
		Where(sq.NotEq{"ExpireAt": 0}).
		Where(sq.Lt{"ExpireAt": currentTime})

	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "plugin_tosql")
	}

	if _, err := ps.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrap(err, "failed to delete all expired PluginKeyValues")
	}
	return nil
}

func (ps SqlPluginStore) List(pluginId string, offset int, limit int) ([]string, error) {
	if limit <= 0 {
		limit = defaultPluginKeyFetchLimit
	}

	if offset <= 0 {
		offset = 0
	}

	var keys []string

	query := ps.GetQueryBuilder().
		Select("Pkey").
		From("PluginKeyValueStore").
		Where(sq.Eq{"PluginId": pluginId}).
		Where(sq.Or{
			sq.Eq{"ExpireAt": int(0)},
			sq.Gt{"ExpireAt": model.GetMillis()},
		}).
		OrderBy("PKey").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "plugin_tosql")
	}

	_, err = ps.GetReplica().Select(&keys, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get PluginKeyValues with pluginId=%s", pluginId)
	}

	return keys, nil
}
