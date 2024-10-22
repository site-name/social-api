// Code generated by "make app-layers"
// DO NOT EDIT

package sub_app_iface

import (
	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
)

// CsvService contains methods for working with csv
type CsvService interface {
	// CommonCreateExportEvent tells store to insert given export event into database then returns the inserted export event
	CommonCreateExportEvent(exportEvent model.ExportEvent) (*model.ExportEvent, *model_helper.AppError)
	// CreateExportFile inserts given export file into database then returns it
	CreateExportFile(file model.ExportFile) (*model.ExportFile, *model_helper.AppError)
	// ExportEventsByOption returns a list of export events filtered using given options
	ExportEventsByOption(options model.ExportEventFilterOption) ([]*model.ExportEvent, *model_helper.AppError)
	// ExportFileById returns an export file found by given id
	ExportFileById(id string) (*model.ExportFile, *model_helper.AppError)
	// ExportProducts is called by product export job, taks needed arguments then exports products
	ExportProducts(input *model.ExportProductsFilterOptions, delimeter string) *model_helper.AppError
	// Get export fields, all headers and headers mapping.
	// Based on export_info returns exported fields, fields to headers mapping and
	// all headers.
	// Headers contains product, variant, attribute and warehouse headers.
	GetExportFieldsAndHeadersInfo(exportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}) ([]string, []string, []string, *model_helper.AppError)
	// Get headers for exported attributes.
	// Headers are build from slug and contains information if it's a product or variant
	// attribute. Respectively for product: "slug-value (product attribute)"
	// and for variant: "slug-value (variant attribute)".
	GetAttributeHeaders(exportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}) ([]string, *model_helper.AppError)
	// Get headers for exported channels.
	//
	// Headers are build from slug and exported field.
	//
	// Example:
	// - currency code data header: "slug-value (channel currency code)"
	// - published data header: "slug-value (channel visible)"
	// - publication date data header: "slug-value (channel publication date)"
	GetChannelsHeaders(exportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}) ([]string, *model_helper.AppError)
	// Get headers for exported warehouses.
	// Headers are build from slug. Example: "slug-value (warehouse quantity)"
	GetWarehousesHeaders(exportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}) ([]string, *model_helper.AppError)
	// GetDefaultExportPayload returns a map for mapping
	GetDefaultExportPayload(exportFile model.ExportFile) (map[string]any, *model_helper.AppError)
	// GetProductsData Create data list of products and their variants with fields values.
	//
	// It return list with product and variant data which can be used as import to
	// csv writer and list of attribute and warehouse headers.
	//
	// TODO: consider improving me
	GetProductsData(products model.ProductSlice, exportFields, attributeIDs, warehouseIDs, channelIDs util.AnyArray[string]) []model_types.JSONString
	ExportProductsInBatches(productQuery squirrel.SelectBuilder, ExportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	}, exportFields []string, headers []string, delimiter string, fileType string) *model_helper.AppError
}
