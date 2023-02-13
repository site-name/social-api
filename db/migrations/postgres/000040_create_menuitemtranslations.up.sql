CREATE TABLE IF NOT EXISTS menuitemtranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(10),
  menuitemid character varying(36),
  name character varying(128)
);

ALTER TABLE ONLY menuitemtranslations
    ADD CONSTRAINT menuitemtranslations_languagecode_menuitemid_key UNIQUE (languagecode, menuitemid);

ALTER TABLE ONLY menuitemtranslations
    ADD CONSTRAINT fk_menuitemtranslations_menuitems FOREIGN KEY (menuitemid) REFERENCES menuitems(id) ON DELETE CASCADE;
