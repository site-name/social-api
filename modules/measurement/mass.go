package measurement

import (
	"errors"
	"fmt"
)

type WeightUnit string

var (
	WEIGHT_UNIT_STRINGS    map[WeightUnit]string  // weight unit aliases to their full name
	WEIGHT_UNIT_CONVERSION map[WeightUnit]float32 // amount of weight units
	ErrInvalidWeightUnit   error                  // Error used when users use weight unit does not match type WeightUnit
)

func init() {
	WEIGHT_UNIT_STRINGS = map[WeightUnit]string{
		G:     "Gram",
		LB:    "Pound",
		OZ:    "Ounce",
		KG:    "Kg",
		TONNE: "Tonne",
	}
	WEIGHT_UNIT_CONVERSION = map[WeightUnit]float32{
		KG:    1.0,
		G:     1000.0,
		OZ:    35.27396195,
		TONNE: 0.001,
		LB:    2.20462262,
	}
	ErrInvalidWeightUnit = errors.New("invalid weight unit, must be either (g|lb|oz|kg|tonne)")
}

// weight units supported by app
const (
	G     WeightUnit = "g"
	LB    WeightUnit = "lb"
	OZ    WeightUnit = "oz"
	KG    WeightUnit = "kg"
	TONNE WeightUnit = "tonne"
)

const STANDARD_WEIGHT_UNIT = KG

type Weight struct {
	Amount *float32   `json:"amount"`
	Unit   WeightUnit `json:"unit"`
}

func (w *Weight) String() string {
	return fmt.Sprintf("%.3f %s", *w.Amount, w.Unit)
}

func newFloat32(f float32) *float32 {
	return &f
}

// Adds weight to current weight and returns new weight.
func (w *Weight) Add(other *Weight) (*Weight, error) {
	// convert other's unit to w's unit
	converted, err := other.ConvertTo(w.Unit)
	if err != nil {
		return nil, err
	}
	return &Weight{
		Amount: newFloat32(*w.Amount + *converted.Amount),
		Unit:   w.Unit,
	}, nil
}

// Subs weight to current weight and returns new weight.
func (w *Weight) Sub(other *Weight) (*Weight, error) {
	// convert other's unit to w's unit
	converted, err := other.ConvertTo(w.Unit)
	if err != nil {
		return nil, err
	}
	return &Weight{
		Amount: newFloat32(*w.Amount - *converted.Amount),
		Unit:   w.Unit,
	}, nil
}

// Multiplies weight to current weight and returns new weight.
func (w *Weight) Mul(quantity float32) *Weight {
	return &Weight{
		Amount: newFloat32(*w.Amount * quantity),
		Unit:   w.Unit,
	}
}

// converts current weight to weight with given unit. Error could be ErrInvalidWeightUnit or nil
func (w *Weight) ConvertTo(unit WeightUnit) (*Weight, error) {

	// check if given unit is supported by system
	if WEIGHT_UNIT_STRINGS[unit] == "" {
		return nil, ErrInvalidWeightUnit
	}
	if unit == w.Unit {
		return w, nil
	}

	resAmount := *w.Amount / WEIGHT_UNIT_CONVERSION[w.Unit] * WEIGHT_UNIT_CONVERSION[unit]
	return &Weight{
		Amount: newFloat32(resAmount),
		Unit:   unit,
	}, nil
}

// Zero weight for unit (kg)
var ZeroWeight = &Weight{
	Amount: newFloat32(0),
	Unit:   KG,
}

// NewWeight returns a customized weight user wants
func NewWeight(amount float32, unit WeightUnit) (*Weight, error) {
	if WEIGHT_UNIT_STRINGS[unit] == "" {
		return nil, ErrInvalidWeightUnit
	}
	return &Weight{
		Amount: newFloat32(amount),
		Unit:   unit,
	}, nil
}
