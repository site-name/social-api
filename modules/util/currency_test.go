package util

import (
	"fmt"
	"testing"
)

func TestGetCurrencyForCountry(t *testing.T) {
	c := GetCurrencyForCountry("ae")
	if c == "" {
		t.Fatal("cannot get currency code")
	} else {
		fmt.Println("currency code: ", c)
	}
}
