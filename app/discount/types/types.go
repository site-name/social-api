package types

// DiscountCalculator number of `args` must be 1 or 2
//
//  // pass 1 argument if you want to calculate fixed discount
//  if len(args) == 1 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange)
//  }
//
//  // pass 2 arguments if you want to calculate percentage discount
//  if len(args) == 2 {
//		args[0].(type) == (*Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange) && args[1].(type) == bool
//  }
type DiscountCalculator func(args ...interface{}) (interface{}, error)
