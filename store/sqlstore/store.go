package sqlstore

import (
	"context"
	"database/sql"
	dbsql "database/sql"
	"fmt"
	"io/fs"
	"log"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	ps "github.com/mattermost/morph/drivers/postgres"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattermost/morph"
	mbindata "github.com/mattermost/morph/sources/embedded"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/db"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
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

var _ store.Store = (*SqlStore)(nil)

type SqlStore struct {
	// rrCounter and srCounter should be kept first.
	rrCounter int64
	srCounter int64

	master         *sqlDBWrapper
	Replicas       []*sqlDBWrapper
	searchReplicas []*sqlDBWrapper

	replicaLagHandles []*dbsql.DB
	stores            *SqlStoreStores
	settings          *model_helper.SqlSettings
	lockedToMaster    bool
	context           context.Context
	metrics           einterfaces.MetricsInterface
}

// New initializes connections to postgresql database
// also migrates all the sql schema using gorp
func New(settings model_helper.SqlSettings, metrics einterfaces.MetricsInterface) *SqlStore {
	store := &SqlStore{
		rrCounter: 0,
		srCounter: 0,
		settings:  &settings,
		metrics:   metrics,
	}

	err := store.initConnection()
	if err != nil {
		slog.Fatal("Error setting up connections", slog.Err(err))
	}
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
	err = store.migrate(migrationsDirectionUp, false)
	if err != nil {
		slog.Fatal("Failed to apply database migrations.", slog.Err(err))
	}

	store.setupStores()
	return store
}

// setupConnection opens connection to database, check if it works by ping
func setupConnection(connType, dataSource string, settings *model_helper.SqlSettings) (*dbsql.DB, error) {
	db, err := dbsql.Open(*settings.DriverName, dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open SQL connection")
	}

	for i := 0; i < DBPingAttempts; i++ {
		slog.Info("Pinging SQL", slog.String("database", connType))
		ctx, cancel := context.WithTimeout(context.Background(), DBPingTimeoutSecs*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			if i == DBPingAttempts-1 {
				return nil, err
			}
			slog.Error("Failed to ping DB", slog.Err(err), slog.Int("retrying in seconds", DBPingTimeoutSecs))
			time.Sleep(DBPingTimeoutSecs * time.Second)
			continue
		}
		break
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

	return db, nil
}

func (ss *SqlStore) SetContext(ctx context.Context) {
	ss.context = ctx
}

func (ss *SqlStore) Context() context.Context {
	return ss.context
}

func (ss *SqlStore) initConnection() (err error) {
	handle, err := setupConnection("master", *ss.settings.DataSource, ss.settings)
	if err != nil {
		return err
	}
	ss.master = newSqlDbWrapper(handle, *ss.settings)

	if len(ss.settings.DataSourceReplicas) > 0 {
		ss.Replicas = make([]*sqlDBWrapper, len(ss.settings.DataSourceReplicas))

		for i, replica := range ss.settings.DataSourceReplicas {
			replicaName := fmt.Sprintf("replica-%v", i)
			handle, err := setupConnection(replicaName, replica, ss.settings)
			if err != nil {
				return errors.Wrapf(err, "failed to setup replica connection: %s", replicaName)
			}
			ss.Replicas[i] = newSqlDbWrapper(handle, *ss.settings)
		}
	}

	if len(ss.settings.DataSourceSearchReplicas) > 0 {
		ss.searchReplicas = make([]*sqlDBWrapper, len(ss.settings.DataSourceSearchReplicas))

		for i, replica := range ss.settings.DataSourceSearchReplicas {
			replicaName := fmt.Sprintf("search-replica-%v", i)
			handle, err := setupConnection(replicaName, replica, ss.settings)
			if err != nil {
				return errors.Wrapf(err, "failed to setup search replica connection: %s", replicaName)
			}
			ss.searchReplicas[i] = newSqlDbWrapper(handle, *ss.settings)
		}
	}

	if len(ss.settings.ReplicaLagSettings) > 0 {
		ss.replicaLagHandles = make([]*dbsql.DB, len(ss.settings.ReplicaLagSettings))

		for i, src := range ss.settings.ReplicaLagSettings {
			if src.DataSource == nil {
				continue
			}
			ss.replicaLagHandles[i], err = setupConnection(fmt.Sprintf(replicaLagPrefix+"-%d", i), *src.DataSource, ss.settings)
			if err != nil {
				slog.Warn("Failed to setup replica lag handle. Skipping..", slog.String("db", fmt.Sprintf(replicaLagPrefix+"-%d", i)), slog.Err(err))
				continue
			}
		}
	}

	return nil
}

func (ss *SqlStore) DriverName() string {
	return *ss.settings.DriverName
}

func (ss *SqlStore) GetCurrentSchemaVersion() string {
	var schemaVersion string
	err := ss.GetMaster().QueryRow("SELECT Value FROM Systems WHERE Name='Version'").Scan(&schemaVersion)
	if err != nil {
		slog.Error("failed to check current schema version", slog.Err(err))
	}

	return schemaVersion
}

// GetDbVersion returns the version of the database being used.
// If numerical is set to true, it attempts to return a numerical version string
// that can be parsed by callers.
func (ss *SqlStore) GetDbVersion(numerical bool) (string, error) {
	var sqlVersionQuery string
	if numerical {
		sqlVersionQuery = `SHOW server_version_num`
	} else {
		sqlVersionQuery = `SHOW server_version`
	}

	var version string

	if err := ss.GetReplica().QueryRow(sqlVersionQuery).Scan(&version); err != nil {
		return "", err
	}

	return version, nil
}

func (ss *SqlStore) GetMaster() store.ContextRunner {
	return ss.master
}

func (ss *SqlStore) GetReplica() boil.ContextExecutor {
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster {
		return ss.GetMaster()
	}
	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.Replicas))
	return ss.Replicas[rrNum]
}

