package sqlstore

import (
	"context"
	dbsql "database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/lib/pq"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/db/migrations"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/store"
)

type migrationDirection string

const (
	IndexTypeFullText      = "full_text"
	IndexTypeFullTextFunc  = "full_text_func"
	IndexTypeDefault       = "default"
	PGDupTableErrorCode    = "42P07"      // see https://github.com/lib/pq/blob/master/error.go#L268
	MySQLDupTableErrorCode = uint16(1050) // see https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_table_exists_error
	DBPingAttempts         = 18
	DBPingTimeoutSecs      = 10
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

type SqlStoreStores struct {
	user            store.UserStore
	job             store.JobStore
	session         store.SessionStore
	preference      store.PreferenceStore
	system          store.SystemStore
	token           store.TokenStore
	status          store.StatusStore
	role            store.RoleStore
	userAccessToken store.UserAccessTokenStore
}

type TraceOnAdapter struct{}

func (t *TraceOnAdapter) Printf(format string, v ...interface{}) {
	originalString := fmt.Sprintf(format, v...)
	newString := strings.ReplaceAll(originalString, "\n", " ")
	newString = strings.ReplaceAll(newString, "\t", " ")
	newString = strings.ReplaceAll(newString, "\"", "")
	log.Debug(newString)
}

type mattermConverter struct{}

func (me mattermConverter) ToDb(val interface{}) (interface{}, error) {

	switch t := val.(type) {
	case model.StringMap:
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

// check if the error is postgresql's unique constraint
func IsUniqueConstraintError(err error, indexName []string) bool {
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

func (me mattermConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *model.StringMap:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("store.sql.convert_string_map")
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *map[string]string:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("store.sql.convert_string_map")
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringArray:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("store.sql.convert_string_array")
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringInterface:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("store.sql.convert_string_interface")
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *map[string]interface{}:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("store.sql.convert_string_interface")
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

type SqlStore struct {
	// rrCounter and srCounter should be kept first.
	// See https://github.com/sitename/sitename/pull/7281
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
	license           *model.License
	licenseMutex      sync.RWMutex
	metrics           einterfaces.MetricsInterface
}

// ColumnInfo holds information about a column.
type ColumnInfo struct {
	DataType          string
	CharMaximumLength int
}

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
		log.Critical("Failed to apply database migrations: %v", err)
		os.Exit(ExitGenericFailure)
	}

	store.stores.user = newSqlUserStore(store, metrics)
	err = store.GetMaster().CreateTablesIfNotExists()
	if err != nil {
		if IsDuplicate(err) {
			log.Warn("Duplicate key error occured; assuming table already created and proceeding: %v", err)
		} else {
			log.Critical("Error creating database tables: %v", err)
			os.Exit(ExitCreateTable)
		}
	}

	// err =

	store.stores.user.(*SqlUserStore).createIndexesIfNotExists()
	return store
}

func (ss *SqlStore) DriverName() string {
	return *ss.settings.DriverName
}

func (ss *SqlStore) getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar) // postgres
}

func (ss *SqlStore) CheckIntegrity() <-chan model.IntegrityCheckResult {
	results := make(chan model.IntegrityCheckResult)
	go CheckRelationalIntegrity(ss, results)
	return results
}

func (ss *SqlStore) GetAllConns() []*gorp.DbMap {
	all := make([]*gorp.DbMap, len(ss.Replicas)+1)
	copy(all, ss.Replicas)
	all[len(ss.Replicas)] = ss.master
	return all
}

// RecycleDBConnections closes active connections by setting the max conn lifetime
// to d, and then resets them back to their original duration.
func (ss *SqlStore) RecycleDBConnections(d time.Duration) {
	// Get old time
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

func (ss *SqlStore) createIndexIfNotExists(indexName, tableName string, columnNames []string, indexType string, unique bool) bool {
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
			log.Critical("Unable to create multi column full text index")
			os.Exit(ExitCreateIndexPostgres)
		}
		columnName := columnNames[0]
		postgresColumnNames := convertMySQLFullTextColumnsToPostgres(columnName)
		query = fmt.Sprintf("CREATE INDEX %s ON %s USING(gin(to_tsvector('english', %s))", indexName, tableName, postgresColumnNames)
	} else if indexType == IndexTypeFullTextFunc {
		if len(columnNames) != 1 {
			log.Critical("Unable to create multi column full text index")
			os.Exit(ExitCreateIndexPostgres)
		}
		columnName := columnNames[0]
		query = "CREATE INDEX " + indexName + " ON " + tableName + " USING gin(to_tsvector('english', " + columnName + "))"
	} else {
		query = "CREATE " + uniqueStr + "INDEX " + indexName + " ON " + tableName + " (" + strings.Join(columnNames, ", ") + ")"
	}

	_, err := ss.GetMaster().Exec(query)
	if err != nil {
		log.Critical("Failed to create index: %v, %v", errExists, err)
		time.Sleep(time.Second)
		os.Exit(ExitCreateIndexPostgres)
	}

	return true
}

