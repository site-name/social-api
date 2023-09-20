package api

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"cmp"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

//go:embed graphql/schemas
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
	entries, err := assets.ReadDir("graphql/schemas")
	if err != nil {
		return "", errors.Wrap(err, "failed to read schema dir")
	}

	var builder strings.Builder
	for _, entry := range entries {
		extname := path.Ext(entry.Name())

		if !entry.IsDir() && (extname == ".graphql" || extname == ".graphqls") {
			data, err := assets.ReadFile(filepath.Join("graphql/schemas", entry.Name()))
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

// operands have format like:
//
//	operands := []any{"Users.Id", "9032nfuy45", "Users.Username", "minhson", "Users.Data": nil}
//	dest := convertGraphqlOperandsToString(operands
//	dest == "Users.Id:9032nfuy45____Users.Username:minhson____Users.Data:nil"
func convertGraphqlOperandsToString(operands []any) string {
	var cursorStrings = make([]string, len(operands)/2)
	var j = 0

	for i := 0; i < len(operands); i += 2 {
		key := operands[i].(string)
		value := operands[i+1]

		if value == nil {
			cursorStrings[j] = key + ":nil"
			j++
			continue
		}

		kind, _ := model.GetModelFieldKind(key)
		if _, ok := model.NormalModelKindsToAccordingPointerKindsMap[kind]; ok {
			value = &value
		}

		switch kind {
		case model.Time, model.TimePtr:
			cursorStrings[j] = fmt.Sprintf("%s:%s", key, (value.(*time.Time)).Format(time.RFC3339Nano))
		case model.Decimal, model.DecimalPtr:
			cursorStrings[j] = fmt.Sprintf("%s:%s", key, (value.(*decimal.Decimal).String()))
		case model.String, model.StringPtr:
			cursorStrings[j] = key + ":" + *(value.(*string))
		case model.Bool, model.BoolPtr:
			cursorStrings[j] = key + ":" + strconv.FormatBool(*(value.(*bool)))
		case model.Int, model.IntPtr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*int)))
		case model.Int8, model.Int8Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*int8)))
		case model.Int16, model.Int16Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*int16)))
		case model.Int32, model.Int32Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*int32)))
		case model.Int64, model.Int64Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*int64)))
		case model.Uint, model.UintPtr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*uint)))
		case model.Uint8, model.Uint8Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*uint8)))
		case model.Uint16, model.Uint16Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*uint16)))
		case model.Uint32, model.Uint32Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*uint32)))
		case model.Uint64, model.Uint64Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*uint64)))
		case model.Float32, model.Float32Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*float32)))
		case model.Float64, model.Float64Ptr:
			cursorStrings[j] = key + ":" + fmt.Sprintf("%v", *(value.(*float64)))

		default:
		}
		j++
	}

	return strings.Join(cursorStrings, cursorPartsSeperator)
}

// NOTE: Don't pass map[any]any values
func compareOperands(a, b any, kind model.ModelFieldKind) int {
	if a == b {
		return 0
	}
	if a == nil && b != nil {
		return -1
	}
	if a != nil && b == nil {
		return 1
	}

	if _, ok := model.NormalModelKindsToAccordingPointerKindsMap[kind]; ok {
		a, b = &a, &b
	}
	// from now on, a and b are pointer values

	switch kind {
	case model.TimePtr, model.Time:
		return a.(*time.Time).Compare(*(b.(*time.Time)))

	case model.Decimal, model.DecimalPtr:
		return a.(*decimal.Decimal).Cmp(*(b.(*decimal.Decimal)))

	case model.String, model.StringPtr:
		return cmp.Compare(*(a.(*string)), *(b.(*string)))

	case model.Bool, model.BoolPtr:
		return cmp.Compare(
			*(*uint8)(unsafe.Pointer(a.(*bool))),
			*(*uint8)(unsafe.Pointer(b.(*bool))),
		)
	case model.Int, model.IntPtr:
		return cmp.Compare(*a.(*int), *b.(*int))
	case model.Int8, model.Int8Ptr:
		return cmp.Compare(*a.(*int8), *b.(*int8))
	case model.Int16, model.Int16Ptr:
		return cmp.Compare(*a.(*int16), *b.(*int16))
	case model.Int32, model.Int32Ptr:
		return cmp.Compare(*a.(*int32), *b.(*int32))
	case model.Int64, model.Int64Ptr:
		return cmp.Compare(*a.(*int64), *b.(*int64))
	case model.Uint, model.UintPtr:
		return cmp.Compare(*a.(*uint), *b.(*uint))
	case model.Uint8, model.Uint8Ptr:
		return cmp.Compare(*a.(*uint8), *b.(*uint8))
	case model.Uint16, model.Uint16Ptr:
		return cmp.Compare(*a.(*uint16), *b.(*uint16))
	case model.Uint32, model.Uint32Ptr:
		return cmp.Compare(*a.(*uint32), *b.(*uint32))
	case model.Uint64, model.Uint64Ptr:
		return cmp.Compare(*a.(*uint64), *b.(*uint64))
	case model.Float32, model.Float32Ptr:
		return cmp.Compare(*a.(*float32), *b.(*float32))
	case model.Float64, model.Float64Ptr:
		return cmp.Compare(*a.(*float64), *b.(*float64))

	default: // this code should never be reached
		return 0
	}
}

