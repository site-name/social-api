package scalars

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/measurement"
)

func MarshalWeightScalar(mass *measurement.Weight) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		var err error
		if mass.Unit != measurement.STANDARD_WEIGHT_UNIT {
			mass, err = mass.ConvertTo(measurement.STANDARD_WEIGHT_UNIT)
		}
		if err != nil {
			w.Write([]byte(measurement.ZeroWeight.String()))
		} else {
			w.Write([]byte(mass.String()))
		}
	})
}

func UnmarshalWeightScalar(v interface{}) (*measurement.Weight, error) {
	var weight *measurement.Weight
	var err error

	switch v := v.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[strings.ToLower(key)] = value
		}
		unit, ok1 := v["unit"]
		amount, ok2 := v["amount"]

		if ok1 && ok2 {
			weight = &measurement.Weight{
				Unit:   measurement.WeightUnit(strings.ToLower(unit.(string))),
				Amount: model.NewFloat32(amount.(float32)),
			}
		} else {
			err = errors.New("both 'amount' and 'unit' must be provided")
		}
	case string:
		err = json.JSON.Unmarshal([]byte(strings.ToLower(v)), &weight)

	default:
		err = fmt.Errorf("value of type %T is not supported", v)
	}

	return weight, err
}
