package sqlstore

import (
	"context"
	dbsql "database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/lib/pq"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/db/migrations"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/seo"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

type migrationDirection string

const (
	IndexTypeFullText              = "full_text"
	IndexTypeFullTextFunc          = "full_text_func"
	IndexTypeDefault               = "default"
	PGDupTableErrorCode            = "42P07" // see https://github.com/lib/pq/blob/master/error.go#L268
	PGForeignKeyViolationErrorCode = "23503"
	PGDuplicateObjectErrorCode     = "42710"
	DBPingAttempts                 = 18
	DBPingTimeoutSecs              = 10
	// This is a numerical version string by postgres. The format is
	// 2 characters for major, minor, and patch version prior to 10.
	// After 10, it's major and minor only.
	// 10.1 would be 100001.
	// 9.6.3 would be 90603.
	MinimumRequiredPostgresVersion = 100000

	migrationsDirectionUp   migrationDirection = "up"
	migrationsDirectionDown migrationDirection = "down"

	replicaLagPrefix = "replica-lag"
)

// os exit codes when error occures
const (
	ExitGenericFailure           = 1
	ExitCreateTable              = 100
	ExitDBOpen                   = 101
	ExitPing                     = 102
	ExitNoDriver                 = 103
	ExitTableExists              = 104
	ExitTableExistsMySQL         = 105
	ExitColumnExists             = 106
	ExitDoesColumnExistsPostgres = 107
	ExitDoesColumnExistsMySQL    = 108
	ExitDoesColumnExistsMissing  = 109
	ExitCreateColumnPostgres     = 110
	ExitCreateColumnMySQL        = 111
	ExitCreateColumnMissing      = 112
	ExitRemoveColumn             = 113
	ExitRenameColumn             = 114
	ExitMaxColumn                = 115
	ExitAlterColumn              = 116
	ExitCreateIndexPostgres      = 117
	ExitCreateIndexMySQL         = 118
	ExitCreateIndexFullMySQL     = 119
	ExitCreateIndexMissing       = 120
	ExitRemoveIndexPostgres      = 121
	ExitRemoveIndexMySQL         = 122
	ExitRemoveIndexMissing       = 123
	ExitRemoveTable              = 134
	ExitAlterPrimaryKey          = 139
)

type SqlStore struct {
	// rrCounter and srCounter should be kept first.
	// See https://github.com/mattermost/mattermost-server/v5/pull/7281
	rrCounter         int64
	srCounter         int64
	master            *gorp.DbMap
	Replicas          []*gorp.DbMap
	searchReplicas    []*gorp.DbMap
	replicaLagHandles []*dbsql.DB
	stores            SqlStoreStores
	settings          *model.SqlSettings
	lockedToMaster    bool
	context           context.Context
	metrics           einterfaces.MetricsInterface
}

type TraceOnAdapter struct{}

// ColumnInfo holds information about a column.
type ColumnInfo struct {
	DataType          string
	CharMaximumLength int
}

func (t *TraceOnAdapter) Printf(format string, v ...interface{}) {
	originalString := fmt.Sprintf(format, v...)
	newString := strings.ReplaceAll(originalString, "\n", " ")
	newString = strings.ReplaceAll(newString, "\t", " ")
	newString = strings.ReplaceAll(newString, "\"", "")
	slog.Debug(newString)
}

// New initializes connections to postgresql database
// also migrates all the sql schema using gorp
func New(settings model.SqlSettings, metrics einterfaces.MetricsInterface) *SqlStore {
	store := &SqlStore{
		rrCounter: 0,
		srCounter: 0,
		settings:  &settings,
		metrics:   metrics,
	}

	store.initConnection()

	err := store.migrate(migrationsDirectionUp)
	if err != nil {
		slog.Critical("Failed to apply database migrations.", slog.Err(err))
		os.Exit(ExitGenericFailure)
	}

	// set up tables before performing migrations:
	store.setupTables()

	// this call is actually do database migration work
	err = store.GetMaster().CreateTablesIfNotExists()

	if err != nil {
		if IsDuplicate(err) {
			slog.Warn("Duplicate key error occurred; assuming table already created and proceeding.", slog.Err(err))
		} else {
			slog.Critical("Error creating database tables.", slog.Err(err))
			os.Exit(ExitCreateTable)
		}
	}

	err = upgradeDatabase(store, model.CurrentVersion)
	if err != nil {
		slog.Critical("Failed to upgrade database.", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitGenericFailure)
	}

	// perform database indexing
	store.indexingTableFields()

	return store
}

// setupConnection opens connection to database, check if it works by ping
func setupConnection(connType string, dataSource string, settings *model.SqlSettings) *gorp.DbMap {
	db, err := dbsql.Open(*settings.DriverName, dataSource)
	if err != nil {
		slog.Critical("Failed to open SQL connection to err.", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitDBOpen)
	}

	for i := 0; i < DBPingAttempts; i++ {
		slog.Info("Pinging SQL", slog.String("database", connType))
		ctx, cancel := context.WithTimeout(context.Background(), DBPingTimeoutSecs*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err == nil {
			break
		} else {
			if i == DBPingAttempts-1 {
				slog.Critical("Failed to ping DB, server will exit.", slog.Err(err))
				time.Sleep(time.Second)
				os.Exit(ExitPing)
			} else {
				slog.Error("Failed to ping DB", slog.Err(err), slog.Int("retrying in seconds", DBPingTimeoutSecs))
				time.Sleep(DBPingTimeoutSecs * time.Second)
			}
		}
	}

	if strings.HasPrefix(connType, replicaLagPrefix) {
		// If this is a replica lag connection, we just open one connection.
		//
		// Arguably, if the query doesn't require a special credential, it does take up
		// one extra connection from the replica DB. But falling back to the replica
		// data source when the replica lag data source is null implies an ordering constraint
		// which makes things brittle and is not a good design.
		// If connections are an overhead, it is advised to use a connection pool.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxIdleConns(*settings.MaxIdleConns)
		db.SetMaxOpenConns(*settings.MaxOpenConns)
	}
	db.SetConnMaxLifetime(time.Duration(*settings.ConnMaxLifetimeMilliseconds) * time.Millisecond)
	db.SetConnMaxIdleTime(time.Duration(*settings.ConnMaxIdleTimeMilliseconds) * time.Millisecond)

	dbmap := &gorp.DbMap{
		Db:            db,
		TypeConverter: siteNameConverter{},
		Dialect:       gorp.PostgresDialect{},
		QueryTimeout:  time.Duration(*settings.QueryTimeout) * time.Second,
	}

	if settings.Trace != nil && *settings.Trace {
		dbmap.TraceOn("sql-trace:", &TraceOnAdapter{})
	}

	return dbmap
}

func (ss *SqlStore) SetContext(context context.Context) {
	ss.context = context
}

func (ss *SqlStore) Context() context.Context {
	return ss.context
}

// connect to postgresql server
func (ss *SqlStore) initConnection() {
	ss.master = setupConnection("master", *ss.settings.DataSource, ss.settings)

	if len(ss.settings.DataSourceReplicas) > 0 {
		ss.Replicas = make([]*gorp.DbMap, len(ss.settings.DataSourceReplicas))
		for i, replica := range ss.settings.DataSourceReplicas {
			ss.Replicas[i] = setupConnection(fmt.Sprintf("replica-%v", i), replica, ss.settings)
		}
	}

	if len(ss.settings.DataSourceSearchReplicas) > 0 {
		ss.searchReplicas = make([]*gorp.DbMap, len(ss.settings.DataSourceSearchReplicas))
		for i, replica := range ss.settings.DataSourceSearchReplicas {
			ss.searchReplicas[i] = setupConnection(fmt.Sprintf("search-replica-%v", i), replica, ss.settings)
		}
	}

	if len(ss.settings.ReplicaLagSettings) > 0 {
		ss.replicaLagHandles = make([]*dbsql.DB, len(ss.settings.ReplicaLagSettings))
		for i, src := range ss.settings.ReplicaLagSettings {
			if src.DataSource == nil {
				continue
			}
			gorpConn := setupConnection(fmt.Sprintf(replicaLagPrefix+"-%d", i), *src.DataSource, ss.settings)
			ss.replicaLagHandles[i] = gorpConn.Db
		}
	}
}

func (ss *SqlStore) DriverName() string {
	return *ss.settings.DriverName
}

func (ss *SqlStore) GetCurrentSchemaVersion() string {
	version, _ := ss.GetMaster().SelectStr("SELECT Value FROM Systems WHERE Name='Version'")
	return version
}

// GetDbVersion returns the version of the database being used.
// If numerical is set to true, it attempts to return a numerical version string
// that can be parsed by callers.
func (ss *SqlStore) GetDbVersion(numerical bool) (string, error) {
	var sqlVersion string

	if numerical {
		sqlVersion = `SHOW server_version_num`
	} else {
		sqlVersion = `SHOW server_version`
	}

	version, err := ss.GetReplica().SelectStr(sqlVersion)
	if err != nil {
		return "", err
	}

	return version, nil
}

func (ss *SqlStore) GetMaster() *gorp.DbMap {
	return ss.master
}

func (ss *SqlStore) GetSearchReplica() *gorp.DbMap {
	if len(ss.settings.DataSourceSearchReplicas) == 0 {
		return ss.GetReplica()
	}

	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicas))
	return ss.searchReplicas[rrNum]
}

