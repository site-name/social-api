// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ExportEvent is an object representing the database table.
type ExportEvent struct {
	ID           string          `boil:"id" json:"id" toml:"id" yaml:"id"`
	Date         int64           `boil:"date" json:"date" toml:"date" yaml:"date"`
	Type         Exporteventtype `boil:"type" json:"type" toml:"type" yaml:"type"`
	Parameters   null.String     `boil:"parameters" json:"parameters,omitempty" toml:"parameters" yaml:"parameters,omitempty"`
	ExportFileID string          `boil:"export_file_id" json:"export_file_id" toml:"export_file_id" yaml:"export_file_id"`
	UserID       null.String     `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`

	R *exportEventR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L exportEventL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ExportEventColumns = struct {
	ID           string
	Date         string
	Type         string
	Parameters   string
	ExportFileID string
	UserID       string
}{
	ID:           "id",
	Date:         "date",
	Type:         "type",
	Parameters:   "parameters",
	ExportFileID: "export_file_id",
	UserID:       "user_id",
}

var ExportEventTableColumns = struct {
	ID           string
	Date         string
	Type         string
	Parameters   string
	ExportFileID string
	UserID       string
}{
	ID:           "export_events.id",
	Date:         "export_events.date",
	Type:         "export_events.type",
	Parameters:   "export_events.parameters",
	ExportFileID: "export_events.export_file_id",
	UserID:       "export_events.user_id",
}

// Generated where

type whereHelperExporteventtype struct{ field string }

func (w whereHelperExporteventtype) EQ(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelperExporteventtype) NEQ(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperExporteventtype) LT(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelperExporteventtype) LTE(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperExporteventtype) GT(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelperExporteventtype) GTE(x Exporteventtype) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperExporteventtype) IN(slice []Exporteventtype) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperExporteventtype) NIN(slice []Exporteventtype) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var ExportEventWhere = struct {
	ID           whereHelperstring
	Date         whereHelperint64
	Type         whereHelperExporteventtype
	Parameters   whereHelpernull_String
	ExportFileID whereHelperstring
	UserID       whereHelpernull_String
}{
	ID:           whereHelperstring{field: "\"export_events\".\"id\""},
	Date:         whereHelperint64{field: "\"export_events\".\"date\""},
	Type:         whereHelperExporteventtype{field: "\"export_events\".\"type\""},
	Parameters:   whereHelpernull_String{field: "\"export_events\".\"parameters\""},
	ExportFileID: whereHelperstring{field: "\"export_events\".\"export_file_id\""},
	UserID:       whereHelpernull_String{field: "\"export_events\".\"user_id\""},
}

// ExportEventRels is where relationship names are stored.
var ExportEventRels = struct {
	ExportFile string
	User       string
}{
	ExportFile: "ExportFile",
	User:       "User",
}

// exportEventR is where relationships are stored.
type exportEventR struct {
	ExportFile *ExportFile `boil:"ExportFile" json:"ExportFile" toml:"ExportFile" yaml:"ExportFile"`
	User       *User       `boil:"User" json:"User" toml:"User" yaml:"User"`
}

// NewStruct creates a new relationship struct
func (*exportEventR) NewStruct() *exportEventR {
	return &exportEventR{}
}

func (r *exportEventR) GetExportFile() *ExportFile {
	if r == nil {
		return nil
	}
	return r.ExportFile
}

func (r *exportEventR) GetUser() *User {
	if r == nil {
		return nil
	}
	return r.User
}

// exportEventL is where Load methods for each relationship are stored.
type exportEventL struct{}

var (
	exportEventAllColumns            = []string{"id", "date", "type", "parameters", "export_file_id", "user_id"}
	exportEventColumnsWithoutDefault = []string{"date", "type", "export_file_id"}
	exportEventColumnsWithDefault    = []string{"id", "parameters", "user_id"}
	exportEventPrimaryKeyColumns     = []string{"id"}
	exportEventGeneratedColumns      = []string{}
)