// operands must have format like:
//
//	[]any{"Products.Slug", "hello-this-is-slug", "Products.CreateAt": 1678089}
//
// compareGraphqlOperands compares according values of two given maps. It returns:
//
// 1 if map1 > map2
//
// 0 if map1 == map2
//
// -1 if map1 < map2
func compareGraphqlOperands(operand1, operand2 []any) int {
	for i := 0; i < len(operand1); i += 2 {
		var (
			key     = operand1[i].(string)
			value1  = operand1[i+1]
			value2  = operand2[i+1]
			kind, _ = model.GetModelFieldKind(key)
			result  = compareOperands(value1, value2, kind)
		)

		if result != 0 {
			return result
		}
	}

	return 0
}

const cursorPartsSeperator = "____"

// If both Before and After are nil (happends when query the first page), returns nil, nil.
//
// Otherwise, return an []any with format like:
//
//	[]any{"Products.Slug", "hello-slug", "Products.CreateAt", 16348956, ...}
//
// and nil error
func parseGraphqlCursor(params *GraphqlParams) ([]any, error) {
	// this case happends when user initially fetch first page
	if params.Before == nil && params.After == nil {
		return nil, nil
	}

	var cursor string
	switch {
	case params.Before != nil:
		data, err := base64.StdEncoding.DecodeString(*params.Before)
		if err != nil {
			return nil, errors.Wrap(err, "invalid cursor before provided")
		}
		cursor = string(data)
	default:
		data, err := base64.StdEncoding.DecodeString(*params.After)
		if err != nil {
			return nil, errors.Wrap(err, "invalid cursor after provided")
		}
		cursor = string(data)
	}

	splitCursor := strings.Split(cursor, cursorPartsSeperator)

	var res = make([]any, 0, len(splitCursor)*2)

	for _, cursorPart := range splitCursor {
		splitCursorsParts := strings.Split(cursorPart, ":")
		if len(splitCursorsParts) != 2 {
			return nil, errors.Errorf("expect cursor part to have 2 component, got %d", len(splitCursorsParts))
		}
		key, value := splitCursorsParts[0], splitCursorsParts[1]

		if value == "nil" {
			res = append(res, key, nil)
			continue
		}

		kind, found := model.GetModelFieldKind(key)
		if !found {
			return nil, errors.Errorf("invalid field key: %s", key)
		}

		switch kind {
		case model.Decimal, model.DecimalPtr:
			deci, err := decimal.NewFromString(value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse decimal value")
			}
			var value any = deci
			if kind == model.DecimalPtr {
				value = &deci
			}
			res = append(res, key, value)

		case model.Time, model.TimePtr:
			tim, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse time value")
			}
			var value any = tim
			if kind == model.TimePtr {
				value = &tim
			}
			res = append(res, key, value)

		case model.Bool, model.BoolPtr:
			boo, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse bool value")
			}
			var value any = boo
			if kind == model.BoolPtr {
				value = &boo
			}
			res = append(res, key, value)

		case model.Int, model.IntPtr:
			in, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse int value")
			}
			var value any = int(in)
			if kind == model.IntPtr {
				value = &value
			}
			res = append(res, key, value)

		case model.Int8, model.Int8Ptr:
			in, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse int8 value")
			}
			var value any = int8(in)
			if kind == model.Int8Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Int16, model.Int16Ptr:
			in, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse int16 value")
			}
			var value any = int16(in)
			if kind == model.Int16Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Int32, model.Int32Ptr:
			in, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse int32 value")
			}
			var value any = int32(in)
			if kind == model.Int32Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Int64, model.Int64Ptr:
			in, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse int64 value")
			}
			var value any = in
			if kind == model.Int64Ptr {
				value = &in
			}
			res = append(res, key, value)

		case model.Uint, model.UintPtr:
			uin, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse uint value")
			}
			var value any = uint(uin)
			if kind == model.UintPtr {
				value = &value
			}
			res = append(res, key, value)

		case model.Uint8, model.Uint8Ptr:
			uin, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse uint8 value")
			}
			var value any = uint8(uin)
			if kind == model.Uint8Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Uint16, model.Uint16Ptr:
			uin, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse uint16 value")
			}
			var value any = uint16(uin)
			if kind == model.Uint16Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Uint32, model.Uint32Ptr:
			uin, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse uint32 value")
			}
			var value any = uint32(uin)
			if kind == model.Uint32Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Uint64, model.Uint64Ptr:
			uin, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse uint64 value")
			}
			var value any = uin
			if kind == model.Uint64Ptr {
				value = &uin
			}
			res = append(res, key, value)

		case model.Float32, model.Float32Ptr:
			float, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse float32 value")
			}
			var value any = float32(float)
			if kind == model.Float32Ptr {
				value = &value
			}
			res = append(res, key, value)

		case model.Float64, model.Float64Ptr:
			float, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse float64 value")
			}
			var value any = float
			if kind == model.Float64Ptr {
				value = &float
			}
			res = append(res, key, value)

		case model.String, model.StringPtr:
			var str any = value
			if kind == model.StringPtr {
				str = &value
			}
			res = append(res, key, str)

		default:
			// NOTE: there is still Map, Struct and Slice types.
			// But we don't sort records using map nor slice nor struct types.
			// So I decide to ignore it here
			return nil, errors.Errorf("unsupported type")
		}
	}

	return res, nil
}

