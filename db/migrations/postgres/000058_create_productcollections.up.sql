CREATE TABLE IF NOT EXISTS product_collections (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  collection_id uuid NOT NULL,
  product_id uuid NOT NULL
);

ALTER TABLE ONLY product_collections
    ADD CONSTRAINT product_collections_collection_id_product_id_key UNIQUE (collection_id, product_id);