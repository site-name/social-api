CREATE TABLE IF NOT EXISTS menu_item_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(10),
  menuitemid character varying(36),
  name character varying(128)
);

ALTER TABLE ONLY menu_item_translations
    ADD CONSTRAINT menu_item_translations_languagecode_menuitemid_key UNIQUE (languagecode, menuitemid);

