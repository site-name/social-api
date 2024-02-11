package util

import (
	"fmt"
	"testing"
)

func TestValidatePhoneNumber(t *testing.T) {
	number, ok := ValidatePhoneNumber("0354575050", "US")
	if !ok {
		t.Fatal("invalid")
	}

	fmt.Println(number)
}
