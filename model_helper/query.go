package model_helper

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/modules/slog"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CommonQueryOptions struct {
	Conditions []qm.QueryMod
}

func NewCommonQueryOptions(conditions ...qm.QueryMod) CommonQueryOptions {
	return CommonQueryOptions{Conditions: conditions}
}

type Or squirrel.Or

var _ qm.QueryMod = Or{}

func (or Or) Apply(q *queries.Query) {
	if or == nil {
		return
	}
	clause, args, err := squirrel.Or(or).ToSql()
	if err != nil {
		slog.Error("Custom Or ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}

type And squirrel.And

var _ qm.QueryMod = And{}

func (and And) Apply(q *queries.Query) {
	if and == nil {
		return
	}
	clause, args, err := squirrel.And(and).ToSql()
	if err != nil {
		slog.Error("Custom And ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
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

	return qm.Where(fmt.Sprintf("%s::jsonb @> ?", field), fmt.Sprintf(template, key, value))
}

// JsonbHasKey builds a query mod that checks if a jsonb field contains a key
func JsonbHasKey(field string, key string) qm.QueryMod {
	return qm.Where(fmt.Sprintf("%s::jsonb -> '%s' IS NOT NULL", field, key))
}

// JsonbHasNoKey builds a query mod that checks if a jsonb field does not contain a key
func JsonbHasNoKey(field string, key string) qm.QueryMod {
	return qm.Where(fmt.Sprintf("%s::jsonb -> '%s' IS NULL", field, key))
}
