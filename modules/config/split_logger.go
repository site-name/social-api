package config

import (
	"fmt"

	"github.com/sitename/sitename/modules/slog"
)

type splitLogger struct {
	wrappedLog *slog.Logger
}

func (s *splitLogger) Error(msg ...any) {
	s.wrappedLog.Error(fmt.Sprint(msg...))
}

func (s *splitLogger) Warning(msg ...any) {
	s.wrappedLog.Warn(fmt.Sprint(msg...))
}

// Ignoring more verbose messages from split
func (s *splitLogger) Info(msg ...any) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}

func (s *splitLogger) Debug(msg ...any) {
	//s.wrappedLog.Debug(fmt.Sprint(msg...))
}

func (s *splitLogger) Verbose(msg ...any) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}
