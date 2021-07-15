package account

import (
	"fmt"
	"testing"
)

func TestGetFont(t *testing.T) {
	fontName := "luximbi.ttf"

	font, err := getFont(fontName)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(font)
}
