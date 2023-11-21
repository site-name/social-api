package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util/fileutils"
)

const (
	LogRotateSize           = 10000
	LogCompress             = true
	LogRotateMaxAge         = 0
	LogRotateMaxBackups     = 0
	LogFilename             = "sitename.log"
	LogNotificationFilename = "notifications.log"
	LogMinLevelLen          = 5
	LogMinMsgLen            = 45
	LogDelim                = " "
	LogEnableCaller         = true
)

type fileLocationFunc func(string) string

func MloggerConfigFromLoggerConfig(s *model_helper.LogSettings, configSrc LogConfigSrc, getFileFunc fileLocationFunc) (slog.LoggerConfiguration, error) {
	cfg := make(slog.LoggerConfiguration)

	var targetCfg slog.TargetCfg
	var err error

	// add the simple logging config
	if *s.EnableConsole {
		targetCfg, err = makeSimpleConsoleTarget(*s.ConsoleLevel, *s.ConsoleJson, *s.EnableColor)
		if err != nil {
			return cfg, err
		}
		cfg["_defConsole"] = targetCfg
	}

	if *s.EnableFile {
		targetCfg, err = makeSimpleFileTarget(getFileFunc(*s.FileLocation), *s.FileLevel, *s.FileJson)
		if err != nil {
			return cfg, err
		}
		cfg["_defFile"] = targetCfg
	}

	if configSrc == nil {
		return cfg, nil
	}

	// add advanced logging config
	cfgAdv := configSrc.Get()
	cfg.Append(cfgAdv)

	return cfg, nil
}

func MloggerConfigFromAuditConfig(auditSettings model_helper.ExperimentalAuditSettings, configSrc LogConfigSrc) (slog.LoggerConfiguration, error) {
	cfg := make(slog.LoggerConfiguration)

	var targetCfg slog.TargetCfg
	var err error

	// add the simple audit config
	if *auditSettings.FileEnabled {
		targetCfg, err = makeSimpleFileTarget(*auditSettings.FileName, "error", true)
		if err != nil {
			return nil, err
		}

		// apply audit specific levels
		targetCfg.Levels = []slog.Level{slog.LvlAuditAPI, slog.LvlAuditContent, slog.LvlAuditPerms, slog.LvlAuditCLI}

		// apply audit specific formatting
		targetCfg.FormatOptions = json.RawMessage(`{"disable_timestamp": true, "disable_msg": true, "disable_stacktrace": true, "disable_level": true}`)

		cfg["_defAudit"] = targetCfg
	}

	if configSrc == nil {
		return cfg, nil
	}

	// add advanced audit config
	cfgAdv := configSrc.Get()
	cfg.Append(cfgAdv)

	return cfg, nil
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = fileutils.FindDir("logs")
	}

	return filepath.Join(fileLocation, LogFilename)
}

func GetNotificationsLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = fileutils.FindDir("logs")
	}

	return filepath.Join(fileLocation, LogNotificationFilename)
}

func GetLogSettingsFromNotificationsLogSettings(notificationLogSettings *model_helper.NotificationLogSettings) *model_helper.LogSettings {
	settings := &model_helper.LogSettings{}
	settings.SetDefaults()
	settings.ConsoleJson = notificationLogSettings.ConsoleJson
	settings.ConsoleLevel = notificationLogSettings.ConsoleLevel
	settings.EnableConsole = notificationLogSettings.EnableConsole
	settings.EnableFile = notificationLogSettings.EnableFile
	settings.FileJson = notificationLogSettings.FileJson
	settings.FileLevel = notificationLogSettings.FileLevel
	settings.FileLocation = notificationLogSettings.FileLocation
	settings.AdvancedLoggingConfig = notificationLogSettings.AdvancedLoggingConfig
	settings.EnableColor = notificationLogSettings.EnableColor
	return settings
}

func makeSimpleConsoleTarget(level string, outputJSON bool, color bool) (slog.TargetCfg, error) {
	levels, err := stdLevels(level)
	if err != nil {
		return slog.TargetCfg{}, err
	}

	target := slog.TargetCfg{
		Type:         "console",
		Levels:       levels,
		Options:      json.RawMessage(`{"out": "stdout"}`),
		MaxQueueSize: 1000,
	}

	if outputJSON {
		target.Format = "json"
		target.FormatOptions = makeJSONFormatOptions()
	} else {
		target.Format = "plain"
		target.FormatOptions = makePlainFormatOptions(color)
	}
	return target, nil
}

func makeSimpleFileTarget(filename string, level string, json bool) (slog.TargetCfg, error) {
	levels, err := stdLevels(level)
	if err != nil {
		return slog.TargetCfg{}, err
	}

	target := slog.TargetCfg{
		Type:         "file",
		Levels:       levels,
		Options:      makeFileOptions(filename),
		MaxQueueSize: 1000,
	}

	if json {
		target.Format = "json"
		target.FormatOptions = makeJSONFormatOptions()
	} else {
		target.Format = "plain"
		target.FormatOptions = makePlainFormatOptions(false)
	}
	return target, nil
}

func stdLevels(level string) ([]slog.Level, error) {
	stdLevel, err := stringToStdLevel(level)
	if err != nil {
		return nil, err
	}

	var levels []slog.Level
	for _, l := range slog.StdAll {
		if l.ID <= stdLevel.ID {
			levels = append(levels, l)
		}
	}
	return levels, nil
}

func stringToStdLevel(level string) (slog.Level, error) {
	level = strings.ToLower(level)
	for _, l := range slog.StdAll {
		if l.Name == level {
			return l, nil
		}
	}
	return slog.Level{}, fmt.Errorf("%s is not a standard level", level)
}

func makeJSONFormatOptions() json.RawMessage {
	str := fmt.Sprintf(`{"enable_caller": %t}`, LogEnableCaller)
	return json.RawMessage(str)
}

func makePlainFormatOptions(enableColor bool) json.RawMessage {
	str := fmt.Sprintf(`{"delim": "%s", "min_level_len": %d, "min_msg_len": %d, "enable_color": %t, "enable_caller": %t}`,
		LogDelim, LogMinLevelLen, LogMinMsgLen, enableColor, LogEnableCaller)
	return json.RawMessage(str)
}

func makeFileOptions(filename string) json.RawMessage {
	str := fmt.Sprintf(`{"filename": "%s", "max_size": %d, "max_age": %d, "max_backups": %d, "compress": %t}`,
		filename, LogRotateSize, LogRotateMaxAge, LogRotateMaxBackups, LogCompress)
	return json.RawMessage(str)
}
