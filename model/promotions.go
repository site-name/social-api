// Code generated by SQLBoiler 4.17.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Promotion is an object representing the database table.
type Promotion struct {
	ID                          string                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name                        string                 `boil:"name" json:"name" toml:"name" yaml:"name"`
	Type                        PromotionType          `boil:"type" json:"type" toml:"type" yaml:"type"`
	Description                 model_types.JSONString `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	StartDate                   int64                  `boil:"start_date" json:"start_date" toml:"start_date" yaml:"start_date"`
	EndDate                     model_types.NullInt64  `boil:"end_date" json:"end_date,omitempty" toml:"end_date" yaml:"end_date,omitempty"`
	CreatedAt                   int64                  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt                   model_types.NullInt64  `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	LastModificationScheduledAt model_types.NullInt64  `boil:"last_modification_scheduled_at" json:"last_modification_scheduled_at,omitempty" toml:"last_modification_scheduled_at" yaml:"last_modification_scheduled_at,omitempty"`
	Metadata                    model_types.JSONString `boil:"metadata" json:"metadata,omitempty" toml:"metadata" yaml:"metadata,omitempty"`
	PrivateMetadata             model_types.JSONString `boil:"private_metadata" json:"private_metadata,omitempty" toml:"private_metadata" yaml:"private_metadata,omitempty"`

	R *promotionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L promotionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PromotionColumns = struct {
	ID                          string
	Name                        string
	Type                        string
	Description                 string
	StartDate                   string
	EndDate                     string
	CreatedAt                   string
	UpdatedAt                   string
	LastModificationScheduledAt string
	Metadata                    string
	PrivateMetadata             string
}{
	ID:                          "id",
	Name:                        "name",
	Type:                        "type",
	Description:                 "description",
	StartDate:                   "start_date",
	EndDate:                     "end_date",
	CreatedAt:                   "created_at",
	UpdatedAt:                   "updated_at",
	LastModificationScheduledAt: "last_modification_scheduled_at",
	Metadata:                    "metadata",
	PrivateMetadata:             "private_metadata",
}

var PromotionTableColumns = struct {
	ID                          string
	Name                        string
	Type                        string
	Description                 string
	StartDate                   string
	EndDate                     string
	CreatedAt                   string
	UpdatedAt                   string
	LastModificationScheduledAt string
	Metadata                    string
	PrivateMetadata             string
}{
	ID:                          "promotions.id",
	Name:                        "promotions.name",
	Type:                        "promotions.type",
	Description:                 "promotions.description",
	StartDate:                   "promotions.start_date",
	EndDate:                     "promotions.end_date",
	CreatedAt:                   "promotions.created_at",
	UpdatedAt:                   "promotions.updated_at",
	LastModificationScheduledAt: "promotions.last_modification_scheduled_at",
	Metadata:                    "promotions.metadata",
	PrivateMetadata:             "promotions.private_metadata",
}

// Generated where

type whereHelperPromotionType struct{ field string }

func (w whereHelperPromotionType) EQ(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelperPromotionType) NEQ(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperPromotionType) LT(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelperPromotionType) LTE(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperPromotionType) GT(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelperPromotionType) GTE(x PromotionType) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelperPromotionType) IN(slice []PromotionType) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperPromotionType) NIN(slice []PromotionType) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var PromotionWhere = struct {
	ID                          whereHelperstring
	Name                        whereHelperstring
	Type                        whereHelperPromotionType
	Description                 whereHelpermodel_types_JSONString
	StartDate                   whereHelperint64
	EndDate                     whereHelpermodel_types_NullInt64
	CreatedAt                   whereHelperint64
	UpdatedAt                   whereHelpermodel_types_NullInt64
	LastModificationScheduledAt whereHelpermodel_types_NullInt64
	Metadata                    whereHelpermodel_types_JSONString
	PrivateMetadata             whereHelpermodel_types_JSONString
}{
	ID:                          whereHelperstring{field: "\"promotions\".\"id\""},
	Name:                        whereHelperstring{field: "\"promotions\".\"name\""},
	Type:                        whereHelperPromotionType{field: "\"promotions\".\"type\""},
	Description:                 whereHelpermodel_types_JSONString{field: "\"promotions\".\"description\""},
	StartDate:                   whereHelperint64{field: "\"promotions\".\"start_date\""},
	EndDate:                     whereHelpermodel_types_NullInt64{field: "\"promotions\".\"end_date\""},
	CreatedAt:                   whereHelperint64{field: "\"promotions\".\"created_at\""},
	UpdatedAt:                   whereHelpermodel_types_NullInt64{field: "\"promotions\".\"updated_at\""},
	LastModificationScheduledAt: whereHelpermodel_types_NullInt64{field: "\"promotions\".\"last_modification_scheduled_at\""},
	Metadata:                    whereHelpermodel_types_JSONString{field: "\"promotions\".\"metadata\""},
	PrivateMetadata:             whereHelpermodel_types_JSONString{field: "\"promotions\".\"private_metadata\""},
}

