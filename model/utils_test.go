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
}

func TestModelFromJson(t *testing.T) {
	// var text = `{"id":"jshd849034bnkjhruieyr","Token":"jhd97847546565","create_at":1620271145022,"expires_at":0,"last_activity_at":0,"user_id":"","device_id":"","roles":"","is_oauth":false,"expired_notify":false,"props":null,"local":false}`
	var ses *Session
	ModelFromJson(&ses, strings.NewReader("sidu3874@"))

	fmt.Println(ses)
}
