package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
)

// ExportProducts is called by product export job, taks needed arguments then exports products
func (s *ServiceCsv) ExportProducts(exportProductsInput *gqlmodel.ExportProductsInput, delimeter string) *model.AppError {
	if delimeter == "" {
		delimeter = ";"
	}

	productFilterQuery := s.srv.Store.Product().AdvancedFilterQueryBuilder(exportProductsInput)

	exportFields, fileHeaders, dataHeaders, appErr := s.GetExportFieldsAndHeadersInfo(*exportProductsInput.ExportInfo)
	if appErr != nil {
		return appErr
	}

	getFileName("product", strings.ToLower(string(exportProductsInput.FileType)))

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