// PromotionRels is where relationship names are stored.
var PromotionRels = struct {
	PromotionEvents string
	PromotionRules  string
}{
	PromotionEvents: "PromotionEvents",
	PromotionRules:  "PromotionRules",
}

// promotionR is where relationships are stored.
type promotionR struct {
	PromotionEvents PromotionEventSlice `boil:"PromotionEvents" json:"PromotionEvents" toml:"PromotionEvents" yaml:"PromotionEvents"`
	PromotionRules  PromotionRuleSlice  `boil:"PromotionRules" json:"PromotionRules" toml:"PromotionRules" yaml:"PromotionRules"`
}

// NewStruct creates a new relationship struct
func (*promotionR) NewStruct() *promotionR {
	return &promotionR{}
}

func (r *promotionR) GetPromotionEvents() PromotionEventSlice {
	if r == nil {
		return nil
	}
	return r.PromotionEvents
}

func (r *promotionR) GetPromotionRules() PromotionRuleSlice {
	if r == nil {
		return nil
	}
	return r.PromotionRules
}

// promotionL is where Load methods for each relationship are stored.
type promotionL struct{}

var (
	promotionAllColumns            = []string{"id", "name", "type", "description", "start_date", "end_date", "created_at", "updated_at", "last_modification_scheduled_at", "metadata", "private_metadata"}
	promotionColumnsWithoutDefault = []string{"id", "name", "start_date", "created_at"}
	promotionColumnsWithDefault    = []string{"type", "description", "end_date", "updated_at", "last_modification_scheduled_at", "metadata", "private_metadata"}
	promotionPrimaryKeyColumns     = []string{"id"}
	promotionGeneratedColumns      = []string{}
)

type (
	// PromotionSlice is an alias for a slice of pointers to Promotion.
	// This should almost always be used instead of []Promotion.
	PromotionSlice []*Promotion

	promotionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	promotionType                 = reflect.TypeOf(&Promotion{})
	promotionMapping              = queries.MakeStructMapping(promotionType)
	promotionPrimaryKeyMapping, _ = queries.BindMapping(promotionType, promotionMapping, promotionPrimaryKeyColumns)
	promotionInsertCacheMut       sync.RWMutex
	promotionInsertCache          = make(map[string]insertCache)
	promotionUpdateCacheMut       sync.RWMutex
	promotionUpdateCache          = make(map[string]updateCache)
	promotionUpsertCacheMut       sync.RWMutex
	promotionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single promotion record from the query.
func (q promotionQuery) One(exec boil.Executor) (*Promotion, error) {
	o := &Promotion{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for promotions")
	}

	return o, nil
}

// All returns all Promotion records from the query.
func (q promotionQuery) All(exec boil.Executor) (PromotionSlice, error) {
	var o []*Promotion

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Promotion slice")
	}

	return o, nil
}

// Count returns the count of all Promotion records in the query.
func (q promotionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count promotions rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q promotionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if promotions exists")
	}

	return count > 0, nil
}

// PromotionEvents retrieves all the promotion_event's PromotionEvents with an executor.
func (o *Promotion) PromotionEvents(mods ...qm.QueryMod) promotionEventQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"promotion_events\".\"promotion_id\"=?", o.ID),
	)

	return PromotionEvents(queryMods...)
}

// PromotionRules retrieves all the promotion_rule's PromotionRules with an executor.
func (o *Promotion) PromotionRules(mods ...qm.QueryMod) promotionRuleQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"promotion_rules\".\"promotion_id\"=?", o.ID),
	)

	return PromotionRules(queryMods...)
}

