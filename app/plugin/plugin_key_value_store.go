package plugin

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func getKeyHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (a *AppPlugin) SetPluginKey(pluginID string, key string, value []byte) *model.AppError {
	return a.SetPluginKeyWithExpiry(pluginID, key, value, 0)
}

func (a *AppPlugin) SetPluginKeyWithExpiry(pluginID string, key string, value []byte, expireInSeconds int64) *model.AppError {
	options := plugins.PluginKVSetOptions{
		ExpireInSeconds: expireInSeconds,
	}
	_, err := a.SetPluginKeyWithOptions(pluginID, key, value, options)
	return err
}

func (a *AppPlugin) CompareAndSetPluginKey(pluginID string, key string, oldValue, newValue []byte) (bool, *model.AppError) {
	options := plugins.PluginKVSetOptions{
		Atomic:   true,
		OldValue: oldValue,
	}
	return a.SetPluginKeyWithOptions(pluginID, key, newValue, options)
}

func (a *AppPlugin) SetPluginKeyWithOptions(pluginID string, key string, value []byte, options plugins.PluginKVSetOptions) (bool, *model.AppError) {
	if err := options.IsValid(); err != nil {
		slog.Debug("Failed to set plugin key value with options", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return false, err
	}

	updated, err := a.Srv().Store.Plugin().SetWithOptions(pluginID, key, value, options)
	if err != nil {
		slog.Error("Failed to set plugin key value with options", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return false, appErr
		default:
			return false, model.NewAppError("SetPluginKeyWithOptions", "app.plugin_store.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Clean up a previous entry using the hashed key, if it exists.
	if err := a.Srv().Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Warn("Failed to clean up previously hashed plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
	}

	return updated, nil
}

func (a *AppPlugin) CompareAndDeletePluginKey(pluginID string, key string, oldValue []byte) (bool, *model.AppError) {
	kv := &plugins.PluginKeyValue{
		PluginId: pluginID,
		Key:      key,
	}

	deleted, err := a.Srv().Store.Plugin().CompareAndDelete(kv, oldValue)
	if err != nil {
		slog.Error("Failed to compare and delete plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return deleted, appErr
		default:
			return false, model.NewAppError("CompareAndDeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Clean up a previous entry using the hashed key, if it exists.
	if err := a.Srv().Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Warn("Failed to clean up previously hashed plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
	}

	return deleted, nil
}

func (a *AppPlugin) GetPluginKey(pluginID string, key string) ([]byte, *model.AppError) {
	if kv, err := a.Srv().Store.Plugin().Get(pluginID, key); err == nil {
		return kv.Value, nil
	} else if nfErr := new(store.ErrNotFound); !errors.As(err, &nfErr) {
		slog.Error("Failed to query plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return nil, model.NewAppError("GetPluginKey", "app.plugin_store.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Lookup using the hashed version of the key for keys written prior to v5.6.
	if kv, err := a.Srv().Store.Plugin().Get(pluginID, getKeyHash(key)); err == nil {
		return kv.Value, nil
	} else if nfErr := new(store.ErrNotFound); !errors.As(err, &nfErr) {
		slog.Error("Failed to query plugin key value using hashed key", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return nil, model.NewAppError("GetPluginKey", "app.plugin_store.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

func (a *AppPlugin) DeletePluginKey(pluginID string, key string) *model.AppError {
	if err := a.Srv().Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Error("Failed to delete plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return model.NewAppError("DeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Also delete the key without hashing
	if err := a.Srv().Store.Plugin().Delete(pluginID, key); err != nil {
		slog.Error("Failed to delete plugin key value using hashed key", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return model.NewAppError("DeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *AppPlugin) DeleteAllKeysForPlugin(pluginID string) *model.AppError {
	if err := a.Srv().Store.Plugin().DeleteAllForPlugin(pluginID); err != nil {
		slog.Error("Failed to delete all plugin key values", slog.String("plugin_id", pluginID), slog.Err(err))
		return model.NewAppError("DeleteAllKeysForPlugin", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *AppPlugin) DeleteAllExpiredPluginKeys() *model.AppError {
	if err := a.Srv().Store.Plugin().DeleteAllExpired(); err != nil {
		slog.Error("Failed to delete all expired plugin key values", slog.Err(err))
		return model.NewAppError("DeleteAllExpiredPluginKeys", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *AppPlugin) ListPluginKeys(pluginID string, page, perPage int) ([]string, *model.AppError) {
	data, err := a.Srv().Store.Plugin().List(pluginID, page*perPage, perPage)

	if err != nil {
		slog.Error("Failed to list plugin key values", slog.Int("page", page), slog.Int("perPage", perPage), slog.Err(err))
		return nil, model.NewAppError("ListPluginKeys", "app.plugin_store.list.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return data, nil
}
