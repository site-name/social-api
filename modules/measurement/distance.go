package measurement

// Distance units supported by system
const (
	CM   = "cm"
	M    = "m"
	KM   = "km"
	FT   = "ft"
	YD   = "yd"
	INCH = "inch"
)

var DISTANCE_UNIT_STRINGS = map[string]string{
	CM:   "Centimeter",
	M:    "Meter",
	KM:   "Kilometers",
	FT:   "Feet",
	YD:   "Yard",
	INCH: "Inch",
}

var DISTANCE_UNIT_CONVERSION = map[string]float32{
	CM:   1.0,
	M:    100.0,
	KM:   10000.0,
	FT:   30.48,
	YD:   91.44,
	INCH: 2.54,
}

const STANDARD_DISTANCE_UNIT = CM
