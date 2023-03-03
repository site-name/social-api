package api

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

//go:embed schemas
var assets embed.FS

// ErrorUnauthorized
const ErrorUnauthorized = "api.unauthorized.app_error"
const ErrorChannelIDQueryParamMissing = "api.channel_id.missing.app_error"

// Unique type to hold our context.
type CTXKey int

const (
	WebCtx CTXKey = iota
	ChannelIdCtx
)

// constructSchema constructs schema from *.graphql files
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
func GetContextValue[T any](ctx context.Context, key CTXKey) (T, error) {
	value := ctx.Value(key)
	if value == nil {
		var res T
		return res, fmt.Errorf("context doesn't store given key")
	}

	cast, ok := value.(T)
	if !ok {
		var res T
		return res, fmt.Errorf("found value has unexpected type: %T", value)
	}

	return cast, nil
}

func MetadataToSlice[T any](m map[string]T) []*MetadataItem {
	return lo.MapToSlice(m, func(k string, v T) *MetadataItem {
		return &MetadataItem{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
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

	tax, _ := money.Tax()
	return &TaxedMoney{
		Currency: money.Currency,
		Gross:    SystemMoneyToGraphqlMoney(money.Gross),
		Net:      SystemMoneyToGraphqlMoney(money.Net),
		Tax:      SystemMoneyToGraphqlMoney(tax),
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

func SystemLanguageToGraphqlLanguageCodeEnum(code string) LanguageCodeEnum {
	if len(code) == 0 {
		return LanguageCodeEnumEn
	}

	upperCaseCode := strings.Map(func(r rune) rune {
		if r == '-' {
			return '_'
		}

		return unicode.ToUpper(r)
	}, code)

	res := LanguageCodeEnum(upperCaseCode)

	if !res.IsValid() {
		return LanguageCodeEnumEn
	}

	return res
}

// DataloaderResultMap converts slice of system models to graphql representations of them
//
// E.g:
//
//	DataloaderResultMap([]*model.Product, func(p *model.Product) *Product) => []*Product
func DataloaderResultMap[S any, D any](slice []S, iteratee func(S) D) []D {
	res := make([]D, len(slice))

	for idx, item := range slice {
		res[idx] = iteratee(item)
	}

	return res
}

// parseOperand can possibly returns (nil, nil)
func (p *graphqlPaginator[_, K]) parseOperand() (*K, error) {
	var res K

	// in case users query resuts for the first time
	if p.Before == nil && p.After == nil {
		return nil, nil
	}

	// convert base64 cursor to string:
	var strCursorValue []byte
	var err error
	if p.Before != nil {
		strCursorValue, err = base64.StdEncoding.DecodeString(*p.Before)
	} else if p.After != nil {
		strCursorValue, err = base64.StdEncoding.DecodeString(*p.After)
	}
	if err != nil {
		return nil, err
	}
	var cursorValue = string(strCursorValue)

	switch any(res).(type) {
	case string:
		return (*K)(unsafe.Pointer(&cursorValue)), nil

	case float64:
		float, err := strconv.ParseFloat(cursorValue, 64)
		if err != nil {
			return nil, err
		}
		return (*K)(unsafe.Pointer(&float)), nil

	case int:
		i32, err := strconv.ParseInt(cursorValue, 10, 32)
		if err != nil {
			return nil, err
		}
		return (*K)(unsafe.Pointer(&i32)), nil

	case int64:
		i64, err := strconv.ParseInt(cursorValue, 10, 64)
		if err != nil {
			return nil, err
		}
		return (*K)(unsafe.Pointer(&i64)), nil

	case uint64:
		ui64, err := strconv.ParseUint(cursorValue, 10, 64)
		if err != nil {
			return nil, nil
		}
		return (*K)(unsafe.Pointer(&ui64)), nil

	case time.Time:
		tim, err := time.Parse(time.RFC3339, cursorValue)
		if err != nil {
			return nil, err
		}
		return (*K)(unsafe.Pointer(&tim)), nil

	case decimal.Decimal:
		deci, err := decimal.NewFromString(cursorValue)
		if err != nil {
			return nil, err
		}
		return (*K)(unsafe.Pointer(&deci)), nil

	default:
		return nil, fmt.Errorf("unknwon dest type: %T", res)
	}
}

// If the type is time.Time, we always parse it in RFC3339 format
type graphqlCursorType interface {
	string | float64 | int | int64 | uint64 | time.Time | decimal.Decimal
}

type CompareOrder int8

const (
	Lesser CompareOrder = iota
	Equal
	Greater
)

func comparePrimitives[T util.Ordered](a, b T) CompareOrder {
	if a < b {
		return Lesser
	} else if a > b {
		return Greater
	}
	return Equal
}

// compare compares a and b and returns CompareOrder.
func compare[K graphqlCursorType](a, b K) CompareOrder {
	anyA, anyB := any(a), any(b)

	switch t := anyA.(type) {
	case time.Time:
		bTime := anyB.(time.Time)
		switch {
		case t.Before(bTime):
			return Lesser
		case t.After(bTime):
			return Greater
		}
		return Equal

	case decimal.Decimal:
		deciB := anyB.(decimal.Decimal)
		switch {
		case t.LessThan(deciB):
			return Lesser
		case t.GreaterThan(deciB):
			return Greater
		}
		return Equal

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

type GraphqlParams struct {
	Before *string
	After  *string
	First  *int32
	Last   *int32
}

func (s *graphqlPaginator[T, K]) Len() int {
	return len(s.data)
}

func (s *graphqlPaginator[T, K]) Less(i, j int) bool {
	return compare(s.keyFunc(s.data[i]), s.keyFunc(s.data[j])) == Lesser
}

func (s *graphqlPaginator[T, K]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

const PaginationError = "api.graphql.pagination_params.invalid.app_error"

// graphqlPaginator implements sort.Interface
type graphqlPaginator[
	T any,
	K graphqlCursorType,
] struct {
	data    []T
	keyFunc func(T) K
	GraphqlParams
}

func (g *graphqlPaginator[T, K]) parse(apiName string) (data []T, hasPreviousPage bool, hasNextPage bool, appErr *model.AppError) {
	if (g.First != nil && g.Last != nil) || (g.First == nil && g.Last == nil) {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "First / Last"}, "provide either First or Last, not both", http.StatusBadRequest)
		return
	}
	if g.First != nil && g.Before != nil {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "First / Before"}, "First and Before can't go together", http.StatusBadRequest)
		return
	}
	if g.Last != nil && g.After != nil {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Last / After"}, "Last and After can't go together", http.StatusBadRequest)
		return
	}
	if g.Before != nil && g.After != nil {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Before / After"}, "Before and After can'g go together", http.StatusBadRequest)
		return
	}

	orderASC := g.First != nil // order ascending or not

	if orderASC {
		sort.Sort(g)
	} else {
		sort.Sort(sort.Reverse(g))
	}

	operand, err := g.parseOperand()
	if err != nil {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "Before / After"}, err.Error(), http.StatusInternalServerError)
		return
	}

	if operand == nil {
		if orderASC {
			if *g.First < int32(g.Len()) { // prevent slicing out of range
				data = g.data[:*g.First]
				hasNextPage = true
			} else {
				data = g.data[:]
			}
			return
		}

		// order desc:
		if *g.Last < int32(g.Len()) { // prevent slicing out of range
			data = g.data[:*g.Last]
			hasNextPage = true
		} else {
			data = g.data[:]
		}
		return
	}

	// case operand provided:
	index := sort.Search(g.Len(), func(i int) bool {
		value := g.keyFunc(g.data[i])
		cmp := compare(value, *operand)

		if orderASC {
			return cmp == Greater || cmp == Equal // >=
		}
		return cmp == Lesser || cmp == Equal // <=
	})

	// if not found, sort.Search returns exactly First int argument passed
	// we need to check it here
	if index >= g.Len() {
		appErr = model.NewAppError(apiName, PaginationError, map[string]interface{}{"Fields": "before / after"}, "invalid before or after provided", http.StatusBadRequest)
		return
	}

	hasPreviousPage = true
	data = g.data[index+1:]

	if orderASC {
		if *g.First < int32(len(data)) {
			data = data[:*g.First]
			hasNextPage = true
		}
		return
	}

	if *g.Last < int32(len(data)) {
		data = data[:*g.Last]
		hasNextPage = true
	}
	return
}

func reportingPeriodToDate(period ReportingPeriod) time.Time {
	now := time.Now()

	switch period {
	case ReportingPeriodToday:
		return util.StartOfDay(now)
	case ReportingPeriodThisMonth:
		return util.StartOfMonth(now)
	default:
		return now
	}
}
