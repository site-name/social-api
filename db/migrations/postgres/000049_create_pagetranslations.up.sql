CREATE TABLE IF NOT EXISTS page_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  page_id varchar(36) NOT NULL,
  title varchar(250) NOT NULL,
  content text,
  seo_title varchar(70),
  seo_description varchar(300)
);

ALTER TABLE ONLY page_translations
    ADD CONSTRAINT page_translations_language_code_page_id_key UNIQUE (language_code, page_id);