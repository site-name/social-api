package measurement

import (
	"errors"
)

type WeightUnit string

// weight units supported by app
const (
	G     WeightUnit = "g"
	LB    WeightUnit = "lb"
	OZ    WeightUnit = "oz"
	KG    WeightUnit = "kg"
	TONNE WeightUnit = "tonne"
)

// weight unit aliases to their full name
var WEIGHT_UNIT_STRINGS = map[WeightUnit]string{
	G:     "Gram",
	LB:    "Pound",
	OZ:    "Ounce",
	KG:    "kg",
	TONNE: "Tonne",
}

// amount of weight units
var WEIGHT_UNIT_CONVERSION = map[WeightUnit]float32{
	KG:    1.0,
	G:     1000.0,
	OZ:    35.27396195,
	TONNE: 0.001,
	LB:    2.20462262,
}

const STANDARD_WEIGHT_UNIT = KG

type Weight struct {
	Amount float32
	Unit   WeightUnit
}

var (
	// Error used when users use weight unit does not match type WeightUnit
	ErrInvalidWeightUnit = errors.New("invalid weight unit, must be either (g|lb|oz|kg|tonne)")
)

// Adds weight to current weight and returns new weight.
func (w *Weight) Add(other *Weight) *Weight {
	// convert other's unit to w's unit
	converted, _ := other.ConvertTo(w.Unit)
	return &Weight{
		Amount: w.Amount + converted.Amount,
		Unit:   w.Unit,
	}
}

// Subs weight to current weight and returns new weight.
func (w *Weight) Sub(other *Weight) *Weight {
	// convert other's unit to w's unit
	converted, _ := other.ConvertTo(w.Unit)
	return &Weight{
		Amount: w.Amount - converted.Amount,
		Unit:   w.Unit,
	}
}

// Multiplies weight to current weight and returns new weight.
func (w *Weight) Mul(quantity float32) *Weight {
	return &Weight{
		Amount: w.Amount * quantity,
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

	resAmount := w.Amount / WEIGHT_UNIT_CONVERSION[w.Unit] * WEIGHT_UNIT_CONVERSION[unit]
	return &Weight{
		Amount: resAmount,
		Unit:   unit,
	}, nil
}

// Zero weight for unit (kg)
var ZeroWeight = &Weight{
	Amount: 0,
	Unit:   KG,
}

// NewWeight returns a customized weight user wants
func NewWeight(amount float32, unit WeightUnit) *Weight {
	return &Weight{
		Amount: amount,
		Unit:   unit,
	}
}
