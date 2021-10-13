package csv

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/warehouse"
)

// Get export fields, all headers and headers mapping.
// Based on export_info returns exported fields, fields to headers mapping and
// all headers.
// Headers contains product, variant, attribute and warehouse headers.
func (a *ServiceCsv) GetExportFieldsAndHeadersInfo(exportInfo map[string][]string) ([]string, []string, []string, *model.AppError) {
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
func (a *ServiceCsv) GetAttributeHeaders(exportInfo map[string][]string) ([]string, *model.AppError) {
	attributeIDs := exportInfo["attributes"]
	if len(attributeIDs) == 0 {
		return []string{}, nil
	}

	var (
		appError      *model.AppError
		attributes_01 []*attribute.Attribute
		attributes_02 []*attribute.Attribute

		// syncSetAppError is used to safely set `appError`
		syncSetAppError = func(err *model.AppError) {
			a.Lock()
			defer a.Unlock()
			if err != nil && appError == nil {
				appError = err
			}
			return
		}
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

	a.Add(len(filterOptions))

	for index, filterOption := range filterOptions {

		go func(idx int, option *attribute.AttributeFilterOption) {
			a.Lock()
			defer a.Unlock()
			defer a.Done()

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

			return

		}(index, filterOption)
	}

	a.Wait()

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
func (a *ServiceCsv) GetWarehousesHeaders(exportInfo map[string][]string) ([]string, *model.AppError) {
	warehouseIDs := exportInfo["warehouses"]
	if len(warehouseIDs) == 0 {
		return []string{}, nil
	}

	warehouses, appErr := a.srv.WarehouseService().WarehousesByOption(&warehouse.WarehouseFilterOption{
		Id: squirrel.Eq{a.srv.Store.Warehouse().TableName("Id"): warehouseIDs},
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
func (a *ServiceCsv) GetChannelsHeaders(exportInfo map[string][]string) ([]string, *model.AppError) {
	channelIDs := exportInfo["channels"]
	if len(channelIDs) == 0 {
		return []string{}, nil
	}

	channels, appErr := a.srv.ChannelService().ChannelsByOption(&channel.ChannelFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: channelIDs,
			},
		},
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
