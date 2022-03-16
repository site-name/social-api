package csv

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func (a *ServiceCsv) GetExportFieldsAndHeadersInfo(exportInfo gqlmodel.ExportInfoInput) ([]string, []string, []string, *model.AppError) {
	exportFields, fileHeaders := GetProductExportFieldsAndHeaders(exportInfo)
	attributeHeaders, appErr := a.GetAttributeHeaders(exportInfo)
	if appErr != nil {
		return nil, nil, nil, appErr
	}
	warehouseHeaders, appErr := a.GetWarehousesHeaders(exportInfo)
	if appErr != nil {
		return nil, nil, nil, appErr
	}
	channelHeaders, appErr := a.GetChannelsHeaders(exportInfo)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	dataHeaders := []string{}
	dataHeaders = append(dataHeaders, exportFields...)
	dataHeaders = append(dataHeaders, attributeHeaders...)
	dataHeaders = append(dataHeaders, warehouseHeaders...)
	dataHeaders = append(dataHeaders, channelHeaders...)

	fileHeaders = append(fileHeaders, attributeHeaders...)
	fileHeaders = append(fileHeaders, warehouseHeaders...)
	fileHeaders = append(fileHeaders, channelHeaders...)

	return exportFields, fileHeaders, dataHeaders, nil
}

// Get headers for exported attributes.
// Headers are build from slug and contains information if it's a product or variant
// attribute. Respectively for product: "slug-value (product attribute)"
// and for variant: "slug-value (variant attribute)".
func (a *ServiceCsv) GetAttributeHeaders(exportInfo gqlmodel.ExportInfoInput) ([]string, *model.AppError) {
	if len(exportInfo.Attributes) == 0 {
		return []string{}, nil
	}

	var (
		appError      *model.AppError
		attributes_01 []*attribute.Attribute
		attributes_02 []*attribute.Attribute

		wg  sync.WaitGroup
		mut sync.Mutex

		// syncSetAppError is used to safely set `appError`
		syncSetAppError = func(err *model.AppError) {
			mut.Lock()
			defer mut.Unlock()
			if err != nil && appError == nil {
				appError = err
			}
		}
	)

	filterOptions := [...]*attribute.AttributeFilterOption{
		{
			Distinct:     true,
			Id:           squirrel.Eq{store.AttributeTableName + ".Id": exportInfo.Attributes},
			ProductTypes: squirrel.NotEq{store.AttributeProductTableName + ".ProductTypeID": nil},
		},
		{
			Distinct:            true,
			Id:                  squirrel.Eq{store.AttributeTableName + ".Id": exportInfo.Attributes},
			ProductVariantTypes: squirrel.NotEq{store.AttributeVariantTableName + ".ProductTypeID": nil},
		},
	}

	wg.Add(len(filterOptions))

	for index, filterOption := range filterOptions {

		go func(idx int, option *attribute.AttributeFilterOption) {

			attributes, appErr := a.srv.AttributeService().AttributesByOption(option)
			if appErr != nil {
				syncSetAppError(appErr)
			} else {
				if idx == 0 {
					attributes_01 = attributes
				} else {
					attributes_02 = attributes
				}
			}

		}(index, filterOption)
	}

	wg.Wait()

	if appError != nil {
		return nil, appError
	}

	productHeaders := []string{}
	variantHeaders := []string{}

	for _, attr := range attributes_01 {
		productHeaders = append(productHeaders, attr.Slug+" (product attribute)")
	}
	for _, attr := range attributes_02 {
		variantHeaders = append(variantHeaders, attr.Slug+" (variant attribute)")
	}

	return append(productHeaders, variantHeaders...), nil
}

// Get headers for exported warehouses.
// Headers are build from slug. Example: "slug-value (warehouse quantity)"
func (a *ServiceCsv) GetWarehousesHeaders(exportInfo gqlmodel.ExportInfoInput) ([]string, *model.AppError) {
	if len(exportInfo.Warehouses) == 0 {
		return []string{}, nil
	}

	warehouses, appErr := a.srv.WarehouseService().WarehousesByOption(&warehouse.WarehouseFilterOption{
		Id: squirrel.Eq{store.WarehouseTableName + ".Id": exportInfo.Warehouses},
	})
	if appErr != nil {
		return nil, appErr
	}

	warehousesHeaders := []string{}
	for _, warehouse := range warehouses {
		warehousesHeaders = append(warehousesHeaders, warehouse.Slug+" (warehouse quantity)")
	}

	return warehousesHeaders, nil
}

// Get headers for exported channels.
//
// Headers are build from slug and exported field.
//
// Example:
// - currency code data header: "slug-value (channel currency code)"
// - published data header: "slug-value (channel visible)"
// - publication date data header: "slug-value (channel publication date)"
func (a *ServiceCsv) GetChannelsHeaders(exportInfo gqlmodel.ExportInfoInput) ([]string, *model.AppError) {
	if len(exportInfo.Channels) == 0 {
		return []string{}, nil
	}

	channels, appErr := a.srv.ChannelService().ChannelsByOption(&channel.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": exportInfo.Channels},
	})

	if appErr != nil {
		return nil, appErr
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

	return channelsHeaders, nil
}

// Get export fields from export info and prepare headers mapping.
// Based on given fields headers from export info, export fields set and
// headers mapping is prepared.
func GetProductExportFieldsAndHeaders(exportInfo gqlmodel.ExportInfoInput) ([]string, []string) {
	var (
		exportFields = []string{"id"}
		fileHeaders  = []string{"id"}
	)

	if len(exportInfo.Fields) == 0 {
		return exportFields, fileHeaders
	}

	fieldsMapping := map[string]string{}
	for _, aMap := range ProductExportFields.HEADERS_TO_FIELDS_MAPPING {
		for key, value := range aMap {
			fieldsMapping[key] = value
		}
	}

	for _, field := range exportInfo.Fields {
		actualField := strings.ToLower(string(field))

		exportFields = append(exportFields, fieldsMapping[actualField])
		fileHeaders = append(fileHeaders, actualField)
	}

	return exportFields, fileHeaders
}
