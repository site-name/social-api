package model_types

import (
	"bytes"
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

type JsonMap map[string]any

func (j *JsonMap) Scan(value any) error {
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

func (j JsonMap) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	if len(data) > maxPropSizeBytes {
		return nil, ErrMaxPropSizeExceeded
	}

	return string(data), err
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

func (n *NullInt) IsNil() bool {
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

	switch t := value.(type) {
	case int:
		deci := decimal.NewFromInt(int64(t))
		n.Decimal = &deci
		return nil

	case int32:
		deci := decimal.NewFromInt32(t)
		n.Decimal = &deci
		return nil

	case int64:
		deci := decimal.NewFromInt(t)
		n.Decimal = &deci
		return nil

	case decimal.Decimal:
		n.Decimal = &t
		return nil

	case float64:
		deci := decimal.NewFromFloat(t)
		n.Decimal = &deci
		return nil

	case float32:
		deci := decimal.NewFromFloat32(t)
		n.Decimal = &deci
		return nil

	case string:
		deci, err := decimal.NewFromString(t)
		if err != nil {
			return err
		}
		n.Decimal = &deci
		return nil

	default:
		return fmt.Errorf("unsupported value with type: %T", value)
	}
}

func (n NullDecimal) Value() (driver.Value, error) {
	if n.Decimal == nil {
		return nil, nil
	}
	return *n.Decimal, nil
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
