package csv

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func (a *ServiceCsv) GetExportFieldsAndHeadersInfo(
	exportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	},
) ([]string, []string, []string, *model.AppError) {
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
func (a *ServiceCsv) GetAttributeHeaders(exportInfo struct {
	Attributes []string
	Warehouses []string
	Channels   []string
	Fields     []string
}) ([]string, *model.AppError) {
	if len(exportInfo.Attributes) == 0 {
		return []string{}, nil
	}

	var (
		appError      *model.AppError
		attributes_01 []*model.Attribute
		attributes_02 []*model.Attribute

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

	wg.Add(2)

	go func() {
		defer wg.Done()

		attributes, appErr := a.srv.AttributeService().AttributesByOption(&model.AttributeFilterOption{
			Distinct:     true,
			Id:           squirrel.Eq{store.AttributeTableName + ".Id": exportInfo.Attributes},
			ProductTypes: squirrel.NotEq{store.AttributeProductTableName + ".ProductTypeID": nil},
		})
		if appErr != nil {
			syncSetAppError(appErr)
		} else {
			attributes_01 = attributes
		}
	}()

	go func() {
		defer wg.Done()

		attributes, appErr := a.srv.AttributeService().AttributesByOption(&model.AttributeFilterOption{
			Distinct:            true,
			Id:                  squirrel.Eq{store.AttributeTableName + ".Id": exportInfo.Attributes},
			ProductVariantTypes: squirrel.NotEq{store.AttributeVariantTableName + ".ProductTypeID": nil},
		})
		if appErr != nil {
			syncSetAppError(appErr)
		} else {
			attributes_02 = attributes
		}
	}()

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
func (a *ServiceCsv) GetWarehousesHeaders(exportInfo struct {
	Attributes []string
	Warehouses []string
	Channels   []string
	Fields     []string
}) ([]string, *model.AppError) {
	if len(exportInfo.Warehouses) == 0 {
		return []string{}, nil
	}

	warehouses, appErr := a.srv.WarehouseService().WarehousesByOption(&model.WarehouseFilterOption{
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
func (a *ServiceCsv) GetChannelsHeaders(exportInfo struct {
	Attributes []string
	Warehouses []string
	Channels   []string
	Fields     []string
}) ([]string, *model.AppError) {
	if len(exportInfo.Channels) == 0 {
		return []string{}, nil
	}

	channels, appErr := a.srv.ChannelService().ChannelsByOption(&model.ChannelFilterOption{
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
func GetProductExportFieldsAndHeaders(exportInfo struct {
	Attributes []string
	Warehouses []string
	Channels   []string
	Fields     []string
}) ([]string, []string) {
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
