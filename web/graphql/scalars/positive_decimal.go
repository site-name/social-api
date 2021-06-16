package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

func MarshalPositiveDecimal(d *decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if d == nil {
			w.Write([]byte("null"))
		} else {
			w.Write([]byte(d.String()))
		}
	})
}

func UnmarshalPositiveDecimal(v interface{}) (*decimal.Decimal, error) {
	var deci decimal.Decimal
	var err error

	switch v := v.(type) {
	case string:
		de, er := decimal.NewFromString(v)
		if er != nil {
			err = er
		} else {
			deci = de
		}
	case int:
		deci = decimal.NewFromInt32(int32(v))
	case int32:
		deci = decimal.NewFromInt32(v)
	case float64:
		deci = decimal.NewFromFloat(v)
	case float32:
		deci = decimal.NewFromFloat32(v)

	default:
		err = fmt.Errorf("%v with type: %T is not a decimal", v, v)
	}

	if err != nil {
		return nil, err
	}

	if deci.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("decimal must be positive")
	}

	return &deci, nil
}
