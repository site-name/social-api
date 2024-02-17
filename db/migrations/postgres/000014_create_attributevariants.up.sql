-- CREATE TABLE IF NOT EXISTS attribute_variants (
--   id varchar(36) NOT NULL PRIMARY KEY,
--   attribute_id varchar(36) NOT NULL,
--   product_type_id varchar(36) NOT NULL,
--   variant_selection boolean NOT NULL DEFAULT false,
--   sort_order integer
-- );

-- ALTER TABLE ONLY attribute_variants
--     ADD CONSTRAINT attribute_variants_attribute_id_product_type_id_key UNIQUE (attribute_id, product_type_id);
