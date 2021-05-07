package measurement

// area units supported by system
const (
	SQ_CM   = "sq_cm"
	SQ_M    = "sq_m"
	SQ_KM   = "sq_km"
	SQ_FT   = "sq_ft"
	SQ_YD   = "sq_yd"
	SQ_INCH = "sq_inch"
)

var AREA_UNIT_STRINGS = map[string]string{
	SQ_CM:   "Square centimeters",
	SQ_M:    "Square meters",
	SQ_KM:   "Square kilometers",
	SQ_FT:   "Square feet",
	SQ_YD:   "Square yards",
	SQ_INCH: "Square inches",
}

var AREA_UNIT_CONVERSION = map[string]float32{
	SQ_CM:   10000.0,
	SQ_M:    1.0,
	SQ_KM:   0.000001,
	SQ_FT:   10.7639104,
	SQ_YD:   1.19599005,
	SQ_INCH: 1550.0031,
}

const STANDARD_AREA_UNIT = SQ_M
