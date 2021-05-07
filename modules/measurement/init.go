package measurement

import (
	"strings"
	"sync"
)

var (
	initOnce sync.Once

	MeasurementUnitMap     map[string]string
	MeasurementUnitChoices [][]string
)

func init() {
	initOnce.Do(func() {
		for k := range DISTANCE_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}
		for k := range AREA_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}
		for k := range WEIGHT_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}
		for k := range VOLUME_UNIT_STRINGS {
			MeasurementUnitMap[strings.ToUpper(k)] = k
		}

		for _, v := range MeasurementUnitMap {
			MeasurementUnitChoices = append(MeasurementUnitChoices, []string{v, v})
		}
	})
}