// GetReplica try getting a slave datasource, if there is no slave or can only choose main source, it returns main datasource.
func (ss *SqlStore) GetReplica() *gorp.DbMap {
	// in case the system does not have slave data source, returns master data source instead
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster {
		return ss.GetMaster()
	}

	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.Replicas))
	return ss.Replicas[rrNum]
}

// returns number of connections to master database
func (ss *SqlStore) TotalMasterDbConnections() int {
	return ss.GetMaster().Db.Stats().OpenConnections
}

// ReplicaLagAbs queries all the replica databases to get the absolute replica lag value
// and updates the Prometheus metric with it.
func (ss *SqlStore) ReplicaLagAbs() error {
	for i, item := range ss.settings.ReplicaLagSettings {
		if item.QueryAbsoluteLag == nil || *item.QueryAbsoluteLag == "" {
			continue
		}
		var binDiff float64
		var node string
		err := ss.replicaLagHandles[i].QueryRow(*item.QueryAbsoluteLag).Scan(&node, &binDiff)
		if err != nil {
			return err
		}
		// There is no nil check needed here because it's called from the metrics store.
		ss.metrics.SetReplicaLagAbsolute(node, binDiff)
	}
	return nil
}

// ReplicaLagAbs queries all the replica databases to get the time-based replica lag value
// and updates the Prometheus metric with it.
func (ss *SqlStore) ReplicaLagTime() error {
	for i, item := range ss.settings.ReplicaLagSettings {
		if item.QueryTimeLag == nil || *item.QueryTimeLag == "" {
			continue
		}
		var timeDiff float64
		var node string
		err := ss.replicaLagHandles[i].QueryRow(*item.QueryTimeLag).Scan(&node, &timeDiff)
		if err != nil {
			return err
		}
		// There is no nil check needed here because it's called from the metrics store.
		ss.metrics.SetReplicaLagTime(node, timeDiff)
	}
	return nil
}

