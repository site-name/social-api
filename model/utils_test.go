package model

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	rds := NewRandomString(20)
	fmt.Println(rds)

	if rds == "" {
		t.Fatal("Failed")
	}
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

// checkNowhereNil checks that the given interface value is not nil, and if a struct, that all of
// its public fields are also nowhere nil
func checkNowhereNil(t *testing.T, name string, value interface{}) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)
	switch v.Type().Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			t.Logf("%s was nil", name)
			return false
		}

		return checkNowhereNil(t, fmt.Sprintf("(*%s)", name), v.Elem().Interface())

	case reflect.Map:
		if v.IsNil() {
			t.Logf("%s was nil", name)
			return false
		}

		// Don't check map values
		return true

	case reflect.Struct:
		nowhereNil := true
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			// Ignore unexported fields
			if v.Type().Field(i).PkgPath != "" {
				continue
			}

			nowhereNil = nowhereNil && checkNowhereNil(t, fmt.Sprintf("%s.%s", name, v.Type().Field(i).Name), f.Interface())
		}

		return nowhereNil

	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.UnsafePointer:
		t.Logf("unhandled field %s, type: %s", name, v.Type().Kind())
		return false

	default:
		return true
	}
}

func TestDraftJSContentToRawText(t *testing.T) {
	data := StringInterface{
		"blocks": []StringInterface{
			{
				"data": StringMap{
					"text": "Hello World",
				},
				"type": "paragraph",
			},
			{
				"data": StringMap{
					"text": "Hello World",
				},
				"type": "paragraph",
			},
		},
	}

	res := DraftJSContentToRawText(data, "")
	fmt.Println(res)
}

func TestPaginationOptionsValidate(t *testing.T) {
	p := &PaginationOptions{}
	expr, appErr := p.ConstructSqlizer()
	if appErr != nil {
		t.Fatal(appErr)
	}

	_, _, err := expr.ToSql()
	if err != nil {
		t.Fatal(err)
	}
}
