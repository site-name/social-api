package csv_export

import (
	"fmt"
	"strings"

	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func (worker *Worker) GetExportFieldsAndHeadersInfo(exportInfo *gqlmodel.ExportInfoInput) ([]string, []string, []string) {
	exportFields, fileHeaders := GetProductExportFieldsAndHeaders(exportInfo)
	attributeHeaders := worker.GetAttributeHeaders(exportInfo)
	warehouseHeaders := worker.GetWarehousesHeaders(exportInfo)
	channelHeaders := worker.GetChannelsHeaders(exportInfo)

	dataHeaders := []string{}
	dataHeaders = append(dataHeaders, exportFields...)
	dataHeaders = append(dataHeaders, attributeHeaders...)
	dataHeaders = append(dataHeaders, warehouseHeaders...)
	dataHeaders = append(dataHeaders, channelHeaders...)

	fileHeaders = append(fileHeaders, attributeHeaders...)
	fileHeaders = append(fileHeaders, warehouseHeaders...)
	fileHeaders = append(fileHeaders, channelHeaders...)

	return exportFields, fileHeaders, dataHeaders
}

// Get headers for exported attributes.
// Headers are build from slug and contains information if it's a product or variant
// attribute. Respectively for product: "slug-value (product attribute)"
// and for variant: "slug-value (variant attribute)".
func (worker *Worker) GetAttributeHeaders(exportInfo *gqlmodel.ExportInfoInput) []string {
	if len(exportInfo.Attributes) == 0 {
		return nil
	}

	headers, err := worker.app.Srv().Store.Attribute().GetProductAndVariantHeaders(exportInfo.Attributes)
	if err != nil {
		slog.Error(
			"worker failed to get attribute headers",
			slog.String("worker", worker.name),
			slog.String("error", err.Error()),
		)
		return nil
	}

	return headers
}

// Get headers for exported warehouses.
// Headers are build from slug. Example: "slug-value (warehouse quantity)"
func (worker *Worker) GetWarehousesHeaders(exportInfo *gqlmodel.ExportInfoInput) []string {
	if len(exportInfo.Warehouses) == 0 {
		return nil
	}

	headers, err := worker.app.Srv().Store.Warehouse().GetWarehousesHeaders(exportInfo.Warehouses)
	if err != nil {
		slog.Error(
			"worker failed to get warehouse headers",
			slog.String("worker", worker.name),
			slog.String("error", err.Error()),
		)
		return nil
	}

	return headers
}

// Get headers for exported channels.
//
// Headers are build from slug and exported field.
//
// Example:
// - currency code data header: "slug-value (channel currency code)"
// - published data header: "slug-value (channel visible)"
// - publication date data header: "slug-value (channel publication date)"
func (worker *Worker) GetChannelsHeaders(exportInfo *gqlmodel.ExportInfoInput) []string {
	if len(exportInfo.Channels) == 0 {
		return nil
	}

	channels, err := worker.app.Srv().Store.Channel().GetChannelsByIdsAndOrder(exportInfo.Channels, "Slug")
	if err != nil {
		slog.Error(
			"worker failed to get channels header",
			slog.String("worker", worker.name),
			slog.String("error", err.Error()),
		)
		return nil
	}

	fields := []string{}
	for k := range ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS {
		fields = append(fields, k)
	}
	for k := range ProductExportFields.VARIANT_CHANNEL_LISTING_FIELDS {
		fields = append(fields, k)
	}

	channelsHeaders := []string{}
	for _, ch := range channels {
		list := []string{}
		for _, field := range fields {
			if field != "slug" && field != "channel_pk" {
				list = append(list, fmt.Sprintf("%s (channel %s)", ch.Slug, strings.ReplaceAll(field, "_", " ")))
			}
		}
		channelsHeaders = append(channelsHeaders, list...)
	}

	return channelsHeaders
}

// Get export fields from export info and prepare headers mapping.
// Based on given fields headers from export info, export fields set and
// headers mapping is prepared.
func GetProductExportFieldsAndHeaders(exportInfo *gqlmodel.ExportInfoInput) (exportFields []string, fileHeaders []string) {
	exportFields = []string{"id"}
	fileHeaders = []string{"id"}

	if len(exportInfo.Fields) == 0 {
		return
	}

	fieldsMapping := make(map[gqlmodel.ProductFieldEnum]string)
	for _, value := range ProductExportFields.HEADERS_TO_FIELDS_MAPPING {
		for k, v := range value {
			fieldsMapping[k] = v
		}
	}

	for _, field := range exportInfo.Fields {
		lookupField := fieldsMapping[field]
		exportFields = append(exportFields, lookupField)
		fileHeaders = append(fileHeaders, strings.ReplaceAll(strings.ToLower(string(field)), "_", " ")) // since fields are upper-cased words concatenated by underscores
	}

	return
}
