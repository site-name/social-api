package goprices

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money represents a money in real life, it includes Amount and currency
type Money struct {
	Amount   *decimal.Decimal
	Currency string
}

// NewMoney returns new Money object
func NewMoney(amount *decimal.Decimal, currency string) (*Money, error) {
	code, err := checkCurrency(currency)
	if err != nil {
		return nil, err
	}
	return &Money{
		Amount:   amount,
		Currency: code,
	}, nil
}

// String implements fmt.Stringer interface
func (m *Money) String() string {
	return fmt.Sprintf("Money{%q, %q}", m.Amount.String(), m.Currency)
}

// LessThan checks if other's amount is greater than m's amount
func (m *Money) LessThan(other *Money) (bool, error) {
	err := m.sameKind(other)
	if err != nil {
		return false, err
	}
	return m.Amount.LessThan(*other.Amount), nil
}

// Equal checks if other's amount is equal to m's amount
func (m *Money) Equal(other *Money) (bool, error) {
	err := m.sameKind(other)
	if err != nil {
		return false, err
	}

	return m.Amount.Equal(*other.Amount), nil
}

// LessThanOrEqual check if m's amount is less than or equal to other's amount
func (m *Money) LessThanOrEqual(other *Money) (bool, error) {
	less, err1 := m.LessThan(other)
	if err1 != nil {
		return false, err1
	}
	eq, err2 := m.Equal(other)
	if err2 != nil {
		return false, err2
	}
	return less || eq, nil
}

// Mul multiplty money with the givent other.
// other must be a float64, float32, int64, int
func (m *Money) Mul(other interface{}) (*Money, error) {
	var d decimal.Decimal

	switch t := other.(type) {
	case float64:
		floatDeci := decimal.NewFromFloat(t)
		d = m.Amount.Mul(floatDeci)
	case float32:
		floatDeci := decimal.NewFromFloat32(t)
		d = m.Amount.Mul(floatDeci)
	case int64:
		intDeci := decimal.NewFromInt(t)
		d = m.Amount.Mul(intDeci)
	case int:
		intDeci := decimal.NewFromInt32(int32(t))
		d = m.Amount.Mul(intDeci)

	default:
		return nil, ErrUnknownType
	}

	return NewMoney(&d, m.Currency)
}

// TrueDiv divides money with the given other.
// other must be a float64, float32, int64, int
func (m *Money) TrueDiv(other interface{}) (*Money, error) {
	var d decimal.Decimal

	switch t := other.(type) {
	case float64:
		floatDeci := decimal.NewFromFloat(t)
		d = m.Amount.Div(floatDeci)
	case float32:
		floatDeci := decimal.NewFromFloat32(t)
		d = m.Amount.Div(floatDeci)
	case int64:
		intDeci := decimal.NewFromInt(t)
		d = m.Amount.Div(intDeci)
	case int:
		intDeci := decimal.NewFromInt32(int32(t))
		d = m.Amount.Div(intDeci)

	default:
		return nil, ErrUnknownType
	}

	return NewMoney(&d, m.Currency)
}

// Add adds two money amount together, returns new money
func (m *Money) Add(other *Money) (*Money, error) {
	if err := m.sameKind(other); err != nil {
		return nil, err
	}
	amount := m.Amount.Add(*other.Amount)
	return &Money{&amount, m.Currency}, nil
}

// Sub subtracts currenct money to given `money`
func (m *Money) Sub(other *Money) (*Money, error) {
	if err := m.sameKind(other); err != nil {
		return nil, err
	}
	amount := m.Amount.Sub(*other.Amount)
	return &Money{&amount, m.Currency}, nil
}

func (m *Money) IsNotZero() bool {
	return !m.Amount.IsZero()
}

// func (m *Money) FlatTax(taxRate *decimal.Decimal, kepGross bool) {
// 	faction := decimal.NewFromInt(1).Add(*taxRate)
// 	if kepGross {
// 		// net :=
// 	}
// 	d := decimal.NewFromInt(12)
// }

// Return a copy of the object with its amount quantized.
// If `exp` is given the resulting exponent will match that of `exp`.
// Otherwise the resulting exponent will be set to the correct exponent
// of the currency if it's known and to default (two decimal places)
// otherwise.
func (m *Money) Quantize() (*Money, error) {
	places, err := GetCurrencyPrecision(m.Currency)
	if err != nil {
		return nil, err
	}
	d := m.Amount.Round(int32(places))
	return &Money{
		Amount:   &d,
		Currency: m.Currency,
	}, nil
}

// Apply a fixed discount to Money type.
func (m *Money) FixedDiscount(discount *Money) (*Money, error) {
	sub, err := m.Sub(discount) // same currencies check included
	if err != nil {
		return nil, err
	}

	if sub.Amount.GreaterThan(decimal.Zero) {
		return sub, nil
	}

	return &Money{
		Currency: m.Currency,
		Amount:   &decimal.Zero,
	}, nil
}
