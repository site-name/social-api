CREATE TABLE IF NOT EXISTS voucher_collections (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  collection_id varchar(36) NOT NULL
);

ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT voucher_collections_voucher_id_collection_id_key UNIQUE (voucher_id, collection_id);
    