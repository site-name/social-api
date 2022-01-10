package measurement

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertTo(t *testing.T) {
	w1 := Weight{
		Amount: 2000,
		Unit:   G,
	}
	res, err := w1.ConvertTo(KG)
	if err != nil {
		t.Fatal(err)
	}
	if res.Unit != KG {
		t.Fatal("res's unit must be 'kg'")
	}
	if res.Amount != float32(2) {
		t.Fatal("res's amount must be 2")
	}
}

func TestAdd(t *testing.T) {
	w1 := Weight{
		Amount: (2000.87345),
		Unit:   OZ,
	}
	w2 := Weight{
		Amount: (2000.3434),
		Unit:   LB,
	}
	addRes, err := w1.Add(&w2)
	require.NoError(t, err)
	if addRes.Unit != w1.Unit {
		t.Fatalf("res's unit must be %s\n", w1.Unit)
	}
	resAmount := w2.Amount/WEIGHT_UNIT_CONVERSION[w2.Unit]*WEIGHT_UNIT_CONVERSION[w1.Unit] + w1.Amount
	if addRes.Amount != resAmount {
		t.Fatal("res's amount is wrong")
	}

	fmt.Println(resAmount)
	fmt.Println(addRes)
}

func TestSub(t *testing.T) {
	w1 := Weight{
		Amount: (2000.87345),
		Unit:   OZ,
	}
	w2 := Weight{
		Amount: (2000.3434),
		Unit:   LB,
	}
	subRes, err := w1.Sub(w2)
	require.NoError(t, err)
	if subRes.Unit != w1.Unit {
		t.Fatalf("res's unit must be %s\n", w1.Unit)
	}
	resAmount := w1.Amount - w2.Amount/WEIGHT_UNIT_CONVERSION[w2.Unit]*WEIGHT_UNIT_CONVERSION[w1.Unit]
	if subRes.Amount != resAmount {
		t.Fatal("res's amount is wrong")
	}

	fmt.Println(resAmount)
	fmt.Println(subRes)
}

func TestMul(t *testing.T) {
	w1 := Weight{
		Amount: (2000.87345),
		Unit:   OZ,
	}
	var quan float32 = 5
	mulRes := w1.Mul(quan)
	if mulRes.Unit != w1.Unit {
		t.Fatalf("res's unit must be %s\n", w1.Unit)
	}
	resAmount := quan * w1.Amount
	if mulRes.Amount != resAmount {
		t.Fatalf("res's amount must be %f\n", resAmount)
	}

	fmt.Println(resAmount)
	fmt.Println(mulRes)
}

func TestMassToString(t *testing.T) {
	w1 := Weight{
		Amount: (3.456),
		Unit:   KG,
	}
	w2, _ := w1.ConvertTo(G)
	str := w2.String()
	fmt.Println(str)
}
