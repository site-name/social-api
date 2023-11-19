CREATE TABLE IF NOT EXISTS page_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(5) NOT NULL,
  page_id uuid NOT NULL,
  title character varying(250) NOT NULL,
  content text,
  seo_title character varying(70),
  seo_description character varying(300)
);

ALTER TABLE ONLY page_translations
    ADD CONSTRAINT page_translations_language_code_page_id_key UNIQUE (language_code, page_id);