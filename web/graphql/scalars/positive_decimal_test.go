package scalars

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/site-name/decimal"
)

func TestMarshalPositiveDecimal(t *testing.T) {
	d, _ := decimal.NewFromString("-3.14")
	msl := MarshalPositiveDecimal(&d)
	buf := bytes.Buffer{}
	msl.MarshalGQL(&buf)

	fmt.Println(buf.String())
}

func TestUnMarshalPositiveDecimal(t *testing.T) {
	str := "-3.14"
	d, err := UnmarshalPositiveDecimal(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d)
}