// GraphqlParams is provided in some resolver methods
type GraphqlParams struct {
	// after base64 decoding, Before must have format like:
	//  "Products.Slug:hello-world____Products.DeleteAt:nil"
	Before *string `json:"before"`
	// after base64 decoding, After must have format like:
	//  "Products.Slug:hello-world____Products.DeleteAt:nil"
	After *string `json:"after"`
	First *int32  `json:"first"`
	Last  *int32  `json:"last"`

	validated   bool
	memoizedErr *model.AppError
}

// if First != nil, returns "ASC". Otherwise return "DESC"
func (g *GraphqlParams) orderDirection() string {
	if g.First != nil {
		return "ASC"
	}
	return "DESC"
}

// If First or Last is provided, return (First || Last) + 1
//
// The trick is to help determine if there is next page available
func (g *GraphqlParams) queryLimit() uint64 {
	switch {
	case g.First != nil:
		return uint64(*g.First) + 1
	case g.Last != nil:
		return uint64(*g.Last) + 1
	default:
		return 0
	}
}

func (g *GraphqlParams) checkNextPageAndPreviousPage(lengthOfRecordSliceFounded int) (hasNextPage, hasPreviousPage bool) {
	queryLimit := g.queryLimit()
	hasNextPage = queryLimit != 0 && lengthOfRecordSliceFounded == int(queryLimit)
	hasPreviousPage = queryLimit != 0 && (g.Before != nil || g.After != nil)
	return
}

