package util

import (
	"path/filepath"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util/fileutils"
)

const (
	LogRotateSize           = 10000
	LogFilename             = "sitename.log"
	LogNotificationFilename = "notifications.log"
)

type fileLocationFunc func(string) string

func MloggerConfigFromLoggerConfig(s *model.LogSettings, getFileFunc fileLocationFunc) *slog.LoggerConfiguration {
	return &slog.LoggerConfiguration{
		EnableConsole: *s.EnableConsole,
		ConsoleJson:   *s.ConsoleJson,
		ConsoleLevel:  strings.ToLower(*s.ConsoleLevel),
		EnableFile:    *s.EnableFile,
		FileJson:      *s.FileJson,
		FileLevel:     strings.ToLower(*s.FileLevel),
		FileLocation:  getFileFunc(*s.FileLocation),
		EnableColor:   *s.EnableColor,
	}
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

func GetLogSettingsFromNotificationsLogSettings(notificationLogSettings *model.NotificationLogSettings) *model.LogSettings {
	return &model.LogSettings{
		ConsoleJson:           notificationLogSettings.ConsoleJson,
		ConsoleLevel:          notificationLogSettings.ConsoleLevel,
		EnableConsole:         notificationLogSettings.EnableConsole,
		EnableFile:            notificationLogSettings.EnableFile,
		FileJson:              notificationLogSettings.FileJson,
		FileLevel:             notificationLogSettings.FileLevel,
		FileLocation:          notificationLogSettings.FileLocation,
		AdvancedLoggingConfig: notificationLogSettings.AdvancedLoggingConfig,
		EnableColor:           notificationLogSettings.EnableColor,
	}
}

// DON'T USE THIS Modify the level on the app logger
func DisableDebugLogForTest() {
	slog.GloballyDisableDebugLogForTest()
}

// DON'T USE THIS Modify the level on the app logger
func EnableDebugLogForTest() {
	slog.GloballyEnableDebugLogForTest()
}
