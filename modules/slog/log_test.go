package slog_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitename/sitename/modules/slog"
)

// Test race condition when shutting down advanced logging. This test must run with the -race flag in order to verify
// that there is no race.
func TestLogger_ShutdownAdvancedLoggingRace(t *testing.T) {
	logger := slog.NewLogger(&slog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleJson:   true,
		EnableFile:    false,
		FileLevel:     slog.LevelInfo,
	})
	started := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		started <- true

		for {
			select {
			case <-ctx.Done():
				return
			default:
				logger.Debug("testing...")
			}
		}
	}()

	<-started

	err := logger.ShutdownAdvancedLogging(ctx)
	require.NoError(t, err)

	cancel()
	wg.Wait()
}
