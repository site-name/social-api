package slog

import (
	"context"
)

// GraphQLLogger is used to log panics that occur during query execution.
type GraphQLLogger struct {
	logger *Logger
}

func NewGraphQLLogger(logger *Logger) *GraphQLLogger {
	return &GraphQLLogger{logger: logger}
}

// LogPanic satisfies the graphql/log.Logger interface.
// It converts the panic into an error.
func (l *GraphQLLogger) LogPanic(_ context.Context, value interface{}) {
	l.logger.Error("Error while executing GraphQL query", Any("error", value))
}
