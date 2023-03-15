package account

import (
	"testing"
)

func TestGetFont(t *testing.T) {
	fontName := "luximbi.ttf"

	_, err := getFont(fontName)
	if err != nil {
		t.Fatal(err)
	}
}
