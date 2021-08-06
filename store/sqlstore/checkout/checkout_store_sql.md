```sql
--
-- all 7 quesries for this call:
-- lines = checkout.lines.prefetch_related(
--         "variant__product__collections",
--         "variant__channel_listings__channel",
--         "variant__product__product_type",
--     )
--
SELECT
  "checkout_checkoutline"."id",
  "checkout_checkoutline"."checkout_id",
  "checkout_checkoutline"."variant_id",
  "checkout_checkoutline"."quantity"
FROM
  "checkout_checkoutline"
WHERE
  "checkout_checkoutline"."checkout_id" = 'eee90544-abf1-4db6-bc6c-7b5c423d8aa5' :: uuid
ORDER BY
  "checkout_checkoutline"."id" ASC
LIMIT
  21;

-- args =(UUID('eee90544-abf1-4db6-bc6c-7b5c423d8aa5'),) [PID:93719:MainThread] DEBUG django.db.backends (0.000)
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
  "product_productvariant"."id" IN (314)
ORDER BY
  "product_productvariant"."sort_order" ASC,
  "product_productvariant"."sku" ASC;

-- args =(314,) [PID:93719:MainThread] DEBUG django.db.backends (0.000)
SELECT
  "product_product"."id",
  "product_product"."private_metadata",
  "product_product"."metadata",
  "product_product"."seo_title",
  "product_product"."seo_description",
  "product_product"."product_type_id",
  "product_product"."name",
  "product_product"."slug",
  "product_product"."description",
  "product_product"."description_plaintext",
  "product_product"."search_vector",
  "product_product"."category_id",
  "product_product"."updated_at",
  "product_product"."charge_taxes",
  "product_product"."weight",
  "product_product"."default_variant_id",
  "product_product"."rating"
FROM
  "product_product"
WHERE
  "product_product"."id" IN (118)
ORDER BY
  "product_product"."slug" ASC;

-- args =(118,) [PID:93719:MainThread] DEBUG django.db.backends (0.001)
SELECT
  ("product_collectionproduct"."product_id") AS "_prefetch_related_val_product_id",
  "product_collection"."id",
  "product_collection"."private_metadata",
  "product_collection"."metadata",
  "product_collection"."seo_title",
  "product_collection"."seo_description",
  "product_collection"."name",
  "product_collection"."slug",
  "product_collection"."background_image",
  "product_collection"."background_image_alt",
  "product_collection"."description"
FROM
  "product_collection"
  INNER JOIN "product_collectionproduct" ON (
    "product_collection"."id" = "product_collectionproduct"."collection_id"
  )
WHERE
  "product_collectionproduct"."product_id" IN (118)
ORDER BY
  "product_collection"."slug" ASC;

-- args =(118,) [PID:93719:MainThread] DEBUG django.db.backends (0.002)
SELECT
  "product_productvariantchannellisting"."id",
  "product_productvariantchannellisting"."variant_id",
  "product_productvariantchannellisting"."channel_id",
  "product_productvariantchannellisting"."currency",
  "product_productvariantchannellisting"."price_amount",
  "product_productvariantchannellisting"."cost_price_amount"
FROM
  "product_productvariantchannellisting"
WHERE
  "product_productvariantchannellisting"."variant_id" IN (314)
ORDER BY
  "product_productvariantchannellisting"."id" ASC;

-- args =(314,) [PID:93719:MainThread] DEBUG django.db.backends (0.000)
SELECT
  "channel_channel"."id",
  "channel_channel"."name",
  "channel_channel"."is_active",
  "channel_channel"."slug",
  "channel_channel"."currency_code"
FROM
  "channel_channel"
WHERE
  "channel_channel"."id" IN (1)
ORDER BY
  "channel_channel"."slug" ASC;

-- args =(1,) [PID:93719:MainThread] DEBUG django.db.backends (0.002)
SELECT
  "product_producttype"."id",
  "product_producttype"."private_metadata",
  "product_producttype"."metadata",
  "product_producttype"."name",
  "product_producttype"."slug",
  "product_producttype"."has_variants",
  "product_producttype"."is_shipping_required",
  "product_producttype"."is_digital",
  "product_producttype"."weight"
FROM
  "product_producttype"
WHERE
  "product_producttype"."id" IN (14)
ORDER BY
  "product_producttype"."slug" ASC;

-- args =(14,)
```