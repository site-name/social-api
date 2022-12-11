package api

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
)

//go:embed schemas
var assets embed.FS

// ErrorUnauthorized
const ErrorUnauthorized = "api.unauthorized.app_error"

// Unique type to hold our context.
type CTXKey int

const (
	WebCtx CTXKey = iota
)

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

		_, err = builder.Write(data)
		if err != nil {
			return "", errors.Wrap(err, "failed to build up schema files")
		}

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

	if len(m) == 0 {
		return []*MetadataItem{}
	}

	i := 0
	res := make([]*MetadataItem, len(m))
	for key, value := range m {
		res[i] = &MetadataItem{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		}
		i++
	}

	return res
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

	res := &TaxedMoney{
		Currency: money.Currency,
		Gross:    SystemMoneyToGraphqlMoney(money.Gross),
		Net:      SystemMoneyToGraphqlMoney(money.Net),
	}
	tax, _ := money.Tax()
	res.Tax = SystemMoneyToGraphqlMoney(tax)

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

type GraphqlPaginationOptions struct {
	Before         *string
	After          *string
	First          *int32
	Last           *int32
	OrderBy        string
	OrderDirection OrderDirection
}

const GraphqlPaginationError = "api.graphql.pagination_params.invalid.app_error"

func (g *GraphqlPaginationOptions) isValid() *model.AppError {
	if strings.TrimSpace(g.OrderBy) == "" {
		return model.NewAppError("GraphqlPaginationOptions.IsValid", GraphqlPaginationError, map[string]interface{}{"Fields": "OrderBy"}, "You must provide order by", http.StatusBadRequest)
	}
	if g.First != nil && g.Last != nil {
		return model.NewAppError("GraphqlPaginationOptions.IsValid", GraphqlPaginationError, map[string]interface{}{"Fields": "Last, First"}, "You must provide either First or Last, not both", http.StatusBadRequest)
	}
	if g.First != nil && g.Before != nil {
		return model.NewAppError("GraphqlPaginationOptions.IsValid", GraphqlPaginationError, map[string]interface{}{"Fields": "First, Before"}, "First and Before can't go together", http.StatusBadRequest)
	}
	if g.Last != nil && g.After != nil {
		return model.NewAppError("GraphqlPaginationOptions.IsValid", GraphqlPaginationError, map[string]interface{}{"Fields": "Last, After"}, "Last and After can't go together", http.StatusBadRequest)
	}

	return nil
}

// Decode decodes before or after from base64 format to initial form, then returns the result
func (g *GraphqlPaginationOptions) decode() (string, *model.AppError) {
	var value string

	if g.Before != nil {
		value = *g.Before
	} else if g.After != nil {
		value = *g.After
	}

	// these before and after values are created by code.
	// they will not be changed by human, so we can safely ignore the error here
	res, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", model.NewAppError("GraphqlPaginationOptions.decode", GraphqlPaginationError, map[string]interface{}{"Fields": "Before/After"}, "Invalid cursor provided", http.StatusBadRequest)
	}

	return string(res), nil
}

// ConstructSqlizer does:
//
// 1) check if arguments are provided properly
//
// 2) decodes given before or after cursor
//
// 3) construct a squirrel expression based on given key
func (g *GraphqlPaginationOptions) ConstructSqlizer() (squirrel.Sqlizer, error) {
	if err := g.isValid(); err != nil {
		return nil, err
	}

	cmp, err := g.decode()
	if err != nil {
		return nil, err
	}

	switch {
	case g.After != nil:
		if g.OrderDirection == OrderDirectionAsc {
			// 1 2 3 4 5 6 (ASC)
			//     | *     (AFTER)
			return squirrel.Gt{g.OrderBy: cmp}, nil
		}

		// 6 5 4 3 2 1 (DESC)
		//       | *   (AFTER)
		return squirrel.Lt{g.OrderBy: cmp}, nil

	case g.Before != nil:
		if g.OrderDirection == OrderDirectionAsc {
			// 1 2 3 4 5 6 (ASC)
			//   * |       (BEFORE)
			return squirrel.Lt{g.OrderBy: cmp}, nil
		}

		// 6 5 4 3 2 1 (DESC)
		//     * |     (BEFORE)
		return squirrel.Gt{g.OrderBy: cmp}, nil

	default:
		return squirrel.Expr(""), nil
	}
}

// If -1, means no limit
func (g *GraphqlPaginationOptions) Limit() int32 {
	if g.First != nil {
		return *g.First
	} else if g.Last != nil {
		return *g.Last
	}

	return -1
}

func (g *GraphqlPaginationOptions) HasPreviousPage() bool {
	return (g.First != nil && g.After != nil) || (g.Last != nil && g.Before != nil)
}
