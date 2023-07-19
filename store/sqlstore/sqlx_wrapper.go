package sqlstore

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/sitename/sitename/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var _ *gorm.DB = (*gormDBWrapper)(nil)

func newGormDBWrapper(handle *sql.DB, sqlSettings *model.SqlSettings) (*gorm.DB, error) {
	var gormLog logger.Interface
	if *sqlSettings.Trace {
		gormLog = logger.New(
			log.New(os.Stderr, "\r", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
			},
		)
	}

	return gorm.Open(postgres.New(postgres.Config{Conn: handle}), &gorm.Config{
		Logger: gormLog,
	})
}

// type gormDBWrapper struct {
// 	DB           *gorm.DB
// 	queryTimeout time.Duration
// }

// func (g *gormDBWrapper) Select(dest any, query string, args ...any) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
// 	defer cancel()

// 	return g.DB.WithContext(ctx).Raw(query, args...).Find(dest).Error
// }

// func (g *gormDBWrapper) First(dest any, conds ...any) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
// 	defer cancel()

// 	return g.DB.WithContext(ctx).First(dest, conds...).Error
// }

// func (g *gormDBWrapper) Stats() sql.DBStats {
// 	sqldb, err := g.DB.DB()
// 	if err != nil {
// 		return sql.DBStats{}
// 	}
// 	return sqldb.Stats()
// }

// func (g *gormDBWrapper) Conn(ctx context.Context) (*sql.Conn, error) {
// 	sqldb, err := g.DB.DB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sqldb.Conn(ctx)
// }

// func (g *gormDBWrapper) Commit() error {
// 	if g.inTransaction() {
// 		return g.DB.Commit().Error
// 	}
// 	return errors.New("not in a transaction")
// }

// func (g *gormDBWrapper) Rollback() error {
// 	if g.inTransaction() {
// 		return g.DB.Rollback().Error
// 	}
// 	return errors.New("not in a transaction")
// }

// func (g *gormDBWrapper) inTransaction() bool {
// 	switch g.DB.Statement.ConnPool.(type) {
// 	case gorm.TxBeginner, gorm.ConnPoolBeginner:
// 		return false
// 	default:
// 		return true
// 	}
// }

// func (g *gormDBWrapper) Begin() (*gorm.DB, error) {
// 	if g.inTransaction() {
// 		return nil, errors.New("already in a transaction")
// 	}

// 	tx := g.DB.Begin()
// 	err := tx.Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &gormDBWrapper{
// 		DB:           tx,
// 		queryTimeout: g.queryTimeout,
// 	}, nil
// }

// func (w *gormDBWrapper) NamedExec(query string, arg any) (sql.Result, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	res := w.DB.WithContext(ctx).Raw(query, arg)
// 	if res.Error != nil {
// 		return nil, res.Error
// 	}

// 	return &sqlResult{res.RowsAffected}, nil
// }

// type sqlResult struct {
// 	rowsAffected int64
// }

// func (s *sqlResult) RowsAffected() (int64, error) { return s.rowsAffected, nil }
// func (s *sqlResult) LastInsertId() (int64, error) { return 0, nil }

// func (w *gormDBWrapper) Table(name string) *gorm.DB {
// 	w.DB = w.DB.Table(name)
// 	return w
// }

// func (w *gormDBWrapper) Model(arg any) *gorm.DB {
// 	w.DB = w.DB.Model(arg)
// 	return w
// }

// func (w *gormDBWrapper) Association(name string) *gorm.Association {
// 	return w.DB.Association(name)
// }

// func (w *gormDBWrapper) Upsert(arg any) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	return w.DB.WithContext(ctx).Save(arg).Error
// }

// var r = strings.NewReplacer()

// func (g *gormDBWrapper) Exec(query string, args ...any) (sql.Result, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), g.queryTimeout)
// 	defer cancel()

// 	res := g.DB.WithContext(ctx).Raw(query, args...)
// 	if res.Error != nil {
// 		return nil, res.Error
// 	}

// 	return &sqlResult{res.RowsAffected}, nil
// }

// func (g *gormDBWrapper) ExecNoTimeout(query string, args ...any) (sql.Result, error) {
// 	sqldb, err := g.DB.DB()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sqldb.ExecContext(context.Background(), query, args...)
// }

// func (w *gormDBWrapper) NamedQuery(query string, arg any) (store_iface.RowsScanner, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	return w.DB.WithContext(ctx).Raw(query, arg).Rows()
// }

// func (w *gormDBWrapper) QueryRow(query string, args ...any) store_iface.Scanner {
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	return w.DB.WithContext(ctx).Raw(query, args...).Row()
// }

// func (w *gormDBWrapper) Query(query string, args ...any) (store_iface.RowsScanner, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), w.queryTimeout)
// 	defer cancel()

// 	return w.DB.WithContext(ctx).Raw(query, args...).Rows()
// }

// func (g *gormDBWrapper) BeginXWithIsolation(opts *sql.TxOptions) (*gorm.DB, error) {
// 	if g.inTransaction() {
// 		return nil, errors.New("already in a transaction")
// 	}

// 	return &gormDBWrapper{
// 		DB:           g.DB.Begin(opts),
// 		queryTimeout: g.queryTimeout,
// 	}, nil
// }

// func removeSpace(r rune) rune {
// 	// Strip everything except ' '
// 	// This also strips out more than one space,
// 	// but we ignore it for now until someone complains.
// 	if unicode.IsSpace(r) && r != ' ' {
// 		return -1
// 	}
// 	return r
// }

// func printArgs(query string, dur time.Duration, args ...any) {
// 	query = strings.Map(removeSpace, query)
// 	fields := make([]slog.Field, 0, len(args)+1)
// 	fields = append(fields, slog.Duration("duration", dur))
// 	for i, arg := range args {
// 		fields = append(fields, slog.Any("arg"+strconv.Itoa(i), arg))
// 	}
// 	slog.Debug(query, fields...)
// }
