package plugin

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
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

func (ps *SqlPluginStore) SaveOrUpdate(kv *model.PluginKeyValue) (*model.PluginKeyValue, error) {
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

	query := ps.GetQueryBuilder().
		Insert(model.PluginKeyValueStoreTableName).
		Columns("PluginId", "PKey", "PValue", "ExpireAt").
		Values(kv.PluginId, kv.Key, kv.Value, kv.ExpireAt).
		SuffixExpr(
			squirrel.Expr("ON CONFLICT (pluginid, pkey) DO UPDATE SET PValue = ?, ExpireAt = ?", kv.Value, kv.ExpireAt),
		)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "plugin_tosql")
	}

	if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
		return nil, errors.Wrap(err, "failed to upsert PluginKeyValue")
	}

	return kv, nil
}

func (ps *SqlPluginStore) CompareAndSet(kv *model.PluginKeyValue, oldValue []byte) (bool, error) {
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

		if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
			return false, errors.Wrap(err, "failed to delete PluginKeyValue")
		}

		// Insert if oldValue is nil
		queryString, args, err = ps.GetQueryBuilder().
			Insert("PluginKeyValueStore").
			Columns("PluginId", "PKey", "PValue", "ExpireAt").
			Values(kv.PluginId, kv.Key, kv.Value, kv.ExpireAt).ToSql()
		if err != nil {
			return false, errors.Wrap(err, "plugin_tosql")
		}
		if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
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

		err = ps.GetMaster().Raw(queryString, args...).Error
		if err != nil {
			return false, errors.Wrap(err, "failed to update PluginKeyValue")
		}
	}

	return true, nil
}

func (ps SqlPluginStore) CompareAndDelete(kv *model.PluginKeyValue, oldValue []byte) (bool, error) {
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

	err = ps.GetMaster().Raw(queryString, args...).Error
	if err != nil {
		return false, errors.Wrap(err, "failed to delete PluginKeyValue")
	}

	return true, nil
}

func (ps SqlPluginStore) SetWithOptions(pluginId string, key string, value []byte, opt model.PluginKVSetOptions) (bool, error) {
	if err := opt.IsValid(); err != nil {
		return false, err
	}

	kv, err := model.NewPluginKeyValueFromOptions(pluginId, key, value, opt)
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

func (ps SqlPluginStore) Get(pluginId, key string) (*model.PluginKeyValue, error) {
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

	var kv model.PluginKeyValue

	err = ps.GetReplica().Raw(queryString, args...).Scan(&kv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
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

	if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
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

	if err := ps.GetMaster().Raw(queryString, args...).Error; err != nil {
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

	err = ps.GetReplica().Raw(queryString, args...).Scan(&keys).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get PluginKeyValues with pluginId=%s", pluginId)
	}

	return keys, nil
}