type (
	// ExportEventSlice is an alias for a slice of pointers to ExportEvent.
	// This should almost always be used instead of []ExportEvent.
	ExportEventSlice []*ExportEvent

	exportEventQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	exportEventType                 = reflect.TypeOf(&ExportEvent{})
	exportEventMapping              = queries.MakeStructMapping(exportEventType)
	exportEventPrimaryKeyMapping, _ = queries.BindMapping(exportEventType, exportEventMapping, exportEventPrimaryKeyColumns)
	exportEventInsertCacheMut       sync.RWMutex
	exportEventInsertCache          = make(map[string]insertCache)
	exportEventUpdateCacheMut       sync.RWMutex
	exportEventUpdateCache          = make(map[string]updateCache)
	exportEventUpsertCacheMut       sync.RWMutex
	exportEventUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single exportEvent record from the query.
func (q exportEventQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ExportEvent, error) {
	o := &ExportEvent{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for export_events")
	}

	return o, nil
}

// All returns all ExportEvent records from the query.
func (q exportEventQuery) All(ctx context.Context, exec boil.ContextExecutor) (ExportEventSlice, error) {
	var o []*ExportEvent

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to ExportEvent slice")
	}

	return o, nil
}

// Count returns the count of all ExportEvent records in the query.
func (q exportEventQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count export_events rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q exportEventQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if export_events exists")
	}

	return count > 0, nil
}

// ExportFile pointed to by the foreign key.
func (o *ExportEvent) ExportFile(mods ...qm.QueryMod) exportFileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ExportFileID),
	}

	queryMods = append(queryMods, mods...)

	return ExportFiles(queryMods...)
}

// User pointed to by the foreign key.
func (o *ExportEvent) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	return Users(queryMods...)
}

// LoadExportFile allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (exportEventL) LoadExportFile(ctx context.Context, e boil.ContextExecutor, singular bool, maybeExportEvent interface{}, mods queries.Applicator) error {
	var slice []*ExportEvent
	var object *ExportEvent

	if singular {
		var ok bool
		object, ok = maybeExportEvent.(*ExportEvent)
		if !ok {
			object = new(ExportEvent)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeExportEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeExportEvent))
			}
		}
	} else {
		s, ok := maybeExportEvent.(*[]*ExportEvent)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeExportEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeExportEvent))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &exportEventR{}
		}
		args = append(args, object.ExportFileID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &exportEventR{}
			}

			for _, a := range args {
				if a == obj.ExportFileID {
					continue Outer
				}
			}

			args = append(args, obj.ExportFileID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`export_files`),
		qm.WhereIn(`export_files.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ExportFile")
	}

	var resultSlice []*ExportFile
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ExportFile")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for export_files")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for export_files")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ExportFile = foreign
		if foreign.R == nil {
			foreign.R = &exportFileR{}
		}
		foreign.R.ExportEvents = append(foreign.R.ExportEvents, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ExportFileID == foreign.ID {
				local.R.ExportFile = foreign
				if foreign.R == nil {
					foreign.R = &exportFileR{}
				}
				foreign.R.ExportEvents = append(foreign.R.ExportEvents, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (exportEventL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeExportEvent interface{}, mods queries.Applicator) error {
	var slice []*ExportEvent
	var object *ExportEvent

	if singular {
		var ok bool
		object, ok = maybeExportEvent.(*ExportEvent)
		if !ok {
			object = new(ExportEvent)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeExportEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeExportEvent))
			}
		}
	} else {
		s, ok := maybeExportEvent.(*[]*ExportEvent)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeExportEvent)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeExportEvent))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &exportEventR{}
		}
		if !queries.IsNil(object.UserID) {
			args = append(args, object.UserID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &exportEventR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.UserID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.UserID) {
				args = append(args, obj.UserID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`users`),
		qm.WhereIn(`users.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.ExportEvents = append(foreign.R.ExportEvents, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.UserID, foreign.ID) {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.ExportEvents = append(foreign.R.ExportEvents, local)
				break
			}
		}
	}

	return nil
}

// SetExportFile of the exportEvent to the related item.
// Sets o.R.ExportFile to related.
// Adds o to related.R.ExportEvents.
func (o *ExportEvent) SetExportFile(ctx context.Context, exec boil.ContextExecutor, insert bool, related *ExportFile) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"export_events\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"export_file_id"}),
		strmangle.WhereClause("\"", "\"", 2, exportEventPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ExportFileID = related.ID
	if o.R == nil {
		o.R = &exportEventR{
			ExportFile: related,
		}
	} else {
		o.R.ExportFile = related
	}

	if related.R == nil {
		related.R = &exportFileR{
			ExportEvents: ExportEventSlice{o},
		}
	} else {
		related.R.ExportEvents = append(related.R.ExportEvents, o)
	}

	return nil
}

