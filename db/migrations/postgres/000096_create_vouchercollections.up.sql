CREATE TABLE IF NOT EXISTS vouchercollections (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  collectionid character varying(36)
);

ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT vouchercollections_voucherid_collectionid_key UNIQUE (voucherid, collectionid);
