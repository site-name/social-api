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

type Or squirrel.Or

var _ qm.QueryMod = Or{}

func (or Or) Apply(q *queries.Query) {
	clause, args, err := squirrel.Or(or).ToSql()
	if err != nil {
		slog.Error("CustomOr ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}

type And squirrel.And

var _ qm.QueryMod = And{}

func (and And) Apply(q *queries.Query) {
	clause, args, err := squirrel.And(and).ToSql()
	if err != nil {
		slog.Error("CustomAnd ToSql", slog.Err(err))
		return
	}
	queries.AppendWhere(q, clause, args...)
}