func (ss *SqlStore) TotalReadDbConnections() int {
	if len(ss.settings.DataSourceReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.Replicas {
		count = count + db.Db.Stats().OpenConnections
	}

	return count
}

// counts all connections to replica source
func (ss *SqlStore) TotalSearchDbConnections() int {
	if len(ss.settings.DataSourceSearchReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.searchReplicas {
		count = count + db.Db.Stats().OpenConnections
	}

	return count
}

func (ss *SqlStore) MarkSystemRanUnitTests() {
	props, err := ss.System().Get()
	if err != nil {
		return
	}

	unitTests := props[model.SYSTEM_RAN_UNIT_TESTS]
	if unitTests == "" {
		systemTests := &model.System{Name: model.SYSTEM_RAN_UNIT_TESTS, Value: "1"}
		ss.System().Save(systemTests)
	}
}

// checks if table does exist in database
func (ss *SqlStore) DoesTableExist(tableName string) bool {
	count, err := ss.GetMaster().SelectInt(
		`SELECT count(relname) FROM pg_class WHERE relname=$1`,
		strings.ToLower(tableName),
	)

	if err != nil {
		slog.Critical("Failed to check if table exists", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitTableExists)
	}

	return count > 0
}

func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
	count, err := ss.GetMaster().SelectInt(
		`SELECT COUNT(0)
			FROM   pg_attribute
			WHERE  attrelid = $1::regclass
			AND    attname = $2
			AND    NOT attisdropped`,
		strings.ToLower(tableName),
		strings.ToLower(columnName),
	)

	if err != nil {
		if err.Error() == "pq: relation \""+strings.ToLower(tableName)+"\" does not exist" {
			return false
		}

		slog.Critical("Failed to check if column exists", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitDoesColumnExistsPostgres)
	}

	return count > 0
}

// GetColumnInfo returns data type information about the given column.
func (ss *SqlStore) GetColumnInfo(tableName, columnName string) (*ColumnInfo, error) {
	var columnInfo ColumnInfo
	err := ss.GetMaster().SelectOne(&columnInfo,
		`SELECT 
			data_type as DataType,
			COALESCE(character_maximum_length, 0) as CharMaximumLength
		FROM 
			information_schema.columns
		WHERE 
			lower(table_name) = lower($1)
		AND 
			lower(column_name) = lower($2)`,
		tableName,
		columnName,
	)
	if err != nil {
		return nil, err
	}
	return &columnInfo, nil
}

// IsVarchar returns true if the column type matches one of the varchar types
// either in MySQL or PostgreSQL.
func (ss *SqlStore) IsVarchar(columnType string) bool {
	return columnType == "character varying"
}

func (ss *SqlStore) DoesTriggerExist(triggerName string) bool {
	count, err := ss.GetMaster().SelectInt(`
			SELECT
				COUNT(0)
			FROM
				pg_trigger
			WHERE
				tgname = $1
		`, triggerName)

	if err != nil {
		slog.Critical("Failed to check if trigger exists", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitGenericFailure)
	}

	return count > 0

}

func (ss *SqlStore) CreateColumnIfNotExists(tableName string, columnName string, mySqlColType string, postgresColType string, defaultValue string) bool {
	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " ADD " + columnName + " " + postgresColType + " DEFAULT '" + defaultValue + "'")
	if err != nil {
		slog.Critical("Failed to create column", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitCreateColumnPostgres)
	}

	return true
}

