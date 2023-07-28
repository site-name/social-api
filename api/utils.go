package api

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

//go:embed schemas
var assets embed.FS

// stringsContainSqlExpr is used to validate strings values contain sql statements or not
var stringsContainSqlExpr = regexp.MustCompile(`(?i)\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER)\b`)

func MakeUnauthorizedError(where string) *model.AppError {
	return model.NewAppError(where, "api.unauthorized.app_error", nil, "you are not allowed to perform this action", http.StatusUnauthorized)
}

// Unique type to hold our context.
type CTXKey int

const WebCtx CTXKey = iota

// constructSchema constructs schema from *.graphql(s) files
func constructSchema() (string, error) {
	entries, err := assets.ReadDir("schemas")
	if err != nil {
		return "", errors.Wrap(err, "failed to read schema dir")
	}

	var builder strings.Builder
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), ".graphql") {
			data, err := assets.ReadFile(filepath.Join("schemas", entry.Name()))
			if err != nil {
				return "", errors.Wrapf(err, "failed to read schema file: %s", filepath.Join("schemas", entry.Name()))
			}

			builder.Write(data)
			builder.WriteByte('\n')
		}
	}

	return builder.String(), nil
}

// GetContextValue extracts according value of given key in given `ctx` and returns the value.
func GetContextValue[T any](ctx context.Context, key CTXKey) T {
	return ctx.Value(key).(T)
}

func MetadataToSlice[T any](m map[string]T) []*MetadataItem {
	return lo.MapToSlice(m, func(k string, v T) *MetadataItem {
		var strValue string
		if impl, ok := any(v).(fmt.Stringer); ok {
			strValue = impl.String()
		} else {
			strValue = fmt.Sprintf("%v", v)
		}
		return &MetadataItem{
			Key:   k,
			Value: strValue,
		}
	})
}

func SystemMoneyToGraphqlMoney(money *goprices.Money) *Money {
	if money == nil {
		return nil
	}
	return &Money{
		Currency: money.Currency,
		Amount:   money.Amount.InexactFloat64(),
	}
}

func SystemTaxedMoneyToGraphqlTaxedMoney(money *goprices.TaxedMoney) *TaxedMoney {
	if money == nil {
		return nil
	}
	return &TaxedMoney{
		Currency: money.Currency,
		Gross:    SystemMoneyToGraphqlMoney(money.Gross),
		Net:      SystemMoneyToGraphqlMoney(money.Net),
		Tax:      SystemMoneyToGraphqlMoney(money.Tax()),
	}
}

func SystemTaxedMoneyRangeToGraphqlTaxedMoneyRange(m *goprices.TaxedMoneyRange) *TaxedMoneyRange {
	if m == nil {
		return nil
	}
	return &TaxedMoneyRange{
		Start: SystemTaxedMoneyToGraphqlTaxedMoney(m.Start),
		Stop:  SystemTaxedMoneyToGraphqlTaxedMoney(m.Stop),
	}
}

func SystemMoneyRangeToGraphqlMoneyRange(money *goprices.MoneyRange) *MoneyRange {
	if money == nil {
		return nil
	}
	return &MoneyRange{
		Start: SystemMoneyToGraphqlMoney(money.Start),
		Stop:  SystemMoneyToGraphqlMoney(money.Stop),
	}
}

// systemRecordsToGraphql converts slice of system models to graphql representations of them
//
// E.g:
//
//	systemRecordsToGraphql([]*model.Product, func(p *model.Product) *Product) => []*Product
func systemRecordsToGraphql[S any, D any](slice []S, iteratee func(S) D) []D {
	res := make([]D, len(slice))

	for idx, item := range slice {
		res[idx] = iteratee(item)
	}

	return res
}

func convertGraphqlOperandToString[C graphqlCursorType](operand C) string {
	switch t := any(operand).(type) {
	case time.Time:
		return t.Format(time.RFC3339)
	case decimal.Decimal:
		return t.String()

	default:
		return fmt.Sprintf("%v", t)
	}
}

// decodeBase64String uses base64 to decode given string to its original value.
//
// NOTE: given value MUST be base64 encoded, otherwise an empty string "" is returned.
// func decodeBase64String(value string) string {
// 	decoded, _ := base64.StdEncoding.DecodeString(value)
// 	return string(decoded)
// }

// // NOTE: You must be sure that given values are base64 encoded.
// func decodeBase64Strings(values ...string) []string {
// 	return lo.Map(values, func(v string, _ int) string { return decodeBase64String(v) })
// }

