package sqlstore

import (
	"context"
	"database/sql"
	dbsql "database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	ps "github.com/mattermost/morph/drivers/postgres"
	mbindata "github.com/mattermost/morph/sources/embedded"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/db"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store/store_iface"
)

type migrationDirection string

const (
	IndexTypeFullText              = "full_text"
	IndexTypeFullTextFunc          = "full_text_func"
	IndexTypeDefault               = "default"
	PGDupTableErrorCode            = "42P07"      // see https://github.com/lib/pq/blob/master/error.go#L268
	MySQLDupTableErrorCode         = uint16(1050) // see https://dev.mysql.com/doc/mysql-errors/5.7/en/server-error-reference.html#error_er_table_exists_error
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

type SqlStore struct {
	// rrCounter and srCounter should be kept first.
	rrCounter int64
	srCounter int64

	// master  *gorp.DbMap
	masterX *sqlxDBWrapper

	// Replicas  []*gorp.DbMap
	ReplicaXs []*sqlxDBWrapper

	// searchReplicas  []*gorp.DbMap
	searchReplicaXs []*sqlxDBWrapper

	replicaLagHandles []*dbsql.DB
	stores            *SqlStoreStores
	settings          *model.SqlSettings
	lockedToMaster    bool
	context           context.Context
	metrics           einterfaces.MetricsInterface
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

	// check if database version is met requirement
	ver, err := store.GetDbVersion(true)
	if err != nil {
		slog.Fatal("Cannot get DB version.", slog.Err(err))
	}

	ok, err := store.ensureMinimumDBVersion(ver)
	if !ok {
		slog.Fatal("Error while checking DB version.", slog.Err(err))
	}

	// migrate tables
	err = store.migrate(migrationsDirectionUp)
	if err != nil {
		slog.Fatal("Failed to apply database migrations.", slog.Err(err))
	}

	// set up tables before performing migrations:
	store.setupTables()
	return store
}

// setupConnection opens connection to database, check if it works by ping
func setupConnection(connType string, dataSource string, settings *model.SqlSettings) *dbsql.DB {
	db, err := dbsql.Open(*settings.DriverName, dataSource)
	if err != nil {
		slog.Fatal("Failed to open SQL connection to err.", slog.Err(err))
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
				slog.Fatal("Failed to ping DB, server will exit.", slog.Err(err))
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

	return db
}

func (ss *SqlStore) SetContext(context context.Context) {
	ss.context = context
}

func (ss *SqlStore) Context() context.Context {
	return ss.context
}

func (ss *SqlStore) initConnection() {
	dataSource := *ss.settings.DataSource

	handle := setupConnection("master", dataSource, ss.settings)
	ss.masterX = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), time.Duration(*ss.settings.QueryTimeout)*time.Second, *ss.settings.Trace)

	if len(ss.settings.DataSourceReplicas) > 0 {
		ss.ReplicaXs = make([]*sqlxDBWrapper, len(ss.settings.DataSourceReplicas))
		for i, replica := range ss.settings.DataSourceReplicas {
			handle := setupConnection(fmt.Sprintf("replica-%v", i), replica, ss.settings)

			ss.ReplicaXs[i] = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), time.Duration(*ss.settings.QueryTimeout)*time.Second, *ss.settings.Trace)
		}
	}

	if len(ss.settings.DataSourceSearchReplicas) > 0 {
		ss.searchReplicaXs = make([]*sqlxDBWrapper, len(ss.settings.DataSourceSearchReplicas))
		for i, replica := range ss.settings.DataSourceSearchReplicas {
			handle := setupConnection(fmt.Sprintf("search-replica-%v", i), replica, ss.settings)

			ss.searchReplicaXs[i] = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), time.Duration(*ss.settings.QueryTimeout)*time.Second, *ss.settings.Trace)
		}
	}

	if len(ss.settings.ReplicaLagSettings) > 0 {
		ss.replicaLagHandles = make([]*dbsql.DB, len(ss.settings.ReplicaLagSettings))
		for i, src := range ss.settings.ReplicaLagSettings {
			if src.DataSource == nil {
				continue
			}
			ss.replicaLagHandles[i] = setupConnection(fmt.Sprintf(replicaLagPrefix+"-%d", i), *src.DataSource, ss.settings)
		}
	}
}

