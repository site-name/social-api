package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// marshalConfig converts the given configuration into JSON bytes for persistence.
func marshalConfig(cfg *model_helper.Config) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "    ")
}

// desanitize replaces fake settings with their actual values.
func desanitize(actual, target *model_helper.Config) {
	if target.LdapSettings.BindPassword != nil && *target.LdapSettings.BindPassword == model_helper.FAKE_SETTING {
		*target.LdapSettings.BindPassword = *actual.LdapSettings.BindPassword
	}

	if *target.FileSettings.PublicLinkSalt == model_helper.FAKE_SETTING {
		*target.FileSettings.PublicLinkSalt = *actual.FileSettings.PublicLinkSalt
	}
	if *target.FileSettings.AmazonS3SecretAccessKey == model_helper.FAKE_SETTING {
		target.FileSettings.AmazonS3SecretAccessKey = actual.FileSettings.AmazonS3SecretAccessKey
	}

	if *target.EmailSettings.SMTPPassword == model_helper.FAKE_SETTING {
		target.EmailSettings.SMTPPassword = actual.EmailSettings.SMTPPassword
	}

	if *target.GitLabSettings.Secret == model_helper.FAKE_SETTING {
		target.GitLabSettings.Secret = actual.GitLabSettings.Secret
	}

	if target.GoogleSettings.Secret != nil && *target.GoogleSettings.Secret == model_helper.FAKE_SETTING {
		target.GoogleSettings.Secret = actual.GoogleSettings.Secret
	}

	// if target.Office365Settings.Secret != nil && *target.Office365Settings.Secret == model_helper.FAKE_SETTING {
	// 	target.Office365Settings.Secret = actual.Office365Settings.Secret
	// }

	if target.OpenIdSettings.Secret != nil && *target.OpenIdSettings.Secret == model_helper.FAKE_SETTING {
		target.OpenIdSettings.Secret = actual.OpenIdSettings.Secret
	}

	if *target.SqlSettings.DataSource == model_helper.FAKE_SETTING {
		*target.SqlSettings.DataSource = *actual.SqlSettings.DataSource
	}
	if *target.SqlSettings.AtRestEncryptKey == model_helper.FAKE_SETTING {
		target.SqlSettings.AtRestEncryptKey = actual.SqlSettings.AtRestEncryptKey
	}

	if *target.ElasticsearchSettings.Password == model_helper.FAKE_SETTING {
		*target.ElasticsearchSettings.Password = *actual.ElasticsearchSettings.Password
	}

	if len(target.SqlSettings.DataSourceReplicas) == len(actual.SqlSettings.DataSourceReplicas) {
		for i, value := range target.SqlSettings.DataSourceReplicas {
			if value == model_helper.FAKE_SETTING {
				target.SqlSettings.DataSourceReplicas[i] = actual.SqlSettings.DataSourceReplicas[i]
			}
		}
	}

	if len(target.SqlSettings.DataSourceSearchReplicas) == len(actual.SqlSettings.DataSourceSearchReplicas) {
		for i, value := range target.SqlSettings.DataSourceSearchReplicas {
			if value == model_helper.FAKE_SETTING {
				target.SqlSettings.DataSourceSearchReplicas[i] = actual.SqlSettings.DataSourceSearchReplicas[i]
			}
		}
	}

	if *target.MessageExportSettings.GlobalRelaySettings.SmtpPassword == model_helper.FAKE_SETTING {
		*target.MessageExportSettings.GlobalRelaySettings.SmtpPassword = *actual.MessageExportSettings.GlobalRelaySettings.SmtpPassword
	}

	if target.ServiceSettings.GfycatApiSecret != nil && *target.ServiceSettings.GfycatApiSecret == model_helper.FAKE_SETTING {
		*target.ServiceSettings.GfycatApiSecret = *actual.ServiceSettings.GfycatApiSecret
	}

	if *target.ServiceSettings.SplitKey == model_helper.FAKE_SETTING {
		*target.ServiceSettings.SplitKey = *actual.ServiceSettings.SplitKey
	}
}

// fixConfig patches invalid or missing data in the configuration.
func fixConfig(cfg *model_helper.Config) {
	// Ensure SiteURL has no trailing slash.
	if strings.HasSuffix(*cfg.ServiceSettings.SiteURL, "/") {
		*cfg.ServiceSettings.SiteURL = strings.TrimRight(*cfg.ServiceSettings.SiteURL, "/")
	}

	// Ensure the directory for a local file store has a trailing slash.
	if *cfg.FileSettings.DriverName == model_helper.IMAGE_DRIVER_LOCAL {
		if *cfg.FileSettings.Directory != "" && !strings.HasSuffix(*cfg.FileSettings.Directory, "/") {
			*cfg.FileSettings.Directory += "/"
		}
	}

	FixInvalidLocales(cfg)
}

