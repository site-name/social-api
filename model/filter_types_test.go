package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func squirrelSelector() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func TestStringFilter(t *testing.T) {
	// b := false
	filter := &StringFilter{

		// Or: &StringOption{
		// 	Eq: "This_is_name",
		// 	In: []string{"one", "two", "three"},
		// },

		Or: &StringOption{
			ExtraExpr: []squirrel.Sqlizer{
				squirrel.Expr("P.Used < P.Haha"),
			},
			Like: "minh",
			In:   []string{"one", "two"},
		},

		// StringOption: &StringOption{
		// 	Eq: "minhSon",
		// 	In: []string{"First", "last"},
		// 	NULL: &b,
		// },
	}

	query, args, err := squirrelSelector().
		Select("*").
		From("Persons").
		Where(filter.ToSquirrel("Name")).
		Suffix("FOR UPDATE ORDER BY Id").
		ToSql()

	require.NoError(t, err)

	fmt.Println("query:", query)
	fmt.Println("args:", args)
}

func Test_TimeFilter(t *testing.T) {
	nul := true
	now := time.Now()
	nowPlusOne := now.Add(time.Hour * 24)

	filter := &TimeFilter{

		// And: &TimeOption{
		// 	Gt:                &now,
		// 	LtE:               &nowPlusOne,
		// 	CompareStartOfDay: true,
		// 	NULL:              &nul,
		// },

		Or: &TimeOption{
			Gt:                &now,
			LtE:               &nowPlusOne,
			CompareStartOfDay: true,
			NULL:              &nul,
		},

		// TimeOption: &TimeOption{
		// Gt: &now,
		// LtE:               &nowPlusOne,
		// CompareStartOfDay: false,
		// NULL:              &nul,
		// },
	}

	query, args, err := squirrelSelector().
		Select("*").
		From("Persons").
		Where(filter.ToSquirrel("StartDate")).
		ToSql()

	require.NoError(t, err)
	fmt.Println("query:", query)
	fmt.Println("args:", args)
}

func TestSomething(t *testing.T) {
	query, args, err := squirrelSelector().
		Select("*").
		From("Persons AS P").
		Where(squirrel.Or{
			squirrel.Eq{"P.CreateAt": nil},
			squirrel.Expr("P.CreateAt < P.UsedAt"),
		}).
		ToSql()

	require.NoError(t, err)

	fmt.Println(query)
	fmt.Println(args...)
}
