package config

import (
	"fmt"

	"github.com/sitename/sitename/modules/slog"
)

type splitLogger struct {
	wrappedLog *slog.Logger
}

func (s *splitLogger) Error(msg ...interface{}) {
	s.wrappedLog.Error(fmt.Sprint(msg...))
}

func (s *splitLogger) Warning(msg ...interface{}) {
	s.wrappedLog.Warn(fmt.Sprint(msg...))
}

// Ignoring more verbose messages from split
func (s *splitLogger) Info(msg ...interface{}) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}

func (s *splitLogger) Debug(msg ...interface{}) {
	//s.wrappedLog.Debug(fmt.Sprint(msg...))
}

func (s *splitLogger) Verbose(msg ...interface{}) {
	//s.wrappedLog.Info(fmt.Sprint(msg...))
}