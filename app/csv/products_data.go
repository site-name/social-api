package csv

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/attribute"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// GetProductsData Create data list of products and their variants with fields values.
//
// It return list with product and variant data which can be used as import to
// csv writer and list of attribute and warehouse headers.
func (a *ServiceCsv) GetProductsData(products product_and_discount.Products, exportFields []string, attributeIDs []string, warehouseIDs []string, channelIDs []string) {
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
func (s *ServiceCsv) getProductsRelationsData(products product_and_discount.Products, exportFields, attributeIDs, channelIDs []string) map[string]model.StringMap {
	var (
		manyToManyFields = ProductExportFields.HEADERS_TO_FIELDS_MAPPING["product_many_to_many"].Values()
		relationFields   = util.StringArrayIntersection(exportFields, manyToManyFields)
	)

	if len(relationFields) > 0 || len(attributeIDs) > 0 || len(channelIDs) > 0 {
		return s.prepareProductsRelationsData(products, relationFields, attributeIDs, channelIDs)
	}

	return map[string]model.StringMap{}
}

func (s *ServiceCsv) prepareProductsRelationsData(products product_and_discount.Products, fields, attributeIDs, channelIDs []string) map[string]model.StringMap {
	var (
		channelFields = ProductExportFields.PRODUCT_CHANNEL_LISTING_FIELDS.DeepCopy()
		resultData    = map[string]map[string][]string{}
	)

	fields = append(fields, "id")
	if len(attributeIDs) > 0 {
		fields = append(fields, ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS.Values()...)
	}
	if len(channelIDs) > 0 {
		fields = append(fields, channelFields.Values()...)
	}

	// var (
	// 	channelPkLookup   = channelFields.Pop("channel_pk")
	// 	channelSlugLookup = channelFields.Pop("slug")
	// )

	for _, data := range products.Flat() {
		var (
			pk         = data["id"].(string)
			collection = data["collections__slug"]
			image      = data.Pop("media__image", nil)
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

		s.handleAttributeData(pk, data, attributeIDs, resultData, ProductExportFields.PRODUCT_ATTRIBUTE_FIELDS, "product attribute")
	}

	panic("not implt")
}

func (s *ServiceCsv) getVariantsRelationsData(products product_and_discount.Products, exportFields, attributeIDs, warehouseIDs, channelIDs []string) {

}

type AttributeData struct {
	Slug       interface{}
	InputType  interface{}
	EntityType interface{}
	Unit       interface{}
	ValueSlug  interface{}
	ValueName  interface{}
	Value      interface{}
	FileUrl    interface{}
	RichText   interface{}
	Boolean    interface{}
	DateTime   interface{}
}

func (s *ServiceCsv) handleAttributeData(pk string, data model.StringInterface, attributeIDs []string, resultData map[string]map[string][]string, attributeFields model.StringMap, attributeOwner string) (map[string]map[string][]string, model.StringInterface) {
	attributePK := data.Pop(attributeFields["attribute_pk"], "").(string)

	attributeData := AttributeData{
		Slug:       data.Pop(attributeFields["slug"], nil),
		InputType:  data.Pop(attributeFields["input_type"], nil),
		FileUrl:    data.Pop(attributeFields["file_url"], nil),
		ValueSlug:  data.Pop(attributeFields["value_slug"], nil),
		ValueName:  data.Pop(attributeFields["value_name"], nil),
		Value:      data.Pop(attributeFields["value"], nil),
		EntityType: data.Pop(attributeFields["entity_type"], nil),
		Unit:       data.Pop(attributeFields["unit"], nil),
		RichText:   data.Pop(attributeFields["rich_text"], nil),
		Boolean:    data.Pop(attributeFields["boolean"], nil),
		DateTime:   data.Pop(attributeFields["date_time"], nil),
	}

	if len(attributeIDs) > 0 && util.StringInSlice(attributePK, attributeIDs) {
		resultData = s.addAttributeInfoToData(pk, attributeData, attributeOwner, resultData)
	}

	return resultData, data
}

func (s *ServiceCsv) addAttributeInfoToData(pk string, attributeData AttributeData, attributeOwner string, resultData map[string]map[string][]string) map[string]map[string][]string {
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
	if attributeData.InputType == nil {
		return ""
	}

	inputType, ok := attributeData.InputType.(string)
	if !ok {
		return ""
	}

	switch inputType {
	case attribute.FILE:
		if attributeData.FileUrl != nil {
			str, ok := attributeData.FileUrl.(string)
			if ok && str != "" {
				return filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, str)
			}

			return ""
		}
		return ""

	case attribute.REFERENCE:
		if attributeData.ValueSlug != nil && attributeData.EntityType != nil {
			return fmt.Sprintf("%v_%s", attributeData.EntityType, strings.Split(attributeData.ValueSlug.(string), "_")[1])
		}
		return ""

	case attribute.NUMERIC:
		value := fmt.Sprintf("%v", attributeData.ValueName)
		if attributeData.Unit != nil {
			value += fmt.Sprintf(" %v", attributeData.Unit)
		}
		return value

	case attribute.RICH_TEXT:
		slog.Warn("this case is not implemented yet")
		return ""

	case attribute.BOOLEAN:
		if attributeData.Boolean != nil {
			return strconv.FormatBool(attributeData.Boolean.(bool))
		}
		return "false"

	case attribute.DATE:
		return ""

	case attribute.DATE_TIME:
		return ""

	case attribute.SWATCH:
		if attributeData.FileUrl != nil {
			str, ok := attributeData.FileUrl.(string)
			if ok && str != "" {
				return filepath.Join(*s.srv.Config().ServiceSettings.SiteURL, str)
			}

			return attributeData.Value.(string)
		}
		return attributeData.Value.(string)

	default:
		if attributeData.ValueName != nil {
			str, ok := attributeData.ValueName.(string)
			if ok && str != "" {
				return str
			}
		}

		if attributeData.ValueSlug != nil {
			str, ok := attributeData.ValueSlug.(string)
			if ok && str != "" {
				return str
			}
		}

		return ""
	}
}
