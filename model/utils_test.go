package model

import (
	"fmt"
	"strings"
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

func TestModelToJson(t *testing.T) {
	model := Session{
		Id:       "jshd849034bnkjhruieyr",
		CreateAt: GetMillis(),
	}
	res := ModelToJson(&model)
	fmt.Println(res)

	m := map[string]string{
		"one": "1",
		"two": "2",
	}
	res = ModelToJson(&m)
	fmt.Println(res)
}

type Person struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func TestModelFromJson(t *testing.T) {
	text := `{"name": "minh", "age": 23}`
	var per *Person

	ModelFromJson(&per, strings.NewReader(text))

	fmt.Println(per)
}
