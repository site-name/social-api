package measurement

type DistanceUnit string

var (
	DISTANCE_UNIT_STRINGS = map[DistanceUnit]string{
		CM:   "Centimeter",
		M:    "Meter",
		KM:   "Kilometers",
		FT:   "Feet",
		YD:   "Yard",
		INCH: "Inch",
	} // map distance unit alias to their fullnames
	DISTANCE_UNIT_CONVERSION = map[DistanceUnit]float32{
		CM:   1.0,
		M:    100.0,
		KM:   10000.0,
		FT:   30.48,
		YD:   91.44,
		INCH: 2.54,
	} // map distance unit to their according value
)

// Distance units supported by system
const (
	CM   DistanceUnit = "cm"
	M    DistanceUnit = "m"
	KM   DistanceUnit = "km"
	FT   DistanceUnit = "ft"
	YD   DistanceUnit = "yd"
	INCH DistanceUnit = "inch"
)

const STANDARD_DISTANCE_UNIT = CM
