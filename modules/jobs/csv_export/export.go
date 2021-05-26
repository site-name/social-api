package csv_export

import (
	"strings"

	"github.com/sitename/sitename/web/model"
)

// Get export fields from export info and prepare headers mapping.
// Based on given fields headers from export info, export fields set and
// headers mapping is prepared.
func GetProductExportFieldsAndHeaders(exportInfo *model.ExportInfoInput) (exportFields []string, fileHeaders []string) {
	exportFields = []string{"id"}
	fileHeaders = []string{"id"}

	if len(exportInfo.Fields) == 0 {
		return
	}

	fieldsMapping := make(map[model.ProductFieldEnum]string)
	for _, value := range ProductExportFields.HEADERS_TO_FIELDS_MAPPING {
		for k, v := range value {
			fieldsMapping[k] = v
		}
	}

	for _, field := range exportInfo.Fields {
		lookupField := fieldsMapping[field]
		exportFields = append(exportFields, lookupField)
		fileHeaders = append(fileHeaders, strings.ReplaceAll(strings.ToLower(string(field)), "_", " "))
	}

	return
}

func (worker *Worker) exportProducts() {

}
