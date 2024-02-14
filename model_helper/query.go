package model_helper

import (
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

var _ qm.QueryMod = (*Or)(nil)

func (or *Or) Apply(q *queries.Query) {
	if or == nil {
		return
	}
	clause, args, err := squirrel.Or(*or).ToSql()
	if err != nil {
		slog.Error("CustomOr ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}

type And squirrel.And

var _ qm.QueryMod = (*And)(nil)

func (and *And) Apply(q *queries.Query) {
	if and == nil {
		return
	}
	clause, args, err := squirrel.And(*and).ToSql()
	if err != nil {
		slog.Error("CustomAnd ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}
