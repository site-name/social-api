```sql
-- Queries created by ApplicableShippingMethods()
-----------------------
--- query with excluded_product_ids argument
-----------------------
SELECT
  "shipping_shippingmethod"."id",
  "shipping_shippingmethod"."private_metadata",
  "shipping_shippingmethod"."metadata",
  "shipping_shippingmethod"."name",
  "shipping_shippingmethod"."type",
  "shipping_shippingmethod"."shipping_zone_id",
  "shipping_shippingmethod"."minimum_order_weight",
  "shipping_shippingmethod"."maximum_order_weight",
  "shipping_shippingmethod"."maximum_delivery_days",
  "shipping_shippingmethod"."minimum_delivery_days",
  "shipping_shippingmethod"."description",
  (
    SELECT
      U0."price_amount"
    FROM
      "shipping_shippingmethodchannellisting" U0
    WHERE
      (
        U0."channel_id" = % s
        AND U0."shipping_method_id" = "shipping_shippingmethod"."id"
      )
  ) AS "price_amount"
FROM
  "shipping_shippingmethod"
  INNER JOIN "shipping_shippingmethodchannellisting" ON (
    "shipping_shippingmethod"."id" = "shipping_shippingmethodchannellisting"."shipping_method_id"
  )
  INNER JOIN "shipping_shippingzone" ON (
    "shipping_shippingmethod"."shipping_zone_id" = "shipping_shippingzone"."id"
  )
  INNER JOIN "shipping_shippingzone_channels" ON (
    "shipping_shippingzone"."id" = "shipping_shippingzone_channels"."shippingzone_id"
  )
WHERE
  (
    (
      "shipping_shippingmethodchannellisting"."channel_id" = % s
      AND "shipping_shippingmethodchannellisting"."currency" = % s
      AND "shipping_shippingzone_channels"."channel_id" = % s
      AND "shipping_shippingzone"."countries" :: text LIKE % s
      AND NOT (
        EXISTS(
          SELECT
            (1) AS "a"
          FROM
            "shipping_shippingmethod_excluded_products" U1
          WHERE
            (
              U1."product_id" IN (% s, % s, % s)
              AND U1."shippingmethod_id" = "shipping_shippingmethod"."id"
            )
          LIMIT
            1
        )
      )
      AND "shipping_shippingmethod"."type" = % s
      AND "shipping_shippingmethod"."id" IN (
        SELECT
          W0."shipping_method_id"
        FROM
          "shipping_shippingmethodchannellisting" W0
        WHERE
          (
            W0."channel_id" = % s
            AND W0."shipping_method_id" IN (
              SELECT
                V0."id"
              FROM
                "shipping_shippingmethod" V0
                INNER JOIN "shipping_shippingmethodchannellisting" V1 ON (V0."id" = V1."shipping_method_id")
                INNER JOIN "shipping_shippingzone" V3 ON (V0."shipping_zone_id" = V3."id")
                INNER JOIN "shipping_shippingzone_channels" V4 ON (V3."id" = V4."shippingzone_id")
              WHERE
                (
                  V1."channel_id" = % s
                  AND V1."currency" = % s
                  AND V4."channel_id" = % s
                  AND V3."countries" :: text LIKE % s
                  AND NOT (
                    EXISTS(
                      SELECT
                        (1) AS "a"
                      FROM
                        "shipping_shippingmethod_excluded_products" U1
                      WHERE
                        (
                          U1."product_id" IN (% s, % s, % s)
                          AND U1."shippingmethod_id" = V0."id"
                        )
                      LIMIT
                        1
                    )
                  )
                  AND V0."type" = % s
                )
            )
            AND W0."minimum_order_price_amount" <= % s
            AND (
              W0."maximum_order_price_amount" IS NULL
              OR W0."maximum_order_price_amount" >= % s
            )
          )
      )
    )
    OR (
      "shipping_shippingmethodchannellisting"."channel_id" = % s
      AND "shipping_shippingmethodchannellisting"."currency" = % s
      AND "shipping_shippingzone_channels"."channel_id" = % s
      AND "shipping_shippingzone"."countries" :: text LIKE % s
      AND NOT (
        EXISTS(
          SELECT
            (1) AS "a"
          FROM
            "shipping_shippingmethod_excluded_products" U1
          WHERE
            (
              U1."product_id" IN (% s, % s, % s)
              AND U1."shippingmethod_id" = "shipping_shippingmethod"."id"
            )
          LIMIT
            1
        )
      )
      AND "shipping_shippingmethod"."type" = % s
      AND (
        "shipping_shippingmethod"."minimum_order_weight" <= % s
        OR "shipping_shippingmethod"."minimum_order_weight" IS NULL
      )
      AND (
        "shipping_shippingmethod"."maximum_order_weight" >= % s
        OR "shipping_shippingmethod"."maximum_order_weight" IS NULL
      )
    )
  )
ORDER BY
  "price_amount" ASC 


-----------------------
--- query without excluded_product_ids passed
-----------------------
SELECT
  "shipping_shippingmethod"."id",
  "shipping_shippingmethod"."private_metadata",
  "shipping_shippingmethod"."metadata",
  "shipping_shippingmethod"."name",
  "shipping_shippingmethod"."type",
  "shipping_shippingmethod"."shipping_zone_id",
  "shipping_shippingmethod"."minimum_order_weight",
  "shipping_shippingmethod"."maximum_order_weight",
  "shipping_shippingmethod"."maximum_delivery_days",
  "shipping_shippingmethod"."minimum_delivery_days",
  "shipping_shippingmethod"."description",
  (
    SELECT
      U0."price_amount"
    FROM
      "shipping_shippingmethodchannellisting" U0
    WHERE
      (
        U0."channel_id" = % s
        AND U0."shipping_method_id" = "shipping_shippingmethod"."id"
      )
  ) AS "price_amount"
FROM
  "shipping_shippingmethod"
  INNER JOIN "shipping_shippingmethodchannellisting" ON (
    "shipping_shippingmethod"."id" = "shipping_shippingmethodchannellisting"."shipping_method_id"
  )
  INNER JOIN "shipping_shippingzone" ON (
    "shipping_shippingmethod"."shipping_zone_id" = "shipping_shippingzone"."id"
  )
  INNER JOIN "shipping_shippingzone_channels" ON (
    "shipping_shippingzone"."id" = "shipping_shippingzone_channels"."shippingzone_id"
  )
WHERE
  (
    (
      "shipping_shippingmethodchannellisting"."channel_id" = % s
      AND "shipping_shippingmethodchannellisting"."currency" = % s
      AND "shipping_shippingzone_channels"."channel_id" = % s
      AND "shipping_shippingzone"."countries" :: text LIKE % s
      AND "shipping_shippingmethod"."type" = % s
      AND "shipping_shippingmethod"."id" IN (
        SELECT
          W0."shipping_method_id"
        FROM
          "shipping_shippingmethodchannellisting" W0
        WHERE
          (
            W0."channel_id" = % s
            AND W0."shipping_method_id" IN (
              SELECT
                V0."id"
              FROM
                "shipping_shippingmethod" V0
                INNER JOIN "shipping_shippingmethodchannellisting" V1 ON (V0."id" = V1."shipping_method_id")
                INNER JOIN "shipping_shippingzone" V3 ON (V0."shipping_zone_id" = V3."id")
                INNER JOIN "shipping_shippingzone_channels" V4 ON (V3."id" = V4."shippingzone_id")
              WHERE
                (
                  V1."channel_id" = % s
                  AND V1."currency" = % s
                  AND V4."channel_id" = % s
                  AND V3."countries" :: text LIKE % s
                  AND V0."type" = % s
                )
            )
            AND W0."minimum_order_price_amount" <= % s
            AND (
              W0."maximum_order_price_amount" IS NULL
              OR W0."maximum_order_price_amount" >= % s
            )
          )
      )
    )
    OR (
      "shipping_shippingmethodchannellisting"."channel_id" = % s
      AND "shipping_shippingmethodchannellisting"."currency" = % s
      AND "shipping_shippingzone_channels"."channel_id" = % s
      AND "shipping_shippingzone"."countries" :: text LIKE % s
      AND "shipping_shippingmethod"."type" = % s
      AND (
        "shipping_shippingmethod"."minimum_order_weight" <= % s
        OR "shipping_shippingmethod"."minimum_order_weight" IS NULL
      )
      AND (
        "shipping_shippingmethod"."maximum_order_weight" >= % s
        OR "shipping_shippingmethod"."maximum_order_weight" IS NULL
      )
    )
  )
ORDER BY
  "price_amount" ASC

```