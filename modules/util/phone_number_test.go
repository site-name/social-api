package util

import (
	"fmt"
	"testing"

	"github.com/ttacon/libphonenumber"
)

func TestIsValidPhoneNumber(t *testing.T) {
	num, err := libphonenumber.Parse("354575050", "")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(num)
}