func (ss *SqlStore) CreateColumnIfNotExistsNoDefault(tableName string, columnName string, mySqlColType string, postgresColType string) bool {
	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " ADD " + columnName + " " + postgresColType)
	if err != nil {
		slog.Critical("Failed to create column", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitCreateColumnPostgres)
	}

	return true
}

func (ss *SqlStore) RemoveColumnIfExists(tableName string, columnName string) bool {
	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " DROP COLUMN " + columnName)
	if err != nil {
		slog.Critical("Failed to drop column", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitRemoveColumn)
	}

	return true
}

func (ss *SqlStore) RemoveTableIfExists(tableName string) bool {
	if !ss.DoesTableExist(tableName) {
		return false
	}
	_, err := ss.GetMaster().Exec("DROP TABLE " + tableName)
	if err != nil {
		slog.Critical("Failed to drop table", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitRemoveTable)
	}

	return true
}

func (ss *SqlStore) RenameColumnIfExists(tableName string, oldColumnName string, newColumnName string, colType string) bool {
	if !ss.DoesColumnExist(tableName, oldColumnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " RENAME COLUMN " + oldColumnName + " TO " + newColumnName)

	if err != nil {
		slog.Critical("Failed to rename column", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitRenameColumn)
	}

	return true
}

func (ss *SqlStore) GetMaxLengthOfColumnIfExists(tableName string, columnName string) string {
	if !ss.DoesColumnExist(tableName, columnName) {
		return ""
	}

	result, err := ss.GetMaster().SelectStr("SELECT character_maximum_length FROM information_schema.columns WHERE table_name = '" + strings.ToLower(tableName) + "' AND column_name = '" + strings.ToLower(columnName) + "'")

	if err != nil {
		slog.Critical("Failed to get max length of column", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitMaxColumn)
	}

	return result
}

func (ss *SqlStore) AlterColumnTypeIfExists(tableName string, columnName string, postgresColType string) bool {
	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + strings.ToLower(tableName) + " ALTER COLUMN " + strings.ToLower(columnName) + " TYPE " + postgresColType)

	if err != nil {
		slog.Critical("Failed to alter column type", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitAlterColumn)
	}

	return true
}

func (ss *SqlStore) AlterColumnDefaultIfExists(tableName string, columnName string, postgresColDefault *string) bool {
	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}
	// Postgres doesn't have the same limitation, but preserve the interface.
	if postgresColDefault == nil {
		return true
	}

	tableName = strings.ToLower(tableName)
	columnName = strings.ToLower(columnName)
	defaultValue := *postgresColDefault
	var err error
	if defaultValue == "" {
		_, err = ss.GetMaster().Exec("ALTER TABLE " + tableName + " ALTER COLUMN " + columnName + " DROP DEFAULT")
	} else {
		_, err = ss.GetMaster().Exec("ALTER TABLE " + tableName + " ALTER COLUMN " + columnName + " SET DEFAULT " + defaultValue)
	}

	if err != nil {
		slog.Critical("Failed to alter column", slog.String("table", tableName), slog.String("column", columnName), slog.String("default value", defaultValue), slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitGenericFailure)
		return false
	}

	return true
}

