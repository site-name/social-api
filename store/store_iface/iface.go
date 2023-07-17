package store_iface

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// SqlxExecutor exposes sqlx operations. It is used to enable some internal store methods to
// accept both transactions (*sqlxTxWrapper) and common db handlers (*sqlxDbWrapper).
type SqlxExecutor interface {
	Get(dest any, query string, args ...any) error
	// NamedExec(query string, arg any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	// NamedQuery(query string, arg any) (RowsScanner, error)
	QueryRowX(query string, args ...any) Scanner
	QueryX(query string, args ...any) (RowsScanner, error)
	Select(dest any, query string, args ...any) error
	ExecNoTimeout(query string, args ...any) (sql.Result, error)
	Beginx() (SqlxExecutor, error)
	Conn(ctx context.Context) (*sql.Conn, error)
	driver.Tx
}

type Scanner interface {
	Scan(dest ...any) error
	Err() error
}

type RowsScanner interface {
	Scanner
	Close() error
	Next() bool
}
