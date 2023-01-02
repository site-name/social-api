package api

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
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
