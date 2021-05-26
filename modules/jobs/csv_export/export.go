package csv_export

import (
	"strings"

	"github.com/sitename/sitename/web/model"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func GetExportFieldsAndHeadersInfo(exportInfo *model.ExportInfoInput) {

}

// Get headers for exported attributes.
// Headers are build from slug and contains information if it's a product or variant
// attribute. Respectively for product: "slug-value (product attribute)"
// and for variant: "slug-value (variant attribute)".
func (worker *Worker) GetAttributeHeaders(exportInfo *model.ExportInfoInput) []string {
	if len(exportInfo.Attributes) == 0 {
		return []string{}
	}

	// query := worker.app.Srv().Store
}

// Get headers for exported warehouses.
// Headers are build from slug. Example: "slug-value (warehouse quantity)"
func GetWarehousesHeaders(exportInfo *model.ExportInfoInput) []string {

}

// Get headers for exported channels.
//
// Headers are build from slug and exported field.
//
// Example:
// - currency code data header: "slug-value (channel currency code)"
// - published data header: "slug-value (channel visible)"
// - publication date data header: "slug-value (channel publication date)"
func GetChannelsHeaders(exportInfo *model.ExportInfoInput) []string {

}

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

func (worker *Worker) ExportProducts() {

}