// LoadPromotionEvents allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (promotionL) LoadPromotionEvents(e boil.Executor, singular bool, maybePromotion interface{}, mods queries.Applicator) error {
	var slice []*Promotion
	var object *Promotion

	if singular {
		var ok bool
		object, ok = maybePromotion.(*Promotion)
		if !ok {
			object = new(Promotion)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePromotion)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePromotion))
			}
		}
	} else {
		s, ok := maybePromotion.(*[]*Promotion)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePromotion)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePromotion))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &promotionR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &promotionR{}
			}
			args[obj.ID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`promotion_events`),
		qm.WhereIn(`promotion_events.promotion_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load promotion_events")
	}

	var resultSlice []*PromotionEvent
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice promotion_events")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on promotion_events")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for promotion_events")
	}

	if singular {
		object.R.PromotionEvents = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &promotionEventR{}
			}
			foreign.R.Promotion = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.PromotionID) {
				local.R.PromotionEvents = append(local.R.PromotionEvents, foreign)
				if foreign.R == nil {
					foreign.R = &promotionEventR{}
				}
				foreign.R.Promotion = local
				break
			}
		}
	}

	return nil
}

// LoadPromotionRules allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (promotionL) LoadPromotionRules(e boil.Executor, singular bool, maybePromotion interface{}, mods queries.Applicator) error {
	var slice []*Promotion
	var object *Promotion

	if singular {
		var ok bool
		object, ok = maybePromotion.(*Promotion)
		if !ok {
			object = new(Promotion)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePromotion)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePromotion))
			}
		}
	} else {
		s, ok := maybePromotion.(*[]*Promotion)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePromotion)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePromotion))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &promotionR{}
		}
		args[object.ID] = struct{}{}
	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &promotionR{}
			}
			args[obj.ID] = struct{}{}
		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`promotion_rules`),
		qm.WhereIn(`promotion_rules.promotion_id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load promotion_rules")
	}

	var resultSlice []*PromotionRule
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice promotion_rules")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on promotion_rules")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for promotion_rules")
	}

	if singular {
		object.R.PromotionRules = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &promotionRuleR{}
			}
			foreign.R.Promotion = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PromotionID {
				local.R.PromotionRules = append(local.R.PromotionRules, foreign)
				if foreign.R == nil {
					foreign.R = &promotionRuleR{}
				}
				foreign.R.Promotion = local
				break
			}
		}
	}

	return nil
}

// AddPromotionEvents adds the given related objects to the existing relationships
// of the promotion, optionally inserting them as new records.
// Appends related to o.R.PromotionEvents.
// Sets related.R.Promotion appropriately.
func (o *Promotion) AddPromotionEvents(exec boil.Executor, insert bool, related ...*PromotionEvent) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.PromotionID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"promotion_events\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"promotion_id"}),
				strmangle.WhereClause("\"", "\"", 2, promotionEventPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.PromotionID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &promotionR{
			PromotionEvents: related,
		}
	} else {
		o.R.PromotionEvents = append(o.R.PromotionEvents, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &promotionEventR{
				Promotion: o,
			}
		} else {
			rel.R.Promotion = o
		}
	}
	return nil
}

// SetPromotionEvents removes all previously related items of the
// promotion replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Promotion's PromotionEvents accordingly.
// Replaces o.R.PromotionEvents with related.
// Sets related.R.Promotion's PromotionEvents accordingly.
func (o *Promotion) SetPromotionEvents(exec boil.Executor, insert bool, related ...*PromotionEvent) error {
	query := "update \"promotion_events\" set \"promotion_id\" = null where \"promotion_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.PromotionEvents {
			queries.SetScanner(&rel.PromotionID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Promotion = nil
		}
		o.R.PromotionEvents = nil
	}

	return o.AddPromotionEvents(exec, insert, related...)
}

// RemovePromotionEvents relationships from objects passed in.
// Removes related items from R.PromotionEvents (uses pointer comparison, removal does not keep order)
// Sets related.R.Promotion.
func (o *Promotion) RemovePromotionEvents(exec boil.Executor, related ...*PromotionEvent) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.PromotionID, nil)
		if rel.R != nil {
			rel.R.Promotion = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("promotion_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.PromotionEvents {
			if rel != ri {
				continue
			}

			ln := len(o.R.PromotionEvents)
			if ln > 1 && i < ln-1 {
				o.R.PromotionEvents[i] = o.R.PromotionEvents[ln-1]
			}
			o.R.PromotionEvents = o.R.PromotionEvents[:ln-1]
			break
		}
	}

	return nil
}

