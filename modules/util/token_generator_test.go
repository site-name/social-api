package util

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestIntToBase36(t *testing.T) {
	s := intToBase36(1)
	fmt.Println(s)
}

func TestBase36ToInt(t *testing.T) {
	n, err := base36ToInt("1")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(n)
}

func TestSaltedHmac(t *testing.T) {
	h := saltedHmac("one", "one", "one", crypto.SHA1)
	res := h.Sum(nil)

	fmt.Println(hex.EncodeToString(res))
}

func TestMakeTokenWithTimestamp(t *testing.T) {
	p := NewTokenGenerator("", crypto.SHA1, "")
	token := p.makeTokenWithTimestamp(A{}, 12345, false)
	fmt.Println(token)
}

type A struct{}

func (a A) GetId() string {
	return "ThisIsID"
}
func (a A) GetPassword() string {
	return "sddddddddddddddddhhhhhhhh@"
}
func (a A) GetLastLogin() time.Time {
	tm, _ := time.Parse("2006-Jan-02", "2014-Feb-04")

	return tm
}
func (a A) GetEmail() string {
	return "leminhson2398@outlook.com"
}

func TestMakeToken(t *testing.T) {
	p := NewTokenGenerator("", crypto.SHA1, "")
	token := p.MakeToken(A{})
	fmt.Println(token)
	// 665580546
}

func TestCheckToken(t *testing.T) {
	p := NewTokenGenerator("", crypto.SHA1, "")
	v := p.CheckToken(A{}, "b09rob-cb1758d1a1755280e7f5")
	fmt.Println(v)
}
