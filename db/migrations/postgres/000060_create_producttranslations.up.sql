CREATE TABLE IF NOT EXISTS product_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5),
  product_id uuid,
  name character varying(250),
  description text,
  seo_title character varying(70),
  seo_description character varying(300)
);

ALTER TABLE ONLY product_translations
    ADD CONSTRAINT product_translations_language_code_product_id_key UNIQUE (language_code, product_id);

ALTER TABLE ONLY product_translations
    ADD CONSTRAINT product_translations_name_key UNIQUE (name);