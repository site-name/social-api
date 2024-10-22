package plugin

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

func getKeyHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (a *ServicePlugin) SetPluginKey(pluginID string, key string, value []byte) *model_helper.AppError {
	return a.SetPluginKeyWithExpiry(pluginID, key, value, 0)
}

func (a *ServicePlugin) SetPluginKeyWithExpiry(pluginID string, key string, value []byte, expireInSeconds int64) *model_helper.AppError {
	options := model.PluginKVSetOptions{
		ExpireInSeconds: expireInSeconds,
	}
	_, err := a.SetPluginKeyWithOptions(pluginID, key, value, options)
	return err
}

func (a *ServicePlugin) CompareAndSetPluginKey(pluginID string, key string, oldValue, newValue []byte) (bool, *model_helper.AppError) {
	options := model.PluginKVSetOptions{
		Atomic:   true,
		OldValue: oldValue,
	}
	return a.SetPluginKeyWithOptions(pluginID, key, newValue, options)
}

func (a *ServicePlugin) SetPluginKeyWithOptions(pluginID string, key string, value []byte, options model.PluginKVSetOptions) (bool, *model_helper.AppError) {
	if err := options.IsValid(); err != nil {
		slog.Debug("Failed to set plugin key value with options", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return false, err
	}

	updated, err := a.srv.Store.Plugin().SetWithOptions(pluginID, key, value, options)
	if err != nil {
		slog.Error("Failed to set plugin key value with options", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		var appErr *model_helper.AppError
		switch {
		case errors.As(err, &appErr):
			return false, appErr
		default:
			return false, model_helper.NewAppError("SetPluginKeyWithOptions", "app.plugin_store.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Clean up a previous entry using the hashed key, if it exists.
	if err := a.srv.Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Warn("Failed to clean up previously hashed plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
	}

	return updated, nil
}

func (a *ServicePlugin) CompareAndDeletePluginKey(pluginID string, key string, oldValue []byte) (bool, *model_helper.AppError) {
	kv := &model.PluginKeyValue{
		PluginId: pluginID,
		Key:      key,
	}

	deleted, err := a.srv.Store.Plugin().CompareAndDelete(kv, oldValue)
	if err != nil {
		slog.Error("Failed to compare and delete plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		var appErr *model_helper.AppError
		switch {
		case errors.As(err, &appErr):
			return deleted, appErr
		default:
			return false, model_helper.NewAppError("CompareAndDeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	// Clean up a previous entry using the hashed key, if it exists.
	if err := a.srv.Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Warn("Failed to clean up previously hashed plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
	}

	return deleted, nil
}

func (a *ServicePlugin) GetPluginKey(pluginID string, key string) ([]byte, *model_helper.AppError) {
	if kv, err := a.srv.Store.Plugin().Get(pluginID, key); err == nil {
		return kv.Value, nil
	} else if nfErr := new(store.ErrNotFound); !errors.As(err, &nfErr) {
		slog.Error("Failed to query plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return nil, model_helper.NewAppError("GetPluginKey", "app.plugin_store.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Lookup using the hashed version of the key for keys written prior to v5.6.
	if kv, err := a.srv.Store.Plugin().Get(pluginID, getKeyHash(key)); err == nil {
		return kv.Value, nil
	} else if nfErr := new(store.ErrNotFound); !errors.As(err, &nfErr) {
		slog.Error("Failed to query plugin key value using hashed key", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return nil, model_helper.NewAppError("GetPluginKey", "app.plugin_store.get.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

func (a *ServicePlugin) DeletePluginKey(pluginID string, key string) *model_helper.AppError {
	if err := a.srv.Store.Plugin().Delete(pluginID, getKeyHash(key)); err != nil {
		slog.Error("Failed to delete plugin key value", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return model_helper.NewAppError("DeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// Also delete the key without hashing
	if err := a.srv.Store.Plugin().Delete(pluginID, key); err != nil {
		slog.Error("Failed to delete plugin key value using hashed key", slog.String("plugin_id", pluginID), slog.String("key", key), slog.Err(err))
		return model_helper.NewAppError("DeletePluginKey", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServicePlugin) DeleteAllKeysForPlugin(pluginID string) *model_helper.AppError {
	if err := a.srv.Store.Plugin().DeleteAllForPlugin(pluginID); err != nil {
		slog.Error("Failed to delete all plugin key values", slog.String("plugin_id", pluginID), slog.Err(err))
		return model_helper.NewAppError("DeleteAllKeysForPlugin", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServicePlugin) DeleteAllExpiredPluginKeys() *model_helper.AppError {
	if err := a.srv.Store.Plugin().DeleteAllExpired(); err != nil {
		slog.Error("Failed to delete all expired plugin key values", slog.Err(err))
		return model_helper.NewAppError("DeleteAllExpiredPluginKeys", "app.plugin_store.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServicePlugin) ListPluginKeys(pluginID string, page, perPage int) ([]string, *model_helper.AppError) {
	data, err := a.srv.Store.Plugin().List(pluginID, page*perPage, perPage)

	if err != nil {
		slog.Error("Failed to list plugin key values", slog.Int("page", page), slog.Int("perPage", perPage), slog.Err(err))
		return nil, model_helper.NewAppError("ListPluginKeys", "app.plugin_store.list.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return data, nil
}