// func (ss *SqlStore) GetSearchReplicaX() *gorm.DB {
// 	if len(ss.settings.DataSourceSearchReplicas) == 0 {
// 		return ss.GetReplica()
// 	}

// 	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicaXs))
// 	return ss.searchReplicaXs[rrNum]
// }

// returns number of connections to master database
func (ss *SqlStore) TotalMasterDbConnections() int {
	return ss.master.sqlDBInterface.(*sql.DB).Stats().OpenConnections
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
		count += db.sqlDBInterface.(*sql.DB).Stats().OpenConnections
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
		count += db.sqlDBInterface.(*sql.DB).Stats().OpenConnections
	}

	return count
}

func (ss *SqlStore) MarkSystemRanUnitTests() {
	props, err := ss.System().Get()
	if err != nil {
		return
	}

	unitTests := props[model_helper.SystemRanUnitTests]
	if unitTests == "" {
		systemTests := model.System{Name: model_helper.SystemRanUnitTests, Value: "1"}
		ss.System().Save(systemTests)
	}
}

// checks if table does exist in database
func (ss *SqlStore) DoesTableExist(tableName string) bool {
	var count int64
	err := ss.GetMaster().QueryRow(`SELECT COUNT(relname) FROM pg_class WHERE relname=$1`, strings.ToLower(tableName)).Scan(&count)
	if err != nil {
		slog.Fatal("Failed to check if table exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
	var count int64
	err := ss.GetMaster().QueryRow(
		`SELECT COUNT(0)
			FROM   pg_attribute
			WHERE  attrelid = $1::regclass
			AND    attname = $2
			AND    NOT attisdropped`,
		strings.ToLower(tableName),
		strings.ToLower(columnName),
	).Scan(&count)

	if err != nil {
		if err.Error() == "pq: relation \""+strings.ToLower(tableName)+"\" does not exist" {
			return false
		}

		slog.Fatal("Failed to check if column exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) IsUniqueConstraintError(err error, indexNames []string) bool {
	for _, contain := range indexNames {
		if strings.Contains(err.Error(), contain) {
			return true
		}
	}

	return false
}

// Get all databases connections
func (ss *SqlStore) GetAllConns() []*sqlDBWrapper {
	all := make([]*sqlDBWrapper, len(ss.Replicas)+1)
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
		if db, ok := conn.sqlDBInterface.(*sql.DB); ok && db != nil {
			db.SetConnMaxLifetime(originalDuration)
		}
	}
	// Wait for that period with an additional 2 seconds of scheduling delay.
	time.Sleep(d + 2*time.Second)
	// Reset max lifetime back to original value.
	for _, conn := range ss.GetAllConns() {
		if db, ok := conn.sqlDBInterface.(*sql.DB); ok && db != nil {
			db.SetConnMaxLifetime(originalDuration)
		}
	}
}

// close all database connections
func (ss *SqlStore) Close() {
	connections := append(ss.Replicas, ss.searchReplicas...)
	connections = append(connections, ss.master)

	for _, conn := range connections {
		if db, ok := conn.sqlDBInterface.(*sql.DB); ok && db != nil {
			db.Close()
		}
	}
}

// constraint db queries to call only master replica
func (ss *SqlStore) LockToMaster() {
	ss.lockedToMaster = true
}

// let db queries free to call any db Replicas
func (ss *SqlStore) UnlockFromMaster() {
	ss.lockedToMaster = false
}

func (ss *SqlStore) DropAllTables() {
	ss.master.Exec(`DO
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

func (ss *SqlStore) GetQueryBuilder(placeholderFormats ...sq.PlaceholderFormat) sq.StatementBuilderType {
	res := sq.StatementBuilder
	if len(placeholderFormats) == 0 {
		return res.PlaceholderFormat(sq.Dollar)
	}
	return res.PlaceholderFormat(placeholderFormats[0])
}

func (ss *SqlStore) CheckIntegrity() <-chan model_helper.IntegrityCheckResult {
	results := make(chan model_helper.IntegrityCheckResult)
	go CheckRelationalIntegrity(ss, results)
	return results
}

func (ss *SqlStore) migrate(direction migrationDirection, drRun bool) error {
	engine, err := ss.initMorph(drRun)
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

func (ss *SqlStore) initMorph(dryRun bool) (*morph.Morph, error) {
	assets := db.Assets()

	assetsList, err := assets.ReadDir(path.Join("migrations", ss.DriverName()))
	if err != nil {
		return nil, err
	}

	assetNamesForDriver := lo.Map(assetsList, func(item fs.DirEntry, _ int) string { return item.Name() })
	src, err := mbindata.WithInstance(&mbindata.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			return assets.ReadFile(path.Join("migrations", ss.DriverName(), name))
		},
	})
	if err != nil {
		return nil, err
	}

	driver, err := ps.WithInstance(ss.GetMaster().(*sqlDBWrapper).sqlDBInterface.(*sql.DB))
	if err != nil {
		return nil, err
	}

	opts := []morph.EngineOption{
		morph.WithLogger(log.New(&morphWriter{}, "", log.Lshortfile)),
		morph.WithLock("sn-lock-key"),
		morph.SetStatementTimeoutInSeconds(*ss.settings.MigrationsStatementTimeoutSeconds),
		morph.SetDryRun(dryRun),
	}

	engine, err := morph.New(context.Background(), driver, src, opts...)
	if err != nil {
		return nil, err
	}

	return engine, nil
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
	if err := ss.GetMaster().QueryRow("SELECT Version FROM db_migrations ORDER BY Version DESC LIMIT 1").Scan(&version); err != nil {
		return 0, errors.Wrap(err, "unable to select from db_migrations")
	}
	return version, nil
}

func (ss *SqlStore) GetAppliedMigrations() ([]model_helper.AppliedMigration, error) {
	migrations := []model_helper.AppliedMigration{}
	err := queries.Raw("SELECT * FROM db_migrations ORDER BY Version DESC").Bind(context.Background(), ss.GetMaster(), &migrations)
	if err != nil {
		return nil, errors.Wrap(err, "unable to select from db_migrations")

	}

	return migrations, nil
}

func (s *SqlStore) FinalizeTransaction(tx store.ContextRunner) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, dbsql.ErrTxDone) {
		slog.Error("failed to rollback a transaction", slog.Err(err))
	}
}
