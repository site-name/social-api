package store_iface

import (
	"database/sql"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
)

type Builder interface {
	ToSql() (string, []interface{}, error)
}

// SqlxExecutor exposes sqlx operations. It is used to enable some internal store methods to
// accept both transactions (*sqlxTxWrapper) and common db handlers (*sqlxDbWrapper).
type SqlxExecutor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetBuilder(dest interface{}, builder Builder) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecBuilder(builder Builder) (sql.Result, error)
	ExecRaw(query string, args ...interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	QueryRowX(query string, args ...interface{}) *sqlx.Row
	QueryX(query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	SelectBuilder(dest interface{}, builder Builder) error
	ExecNoTimeout(query string, args ...interface{}) (sql.Result, error)
	Beginx() (SqlxTxExecutor, error)
}

type SqlxTxExecutor interface {
	SqlxExecutor
	driver.Tx
}
