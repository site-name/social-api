package storetest

import (
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
)

const (
	defaultPostgresqlDSN = "postgres://minh:anhyeuem98@localhost:5432/sitename_test?sslmode=disable&connect_timeout=10"
)

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
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

func getDefaultPostgresqlDSN() string {
	if os.Getenv("IS_CI") == "true" {
		return strings.ReplaceAll(defaultPostgresqlDSN, "localhost", "postgres")
	}
	return defaultPostgresqlDSN
}

// PostgresSQLSettings returns the database settings to connect to the PostgreSQL unittesting database.
// The database name is generated randomly and must be created before use.
func PostgreSQLSettings() *model.SqlSettings {
	dsn := os.Getenv("TEST_DATABASE_POSTGRESQL_DSN")
	if dsn == "" {
		dsn = getDefaultPostgresqlDSN()
		slog.Info("No TEST_DATABASE_POSTGRESQL_DSN override, using default", slog.String("default_dsn", dsn))
	} else {
		slog.Info("Using TEST_DATABASE_POSTGRESQL_DSN override", slog.String("dsn", dsn))
	}

	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	// Generate a random database name
	dsnURL.Path = "db" + model.NewRandomString(26)

	return databaseSettings("postgres", dsnURL.String())
}

func databaseSettings(driver, dataSource string) *model.SqlSettings {
	settings := &model.SqlSettings{
		DriverName:                        &driver,
		DataSource:                        &dataSource,
		DataSourceReplicas:                []string{},
		DataSourceSearchReplicas:          []string{},
		MaxIdleConns:                      new(int),
		ConnMaxLifetimeMilliseconds:       new(int),
		ConnMaxIdleTimeMilliseconds:       new(int),
		MaxOpenConns:                      new(int),
		Trace:                             model.GetPointerOfValue(false),
		AtRestEncryptKey:                  model.GetPointerOfValue(model.NewRandomString(32)),
		QueryTimeout:                      new(int),
		MigrationsStatementTimeoutSeconds: new(int),
	}
	*settings.MaxIdleConns = 10
	*settings.ConnMaxLifetimeMilliseconds = 3600000
	*settings.ConnMaxIdleTimeMilliseconds = 300000
	*settings.MaxOpenConns = 100
	*settings.QueryTimeout = 60
	*settings.MigrationsStatementTimeoutSeconds = 60

	return settings
}

func postgreSQLRootDSN(dsn string) string {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}
	// dsnUrl.User = url.UserPassword("", password)
	dsnURL.Path = "postgres"

	return dsnURL.String()
}

func execAsRoot(settings *model.SqlSettings, sqlCommand string) error {
	dsn := postgreSQLRootDSN(*settings.DataSource)
	db, err := sql.Open(*settings.DriverName, dsn)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s database as root", *settings.DriverName)
	}
	defer db.Close()
	if _, err = db.Exec(sqlCommand); err != nil {
		return errors.Wrapf(err, "failed to execute `%s` against %s database as root", sqlCommand, *settings.DriverName)
	}

	return nil
}

func postgreSQLDSNDatabase(dsn string) string {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	return path.Base(dsnURL.Path)
}

func MakeSqlSettings(driver string, withReplica bool) *model.SqlSettings {
	settings := PostgreSQLSettings()
	dbName := postgreSQLDSNDatabase(*settings.DataSource)

	if err := execAsRoot(settings, "CREATE DATABASE "+dbName); err != nil {
		panic("failed to create temporary database " + dbName + ": " + err.Error())
	}

	if err := execAsRoot(settings, "GRANT ALL PRIVILEGES ON DATABASE \""+dbName+"\" TO minh"); err != nil {
		panic("failed to grant minh permission to " + dbName + ":" + err.Error())
	}

	log("Created temporary " + driver + " database " + dbName)
	settings.ReplicaMonitorIntervalSeconds = model.GetPointerOfValue(5)

	return settings
}

func CleanupSqlSettings(settings *model.SqlSettings) {
	dbName := postgreSQLDSNDatabase(*settings.DataSource)
	if err := execAsRoot(settings, "DROP DATABASE "+dbName); err != nil {
		panic("failed to drop temporary database " + dbName + ": " + err.Error())
	}

	log("Dropped temporary database " + dbName)
}
