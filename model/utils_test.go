package model

import (
	"fmt"
	"testing"

	"github.com/nyaruka/phonenumbers"
)

func TestNewRandomString(t *testing.T) {
	rds := NewRandomString(20)
	fmt.Println(rds)

	if rds == "" {
		t.Fatal("Failed")
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	phone := "0354575050"
	country := ""

	num, err := phonenumbers.Parse(phone, country)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(num.String())
}
