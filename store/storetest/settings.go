package storetest

import (
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
)

type SqlStore interface {
	GetMasterX() SqlXExecutor
	DriverName() string
}

type SqlXExecutor interface {
	Get(dest any, query string, args ...any) error
	NamedExec(query string, arg any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	ExecRaw(query string, args ...any) (sql.Result, error)
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
	QueryRowX(query string, args ...any) *sqlx.Row
	QueryX(query string, args ...any) (*sqlx.Rows, error)
	Select(dest any, query string, args ...any) error
}

func log(message string) {
	verbose := false
	if verboseFlag := flag.Lookup("test.v"); verboseFlag != nil {
		verbose = verboseFlag.Value.String() != ""
	}
	if verboseFlag := flag.Lookup("v"); verboseFlag != nil {
		verbose = verboseFlag.Value.String() != ""
	}

	if verbose {
		fmt.Println(message)
	}
}

func postgreSQLDSNDatabase(dsn string) string {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	return path.Base(dsnURL.Path)
}

func postgreSQLRootDSN(dsn string) string {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}
	dsnURL.Path = "postgres"

	return dsnURL.String()
}

// execAsRoot executes the given sql as root against the testing database
func execAsRoot(settings *model.SqlSettings, sqlCommand string) error {
	dsn := postgreSQLRootDSN(*settings.DataSource)

	db, err := sql.Open(model.DATABASE_DRIVER_POSTGRES, dsn)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s database as root", model.DATABASE_DRIVER_POSTGRES)
	}
	defer db.Close()
	if _, err = db.Exec(sqlCommand); err != nil {
		return errors.Wrapf(err, "failed to execute `%s` against %s database as root", sqlCommand, model.DATABASE_DRIVER_POSTGRES)
	}

	return nil
}

func CleanupSqlSettings(settings *model.SqlSettings) {
	dbName := postgreSQLDSNDatabase(*settings.DataSource)

	if err := execAsRoot(settings, "DROP DATABASE "+dbName); err != nil {
		panic("failed to drop temporary database " + dbName + ": " + err.Error())
	}

	log("Dropped temporary database " + dbName)
}
