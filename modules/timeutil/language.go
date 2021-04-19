package timeutil

import (
	"time"

	"github.com/sitename/sitename/modules/setting"
)

var (
	langTimeFormats = map[string]string{
		"zh-CN": "2006年01月02日 15时04分05秒",
		"en-US": time.RFC1123,
		"lv-LV": "02.01.2006. 15:04:05",
	}
)

// GetLangTimeFormat represents the default time format for the language
func GetLangTimeFormat(lang string) string {
	return langTimeFormats[lang]
}

// GetTimeFormat represents the
func GetTimeFormat(lang string) string {
	if setting.TimeFormat == "" {
		format := GetLangTimeFormat(lang)
		if format == "" {
			format = time.RFC1123
		}
		return format
	}
	return setting.TimeFormat
}
