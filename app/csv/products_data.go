package csv

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// GetProductsData Create data list of products and their variants with fields values.
//
// It return list with product and variant data which can be used as import to
// csv writer and list of attribute and warehouse headers.
//
// TODO: consider improving me
func (a *ServiceCsv) GetProductsData(products model.ProductSlice, exportFields, attributeIDs, warehouseIDs, channelIDs util.AnyArray[string]) []model_types.JSONString {
	var (
		exportVariantID     = exportFields.Contains("variants__id")
		productFields       = lo.Values(ProductExportFields.HEADERS_TO_FIELDS_MAPPING["fields"])
		productExportFields = exportFields.InterSection(productFields)
	)

	if !exportVariantID {
		productExportFields = append(productExportFields, "variants__id")
	}

	productsRelationsData := a.getProductsRelationsData(products, exportFields, attributeIDs, channelIDs)
	variantsRelationsData := a.getVariantsRelationsData(products, exportFields, attributeIDs, warehouseIDs, channelIDs)

	res := []model_types.JSONString{}

	for _, productData := range products.Flat() {
		var pk = productData["id"].(string)
		var variantPK string

		if exportVariantID {
			variantPK = productData.Get("variants__id", "").(string)
		} else {
			variantPK = productData.Pop("variants__id", "").(string)
		}

		productRelationsData := productsRelationsData[pk]
		if productRelationsData == nil {
			productRelationsData = model_helper.StringMap{}
		}

		variantRelationsData := variantsRelationsData[variantPK]
		if variantRelationsData == nil {
			variantRelationsData = model_helper.StringMap{}
		}

		if exportVariantID {
			productData["variants__id"] = variantPK
		}

		data := model_types.JSONString{}
		data.Merge(productData)
		for k, v := range productRelationsData {
			data[k] = v
		}
		for k, v := range variantRelationsData {
			data[k] = v
		}

		res = append(res, data)
	}

	return res
}

// Get data about product relations fields.
// If any many to many fields are in export_fields or some attribute_ids exists then
// dict with product relations fields is returned.
// Otherwise it returns empty di`ct.
func (s *ServiceCsv) getProductsRelationsData(products model.ProductSlice, exportFields, attributeIDs, channelIDs util.AnyArray[string]) map[string]model_helper.StringMap {
	var (
		manyToManyFields = lo.Values(ProductExportFields.HEADERS_TO_FIELDS_MAPPING["product_many_to_many"])
		relationFields   = exportFields.InterSection(manyToManyFields)
	)

	if len(relationFields) > 0 || len(attributeIDs) > 0 || len(channelIDs) > 0 {
		return s.prepareProductsRelationsData(products, relationFields, attributeIDs, channelIDs)
	}

	return map[string]model_helper.StringMap{}
}

func (s *ServiceCsv) prepareProductsRelationsData(products model.ProductSlice, fields util.AnyArray[string], attributeIDs, channelIDs []string) map[string]model_helper.StringMap {
	var (
		channelFields = ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS.DeepCopy()
		resultData    = map[string]map[string][]any{}
	)

	// fields = append(fields, "id")
	// if len(attributeIDs) > 0 {
	// 	fields = fields.AddNoDup(ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS.Values()...)
	// }
	// if len(channelIDs) > 0 {
	// 	fields = fields.AddNoDup(channelFields.Values()...)
	// }

	var (
		channelPkLookup   = channelFields.Pop("channel_pk")
		channelSlugLookup = channelFields.Pop("slug")
	)

	for _, data := range products.Flat() {
		var (
			pk         = data["id"].(string)
			collection = data["collections__slug"]
			image      = data.Pop("media__image", "")
		)

		if image != nil {
			resultData[pk]["media__image"] = append(
				resultData[pk]["media__image"],
				filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, image.(string)),
			)
		}

		if collection != nil {
			resultData[pk]["collections__slug"] = append(
				resultData[pk]["collections__slug"],
				collection.(string),
			)
		}

		resultData, data = s.handleAttributeData(pk, data, attributeIDs, resultData, ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS, "product attribute")
		resultData, data = s.handleChannelData(pk, data, channelIDs, resultData, channelPkLookup, channelSlugLookup, channelFields)
	}

	result := map[string]model_helper.StringMap{}
	for pk, data := range resultData {
		result[pk] = model_helper.StringMap{}

		for header, values := range data {
			var str string

			for idx, value := range values {
				if idx < len(values)-1 {
					str += fmt.Sprintf("%v, ", value)
					continue
				}
				str += fmt.Sprintf("%v", value)
			}

			result[pk][header] = str
		}
	}

	return result
}