func (g *GraphqlParams) validate(where string) *model.AppError {
	if g.validated {
		return g.memoizedErr
	}
	g.validated = true

	switch {
	case (g.First != nil && *g.First < 0) || (g.Last != nil && *g.Last < 0):
		g.memoizedErr = model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "First / Last"}, "First and Last cannot be negative", http.StatusBadRequest)
	case (g.First != nil && g.Last != nil) || (g.First == nil && g.Last == nil):
		g.memoizedErr = model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "First / Last"}, "provide either First or Last, not both", http.StatusBadRequest)
	case g.First != nil && g.Before != nil:
		g.memoizedErr = model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "First / Before"}, "First and Before can not go together", http.StatusBadRequest)
	case g.Last != nil && g.After != nil:
		g.memoizedErr = model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "Last / After"}, "Last and After can not go together", http.StatusBadRequest)
	case g.Before != nil && g.After != nil:
		g.memoizedErr = model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "Before / After"}, "Before and After can not go together", http.StatusBadRequest)
	default:
		g.memoizedErr = nil
	}

	return g.memoizedErr
}

func (s *graphqlPaginator[RawT, DestT]) Len() int {
	return len(s.data)
}

func (s *graphqlPaginator[RawT, DestT]) Less(i, j int) bool {
	return compareGraphqlOperands(s.keyFunc(s.data[i]), s.keyFunc(s.data[j])) == -1
}

