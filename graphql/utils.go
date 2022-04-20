package graphql

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web/shared"
)

// common id strings for creating AppErrors
const (
	UserUnauthenticatedId    = "graphql.account.user_unauthenticated.app_error"
	PermissionDeniedId       = "app.account.permission_denied.app_error"
	GraphqlArgumentInvalidID = "util.invalid_graphql_arguments.app_error"
)

// CheckUserAuthenticated is an utility function that check if session contained inside context is authenticated:
//
// 1) extracts value embedded in given ctx
//
// 2) checks whether session inside that value is nil or concrete
//
// 3) checks whether UserId property of session is valid uuid
func CheckUserAuthenticated(where string, ctx context.Context) (*model.Session, *model.AppError) {
	session := ctx.Value(shared.APIContextKey).(*shared.Context).AppContext.Session()

	if session == nil || !model.IsValidId(session.UserId) {
		return nil, model.NewAppError(where, UserUnauthenticatedId, nil, "", http.StatusForbidden)
	}
	return session, nil
}

// GraphqlArgumentsParser validates against these rules:
//
// 1) Either First or Last must be provided, not both
//
// 2) First and Before can not go together
//
// 3) Last and After can not go together
type GraphqlArgumentsParser struct {
	First          *int
	Last           *int
	Before         *string
	After          *string
	OrderDirection gqlmodel.OrderDirection
}

func (g *GraphqlArgumentsParser) IsValid() *model.AppError {
	if (g.First == nil && g.Last == nil) || (g.First != nil && g.Last != nil) {
		return model.NewAppError("GraphqlArgumentsParser.IsValid", GraphqlArgumentInvalidID, map[string]interface{}{"Fields": "Last, First"}, "You must provide either First or Last, not both", http.StatusBadRequest)
	}
	if g.First != nil && g.Before != nil {
		return model.NewAppError("GraphqlArgumentsParser.IsValid", GraphqlArgumentInvalidID, map[string]interface{}{"Fields": "First, Before"}, "First and Before can not go together", http.StatusBadRequest)
	}
	if g.Last != nil && g.After != nil {
		return model.NewAppError("GraphqlArgumentsParser.IsValid", GraphqlArgumentInvalidID, map[string]interface{}{"Fields": "Last, After"}, "Last and After can not go together", http.StatusBadRequest)
	}

	return nil
}

// Decode decodes before or after from base64 format to initial form, then returns the result
func (g *GraphqlArgumentsParser) decode() (string, *model.AppError) {
	var value string

	if g.Before != nil {
		value = *g.Before
	} else if g.After != nil {
		value = *g.After
	}

	res, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", model.NewAppError("GraphqlArgumentsParser.Decode", "graphql.error_decoding_base64.app_error", map[string]interface{}{"Value": value}, err.Error(), http.StatusBadRequest)
	}

	return string(res), nil
}

// ConstructSqlExpr does:
//
// 1) check if arguments are provided properly
//
// 2) decodes given before or after cursor
//
// 3) construct a squirrel expression based on given key
//  Eg
//  ConstructSqlExpr("TableName.FieldName") => squirrel.Gt{"TableName.FieldName": ...}
//
// NOTE: Call IsValid() before calling this
func (g *GraphqlArgumentsParser) ConstructSqlExpr(key string) (squirrel.Sqlizer, *model.AppError) {

	cmp, err := g.decode()
	if err != nil {
		return nil, err
	}

	if g.After != nil {
		if g.OrderDirection == gqlmodel.OrderDirectionAsc {
			// 1 2 3 4 5 6 (ASC)
			//     | *     (AFTER)
			return squirrel.Gt{key: cmp}, nil
		}

		// 6 5 4 3 2 1 (DESC)
		//       | *   (AFTER)
		return squirrel.Lt{key: cmp}, nil
	}

	if g.OrderDirection == gqlmodel.OrderDirectionAsc {
		// 1 2 3 4 5 6 (ASC)
		//   * |       (BEFORE)
		return squirrel.Lt{key: cmp}, nil
	}

	// 6 5 4 3 2 1 (DESC)
	//     * |     (BEFORE)
	return squirrel.Gt{key: cmp}, nil
}

// Limit must be called after
//
// NOTE: Call IsValid() before calling this
func (g *GraphqlArgumentsParser) Limit() int {
	if g.First != nil {
		return *g.First
	} else if g.Last != nil {
		return *g.Last
	}

	return 0
}

// HasPreviousPage determines if there is previous page
//
// NOTE: Call IsValid() before calling this
func (g *GraphqlArgumentsParser) HasPreviousPage() bool {
	return (g.First != nil && g.After != nil) || (g.Last != nil && g.Before != nil)
}
