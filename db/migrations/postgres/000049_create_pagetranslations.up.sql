CREATE TABLE IF NOT EXISTS page_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  pageid character varying(36),
  title character varying(250),
  content text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY page_translations
    ADD CONSTRAINT page_translations_languagecode_pageid_key UNIQUE (languagecode, pageid);
