CREATE TABLE IF NOT EXISTS collection_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  collection_id varchar(36) NOT NULL,
  name varchar(250) NOT NULL,
  description text NOT NULL,
  seo_title varchar(70),
  seo_description varchar(300)
);

ALTER TABLE ONLY collection_translations
    ADD CONSTRAINT collection_translations_language_code_collection_id_key UNIQUE (language_code, collection_id);