package csv_export

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/slog"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func (worker *Worker) GetExportFieldsAndHeadersInfo(exportInfo map[string][]string) ([]string, []string, []string) {
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
func (worker *Worker) GetAttributeHeaders(exportInfo map[string][]string) []string {
	attributeIDs := exportInfo["attributes"]
	if len(attributeIDs) == 0 {
		return []string{}
	}

	var (
		wg            sync.WaitGroup
		mutex         sync.Mutex
		err           error
		attributes_01 []*attribute.Attribute
		attributes_02 []*attribute.Attribute
	)

	filterOptions := [...]*attribute.AttributeFilterOption{
		{
			Distinct: true,
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: attributeIDs,
				},
			},
			ProductTypes: &model.StringFilter{
				StringOption: &model.StringOption{
					NULL: model.NewBool(false),
				},
			},
		},
		{
			Distinct: true,
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: attributeIDs,
				},
			},
			ProductVariantTypes: &model.StringFilter{
				StringOption: &model.StringOption{
					NULL: model.NewBool(false),
				},
			},
		},
	}

	wg.Add(len(filterOptions))

	for index, filterOption := range filterOptions {

		go func(idx int, option *attribute.AttributeFilterOption) {
			attributes, er := worker.app.Srv().Store.Attribute().FilterbyOption(filterOption)
			mutex.Lock()
			if er != nil {
				if err == nil {
					err = er
				}
			} else {
				if idx == 0 {
					attributes_01 = attributes
				} else {
					attributes_02 = attributes
				}
			}
			mutex.Unlock()

			wg.Done()
		}(index, filterOption)
	}

	wg.Wait()

	if err != nil {
		slog.Error(
			"worker failed to get attribute headers",
			slog.String("worker", worker.name),
			slog.String("error", err.Error()),
		)
		return nil
	}

	productHeaders := []string{}
	variantHeaders := []string{}

	for _, attr := range attributes_01 {
		productHeaders = append(productHeaders, attr.Slug+" (product attribute)")
	}
	for _, attr := range attributes_02 {
		variantHeaders = append(variantHeaders, attr.Slug+" (variant attribute)")
	}

	return append(productHeaders, variantHeaders...)
}

// Get headers for exported warehouses.
// Headers are build from slug. Example: "slug-value (warehouse quantity)"
func (worker *Worker) GetWarehousesHeaders(exportInfo map[string][]string) []string {
	warehouseIDs := exportInfo["warehouses"]
	if len(warehouseIDs) == 0 {
		return []string{}
	}

	warehouses, err := worker.app.Srv().Store.Warehouse().FilterByOprion(&warehouse.WarehouseFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: warehouseIDs,
			},
		},
	})
	if err != nil {
		slog.Error(
			"worker failed to get warehouse headers",
			slog.String("worker", worker.name),
			slog.String("error", err.Error()),
		)
		return nil
	}

	warehousesHeaders := []string{}
	for _, warehouse := range warehouses {
		warehousesHeaders = append(warehousesHeaders, warehouse.Slug+" (warehouse quantity)")
	}

	return warehousesHeaders
}

// Get headers for exported channels.
//
// Headers are build from slug and exported field.
//
// Example:
// - currency code data header: "slug-value (channel currency code)"
// - published data header: "slug-value (channel visible)"
// - publication date data header: "slug-value (channel publication date)"
func (worker *Worker) GetChannelsHeaders(exportInfo map[string][]string) []string {
	channelIDs := exportInfo["channels"]
	if len(channelIDs) == 0 {
		return []string{}
	}

	channels, err := worker.app.Srv().Store.Channel().FilterByOption(&channel.ChannelFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: channelIDs,
			},
		},
	})

	if err != nil {
		slog.Error(
			"worker failed to find channels",
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
	for _, channel := range channels {
		for _, field := range fields {
			if field != "slug" && field != "channel_pk" {
				channelsHeaders = append(channelsHeaders, fmt.Sprintf("%s (channel %s)", channel.Slug, strings.ReplaceAll(field, "_", " ")))
			}
		}
	}

	return channelsHeaders
}

// Get export fields from export info and prepare headers mapping.
// Based on given fields headers from export info, export fields set and
// headers mapping is prepared.
func GetProductExportFieldsAndHeaders(exportInfo map[string][]string) ([]string, []string) {
	var (
		exportFields = []string{"id"}
		fileHeaders  = []string{"id"}
	)

	fields := exportInfo["fields"]
	if len(fields) == 0 {
		return exportFields, fileHeaders
	}

	fieldsMapping := map[string]string{}
	for _, aMap := range ProductExportFields.HEADERS_TO_FIELDS_MAPPING {
		for key, value := range aMap {
			fieldsMapping[key] = value
		}
	}

	for _, field := range fields {
		exportFields = append(exportFields, fieldsMapping[field])
		fileHeaders = append(fileHeaders, field)
	}

	return exportFields, fileHeaders
}