// AddPromotionRules adds the given related objects to the existing relationships
// of the promotion, optionally inserting them as new records.
// Appends related to o.R.PromotionRules.
// Sets related.R.Promotion appropriately.
func (o *Promotion) AddPromotionRules(exec boil.Executor, insert bool, related ...*PromotionRule) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PromotionID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"promotion_rules\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"promotion_id"}),
				strmangle.WhereClause("\"", "\"", 2, promotionRulePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PromotionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &promotionR{
			PromotionRules: related,
		}
	} else {
		o.R.PromotionRules = append(o.R.PromotionRules, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &promotionRuleR{
				Promotion: o,
			}
		} else {
			rel.R.Promotion = o
		}
	}
	return nil
}

// Promotions retrieves all the records using an executor.
func Promotions(mods ...qm.QueryMod) promotionQuery {
	mods = append(mods, qm.From("\"promotions\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"promotions\".*"})
	}

	return promotionQuery{q}
}

// FindPromotion retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPromotion(exec boil.Executor, iD string, selectCols ...string) (*Promotion, error) {
	promotionObj := &Promotion{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"promotions\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, promotionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from promotions")
	}

	return promotionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Promotion) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no promotions provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(promotionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	promotionInsertCacheMut.RLock()
	cache, cached := promotionInsertCache[key]
	promotionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			promotionAllColumns,
			promotionColumnsWithDefault,
			promotionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(promotionType, promotionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(promotionType, promotionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"promotions\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"promotions\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "model: unable to insert into promotions")
	}

	if !cached {
		promotionInsertCacheMut.Lock()
		promotionInsertCache[key] = cache
		promotionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Promotion.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Promotion) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	promotionUpdateCacheMut.RLock()
	cache, cached := promotionUpdateCache[key]
	promotionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			promotionAllColumns,
			promotionPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update promotions, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"promotions\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, promotionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(promotionType, promotionMapping, append(wl, promotionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "model: unable to update promotions row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for promotions")
	}

	if !cached {
		promotionUpdateCacheMut.Lock()
		promotionUpdateCache[key] = cache
		promotionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q promotionQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for promotions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for promotions")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PromotionSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), promotionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"promotions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, promotionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in promotion slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all promotion")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Promotion) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no promotions provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(promotionColumnsWithDefault, o)

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

	promotionUpsertCacheMut.RLock()
	cache, cached := promotionUpsertCache[key]
	promotionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			promotionAllColumns,
			promotionColumnsWithDefault,
			promotionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			promotionAllColumns,
			promotionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert promotions, could not build update column list")
		}

		ret := strmangle.SetComplement(promotionAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(promotionPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert promotions, could not build conflict column list")
			}

			conflict = make([]string, len(promotionPrimaryKeyColumns))
			copy(conflict, promotionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"promotions\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(promotionType, promotionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(promotionType, promotionMapping, ret)
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
		return errors.Wrap(err, "model: unable to upsert promotions")
	}

	if !cached {
		promotionUpsertCacheMut.Lock()
		promotionUpsertCache[key] = cache
		promotionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Promotion record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Promotion) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no Promotion provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), promotionPrimaryKeyMapping)
	sql := "DELETE FROM \"promotions\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from promotions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for promotions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q promotionQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no promotionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from promotions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for promotions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PromotionSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), promotionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"promotions\" WHERE " +
		strmangle.WhereInClause(string(dialect.LQ), string(dialect.RQ), 1, promotionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from promotion slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for promotions")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Promotion) Reload(exec boil.Executor) error {
	ret, err := FindPromotion(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PromotionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PromotionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), promotionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"promotions\".* FROM \"promotions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, promotionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in PromotionSlice")
	}

	*o = slice

	return nil
}

// PromotionExists checks if the Promotion row exists.
func PromotionExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"promotions\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if promotions exists")
	}

	return exists, nil
}

// Exists checks if the Promotion row exists.
func (o *Promotion) Exists(exec boil.Executor) (bool, error) {
	return PromotionExists(exec, o.ID)
}
