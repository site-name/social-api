package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

var _ store.ContextRunner = (*sqlDBWrapper)(nil)
var _ sqlDBInterface = (*sql.DB)(nil)
var _ sqlDBInterface = (*sql.Tx)(nil)

type sqlDBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type sqlDBWrapper struct {
	sqlDBInterface

	queryTimeout time.Duration
	trace        bool
}

func newSqlDbWrapper(db *sql.DB, timeout time.Duration, trace bool) *sqlDBWrapper {
	return &sqlDBWrapper{
		sqlDBInterface: db,
		queryTimeout:   timeout,
		trace:          trace,
	}
}

func (w *sqlDBWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (store.ContextRunner, error) {
	switch t := w.sqlDBInterface.(type) {
	case *sql.Tx:
		return w, nil
	case *sql.DB:
		tx, err := t.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		res := &sqlDBWrapper{
			sqlDBInterface: tx,
			queryTimeout:   w.queryTimeout,
			trace:          w.trace,
		}
		return res, nil

	default:
		return nil, fmt.Errorf("upsupported type %T", t)
	}
}

func (w *sqlDBWrapper) Commit() error {
	switch t := w.sqlDBInterface.(type) {
	case *sql.DB:
		return nil
	case *sql.Tx:
		return t.Commit()
	default:
		return fmt.Errorf("upsupported type %T", t)
	}
}

func (w *sqlDBWrapper) Rollback() error {
	switch t := w.sqlDBInterface.(type) {
	case *sql.DB:
		return nil
	case *sql.Tx:
		return t.Commit()
	default:
		return fmt.Errorf("upsupported type %T", t)
	}
}

// Exec implements boil.ContextExecutor.
func (w *sqlDBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.ExecContext(ctx, query, args...)
}

// ExecContext implements boil.ContextExecutor.
func (w *sqlDBWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.ExecContext(ctx, query, args...)
}

// Query implements boil.ContextExecutor.
func (w *sqlDBWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.QueryContext(ctx, query, args...)
}

// QueryContext implements boil.ContextExecutor.
func (w *sqlDBWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.QueryContext(ctx, query, args...)
}

// QueryRow implements boil.ContextExecutor.
func (w *sqlDBWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.QueryRowContext(ctx, query, args...)
}

// QueryRowContext implements boil.ContextExecutor.
func (w *sqlDBWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(ctx, w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.sqlDBInterface.QueryRowContext(ctx, query, args...)
}

func removeSpace(r rune) rune {
	// Strip everything except ' '
	// This also strips out more than one space,
	// but we ignore it for now until someone complains.
	if unicode.IsSpace(r) && r != ' ' {
		return -1
	}
	return r
}

func printArgs(query string, dur time.Duration, args ...any) {
	query = strings.Map(removeSpace, query)
	fields := make([]slog.Field, 0, len(args)+1)
	fields = append(fields, slog.Duration("duration", dur))
	for i, arg := range args {
		fields = append(fields, slog.Any("arg"+strconv.Itoa(i), arg))
	}
	slog.Debug(query, fields...)
}
