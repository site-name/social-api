```sql

-- 1) perform the following query first.
SELECT
  "order_orderline"."id",
  "order_orderline"."order_id",
  "order_orderline"."variant_id",
  "order_orderline"."product_name",
  "order_orderline"."variant_name",
  "order_orderline"."translated_product_name",
  "order_orderline"."translated_variant_name",
  "order_orderline"."product_sku",
  "order_orderline"."is_shipping_required",
  "order_orderline"."quantity",
  "order_orderline"."quantity_fulfilled",
  "order_orderline"."currency",
  "order_orderline"."unit_discount_amount",
  "order_orderline"."unit_discount_type",
  "order_orderline"."unit_discount_reason",
  "order_orderline"."unit_price_net_amount",
  "order_orderline"."unit_discount_value",
  "order_orderline"."unit_price_gross_amount",
  "order_orderline"."total_price_net_amount",
  "order_orderline"."total_price_gross_amount",
  "order_orderline"."undiscounted_unit_price_gross_amount",
  "order_orderline"."undiscounted_unit_price_net_amount",
  "order_orderline"."undiscounted_total_price_gross_amount",
  "order_orderline"."undiscounted_total_price_net_amount",
  "order_orderline"."tax_rate"
FROM
  "order_orderline"
  INNER JOIN "product_productvariant" ON (
    "order_orderline"."variant_id" = "product_productvariant"."id"
  )
  INNER JOIN "product_digitalcontent" ON (
    "product_productvariant"."id" = "product_digitalcontent"."product_variant_id"
  )
WHERE
  (
    "order_orderline"."order_id" = 14
    AND NOT "order_orderline"."is_shipping_required"
    AND "product_digitalcontent"."id" IS NOT NULL
  )
ORDER BY
  "order_orderline"."id" ASC;

-- If the above query found order lines, perform the following 2 queries

SELECT
  "product_productvariant"."id",
  "product_productvariant"."sort_order",
  "product_productvariant"."private_metadata",
  "product_productvariant"."metadata",
  "product_productvariant"."sku",
  "product_productvariant"."name",
  "product_productvariant"."product_id",
  "product_productvariant"."track_inventory",
  "product_productvariant"."weight"
FROM
  "product_productvariant"
WHERE
  "product_productvariant"."id" IN (196, 197)
ORDER BY
  "product_productvariant"."sort_order" ASC,
  "product_productvariant"."sku" ASC;

-- 
SELECT
  "product_digitalcontent"."id",
  "product_digitalcontent"."private_metadata",
  "product_digitalcontent"."metadata",
  "product_digitalcontent"."use_default_settings",
  "product_digitalcontent"."automatic_fulfillment",
  "product_digitalcontent"."content_type",
  "product_digitalcontent"."product_variant_id",
  "product_digitalcontent"."content_file",
  "product_digitalcontent"."max_downloads",
  "product_digitalcontent"."url_valid_days"
FROM
  "product_digitalcontent"
WHERE
  "product_digitalcontent"."product_variant_id" IN (197, 196);

```
