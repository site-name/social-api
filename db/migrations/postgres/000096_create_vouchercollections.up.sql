CREATE TABLE IF NOT EXISTS voucher_collections (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  voucher_id uuid NOT NULL,
  collection_id uuid NOT NULL
);

ALTER TABLE ONLY voucher_collections
    ADD CONSTRAINT voucher_collections_voucher_id_collection_id_key UNIQUE (voucher_id, collection_id);
    