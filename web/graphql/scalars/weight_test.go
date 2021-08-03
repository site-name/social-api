package scalars

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
)

func TestMarshalWeightScalar(t *testing.T) {
	w := measurement.Weight{
		Unit:   measurement.G,
		Amount: model.NewFloat32(2457),
	}
	msl := MarshalWeightScalar(&w)
	buf := bytes.Buffer{}
	msl.MarshalGQL(&buf)

	fmt.Println(buf.String())
}

func TestUnMarshalWeightScalar(t *testing.T) {
	// txt := map[string]interface{}{
	// 	"Unit":   "G",
	// 	"AMOUNT": 34.56,
	// }
	jsonTxt := `{"amount": 34.56, "unit": "KG"}`
	w, err := UnmarshalWeightScalar(jsonTxt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(w)
}
