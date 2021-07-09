package plugins

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// PluginKVSetOptions contains information on how to store a value in the plugin KV store.
type PluginKVSetOptions struct {
	Atomic          bool   // Only store the value if the current value matches the oldValue
	OldValue        []byte // The value to compare with the current value. Only used when Atomic is true
	ExpireInSeconds int64  // Set an expire counter
}

// IsValid returns nil if the chosen options are valid.
func (opt *PluginKVSetOptions) IsValid() *model.AppError {
	if !opt.Atomic && opt.OldValue != nil {
		return model.NewAppError(
			"PluginKVSetOptions.IsValid",
			"model.plugin_kvset_options.is_valid.old_value.app_error",
			nil,
			"",
			http.StatusBadRequest,
		)
	}

	return nil
}

// NewPluginKeyValueFromOptions return a PluginKeyValue given a pluginID, a KV pair and options.
func NewPluginKeyValueFromOptions(pluginId, key string, value []byte, opt PluginKVSetOptions) (*PluginKeyValue, *model.AppError) {
	expireAt := int64(0)
	if opt.ExpireInSeconds != 0 {
		expireAt = model.GetMillis() + (opt.ExpireInSeconds * 1000)
	}

	kv := &PluginKeyValue{
		PluginId: pluginId,
		Key:      key,
		Value:    value,
		ExpireAt: expireAt,
	}

	return kv, nil
}
