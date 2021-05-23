package scalars

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestUnmarshalGQL(t *testing.T) {
	var d Decimal
	err := d.UnmarshalGQL(-12.45)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d)
}

func TestMarshalGQL(t *testing.T) {
	d := Decimal{
		Decimal: decimal.NewFromInt(12),
	}
	b := strings.Builder{}
	d.MarshalGQL(&b)
	if b.String() != "12" {
		t.Fatal("must be 12")
	}
	fmt.Println(b.String())
}

// func TestPositiveDecimalUnmarshalGQL(t *testing.T) {
// 	var d PositiveDecimal
// 	err := d.UnmarshalGQL("12.55")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(d)
// }

// func TestPositiveDecimalMarshalGQL(t *testing.T) {
// 	d := PositiveDecimal{
// 		Decimal: Decimal{
// 			Decimal: decimal.NewFromInt(-12),
// 		},
// 	}
// 	b := strings.Builder{}
// 	d.MarshalGQL(&b)
// 	if b.String() != "-12" {
// 		t.Fatal("must be 12")
// 	}
// 	fmt.Println(b.String())
// }

func TestWeightScalarUnMarshalQGL(t *testing.T) {
	var w WeightScalar
	err := w.UnmarshalGQL("23")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(w.Weight)
}
