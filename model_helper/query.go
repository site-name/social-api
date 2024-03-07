package model_helper

import (
	"fmt"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/modules/slog"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

type CommonQueryOptions struct {
	Conditions []qm.QueryMod
}

func NewCommonQueryOptions(conditions ...qm.QueryMod) CommonQueryOptions {
	return CommonQueryOptions{Conditions: conditions}
}

var _ qm.QueryMod = (*SelectBuilder)(nil)

type SelectBuilder squirrel.SelectBuilder

func (s SelectBuilder) Apply(q *queries.Query) {
	query, args, err := squirrel.SelectBuilder(s).ToSql()
	if err != nil {
		slog.Error("Custom SelectBuilder ToSql", slog.Err(err))
		return
	}

	queries.AppendWhere(q, query, args...)
}

type Or squirrel.Or

var _ qm.QueryMod = (*Or)(nil)

func (or Or) Apply(q *queries.Query) {
	clause, args, err := squirrel.Or(or).ToSql()
	if err != nil {
		slog.Error("Custom Or ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}

func (or Or) ToSql() (string, []any, error) {
	return squirrel.Or(or).ToSql()
}

type And squirrel.And

var _ qm.QueryMod = (*And)(nil)

func (and And) Apply(q *queries.Query) {
	clause, args, err := squirrel.And(and).ToSql()
	if err != nil {
		slog.Error("Custom And ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}

func (and And) ToSql() (string, []any, error) {
	return squirrel.And(and).ToSql()
}

// JsonbContains buils a query mod that checks if a jsonb field contains a key-value pair
func JsonbContains(field string, key string, value any) qm.QueryMod {
	var template string

	switch value.(type) {
	case string:
		template = "{%q:%q}"
	default:
		template = "{%q:%v}"
	}

	return qmhelper.WhereQueryMod{
		Clause: fmt.Sprintf("%s::jsonb @> ?", field),
		Args:   []any{fmt.Sprintf(template, key, value)},
	}
}

// JsonbHasKey builds a query mod that checks if a jsonb field contains a key
func JsonbHasKey(field string, key string) qm.QueryMod {
	return qmhelper.WhereQueryMod{
		Clause: fmt.Sprintf("%s::jsonb -> '%s' IS NOT NULL", field, key),
	}
}

// JsonbHasNoKey builds a query mod that checks if a jsonb field does not contain a key
func JsonbHasNoKey(field string, key string) qm.QueryMod {
	return qmhelper.WhereQueryMod{
		Clause: fmt.Sprintf("%s::jsonb -> '%s' IS NULL", field, key),
	}
}

// AnnotationAggregator is a query modifier that adds annotations to the query.
// E.g:
//
//	AnnotationAggregator{
//		"another": `1 + 2`,
//	}
//
// When applied to a query, it will produce:
//
//	`SELECT JSON_BUILD_OBJECT('another', 1 + 2) AS "annotations"``
type AnnotationAggregator map[string]any

func (a AnnotationAggregator) Apply(q *queries.Query) {
	if len(a) == 0 {
		return
	}

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	counter := 0
	for key, value := range a {
		if counter > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("'%s', %v", key, value))
		counter++
	}

	queries.AppendSelect(q, fmt.Sprintf(`JSON_BUILD_OBJECT(%s) AS "annotations"`, buf.String()))
}