func (s *ServiceCsv) getVariantsRelationsData(products model.ProductSlice, exportFields, attributeIDs, warehouseIDs, channelIDs util.AnyArray[string]) map[string]model_helper.StringMap {
	manyToManyFields := lo.Values(ProductExportFields.HEADERS_TO_FIELDS_MAPPING["variant_many_to_many"])
	relationsFields := exportFields.InterSection(manyToManyFields)

	if len(relationsFields) > 0 || len(attributeIDs) > 0 || len(channelIDs) > 0 {
		return s.prepareVariantsRelationsData(products, relationsFields, attributeIDs, warehouseIDs, channelIDs)
	}

	return map[string]model_helper.StringMap{}
}

func (s *ServiceCsv) prepareVariantsRelationsData(products model.ProductSlice, fields util.AnyArray[string], attributeIDs, warehouseIDs, channelIDs []string) map[string]model_helper.StringMap {
	var channelFields = ProductExportFields.VARIANT_CHANNEL_LISTING_FIELDS.DeepCopy()

	// fields = append(fields, "variants__id")
	// if len(attributeIDs) > 0 {
	// 	fields = fields.AddNoDup(ProductExportFields.VARIANT_ATTRIBUTE_FIELDS.Values()...)
	// }
	// if len(warehouseIDs) > 0 {
	// 	fields = fields.AddNoDup(ProductExportFields.WAREHOUSE_FIELDS.Values()...)
	// }
	// if len(channelIDs) > 0 {
	// 	fields = fields.AddNoDup(channelFields.Values()...)
	// }

	var (
		resultData        = map[string]map[string][]any{}
		channelPKLookup   = channelFields.Pop("channel_pk")
		channelSlugLookup = channelFields.Pop("slug")
	)

	for _, data := range products.Flat() {
		pk := data.Get("variants__id").(string)
		image := data.Pop("variants__media__image", nil)

		if image != nil {
			resultData[pk]["variants__media__image"] = append(resultData[pk]["variants__media__image"], filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, image.(string)))
		}

		resultData, data = s.handleAttributeData(pk, data, attributeIDs, resultData, ProductExportFields.VARIANT_ATTRIBUTE_FIELDS, "variant attribute")
		resultData, data = s.handleChannelData(pk, data, channelIDs, resultData, channelPKLookup, channelSlugLookup, channelFields)
		resultData, data = s.handleWarehouseData(pk, data, warehouseIDs, resultData, ProductExportFields.WAREHOUSE_FIELDS)
	}

	result := map[string]model_helper.StringMap{}
	for pk, data := range resultData {
		result[pk] = model_helper.StringMap{}

		for header, values := range data {
			var str string

			for idx, item := range values {
				if idx < len(values)-1 {
					str += fmt.Sprintf("%v, ", item)
					continue
				}
				str += fmt.Sprintf("%v", item)
			}

			result[pk][header] = str
		}
	}

	return result
}

func (s *ServiceCsv) handleWarehouseData(pk string, data model_types.JSONString, warehouseIDs util.AnyArray[string], resultData map[string]map[string][]any, warehouseFields model_helper.StringMap) (map[string]map[string][]any, model_types.JSONString) {
	warehousePK := data.Pop(warehouseFields["warehouse_pk"], "").(string)
	warehouseData := model_types.JSONString{
		"slug": data.Pop(warehouseFields["slug"], nil),
		"qty":  data.Pop(warehouseFields["quantity"], nil),
	}

	if warehouseIDs.Contains(warehousePK) {
		resultData = s.addWarehouseInfoToData(pk, warehouseData, resultData)
	}

	return resultData, data
}

func (s *ServiceCsv) addWarehouseInfoToData(pk string, warehouseData model_types.JSONString, resultData map[string]map[string][]any) map[string]map[string][]any {
	slug, ok := warehouseData["slug"]
	if ok && slug != nil {
		warehouseQtyHeader := fmt.Sprintf("%v (warehouse quantity)", slug)
		if _, ok := resultData[pk][warehouseQtyHeader]; !ok {
			resultData[pk][warehouseQtyHeader] = []any{warehouseData["qty"]}
		}
	}

	return resultData
}

type AttributeData struct {
	Slug       any
	InputType  any
	EntityType any
	Unit       any
	ValueSlug  any
	ValueName  any
	Value      any
	FileUrl    any
	RichText   any
	Boolean    any
	DateTime   any
}

