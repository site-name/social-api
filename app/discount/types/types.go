package types

// DiscountCalculator
// first must be *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange
type DiscountCalculator func(first any, fromGross *bool) (any, error)