// GetDbVersion returns the version of the database being used.
// If numerical is set to true, it attempts to return a numerical version string
// that can be parsed by callers.
func (ss *SqlStore) GetDbVersion(numerical bool) (string, error) {
	var sqlVersion string
	if ss.DriverName() == model.DATABASE_DRIVER_POSTGRES {
		if numerical {
			sqlVersion = `SHOW server_version_num`
		} else {
			sqlVersion = `SHOW server_version`
		}
	} else if ss.DriverName() == model.DATABASE_DRIVER_MYSQL {
		sqlVersion = `SELECT version()`
	} else {
		return "", errors.New("Not supported driver")
	}

	version, err := ss.GetReplica().SelectStr(sqlVersion)
	if err != nil {
		return "", err
	}

	return version, nil
}

func (ss *SqlStore) CreateUniqueIndexIfNotExists(indexName, tableName, columnName string) bool {
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

func (ss *SqlStore) RemoveIndexIfExists(indexName string, tableName string) bool {
	_, err := ss.GetMaster().SelectStr("SELECT $1::regclass", indexName)
	// It should fail if the index does not exist
	if err != nil {
		return false
	}

	_, err = ss.GetMaster().Exec("DROP INDEX " + indexName)
	if err != nil {
		log.Critical("Failed to remove index %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitRemoveIndexPostgres)
	}

	return true
}

func (ss *SqlStore) LockToMaster() {
	ss.lockedToMaster = true
}

func (ss *SqlStore) UnlockFromMaster() {
	ss.lockedToMaster = false
}

func (ss *SqlStore) DropAllTables() {
	ss.master.TruncateTables()
}

func (ss *SqlStore) User() store.UserStore {
	return ss.stores.user
}

func (ss *SqlStore) Job() store.JobStore {
	return ss.stores.job
}

func (ss *SqlStore) System() store.SystemStore {
	return ss.stores.system
}

func (ss *SqlStore) Session() store.SessionStore {
	return ss.stores.session
}

func (ss *SqlStore) Preference() store.PreferenceStore {
	return ss.stores.preference
}

func (ss *SqlStore) Token() store.TokenStore {
	return ss.stores.token
}

func (ss *SqlStore) Status() store.StatusStore {
	return ss.stores.status
}

func (ss *SqlStore) UserAccessToken() store.UserAccessTokenStore {
	return ss.stores.userAccessToken
}

func (ss *SqlStore) Role() store.RoleStore {
	return ss.stores.role
}

// Close databse and every replications
func (ss *SqlStore) Close() {
	ss.master.Db.Close()
	for _, replica := range ss.Replicas {
		replica.Db.Close()
	}

	for _, replica := range ss.searchReplicas {
		replica.Db.Close()
	}
}

func (ss *SqlStore) GetMaster() *gorp.DbMap {
	return ss.master
}

func (ss *SqlStore) GetSearchReplica() *gorp.DbMap {
	ss.licenseMutex.RLock()
	license := ss.license
	ss.licenseMutex.RUnlock()
	if license == nil {
		return ss.GetMaster()
	}

	if len(ss.settings.DataSourceReplicas) == 0 {
		return ss.GetReplica()
	}

	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicas))
	return ss.searchReplicas[rrNum]
}

func (ss *SqlStore) GetReplica() *gorp.DbMap {
	ss.licenseMutex.RLock()
	license := ss.license
	ss.licenseMutex.RUnlock()
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster || license == nil {
		return ss.GetMaster()
	}

	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.Replicas))
	return ss.Replicas[rrNum]
}

func (ss *SqlStore) GetCurrentSchemaVersion() string {
	version, _ := ss.GetMaster().SelectStr("SELECT Value FROM Systems WHERE Name='Version'")
	return version
}

func (ss *SqlStore) TotalMasterDbConnections() int {
	return ss.GetMaster().Db.Stats().OpenConnections
}

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

