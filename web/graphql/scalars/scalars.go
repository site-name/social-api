package scalars

import (
	"fmt"
	"strconv"

	"github.com/sitename/sitename/modules/measurement"
)

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
