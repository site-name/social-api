package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	// "github.com/sitename/sitename/model/csv"
)

// func ExportProducts(exportFile *csv.ExportFile, scope map[string]interface{}, exportInfo map[string][]string, fileType string, delimiter string) {

// 	if delimiter == "" {
// 		delimiter = ";"
// 	}
// 	fileName := GetFileName("product", fileType)

// }

func GetFileName(modelName string, fileType string) string {
	utcNow := time.Now().UTC()
	year, month, day := utcNow.Date()
	hour := utcNow.Hour()
	min := utcNow.Minute()
	sec := utcNow.Second()

	timeStr := fmt.Sprintf("%d_%s_%d_%d_%d_%d", day, month.String(), year, hour, min, sec)

	return fmt.Sprintf(
		"%s_data_%s_%s.%s",
		modelName,
		timeStr,
		strings.ReplaceAll(model.NewId(), "-", ""),
		fileType,
	)
}

func GetProductQueryset(scope map[string]interface{}) {

}

func GetProductsInPatches() {

}