// FixInvalidLocales checks and corrects the given config for invalid locale-related settings.
//
// Ideally, this function would be completely internal, but it's currently exposed to allow the cli
// to test the config change before allowing the save.
func FixInvalidLocales(cfg *model_helper.Config) bool {
	var changed bool

	locales := i18n.GetSupportedLocales()
	if _, ok := locales[cfg.LocalizationSettings.DefaultServerLocale.String()]; !ok {
		*cfg.LocalizationSettings.DefaultServerLocale = model_helper.DEFAULT_LOCALE
		slog.Warn("DefaultServerLocale must be one of the supported locales. Setting DefaultServerLocale to en as default value.")
		changed = true
	}

	if _, ok := locales[cfg.LocalizationSettings.DefaultClientLocale.String()]; !ok {
		*cfg.LocalizationSettings.DefaultClientLocale = model_helper.DEFAULT_LOCALE
		slog.Warn("DefaultClientLocale must be one of the supported locales. Setting DefaultClientLocale to en as default value.")
		changed = true
	}

	if *cfg.LocalizationSettings.AvailableLocales != "" {
		isDefaultClientLocaleInAvailableLocales := false
		for _, word := range strings.Split(*cfg.LocalizationSettings.AvailableLocales, ",") {
			if _, ok := locales[word]; !ok {
				*cfg.LocalizationSettings.AvailableLocales = ""
				isDefaultClientLocaleInAvailableLocales = true
				slog.Warn("AvailableLocales must include DefaultClientLocale. Setting AvailableLocales to all locales as default value.")
				changed = true
				break
			}

			if word == cfg.LocalizationSettings.DefaultClientLocale.String() {
				isDefaultClientLocaleInAvailableLocales = true
			}
		}

		availableLocales := *cfg.LocalizationSettings.AvailableLocales

		if !isDefaultClientLocaleInAvailableLocales {
			availableLocales += "," + cfg.LocalizationSettings.DefaultClientLocale.String()
			slog.Warn("Adding DefaultClientLocale to AvailableLocales.")
			changed = true
		}

		*cfg.LocalizationSettings.AvailableLocales = util.AnyArray[string](strings.Fields(availableLocales)).Dedup().Join(",")
	}

	return changed
}

// Merge merges two configs together. The receiver's values are overwritten with the patch's
// values except when the patch's values are nil.
func Merge(cfg *model_helper.Config, patch *model_helper.Config, mergeConfig *util.MergeConfig) (*model_helper.Config, error) {
	ret, err := util.Merge(cfg, patch, mergeConfig)
	if err != nil {
		return nil, err
	}

	retCfg := ret.(model_helper.Config)
	return &retCfg, nil
}

func IsDatabaseDSN(dsn string) bool {
	return strings.HasPrefix(dsn, "mysql://") || strings.HasPrefix(dsn, "postgres://")
}

// stripPassword remove the password from a given DSN
func stripPassword(dsn, schema string) string {
	prefix := schema + "://"
	dsn = strings.TrimPrefix(dsn, prefix)

	i := strings.Index(dsn, ":")
	j := strings.LastIndex(dsn, "@")

	// Return error if no @ sign is found
	if j < 0 {
		return "(omitted due to error parsing the DSN)"
	}

	// Return back the input if no password is found
	if i < 0 || i > j {
		return prefix + dsn
	}

	return prefix + dsn[:i+1] + dsn[j:]
}

func isJSONMap(data string) bool {
	var m map[string]any
	return json.Unmarshal([]byte(data), &m) == nil
}

func GetValueByPath(path []string, obj any) (any, bool) {
	r := reflect.ValueOf(obj)
	var val reflect.Value
	if r.Kind() == reflect.Map {
		val = r.MapIndex(reflect.ValueOf(path[0]))
		if val.IsValid() {
			val = val.Elem()
		}
	} else {
		val = r.FieldByName(path[0])
	}

	if !val.IsValid() {
		return nil, false
	}

	switch {
	case len(path) == 1:
		return val.Interface(), true
	case val.Kind() == reflect.Struct:
		return GetValueByPath(path[1:], val.Interface())
	case val.Kind() == reflect.Map:
		remainingPath := strings.Join(path[1:], ".")
		mapIter := val.MapRange()
		for mapIter.Next() {
			key := mapIter.Key().String()
			if strings.HasPrefix(remainingPath, key) {
				i := strings.Count(key, ".") + 2 // number of dots + a dot on each side
				mapVal := mapIter.Value()
				// if no sub field path specified, return the object
				if len(path[i:]) == 0 {
					return mapVal.Interface(), true
				}
				data := mapVal.Interface()
				if mapVal.Kind() == reflect.Ptr {
					data = mapVal.Elem().Interface() // if value is a pointer, dereference it
				}
				// pass subpath
				return GetValueByPath(path[i:], data)
			}
		}
	}
	return nil, false
}

func equal(oldCfg, newCfg *model_helper.Config) (bool, error) {
	oldCfgBytes, err := json.Marshal(oldCfg)
	if err != nil {
		return false, fmt.Errorf("failed to marshal old config: %w", err)
	}
	newCfgBytes, err := json.Marshal(newCfg)
	if err != nil {
		return false, fmt.Errorf("failed to marshal new config: %w", err)
	}
	return !bytes.Equal(oldCfgBytes, newCfgBytes), nil
}
