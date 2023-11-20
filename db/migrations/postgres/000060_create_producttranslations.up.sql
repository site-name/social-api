CREATE TABLE IF NOT EXISTS product_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code varchar(10) NOT NULL,
  product_id uuid NOT NULL,
  name varchar(250) NOT NULL,
  description text NOT NULL,
  seo_title varchar(70),
  seo_description varchar(300)
);

ALTER TABLE ONLY product_translations
    ADD CONSTRAINT product_translations_language_code_product_id_key UNIQUE (language_code, product_id);

ALTER TABLE ONLY product_translations
    ADD CONSTRAINT product_translations_name_key UNIQUE (name);