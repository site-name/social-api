package gqlmodel

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// JSONString implements JSONString custom graphql scalar type
type JSONString map[string]interface{}

func (JSONString) ImplementsGraphQLType(name string) bool {
	return name == "JSONString"
}

func (j *JSONString) UnmarshalGraphQL(input interface{}) error {
	switch t := input.(type) {
	case model.StringInterface:
		*j = JSONString(t)
	case map[string]interface{}:
		*j = t

	default:
		return fmt.Errorf("wrong type: %T", t)
	}

	return nil
}

// PositiveDecimal implements JSONString custom graphql scalar type
type PositiveDecimal decimal.Decimal

func (PositiveDecimal) ImplementsGraphQLType(name string) bool {
	return name == "PositiveDecimal"
}

func (j *PositiveDecimal) UnmarshalGraphQL(input interface{}) error {
	if input == nil {
		return errors.New("input can't be nil")
	}

	var (
		deci decimal.Decimal
		err  error
	)

	value := reflect.ValueOf(input)

	switch value.Kind() {
	case reflect.String:
		deci, err = decimal.NewFromString(value.String())

	case reflect.Int,
		reflect.Int16,
		reflect.Int8,
		reflect.Int32,
		reflect.Int64:
		deci = decimal.NewFromInt(value.Int())

	case reflect.Float32,
		reflect.Float64:
		deci = decimal.NewFromFloat(value.Float())

	default:
		return fmt.Errorf("invalid input type: %T", input)
	}

	if err != nil {
		return err
	}
	if deci.LessThan(decimal.Zero) {
		return errors.New("positive decimal can't be less then zero")
	}

	*j = PositiveDecimal(deci)

	return nil
}

// Date implementes custom graphql scalar Date
// Date includes (Year, Month, Date) only
type Date struct {
	DateTime
}

// DateTime implementes custom graphql scalar DateTime
type DateTime struct {
	time.Time
}

func (Date) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

func (DateTime) ImplementsGraphQLType(name string) bool {
	return name == "DateTime"
}

func (j *DateTime) UnmarshalGraphQL(input interface{}) error {
	if input == nil {
		return errors.New("input can't be nil")
	}

	var err error

	switch t := input.(type) {
	case string:
		j.Time, err = time.Parse(time.RFC3339, t)
	case []byte:
		j.Time, err = time.Parse(time.RFC3339, string(t))
	default:
		return fmt.Errorf("invalid input type: %T", input)
	}

	return err
}

func (j *Date) UnmarshalGraphQL(input interface{}) error {
	err := j.DateTime.UnmarshalGraphQL(input)
	if err != nil {
		return err
	}

	j.Time = util.StartOfDay(j.Time)

	return nil
}
