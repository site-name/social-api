CREATE TABLE IF NOT EXISTS collection_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5),
  collection_id uuid,
  name character varying(250),
  description text,
  seo_title character varying(70),
  seo_description character varying(300)
);

ALTER TABLE ONLY collection_translations
    ADD CONSTRAINT collection_translations_language_code_collection_id_key UNIQUE (language_code, collection_id);