package api

import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// JSONString implements JSONString custom graphql scalar type
type JSONString map[string]any

func (JSONString) ImplementsGraphQLType(name string) bool {
	return name == "JSONString"
}

func (j *JSONString) UnmarshalGraphQL(input any) error {
	switch t := input.(type) {
	case model.StringInterface:
		*j = JSONString(t)
	case map[string]any:
		*j = t

	default:
		return fmt.Errorf("wrong type: %T", t)
	}

	return nil
}

type UUID string

func (UUID) ImplementsGraphQLType(name string) bool {
	return name == "UUID"
}

func (u UUID) String() string {
	return *(*string)(unsafe.Pointer(&u))
}

func (j *UUID) UnmarshalGraphQL(input any) error {
	switch t := input.(type) {
	case string:
		uid, err := uuid.Parse(t)
		if err != nil {
			return errors.Wrap(err, "failed to parse uuid value")
		}
		strUid := uid.String()
		*j = *(*UUID)(unsafe.Pointer(&strUid))
		return nil

	case []byte:
		uid, err := uuid.ParseBytes(t)
		if err != nil {
			return errors.Wrap(err, "failed to parse uuid value")
		}
		strUid := uid.String()
		*j = *(*UUID)(unsafe.Pointer(&strUid))
		return nil

	default:
		return fmt.Errorf("unsupported input type: %T", input)
	}
}

// PositiveDecimal implements custom graphql scalar type
type PositiveDecimal decimal.Decimal

func (p PositiveDecimal) ToDecimal() decimal.Decimal {
	return *(*decimal.Decimal)(unsafe.Pointer(&p))
}

func (PositiveDecimal) ImplementsGraphQLType(name string) bool {
	return name == "PositiveDecimal"
}

func (j *PositiveDecimal) UnmarshalGraphQL(input any) error {
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

func (DateTime) ImplementsGraphQLType(name string) bool {
	return name == "DateTime"
}

func (j *DateTime) UnmarshalGraphQL(input any) error {
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

func (j *Date) UnmarshalGraphQL(input any) error {
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

type WeightScalar struct {
	Value float64         `json:"value"`
	Unit  WeightUnitsEnum `json:"unit"`
}

func (w *WeightScalar) UnmarshalGraphQL(input any) error {
	switch v := input.(type) {
	case []byte:
		return json.Unmarshal(v, w)
	case string:
		return json.Unmarshal([]byte(v), w)
	default:
		return fmt.Errorf("unsupported value type: %T", input)
	}
}

func (w WeightScalar) ImplementsGraphQLType(name string) bool {
	return name == "WeightScalar"
}
