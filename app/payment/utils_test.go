package payment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PriceFromMinorUnit(t *testing.T) {
	value := "1000"
	currency := "USD"

	res, err := PriceFromMinorUnit(value, currency)
	require.NoError(t, err, err.Error())

	fmt.Println(res == nil)
}
