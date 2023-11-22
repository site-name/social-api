package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/spf13/cobra"
)

const (
	CustomDefaultsEnvVar = "SN_CUSTOM_DEFAULTS_PATH"
)

// prettyPrintStruct will return a prettyPrint version of a given struct
func prettyPrintStruct(t interface{}) string {
	return prettyPrintMap(structToMap(t))
}

// structToMap converts a struct into a map
func structToMap(t interface{}) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Panicked in structToMap. This should never happen.", slog.Any("recover", r))
		}
	}()

	val := reflect.ValueOf(t)

	if val.Kind() != reflect.Struct {
		return nil
	}

	out := map[string]interface{}{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		var value interface{}

		switch field.Kind() {
		case reflect.Struct:
			value = structToMap(field.Interface())
		case reflect.Ptr:
			indirectType := field.Elem()

			if indirectType.Kind() == reflect.Struct {
				value = structToMap(indirectType.Interface())
			} else if indirectType.Kind() != reflect.Invalid {
				value = indirectType.Interface()
			}
		default:
			value = field.Interface()
		}

		out[val.Type().Field(i).Name] = value
	}

	return out
}

// prettyPrintMap will return a prettyPrint version of a given map
func prettyPrintMap(configMap map[string]interface{}) string {
	value := reflect.ValueOf(configMap)
	return printStringMap(value, 0)
}

// printStringMap takes a reflect.Value and prints it out alphabetically based on key values, which must be strings.
// This is done recursively if it's a map, and uses the given tab settings.
func printStringMap(value reflect.Value, tabVal int) string {
	out := &bytes.Buffer{}

	var sortedKeys []string
	stringToKeyMap := make(map[string]reflect.Value)
	for _, k := range value.MapKeys() {
		sortedKeys = append(sortedKeys, k.String())
		stringToKeyMap[k.String()] = k
	}

	sort.Strings(sortedKeys)

	for _, keyString := range sortedKeys {
		key := stringToKeyMap[keyString]
		val := value.MapIndex(key)
		if newVal, ok := val.Interface().(map[string]interface{}); !ok {
			fmt.Fprintf(out, "%s", strings.Repeat("\t", tabVal))
			fmt.Fprintf(out, "%v: \"%v\"\n", key.Interface(), val.Interface())
		} else {
			fmt.Fprintf(out, "%s", strings.Repeat("\t", tabVal))
			fmt.Fprintf(out, "%v:\n", key.Interface())
			// going one level in, increase the tab
			tabVal++
			fmt.Fprintf(out, "%s", printStringMap(reflect.ValueOf(newVal), tabVal))
			// coming back one level, decrease the tab
			tabVal--
		}
	}

	return out.String()
}

func getConfigDSN(command *cobra.Command, env map[string]string) string {
	configDSN, _ := command.Flags().GetString("config")

	// Config not supplied in flag, check env
	if configDSN == "" {
		configDSN = env["SN_CONFIG"]
	}

	// Config not supplied in env or flag use default
	if configDSN == "" {
		configDSN = "config.json"
	}

	return configDSN
}

func loadCustomDefaults() (*model_helper.Config, error) {
	customDefaultsPath := os.Getenv(CustomDefaultsEnvVar)
	if customDefaultsPath == "" {
		return nil, nil
	}

	file, err := os.Open(customDefaultsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open custom defaults file at %q: %w", customDefaultsPath, err)
	}
	defer file.Close()

	var customDefaults *model_helper.Config
	err = json.NewDecoder(file).Decode(&customDefaults)
	if err != nil {
		return nil, fmt.Errorf("unable to decode custom defaults configuration: %w", err)
	}

	return customDefaults, nil
}
