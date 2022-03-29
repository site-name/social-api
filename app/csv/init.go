package csv

import (
	"github.com/sitename/sitename/model"
)

type Fields struct {
	HEADERS_TO_FIELDS_MAPPING      map[string]model.StringMap
	PRODUCT_ATTRIBUTE_FIELDS       model.StringMap
	PRODUCT_CHANNEL_LISTING_FIELDS model.StringMap
	WAREHOUSE_FIELDS               model.StringMap
	VARIANT_ATTRIBUTE_FIELDS       model.StringMap
	VARIANT_CHANNEL_LISTING_FIELDS model.StringMap
}

var ProductExportFields = &Fields{
	HEADERS_TO_FIELDS_MAPPING: map[string]model.StringMap{
		"fields": {
			"id":                                "id",
			"name":                              "name",
			"description":                       "description_as_str",
			"category":                          "category__slug",
			"product_type":                      "product_type__name",
			"charge_taxes":                      "charge_taxes",
			"product_weight":                    "product_weight",
			"variant_id":                        "variants__id",
			"variant_sku":                       "variants__sku",
			"variant_weight":                    "variant_weight",
			"variant_is_preorder":               "variants__is_preorder",
			"variant_preorder_global_threshold": "variants__preorder_global_threshold",
			"variant_preorder_end_date":         "variants__preorder_end_date",
		},
		"product_many_to_many": {
			"collections":   "collections__slug",
			"product_media": "media__image",
		},
		"variant_many_to_many": {
			"variant_media": "variants__media__image",
		},
	},
	PRODUCT_ATTRIBUTE_FIELDS: model.StringMap{
		"value_slug":   "attributes__values__slug",
		"value_name":   "attributes__values__name",
		"file_url":     "attributes__values__file_url",
		"rich_text":    "attributes__values__rich_text",
		"value":        "attributes__values__value",
		"boolean":      "attributes__values__boolean",
		"date_time":    "attributes__values__date_time",
		"slug":         "attributes__assignment__attribute__slug",
		"input_type":   "attributes__assignment__attribute__input_type",
		"entity_type":  "attributes__assignment__attribute__entity_type",
		"unit":         "attributes__assignment__attribute__unit",
		"attribute_pk": "attributes__assignment__attribute__pk",
	},
	PRODUCT_CHANNEL_LISTING_FIELDS: model.StringMap{
		"channel_pk":             "channel_listings__channel__pk",
		"slug":                   "channel_listings__channel__slug",
		"product_currency_code":  "channel_listings__currency",
		"published":              "channel_listings__is_published",
		"publication_date":       "channel_listings__publication_date",
		"searchable":             "channel_listings__visible_in_listings",
		"available_for_purchase": "channel_listings__available_for_purchase",
	},
	WAREHOUSE_FIELDS: model.StringMap{
		"slug":         "variants__stocks__warehouse__slug",
		"quantity":     "variants__stocks__quantity",
		"warehouse_pk": "variants__stocks__warehouse__id",
	},
	VARIANT_ATTRIBUTE_FIELDS: model.StringMap{
		"value_slug":   "variants__attributes__values__slug",
		"value_name":   "variants__attributes__values__name",
		"file_url":     "variants__attributes__values__file_url",
		"rich_text":    "variants__attributes__values__rich_text",
		"value":        "variants__attributes__values__value",
		"boolean":      "variants__attributes__values__boolean",
		"date_time":    "variants__attributes__values__date_time",
		"slug":         "variants__attributes__assignment__attribute__slug",
		"input_type":   "variants__attributes__assignment__attribute__input_type",
		"entity_type":  "variants__attributes__assignment__attribute__entity_type",
		"unit":         "variants__attributes__assignment__attribute__unit",
		"attribute_pk": "variants__attributes__assignment__attribute__pk",
	},
	VARIANT_CHANNEL_LISTING_FIELDS: model.StringMap{
		"channel_pk":                          "variants__channel_listings__channel__pk",
		"slug":                                "variants__channel_listings__channel__slug",
		"price_amount":                        "variants__channel_listings__price_amount",
		"variant_currency_code":               "variants__channel_listings__currency",
		"variant_cost_price":                  "variants__channel_listings__cost_price_amount",
		"variant_preorder_quantity_threshold": "variants__channel_listings__preorder_quantity_threshold",
	},
}
