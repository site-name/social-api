package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	// Needed for Postgresql driver
	_ "github.com/lib/pq"
	"github.com/sitename/sitename/modules/setting"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"
)

// Engine represents a xorm engine or session.
type Engine interface {
	Table(tableNameOrBean interface{}) *xorm.Session
	Count(...interface{}) (int64, error)
	Decr(column string, arg ...interface{}) *xorm.Session
	Delete(interface{}) (int64, error)
	Exec(...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	ID(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Incr(column string, arg ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session
	SQL(interface{}, ...interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
	Asc(colNames ...string) *xorm.Session
	Desc(colNames ...string) *xorm.Session
	Limit(limit int, start ...int) *xorm.Session
	SumInt(bean interface{}, columnName string) (res int64, err error)
}

const (
	// When queries are broken down in parts because of the number
	// of parameters, attempt to break by this amount
	maxQueryParameters = 300
)

var (
	x      *xorm.Engine
	tables []interface{}

	// HasEngine specifies if we have a xorm.Engine
	HasEngine bool
)

func init() {
	tables = append(tables,
		new(User),
	)

	gonicNames := []string{"SSL", "UID"}
	for _, name := range gonicNames {
		names.LintGonicMapper[name] = true
	}
}

func getEngine() (*xorm.Engine, error) {
	// Get db connection URL
	connStr, err := setting.DBConnStr()
	if err != nil {
		return nil, err
	}

	var engine *xorm.Engine

	if len(setting.Database.Schema) > 0 {
		// OK whilst we sort out our schema issues - create a schema aware postgres
		registerPostgresSchemaDriver()
		engine, err = xorm.NewEngine("postgresschema", connStr)
	}
	if err != nil {
		return nil, err
	}
	engine.SetSchema(setting.Database.Schema)
	return engine, nil
}

// NewTestEngine sets a new test xorm.Engine
func NewTestEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Connect to database: %v", err)
	}

	x.SetMapper(names.GonicMapper{})
	x.SetLogger(NewXORMLogger(!setting.IsProd()))
	x.ShowSQL(!setting.IsProd())
	return x.StoreEngine("InnoDB").Sync2(tables...)
}

// SetEngine sets the xorm engine
func SetEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}

	x.SetMapper(names.GonicMapper{})
	// WARNING: for serv command, MUST remove the output to os.stdout,
	// so use log file to instead print to stdout.
	x.SetLogger(NewXORMLogger(setting.Database.LogSQL))
	x.ShowSQL(setting.Database.LogSQL)
	x.SetMaxIdleConns(setting.Database.MaxIdleConns)
	x.SetMaxOpenConns(setting.Database.MaxOpenConns)
	x.SetConnMaxLifetime(setting.Database.ConnMaxLifetime)
	return nil
}

// NewEngine initializes a new xorm.Engine
// This function must never call .Sync2() if the provided migration function fails.
// When called from the "doctor" command, the migration function is a version check
// that prevents the doctor from fixing anything in the database if the migration level
// is different from the expected value.
func NewEngine(ctx context.Context, migrateFunc func(*xorm.Engine) error) (err error) {
	if err = SetEngine(); err != nil {
		return err
	}

	x.SetDefaultContext(ctx)

	if err = x.Ping(); err != nil {
		return err
	}

	if err = migrateFunc(x); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v", err)
	}

	return nil
}

// NamesToBean return a list of beans or an error
func NamesToBean(names ...string) ([]interface{}, error) {
	beans := []interface{}{}
	if len(names) == 0 {
		beans = append(beans, tables...)
		return beans, nil
	}
	// Need to map provided names to beans...
	beanMap := make(map[string]interface{})
	for _, bean := range tables {
		beanMap[strings.ToLower(reflect.Indirect(reflect.ValueOf(bean)).Type().Name())] = bean
		beanMap[strings.ToLower(x.TableName(bean))] = bean
		beanMap[strings.ToLower(x.TableName(bean, true))] = bean
	}

	gotBean := make(map[interface{}]bool)
	for _, name := range names {
		bean, ok := beanMap[strings.ToLower(strings.TrimSpace(name))]
		if !ok {
			return nil, fmt.Errorf("No table found that matches: %s", name)
		}
		if !gotBean[bean] {
			beans = append(beans, bean)
			gotBean[bean] = true
		}
	}
	return beans, nil
}

// Statistic contains the database statistics
type Statistic struct {
	Counter struct {
		User, Org, PublicKey,
		Repo, Watch, Star, Action, Access,
		Issue, Comment, Oauth, Follow,
		Mirror, Release, LoginSource, Webhook,
		Milestone, Label, HookTask,
		Team, UpdateTask, Attachment int64
	}
}

// GetStatistic returns the database statistics
func GetStatistic() (stats Statistic) {
	stats.Counter.User = CountUsers()
	stats.Counter.Org = CountOrganizations()
	stats.Counter.Oauth = 0

	return
}

// Ping tests if database is alive
func Ping() error {
	if x != nil {
		return x.Ping()
	}

	return errors.New("database not configured")
}

// DumpDatabase dumps all data from database according the special database SQL syntax to file system.
func DumpDatabase(filepath string, dbType string) error {
	var tbs []*schemas.Table
	for _, t := range tables {
		t, err := x.TableInfo(t)
		if err != nil {
			return err
		}
		tbs = append(tbs, t)
	}

	type Version struct {
		ID      int64 `xorm:"pk autoincr"`
		Version int64
	}
	t, err := x.TableInfo(Version{})
	if err != nil {
		return err
	}
	tbs = append(tbs, t)

	if len(dbType) > 0 {
		return x.DumpTablesToFile(tbs, filepath, schemas.DBType(dbType))
	}
	return x.DumpTablesToFile(tbs, filepath)
}

// MaxBatchInsertSize returns the table's max batch insert size
func MaxBatchInsertSize(bean interface{}) int {
	t, err := x.TableInfo(bean)
	if err != nil {
		return 50
	}
	return 999 / len(t.ColumnsSeq())
}

// Count returns records number according struct's fields as database query conditions
func Count(bean interface{}) (int64, error) {
	return x.Count(bean)
}

// IsTableNotEmpty returns true if table has at least one record
func IsTableNotEmpty(tableName string) (bool, error) {
	return x.Table(tableName).Exist()
}

// DeleteAllRecords will delete all the records of this table
func DeleteAllRecords(tableName string) error {
	_, err := x.Exec(fmt.Printf("DELETE FROM %s", tableName))
	return err
}

// GetMaxID will return max id of the table
func GetMaxID(beanOrTableName interface{}) (maxID int64, err error) {
	_, err = x.Select("MAX(id)").Table(beanOrTableName).Get(&maxID)
	return
}

// FindByMaxID filled results as the condition from database
func FindByMaxID(maxID int64, limit int, results interface{}) error {
	return x.Where("id <= ?", maxID).OrderBy("id DESC").Limit(limit).Find(results)
}
