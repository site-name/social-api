CREATE TABLE IF NOT EXISTS collection_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  collectionid character varying(36),
  name character varying(250),
  description text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY collection_translations
    ADD CONSTRAINT collection_translations_languagecode_collectionid_key UNIQUE (languagecode, collectionid);

