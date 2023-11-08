CREATE TABLE IF NOT EXISTS product_collections (
  id character varying(36) NOT NULL PRIMARY KEY,
  collectionid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY product_collections
    ADD CONSTRAINT product_collections_collectionid_productid_key UNIQUE (collectionid, productid);
