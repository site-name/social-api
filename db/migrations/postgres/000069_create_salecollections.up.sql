CREATE TABLE IF NOT EXISTS sale_collections (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid NOT NULL,
  collection_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT sale_collections_sale_id_collection_id_key UNIQUE (sale_id, collection_id);