func (ss *SqlStore) DriverName() string {
	return *ss.settings.DriverName
}

func (ss *SqlStore) GetCurrentSchemaVersion() string {
	var schemaVersion string
	err := ss.GetMasterX().Get(&schemaVersion, "SELECT Value FROM Systems WHERE Name='Version'")
	if err != nil {
		slog.Error("failed to check current schema version", slog.Err(err))
	}

	return schemaVersion
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

	var version string

	if err := ss.GetReplicaX().Get(&version, sqlVersion); err != nil {
		return "", err
	}

	return version, nil
}

func (ss *SqlStore) GetMasterX() store_iface.SqlxExecutor {
	return ss.masterX
}

func (ss *SqlStore) SetMasterX(db *sql.DB) {
	ss.masterX = newSqlxDBWrapper(
		sqlx.NewDb(db, ss.DriverName()),
		time.Duration(*ss.settings.QueryTimeout)*time.Second,
		*ss.settings.Trace,
	)
}

// func (ss *SqlStore) GetInternalMasterDB() *sql.DB {
// 	return ss.GetMasterX().DB.DB
// }

func (ss *SqlStore) GetSearchReplicaX() store_iface.SqlxExecutor {
	if len(ss.settings.DataSourceSearchReplicas) == 0 {
		return ss.GetReplicaX()
	}

	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicaXs))
	return ss.searchReplicaXs[rrNum]
}

func (ss *SqlStore) GetReplicaX() store_iface.SqlxExecutor {
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster {
		return ss.GetMasterX()
	}
	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.ReplicaXs))
	return ss.ReplicaXs[rrNum]
}

func (ss *SqlStore) GetInternalReplicaDBs() []*sql.DB {
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster {
		return []*sql.DB{
			ss.masterX.DB.DB,
		}
	}

	dbs := make([]*sql.DB, len(ss.ReplicaXs))
	for i, rx := range ss.ReplicaXs {
		dbs[i] = rx.DB.DB
	}

	return dbs
}

// returns number of connections to master database
func (ss *SqlStore) TotalMasterDbConnections() int {
	return ss.masterX.Stats().OpenConnections
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
	for _, db := range ss.ReplicaXs {
		count = count + db.Stats().OpenConnections
	}

	return count
}

