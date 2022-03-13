package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/csv"
)

// ExportProducts is called by product export job, taks needed arguments then exports products
func (s *ServiceCsv) ExportProducts(exportFile *csv.ExportFile, exportProductsInput *gqlmodel.ExportProductsInput, delimeter string) *model.AppError {
	if delimeter == "" {
		delimeter = ";"
	}

	// parse export info

	// s.GetExportFieldsAndHeadersInfo()

	getFileName("product", strings.ToLower(string(exportProductsInput.FileType)))

	panic("not implemented")
}

// getFileName returns a file name for exported file
func getFileName(modelName string, fileType string) string {
	return fmt.Sprintf(
		"%s_data_%s_%s.%s",
		modelName,
		time.Now().UTC().Format("02_Jan_2006_15_04_05"),
		model.NewRandomString(16),
		fileType,
	)
}
