CREATE TABLE IF NOT EXISTS product_variant_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  product_variant_id varchar(36) NOT NULL,
  name varchar(255) NOT NULL
);

ALTER TABLE ONLY product_variant_translations
    ADD CONSTRAINT product_variant_translations_language_code_product_variant_id_key UNIQUE (language_code, product_variant_id);