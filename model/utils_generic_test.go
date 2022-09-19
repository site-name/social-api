package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddNoDup(t *testing.T) {
	data := AnyArray[string]{"one", "two", "three"}
	add := []string{"three", "four", "five"}

	res := data.AddNoDup(add...)

	if !reflect.DeepEqual(res, AnyArray[string]{"one", "two", "three", "four", "five"}) {
		t.Fatal("wrong implementation")
	}

	fmt.Println(res)
}
