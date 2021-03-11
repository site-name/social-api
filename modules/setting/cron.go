package setting

import "reflect"

// GetCronSettings maps the cron subsection to the provided config
func GetCronSettings(name string, config interface{}) (interface{}, error) {
	if err := Cfg.Section("cron." + name).MapTo(config); err != nil {
		return config, err
	}

	typ := reflect.TypeOf(config).Elem()
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		tpField := typ.Field(i)
		if tpField.Type.Kind() == reflect.Struct && tpField.Anonymous {
			if err := Cfg.Section("cron." + name).MapTo(field.Addr().Interface()); err != nil {
				return config, err
			}
		}
	}

	return config, nil
}
