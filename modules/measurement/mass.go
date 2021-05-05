package measurement

// weight units supported by app
const (
	G     = "g"
	LB    = "lb"
	OZ    = "oz"
	KG    = "kg"
	TONNE = "tonne"
)

var WEIGHT_UNIT_STRINGS = map[string]string{
	G:     "Gram",
	LB:    "Pound",
	OZ:    "Ounce",
	KG:    "kg",
	TONNE: "Tonne",
}

var WEIGHT_UNIT_CONVERSION = map[string]float32{
	KG:    1.0,
	G:     1000.0,
	OZ:    35.27396195,
	TONNE: 0.001,
	LB:    2.2046226218,
}

const STANDARD_WEIGHT_UNIT = KG