func (s *ServiceCsv) handleAttributeData(pk string, data model_types.JSONString, attributeIDs util.AnyArray[string], resultData map[string]map[string][]any, attributeFields model_helper.StringMap, attributeOwner string) (map[string]map[string][]any, model_types.JSONString) {
	attributePK := data.Pop(attributeFields["attribute_pk"], "").(string)

	attributeData := AttributeData{
		Slug:       data.Pop(attributeFields["slug"], nil),
		InputType:  data.Pop(attributeFields["input_type"], nil),
		FileUrl:    data.Pop(attributeFields["file_url"], nil),
		ValueSlug:  data.Pop(attributeFields["value_slug"], ""),
		ValueName:  data.Pop(attributeFields["value_name"], nil),
		Value:      data.Pop(attributeFields["value"], nil),
		EntityType: data.Pop(attributeFields["entity_type"], nil),
		Unit:       data.Pop(attributeFields["unit"], nil),
		RichText:   data.Pop(attributeFields["rich_text"], nil),
		Boolean:    data.Pop(attributeFields["boolean"], nil),
		DateTime:   data.Pop(attributeFields["date_time"], nil),
	}

	if attributeIDs.Contains(attributePK) {
		resultData = s.addAttributeInfoToData(pk, attributeData, attributeOwner, resultData)
	}

	return resultData, data
}

func (s *ServiceCsv) handleChannelData(pk string, data model_types.JSONString, channelIDs util.AnyArray[string], resultData map[string]map[string][]any, pkLookup, slugLookup string, fields model_helper.StringMap) (map[string]map[string][]any, model_types.JSONString) {
	channelPK := data.Pop(pkLookup, "").(string)
	channelData := model_types.JSONString{
		"slug": data.Pop(slugLookup, nil),
	}

	for field, lookup := range fields {
		channelData[field] = data.Pop(lookup, nil)
	}

	if channelIDs.Contains(channelPK) {
		resultData = s.addChannelInfoToData(pk, channelData, resultData, lo.Keys(fields))
	}

	return resultData, data
}

func (s *ServiceCsv) addAttributeInfoToData(pk string, attributeData AttributeData, attributeOwner string, resultData map[string]map[string][]any) map[string]map[string][]any {
	if attributeData.Slug == nil {
		return resultData
	}

	var (
		header = fmt.Sprintf("%v (%s)", attributeData.Slug, attributeOwner)
		value  = s.prepareAttributeValue(attributeData)
	)

	resultData[pk][header] = append(resultData[pk][header], value)

	return resultData
}

func (s *ServiceCsv) prepareAttributeValue(attributeData AttributeData) string {

	t, ok := attributeData.InputType.(model.AttributeInputType)
	if !ok {
		return ""
	}

	switch t {
	case model.AttributeInputTypeFile:
		if url := attributeData.FileUrl; url != nil {
			if strURL, ok := url.(string); ok {
				return filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, strURL)
			}
			return ""
		}
		return ""

	case model.AttributeInputTypeReference:
		if slug := attributeData.ValueSlug; slug != nil {
			if strSlug, ok := slug.(string); ok {
				return fmt.Sprintf("%v_%s", attributeData.EntityType, strings.Split(strSlug, "_")[1])
			}
			return ""
		}
		return ""

	case model.AttributeInputTypeNumeric:
		value := fmt.Sprintf("%v", attributeData.ValueName)
		if attributeData.Unit != nil {
			value += fmt.Sprintf(" %v", attributeData.Unit)
		}
		return value

	case model.AttributeInputTypeRichText:
		slog.Warn("this case is not implemented yet")
		return ""

	case model.AttributeInputTypeBoolean:
		if attributeData.Boolean != nil {
			if b, ok := attributeData.Boolean.(bool); ok {
				return strconv.FormatBool(b)
			}
			return ""
		}
		return ""

	case model.AttributeInputTypeDate:
		if tim, ok := attributeData.DateTime.(time.Time); ok {
			return tim.Format("2006-01-02")
		}
		return ""

	case model.AttributeInputTypeDateTime:
		if tim, ok := attributeData.DateTime.(time.Time); ok {
			return tim.Format("2006-01-02 15:04:05")
		}
		return ""

	case model.AttributeInputTypeSwatch:
		if attributeData.FileUrl != nil {
			if strURL, ok := attributeData.FileUrl.(string); ok {
				return filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, strURL)
			}
			return ""
		}
		if attributeData.Value != nil {
			return fmt.Sprintf("%v", attributeData.Value)
		}

	default:
		if attributeData.ValueName != nil {
			return fmt.Sprintf("%v", attributeData.ValueName)
		} else if attributeData.ValueSlug != nil {
			return fmt.Sprintf("%v", attributeData.ValueSlug)
		}
		return ""
	}

	return ""
}

func (s *ServiceCsv) addChannelInfoToData(pk string, channelData model_types.JSONString, resultData map[string]map[string][]any, fields []string) map[string]map[string][]any {
	slug, ok := channelData["slug"]
	if ok && slug != nil {
		for _, field := range fields {
			header := fmt.Sprintf("%v (channel %s)", slug, strings.ReplaceAll(field, "_", " "))

			if _, ok := resultData[header]; !ok {
				resultData[pk][header] = append(resultData[pk][header], channelData[field])
			}
		}
	}

	return resultData
}
