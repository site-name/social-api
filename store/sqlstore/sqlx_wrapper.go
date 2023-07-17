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
	"gorm.io/gorm"
)

// SqlxCommonIface contains common methods implemented by sqlx's *DB, *Tx types.
type SqlxDBIface interface {
	BindNamed(query string, arg any) (string, []any, error)
	DriverName() string
	Get(dest any, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExec(query string, arg any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
	QueryRowx(query string, args ...any) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	Select(dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	Rebind(query string) string
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
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

func (w *sqlxDBWrapper) Get(dest any, query string, args ...any) error {
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

// func (w *sqlxDBWrapper) NamedExec(query string, arg any) (sql.Result, error) {
// 	if w.DB.DriverName() == model.DATABASE_DRIVER_POSTGRES {
// 		query = namedParamRegex.ReplaceAllStringFunc(query, strings.ToLower)
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	if w.trace {
// 		defer func(then time.Time) {
// 			printArgs(query, time.Since(then), arg)
// 		}(time.Now())
// 	}

// 	return w.DB.NamedExecContext(ctx, query, arg)
// }

func (w *sqlxDBWrapper) Exec(query string, args ...any) (sql.Result, error) {
	query = w.DB.Rebind(query)

	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.ExecContext(ctx, query, args...)
}

func (w *sqlxDBWrapper) ExecNoTimeout(query string, args ...any) (sql.Result, error) {
	query = w.DB.Rebind(query)

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.ExecContext(context.Background(), query, args...)
}

func (w *sqlxDBWrapper) NamedQuery(query string, arg any) (store_iface.RowsScanner, error) {
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

func (w *sqlxDBWrapper) QueryRowX(query string, args ...any) store_iface.Scanner {
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

func (w *sqlxDBWrapper) QueryX(query string, args ...any) (store_iface.RowsScanner, error) {
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

func (w *sqlxDBWrapper) Select(dest any, query string, args ...any) error {
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

var _ store_iface.SqlxExecutor = (*gormDBWrapper)(nil)

type gormDBWrapper struct {
	DB           *gorm.DB
	queryTimeout time.Duration
	trace        bool
}

func (g *gormDBWrapper) Select(dest any, query string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
	defer cancel()

	if g.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return g.DB.WithContext(ctx).Raw(query, args...).Find(dest).Error
}

func (g *gormDBWrapper) Get(dest any, query string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
	defer cancel()

	if g.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return g.DB.WithContext(ctx).Raw(query, args...).First(dest).Error
}

func (g *gormDBWrapper) Stats() sql.DBStats {
	sqldb, _ := g.DB.DB()
	return sqldb.Stats()
}

func (g *gormDBWrapper) Conn(ctx context.Context) (*sql.Conn, error) {
	sqldb, err := g.DB.DB()
	if err != nil {
		return nil, err
	}
	return sqldb.Conn(ctx)
}

func (g *gormDBWrapper) Commit() error {
	if g.inTransaction() {
		return g.DB.Commit().Error
	}
	return errors.New("not in a transaction")
}

func (g *gormDBWrapper) Rollback() error {
	if g.inTransaction() {
		return g.DB.Rollback().Error
	}
	return errors.New("not in a transaction")
}

func (g *gormDBWrapper) inTransaction() bool {
	switch g.DB.Statement.ConnPool.(type) {
	case gorm.TxBeginner, gorm.ConnPoolBeginner:
		return false
	default:
		return true
	}
}

func (g *gormDBWrapper) Beginx() (store_iface.SqlxExecutor, error) {
	if g.inTransaction() {
		return nil, errors.New("already in a transaction")
	}

	tx := g.DB.Begin()
	err := tx.Error
	if err != nil {
		return nil, err
	}

	return &gormDBWrapper{
		DB:           tx,
		trace:        g.trace,
		queryTimeout: g.queryTimeout,
	}, nil
}

// func (g *gormDBWrapper) NamedExec(query string, arg any) (sql.Result, error) {

// }

var r = strings.NewReplacer()

func (g *gormDBWrapper) Exec(query string, args ...any) (sql.Result, error) {
	sqldb, err := g.DB.DB()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
	defer cancel()

	if g.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return sqldb.ExecContext(ctx, query, args...)
}

func (g *gormDBWrapper) ExecNoTimeout(query string, args ...any) (sql.Result, error) {
	sqldb, err := g.DB.DB()
	if err != nil {
		return nil, err
	}

	if g.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return sqldb.ExecContext(context.Background(), query, args...)
}

func (w *gormDBWrapper) NamedQuery(query string, arg any) (store_iface.RowsScanner, error) {
	return sqlx.NamedQuery()
}

func (w *gormDBWrapper) QueryRowX(query string, args ...any) store_iface.Scanner {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.WithContext(ctx).Raw(query, args...).Row()
}

func (w *gormDBWrapper) QueryX(query string, args ...any) (store_iface.RowsScanner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
	defer cancel()

	if w.trace {
		defer func(then time.Time) {
			printArgs(query, time.Since(then), args)
		}(time.Now())
	}

	return w.DB.WithContext(ctx).Raw(query, args...).Rows()
}

func (g *gormDBWrapper) BeginXWithIsolation(opts *sql.TxOptions) (store_iface.SqlxExecutor, error) {
	if g.inTransaction() {
		return nil, errors.New("already in a transaction")
	}

	return &gormDBWrapper{
		DB:           g.DB.Begin(opts),
		trace:        g.trace,
		queryTimeout: g.queryTimeout,
	}, nil
}

func bindQueryAndArgs(query string, args []any) []any {
	res := make([]any, 0, len(args)+1)
	res[0] = query
	return append(res, args...)
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