// parseGraphqlOperand can possibly returns (nil, nil)
func parseGraphqlOperand[C graphqlCursorType](params GraphqlParams) (*C, *model.AppError) {
	// in case users query resuts for the first time
	if params.Before == nil && params.After == nil {
		return nil, nil
	}

	// convert base64 cursor to string:
	var byteCursorValue []byte
	var err error
	if params.Before != nil {
		byteCursorValue, err = base64.StdEncoding.DecodeString(*params.Before)
	} else if params.After != nil {
		byteCursorValue, err = base64.StdEncoding.DecodeString(*params.After)
	}

	if err != nil {
		return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not valid base64 encoded string", http.StatusBadRequest)
	}
	var cursorValue = string(byteCursorValue)

	var res C
	switch any(res).(type) {
	case string:
		return (*C)(unsafe.Pointer(&cursorValue)), nil

	case float64:
		float, err := strconv.ParseFloat(cursorValue, 64)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not float64 based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&float)), nil

	case int:
		i32, err := strconv.ParseInt(cursorValue, 10, 32)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not int based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&i32)), nil

	case int64:
		i64, err := strconv.ParseInt(cursorValue, 10, 64)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not int64 based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&i64)), nil

	case uint64:
		ui64, err := strconv.ParseUint(cursorValue, 10, 64)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not uint64 based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&ui64)), nil

	case time.Time:
		tim, err := time.Parse(time.RFC3339, cursorValue)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not RFC3339 time based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&tim)), nil

	default:
		deci, err := decimal.NewFromString(cursorValue)
		if err != nil {
			return nil, model.NewAppError("parseGraphqlOperand", PaginationError, map[string]interface{}{"Fields": "before / after"}, "before or after is not decimal based string", http.StatusBadRequest)
		}
		return (*C)(unsafe.Pointer(&deci)), nil
	}
}

// If the type is time.Time, we always parse it in RFC3339 format
type graphqlCursorType interface {
	string | float64 | int | int64 | uint64 | time.Time | decimal.Decimal
}

// It returns -1 if a < b, 0 if a == b and +1 if a > b
func comparePrimitives[T util.Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// compareGraphqlOperands compares a and b and returns int.
//
// It returns -1 if a < b, 0 if a == b and +1 if a > b
func compareGraphqlOperands[K graphqlCursorType](a, b K) int {
	anyB := any(b)
	anyA := any(a)

	switch t := anyA.(type) {
	case time.Time:
		return t.Compare(anyB.(time.Time))
	case decimal.Decimal:
		return t.Cmp(anyB.(decimal.Decimal))
	case string:
		return comparePrimitives(t, anyB.(string))
	case int:
		return comparePrimitives(t, anyB.(int))
	case float64:
		return comparePrimitives(t, anyB.(float64))
	case int64:
		return comparePrimitives(t, anyB.(int64))

	default:
		return comparePrimitives(t.(uint64), anyB.(uint64))
	}
}

// GraphqlParams is provided in some resolver methods
type GraphqlParams struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
	First  *int32  `json:"first"`
	Last   *int32  `json:"last"`

	validated bool
}

// ParseGraphqlParams parse given graphql params to produce *model.PaginationValues for used in database query pagination
func parseGraphqlParams[C graphqlCursorType](params *GraphqlParams, api, orderKey string) (*model.PaginationValues, *model.AppError) {
	if !params.validated {
		appErr := params.Validate(api)
		if appErr != nil {
			return nil, appErr
		}
	}

	res := &model.PaginationValues{}

	operand, appErr := parseGraphqlOperand[C](*params)
	if appErr != nil {
		return nil, appErr
	}

	switch {
	case params.First != nil:
		res.OrderBy = orderKey + " ASC"
		res.Limit = *(*uint64)(unsafe.Pointer(params.First))
		if operand != nil {
			res.Condition = squirrel.Gt{orderKey: *operand}
		}

	default:
		res.OrderBy = orderKey + " DESC"
		res.Limit = *(*uint64)(unsafe.Pointer(params.Last))
		if operand != nil {
			res.Condition = squirrel.Lt{orderKey: *operand}
		}
	}

	return res, nil
}

func (g *GraphqlParams) Validate(apiName string) *model.AppError {
	g.validated = true
	if (g.First != nil && *g.First < 0) || (g.Last != nil && *g.Last < 0) {
		return model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "First / Last"}, "first and last cannot be negative", http.StatusBadRequest)
	}
	if (g.First != nil && g.Last != nil) || (g.First == nil && g.Last == nil) {
		return model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "First / Last"}, "provide either First or Last, not both", http.StatusBadRequest)
	}
	if g.First != nil && g.Before != nil {
		return model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "First / Before"}, "First and Before can't go together", http.StatusBadRequest)
	}
	if g.Last != nil && g.After != nil {
		return model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Last / After"}, "Last and After can't go together", http.StatusBadRequest)
	}
	if g.Before != nil && g.After != nil {
		return model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Before / After"}, "Before and After can'g go together", http.StatusBadRequest)
	}

	return nil
}