func (ss *SqlStore) AlterPrimaryKey(tableName string, columnNames []string) bool {
	// get the current primary key as a comma separated list of columns
	query := `
			SELECT string_agg(a.attname, ',') AS pk
		FROM
			pg_constraint AS c
		CROSS JOIN
			(SELECT unnest(conkey) FROM pg_constraint WHERE conrelid='` + strings.ToLower(tableName) + `'::REGCLASS AND contype='p') AS cols(colnum)
		INNER JOIN
			pg_attribute AS a ON a.attrelid = c.conrelid
		AND cols.colnum = a.attnum
		WHERE
			c.contype = 'p'
		AND c.conrelid = '` + strings.ToLower(tableName) + `'::REGCLASS`
	currentPrimaryKey, err := ss.GetMaster().SelectStr(query)
	if err != nil {
		slog.Critical("Failed to get current primary key", slog.String("table", tableName), slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitAlterPrimaryKey)
	}

	primaryKey := strings.Join(columnNames, ",")
	if strings.EqualFold(currentPrimaryKey, primaryKey) {
		return false
	}
	// alter primary key
	alterQuery := "ALTER TABLE " + tableName + " DROP CONSTRAINT " + strings.ToLower(tableName) + "_pkey, ADD PRIMARY KEY (" + strings.ToLower(primaryKey) + ")"

	_, err = ss.GetMaster().Exec(alterQuery)
	if err != nil {
		slog.Critical("Failed to alter primary key", slog.String("table", tableName), slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitAlterPrimaryKey)
	}
	return true
}

func (ss *SqlStore) CreateUniqueIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, []string{columnName}, IndexTypeDefault, true)
}

func (ss *SqlStore) CreateIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, []string{columnName}, IndexTypeDefault, false)
}

func (ss *SqlStore) CreateCompositeIndexIfNotExists(indexName string, tableName string, columnNames []string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, columnNames, IndexTypeDefault, false)
}

func (ss *SqlStore) CreateUniqueCompositeIndexIfNotExists(indexName string, tableName string, columnNames []string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, columnNames, IndexTypeDefault, true)
}

func (ss *SqlStore) CreateFullTextIndexIfNotExists(indexName string, tableName string, columnName string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, []string{columnName}, IndexTypeFullText, false)
}

func (ss *SqlStore) CreateFullTextFuncIndexIfNotExists(indexName string, tableName string, function string) bool {
	return ss.createIndexIfNotExists(indexName, tableName, []string{function}, IndexTypeFullTextFunc, false)
}

func (ss *SqlStore) createIndexIfNotExists(indexName string, tableName string, columnNames []string, indexType string, unique bool) bool {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}

	_, errExists := ss.GetMaster().SelectStr("SELECT $1::regclass", indexName)
	// It should fail if the index does not exist
	if errExists == nil {
		return false
	}

	query := ""
	if indexType == IndexTypeFullText {
		if len(columnNames) != 1 {
			slog.Critical("Unable to create multi column full text index")
			os.Exit(ExitCreateIndexPostgres)
		}
		columnName := columnNames[0]
		postgresColumnNames := convertMySQLFullTextColumnsToPostgres(columnName)
		query = "CREATE INDEX " + indexName + " ON " + tableName + " USING gin(to_tsvector('english', " + postgresColumnNames + "))"
	} else if indexType == IndexTypeFullTextFunc {
		if len(columnNames) != 1 {
			slog.Critical("Unable to create multi column full text index")
			os.Exit(ExitCreateIndexPostgres)
		}
		columnName := columnNames[0]
		query = "CREATE INDEX " + indexName + " ON " + tableName + " USING gin(to_tsvector('english', " + columnName + "))"
	} else {
		query = "CREATE " + uniqueStr + "INDEX " + indexName + " ON " + tableName + " (" + strings.Join(columnNames, ", ") + ")"
	}

	_, err := ss.GetMaster().Exec(query)
	if err != nil {
		slog.Critical("Failed to create index", slog.Err(errExists), slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitCreateIndexPostgres)
	}

	return true
}

