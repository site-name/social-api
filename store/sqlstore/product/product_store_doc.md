```SQL
---------------published--------------
SELECT
  "*"
FROM
  "product_product"
WHERE
  EXISTS(
    SELECT
      (1) AS "a"
    FROM
      "product_productchannellisting" V0
    WHERE
      (
        (
          V0."publication_date" <= '2021-08-19' :: date
          OR V0."publication_date" IS NULL
        )
        AND EXISTS(
          SELECT
            (1) AS "a"
          FROM
            "channel_channel" U0
          WHERE
            (
              U0."is_active"
              AND U0."slug" = 'hello-world'
              AND U0."id" = V0."channel_id"
            )
          LIMIT
            1
        )
        AND V0."is_published"
        AND V0."product_id" = "product_product"."id"
      )
    LIMIT
      1
  )
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

-----------END published---------------
-----------not_published---------------
SELECT
  "*",
  (
    SELECT
      U0."is_published"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "is_published",
  (
    SELECT
      U0."publication_date"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "publication_date"
FROM
  "product_product"
WHERE
  (
    (
      (
        SELECT
          U0."publication_date"
        FROM
          "product_productchannellisting" U0
          INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
        WHERE
          (
            U1."slug" = 'hello-world'
            AND U0."product_id" = "product_product"."id"
          )
        ORDER BY
          U0."id" ASC
        LIMIT
          1
      ) > '2021-08-19' :: date
      AND (
        SELECT
          U0."is_published"
        FROM
          "product_productchannellisting" U0
          INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
        WHERE
          (
            U1."slug" = 'hello-world'
            AND U0."product_id" = "product_product"."id"
          )
        ORDER BY
          U0."id" ASC
        LIMIT
          1
      )
    )
    OR NOT (
      SELECT
        U0."is_published"
      FROM
        "product_productchannellisting" U0
        INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
      WHERE
        (
          U1."slug" = 'hello-world'
          AND U0."product_id" = "product_product"."id"
        )
      ORDER BY
        U0."id" ASC
      LIMIT
        1
    )
    OR (
      SELECT
        U0."is_published"
      FROM
        "product_productchannellisting" U0
        INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
      WHERE
        (
          U1."slug" = 'hello-world'
          AND U0."product_id" = "product_product"."id"
        )
      ORDER BY
        U0."id" ASC
      LIMIT
        1
    ) IS NULL
  )
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

-----------END not_published---------------
--------published_with_variants------------
SELECT
  "*"
FROM
  "product_product"
WHERE
  (
    EXISTS(
      SELECT
        (1) AS "a"
      FROM
        "product_productchannellisting" V0
      WHERE
        (
          (
            V0."publication_date" <= '2021-08-19' :: date
            OR V0."publication_date" IS NULL
          )
          AND EXISTS(
            SELECT
              (1) AS "a"
            FROM
              "channel_channel" U0
            WHERE
              (
                U0."is_active"
                AND U0."slug" = 'hello-world'
                AND U0."id" = V0."channel_id"
              )
            LIMIT
              1
          )
          AND V0."is_published"
          AND V0."product_id" = "product_product"."id"
        )
      LIMIT
        1
    )
    AND EXISTS(
      SELECT
        (1) AS "a"
      FROM
        "product_productvariant" W0
      WHERE
        (
          EXISTS(
            SELECT
              (1) AS "a"
            FROM
              "product_productvariantchannellisting" V0
            WHERE
              (
                EXISTS(
                  SELECT
                    (1) AS "a"
                  FROM
                    "channel_channel" U0
                  WHERE
                    (
                      U0."is_active"
                      AND U0."slug" = 'hello-world'
                      AND U0."id" = V0."channel_id"
                    )
                  LIMIT
                    1
                )
                AND V0."price_amount" IS NOT NULL
                AND V0."variant_id" = W0."id"
              )
            LIMIT
              1
          )
          AND W0."product_id" = "product_product"."id"
        )
      LIMIT
        1
    )
  )
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

--------END published_with_variants------------
-------------visible_to_user-------------------

--- CASE 1: requestor is shop staff and channel-slug is provided:
SELECT
  "*"
FROM
  "product_product"
WHERE
  EXISTS(
    SELECT
      (1) AS "a"
    FROM
      "product_productchannellisting" V0
    WHERE
      (
        EXISTS(
          SELECT
            (1) AS "a"
          FROM
            "channel_channel" U0
          WHERE
            (
              U0."slug" = 'hello-world'
              AND U0."id" = V0."channel_id"
            )
          LIMIT
            1
        )
        AND V0."product_id" = "product_product"."id"
      )
    LIMIT
      1
  )
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;
----- CASE 2: requestor is shop staff and channel-slug empty
SELECT "*" FROM "product_product";

----- CASE 3: requester is shop visitor, do the same as `published_with_variants (line 157)`
SELECT
  "*"
FROM
  "product_product"
WHERE
  (
    EXISTS(
      SELECT
        (1) AS "a"
      FROM
        "product_productchannellisting" V0
      WHERE
        (
          (
            V0."publication_date" <= '2021-08-19' :: date
            OR V0."publication_date" IS NULL
          )
          AND EXISTS(
            SELECT
              (1) AS "a"
            FROM
              "channel_channel" U0
            WHERE
              (
                U0."is_active"
                AND U0."slug" = 'hello-world'
                AND U0."id" = V0."channel_id"
              )
            LIMIT
              1
          )
          AND V0."is_published"
          AND V0."product_id" = "product_product"."id"
        )
      LIMIT
        1
    )
    AND EXISTS(
      SELECT
        (1) AS "a"
      FROM
        "product_productvariant" W0
      WHERE
        (
          EXISTS(
            SELECT
              (1) AS "a"
            FROM
              "product_productvariantchannellisting" V0
            WHERE
              (
                EXISTS(
                  SELECT
                    (1) AS "a"
                  FROM
                    "channel_channel" U0
                  WHERE
                    (
                      U0."is_active"
                      AND U0."slug" = 'hello-world'
                      AND U0."id" = V0."channel_id"
                    )
                  LIMIT
                    1
                )
                AND V0."price_amount" IS NOT NULL
                AND V0."variant_id" = W0."id"
              )
            LIMIT
              1
          )
          AND W0."product_id" = "product_product"."id"
        )
      LIMIT
        1
    )
  )
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

-------------END visible_to_user-------------------
----------annotate_publication_info----------------
SELECT
  "*",
  (
    SELECT
      U0."is_published"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "is_published",
  (
    SELECT
      U0."publication_date"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "publication_date"
FROM
  "product_product"
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

----------END annotate_publication_info----------
---------annotate_is_published-------------------
SELECT
  "*",
  (
    SELECT
      U0."is_published"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "is_published"
FROM
  "product_product"
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

---------END annotate_is_published--------------
---------annotate_publication_date--------------
SELECT
  "*",
  (
    SELECT
      U0."publication_date"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "publication_date"
FROM
  "product_product"
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

---------END annotate_publication_date--------------
-----------annotate_visible_in_listings-------------
SELECT
  "*",
  (
    SELECT
      U0."visible_in_listings"
    FROM
      "product_productchannellisting" U0
      INNER JOIN "channel_channel" U1 ON (U0."channel_id" = U1."id")
    WHERE
      (
        U1."slug" = 'hello-world'
        AND U0."product_id" = "product_product"."id"
      )
    ORDER BY
      U0."id" ASC
    LIMIT
      1
  ) AS "visible_in_listings"
FROM
  "product_product"
ORDER BY
  "product_product"."slug" ASC
LIMIT
  21;

-----------END annotate_visible_in_listings-------------
------------------sort_by_attribute---------------------
SELECT
  "attribute_attributeproduct"."id",
  "attribute_attributeproduct"."product_type_id"
FROM
  "attribute_attributeproduct"
WHERE
  "attribute_attributeproduct"."attribute_id" = 1
ORDER BY
  "attribute_attributeproduct"."sort_order" ASC,
  "attribute_attributeproduct"."id" ASC;

SELECT
  "*",
  NULL AS "concatenated_values_order",
  NULL AS "concatenated_values"
FROM
  "product_product"
ORDER BY
  "concatenated_values_order" DESC,
  "concatenated_values" DESC,
  "product_product"."name" DESC
LIMIT
  21;

------------END------------------------------------
```