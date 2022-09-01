CREATE TABLE IF NOT EXISTS salecollections (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  collectionid character varying(36),
  createat bigint
);

ALTER TABLE ONLY salecollections
    ADD CONSTRAINT salecollections_saleid_collectionid_key UNIQUE (saleid, collectionid);

