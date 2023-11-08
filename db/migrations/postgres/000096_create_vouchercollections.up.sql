CREATE TABLE IF NOT EXISTS voucher_collections (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  collectionid character varying(36)
);

ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT voucher_collections_voucherid_collectionid_key UNIQUE (voucherid, collectionid);
