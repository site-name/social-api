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

func TestMap(t *testing.T) {
	data := AnyArray[int]{1, 2, 3, 4, 5}
	newArr := data.Map(func(_ int, item int) int {
		return item + 1
	})

	if !reflect.DeepEqual(newArr, AnyArray[int]{2, 3, 4, 5, 6}) {
		t.Fatal("wrong implementation")
	}

	fmt.Println(newArr)
}
