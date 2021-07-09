package plugin

import (
	"database/sql/driver"
)

type Helpers interface {
}

// HelpersImpl implements the helpers interface with an API that retrieves data on behalf of the plugin.
type HelpersImpl struct {
	API API
}

// ResultContainer contains the output from the LastInsertID
// and RowsAffected methods for a given set of rows.
// It is used to embed another round-trip to the server,
// and helping to avoid tracking results on the server.
type ResultContainer struct {
	LastID            int64
	LastIDError       error
	RowsAffected      int64
	RowsAffectedError error
}

// Driver is a sql driver interface that is used by plugins to perform
// raw SQL queries without opening DB connections by themselves. This interface
// is not subject to backward compatibility guarantees and is only meant to be
// used by plugins built by the Sitename team.
type Driver interface {
	// Connection
	Conn(isMaster bool) (string, error)
	ConnPing(connID string) error
	ConnClose(connID string) error
	ConnQuery(connID, q string, args []driver.NamedValue) (string, error)         // rows
	ConnExec(connID, q string, args []driver.NamedValue) (ResultContainer, error) // result

	// Transaction
	Tx(connID string, opts driver.TxOptions) (string, error)
	TxCommit(txID string) error
	TxRollback(txID string) error

	// Statement
	Stmt(connID, q string) (string, error)
	StmtClose(stID string) error
	StmtNumInput(stID string) int
	StmtQuery(stID string, args []driver.NamedValue) (string, error)         // rows
	StmtExec(stID string, args []driver.NamedValue) (ResultContainer, error) // result

	// Rows
	RowsColumns(rowsID string) []string
	RowsClose(rowsID string) error
	RowsNext(rowsID string, dest []driver.Value) error
	RowsHasNextResultSet(rowsID string) bool
	RowsNextResultSet(rowsID string) error
	RowsColumnTypeDatabaseTypeName(rowsID string, index int) string
	RowsColumnTypePrecisionScale(rowsID string, index int) (int64, int64, bool)
}
