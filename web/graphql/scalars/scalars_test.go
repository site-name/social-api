package scalars

import (
	"fmt"
	"testing"
)

func TestWeightScalarUnMarshalQGL(t *testing.T) {
	var w WeightScalar
	err := w.UnmarshalGQL("23")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(w.Weight)
}
