package model_types

import (
	"bytes"
	"cmp"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/site-name/decimal"
)

const maxPropSizeBytes = 1024 * 1024

var NullBytes = []byte("null")
var ErrMaxPropSizeExceeded = fmt.Errorf("max prop size of %d exceeded", maxPropSizeBytes)

type JSONString map[string]any

func (j *JSONString) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case []byte:
		return json.Unmarshal(t, j)
	case string:
		return json.Unmarshal([]byte(t), j)
	default:
		return errors.New("received value is neither a byte slice or sttring")
	}
}

func (j JSONString) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	if len(data) > maxPropSizeBytes {
		return nil, ErrMaxPropSizeExceeded
	}

	return string(data), err
}

func (j JSONString) Get(key string, defaultValue ...any) any {
	value, exist := j[key]
	if exist {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (JSONString) ImplementsGraphQLType(name string) bool {
	return name == "JSONString"
}

func (j *JSONString) UnmarshalGraphQL(input any) error {
	switch t := input.(type) {
	case JSONString:
		*j = t
	case map[string]any:
		*j = t

	default:
		return fmt.Errorf("wrong type: %T", t)
	}

	return nil
}

type NullInt64 struct {
	Int64 *int64
}

func (n NullInt64) IsNil() bool {
	return n.Int64 == nil
}

func NewNullInt64(value int64) NullInt64 {
	return NullInt64{&value}
}

func (n NullInt64) IsZero() bool {
	return n.IsNil()
}

func (n *NullInt64) Scan(value any) error {
	if value == nil {
		n.Int64 = nil
		return nil
	}

	switch t := value.(type) {
	case int:
		t64 := int64(t)
		n.Int64 = &t64
		return nil
	case int32:
		t64 := int64(t)
		n.Int64 = &t64
		return nil
	case int64:
		n.Int64 = &t
		return nil

	default:
		return fmt.Errorf("unsupported value with type: %T", value)
	}
}

func (n NullInt64) Value() (driver.Value, error) {
	if n.Int64 == nil {
		return nil, nil
	}
	return *n.Int64, nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Int64 == nil {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(*n.Int64, 10)), nil
}

func (n *NullInt64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(NullBytes, data) {
		n.Int64 = nil
		return nil
	}
	return json.Unmarshal(data, &n.Int64)
}

type NullInt struct {
	Int *int
}

func (n NullInt) IsNil() bool {
	return n.Int == nil
}

func NewNullInt(value int) NullInt {
	return NullInt{&value}
}

func (n *NullInt) Scan(value any) error {
	if value == nil {
		n.Int = nil
		return nil
	}

	switch t := value.(type) {
	case int:
		n.Int = &t
		return nil
	case int32:
		tint := int(t)
		n.Int = &tint
		return nil
	case int64:
		tint := int(t)
		n.Int = &tint
		return nil

	default:
		return fmt.Errorf("unsupported value with type: %T", value)
	}
}

func (n NullInt) IsZero() bool {
	return n.IsNil()
}

func (n NullInt) Value() (driver.Value, error) {
	if n.Int == nil {
		return nil, nil
	}
	return *n.Int, nil
}

func (n NullInt) MarshalJSON() ([]byte, error) {
	if n.Int == nil {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(*n.Int), 10)), nil
}

func (n *NullInt) UnmarshalJSON(data []byte) error {
	if bytes.Equal(NullBytes, data) {
		n.Int = nil
		return nil
	}
	return json.Unmarshal(data, &n.Int)
}

type NullDecimal struct {
	Decimal *decimal.Decimal
}

func (n NullDecimal) IsZero() bool {
	return n.IsNil()
}

func NewNullDecimal(value decimal.Decimal) NullDecimal {
	return NullDecimal{&value}
}

func (n NullDecimal) IsNil() bool {
	return n.Decimal == nil
}

func (n *NullDecimal) Scan(value any) error {
	if value == nil {
		n.Decimal = nil
		return nil
	}

	var deci decimal.Decimal
	err := deci.Scan(value)
	if err != nil {
		return err
	}

	n.Decimal = &deci
	return nil
}

func (n NullDecimal) Value() (driver.Value, error) {
	if n.Decimal == nil {
		return nil, nil
	}
	return n.Decimal.InexactFloat64(), nil
}

func (n NullDecimal) MarshalJSON() ([]byte, error) {
	if n.Decimal == nil {
		return NullBytes, nil
	}
	return []byte(n.Decimal.String()), nil
}

func (n *NullDecimal) UnmarshalJSON(data []byte) error {
	if bytes.Equal(NullBytes, data) {
		n.Decimal = nil
		return nil
	}
	return json.Unmarshal(data, &n.Decimal)
}

type NullString struct {
	String *string
}

func (n NullString) IsNil() bool {
	return n.String == nil
}

func NewNullString(value string) NullString {
	return NullString{&value}
}

func (n NullString) IsZero() bool {
	return n.IsNil()
}

func (n *NullString) Scan(value any) error {
	if value == nil {
		n.String = nil
		return nil
	}

	switch t := value.(type) {
	case string:
		n.String = &t
		return nil
	case []byte:
		str := string(t)
		n.String = &str
		return nil
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
}

func (n NullString) Value() (driver.Value, error) {
	if n.String == nil {
		return nil, nil
	}
	return *n.String, nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.String == nil {
		return NullBytes, nil
	}
	return []byte(*n.String), nil
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(NullBytes, data) {
		n.String = nil
		return nil
	}
	return json.Unmarshal(data, &n.String)
}

type NullFloat32 struct {
	Float32 *float32
}

func (n NullFloat32) IsZero() bool {
	return n.IsNil()
}

func (f NullFloat32) IsNil() bool {
	return f.Float32 == nil
}

func NewNullFloat32(value float32) NullFloat32 {
	return NullFloat32{&value}
}

func (f *NullFloat32) Scan(value any) error {
	if value == nil {
		f.Float32 = nil
		return nil
	}

	switch t := value.(type) {
	case float32:
		f.Float32 = &t
		return nil
	case int:
		tf32 := float32(t)
		f.Float32 = &tf32
		return nil
	case int32:
		tf32 := float32(t)
		f.Float32 = &tf32
		return nil
	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
}

func (f NullFloat32) Value() (driver.Value, error) {
	if f.Float32 == nil {
		return nil, nil
	}

	return *f.Float32, nil
}

func (n NullFloat32) MarshalJSON() ([]byte, error) {
	if n.Float32 == nil {
		return NullBytes, nil
	}
	return []byte(fmt.Sprintf("%f", *n.Float32)), nil
}

func (n *NullFloat32) UnmarshalJSON(data []byte) error {
	if bytes.Equal(NullBytes, data) {
		n.Float32 = nil
		return nil
	}
	return json.Unmarshal(data, &n.Float32)
}

// PrimitiveIsNotNilAndEqual returns v != nil && *v == other
func PrimitiveIsNotNilAndEqual[T cmp.Ordered | ~bool](v *T, other T) bool {
	return v != nil && *v == other
}

// PrimitiveIsNotNilAndNotEqual returns v != nil && *v != other
func PrimitiveIsNotNilAndNotEqual[T cmp.Ordered | ~bool](v *T, other T) bool {
	return v != nil && *v != other
}

func (n NullBool) IsZero() bool {
	return n.IsNil()
}

type NullBool struct {
	Bool *bool
}

func (n NullBool) IsNil() bool {
	return n.Bool == nil
}

func NewNullBool(value bool) NullBool {
	return NullBool{&value}
}

func (n *NullBool) Scan(value any) error {
	if value == nil {
		n.Bool = nil
		return nil
	}

	b, ok := value.(bool)
	if ok {
		n.Bool = &b
		return nil
	}

	return fmt.Errorf("unsupported value type: %T", value)
}

func (n NullBool) Value() (driver.Value, error) {
	if n.Bool == nil {
		return nil, nil
	}

	return *n.Bool, nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.IsNil() {
		return NullBytes, nil
	}
	res := strconv.FormatBool(*n.Bool)
	return []byte(res), nil
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		n.Bool = nil
		return nil
	}
	return json.Unmarshal(data, &n.Bool)
}

type NullTime struct {
	Time *time.Time
}

func (n NullTime) IsZero() bool {
	return n.IsNil()
}

func (n NullTime) IsNil() bool {
	return n.Time == nil
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{&value}
}

func (n *NullTime) Scan(value any) error {
	if value == nil {
		n.Time = nil
		return nil
	}

	switch t := value.(type) {
	case time.Time:
		n.Time = &t
		return nil
	case string:
		tim, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return err
		}
		n.Time = &tim
		return nil
	case []byte:
		tim, err := time.Parse(time.RFC3339, string(t))
		if err != nil {
			return err
		}
		n.Time = &tim
		return nil

	default:
		return fmt.Errorf("unsupported value type: %T", value)
	}
}

func (n NullTime) Value() (driver.Value, error) {
	if n.IsNil() {
		return nil, nil
	}
	return *n.Time, nil
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.IsNil() {
		return NullBytes, nil
	}
	return []byte(n.Time.String()), nil
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		n.Time = nil
		return nil
	}
	return json.Unmarshal(data, &n.Time)
}

// NilTypeIsNotNilAndNotZero checks if given primitive pointer is not nil AND its point to value is not zero.
// E.g
//
//	var number int = 10
//	NilTypeIsNotNilAndNotZero(&number) == true
//
//	str := ""
//	NilTypeIsNotNilAndNotZero(&str) == false
func NilTypeIsNotNilAndNotZero[T cmp.Ordered](v *T) bool {
	var zeroValue T
	return v != nil && *v != zeroValue
}
