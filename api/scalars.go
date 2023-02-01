package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	switch t := input.(type) {
	case string:
		deci, err = decimal.NewFromString(t)
	case int:
		deci = decimal.NewFromInt32(int32(t))
	case int32:
		deci = decimal.NewFromInt32(t)
	case float64:
		deci = decimal.NewFromFloat(t)
	case decimal.Decimal:
		deci = t
	default:
		err = fmt.Errorf("unexpected input value's type: %T", input)
	}

	if err != nil {
		return err
	}
	if deci.LessThan(decimal.Zero) {
		return errors.New("positive decimal can't be less than zero")
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
	case Warehouse, ShippingMethod, *Warehouse, *ShippingMethod:
		return true

	default:
		return false
	}
}