// counts all connections to replica source
func (ss *SqlStore) TotalSearchDbConnections() int {
	if len(ss.settings.DataSourceSearchReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.searchReplicaXs {
		count = count + db.Stats().OpenConnections
	}

	return count
}

func (ss *SqlStore) MarkSystemRanUnitTests() {
	props, err := ss.System().Get()
	if err != nil {
		return
	}

	unitTests := props[model.SystemRanUnitTests]
	if unitTests == "" {
		systemTests := &model.System{Name: model.SystemRanUnitTests, Value: "1"}
		ss.System().Save(systemTests)
	}
}

// checks if table does exist in database
func (ss *SqlStore) DoesTableExist(tableName string) bool {
	var count int64
	err := ss.GetMasterX().Get(&count, `SELECT count(relname) FROM pg_class WHERE relname=$1`, strings.ToLower(tableName))
	if err != nil {
		slog.Fatal("Failed to check if table exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
	var count int64
	err := ss.GetMasterX().Get(
		&count,
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

		slog.Fatal("Failed to check if column exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) DoesTriggerExist(triggerName string) bool {
	var count int64
	err := ss.GetMasterX().Get(&count, `SELECT COUNT(0) FROM pg_trigger WHERE tgname = $1`, triggerName)

	if err != nil {
		slog.Fatal("Failed to check if trigger exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) CreateColumnIfNotExists(tableName string, columnName string, mySqlColType string, postgresColType string, defaultValue string) bool {
	if ss.DoesColumnExist(tableName, columnName) {
		return false
	}

	_, err := ss.GetMasterX().ExecNoTimeout("ALTER TABLE " + tableName + " ADD " + columnName + " " + postgresColType + " DEFAULT '" + defaultValue + "'")
	if err != nil {
		slog.Fatal("Failed to create column", slog.Err(err))
	}

	return true
}

func (ss *SqlStore) RemoveTableIfExists(tableName string) bool {
	if !ss.DoesTableExist(tableName) {
		return false
	}

	_, err := ss.GetMasterX().ExecNoTimeout("DROP TABLE " + tableName)
	if err != nil {
		slog.Fatal("Failed to drop table", slog.Err(err))
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
func (ss *SqlStore) GetAllConns() []*sqlxDBWrapper {
	all := make([]*sqlxDBWrapper, len(ss.ReplicaXs)+1)
	copy(all, ss.ReplicaXs)
	all[len(ss.ReplicaXs)] = ss.masterX
	return all
}

// RecycleDBConnections closes active connections by setting the max conn lifetime
// to d, and then resets them back to their original duration.
func (ss *SqlStore) RecycleDBConnections(d time.Duration) {
	// Get old time.
	originalDuration := time.Duration(*ss.settings.ConnMaxLifetimeMilliseconds) * time.Millisecond
	// Set the max lifetimes for all connections.
	for _, conn := range ss.GetAllConns() {
		conn.SetConnMaxLifetime(d)
	}
	// Wait for that period with an additional 2 seconds of scheduling delay.
	time.Sleep(d + 2*time.Second)
	// Reset max lifetime back to original value.
	for _, conn := range ss.GetAllConns() {
		conn.SetConnMaxLifetime(originalDuration)
	}
}

// close all database connections
func (ss *SqlStore) Close() {
	ss.masterX.Close()
	for _, replica := range ss.ReplicaXs {
		replica.Close()
	}

	for _, replica := range ss.searchReplicaXs {
		replica.Close()
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
	ss.masterX.Exec(`DO
			$func$
			BEGIN
			   EXECUTE
			   (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
			    FROM   pg_class
			    WHERE  relkind = 'r'  -- only tables
			    AND    relnamespace = 'public'::regnamespace
			   );
			END
			$func$;`)
}

func (ss *SqlStore) GetQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (ss *SqlStore) CheckIntegrity() <-chan model.IntegrityCheckResult {
	results := make(chan model.IntegrityCheckResult)
	go CheckRelationalIntegrity(ss, results)
	return results
}

// func (ss *SqlStore) migrate(direction migrationDirection) error {
// 	assets := db.Assets()

// 	assetsList, err := assets.ReadDir(filepath.Join("migrations", ss.DriverName()))
// 	if err != nil {
// 		return err
// 	}

// 	driver, err := postgres.WithInstance(ss.masterX.DB.DB, &postgres.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	var assetNamesForDriver []string
// 	for _, entry := range assetsList {
// 		assetNamesForDriver = append(assetNamesForDriver, entry.Name())
// 	}

// 	source := bindata.Resource(assetNamesForDriver, func(name string) ([]byte, error) {
// 		return assets.ReadFile(filepath.Join("migrations", ss.DriverName(), name))
// 	})

// 	sourceDriver, err := bindata.WithInstance(source)
// 	if err != nil {
// 		return err
// 	}

// 	migrations, err := migrate.NewWithInstance(
// 		"go-bindata",
// 		sourceDriver,
// 		ss.DriverName(),
// 		driver,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer migrations.Close()

// 	switch direction {
// 	case migrationsDirectionUp:
// 		err = migrations.Up()
// 	case migrationsDirectionDown:
// 		err = migrations.Down()
// 	default:
// 		return fmt.Errorf("un supported migration direction %s", direction)
// 	}

// 	if err != nil && err != migrate.ErrNoChange && !errors.Is(err, os.ErrNotExist) {
// 		return err
// 	}

// 	return nil
// }

func (ss *SqlStore) migrate(direction migrationDirection) error {
	assets := db.Assets()

	assetsList, err := assets.ReadDir(filepath.Join("migrations", ss.DriverName()))
	if err != nil {
		return err
	}

	assetNamesForDriver := make([]string, len(assetsList))
	for i, entry := range assetsList {
		assetNamesForDriver[i] = entry.Name()
	}

	src, err := mbindata.WithInstance(&mbindata.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			return assets.ReadFile(filepath.Join("migrations", ss.DriverName(), name))
		},
	})
	if err != nil {
		return err
	}

	driver, err := ps.WithInstance(ss.masterX.DB.DB, &ps.Config{
		Config: drivers.Config{
			StatementTimeoutInSecs: *ss.settings.MigrationsStatementTimeoutSeconds,
		},
	})
	if err != nil {
		return err
	}

	opts := []morph.EngineOption{
		morph.WithLogger(log.New(&morphWriter{}, "", log.Lshortfile)),
		morph.WithLock("sn-lock-key"),
	}
	engine, err := morph.New(context.Background(), driver, src, opts...)
	if err != nil {
		return err
	}
	defer engine.Close()

	switch direction {
	case migrationsDirectionDown:
		_, err = engine.ApplyDown(-1)
		return err
	default:
		return engine.ApplyAll()
	}
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

// ensureMinimumDBVersion gets the DB version and ensures it is
// above the required minimum version requirements.
func (ss *SqlStore) ensureMinimumDBVersion(ver string) (bool, error) {
	intVer, err2 := strconv.Atoi(ver)
	if err2 != nil {
		return false, fmt.Errorf("cannot parse DB version: %v", err2)
	}
	if intVer < MinimumRequiredPostgresVersion {
		return false, fmt.Errorf("minimum Postgres version requirements not met. Found: %s, Wanted: %s", versionString(intVer, *ss.settings.DriverName), versionString(MinimumRequiredPostgresVersion, *ss.settings.DriverName))
	}
	return true, nil
}

// versionString converts an integer representation of a DB version
// to a pretty-printed string.
// Postgres doesn't follow three-part version numbers from 10.0 onwards:
// https://www.postgresql.org/docs/13/libpq-status.html#LIBPQ-PQSERVERVERSION.
// For MySQL, we consider a major*1000 + minor*100 + patch format.
func versionString(v int, driver string) string {
	minor := v % 10000
	major := v / 10000
	return strconv.Itoa(major) + "." + strconv.Itoa(minor)
}

func (ss *SqlStore) GetDBSchemaVersion() (int, error) {
	var version int
	if err := ss.GetMasterX().Get(&version, "SELECT Version FROM db_migrations ORDER BY Version DESC LIMIT 1"); err != nil {
		return 0, errors.Wrap(err, "unable to select from db_migrations")
	}
	return version, nil
}

func (ss *SqlStore) GetAppliedMigrations() ([]model.AppliedMigration, error) {
	migrations := []model.AppliedMigration{}
	if err := ss.GetMasterX().Select(&migrations, "SELECT Version, Name FROM db_migrations ORDER BY Version DESC"); err != nil {
		return nil, errors.Wrap(err, "unable to select from db_migrations")
	}

	return migrations, nil
}

// finalizeTransaction ensures a transaction is closed after use, rolling back if not already committed.
func (s *SqlStore) FinalizeTransaction(transaction driver.Tx) {
	if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
		slog.Error("Failed to rollback transaction", slog.Err(err))
	}
}
