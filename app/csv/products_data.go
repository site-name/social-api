package csv

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// GetProductsData Create data list of products and their variants with fields values.
//
// It return list with product and variant data which can be used as import to
// csv writer and list of attribute and warehouse headers.
func (a *ServiceCsv) GetProductsData(query squirrel.SelectBuilder, exportFields []string, attributeIDs []string, warehouseIDs []string, channelIDs []string) {

	var (
		exportVariantID     = util.StringInSlice("variants__id", exportFields)
		productFields       = ProductExportFields.HEADERS_TO_FIELDS_MAPPING["fields"].Values()
		productExportFields = util.StringArrayIntersection(exportFields, productFields)
	)

	if !exportVariantID {
		productExportFields = append(productExportFields, "variants__id")
	}
}

// Get data about product relations fields.
// If any many to many fields are in export_fields or some attribute_ids exists then
// dict with product relations fields is returned.
// Otherwise it returns empty dict.
func (s *ServiceCsv) getProductsRelationsData(query squirrel.SelectBuilder, exportFields, attributeIDs, channelIDs []string) map[string]model.StringMap {
	var (
		manyToManyFields = ProductExportFields.HEADERS_TO_FIELDS_MAPPING["product_many_to_many"].Values()
		relationFields   = util.StringArrayIntersection(exportFields, manyToManyFields)
	)

	if len(relationFields) > 0 || len(attributeIDs) > 0 || len(channelIDs) > 0 {
		return s.prepareProductsRelationsData(query, relationFields, attributeIDs, channelIDs)
	}

	return map[string]model.StringMap{}
}

func (s *ServiceCsv) prepareProductsRelationsData(query squirrel.SelectBuilder, fields, attributeIDs, channelIDs []string) map[string]model.StringMap {
	var (
		channelFields = ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS.DeepCopy()
		resultData    = map[string]model.StringMap{}
	)

	fields = append(fields, "id")
	if len(ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS) > 0 {
		fields = append(fields, ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS.Values()...)
	}
	if len(ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS) > 0 {
		fields = append(fields, ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS.Values()...)
	}

	var (
		channelPkLookup   = channelFields["channel_pk"]
		channelSlugLookup = channelFields["slug"]
	)
	delete(channelFields, "channel_pk")
	delete(channelFields, "slug")

}

func (s *ServiceCsv) getVariantsRelationsData(query squirrel.SelectBuilder, exportFields, attributeIDs, warehouseIDs, channelIDs []string) {

}
