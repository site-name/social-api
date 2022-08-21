CREATE TABLE IF NOT EXISTS collectiontranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  collectionid character varying(36),
  name character varying(250),
  description text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY public.collectiontranslations
    ADD CONSTRAINT collectiontranslations_languagecode_collectionid_key UNIQUE (languagecode, collectionid);

ALTER TABLE ONLY public.collectiontranslations
    ADD CONSTRAINT fk_collectiontranslations_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;