// SetUser of the exportEvent to the related item.
// Sets o.R.User to related.
// Adds o to related.R.ExportEvents.
func (o *ExportEvent) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"export_events\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, exportEventPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.UserID, related.ID)
	if o.R == nil {
		o.R = &exportEventR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			ExportEvents: ExportEventSlice{o},
		}
	} else {
		related.R.ExportEvents = append(related.R.ExportEvents, o)
	}

	return nil
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct.
func (o *ExportEvent) RemoveUser(ctx context.Context, exec boil.ContextExecutor, related *User) error {
	var err error

	queries.SetScanner(&o.UserID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("user_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.User = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.ExportEvents {
		if queries.Equal(o.UserID, ri.UserID) {
			continue
		}

		ln := len(related.R.ExportEvents)
		if ln > 1 && i < ln-1 {
			related.R.ExportEvents[i] = related.R.ExportEvents[ln-1]
		}
		related.R.ExportEvents = related.R.ExportEvents[:ln-1]
		break
	}
	return nil
}

// ExportEvents retrieves all the records using an executor.
func ExportEvents(mods ...qm.QueryMod) exportEventQuery {
	mods = append(mods, qm.From("\"export_events\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"export_events\".*"})
	}

	return exportEventQuery{q}
}

// FindExportEvent retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindExportEvent(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ExportEvent, error) {
	exportEventObj := &ExportEvent{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"export_events\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, exportEventObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from export_events")
	}

	return exportEventObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ExportEvent) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no export_events provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(exportEventColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	exportEventInsertCacheMut.RLock()
	cache, cached := exportEventInsertCache[key]
	exportEventInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			exportEventAllColumns,
			exportEventColumnsWithDefault,
			exportEventColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(exportEventType, exportEventMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(exportEventType, exportEventMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"export_events\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"export_events\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into export_events")
	}

	if !cached {
		exportEventInsertCacheMut.Lock()
		exportEventInsertCache[key] = cache
		exportEventInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ExportEvent.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ExportEvent) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	exportEventUpdateCacheMut.RLock()
	cache, cached := exportEventUpdateCache[key]
	exportEventUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			exportEventAllColumns,
			exportEventPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update export_events, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"export_events\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, exportEventPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(exportEventType, exportEventMapping, append(wl, exportEventPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update export_events row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for export_events")
	}

	if !cached {
		exportEventUpdateCacheMut.Lock()
		exportEventUpdateCache[key] = cache
		exportEventUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q exportEventQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for export_events")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for export_events")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ExportEventSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exportEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"export_events\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, exportEventPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in exportEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all exportEvent")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ExportEvent) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no export_events provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(exportEventColumnsWithDefault, o)

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

	exportEventUpsertCacheMut.RLock()
	cache, cached := exportEventUpsertCache[key]
	exportEventUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			exportEventAllColumns,
			exportEventColumnsWithDefault,
			exportEventColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			exportEventAllColumns,
			exportEventPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert export_events, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(exportEventPrimaryKeyColumns))
			copy(conflict, exportEventPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"export_events\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(exportEventType, exportEventMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(exportEventType, exportEventMapping, ret)
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert export_events")
	}

	if !cached {
		exportEventUpsertCacheMut.Lock()
		exportEventUpsertCache[key] = cache
		exportEventUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ExportEvent record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ExportEvent) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no ExportEvent provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), exportEventPrimaryKeyMapping)
	sql := "DELETE FROM \"export_events\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from export_events")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for export_events")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q exportEventQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no exportEventQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from export_events")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for export_events")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ExportEventSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exportEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"export_events\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, exportEventPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from exportEvent slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for export_events")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ExportEvent) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindExportEvent(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ExportEventSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ExportEventSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exportEventPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"export_events\".* FROM \"export_events\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, exportEventPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ExportEventSlice")
	}

	*o = slice

	return nil
}

// ExportEventExists checks if the ExportEvent row exists.
func ExportEventExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"export_events\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if export_events exists")
	}

	return exists, nil
}

// Exists checks if the ExportEvent row exists.
func (o *ExportEvent) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ExportEventExists(ctx, exec, o.ID)
}