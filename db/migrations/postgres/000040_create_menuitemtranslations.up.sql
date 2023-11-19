CREATE TABLE IF NOT EXISTS menu_item_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(10) NOT NULL,
  menu_item_id uuid NOT NULL,
  name character varying(128) NOT NULL
);

ALTER TABLE ONLY menu_item_translations
    ADD CONSTRAINT menu_item_translations_language_code_menu_item_id_key UNIQUE (language_code, menu_item_id);