package model_types

import "testing"

func TestJsonMapGet(t *testing.T) {
	jm := JsonMap{
		"key1": "value1",
		"key2": "value2",
	}

	cases := []struct {
		key          string
		expected     any
		defaultValue any
	}{
		{"key1", "value1", ""},
		{"key2", "value2", ""},
		{"key3", "default", "default"},
		{"key3", "default", "default"},
	}

	for _, tc := range cases {
		var actual any
		if tc.defaultValue != "" {
			actual = jm.Get(tc.key, tc.defaultValue)
		} else {
			actual = jm.Get(tc.key)
		}
		if actual != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, actual)
		}
	}
}
