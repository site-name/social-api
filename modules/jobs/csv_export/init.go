package csv_export

type Fields struct {
	HEADERS_TO_FIELDS_MAPPING      map[string]map[string]string
	PRODUCT_ATTRIBUTE_FIELDS       map[string]string
	PRODUCT_CHANNEL_LISTING_FIELDS map[string]string
	WAREHOUSE_FIELDS               map[string]string
	VARIANT_ATTRIBUTE_FIELDS       map[string]string
	VARIANT_CHANNEL_LISTING_FIELDS map[string]string
}

var (
	ProductExportFields *Fields // Data structure with fields for product export
)

func init() {
	ProductExportFields = &Fields{
		HEADERS_TO_FIELDS_MAPPING: map[string]map[string]string{
			"fields": {
				"id":             "id",
				"name":           "name",
				"description":    "description_as_str",
				"category":       "category__slug",
				"product_type":   "product_type__name",
				"charge_taxes":   "charge_taxes",
				"product_weight": "product_weight",
				"variant_sku":    "variants__sku",
				"variant_weight": "variant_weight",
			},
			"product_many_to_many": {
				"collections":   "collections__slug",
				"product_media": "media__image",
			},
			"variant_many_to_many": {
				"variant_media": "variants__media__image",
			},
		},
		PRODUCT_ATTRIBUTE_FIELDS: map[string]string{
			"value":        "attributes__values__slug",
			"file_url":     "attributes__values__file_url",
			"rich_text":    "attributes__values__rich_text",
			"slug":         "attributes__assignment__attribute__slug",
			"input_type":   "attributes__assignment__attribute__input_type",
			"entity_type":  "attributes__assignment__attribute__entity_type",
			"unit":         "attributes__assignment__attribute__unit",
			"attribute_pk": "attributes__assignment__attribute__pk",
		},
		PRODUCT_CHANNEL_LISTING_FIELDS: map[string]string{
			"channel_pk":             "channel_listings__channel__pk",
			"slug":                   "channel_listings__channel__slug",
			"product_currency_code":  "channel_listings__currency",
			"published":              "channel_listings__is_published",
			"publication_date":       "channel_listings__publication_date",
			"searchable":             "channel_listings__visible_in_listings",
			"available for purchase": "channel_listings__available_for_purchase",
		},
		WAREHOUSE_FIELDS: map[string]string{
			"slug":         "variants__stocks__warehouse__slug",
			"quantity":     "variants__stocks__quantity",
			"warehouse_pk": "variants__stocks__warehouse__id",
		},
		VARIANT_ATTRIBUTE_FIELDS: map[string]string{
			"value":        "variants__attributes__values__slug",
			"file_url":     "variants__attributes__values__file_url",
			"rich_text":    "variants__attributes__values__rich_text",
			"slug":         "variants__attributes__assignment__attribute__slug",
			"input_type":   "variants__attributes__assignment__attribute__input_type",
			"entity_type":  "variants__attributes__assignment__attribute__entity_type",
			"unit":         "variants__attributes__assignment__attribute__unit",
			"attribute_pk": "variants__attributes__assignment__attribute__pk",
		},
		VARIANT_CHANNEL_LISTING_FIELDS: map[string]string{
			"channel_pk":            "variants__channel_listings__channel__pk",
			"slug":                  "variants__channel_listings__channel__slug",
			"price_amount":          "variants__channel_listings__price_amount",
			"variant_currency_code": "variants__channel_listings__currency",
			"variant_cost_price":    "variants__channel_listings__cost_price_amount",
		},
	}
}
