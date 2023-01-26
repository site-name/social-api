package api

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

//go:embed schemas
var assets embed.FS

// ErrorUnauthorized
const ErrorUnauthorized = "api.unauthorized.app_error"

// Unique type to hold our context.
type CTXKey int

const WebCtx CTXKey = iota

// constructSchema constructs schema from *.graphql files
func constructSchema() (string, error) {
	entries, err := assets.ReadDir("schemas")
	if err != nil {
		return "", errors.Wrap(err, "failed to read schema dir")
	}

	var builder strings.Builder
	for _, entry := range entries {
		if entry.IsDir() || !(strings.HasSuffix(entry.Name(), ".graphql") || strings.HasSuffix(entry.Name(), ".graphqls")) {
			continue
		}
		data, err := assets.ReadFile(filepath.Join("schemas", entry.Name()))
		if err != nil {
			return "", errors.Wrapf(err, "failed to read schema file: %s", filepath.Join("schemas", entry.Name()))
		}

		builder.Write(data)
		builder.WriteByte('\n')
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

	res := &Money{
		Currency: money.Currency,
	}
	res.Amount, _ = money.Amount.Float64()

	return res
}

func SystemTaxedMoneyToGraphqlTaxedMoney(money *goprices.TaxedMoney) *TaxedMoney {
	if money == nil {
		return nil
	}

	tax, _ := money.Tax()
	res := &TaxedMoney{
		Currency: money.Currency,
		Gross:    SystemMoneyToGraphqlMoney(money.Gross),
		Net:      SystemMoneyToGraphqlMoney(money.Net),
		Tax:      SystemMoneyToGraphqlMoney(tax),
	}

	return res
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

func (s *graphqlPaginator[T, K]) Len() int {
	return len(s.data)
}

func (s *graphqlPaginator[T, K]) Less(i, j int) bool {
	return s.keyFunc(s.data[i]) < s.keyFunc(s.data[j])
}

func (s *graphqlPaginator[T, K]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

// graphqlPaginator implements sort.Interface
type graphqlPaginator[T any, K util.Ordered] struct {
	data    []T
	keyFunc func(T) K

	before *K
	after  *K
	first  *int32
	last   *int32
}

func (g *graphqlPaginator[T, K]) parse(apiName string) (data []T, hasPreviousPage bool, hasNextPage bool, err *model.AppError) {
	if (g.first != nil && g.last != nil) || (g.first == nil && g.last == nil) {
		return nil, false, false, model.NewAppError(apiName, model.PaginationError, map[string]interface{}{"Fields": "first / last"}, "provide either first or last, not both", http.StatusBadRequest)
	}
	if g.first != nil && g.before != nil {
		return nil, false, false, model.NewAppError(apiName, model.PaginationError, map[string]interface{}{"Fields": "first / before"}, "first and before can't go together", http.StatusBadRequest)
	}
	if g.last != nil && g.after != nil {
		return nil, false, false, model.NewAppError(apiName, model.PaginationError, map[string]interface{}{"Fields": "last / after"}, "last and after can't go together", http.StatusBadRequest)
	}
	if g.before != nil && g.after != nil {
		return nil, false, false, model.NewAppError(apiName, model.PaginationError, map[string]interface{}{"Fields": "before / after"}, "before and after can'g go together", http.StatusBadRequest)
	}

	orderASC := g.first != nil // order ascending or not

	if orderASC {
		sort.Sort(g)
	} else {
		sort.Sort(sort.Reverse(g))
	}

	var operand K
	if g.before != nil {
		operand = *g.before
	} else if g.after != nil {
		operand = *g.after
	}

	var emptyK K
	if operand == emptyK {
		if orderASC {
			if *g.first < int32(g.Len()) { // prevent slicing out of range
				data = g.data[:*g.first]
				hasNextPage = true
			} else {
				data = g.data[:]
			}

			return
		}

		// order desc:

		if *g.last < int32(g.Len()) { // prevent slicing out of range
			data = g.data[:*g.last]
			hasPreviousPage = true
		} else {
			data = g.data[:]
		}

		return
	}

	// case operand provided:

	index := sort.Search(g.Len(), func(i int) bool {
		value := g.keyFunc(g.data[i])
		if orderASC {
			return value >= operand
		}
		return value <= operand
	})

	// if not found, sort.Search returns exactly first int argument passed
	// we need to check it here
	if index >= g.Len() {
		return nil, false, false, model.NewAppError(apiName, model.PaginationError, map[string]interface{}{"Fields": "before / after"}, "invalid before or after provided", http.StatusBadRequest)
	}

	if orderASC {
		data = g.data[index+1:]
		hasPreviousPage = true

		if *g.first < int32(len(data)) {
			data = data[:*g.first]
			hasNextPage = true
		}
		return
	}

	data = g.data[index+1:]
	hasPreviousPage = true

	if *g.last < int32(len(data)) {
		data = data[:*g.last]
		hasNextPage = true
	}
	return
}
