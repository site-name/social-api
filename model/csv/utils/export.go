package utils

import (
	"fmt"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
)

func ExportProducts(exportFile *csv.ExportFile, scope map[string]interface{}, exportInfo map[string][]string, fileType string, delimiter string) {
	// fileName := GetFileName("product", fileType)

}

func GetFileName(modelName string, fileType string) string {
	return fmt.Sprintf(
		"%s_data_%s_%s.%s",
		modelName,
		time.Now().Format(time.RFC3339),
		model.NewId(),
		fileType,
	)
}

// func GetProductQueryset(scope map[string]interface{})

func GetProductsInPatches() {

}
