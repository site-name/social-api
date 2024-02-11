package measurement

// volume units supported by system
const (
	CUBIC_MILLIMETER = "cubic_millimeter"
	CUBIC_CENTIMETER = "cubic_centimeter"
	CUBIC_DECIMETER  = "cubic_decimeter"
	CUBIC_METER      = "cubic_meter"
	LITER            = "liter"
	CUBIC_FOOT       = "cubic_foot"
	CUBIC_INCH       = "cubic_inch"
	CUBIC_YARD       = "cubic_yard"
	QT               = "qt"
	PINT             = "pint"
	FL_OZ            = "fl_oz"
	ACRE_IN          = "acre_in"
	ACRE_FT          = "acre_ft"
)

var VOLUME_UNIT_STRINGS = map[string]string{
	CUBIC_MILLIMETER: "Cubic millimeter",
	CUBIC_CENTIMETER: "Cubic centimeter",
	CUBIC_DECIMETER:  "Cubic decimeter",
	CUBIC_METER:      "Cubic meter",
	LITER:            "Liter",
	CUBIC_FOOT:       "Cubic foot",
	CUBIC_INCH:       "Cubic inch",
	CUBIC_YARD:       "Cubic yard",
	QT:               "Quart",
	PINT:             "Pint",
	FL_OZ:            "Fluid ounce",
	ACRE_IN:          "Acre inch",
	ACRE_FT:          "Acre feet",
}

var VOLUME_UNITS_CONVERSION = map[string]float64{
	CUBIC_MILLIMETER: 0.000000001,
	CUBIC_CENTIMETER: 0.000001,
	CUBIC_DECIMETER:  0.001,
	CUBIC_METER:      1.0,
	LITER:            0.001,
	CUBIC_FOOT:       0.0283168,
	CUBIC_INCH:       1.6387e-5,
	CUBIC_YARD:       0.76455486121558,
	QT:               1056.6882094,
	PINT:             2113.3764,
	FL_OZ:            2.8413e-5,
	ACRE_IN:          0.0097285583,
	ACRE_FT:          0.0008107132,
}

const STANDARD_VOLUME_UNIT = CUBIC_METER
