// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// FileInfo is an object representing the database table.
type FileInfo struct {
	ID              string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatorID       string                 `boil:"creator_id" json:"creator_id" toml:"creator_id" yaml:"creator_id"`
	ParentID        string                 `boil:"parent_id" json:"parent_id" toml:"parent_id" yaml:"parent_id"`
	CreatedAt       int64                  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt       int64                  `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	DeleteAt        model_types.NullInt64  `boil:"delete_at" json:"delete_at,omitempty" toml:"delete_at" yaml:"delete_at,omitempty"`
	Path            string                 `boil:"path" json:"path" toml:"path" yaml:"path"`
	ThumbnailPath   string                 `boil:"thumbnail_path" json:"thumbnail_path" toml:"thumbnail_path" yaml:"thumbnail_path"`
	PreviewPath     string                 `boil:"preview_path" json:"preview_path" toml:"preview_path" yaml:"preview_path"`
	Name            string                 `boil:"name" json:"name" toml:"name" yaml:"name"`
	Extension       string                 `boil:"extension" json:"extension" toml:"extension" yaml:"extension"`
	Size            int64                  `boil:"size" json:"size" toml:"size" yaml:"size"`
	MimeType        string                 `boil:"mime_type" json:"mime_type" toml:"mime_type" yaml:"mime_type"`
	Width           model_types.NullInt    `boil:"width" json:"width,omitempty" toml:"width" yaml:"width,omitempty"`
	Height          model_types.NullInt    `boil:"height" json:"height,omitempty" toml:"height" yaml:"height,omitempty"`
	HasPreviewImage bool                   `boil:"has_preview_image" json:"has_preview_image" toml:"has_preview_image" yaml:"has_preview_image"`
	MiniPreview     null.Bytes             `boil:"mini_preview" json:"mini_preview,omitempty" toml:"mini_preview" yaml:"mini_preview,omitempty"`
	Content         string                 `boil:"content" json:"content" toml:"content" yaml:"content"`
	RemoteID        model_types.NullString `boil:"remote_id" json:"remote_id,omitempty" toml:"remote_id" yaml:"remote_id,omitempty"`

	R *fileInfoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fileInfoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var FileInfoColumns = struct {
	ID              string
	CreatorID       string
	ParentID        string
	CreatedAt       string
	UpdatedAt       string
	DeleteAt        string
	Path            string
	ThumbnailPath   string
	PreviewPath     string
	Name            string
	Extension       string
	Size            string
	MimeType        string
	Width           string
	Height          string
	HasPreviewImage string
	MiniPreview     string
	Content         string
	RemoteID        string
}{
	ID:              "id",
	CreatorID:       "creator_id",
	ParentID:        "parent_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeleteAt:        "delete_at",
	Path:            "path",
	ThumbnailPath:   "thumbnail_path",
	PreviewPath:     "preview_path",
	Name:            "name",
	Extension:       "extension",
	Size:            "size",
	MimeType:        "mime_type",
	Width:           "width",
	Height:          "height",
	HasPreviewImage: "has_preview_image",
	MiniPreview:     "mini_preview",
	Content:         "content",
	RemoteID:        "remote_id",
}

var FileInfoTableColumns = struct {
	ID              string
	CreatorID       string
	ParentID        string
	CreatedAt       string
	UpdatedAt       string
	DeleteAt        string
	Path            string
	ThumbnailPath   string
	PreviewPath     string
	Name            string
	Extension       string
	Size            string
	MimeType        string
	Width           string
	Height          string
	HasPreviewImage string
	MiniPreview     string
	Content         string
	RemoteID        string
}{
	ID:              "file_infos.id",
	CreatorID:       "file_infos.creator_id",
	ParentID:        "file_infos.parent_id",
	CreatedAt:       "file_infos.created_at",
	UpdatedAt:       "file_infos.updated_at",
	DeleteAt:        "file_infos.delete_at",
	Path:            "file_infos.path",
	ThumbnailPath:   "file_infos.thumbnail_path",
	PreviewPath:     "file_infos.preview_path",
	Name:            "file_infos.name",
	Extension:       "file_infos.extension",
	Size:            "file_infos.size",
	MimeType:        "file_infos.mime_type",
	Width:           "file_infos.width",
	Height:          "file_infos.height",
	HasPreviewImage: "file_infos.has_preview_image",
	MiniPreview:     "file_infos.mini_preview",
	Content:         "file_infos.content",
	RemoteID:        "file_infos.remote_id",
}

// Generated where

type whereHelpermodel_types_NullInt64 struct{ field string }

func (w whereHelpermodel_types_NullInt64) EQ(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpermodel_types_NullInt64) NEQ(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpermodel_types_NullInt64) LT(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpermodel_types_NullInt64) LTE(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpermodel_types_NullInt64) GT(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpermodel_types_NullInt64) GTE(x model_types.NullInt64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpermodel_types_NullInt64) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpermodel_types_NullInt64) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

type whereHelpernull_Bytes struct{ field string }

func (w whereHelpernull_Bytes) EQ(x null.Bytes) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bytes) NEQ(x null.Bytes) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bytes) LT(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bytes) LTE(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bytes) GT(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bytes) GTE(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Bytes) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bytes) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var FileInfoWhere = struct {
	ID              whereHelperstring
	CreatorID       whereHelperstring
	ParentID        whereHelperstring
	CreatedAt       whereHelperint64
	UpdatedAt       whereHelperint64
	DeleteAt        whereHelpermodel_types_NullInt64
	Path            whereHelperstring
	ThumbnailPath   whereHelperstring
	PreviewPath     whereHelperstring
	Name            whereHelperstring
	Extension       whereHelperstring
	Size            whereHelperint64
	MimeType        whereHelperstring
	Width           whereHelpermodel_types_NullInt
	Height          whereHelpermodel_types_NullInt
	HasPreviewImage whereHelperbool
	MiniPreview     whereHelpernull_Bytes
	Content         whereHelperstring
	RemoteID        whereHelpermodel_types_NullString
}{
	ID:              whereHelperstring{field: "\"file_infos\".\"id\""},
	CreatorID:       whereHelperstring{field: "\"file_infos\".\"creator_id\""},
	ParentID:        whereHelperstring{field: "\"file_infos\".\"parent_id\""},
	CreatedAt:       whereHelperint64{field: "\"file_infos\".\"created_at\""},
	UpdatedAt:       whereHelperint64{field: "\"file_infos\".\"updated_at\""},
	DeleteAt:        whereHelpermodel_types_NullInt64{field: "\"file_infos\".\"delete_at\""},
	Path:            whereHelperstring{field: "\"file_infos\".\"path\""},
	ThumbnailPath:   whereHelperstring{field: "\"file_infos\".\"thumbnail_path\""},
	PreviewPath:     whereHelperstring{field: "\"file_infos\".\"preview_path\""},
	Name:            whereHelperstring{field: "\"file_infos\".\"name\""},
	Extension:       whereHelperstring{field: "\"file_infos\".\"extension\""},
	Size:            whereHelperint64{field: "\"file_infos\".\"size\""},
	MimeType:        whereHelperstring{field: "\"file_infos\".\"mime_type\""},
	Width:           whereHelpermodel_types_NullInt{field: "\"file_infos\".\"width\""},
	Height:          whereHelpermodel_types_NullInt{field: "\"file_infos\".\"height\""},
	HasPreviewImage: whereHelperbool{field: "\"file_infos\".\"has_preview_image\""},
	MiniPreview:     whereHelpernull_Bytes{field: "\"file_infos\".\"mini_preview\""},
	Content:         whereHelperstring{field: "\"file_infos\".\"content\""},
	RemoteID:        whereHelpermodel_types_NullString{field: "\"file_infos\".\"remote_id\""},
}

// FileInfoRels is where relationship names are stored.
var FileInfoRels = struct {
}{}

// fileInfoR is where relationships are stored.
type fileInfoR struct {
}

// NewStruct creates a new relationship struct
func (*fileInfoR) NewStruct() *fileInfoR {
	return &fileInfoR{}
}

// fileInfoL is where Load methods for each relationship are stored.
type fileInfoL struct{}

var (
	fileInfoAllColumns            = []string{"id", "creator_id", "parent_id", "created_at", "updated_at", "delete_at", "path", "thumbnail_path", "preview_path", "name", "extension", "size", "mime_type", "width", "height", "has_preview_image", "mini_preview", "content", "remote_id"}
	fileInfoColumnsWithoutDefault = []string{"id", "creator_id", "parent_id", "created_at", "updated_at", "path", "thumbnail_path", "preview_path", "name", "extension", "size", "mime_type", "has_preview_image", "content"}
	fileInfoColumnsWithDefault    = []string{"delete_at", "width", "height", "mini_preview", "remote_id"}
	fileInfoPrimaryKeyColumns     = []string{"id"}
	fileInfoGeneratedColumns      = []string{}
)

type (
	// FileInfoSlice is an alias for a slice of pointers to FileInfo.
	// This should almost always be used instead of []FileInfo.
	FileInfoSlice []*FileInfo

	fileInfoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fileInfoType                 = reflect.TypeOf(&FileInfo{})
	fileInfoMapping              = queries.MakeStructMapping(fileInfoType)
	fileInfoPrimaryKeyMapping, _ = queries.BindMapping(fileInfoType, fileInfoMapping, fileInfoPrimaryKeyColumns)
	fileInfoInsertCacheMut       sync.RWMutex
	fileInfoInsertCache          = make(map[string]insertCache)
	fileInfoUpdateCacheMut       sync.RWMutex
	fileInfoUpdateCache          = make(map[string]updateCache)
	fileInfoUpsertCacheMut       sync.RWMutex
	fileInfoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single fileInfo record from the query.
func (q fileInfoQuery) One(exec boil.Executor) (*FileInfo, error) {
	o := &FileInfo{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for file_infos")
	}

	return o, nil
}

// All returns all FileInfo records from the query.
func (q fileInfoQuery) All(exec boil.Executor) (FileInfoSlice, error) {
	var o []*FileInfo

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to FileInfo slice")
	}

	return o, nil
}

// Count returns the count of all FileInfo records in the query.
func (q fileInfoQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count file_infos rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q fileInfoQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if file_infos exists")
	}

	return count > 0, nil
}

// FileInfos retrieves all the records using an executor.
func FileInfos(mods ...qm.QueryMod) fileInfoQuery {
	mods = append(mods, qm.From("\"file_infos\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"file_infos\".*"})
	}

	return fileInfoQuery{q}
}

// FindFileInfo retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFileInfo(exec boil.Executor, iD string, selectCols ...string) (*FileInfo, error) {
	fileInfoObj := &FileInfo{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"file_infos\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, fileInfoObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from file_infos")
	}

	return fileInfoObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *FileInfo) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no file_infos provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(fileInfoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	fileInfoInsertCacheMut.RLock()
	cache, cached := fileInfoInsertCache[key]
	fileInfoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			fileInfoAllColumns,
			fileInfoColumnsWithDefault,
			fileInfoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(fileInfoType, fileInfoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fileInfoType, fileInfoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"file_infos\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"file_infos\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into file_infos")
	}

	if !cached {
		fileInfoInsertCacheMut.Lock()
		fileInfoInsertCache[key] = cache
		fileInfoInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the FileInfo.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *FileInfo) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	fileInfoUpdateCacheMut.RLock()
	cache, cached := fileInfoUpdateCache[key]
	fileInfoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			fileInfoAllColumns,
			fileInfoPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update file_infos, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"file_infos\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, fileInfoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fileInfoType, fileInfoMapping, append(wl, fileInfoPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update file_infos row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for file_infos")
	}

	if !cached {
		fileInfoUpdateCacheMut.Lock()
		fileInfoUpdateCache[key] = cache
		fileInfoUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q fileInfoQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for file_infos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for file_infos")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FileInfoSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"file_infos\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, fileInfoPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in fileInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all fileInfo")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *FileInfo) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no file_infos provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(fileInfoColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	fileInfoUpsertCacheMut.RLock()
	cache, cached := fileInfoUpsertCache[key]
	fileInfoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			fileInfoAllColumns,
			fileInfoColumnsWithDefault,
			fileInfoColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			fileInfoAllColumns,
			fileInfoPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert file_infos, could not build update column list")
		}

		ret := strmangle.SetComplement(fileInfoAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(fileInfoPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert file_infos, could not build conflict column list")
			}

			conflict = make([]string, len(fileInfoPrimaryKeyColumns))
			copy(conflict, fileInfoPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"file_infos\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(fileInfoType, fileInfoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fileInfoType, fileInfoMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert file_infos")
	}

	if !cached {
		fileInfoUpsertCacheMut.Lock()
		fileInfoUpsertCache[key] = cache
		fileInfoUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single FileInfo record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *FileInfo) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no FileInfo provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), fileInfoPrimaryKeyMapping)
	sql := "DELETE FROM \"file_infos\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from file_infos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for file_infos")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q fileInfoQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no fileInfoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from file_infos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for file_infos")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FileInfoSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"file_infos\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, fileInfoPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from fileInfo slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for file_infos")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *FileInfo) Reload(exec boil.Executor) error {
	ret, err := FindFileInfo(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileInfoSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FileInfoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileInfoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"file_infos\".* FROM \"file_infos\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, fileInfoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in FileInfoSlice")
	}

	*o = slice

	return nil
}

// FileInfoExists checks if the FileInfo row exists.
func FileInfoExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"file_infos\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if file_infos exists")
	}

	return exists, nil
}

// Exists checks if the FileInfo row exists.
func (o *FileInfo) Exists(exec boil.Executor) (bool, error) {
	return FileInfoExists(exec, o.ID)
}
