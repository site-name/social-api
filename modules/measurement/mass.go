package measurement

import (
	"errors"
	"fmt"
)

type WeightUnit string

var (
	WEIGHT_UNIT_STRINGS = map[WeightUnit]string{
		G:     "Gram",
		LB:    "Pound",
		OZ:    "Ounce",
		KG:    "Kg",
		TONNE: "Tonne",
	} // weight unit aliases to their full name
	WEIGHT_UNIT_CONVERSION = map[WeightUnit]float64{
		KG:    1.0,
		G:     1000.0,
		OZ:    35.27396195,
		TONNE: 0.001,
		LB:    2.20462262,
	} // amount of weight units
	ErrInvalidWeightUnit = errors.New("invalid weight unit, must be either (g|lb|oz|kg|tonne)") // Error used when users use weight unit does not match type WeightUnit
)

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
	Amount float64    `json:"amount"`
	Unit   WeightUnit `json:"unit"`
}

func (w Weight) String() string {
	return fmt.Sprintf("%.3f %s", w.Amount, w.Unit)
}

// Adds weight to current weight and returns new weight.
func (w Weight) Add(other Weight) (*Weight, error) {
	// convert other's unit to w's unit
	converted, err := other.ConvertTo(w.Unit)
	if err != nil {
		return nil, err
	}
	return &Weight{
		Amount: w.Amount + converted.Amount,
		Unit:   w.Unit,
	}, nil
}

// Subs weight to current weight and returns new weight.
func (w Weight) Sub(other Weight) (*Weight, error) {
	// convert other's unit to w's unit
	converted, err := other.ConvertTo(w.Unit)
	if err != nil {
		return nil, err
	}
	return &Weight{
		Amount: w.Amount - converted.Amount,
		Unit:   w.Unit,
	}, nil
}

// Multiplies weight to current weight and returns new weight.
func (w Weight) Mul(quantity int) Weight {
	return Weight{
		Amount: w.Amount * float64(quantity),
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

	w.Amount = w.Amount / WEIGHT_UNIT_CONVERSION[w.Unit] * WEIGHT_UNIT_CONVERSION[unit]
	w.Unit = unit
	return w, nil
}

// Zero weight for unit (kg)
var ZeroWeight = Weight{
	Amount: 0,
	Unit:   KG,
}

// NewWeight returns a customized weight user wants
func NewWeight(amount float64, unit WeightUnit) (*Weight, error) {
	if WEIGHT_UNIT_STRINGS[unit] == "" {
		return nil, ErrInvalidWeightUnit
	}
	return &Weight{
		Amount: amount,
		Unit:   unit,
	}, nil
}
