package slog_test

import (
	"encoding/json"
	"fmt"
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
		description  string
		cfg          slog.TargetCfg
		expectedLogs []string
	}{
		{
			"file logging, json, debug",
			slog.TargetCfg{
				Type:          "file",
				Format:        "json",
				FormatOptions: json.RawMessage(`{"enable_caller":true}`),
				Levels:        []slog.Level{slog.LvlCritical, slog.LvlError, slog.LvlWarn, slog.LvlInfo, slog.LvlDebug},
			},
			[]string{
				`{"timestamp":0,"level":"debug","msg":"real debug log","caller":"slog/global_test.go:0"}`,
				`{"timestamp":0,"level":"info","msg":"real info log","caller":"slog/global_test.go:0"}`,
				`{"timestamp":0,"level":"warn","msg":"real warning log","caller":"slog/global_test.go:0"}`,
				`{"timestamp":0,"level":"error","msg":"real error log","caller":"slog/global_test.go:0"}`,
				`{"timestamp":0,"level":"critical","msg":"real critical log","caller":"slog/global_test.go:0"}`,
			},
		},
		{
			"file logging, json, error",
			slog.TargetCfg{
				Type:          "file",
				Format:        "json",
				FormatOptions: json.RawMessage(`{"enable_caller":true}`),
				Levels:        []slog.Level{slog.LvlCritical, slog.LvlError},
			},
			[]string{
				`{"timestamp":0,"level":"error","msg":"real error log","caller":"slog/global_test.go:0"}`,
				`{"timestamp":0,"level":"critical","msg":"real critical log","caller":"slog/global_test.go:0"}`,
			},
		},
		{
			"file logging, non-json, debug",
			slog.TargetCfg{
				Type:          "file",
				Format:        "plain",
				FormatOptions: json.RawMessage(`{"delim":" | ", "enable_caller":true}`),
				Levels:        []slog.Level{slog.LvlCritical, slog.LvlError, slog.LvlWarn, slog.LvlInfo, slog.LvlDebug},
			},
			[]string{
				`debug | TIME | real debug log | caller="slog/global_test.go:0"`,
				`info | TIME | real info log | caller="slog/global_test.go:0"`,
				`warn | TIME | real warning log | caller="slog/global_test.go:0"`,
				`error | TIME | real error log | caller="slog/global_test.go:0"`,
				`critical | TIME | real critical log | caller="slog/global_test.go:0"`,
			},
		},
		{
			"file logging, non-json, error",
			slog.TargetCfg{
				Type:          "file",
				Format:        "plain",
				FormatOptions: json.RawMessage(`{"delim":" | ", "enable_caller":true}`),
				Levels:        []slog.Level{slog.LvlCritical, slog.LvlError},
			},
			[]string{
				`error | TIME | real error log | caller="slog/global_test.go:0"`,
				`critical | TIME | real critical log | caller="slog/global_test.go:0"`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var filePath string
			if testCase.cfg.Type == "file" {
				tempDir, err := ioutil.TempDir(os.TempDir(), "TestLoggingAfterInitialized")
				require.NoError(t, err)
				defer os.Remove(tempDir)

				filePath = filepath.Join(tempDir, "file.log")
				testCase.cfg.Options = json.RawMessage(fmt.Sprintf(`{"filename": "%s"}`, filePath))
			}

			logger, _ := slog.NewLogger()
			err := logger.ConfigureTargets(map[string]slog.TargetCfg{testCase.description: testCase.cfg}, nil)
			require.NoError(t, err)

			slog.InitGlobalLogger(logger)

			slog.Debug("real debug log")
			slog.Info("real info log")
			slog.Warn("real warning log")
			slog.Error("real error log")
			slog.Critical("real critical log")

			logger.Shutdown()

			if testCase.cfg.Type == "file" {
				logs, err := ioutil.ReadFile(filePath)
				require.NoError(t, err)

				actual := strings.TrimSpace(string(logs))

				if testCase.cfg.Format == "json" {
					reTs := regexp.MustCompile(`"timestamp":"[0-9\.\-\:\sZ]+"`)
					reCaller := regexp.MustCompile(`"caller":"([^"]+):[0-9\.]+"`)
					actual = reTs.ReplaceAllString(actual, `"timestamp":0`)
					actual = reCaller.ReplaceAllString(actual, `"caller":"$1:0"`)
				} else {
					reTs := regexp.MustCompile(`\[\d\d\d\d-\d\d-\d\d\s[0-9\:\.\s\-Z]+\]`)
					reCaller := regexp.MustCompile(`caller="([^"]+):[0-9\.]+"`)
					actual = reTs.ReplaceAllString(actual, "TIME")
					actual = reCaller.ReplaceAllString(actual, `caller="$1:0"`)
				}
				require.ElementsMatch(t, testCase.expectedLogs, strings.Split(actual, "\n"))
			}
		})
	}
}
