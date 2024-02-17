CREATE TABLE IF NOT EXISTS sale_collections (
  id varchar(36) NOT NULL PRIMARY KEY,
  sale_id varchar(36) NOT NULL,
  collection_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT sale_collections_sale_id_collection_id_key UNIQUE (sale_id, collection_id);