```sql
--- sql for FilterVariantStocksForCountry
SELECT
  "warehouse_stock"."id",
  "warehouse_stock"."warehouse_id",
  "warehouse_stock"."product_variant_id",
  "warehouse_stock"."quantity",
  "warehouse_warehouse"."private_metadata",
  "warehouse_warehouse"."metadata",
  "warehouse_warehouse"."id",
  "warehouse_warehouse"."name",
  "warehouse_warehouse"."slug",
  "warehouse_warehouse"."address_id",
  "warehouse_warehouse"."email",
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
  "warehouse_stock"
  INNER JOIN "warehouse_warehouse" ON (
    "warehouse_stock"."warehouse_id" = "warehouse_warehouse"."id"
  )
  INNER JOIN "product_productvariant" ON (
    "warehouse_stock"."product_variant_id" = "product_productvariant"."id"
  )
WHERE
  (
    "warehouse_stock"."warehouse_id" IN (
      SELECT
        U0."id"
      FROM
        "warehouse_warehouse" U0
        INNER JOIN "warehouse_warehouse_shipping_zones" U1 ON (U0."id" = U1."warehouse_id")
        INNER JOIN "shipping_shippingzone" U2 ON (U1."shippingzone_id" = U2."id")
        INNER JOIN "shipping_shippingzone_channels" U3 ON (U2."id" = U3."shippingzone_id")
        INNER JOIN "channel_channel" U4 ON (U3."channel_id" = U4."id")
      WHERE
        (
          U4."slug" = % s
          AND U2."countries" :: text LIKE % s
        )
    )
    AND "warehouse_stock"."product_variant_id" = % s
  )
ORDER BY
  "warehouse_stock"."id" ASC 
```

```sql
---sql for FilterProductStocksForCountryAndChannel

SELECT
  "warehouse_stock"."id",
  "warehouse_stock"."warehouse_id",
  "warehouse_stock"."product_variant_id",
  "warehouse_stock"."quantity",
  "warehouse_warehouse"."private_metadata",
  "warehouse_warehouse"."metadata",
  "warehouse_warehouse"."id",
  "warehouse_warehouse"."name",
  "warehouse_warehouse"."slug",
  "warehouse_warehouse"."address_id",
  "warehouse_warehouse"."email",
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
  "warehouse_stock"
  INNER JOIN "warehouse_warehouse" ON (
    "warehouse_stock"."warehouse_id" = "warehouse_warehouse"."id"
  )
  INNER JOIN "product_productvariant" ON (
    "warehouse_stock"."product_variant_id" = "product_productvariant"."id"
  )
WHERE
  (
    "warehouse_stock"."warehouse_id" IN (
      SELECT
        U0."id"
      FROM
        "warehouse_warehouse" U0
        INNER JOIN "warehouse_warehouse_shipping_zones" U1 ON (U0."id" = U1."warehouse_id")
        INNER JOIN "shipping_shippingzone" U2 ON (U1."shippingzone_id" = U2."id")
        INNER JOIN "shipping_shippingzone_channels" U3 ON (U2."id" = U3."shippingzone_id")
        INNER JOIN "channel_channel" U4 ON (U3."channel_id" = U4."id")
      WHERE
        (
          U4."slug" = % s
          AND U2."countries" :: text LIKE % s
        )
    )
    AND "product_productvariant"."product_id" = % s
  )
ORDER BY
  "warehouse_stock"."id" ASC
```
