package sqlstore

import (
	"context"
	"database/sql"
	dbsql "database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	jackcpgconn "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"gorm.io/gorm"
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

type SqlStore struct {
	// rrCounter and srCounter should be kept first.
	rrCounter int64
	srCounter int64

	master         *gorm.DB
	replicas       []*gorm.DB
	searchReplicas []*gorm.DB

	// masterX         *sqlxDBWrapper
	// replicaXs       []*sqlxDBWrapper
	// searchReplicaXs []*sqlxDBWrapper

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
	err = store.migrate(migrationsDirectionUp)
	if err != nil {
		slog.Fatal("Failed to apply database migrations.", slog.Err(err))
	}

	// set up tables before performing migrations:
	store.setupTables()
	return store
}

// setupConnection opens connection to database, check if it works by ping
func setupConnection(connType string, dataSource string, settings *model.SqlSettings) (*dbsql.DB, error) {
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

func (ss *SqlStore) SetContext(context context.Context) {
	ss.context = context
}

func (ss *SqlStore) Context() context.Context {
	return ss.context
}

func (ss *SqlStore) initConnection() error {
	dataSource := *ss.settings.DataSource

	handle, err := setupConnection("master", dataSource, ss.settings)
	if err != nil {
		return err
	}
	// ss.masterX = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), ss.settings)
	ss.master, err = newGormDBWrapper(handle, ss.settings)
	if err != nil {
		return err
	}

	if len(ss.settings.DataSourceReplicas) > 0 {
		// ss.replicaXs = make([]*sqlxDBWrapper, len(ss.settings.DataSourceReplicas))
		ss.replicas = make([]*gorm.DB, len(ss.settings.DataSourceReplicas))

		for i, replica := range ss.settings.DataSourceReplicas {
			replicaName := fmt.Sprintf("replica-%v", i)
			handle, err := setupConnection(replicaName, replica, ss.settings)
			if err != nil {
				return errors.Wrapf(err, "failed to setup replica connection: %s", replicaName)
			}
			// ss.replicaXs[i] = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), ss.settings)
			ss.replicas[i], err = newGormDBWrapper(handle, ss.settings)
			if err != nil {
				return err
			}
		}
	}

	if len(ss.settings.DataSourceSearchReplicas) > 0 {
		// ss.searchReplicaXs = make([]*sqlxDBWrapper, len(ss.settings.DataSourceSearchReplicas))
		ss.searchReplicas = make([]*gorm.DB, len(ss.settings.DataSourceSearchReplicas))

		for i, replica := range ss.settings.DataSourceSearchReplicas {
			replicaName := fmt.Sprintf("search-replica-%v", i)
			handle, err := setupConnection(replicaName, replica, ss.settings)
			if err != nil {
				return errors.Wrapf(err, "failed to setup search replica connection: %s", replicaName)
			}
			// ss.searchReplicaXs[i] = newSqlxDBWrapper(sqlx.NewDb(handle, ss.DriverName()), ss.settings)
			ss.searchReplicas[i], err = newGormDBWrapper(handle, ss.settings)
			if err != nil {
				return err
			}
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
	err := ss.GetMaster().Raw("SELECT Value FROM Systems WHERE Name='Version'").Scan(&schemaVersion).Error
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

	if err := ss.GetReplica().Raw(sqlVersion).Scan(&sqlVersion).Error; err != nil {
		return "", err
	}

	return version, nil
}

func (ss *SqlStore) GetMaster(noTimeout ...bool) *gorm.DB {
	if len(noTimeout) > 0 && noTimeout[0] {
		return ss.master
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ss.settings.QueryTimeout)*time.Second)
	defer cancel()
	return ss.master.WithContext(ctx)
}

func (ss *SqlStore) GetReplica(noTimeout ...bool) *gorm.DB {
	if len(ss.settings.DataSourceReplicas) == 0 || ss.lockedToMaster {
		return ss.GetMaster(noTimeout...)
	}
	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.replicas))
	db := ss.replicas[rrNum]

	if len(noTimeout) > 0 && noTimeout[0] {
		return db
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*ss.settings.QueryTimeout)*time.Second)
	defer cancel()
	return db.WithContext(ctx)
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
	db, err := ss.master.DB()
	if err != nil {
		slog.Error("failed to retrieve underlying *sql.DB instance", slog.Err(err))
		return 0
	}
	return db.Stats().OpenConnections
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
	for _, gormDB := range ss.replicas {
		db, err := gormDB.DB()
		if err != nil {
			slog.Error("failed to retrieve underlying replica *sql.DB instance", slog.Err(err))
			continue
		}
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
	for _, gormDB := range ss.searchReplicas {
		db, err := gormDB.DB()
		if err != nil {
			slog.Error("failed to retrieve underlying search replica *sql.DB instance", slog.Err(err))
			continue
		}
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
	err := ss.GetMaster().Raw(`SELECT COUNT(relname) FROM pg_class WHERE relname=$1`, strings.ToLower(tableName)).Scan(&count).Error
	if err != nil {
		slog.Fatal("Failed to check if table exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) DoesColumnExist(tableName string, columnName string) bool {
	var count int64
	err := ss.GetMaster().Raw(
		`SELECT COUNT(0)
			FROM   pg_attribute
			WHERE  attrelid = $1::regclass
			AND    attname = $2
			AND    NOT attisdropped`,
		strings.ToLower(tableName),
		strings.ToLower(columnName),
	).Scan(&count).Error

	if err != nil {
		if err.Error() == "pq: relation \""+strings.ToLower(tableName)+"\" does not exist" {
			return false
		}

		slog.Fatal("Failed to check if column exists", slog.Err(err))
	}

	return count > 0
}

func (ss *SqlStore) IsUniqueConstraintError(err error, indexNames []string) bool {
	unique := false

	switch errT := err.(type) {
	case *pq.Error:
		unique = errT.Code == "23505"
	case *pgconn.PgError:
		unique = errT.Code == "23505"
	case *jackcpgconn.PgError:
		unique = errT.Code == "23505"
	}

	strErr := err.Error()
	return unique && lo.SomeBy(indexNames, func(index string) bool { return strings.Contains(strErr, index) })
}

// Get all databases connections
func (ss *SqlStore) GetAllConns() []*gorm.DB {
	all := make([]*gorm.DB, len(ss.replicas)+1)
	copy(all, ss.replicas)
	all[len(ss.replicas)] = ss.master
	return all
}

// RecycleDBConnections closes active connections by setting the max conn lifetime
// to d, and then resets them back to their original duration.
func (ss *SqlStore) RecycleDBConnections(d time.Duration) {
	// Get old time.
	originalDuration := time.Duration(*ss.settings.ConnMaxLifetimeMilliseconds) * time.Millisecond
	// Set the max lifetimes for all connections.
	for _, conn := range ss.GetAllConns() {
		db, err := conn.DB()
		if err == nil && db != nil {
			db.SetConnMaxLifetime(d)
		}
	}
	// Wait for that period with an additional 2 seconds of scheduling delay.
	time.Sleep(d + 2*time.Second)
	// Reset max lifetime back to original value.
	for _, conn := range ss.GetAllConns() {
		db, err := conn.DB()
		if err == nil && db != nil {
			db.SetConnMaxLifetime(originalDuration)
		}
	}
}

// close all database connections
func (ss *SqlStore) Close() {
	connections := append(ss.replicas, ss.searchReplicas...)
	connections = append(connections, ss.master)

	for _, conn := range connections {
		db, err := conn.DB()
		if err == nil && db != nil {
			db.Close()
		}
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

func (ss *SqlStore) GetQueryBuilder(placeholderFormats ...squirrel.PlaceholderFormat) sq.StatementBuilderType {
	res := sq.StatementBuilder
	if len(placeholderFormats) == 0 {
		return res.PlaceholderFormat(squirrel.Dollar)
	}
	return res.PlaceholderFormat(placeholderFormats[0])
}

func (ss *SqlStore) CheckIntegrity() <-chan model.IntegrityCheckResult {
	results := make(chan model.IntegrityCheckResult)
	go CheckRelationalIntegrity(ss, results)
	return results
}

type m2mRelation struct {
	model     any
	field     string
	joinTable any
}

func (ss *SqlStore) migrate(direction migrationDirection) error {
	// account
	for _, model := range []any{&model.User{}, &model.Address{}, &model.Status{}, &model.UserAccessToken{}, &model.CustomerEvent{}, &model.CustomerNote{}, &model.AppToken{}} {
		if err := ss.master.AutoMigrate(model); err != nil {
			return err
		}
	}

	// attribute
	for _, m2mRel := range []m2mRelation{
		{&model.AttributeVariant{}, "AssignedVariants", &model.AssignedVariantAttribute{}},
		{&model.AssignedVariantAttribute{}, "Values", &model.AssignedVariantAttributeValue{}},
		{&model.AttributePage{}, "AssignedPages", &model.AssignedPageAttribute{}},
		{&model.AssignedPageAttribute{}, "Values", &model.AssignedPageAttributeValue{}},
		{&model.AttributeProduct{}, "AssignedProducts", &model.AssignedProductAttribute{}},
		{&model.AssignedProductAttribute{}, "Values", &model.AssignedProductAttributeValue{}},
	} {
		if err := ss.master.SetupJoinTable(m2mRel.model, m2mRel.field, m2mRel.joinTable); err != nil {
			return err
		}
	}

	return nil
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
	if err := ss.GetMaster().Raw("SELECT Version FROM db_migrations ORDER BY Version DESC LIMIT 1").Scan(&version).Error; err != nil {
		return 0, errors.Wrap(err, "unable to select from db_migrations")
	}
	return version, nil
}

func (ss *SqlStore) GetAppliedMigrations() ([]model.AppliedMigration, error) {
	migrations := []model.AppliedMigration{}
	if err := ss.GetMaster().Table("db_migrations").Order("Version DESC").Find(&migrations).Error; err != nil {
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
