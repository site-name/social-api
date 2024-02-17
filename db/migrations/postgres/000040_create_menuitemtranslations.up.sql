CREATE TABLE IF NOT EXISTS menu_item_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  menu_item_id varchar(36) NOT NULL,
  name varchar(128) NOT NULL
);

ALTER TABLE ONLY menu_item_translations
    ADD CONSTRAINT menu_item_translations_language_code_menu_item_id_key UNIQUE (language_code, menu_item_id);