// create foreign keys
func (ss *SqlStore) CreateForeignKeyIfNotExists(tableName, columnName, refTableName, refColumnName string, onDeleteCascade bool) (err error) {
	deleteClause := ""
	if onDeleteCascade {
		deleteClause = "ON DELETE CASCADE"
	}
	constraintName := "FK_" + tableName + "_" + refTableName
	sQuery := `ALTER TABLE ` + tableName + ` ADD CONSTRAINT ` + constraintName + ` FOREIGN KEY (` + columnName + `) REFERENCES ` + refTableName + ` (` + refColumnName + `) ` + deleteClause + `;`
	_, err = ss.GetMaster().Exec(sQuery)
	if IsConstraintAlreadyExistsError(err) {
		err = nil
	}
	if err != nil {
		slog.Warn("Could not create foreign key: " + err.Error())
	}
	return
}

func (ss *SqlStore) RemoveIndexIfExists(indexName string, tableName string) bool {
	_, err := ss.GetMaster().SelectStr("SELECT $1::regclass", indexName)
	// It should fail if the index does not exist
	if err != nil {
		return false
	}

	_, err = ss.GetMaster().Exec("DROP INDEX " + indexName)
	if err != nil {
		slog.Critical("Failed to remove index", slog.Err(err))
		time.Sleep(time.Second)
		os.Exit(ExitRemoveIndexPostgres)
	}

	return true
}

// check if given err is postgres's duplicate error
func IsConstraintAlreadyExistsError(err error) bool {
	if dbErr, ok := err.(*pq.Error); ok {
		if dbErr.Code == PGDuplicateObjectErrorCode {
			return true
		}
	}
	return false
}

// check if given err is postgres's unique error if it satisfy the following:
//
// 1) given error's type is pq.Error and err's code is "23505"(uniqueless violation code)
//
// 2) error's string contains at least one value in given indexName
func (ss *SqlStore) IsUniqueConstraintError(err error, indexName []string) bool {
	unique := false
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		unique = true
	}

	field := false
	for _, contain := range indexName {
		if strings.Contains(err.Error(), contain) {
			field = true
			break
		}
	}

	return unique && field
}

// Get all databases connections
func (ss *SqlStore) GetAllConns() []*gorp.DbMap {
	all := make([]*gorp.DbMap, len(ss.Replicas)+1)
	copy(all, ss.Replicas)
	all[len(ss.Replicas)] = ss.master
	return all
}

// RecycleDBConnections closes active connections by setting the max conn lifetime
// to d, and then resets them back to their original duration.
func (ss *SqlStore) RecycleDBConnections(d time.Duration) {
	// Get old time.
	originalDuration := time.Duration(*ss.settings.ConnMaxLifetimeMilliseconds) * time.Millisecond
	// Set the max lifetimes for all connections.
	for _, conn := range ss.GetAllConns() {
		conn.Db.SetConnMaxLifetime(d)
	}
	// Wait for that period with an additional 2 seconds of scheduling delay.
	time.Sleep(d + 2*time.Second)
	// Reset max lifetime back to original value.
	for _, conn := range ss.GetAllConns() {
		conn.Db.SetConnMaxLifetime(originalDuration)
	}
}

// close all database connections
func (ss *SqlStore) Close() {
	ss.master.Db.Close()
	for _, replica := range ss.Replicas {
		replica.Db.Close()
	}

	for _, replica := range ss.searchReplicas {
		replica.Db.Close()
	}
}

// constraint db queries to call only master replica
func (ss *SqlStore) LockToMaster() {
	ss.lockedToMaster = true
}

// let db queries free to call any db replicas
func (ss *SqlStore) UnlockFromMaster() {
	ss.lockedToMaster = false
}

func (ss *SqlStore) DropAllTables() {
	ss.master.TruncateTables()
}

func (ss *SqlStore) GetQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (ss *SqlStore) CheckIntegrity() <-chan model.IntegrityCheckResult {
	results := make(chan model.IntegrityCheckResult)
	go CheckRelationalIntegrity(ss, results)
	return results
}

