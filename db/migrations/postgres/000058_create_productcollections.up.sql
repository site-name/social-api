CREATE TABLE IF NOT EXISTS product_collections (
  id varchar(36) NOT NULL PRIMARY KEY,
  collection_id varchar(36) NOT NULL,
  product_id varchar(36) NOT NULL,
  sort_order integer NOT NULL
);

ALTER TABLE ONLY product_collections
    ADD CONSTRAINT product_collections_collection_id_product_id_key UNIQUE (collection_id, product_id);

CREATE INDEX IF NOT EXISTS product_collections_sort_order_idx ON product_collections (sort_order);