func (s *graphqlPaginator[R, C, D]) Len() int {
	return len(s.data)
}

func (s *graphqlPaginator[R, C, D]) Less(i, j int) bool {
	return compareGraphqlOperands(s.keyFunc(s.data[i]), s.keyFunc(s.data[j])) == -1
}

func (s *graphqlPaginator[R, C, D]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

const PaginationError = "api.graphql.pagination_params_invalid.app_error"

// graphqlPaginator implements sort.Interface
type graphqlPaginator[RawT any, CurT graphqlCursorType, DestT any] struct {
	data                  []RawT           // E.g []*model.Product
	keyFunc               func(RawT) CurT  // extract value from system model types
	rawTypeToDestTypeFunc func(RawT) DestT // convert raw system model types to their according graphql type
	GraphqlParams
}

// newGraphqlPaginator returns *graphqlPaginator formed using given arguments.
// Use this instead of manually construct &graphqlPaginator{} to prevent missing some fields.
func newGraphqlPaginator[RawT any, CurT graphqlCursorType, DestT any](
	data []RawT,
	keyFunc func(RawT) CurT,
	rawTypeToDestTypeFunc func(RawT) DestT,
	params GraphqlParams) *graphqlPaginator[RawT, CurT, DestT] {
	return &graphqlPaginator[RawT, CurT, DestT]{data, keyFunc, rawTypeToDestTypeFunc, params}
}

// CountableConnection shares similar memory layout as all graphql api Connections.
type CountableConnection[D any] struct {
	PageInfo   *PageInfo
	Edges      []*CountableConnectionEdge[D]
	TotalCount *int32
}

type CountableConnectionEdge[D any] struct {
	Node   D
	Cursor string
}

func (g *graphqlPaginator[R, C, D]) parse(apiName string) (*CountableConnection[D], *model.AppError) {
	if !g.validated {
		appErr := g.Validate(apiName)
		if appErr != nil {
			return nil, appErr
		}
	}

	orderASC := g.First != nil // order ascending or not

	if orderASC {
		sort.Sort(g)
	} else {
		sort.Sort(sort.Reverse(g))
	}

	operand, err := parseGraphqlOperand[C](g.GraphqlParams)
	if err != nil {
		return nil, model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Before / After"}, err.Error(), http.StatusInternalServerError)
	}

	var (
		resultData                   []R
		hasNextPage, hasPreviousPage bool
		index                        int
		limit                        = g.First
		totalCount                   = g.Len()
	)

	// return immediately when no data passed
	if totalCount == 0 {
		goto returnLabel
	}

	if limit == nil {
		limit = g.Last
	}

	if operand == nil {
		if *limit < int32(totalCount) { // prevent slicing out of range
			resultData = g.data[:*limit]
			hasNextPage = true
		} else {
			resultData = g.data
		}
		goto returnLabel
	}

	// case operand provided:
	index = sort.Search(totalCount, func(i int) bool {
		value := g.keyFunc(g.data[i])
		cmp := compareGraphqlOperands(value, *operand)

		// order ASC && >= || order DESC && <=
		return (orderASC && cmp >= 0) || cmp <= 0
	})

	// if not found, sort.Search returns exactly first int argument passed. We need to check it here
	if index >= totalCount {
		return nil, model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "before / after"}, "invalid before or after provided", http.StatusBadRequest)
	}

	hasPreviousPage = true
	resultData = g.data[index+1:]

	if *limit < int32(len(resultData)) {
		resultData = resultData[:*limit]
		hasNextPage = true
	}

returnLabel:
	res := &CountableConnection[D]{
		TotalCount: (*int32)(unsafe.Pointer(&totalCount)),
		Edges: lo.Map(resultData, func(item R, _ int) *CountableConnectionEdge[D] {
			stringRawCursor := convertGraphqlOperandToString(g.keyFunc(item))

			return &CountableConnectionEdge[D]{
				Cursor: base64.StdEncoding.EncodeToString([]byte(stringRawCursor)),
				Node:   g.rawTypeToDestTypeFunc(item),
			}
		}),
	}
	res.PageInfo = &PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
	if len(res.Edges) > 0 {
		res.PageInfo.StartCursor = &res.Edges[0].Cursor
		res.PageInfo.EndCursor = &res.Edges[len(res.Edges)-1].Cursor
	}

	return res, nil
}

func reportingPeriodToDate(period ReportingPeriod) time.Time {
	now := time.Now()

	switch period {
	case ReportingPeriodToday:
		return util.StartOfDay(now)
	default:
		return util.StartOfMonth(now)
	}
}