// migrate performs database migration with go-migrate. Make sure to add .sql schemas to /db/migrations folder so it can work.
// If you don't add .sql schemas, still migration can be done using gorp
func (ss *SqlStore) migrate(direction migrationDirection) error {
	// When WithInstance is used in golang-migrate, the underlying driver connections are not tracked.
	// So we will have to open a fresh connection for migrations and explicitly close it when all is done.
	conn := setupConnection("migrations", *ss.settings.DataSource, ss.settings)
	defer conn.Db.Close()

	driver, err := postgres.WithInstance(conn.Db, &postgres.Config{})
	if err != nil {
		return err
	}

	var assetNamesForDriver []string
	for _, assetName := range migrations.AssetNames() {
		if strings.HasPrefix(assetName, ss.DriverName()) {
			assetNamesForDriver = append(assetNamesForDriver, filepath.Base(assetName))
		}
	}

	source := bindata.Resource(assetNamesForDriver, func(name string) ([]byte, error) {
		return migrations.Asset(filepath.Join(ss.DriverName(), name))
	})

	sourceDriver, err := bindata.WithInstance(source)
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithInstance(
		"go-bindata",
		sourceDriver,
		ss.DriverName(),
		driver)

	if err != nil {
		return err
	}
	defer migrations.Close()

	switch direction {
	case migrationsDirectionUp:
		err = migrations.Up()
	case migrationsDirectionDown:
		err = migrations.Down()
	default:
		return fmt.Errorf("unsupported migration direction %s", direction)
	}

	if err != nil && err != migrate.ErrNoChange && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

// siteNameConverter make tables able to have fields with custom types
//
// Example:
//  map[string]string, []string, map[string]interface{}, ...
type siteNameConverter struct{}

func (me siteNameConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case model.StringMap:
		return model.MapToJson(t), nil
	case account.StringMap: // this is needed
		return model.MapToJson(t), nil
	case map[string]string:
		return model.MapToJson(model.StringMap(t)), nil
	case model.StringArray:
		return model.ArrayToJson(t), nil
	case model.StringInterface:
		return model.StringInterfaceToJson(t), nil
	case map[string]interface{}:
		return model.StringInterfaceToJson(model.StringInterface(t)), nil
	case JSONSerializable:
		return t.ToJson(), nil
	case *opengraph.OpenGraph:
		return json.JSON.Marshal(t)
	}

	return val, nil
}

func (me siteNameConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {

	switch target.(type) {
	case *model.StringMap, *account.StringMap, *map[string]string:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_map"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringArray, *[]string:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_array"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringInterface, *map[string]interface{}:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_interface"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	}

	return gorp.CustomScanner{}, false
}

type JSONSerializable interface {
	ToJson() string
}

func convertMySQLFullTextColumnsToPostgres(columnNames string) string {
	columns := strings.Split(columnNames, ", ")
	concatenatedColumnNames := ""
	for i, c := range columns {
		concatenatedColumnNames += c
		if i < len(columns)-1 {
			concatenatedColumnNames += " || ' ' || "
		}
	}

	return concatenatedColumnNames
}

// IsDuplicate checks whether an error is a duplicate key error, which comes when processes are competing on creating the same
// tables in the database.
func IsDuplicate(err error) bool {
	var pqErr *pq.Error
	if errors.As(errors.Cause(err), &pqErr) {
		if pqErr.Code == PGDupTableErrorCode {
			return true
		}
	}

	return false
}

// VersionString converts an integer representation of a DB version
// to a pretty-printed string.
// Postgres doesn't follow three-part version numbers from 10.0 onwards:
// https://www.postgresql.org/docs/13/libpq-status.html#LIBPQ-PQSERVERVERSION.
func VersionString(v int) string {
	minor := v % 10000
	major := v / 10000
	return strconv.Itoa(major) + "." + strconv.Itoa(minor)
}

// indexing metadata fields for models
func (ss *SqlStore) CommonMetaDataIndex(tableName string) {
	lowerTableName := strings.ToLower(tableName)
	ss.CreateIndexIfNotExists("idx_"+lowerTableName+"_private_metadata", tableName, "PrivateMetadata")
	ss.CreateIndexIfNotExists("idx_"+lowerTableName+"_metadata", tableName, "Metadata")
}

// common method, set max size for model's sep fields
func (ss *SqlStore) CommonSeoMaxLength(table *gorp.TableMap) {
	table.ColMap("SeoTitle").SetMaxSize(seo.SEO_TITLE_MAX_LENGTH)
	table.ColMap("SeoDescription").SetMaxSize(seo.SEO_DESCRIPTION_MAX_LENGTH)
}
