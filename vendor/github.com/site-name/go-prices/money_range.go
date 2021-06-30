package goprices

import (
	"errors"
	"fmt"
)

var (
	// ErrStopLessThanStart is used when input stop money less than input start money
	ErrStopLessThanStart = errors.New("stop must be greater than start")
)

// MoneyRange has start and stop ends
type MoneyRange struct {
	Start    *Money
	Stop     *Money
	Currency string
}

// NewMoneyRange returns a new range. If start is greater than stop or start and stop have different
// currencies, return nil and non nil error
func NewMoneyRange(start, stop *Money) (*MoneyRange, error) {
	ok, err := stop.LessThan(start) // checking for same currency included
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, ErrStopLessThanStart
	}

	return &MoneyRange{
		Start:    start,
		Stop:     stop,
		Currency: start.Currency,
	}, nil
}

// String implements fmt.Stringer interface{}
func (m *MoneyRange) String() string {
	return fmt.Sprintf("Money{%q, %q}", m.Start.String(), m.Stop.String())
}

// Add adds a Money or MoneyRange to this MoneyRange
func (m *MoneyRange) Add(other interface{}) (*MoneyRange, error) {
	switch v := other.(type) {
	case *Money:
		start, err := m.Start.Add(v)
		if err != nil {
			return nil, err
		}
		stop, err := m.Stop.Add(v)
		if err != nil {
			return nil, err
		}
		return &MoneyRange{start, stop, m.Currency}, nil
	case *MoneyRange:
		start, err := m.Start.Add(v.Start)
		if err != nil {
			return nil, err
		}
		stop, err := m.Stop.Add(v.Stop)
		if err != nil {
			return nil, err
		}
		return &MoneyRange{start, stop, m.Currency}, nil
	default:
		return nil, ErrUnknownType
	}
}

// Sub subtracts current money to given `other`.
// `other` can be either `*Money` or `*MoneyRange`
func (m *MoneyRange) Sub(other interface{}) (*MoneyRange, error) {
	switch v := other.(type) {
	case *Money:
		start, err := m.Start.Sub(v)
		if err != nil {
			return nil, err
		}
		stop, err := m.Stop.Sub(v)
		if err != nil {
			return nil, err
		}
		return &MoneyRange{start, stop, m.Currency}, nil
	case *MoneyRange:
		start, err := m.Start.Sub(v.Start)
		if err != nil {
			return nil, err
		}
		stop, err := m.Stop.Sub(v.Stop)
		if err != nil {
			return nil, err
		}
		return &MoneyRange{start, stop, m.Currency}, nil
	default:
		return nil, ErrUnknownType
	}
}

// Equal Checks if two MoneyRange are equal both `Start`, `Stop` and `Currency`
func (m *MoneyRange) Equal(other *MoneyRange) bool {
	b1, err := m.Start.Equal(other.Start)
	if err != nil {
		return false
	}
	b2, err := m.Stop.Equal(other.Stop)
	if err != nil {
		return false
	}
	return b1 && b2
}

// Contains check if a Money is between this MoneyRange's two ends
func (m *MoneyRange) Contains(item *Money) bool {
	itemGreaterThanStart, err := m.Start.LessThanOrEqual(item)
	if err != nil {
		return false
	}
	itemLessThanStop, err := item.LessThanOrEqual(m.Stop)
	if err != nil {
		return false
	}
	return itemGreaterThanStart && itemLessThanStop
}

//Return a copy of the range with start and stop quantized.
// All arguments are passed to `Money.quantize
func (m *MoneyRange) Quantize() (*MoneyRange, error) {
	start, err := m.Start.Quantize()
	if err != nil {
		return nil, err
	}
	stop, err := m.Stop.Quantize()
	if err != nil {
		return nil, err
	}
	return &MoneyRange{
		Start:    start,
		Stop:     stop,
		Currency: m.Currency,
	}, nil
}

// Replace replace Start and Stop of currenct MoneyRagne With two given `start` and `stop` respectively.
func (m *MoneyRange) Replace(start, stop *Money) (*MoneyRange, error) {
	if start == nil {
		start = m.Start
	}
	if stop == nil {
		stop = m.Stop
	}
	return NewMoneyRange(start, stop)
}

// Apply a fixed discount to MoneyRange.
func (m *MoneyRange) FixedDiscount(discount *Money) (*MoneyRange, error) {
	baseStart, err := m.Start.FixedDiscount(discount)
	if err != nil {
		return nil, err
	}
	baseStop, err := m.Stop.FixedDiscount(discount)
	if err != nil {
		return nil, err
	}
	return NewMoneyRange(baseStart, baseStop)
}
