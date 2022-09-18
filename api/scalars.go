package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// UnmarshalGQL for gqlgen compartible
func (j *JSONString) UnmarshalGQL(v any) error {
	return j.UnmarshalGraphQL(v)
}

// MarshalGQL for gqlgen compartible
func (j *JSONString) MarshalGQL(w io.Writer) {
	data, err := json.Marshal(j)
	if err != nil {
		w.Write([]byte{'{', '}'})
		return
	}
	w.Write(data)
}

// PositiveDecimal implements custom graphql scalar type
type PositiveDecimal decimal.Decimal

// UnmarshalGQL for gqlgen compartible
func (j *PositiveDecimal) UnmarshalGQL(v any) error {
	return j.UnmarshalGraphQL(v)
}

// MarshalGQL for gqlgen compartible
func (j *PositiveDecimal) MarshalGQL(w io.Writer) {
	w.Write([]byte(decimal.Decimal(*j).String()))
}

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

	case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		deci = decimal.NewFromInt(int64(value.Uint()))

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

func (Date) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

// DateTime implementes custom graphql scalar DateTime
type DateTime struct {
	time.Time
}

// UnmarshalGQL for gqlgen compartible
func (j *DateTime) UnmarshalGQL(v any) error {
	return j.UnmarshalGraphQL(v)
}

// MarshalGQL for gqlgen compartible
func (j *DateTime) MarshalGQL(w io.Writer) {
	w.Write([]byte(j.Format(time.RFC3339)))
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

type TranslatableItem any

func IsValidTranslatableItem(v TranslatableItem) bool {
	switch v.(type) {
	case PageTranslatableContent,
		SaleTranslatableContent,
		VoucherTranslatableContent,
		ProductTranslatableContent,
		CategoryTranslatableContent,
		AttributeTranslatableContent,
		CollectionTranslatableContent,
		AttributeValueTranslatableContent,
		ProductVariantTranslatableContent,
		ShippingMethodTranslatableContent,
		MenuItemTranslatableContent:
		return true

	default:
		return false
	}
}

type DeliveryMethod any

func IsValidDeliveryMethod(v DeliveryMethod) bool {
	switch v.(type) {
	case Warehouse, ShippingMethod:
		return true

	default:
		return false
	}
}