func (s *graphqlPaginator[RawT, DestT]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

const PaginationError = "api.graphql.pagination_params_invalid.app_error"

// graphqlPaginator is used to paginate the whole set of model records in graphql's way.
type graphqlPaginator[RawT any, DestT any] struct {
	// E.g
	//  []*model.Product{...}
	data []RawT
	// extract values from a model record.
	// E.g
	//  func(c *model.Category) []any {
	//     return []any{
	//         "Categories.Slug",
	//         "category-slug-value",
	//         "Categories.Name",
	//         "This is category name",
	//     }
	//  }
	keyFunc func(RawT) []any
	// E.g
	//  func(c *model.Category) *GraphqlCategory {...}
	modelTypeToGraphqlTypeFunc func(RawT) DestT // convert raw system model types to their according graphql type
	GraphqlParams
}

// newGraphqlPaginator returns *graphqlPaginator formed using given arguments.
// Use this instead of manually construct &graphqlPaginator{} to prevent missing some fields.
func newGraphqlPaginator[RawT any, DestT any](
	data []RawT,
	keyFunc func(RawT) []any,
	modelTypeToGraphqlTypeFunc func(RawT) DestT,
	params GraphqlParams) *graphqlPaginator[RawT, DestT] {
	return &graphqlPaginator[RawT, DestT]{data, keyFunc, modelTypeToGraphqlTypeFunc, params}
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

func constructCountableConnection[R any, D any](
	data []R,
	totalCount int64,
	hasNextPage, hasPreviousPage bool,
	keyFunc func(R) []any, // E.g func(p *model.Product) []any{"Products.CreateAt", 1674545, "Products.Name", "hello world"}
	modelTypeToGraphqlTypeFunc func(R) D,
) *CountableConnection[D] {
	res := &CountableConnection[D]{
		TotalCount: (*int32)(unsafe.Pointer(&totalCount)),
		Edges: lo.Map(data, func(item R, _ int) *CountableConnectionEdge[D] {
			stringRawCursor := convertGraphqlOperandsToString(keyFunc(item))

			return &CountableConnectionEdge[D]{
				Cursor: base64.StdEncoding.EncodeToString([]byte(stringRawCursor)),
				Node:   modelTypeToGraphqlTypeFunc(item),
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

	return res
}

func (g *graphqlPaginator[RawT, DestT]) parse(where string) (*CountableConnection[DestT], *model.AppError) {
	appErr := g.validate(where)
	if appErr != nil {
		return nil, appErr
	}

	orderASC := g.orderDirection() == "ASC" // order ascending or not

	if orderASC {
		sort.Sort(g)
	} else {
		sort.Sort(sort.Reverse(g))
	}

	operand, err := parseGraphqlCursor(&g.GraphqlParams)
	if err != nil {
		return nil, model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "Before / After"}, err.Error(), http.StatusInternalServerError)
	}

	var (
		resultData                   []RawT
		hasNextPage, hasPreviousPage bool
		index                        int
		limit                        = g.First
		totalCount                   = g.Len()
	)

	// return immediately when no data passed
	if totalCount == 0 {
		goto result
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
		goto result
	}

	// case operand provided:
	index = sort.Search(totalCount, func(i int) bool {
		value := g.keyFunc(g.data[i])
		cmp := compareGraphqlOperands(value, operand)

		// order ASC && >= || order DESC && <=
		return (orderASC && cmp >= 0) || cmp <= 0
	})

	// if not found, sort.Search returns exactly first int argument passed. We need to check it here
	if index >= totalCount {
		return nil, model.NewAppError(where, PaginationError, map[string]interface{}{"Fields": "before / after"}, "invalid before or after provided", http.StatusBadRequest)
	}

	hasPreviousPage = true
	resultData = g.data[index+1:]

	if *limit < int32(len(resultData)) {
		resultData = resultData[:*limit]
		hasNextPage = true
	}

result:
	return constructCountableConnection(resultData, int64(totalCount), hasNextPage, hasPreviousPage, g.keyFunc, g.modelTypeToGraphqlTypeFunc), nil
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

func prepareFilterExpression(fieldName string, index int, cursors []any, sortingFields []string, sortAscending bool) (squirrel.Or, squirrel.And) {
	var fieldExpression = squirrel.And{}
	var extraExpression = squirrel.Or{}

	for idx, cursorValue := range cursors[:index] {
		fieldExpression = append(fieldExpression, squirrel.Expr(sortingFields[idx]+" = ?", cursorValue))
	}

	if sortAscending {
		extraExpression = append(
			extraExpression,
			squirrel.Expr(fieldName+" > ?", cursors[index]),
			squirrel.Expr(fieldName+" IS NULL"),
		)
	} else if cursors[index] != nil {
		var expr squirrel.Sqlizer = squirrel.Expr(fieldName+" > ?", cursors[index])
		if !sortAscending {
			expr = squirrel.Expr(fieldName+" < ?", cursors[index])
		}
		fieldExpression = append(fieldExpression, expr)
	} else {
		fieldExpression = append(fieldExpression, squirrel.Expr(fieldName+" IS NOT NULL"))
	}

	return extraExpression, fieldExpression
}

func (g *GraphqlParams) Parse(where string) (*model.GraphqlPaginationValues, *model.AppError) {
	appErr := g.validate(where)
	if appErr != nil {
		return nil, appErr
	}

	operand, err := parseGraphqlCursor(g)
	if err != nil {
		return nil, model.NewAppError(where, model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "GraphqlParams"}, err.Error(), http.StatusBadRequest)
	}

	var (
		cursors          = make([]any, 0, len(operand)/2)
		sortingFields    = make(util.AnyArray[string], 0, len(operand)/2)
		orderDirection   = g.orderDirection()
		sortingAscending = orderDirection == "ASC"
		conditions       = squirrel.Or{}
	)

	if len(operand) > 0 {
		for i := 0; i < len(operand); i += 2 {
			cursors = append(cursors, operand[i+1])
			sortingFields = append(sortingFields, operand[i].(string))
		}

		for idx, fieldName := range sortingFields {
			if cursors[idx] == nil && sortingAscending {
				continue
			}

			extraExpr, fieldExpr := prepareFilterExpression(fieldName, idx, cursors, sortingFields, sortingAscending)
			conditions = append(conditions, squirrel.And{extraExpr, fieldExpr})
		}
	}

	res := &model.GraphqlPaginationValues{
		Limit:   g.queryLimit(),
		OrderBy: sortingFields.Map(func(_ int, item string) string { return item + " " + orderDirection }).Join(","),
	}
	if len(conditions) > 0 {
		res.Condition = conditions
	}

	return res, nil
}

// getMax returns the largest item from given items
func getMax[T cmp.Ordered](items ...T) T {
	var max T
	if len(items) == 0 {
		return max
	}
	if len(items) == 1 {
		return items[0]
	}
	max = items[0]
	for i := 1; i < len(items); i++ {
		if itemAtI := items[i]; itemAtI > max {
			max = itemAtI
		}
	}

	return max
}
