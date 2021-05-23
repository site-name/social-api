package scalars

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/slog"
)

// ustom Decimal implementation.
//
// Returns Decimal as a float in the API,
// parses float to the Decimal on the way back
type Decimal struct {
	decimal.Decimal
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Decimal) UnmarshalGQL(v interface{}) error {
	var err error
	switch t := v.(type) {
	case decimal.Decimal:
		*d = Decimal{Decimal: t}
	case string:
		parsed, er := decimal.NewFromString(t)
		if er != nil {
			err = er
		} else {
			*d = Decimal{Decimal: parsed}
		}
	case int64:
		*d = Decimal{Decimal: decimal.NewFromInt(t)}
	case int32:
		*d = Decimal{Decimal: decimal.NewFromInt32(t)}
	case float32:
		*d = Decimal{Decimal: decimal.NewFromFloat32(t)}
	case float64:
		*d = Decimal{Decimal: decimal.NewFromFloat(t)}
	case int:
		*d = Decimal{Decimal: decimal.NewFromInt32(int32(t))}
	case int16:
		*d = Decimal{Decimal: decimal.NewFromInt32(int32(t))}
	default:
		err = fmt.Errorf("unknown type %T cannot be unmarshaled", v)
	}

	return err
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Decimal) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(d.String()))
	if err != nil {
		slog.Error("error marshaling decimal", slog.Err(err))
	}
}

// Positive Decimal scalar implementation.
//
// Should be used in places where value must be positive.
// type PositiveDecimal struct {
// 	Decimal
// }

// func (d *PositiveDecimal) UnmarshalGQL(v interface{}) error {
// 	err := d.Decimal.UnmarshalGQL(v)
// 	if err != nil {
// 		return err
// 	}
// 	if d.LessThan(decimal.Zero) {
// 		return fmt.Errorf("decimal must be positive")
// 	}
// 	return nil
// }

type WeightScalar struct {
	measurement.Weight
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *WeightScalar) UnmarshalGQL(v interface{}) error {
	var err error
	var amount float32

	switch t := v.(type) {
	case float32:
		amount = t
	case int:
		amount = float32(t)
	case int32:
		amount = float32(t)
	case string:
		parsed, er := strconv.ParseFloat(t, 32)
		if er != nil {
			err = er
		} else {
			amount = float32(parsed)
		}
	default:
		err = fmt.Errorf("type %T cannot be unmarshaled", v)
	}

	if err == nil {
		*d = WeightScalar{
			Weight: measurement.Weight{
				Amount: amount,
				Unit:   measurement.STANDARD_WEIGHT_UNIT,
			},
		}
	}

	return err
}

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
		err = fmt.Errorf("%T is not a decimal", v)
	}

	if err != nil {
		return nil, err
	}

	if deci.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("decimal must be positive")
	}

	return &deci, nil
}