func setupConnection(connType string, dataSource string, settings *model.SqlSettings) *gorp.DbMap {
	db, err := dbsql.Open(*settings.DriverName, dataSource)
	if err != nil {
		log.Critical("Failed to open SQL connection: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitDBOpen)
	}

	for i := 0; i < DBPingAttempts; i++ {
		log.Info("Pinging SQL #%d/%d...", i+1, DBPingAttempts)
		ctx, cancel := context.WithTimeout(context.Background(), DBPingTimeoutSecs*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err == nil {
			break
		} else {
			if i == DBPingAttempts-1 {
				log.Critical("Failed to ping DB, server will exit. Err: %v", err)
				time.Sleep(time.Second)
				os.Exit(ExitPing)
			} else {
				log.Error("Failed to ping: %v", err)
				time.Sleep(DBPingTimeoutSecs * time.Second)
			}
		}
	}

	if strings.HasPrefix(connType, replicaLagPrefix) {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
	} else {
		db.SetMaxIdleConns(*settings.MaxIdleConns)
		db.SetMaxOpenConns(*settings.MaxOpenConns)
	}
	db.SetConnMaxLifetime(time.Duration(*settings.ConnMaxLifetimeMilliseconds) * time.Millisecond)

	// only go 1.15 or above support this:
	db.SetConnMaxIdleTime(time.Duration(*settings.ConnMaxIdleTimeMilliseconds) * time.Millisecond)

	var dbmap *gorp.DbMap

	if *settings.DriverName == model.DATABASE_DRIVER_POSTGRES {
		dbmap = &gorp.DbMap{
			Db:            db,
			TypeConverter: mattermConverter{},
			Dialect:       gorp.PostgresDialect{},
			// QueryTimeout: time.Duration,
		}
	} else {
		log.Critical("Failed to create dialect specific driver")
		time.Sleep(time.Second)
		os.Exit(ExitNoDriver)
	}

	// Check if need to perform database logging
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

func (ss *SqlStore) appendMultipleStatementsFlag(dataSource string) string {
	// We need to tell the MySQL driver that we want to use multiStatements
	// in order to make migrations work.
	if ss.DriverName() == model.DATABASE_DRIVER_MYSQL {
		u, err := url.Parse(dataSource)
		if err != nil {
			log.Critical("Invalid database url found: %v", err)
			os.Exit(ExitGenericFailure)
		}
		q := u.Query()
		q.Set("multiStatements", "true")
		u.RawQuery = q.Encode()
		return u.String()
	}

	return dataSource
}

func (ss *SqlStore) migrate(direction migrationDirection) error {
	var driver database.Driver
	var err error

	// When WithInstance is used in golang-migrate, the underlying driver connections are not tracked.
	// So we will have to open a fresh connection for migrations and explicitly close it when all is done.
	dataSource := ss.appendMultipleStatementsFlag(*ss.settings.DataSource)
	conn := setupConnection("migrations", dataSource, ss.settings)
	defer conn.Db.Close()

	driver, err = postgres.WithInstance(conn.Db, &postgres.Config{})
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
		driver,
	)
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
		return errors.New(fmt.Sprintf("unsupported migration direction %s", direction))
	}

	if err != nil && err != migrate.ErrNoChange && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

// DoesColumnExist check whether the column with given column name does exist in given table
func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
	count, err := ss.GetMaster().SelectInt(
		`SELECT COUNT(0)
		FROM pg_attribute
		WHERE attrelid = $1::regclass
		AND attname = $2
		AND NOT attisdropped`,
		strings.ToLower(tableName),
		strings.ToLower(columnName),
	)

	if err != nil {
		if err.Error() == "pq: relation \""+strings.ToLower(tableName)+"\" does not exist" {
			return false
		}
		log.Critical("Failed to check if column exists: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitDoesColumnExistsPostgres)
	}

	return count > 0
}

// GetColumnInfo returns data type information about the given column.
func (ss *SqlStore) GetColumnInfo(tableName, columnName string) (*ColumnInfo, error) {
	var columnInfo ColumnInfo
	err := ss.GetMaster().SelectOne(&columnInfo,
		`SELECT data_type as DataType, 
		COALESCE(character_maximum_length, 0) as CharMaximumLength
		FROM information_schema.columns
		WHERE lower(table_name) = lower($1)
		AND lower(column_name) = lower($2)`, tableName, columnName)
	if err != nil {
		return nil, err
	}
	return &columnInfo, nil
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

// IsVarchar returns true if the column type matches one of the varchar types
// either in MySQL or PostgreSQL.
func (ss *SqlStore) IsVarChar(columnType string) bool {
	return columnType == "character varying"
}

func (ss *SqlStore) DoesTriggerExist(triggerName string) bool {
	count, err := ss.GetMaster().SelectInt(`
	SELECT
		COUNT(0)
	FROM
		pg_trigger
	WHERE
		tgname = $1`, triggerName)
	if err != nil {
		log.Critical("Failed to check if trigger exists: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitGenericFailure)
	}

	return count > 0
}

func (ss *SqlStore) CreateColumnIfNotExists(tableName string, columnName string, mySqlColType string, postgresColType string, defaultValue string) bool {
	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " AND " + columnName + " " + postgresColType + " DEFAULT '" + defaultValue + "'")
	if err != nil {
		log.Critical("Failed to create column: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitCreateColumnPostgres)
	}
	return true
}

func (ss *SqlStore) CreateColumnIfNotExistsNoDefault(tableName string, columnName string, mySqlColType string, postgresColType string) bool {
	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + "AND " + columnName + " " + postgresColType)
	if err != nil {
		log.Critical("Failed to create column: %v", err)
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
		log.Critical("Failed to drop column: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitRemoveColumn)
	}

	return true
}

func (ss *SqlStore) DoesTableExist(tableName string) bool {
	count, err := ss.GetMaster().SelectInt(`SELECT count(relname) FROM pg_class WHERE relname=$1`, strings.ToLower(tableName))
	if err != nil {
		log.Critical("Failed to check if table exists: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitTableExists)
	}

	return count > 0
}

func (ss *SqlStore) GetMaxLengthOfColumnIfExists(tableName string, columnName string) string {
	if !ss.DoesColumnExist(tableName, columnName) {
		return ""
	}

	res, err := ss.GetMaster().SelectStr("SELECT character_maximum_length FROM information_schema.columns WHERE table_name = '" + strings.ToLower(tableName) + "' AND column_name = '" + strings.ToLower(columnName) + "'")
	if err != nil {
		log.Critical("Failed to get max length of column: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitMaxColumn)
	}

	return res
}

func (ss *SqlStore) AlterColumnTypeIfExists(tableName string, columnName string, mySqlColType string, postgresColType string) bool {
	if !ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	if _, err := ss.GetMaster().Exec("ALTER TABLE " + strings.ToLower(tableName) + " ALTER COLUMN " + strings.ToLower(columnName) + " TYPE " + postgresColType); err != nil {
		log.Critical("Failed to alter column type: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitAlterColumn)
	}

	return true
}

func (ss *SqlStore) RemoveTableIfExists(tableName string) bool {
	if !ss.DoesTableExist(tableName) {
		return false
	}

	_, err := ss.GetMaster().Exec("DROP TABLE " + tableName)
	if err != nil {
		log.Critical("Failed to drop table: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitRemoveTable)
	}

	return true
}

func (ss *SqlStore) RenameColumnIfExists(tableName string, oldColumnName string, newColumnName string, colType string) bool {
	if !ss.DoesColumnExist(tableName, oldColumnName) {
		return false
	}

	if _, err := ss.GetMaster().Exec("ALTER TABLE " + tableName + " RENAME COLUMN " + oldColumnName + " TO " + newColumnName); err != nil {
		log.Critical("Failed to rename column: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitRenameColumn)
	}

	return true
}

func (ss *SqlStore) AlterColumnDefaultIfExists(tableName string, columnName string, mySqlColDefault *string, postgresColDefault *string) bool {
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
		log.Critical("Failed to alter column: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitGenericFailure)
		return false
	}

	return true
}

func (ss *SqlStore) AlterPrimaryKey(tableName string, columnNames []string) bool {
	query := `
		SELECT GROUP_CONCAT(column_name ORDER BY seq_in_index) AS PK
	FROM
		information_schema.statistics
	WHERE
		table_schema = DATABASE()
	AND table_name = ?
	AND index_name = 'PRIMARY'
	GROUP BY
		index_name`
	currentPrimaryKey, err := ss.GetMaster().SelectStr(query, tableName)
	if err != nil {
		log.Critical("Failed to get current primary key: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitAlterPrimaryKey)
	}

	primaryKey := strings.Join(columnNames, ",")
	if strings.EqualFold(currentPrimaryKey, primaryKey) {
		return false
	}

	// alter primary key
	alterQuery := "ALTER TABLE " + tableName + " DROP CONSTRAINT " + strings.ToLower(tableName) + "_pkey, AND PRIMARY KEY (" + strings.ToLower(primaryKey) + ")"
	if _, err := ss.GetMaster().Exec(alterQuery); err != nil {
		log.Critical("Failed to alter primary key: %v", err)
		time.Sleep(time.Second)
		os.Exit(ExitAlterPrimaryKey)
	}
	return true
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

// VersionString converts an integer representation of a DB version
// to a pretty-printed string.
// Postgres doesn't follow three-part version numbers from 10.0 onwards:
// https://www.postgresql.org/docs/13/libpq-status.html#LIBPQ-PQSERVERVERSION.
func VersionString(v int) string {
	minor := v % 10000
	major := v / 10000
	return strconv.Itoa(major) + "." + strconv.Itoa(minor)
}
