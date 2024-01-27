package measurement

var (
	MeasurementUnitMap map[string]string // MeasurementUnitMap contains all measurement unit notations supported by this aplication
)

func init() {
	MeasurementUnitMap = make(map[string]string)
	for k := range DISTANCE_UNIT_STRINGS {
		MeasurementUnitMap[string(k)] = string(k)
	}
	for k := range AREA_UNIT_STRINGS {
		MeasurementUnitMap[k] = k
	}
	for k := range WEIGHT_UNIT_STRINGS {
		MeasurementUnitMap[string(k)] = string(k)
	}
	for k := range VOLUME_UNIT_STRINGS {
		MeasurementUnitMap[k] = k
	}

}
