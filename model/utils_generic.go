package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sitename/sitename/modules/util"
)

// AnyArray if a generic slice with a set of member methods that can be chained
type AnyArray[T util.Ordered] []T

// Remove removes input from the array
func (a AnyArray[T]) Remove(item T) AnyArray[T] {
	return util.RemoveItemsFromSlice(a, item)
}

// Map loops through current string slice and applies mapFunc to each index-item pair
//
// E.g
//
//	StringArray{"a", "b", "c"}.Map(func(_ int, s string) string { return s + s })
func (a AnyArray[T]) Map(fn func(index int, item T) T) AnyArray[T] {
	res := make([]T, len(a), cap(a))

	for idx, item := range a {
		res[idx] = fn(idx, item)
	}

	return res
}

// check if array of strings contains given input
func (sa AnyArray[T]) Contains(input T) bool {
	return util.ItemInSlice(input, sa)
}

// Equals checks if two arrays of strings have same length and contains the same elements at each index
func (sa AnyArray[T]) Equals(input []T) bool {
	return reflect.DeepEqual(sa, input)
}

// Join
func (sa AnyArray[T]) Join(sep string) string {
	var builder strings.Builder

	for i, item := range sa {
		builder.WriteString(fmt.Sprintf("%v", item))
		if i < len(sa) {
			builder.WriteString(sep)
		}
	}

	return builder.String()
}
