package payment

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/stretchr/testify/require"
)

func Test_PriceFromMinorUnit(t *testing.T) {
	value := "1000"
	currency := "USD"

	res, err := PriceFromMinorUnit(value, currency)
	require.NoError(t, err, err.Error())

	fmt.Println(res == nil)
}

func Test_PriceToMinorUnit(t *testing.T) {
	decimal := goprices.NewDecimal(decimal.NewFromFloat(12.34))
	currency := "USD"

	res, err := PriceToMinorUnit(decimal, currency)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}