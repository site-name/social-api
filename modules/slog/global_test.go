package slog_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitename/sitename/modules/slog"
)

func TestLoggingBeforeInitialized(t *testing.T) {
	require.NotPanics(t, func() {
		// None of these should segfault before slog is globally configured
		slog.Info("info log")
		slog.Debug("debug log")
		slog.Warn("warning log")
		slog.Error("error log")
		slog.Critical("critical log")
	})
}

func TestLoggingAfterInitialized(t *testing.T) {
	testCases := []struct {
		Description         string
		LoggerConfiguration *slog.LoggerConfiguration
		ExpectedLogs        []string
	}{
		{
			"file logging, json, debug",
			&slog.LoggerConfiguration{
				EnableConsole: false,
				EnableFile:    true,
				FileJson:      true,
				FileLevel:     slog.LevelDebug,
			},
			[]string{
				`{"level":"debug","ts":0,"caller":"slog/global_test.go:0","msg":"real debug log"}`,
				`{"level":"info","ts":0,"caller":"slog/global_test.go:0","msg":"real info log"}`,
				`{"level":"warn","ts":0,"caller":"slog/global_test.go:0","msg":"real warning log"}`,
				`{"level":"error","ts":0,"caller":"slog/global_test.go:0","msg":"real error log"}`,
				`{"level":"error","ts":0,"caller":"slog/global_test.go:0","msg":"real critical log"}`,
			},
		},
		{
			"file logging, json, error",
			&slog.LoggerConfiguration{
				EnableConsole: false,
				EnableFile:    true,
				FileJson:      true,
				FileLevel:     slog.LevelError,
			},
			[]string{
				`{"level":"error","ts":0,"caller":"slog/global_test.go:0","msg":"real error log"}`,
				`{"level":"error","ts":0,"caller":"slog/global_test.go:0","msg":"real critical log"}`,
			},
		},
		{
			"file logging, non-json, debug",
			&slog.LoggerConfiguration{
				EnableConsole: false,
				EnableFile:    true,
				FileJson:      false,
				FileLevel:     slog.LevelDebug,
			},
			[]string{
				`TIME	debug	slog/global_test.go:0	real debug log`,
				`TIME	info	slog/global_test.go:0	real info log`,
				`TIME	warn	slog/global_test.go:0	real warning log`,
				`TIME	error	slog/global_test.go:0	real error log`,
				`TIME	error	slog/global_test.go:0	real critical log`,
			},
		},
		{
			"file logging, non-json, error",
			&slog.LoggerConfiguration{
				EnableConsole: false,
				EnableFile:    true,
				FileJson:      false,
				FileLevel:     slog.LevelError,
			},
			[]string{
				`TIME	error	slog/global_test.go:0	real error log`,
				`TIME	error	slog/global_test.go:0	real critical log`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			var filePath string
			if testCase.LoggerConfiguration.EnableFile {
				tempDir, err := ioutil.TempDir(os.TempDir(), "TestLoggingAfterInitialized")
				require.NoError(t, err)
				defer os.Remove(tempDir)

				filePath = filepath.Join(tempDir, "file.log")
				testCase.LoggerConfiguration.FileLocation = filePath
			}

			logger := slog.NewLogger(testCase.LoggerConfiguration)
			slog.InitGlobalLogger(logger)

			slog.Debug("real debug log")
			slog.Info("real info log")
			slog.Warn("real warning log")
			slog.Error("real error log")
			slog.Critical("real critical log")

			if testCase.LoggerConfiguration.EnableFile {
				logs, err := ioutil.ReadFile(filePath)
				require.NoError(t, err)

				actual := strings.TrimSpace(string(logs))

				if testCase.LoggerConfiguration.FileJson {
					reTs := regexp.MustCompile(`"ts":[0-9\.]+`)
					reCaller := regexp.MustCompile(`"caller":"([^"]+):[0-9\.]+"`)
					actual = reTs.ReplaceAllString(actual, `"ts":0`)
					actual = reCaller.ReplaceAllString(actual, `"caller":"$1:0"`)
				} else {
					actualRows := strings.Split(actual, "\n")
					for i, actualRow := range actualRows {
						actualFields := strings.Split(actualRow, "\t")
						if len(actualFields) > 3 {
							actualFields[0] = "TIME"
							reCaller := regexp.MustCompile(`([^"]+):[0-9\.]+`)
							actualFields[2] = reCaller.ReplaceAllString(actualFields[2], "$1:0")
							actualRows[i] = strings.Join(actualFields, "\t")
						}
					}

					actual = strings.Join(actualRows, "\n")
				}
				require.ElementsMatch(t, testCase.ExpectedLogs, strings.Split(actual, "\n"))
			}
		})
	}
}
