CREATE TABLE IF NOT EXISTS pagetranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  pageid character varying(36),
  title character varying(250),
  content text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY pagetranslations
    ADD CONSTRAINT pagetranslations_languagecode_pageid_key UNIQUE (languagecode, pageid);
ALTER TABLE ONLY pagetranslations
    ADD CONSTRAINT fk_pagetranslations_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;
