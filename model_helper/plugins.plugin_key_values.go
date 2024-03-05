package model_helper

import (
	"net/http"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/null/v8"
)

const (
	KeyValuePluginIdMaxRunes = 190
	KeyValueKeyMaxRunes      = 150
)

// PluginKVSetOptions contains information on how to store a value in the plugin KV store.
type PluginKVSetOptions struct {
	Atomic          bool   // Only store the value if the current value matches the oldValue
	OldValue        []byte // The value to compare with the current value. Only used when Atomic is true
	ExpireInSeconds int64  // Set an expire counter
}

// IsValid returns nil if the chosen options are valid.
func (opt *PluginKVSetOptions) IsValid() *AppError {
	if !opt.Atomic && opt.OldValue != nil {
		return NewAppError(
			"PluginKVSetOptions.IsValid",
			"plugin_kvset_options.is_valid.old_value.app_error",
			nil,
			"",
			http.StatusBadRequest,
		)
	}

	return nil
}

func NewPluginKeyValueFromOptions(pluginId, key string, value []byte, opt PluginKVSetOptions) *model.PluginKeyValue {
	expireAt := int64(0)
	if opt.ExpireInSeconds != 0 {
		expireAt = GetMillis() + (opt.ExpireInSeconds * 1000)
	}

	kv := &model.PluginKeyValue{
		PluginID: pluginId,
		Pkey:     key,
		Pvalue: null.Bytes{
			Bytes: value,
			Valid: value != nil,
		},
		ExpireAt: model_types.NewNullInt64(expireAt),
	}

	return kv
}

func PluginKeyValueIsValid(kv model.PluginKeyValue) *AppError {
	if kv.PluginID == "" || utf8.RuneCountInString(kv.PluginID) > KeyValuePluginIdMaxRunes {
		return NewAppError("PluginKeyValue.IsValid", "model.plugin_key_value.is_valid.plugin_id.app_error", map[string]any{"Max": KeyValueKeyMaxRunes, "Min": 0}, "key="+kv.Pkey, http.StatusBadRequest)
	}

	if kv.Pkey == "" || utf8.RuneCountInString(kv.Pkey) > KeyValueKeyMaxRunes {
		return NewAppError("PluginKeyValue.IsValid", "model.plugin_key_value.is_valid.key.app_error", map[string]any{"Max": KeyValueKeyMaxRunes, "Min": 0}, "key="+kv.Pkey, http.StatusBadRequest)
	}

	return nil
}
