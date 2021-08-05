package measurement

import (
	"strings"
	"sync"
)

var (
	initOnce sync.Once

	MeasurementUnitMap     map[string]string // MeasurementUnitMap contains all measurement unit notations supported by this aplication
	MeasurementUnitChoices [][]string        // MeasurementUnitChoices contains all measurements supported by this application
)

func init() {
	initOnce.Do(func() {
		MeasurementUnitMap = make(map[string]string)
		for k := range DISTANCE_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(string(k))] = string(k)
		}
		for k := range AREA_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}
		for k := range WEIGHT_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(string(k))] = string(k)
		}
		for k := range VOLUME_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}

		for _, v := range MeasurementUnitMap {
			MeasurementUnitChoices = append(MeasurementUnitChoices, []string{v, v})
		}
	})
}
