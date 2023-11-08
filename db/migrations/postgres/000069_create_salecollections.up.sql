CREATE TABLE IF NOT EXISTS sale_collections (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  collectionid character varying(36),
  createat bigint
);

ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT sale_collections_saleid_collectionid_key UNIQUE (saleid, collectionid);

