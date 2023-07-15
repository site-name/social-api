package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store/store_iface"
)

// SqlxCommonIface contains common methods implemented by sqlx's *DB, *Tx types.
type SqlxDBIface interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// SqlxTxCreator contains methods exclusively implemented by sqlx's *DB type.
type SqlxTxCreator interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Beginx() (*sqlx.Tx, error)
}

// namedParamRegex is used to capture all named parameters and convert them
// to lowercase. This is necessary to be able to use a single query for both
// Postgres and MySQL.
// This will also lowercase any constant strings containing a :, but sqlx
// will fail the query, so it won't be checked in inadvertently.
var namedParamRegex = regexp.MustCompile(`:\w+`)

type sqlxDBWrapper struct {
	DB           SqlxDBIface
	queryTimeout time.Duration
	trace        bool
}

// type check
var _ store_iface.SqlxExecutor = (*sqlxDBWrapper)(nil)

func newSqlxDBWrapper(db *sqlx.DB, timeout time.Duration, trace bool) *sqlxDBWrapper {
	return &sqlxDBWrapper{
		DB:           db,
		queryTimeout: timeout,
		trace:        trace,
	}
}

func (w *sqlxDBWrapper) Stats() sql.DBStats {
	if db, ok := w.DB.(*sqlx.DB); ok {
		return db.Stats()
	}
	return sql.DBStats{}
}

func (w *sqlxDBWrapper) Conn(ctx context.Context) (*sql.Conn, error) {
	if db, ok := w.DB.(*sqlx.DB); ok {
		return db.Conn(ctx)
	}
	return nil, errors.New("underlying type is not a db type")
}

func (w *sqlxDBWrapper) Commit() error {
	tx, ok := w.DB.(*sqlx.Tx)
	if ok {
		return tx.Commit()
	}

	return errors.New("the underlying type is not a transaction")
}

func (w *sqlxDBWrapper) Rollback() error {
	if tx, ok := w.DB.(*sqlx.Tx); ok {
		return tx.Rollback()
	}
	return errors.New("the underlying type is not a transaction")
}

func (w *sqlxDBWrapper) Beginx() (store_iface.SqlxExecutor, error) {
	if db, ok := w.DB.(*sqlx.DB); ok {
		tx, err := db.Beginx()
		if err != nil {
			return nil, err
		}

		return &sqlxDBWrapper{
			queryTimeout: w.queryTimeout,
			trace:        w.trace,
			DB:           tx,
		}, nil
	}

	return w, nil
}

func (w *sqlxDBWrapper) BeginXWithIsolation(opts *sql.TxOptions) (store_iface.SqlxExecutor, error) {
	if db, ok := w.DB.(*sqlx.DB); ok {
		tx, err := db.BeginTxx(context.Background(), opts)
		if err != nil {
			return nil, err
		}

		return &sqlxDBWrapper{
			queryTimeout: w.queryTimeout,
			trace:        w.trace,
			DB:           tx,
		}, nil
	}

	return nil, errors.New("already on a transaction")
}

func (w *sqlxDBWrapper) Get(dest interface{}, query string, args ...interface{}) error {
	query = w.DB.Rebind(query)
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.GetContext(ctx, dest, query, args...)
}

func (w *sqlxDBWrapper) NamedExec(query string, arg interface{}) (sql.Result, error) {
	if w.DB.DriverName() == model.DATABASE_DRIVER_POSTGRES {
		query = namedParamRegex.ReplaceAllStringFunc(query, strings.ToLower)
	}
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), arg)
		}(time.Now())
	}

	return w.DB.NamedExecContext(ctx, query, arg)
}

func (w *sqlxDBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	query = w.DB.Rebind(query)

	return w.ExecRaw(query, args...)
}

func (w *sqlxDBWrapper) ExecNoTimeout(query string, args ...interface{}) (sql.Result, error) {
	query = w.DB.Rebind(query)

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.ExecContext(context.Background(), query, args...)
}

// ExecRaw is like Exec but without any rebinding of params. You need to pass
// the exact param types of your target database.
func (w *sqlxDBWrapper) ExecRaw(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.ExecContext(ctx, query, args...)
}

func (w *sqlxDBWrapper) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	if w.DB.DriverName() == model.DATABASE_DRIVER_POSTGRES {
		query = namedParamRegex.ReplaceAllStringFunc(query, strings.ToLower)
	}
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), arg)
		}(time.Now())
	}

	switch dbOrTx := w.DB.(type) {
	case *sqlx.DB:
		return dbOrTx.NamedQueryContext(ctx, query, arg)

	default: // NOTE: only *sqlx.Tx for this case
		// There is no tx.NamedQueryContext support in the sqlx API. (https://github.com/jmoiron/sqlx/issues/447)
		// So we need to implement this ourselves.
		type result struct {
			rows *sqlx.Rows
			err  error
		}

		// Need to add a buffer of 1 to prevent goroutine leak.
		resChan := make(chan *result, 1)
		go func() {
			rows, err := dbOrTx.NamedQuery(query, arg)
			resChan <- &result{
				rows: rows,
				err:  err,
			}
		}()

		// staticcheck fails to check that res gets re-assigned later.
		res := &result{} //nolint:staticcheck
		select {
		case res = <-resChan:
		case <-ctx.Done():
			res = &result{
				rows: nil,
				err:  ctx.Err(),
			}
		}

		return res.rows, res.err
	}
}

func (w *sqlxDBWrapper) QueryRowX(query string, args ...interface{}) *sqlx.Row {
	query = w.DB.Rebind(query)
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.QueryRowxContext(ctx, query, args...)
}

func (w *sqlxDBWrapper) QueryX(query string, args ...interface{}) (*sqlx.Rows, error) {
	query = w.DB.Rebind(query)
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.QueryxContext(ctx, query, args)
}

func (w *sqlxDBWrapper) Select(dest interface{}, query string, args ...interface{}) error {
	query = w.DB.Rebind(query)
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.SelectContext(ctx, dest, query, args...)
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

func printArgs(query string, dur time.Duration, args ...interface{}) {
	query = strings.Map(removeSpace, query)
	fields := make([]slog.Field, 0, len(args)+1)
	fields = append(fields, slog.Duration("duration", dur))
	for i, arg := range args {
		fields = append(fields, slog.Any("arg"+strconv.Itoa(i), arg))
	}
	slog.Debug(query, fields...)
}
