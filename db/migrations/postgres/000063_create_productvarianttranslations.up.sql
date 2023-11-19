CREATE TABLE IF NOT EXISTS product_variant_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5) NOT NULL,
  product_variant_id uuid NOT NULL,
  name character varying(255) NOT NULL
);

ALTER TABLE ONLY product_variant_translations
    ADD CONSTRAINT product_variant_translations_language_code_product_variant_id_key UNIQUE (language_code, product_variant_id);