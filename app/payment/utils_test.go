package payment

import (
	"fmt"
	"log"
	"testing"
)

func Test_PriceFromMinorUnit(t *testing.T) {
	value := "12.345"
	currency := "VND"

	res, err := PriceFromMinorUnit(value, currency)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.String())
}
