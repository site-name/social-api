CREATE TABLE IF NOT EXISTS vouchercollections (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  collectionid character varying(36)
);

ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT vouchercollections_voucherid_collectionid_key UNIQUE (voucherid, collectionid);

ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT fk_vouchercollections_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;

ALTER TABLE ONLY vouchercollections
    ADD CONSTRAINT fk_vouchercollections_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;
